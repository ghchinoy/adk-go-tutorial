# Tutorial 07: Generating Artifacts

In this tutorial, you will learn how to let your agent generate and save files (artifacts) that are attached to the user's session. We will build a "Reporter" agent that writes a text report and saves it as a file.

## Core Concepts

*   **Artifacts**: Files (text, images, etc.) generated during a session that need to be persisted separately from the chat history.
*   **`ctx.Artifacts()`**: The interface in `tool.Context` for managing these files.
*   **`artifact.Service`**: The backend service that stores the files (In-Memory, GCS, etc.).

## Prerequisites

*   Completion of [Tutorial 02: Custom Tools](02_custom_tools.md).

## The Code

We will build this in `main.go`.

### 1. The Save Tool

We create a custom tool that takes a filename and content, and uses `ctx.Artifacts().Save()` to store it.

```go
func saveReportHandler(ctx tool.Context, input SaveReportInput) SaveReportOutput {
	// ctx.Artifacts() automatically handles the current session ID.
	// We wrap the content in a genai.Part.
	_, err := ctx.Artifacts().Save(context.Background(), input.Filename, genai.NewPartFromText(input.Content))
	if err != nil {
		return SaveReportOutput{Success: false}
	}
	return SaveReportOutput{Success: true}
}
```

### 2. Configuring the Launcher

Crucially, we must tell the ADK launcher *where* to store these artifacts by providing an `ArtifactService`. For this tutorial, we use the in-memory version.

```go
	config := &adk.Config{
		AgentLoader:     services.NewSingleAgentLoader(agent),
		// Enable artifact storage
		ArtifactService: artifact.InMemoryService(),
	}
```

## Running the Agent

```bash
go run main.go "Write a poem about compilation errors and save it as poem.txt"
```

**Expected Output:**
The agent will generate the poem and call the tool. You should see our debug print confirming the save.

```text
[SYSTEM] Saved artifact 'poem.txt' (Content length: 142)
[reporter]: I have written the poem and saved it as poem.txt.
```

## Concept Deep Dive: Artifact Services

Just like Session Services, Artifact Services can be swapped out.
*   **`artifact.InMemoryService()`**: Good for testing, lost on restart.
*   **GCS (Google Cloud Storage)**: For production, you would typically use a GCS-backed service (not shown here, but follows the same pattern as Vertex AI sessions) to store files permanently.
