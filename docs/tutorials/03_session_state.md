# Tutorial 03: Session State & Memory

In this tutorial, you will learn how to make your agent "remember" information across different turns of a conversation. We will build an agent that can save and recall a user's favorite color.

## Core Concepts

*   **`session.State`**: A key-value store attached to the current session.
*   **`tool.Context`**: The interface passed to every tool handler, providing access to the state.
*   **State Scopes**: Understanding how `KeyPrefixUser` affects data visibility.

## Prerequisites

*   Completion of [Tutorial 02: Custom Tools](02_custom_tools.md).

## The Code

We will build an interactive CLI agent in `main.go`.

### 1. Define Tools for State Access

Agents don't magically "remember" things in ADK. You must give them tools to read and write to their state.

```go
package main

import (
    // ... imports
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
)

// ... Input/Output structs (see full code) ...

func saveColorHandler(ctx tool.Context, input SaveColorInput) SaveColorOutput {
	// We use KeyPrefixUser to make this a user-level preference.
	// It will be available in ANY session for this specific user_id.
	key := session.KeyPrefixUser + "fav_color"
	err := ctx.State().Set(key, input.Color)
    // ... handle error and return success ...
}

func getColorHandler(ctx tool.Context, _ GetColorInput) GetColorOutput {
	key := session.KeyPrefixUser + "fav_color"
	val, err := ctx.State().Get(key)
	if err != nil {
		return GetColorOutput{Color: "unknown"}
	}
    // ... cast val to string and return ...
}
```

### 2. The Agent's Instructions

You must explicitly tell the LLM *when* to use these tools in its system instructions.

```go
	agent, err := llmagent.New(llmagent.Config{
        // ...
		Instruction: "You are a helpful assistant that remembers user preferences. " +
			"If the user tells you their favorite color, save it using the save_favorite_color tool. " +
			"If asked about it, use the get_favorite_color tool to look it up.",
		Tools: []tool.Tool{saveTool, getTool},
	})
```

## Running the Agent (Interactive Mode)

For this tutorial, we need a multi-turn conversation. Run the program *without* arguments to enter the interactive console.

```bash
go run main.go
```

**Interaction Example:**

```text
User -> Hi, my favorite color is blue.
Agent -> Okay, I've saved that your favorite color is blue. [Calls save_favorite_color]

User -> What is my favorite color?
Agent -> Your favorite color is blue. [Calls get_favorite_color]
```

## Concept Deep Dive: State Scopes

ADK supports different scopes for state keys, managed via prefixes:

*   **No Prefix** (e.g., `"current_task"`): **Session Scope**. Visible only within the current conversation thread.
*   **`session.KeyPrefixUser`** (e.g., `"user:fav_color"`): **User Scope**. Visible across *all* sessions for the same `user_id`.
*   **`session.KeyPrefixApp`** (e.g., `"app:global_config"`): **App Scope**. Visible to all users of the application.

By using `session.KeyPrefixUser`, we ensure that if this same user starts a *new* chat session tomorrow, the agent will still know their favorite color (assuming you are using a persistent `SessionService` like Vertex AI, as discussed in the [Sessions Explainer](../explainer_sessions_and_backends.md)).
