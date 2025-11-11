# Tutorial 06: Loops & Conditions

In this tutorial, you will learn how to create dynamic workflows that repeat until a specific condition is met. We will build a "Writer's Room" where one agent writes a joke and another critiques it, looping until the joke is good enough.

## Core Concepts

*   **`loopagent`**: A workflow agent that repeats its sub-agents.
*   **`exitlooptool`**: A special tool that allows an LLM agent to signal that the loop should terminate.
*   **Iterative Refinement**: Using loops to improve output quality.

## Prerequisites

*   Completion of [Tutorial 04: Sequential Orchestration](04_sequential_orchestration.md).

## The Code

We will build this in `main.go`.

### 1. The Writer Agent

This agent is simple: it writes a joke and, crucially, is instructed to *improve* it if it sees feedback in the history.

```go
	writer, _ := llmagent.New(llmagent.Config{
		Name:        "writer",
		Model:       model,
		Instruction: "Write a joke about the topic. If you see feedback, improve your joke.",
	})
```

### 2. The Critic Agent (with Exit Tool)

This agent decides when the loop ends. We give it the `exit_loop` tool.

```go
	// Create the special exit tool
	exitTool, _ := exitlooptool.New()

	critic, _ := llmagent.New(llmagent.Config{
		Name:  "critic",
		Model: model,
		Instruction: "Rate the joke 1-10. If rating >= 8, call exit_loop. Otherwise, give feedback.",
		Tools: []tool.Tool{exitTool},
	})
```

### 3. The Loop Orchestrator

We wrap them in a `loopagent`. We also set `MaxIterations` as a safety net.

```go
	loop, _ := loopagent.New(loopagent.Config{
		AgentConfig: agent.Config{
			Name:      "writers_room",
			SubAgents: []agent.Agent{writer, critic},
		},
		MaxIterations: 3, // Stop after 3 tries even if not satisfied
	})
```

## Running the Agent

```bash
go run main.go "Recursion"
```

**Expected Output (Simulated):**

*Iteration 1:*
```text
[writer]: To understand recursion, you must first understand recursion.
[critic]: Rating: 4/10. Too cliché. Try something fresher.
```

*Iteration 2 (Loop continues because exit_loop wasn't called):*
```text
[writer]: My friend fell into a recursive function. We're still waiting for him to return.
[critic]: Rating: 8/10. Much better! [Calls exit_loop]
```

*Loop terminates.*

## Concept Deep Dive: Termination Signals

The `loopagent` continues indefinitely (or until `MaxIterations`) unless it receives a specific signal.
The `exitlooptool` provides this signal by setting a special flag (`Actions.Escalate = true`) on the event it generates. The `loopagent` checks for this flag after every sub-agent runs and terminates immediately if it's found.
