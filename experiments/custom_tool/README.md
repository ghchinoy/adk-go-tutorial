# Tutorial 02: Custom Tools

In this tutorial, you will learn how to extend an agent's capabilities by giving it a custom tool written in Go. We will build an agent that can roll virtual dice.

## Core Concepts

*   **`functiontool`**: The ADK package for turning standard Go functions into tools an LLM can call.
*   **Schema Inference**: How ADK uses Go structs and reflection to automatically generate the JSON schema required by the LLM.

## Prerequisites

*   Completion of [Tutorial 01: Quickstart](01_quickstart.md).

## The Code

We will build a new agent in `main.go`.

### 1. Define Input/Output Structs

The most idiomatic way to define a tool in ADK-Go is to start with its input and output. We define standard Go structs. ADK uses the `json` tags to name the fields in the schema sent to the LLM, and `jsonschema` tags to provide descriptions.

```go
package main

import (
    // ... imports
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// Input defines what the LLM needs to provide to call this tool.
type RollDiceInput struct {
	NumDice int `json:"num_dice" jsonschema:"description=The number of dice to roll, default is 1"`
	Sides   int `json:"sides" jsonschema:"description=The number of sides on each die, default is 6"`
}

// Output defines what the tool returns to the LLM.
type RollDiceOutput struct {
	Rolls []int `json:"rolls"`
	Total int   `json:"total"`
}
```

### 2. Write the Handler Function

Next, write the actual Go function that performs the logic. It must match a specific signature: it takes a `tool.Context` and your Input struct, and returns your Output struct.

```go
func rollDiceHandler(ctx tool.Context, input RollDiceInput) RollDiceOutput {
	// Set defaults if the LLM didn't provide them
	if input.NumDice <= 0 { input.NumDice = 1 }
	if input.Sides <= 0 { input.Sides = 6 }

	rolls := make([]int, input.NumDice)
	total := 0
	for i := 0; i < input.NumDice; i++ {
		roll := rand.Intn(input.Sides) + 1
		rolls[i] = roll
		total += roll
	}

	return RollDiceOutput{Rolls: rolls, Total: total}
}
```

### 3. Register the Tool and Agent

In your `main` function, use `functiontool.New` to wrap your handler. Go's generics automatically infer the input/output types and generate the schema.

```go
func main() {
    // ... model initialization ...

	// Create the tool
	diceTool, err := functiontool.New(functiontool.Config{
		Name:        "roll_dice",
		Description: "Rolls one or more dice and returns the results.",
	}, rollDiceHandler)
	if err != nil {
		log.Fatalf("Failed to create tool: %v", err)
	}

	// Give the tool to the agent
	agent, err := llmagent.New(llmagent.Config{
		Name:        "gambler",
		Model:       model,
		Instruction: "You are a helpful assistant that can roll dice. Call the roll_dice tool when asked.",
		Tools: []tool.Tool{

			diceTool,
		},
	})
    // ... launch as usual ...
}
```

## Running the Agent

```bash
go run main.go "Roll 3 d20s for me"
```

**Expected Output:**
The agent will call your Go function, get the random numbers, and then write a response.

```text
I rolled 3 d20s. The results were [15, 2, 19], for a total of 36.
```

## Concept Deep Dive: Schema Inference

How does `functiontool.New` know what to tell Gemini? It uses reflection on your `RollDiceInput` struct.

*   `NumDice int` tells Gemini this parameter is an integer.
*   `json:"num_dice"` tells Gemini the parameter name is `num_dice`.
*   `jsonschema:"description=..."` provides the documentation Gemini reads to understand *when* and *how* to use this parameter.

This "Code-First" approach means your Go code *is* your tool definition, keeping everything in sync.
