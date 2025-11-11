// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

type RecallArgs struct {
	Query string `json:"query"`
}

type RecallResult struct {
	Memories []string `json:"memories"`
}

func recall(ctx tool.Context, args RecallArgs) RecallResult {
	fmt.Printf("  [Tool] Recalling memory for query: '%s'\n", args.Query)
	resp, err := ctx.SearchMemory(context.Background(), args.Query)
	if err != nil {
		fmt.Printf("  [Tool] Error searching memory: %v\n", err)
		return RecallResult{Memories: []string{fmt.Sprintf("Error searching memory: %v", err)}}
	}

	fmt.Printf("  [Tool] Found %d raw memories.\n", len(resp.Memories))

	var memories []string
	for i, m := range resp.Memories {
		if m.Content == nil {
			continue
		}
		text := ""
		for _, p := range m.Content.Parts {
			text += p.Text
		}
		if text != "" {
			fmt.Printf("  [Tool] Memory %d: %s\n", i, text)
			memories = append(memories, text)
		}
	}

	if len(memories) == 0 {
		fmt.Println("  [Tool] No relevant memories found after filtering.")
		return RecallResult{Memories: []string{"No relevant memories found."}}
	}
	return RecallResult{Memories: memories}
}

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT environment variable must be set")
	}
	location := os.Getenv("GOOGLE_CLOUD_LOCATION")

	// 1. Initialize Services
	cfg := &genai.ClientConfig{
		Backend: genai.BackendVertexAI,
		Project: projectID,
	}
	if location != "" {
		cfg.Location = location
	}

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", cfg)
	if err != nil {
		log.Fatal(err)
	}

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

	// 4. Initialize Runner
	appName := "memory_experiment"
	r, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          myAgent,
		SessionService: sessionService,
		MemoryService:  memService, // Important: Pass memory service to runner
	})
	if err != nil {
		log.Fatal(err)
	}

	userID := "test_user"

	// --- Session 1: Storing Information ---
	fmt.Println("--- Session 1 ---")
	session1Resp, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName: appName,
		UserID:  userID,
	})
	if err != nil {
		log.Fatal(err)
	}
	session1ID := session1Resp.Session.ID()

	runTurn(ctx, r, session1ID, userID, "my favorite color is blue")

	// Manually ingest session 1 into memory.
	// CRITICAL: Must re-fetch the session to get the latest events!
	fmt.Println("  [System] Re-fetching Session 1 to get latest events...")
	session1UpdatedResp, err := sessionService.Get(ctx, &session.GetRequest{
		AppName:   appName,
		UserID:    userID,
		SessionID: session1ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("  [System] Ingesting Session 1 into memory...")
	if err := memService.AddSession(ctx, session1UpdatedResp.Session); err != nil {
		log.Fatal(err)
	}

	// --- Session 2: Recalling Information ---
	fmt.Println("\n--- Session 2 ---")
	session2Resp, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName: appName,
		UserID:  userID,
	})
	if err != nil {
		log.Fatal(err)
	}

	// This should trigger the 'recall' tool because the information is in a *different* session.
	runTurn(ctx, r, session2Resp.Session.ID(), userID, "what is my favorite color")
}

func runTurn(ctx context.Context, r *runner.Runner, sessionID, userID, prompt string) {
	fmt.Printf("User: %s\n", prompt)
	fmt.Print("Agent: ")
	userMsg := genai.NewContentFromText(prompt, genai.RoleUser)
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
