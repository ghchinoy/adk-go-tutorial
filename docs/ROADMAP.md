# ADK Go Samples Roadmap

This roadmap outlines the plan for creating a comprehensive set of tutorials and samples for the ADK Go SDK.

## Phase 1: Core Concepts (Complete)

- [x] **Tutorial 01: Quickstart**
    - Goal: Get a basic agent running.
    - Covers: `llmagent`, `runner`, `gemini.NewClient`, basic CLI interaction.
- [x] **Tutorial 02: Custom Tools**
    - Goal: Extend agent capabilities with local functions.
    - Covers: `functiontool`, defining tools, registering them with an agent.
- [x] **Tutorial 03: Session State**
    - Goal: Manage state across turns in a conversation.
    - Covers: `session.Service` (in-memory), accessing state in tools via `tool.Context`.

## Phase 2: Orchestration Patterns (Complete)

- [x] **Tutorial 04: Sequential Orchestration**
    - Goal: Chain agents together.
    - Covers: `sequentialagent`, passing context between agents.
- [x] **Tutorial 05: Parallel Orchestration**
    - Goal: Run multiple agents concurrently and aggregate results.
    - Covers: `parallelagent`, fan-out/fan-in pattern.
- [x] **Tutorial 06: Loops and Conditions**
    - Goal: Implement complex flows with decision making.
    - Covers: `loopagent`, routing based on agent output or state.

## Phase 3: Advanced Features (Complete)

- [x] **Tutorial 07: Artifacts**
    - Goal: Manage large or binary data.
    - Covers: `artifact.Service`, saving/retrieving artifacts in tools.
- [x] **Tutorial 08: Human-in-the-Loop (HITL)**
    - Goal: Involve users in agent decisions.
    - Covers: Implementing a "ask human" tool, pausing/resuming flows (conceptual for CLI).
- [x] **Tutorial 09: Long-term Memory**
    - Goal: Retain information across sessions.
    - Covers: `memory.Service`, ingesting sessions, searching memory.

## Future Considerations

- [ ] **Web-based Agents:** Using the `web` launcher for HTTP/SSE interfaces.
- [ ] **Persistent Storage:** Using Firestore or other databases for sessions and memory.
- [ ] **Observability:** Integrating OpenTelemetry for tracing agent execution.
- [ ] **Multi-modal Inputs:** Handling images and audio in agent conversations.