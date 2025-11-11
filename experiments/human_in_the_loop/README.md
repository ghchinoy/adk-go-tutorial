# Tutorial 08: Human-in-the-Loop (HITL)

In this tutorial, you will learn how to create workflows where the agent pauses execution to ask the human user for input or confirmation before proceeding. We will build a "Careful Agent" that asks for permission before doing anything "dangerous".

## Core Concepts

*   **HITL (Human-in-the-Loop)**: Integrating human judgment into automated AI workflows.
*   **Blocking Tools**: A simple way to implement HITL in CLI agents is by creating a tool that synchronously waits for user input.

## Prerequisites

*   Completion of [Tutorial 02: Custom Tools](02_custom_tools.md).

## The Code

We will build this in `main.go`.

### 1. The 'Ask Human' Tool

We create a tool that prints a question to the console and uses `bufio` to read a line from standard input. This effectively pauses the agent until the user responds.

```go
func askHumanHandler(ctx tool.Context, input AskHumanInput) AskHumanOutput {
	fmt.Printf("\n[AGENT ASKS]: %s\n[YOU ANSWER] > ", input.Question)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	return AskHumanOutput{Answer: strings.TrimSpace(answer)}
}
```

### 2. The Careful Agent

We instruct the agent to *always* use this tool before taking specific actions.

```go
	agent, _ := llmagent.New(llmagent.Config{
		Name:  "careful_agent",
		Model: model,
		Instruction: `You are a helpful assistant.
If the user asks you to do something "dangerous" (like deleting files),
    you MUST first use the 'ask_human' tool to get explicit confirmation.`,
		Tools: []tool.Tool{askTool},
	})
```

## Running the Agent

Run the agent and ask it to do something dangerous.

```bash
go run main.go "Please delete all my files."
```

**Expected Output:**

```text
[AGENT ASKS]: Are you sure you want to delete all your files? This action cannot be undone.
[YOU ANSWER] > no
[careful_agent]: Okay, I will not delete your files.
```

Try running it again and answering "yes" to see how it proceeds (it should "pretend" to delete them based on our instructions).

## Concept Deep Dive: Async HITL

While this synchronous approach works great for CLIs, real-world web applications often need **asynchronous HITL**.
In that pattern, the tool would return a special "pending" status, the ADK would suspend the session, and a separate API call (e.g., from a web UI) would later resume the session with the human's answer.
ADK Go's architecture supports this via its event-driven model, but it requires more complex setup (custom runners and persistence) than this tutorial covers.


