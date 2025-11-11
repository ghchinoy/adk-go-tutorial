# Tutorial 09: Long-term Memory

In previous tutorials, we've seen how agents maintain context within a single session. However, once a session ends, that context is typically lost. Long-term memory allows agents to retain information across different sessions, enabling them to learn about users over time and provide more personalized experiences.

In this tutorial, you will learn how to:
1.  Use the `memory.Service` to store and retrieve information.
2.  Create a tool that allows the agent to search its long-term memory.
3.  Ingest completed sessions into memory for future recall.

## Prerequisites

- Completed Tutorial 08 (Human-in-the-Loop).
- Basic understanding of ADK tools and sessions.

## 1. Project Setup

Create a new directory for this tutorial and initialize the Go module:

```bash
mkdir -p adk-tutorial-09
cd adk-tutorial-09
go mod init adk-tutorial-09
# Replace with the actual path to your local ADK copy if necessary
go mod edit -replace google.golang.org/adk=../adk-go
go mod tidy
```

## 2. The Memory Service

ADK provides a `memory.Service` interface for managing long-term memory. For this tutorial, we'll use the built-in `memory.InMemoryService`, which stores data in RAM. In a production environment, you would likely use a persistent implementation (e.g., backed by a vector database).

The memory service works by "ingesting" entire sessions. It indexes the events within those sessions so they can be searched later.

## 3. Creating the Recall Tool

While the `runner` can be configured with a memory service, the agent doesn't automatically know *when* or *how* to search it. We need to provide a tool for this.

ADK's `tool.Context` provides a `SearchMemory` method that tools can use to access the memory service associated with the current runner.

Create a `main.go` file and start by defining the tool:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// RecallArgs defines the input for the recall tool.
type RecallArgs struct {
	Query string `json:"query"`
}

// RecallResult defines the output of the recall tool.
type RecallResult struct {
	Memories []string `json:"memories"`
}

// recall is the function that will be wrapped as a tool.
// It uses ctx.SearchMemory to find relevant past interactions.
func recall(ctx tool.Context, args RecallArgs) RecallResult {
	log.Printf("[Tool] Recalling memory for query: '%s'", args.Query)
	// SearchMemory uses the memory service configured in the runner.
	resp, err := ctx.SearchMemory(context.Background(), args.Query)
	if err != nil {
		log.Printf("[Tool] Error searching memory: %v", err)
		return RecallResult{Memories: []string{fmt.Sprintf("Error searching memory: %v", err)}}
	}

	var memories []string
	for _, m := range resp.Memories {
		if m.Content == nil {
			continue
		}
		// Extract text from the memory content.
		text := ""
		for _, p := range m.Content.Parts {
			text += p.Text
		}
		if text != "" {
			memories = append(memories, text)
		}
	}

	if len(memories) == 0 {
		return RecallResult{Memories: []string{"No relevant memories found."}}
	}
	return RecallResult{Memories: memories}
}
```

## 4. Simulating Multiple Sessions

To demonstrate long-term memory, we need to simulate at least two sessions:
1.  **Session 1:** The user provides some information (e.g., "My favorite color is blue").
2.  **Ingestion:** We manually ingest Session 1 into the memory service.
3.  **Session 2:** The user asks for that information, and the agent uses the `recall` tool to find it.

Add the `main` function to `main.go` to orchestrate this flow:

```go
func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT environment variable must be set")
	}
	// Use a specific model version known to work in your region.
	// You might need to set GOOGLE_CLOUD_LOCATION if it's not us-central1.
	model, err := gemini.NewModel(ctx, "gemini-2.0-flash-001", &genai.ClientConfig{
		Project: projectID,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 1. Initialize Services
	memService := memory.InMemoryService()
	sessionService := session.InMemoryService()

	// 2. Define Tools
	recallTool, err := functiontool.New(functiontool.Config{
		Name:        "recall",
		Description: "Recalls information from previous conversations based on a query. Use this when asked about things you might have learned in the past.",
	}, recall)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Define Agent
	// Give the agent clear instructions on when to use memory.
	myAgent, err := llmagent.New(llmagent.Config{
		Name:  "memory_agent",
		Model: model,
		Tools: []tool.Tool{recallTool},
		Instruction: `You have a memory of past conversations.
Use the 'recall' tool to find information from previous sessions if you don't know the answer immediately.
Always check your memory before saying you don't know something about the user.`,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 4. Initialize Runner with Memory Service
	appName := "memory_tutorial"
	r, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          myAgent,
		SessionService: sessionService,
		MemoryService:  memService, // CRITICAL: Pass memory service here
	})
	if err != nil {
		log.Fatal(err)
	}

	userID := "tutorial_user"

	// --- Session 1: Storing Information ---
	fmt.Println("\n--- Session 1 ---")
	session1Resp, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName: appName,
		UserID:  userID,
	})
	if err != nil {
		log.Fatal(err)
	}
	session1ID := session1Resp.Session.ID()

	// Run a single turn where the user shares information.
	runTurn(ctx, r, session1ID, userID, "Hi, my favorite color is blue. Remember that!")

	// --- Ingestion ---
	fmt.Println("\n[System] Ingesting Session 1 into memory...")
	// IMPORTANT: When using InMemoryService for sessions, we must re-fetch
	// the session to ensure we have the latest events before ingestion,
	// as Create() might return a copy that doesn't get updated.
	session1Updated, err := sessionService.Get(ctx, &session.GetRequest{
		AppName:   appName,
		UserID:    userID,
		SessionID: session1ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := memService.AddSession(ctx, session1Updated.Session); err != nil {
		log.Fatal(err)
	}

	// --- Session 2: Recalling Information ---
	fmt.Println("\n--- Session 2 ---")
	// Create a NEW session for the same user.
	session2Resp, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName: appName,
		UserID:  userID,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Ask a question that requires memory from the previous session.
	runTurn(ctx, r, session2Resp.Session.ID(), userID, "What is my favorite color?")
}

// Helper function to run a single turn of conversation.
func runTurn(ctx context.Context, r *runner.Runner, sessionID, userID, prompt string) {
	fmt.Printf("User: %s\n", prompt)
	fmt.Print("Agent: ")
	userMsg := genai.NewContentFromText(prompt, genai.RoleUser)
	// We use a simplified loop here to just print the response.
	for event, err := range r.Run(ctx, userID, sessionID, userMsg, agent.RunConfig{}) {
		if err != nil {
			log.Printf("\nError during turn: %v\n", err)
			return
		}
		if event.LLMResponse.Content != nil {
			for _, part := range event.LLMResponse.Content.Parts {
				fmt.Print(part.Text)
			}
		}
	}
	fmt.Println()
}
