# ADK Go Samples - Assistant Guide

**Note:** This tutorial series utilizes Google Cloud's Vertex AI for its Generative AI models and requires a Google Cloud Project with the Vertex AI API enabled.

This repository contains a structured set of tutorials for the Go ADK. When assisting a user with this repository, follow these guidelines to ensure a smooth experience.

## Tutorial Structure

The tutorials are numbered 01-09 in `docs/tutorials/` and have corresponding code in `experiments/`.

| Experiment Directory | Key Concept |
|---|---|
| `quickstart` | Basic agent, `console` launcher. |
| `custom_tool` | `functiontool` implementation. |
| `session_state` | Using `ctx.State()` for short-term memory. |
| `sequential_jokes` | `sequentialagent` (chaining). |
| `parallel_perspectives` | `parallelagent` (concurrency). |
| `loop_improver` | `loopagent` with `exitlooptool`. |
| `generate_artifact` | Using `ctx.Artifacts()` for file generation. |
| `human_in_the_loop` | Pausing for user input via tools. |
| `long_term_memory` | Using `memory.Service` across sessions (no launcher). |

## Teaching Workflow

When guiding a user through the tutorials, follow this sequence for *each* numbered tutorial:

1.  **Concept Introduction:** Read the corresponding `docs/tutorials/XX_*.md` file. Summarize the "Core Concepts" for the user *before* running the code.
2.  **Code Walkthrough:** Briefly highlight how those concepts are implemented in the `experiments/XX_*/main.go` file.
3.  **Execution:** Run the experiment using the appropriate method (usually `console` mode for interactivity).
4.  **Pause for Understanding:** Explicitly ask the user if they have questions about the current tutorial's concepts before offering to move to the next one.

## Operational Guidelines

### 1. Initialization & Setup
All experiment directories have pre-initialized Go modules. Before running any experiment, ensure dependencies are up to date:

```bash
cd experiments/<experiment_name>
go mod tidy
```

### 2. Authentication (Vertex AI Standard)
All tutorials in this repository have been standardized to use the **Vertex AI** backend.
*   Ensure the user has provided `GOOGLE_CLOUD_PROJECT` and `GOOGLE_CLOUD_LOCATION`.
*   Export these as environment variables before running any experiment.

### 3. Running Agents
The `full.NewLauncher` used in these samples primarily supports `console` and `web` modes. It does **not** support a standalone `run` command for single-turn input in standard `os.Args`.

*   **Interactive Mode:** Use `go run main.go console`.
*   **Single-Turn Testing:** Pipe input to the console mode for reliable automated testing:
    ```bash
    printf "Your input here\n" | go run main.go console
    ```
*   **Multi-Turn Testing:** Use `printf` with multiple lines:
    ```bash
    printf "First turn\nSecond turn\n" | go run main.go console
    ```
*   **Interactive Actions with Environment Variables:** When performing interactive actions that require environment variables, use the following pattern:
    ```bash
    export GOOGLE_CLOUD_PROJECT=<your-project-id>
    export GOOGLE_CLOUD_LOCATION=<your-location>
    cd experiments/<experiment_name>
    printf "Your input here\n" | go run main.go console
    ```

### 4. Known Issues & Fixes
*   **JSON Schema Tags:** Some `jsonschema` tags in the samples (e.g., `description=...`) may cause runtime errors with strict parsers.
    *   *Fix:* If an error like `tag must not begin with 'WORD='` occurs, simplify or remove the `jsonschema` struct tags in the experiment's `main.go`.