# ADK for Go: Primer

This is a living document explaining the core concepts and architecture of the Agent Development Kit (ADK) for Go.

## High-Level Architecture

ADK-Go is designed as a standard Go library, avoiding heavy framework magic in favor of explicit composition.

### Core Components

1.  **Agent (`agent.Agent`)**
    *   The fundamental interface. An Agent is anything that can take an input and produce a sequence of events (responses, tool calls, etc.).
    *   **`llmagent`**: The most common implementation. It wraps an LLM (like Gemini) and manages the loop of sending user input to the model and processing the model's response (including executing tool calls automatically).
    *   **Workflow Agents**: Agents that don't call LLMs directly but orchestrate other agents (e.g., `sequentialagent`, `parallelagent`).

2.  **Model (`model.Model`)**
    *   An abstraction over the raw LLM API.
    *   ADK provides a robust `gemini` implementation out-of-the-box.
    *   Responsible for formatting standard ADK requests into the specific format required by the underlying model provider.

3.  **Tool (`tool.Tool`)**
    *   Defines capabilities an agent can use.
    *   Tools are standard Go structs with methods that can be exposed to the LLM.
    *   ADK handles generating the schema for the LLM and invoking the Go method when the LLM requests it.

4.  **Launcher (`cmd/launcher`)**
    *   The runtime environment for an agent.
    *   It standardizes how agents are started, how configuration is loaded, and how they interact with the outside world (CLI stdio, HTTP, etc.).
    *   Using a launcher is recommended for production consistency but optional for simple tests.

## Key Concepts

### The Agent Loop
When an `llmagent` runs:
1.  It receives `InvocationContext` (user input, history, state).
2.  It sends this context to the configured `Model`.
3.  The `Model` returns a response.
4.  If the response is a **Tool Call**, the agent executes the corresponding Go function and feeds the result back to the Model (repeating the loop).
5.  If the response is a **Final Answer** (text), the agent yields it to the caller and ends the turn.

### Composability
Because every component (LLM agents, workflow agents) implements the same `agent.Agent` interface, they can be nested arbitrarily. A `parallelagent` can run three `llmagent`s, or it could run two `llmagent`s and one `sequentialagent`.
