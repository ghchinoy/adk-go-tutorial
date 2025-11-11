# Tutorial 04: Sequential Orchestration

In this tutorial, you will learn how to compose multiple specialized agents into a single workflow. We will build a "Joke Machine" where one agent invents a topic and another agent writes a joke about it.

## Core Concepts

*   **Workflow Agents**: Special agents that don't call LLMs directly but instead manage other agents.
*   **`sequentialagent`**: A workflow agent that runs its sub-agents one after another in a fixed order.
*   **Shared History**: How sub-agents in a sequence can see each other's outputs.

## Prerequisites

*   Completion of [Tutorial 01: Quickstart](01_quickstart.md).

## The Code

We will build this in `main.go`.

### 1. Define Specialized Sub-Agents

Instead of one general-purpose agent, we create two highly specialized ones. This is a key pattern in building reliable AI systems: smaller, focused agents are easier to prompt and test.

```go
	// Agent 1: specialized in being creative and random.
	ideaAgent, _ := llmagent.New(llmagent.Config{
		Name:        "idea_generator",
		Model:       model,
		Instruction: "Generate ONE random, funny, specific topic. Output ONLY the topic.",
	})

	// Agent 2: specialized in humor writing.
	jokeAgent, _ := llmagent.New(llmagent.Config{
		Name:        "joke_writer",
		Model:       model,
		Instruction: "Write a short, punchy joke about the topic provided by the previous agent.",
	})
```

### 2. Create the Orchestrator

We use `sequentialagent.New` to wrap them. The order in `SubAgents` defines the execution order.

```go
	orchestrator, _ := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "joke_machine",
			SubAgents:   []agent.Agent{ideaAgent, jokeAgent},
		},
	})
```

### 3. Launch the Orchestrator

We launch the `orchestrator` just like any other agent. The launcher doesn't need to know it's composed of multiple parts.

```go
	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(orchestrator),
	}
    // ... launch ...
```

## Running the Agent

```bash
go run main.go "Make me laugh"
```

**Expected Output:**
You will see the output from both agents in sequence.

```text
[idea_generator]: Underwater basket weaving for squirrels.
[joke_writer]: Why did the squirrel fail his underwater basket weaving class?
Because he kept trying to bury the bubbles!
```

## Concept Deep Dive: Shared History

By default, agents in a `sequentialagent` share the same conversation history.
1.  User says "Make me laugh".
2.  `ideaAgent` sees "Make me laugh" and outputs a topic.
3.  `jokeAgent` sees "Make me laugh" AND the topic from `ideaAgent`.

This implicit data passing is what makes sequential workflows so powerful for multi-step reasoning.
