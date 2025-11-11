# Tutorial 01: Quickstart (Vertex AI)

This tutorial will guide you through building your first agent using ADK for Go, utilizing the **Vertex AI** backend. We will create a simple CLI agent capable of answering questions and using a tool (Google Search) to find real-time information.

## Prerequisites

* Go 1.23 or higher.
* A Google Cloud Project with the **Vertex AI API** enabled.
* Authentication set up for your environment (e.g., Application Default Credentials via `gcloud auth application-default login` or a service account).
* The following environment variables set:
    * `GOOGLE_CLOUD_PROJECT`: Your Google Cloud Project ID.
    * `GOOGLE_CLOUD_LOCATION`: The region for Vertex AI (e.g., `us-central1`).

## The Code

We will build this agent in a single `main.go` file.

### 1. Setup and Imports

First, we import the necessary packages. Notice we are importing standard ADK components for Gemini, the LLM agent implementation, and the standard launcher.

```go
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)
```

### 2. Initialize the Model (Vertex AI)

Inside `main`, we initialize the Gemini model, explicitly configuring it to use the Vertex AI backend.

```go
func main() {
	ctx := context.Background()

	// Initialize the Gemini model using the Vertex AI backend.
	// We read the required project and location from environment variables.
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		Backend:  genai.BackendVertexAI,
		Project:  os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}
    // ...
```

### 3. Define the Agent

Next, we create the agent itself using `llmagent.New`. We give it a name, a description, instructions on how to behave, and equip it with the `GoogleSearch` tool.

```go
    // ...
	agent, err := llmagent.New(llmagent.Config{
		Name:        "hello_time_agent",
		Model:       model,
		Description: "Tells the current time in a specified city.",
		Instruction: "You are a helpful assistant that tells the current time in a city.",
		Tools: []tool.Tool{
			// ADK includes pre-built tools like GoogleSearch.
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}
    // ...
```

### 4. Launch the Agent

Finally, we use the `full.Launcher` to run our agent. The launcher handles standard CLI argument parsing and sets up the runtime environment.

```go
    // ...
	config := &adk.Config{
		// We tell the launcher to load our single agent.
		AgentLoader: services.NewSingleAgentLoader(agent),
	}

	l := full.NewLauncher()
	// Execute runs the agent loop, processing input from os.Args or stdin depending on flags.
	err = l.Execute(ctx, config, os.Args[1:])
	if err != nil {
		log.Fatalf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
```

## Running the Agent

With the code saved to `main.go`, configure your environment and run it:

```bash
# Replace with your actual project ID and preferred location
export GOOGLE_CLOUD_PROJECT="your-project-id"
export GOOGLE_CLOUD_LOCATION="us-central1"

# Ensure you are authenticated
# gcloud auth application-default login

go run main.go "What is the current time in Tokyo?"
```

**Expected Output:**
The agent will likely use the Google Search tool to find the current time and then synthesize an answer for you.

```text
The current time in Tokyo, Japan is [current time].
```
