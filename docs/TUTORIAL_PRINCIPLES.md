# ADK for Go: Tutorial Principles

This document outlines the guiding principles for creating tutorials and documentation for the Go Agent Development Kit (ADK).

## Core Philosophy

1.  **Code-First & Idiomatic Go**: Tutorials should emphasize ADK's nature as a standard Go library. Avoid over-abstracting; show standard Go patterns (explicit error handling, context propagation, standard project layouts).
2.  **Composability**: Highlight how small, focused agents can be composed into larger systems. Start simple, then show how to wrap or combine agents.
3.  **Progressive Disclosure**: Start with the absolute minimum code to get a working agent (the Quickstart). Introduce complex concepts (Tools, Memory, Workflow Agents) only after the basics are solid.
4.  **Runnable Examples**: Every tutorial must have accompanying, complete, runnable code in the repository. Snippets in docs must match the real code.
5.  **Concept-Driven**: Each tutorial must have a primary "learning objective" tied to a core ADK concept (e.g., "Session State", "Custom Tools", "Parallel Execution").

## Tutorial Structure Standard

*   **Goal**: State clearly what the user will build and which **Core Concepts** they will learn.
*   **Prerequisites**: List required tools and prior knowledge.
*   **Step-by-Step Implementation**:
    *   Break down the code into logical blocks.
    *   Explain the *why* behind each block, tying it back to the core concepts.
*   **Run & Verify**: Show exactly how to run the example and expected output.
*   **Concept Deep Dive**: A dedicated section summarizing *how* ADK implements the concepts used in the tutorial (e.g., "How ADK marshals Go structs to Tool schemas").
*   **Next Steps**: Link to the next logical tutorial.

## Tone and Style

*   **Professional & Direct**: Write for competent Go developers.
*   **Opinionated but Flexible**: Show the recommended "happy path" while acknowledging alternatives.