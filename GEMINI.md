# ADK Go Samples - Assistant Guide

**Note:** This tutorial series utilizes Google Cloud's Vertex AI for its Generative AI models and requires a Google Cloud Project with the Vertex AI API enabled.

This repository contains a structured set of tutorials for the Go ADK. When assisting a user with this repository, follow these guidelines to ensure a smooth experience.

## Tutorial Structure

The tutorials are numbered 01-09 in `docs/tutorials/` and have corresponding code in `experiments/`. Each `experiments/` directory now also contains a `README.md` file summarizing its tutorial.

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

## Operational Guidelines

### 1. Initialization & Setup
All experiment directories have pre-initialized Go modules. Before running any experiment, ensure dependencies are up to date:

```bash
cd experiments/<experiment_name>
go mod tidy
```

### 2. Authentication Standardization
For consistency in enterprise/cloud environments, standardizing on Vertex AI is recommended.

**Action:** Before running any experiment, you MUST ask the user for their Google Cloud Project ID and Location. Then, check `main.go` and replace the `gemini.NewModel` initialization to use Vertex AI.

*Old:*
```go
model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
    APIKey: os.Getenv("GOOGLE_API_KEY"),
})
```

*New:*
```go
model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
    Backend:  genai.BackendVertexAI,
    Project:  os.Getenv("GOOGLE_CLOUD_PROJECT"),
    Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
})
```

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

