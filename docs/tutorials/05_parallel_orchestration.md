# Tutorial 05: Parallel Orchestration

In this tutorial, you will learn how to run multiple agents concurrently. We will build a "Debate Team" where an optimist and a pessimist give their takes on a topic at the same time.

## Core Concepts

*   **`parallelagent`**: A workflow agent that executes all its sub-agents concurrently.
*   **Branched History**: How ADK isolates parallel agents so they don't interfere with each other.
*   **Performance**: The latency benefits of parallelization.

## Prerequisites

*   Completion of [Tutorial 04: Sequential Orchestration](04_sequential_orchestration.md).

## The Code

We will build this in `main.go`.

### 1. Define Sub-Agents

We create two simple agents with opposing personalities.

```go
	optimist, _ := llmagent.New(llmagent.Config{
		Name:        "optimist",
		Model:       model,
		Instruction: "You are an eternal optimist. Give a positive take.",
	})

	pessimist, _ := llmagent.New(llmagent.Config{
		Name:        "pessimist",
		Model:       model,
		Instruction: "You are a grumpy pessimist. Give a negative take.",
	})
```

### 2. Create the Parallel Orchestrator

We use `parallelagent.New`.

```go
	orchestrator, _ := parallelagent.New(parallelagent.Config{
		AgentConfig: agent.Config{
			Name:        "debate_team",
			SubAgents:   []agent.Agent{optimist, pessimist},
		},
	})
```

## Running the Agent

```bash
go run main.go "Remote Work"
```

**Expected Output:**
You will see outputs from both agents. Because they run in parallel, the order they appear in the console might vary slightly depending on which one finishes first, though the launcher typically tries to stream them as they come in.

```text
[optimist]: Remote work is amazing! It gives people flexibility and better work-life balance.
[pessimist]: Remote work is isolating. You lose all sense of company culture and human connection.
```

## Concept Deep Dive: Branched History

When `parallelagent` runs, it creates a **branch** of the conversation history for each sub-agent.
*   The `optimist` sees the user's prompt, but it *does not* see what the `pessimist` is generating, and vice-versa.
*   This isolation is crucial. If they shared history while running in parallel, they might get confused by each other's partial outputs.
*   Once both finish, their final responses are merged back into the main history so subsequent agents (if any) can see both perspectives.
