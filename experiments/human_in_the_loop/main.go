// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

type AskHumanInput struct {
	Question string `json:"question"`
}

type AskHumanOutput struct {
	Answer string `json:"answer"`
}

// askHumanHandler pauses execution until the user types an answer in the console.
// This is a simple way to implement HITL for CLI agents.
func askHumanHandler(ctx tool.Context, input AskHumanInput) AskHumanOutput {
	fmt.Printf("\n[AGENT ASKS]: %s\n[YOU ANSWER] > ", input.Question)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	return AskHumanOutput{Answer: strings.TrimSpace(answer)}
}

func main() {
	ctx := context.Background()
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		Backend:  genai.BackendVertexAI,
		Project:  os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
	})
	if err != nil {
		log.Fatal(err)
	}

	askTool, err := functiontool.New(functiontool.Config{
		Name:        "ask_human",
		Description: "Asks the human user a question and waits for their response. Use this before taking any 'dangerous' action.",
	}, askHumanHandler)
	if err != nil {
		log.Fatal(err)
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:  "careful_agent",
		Model: model,
		Instruction: `You are a helpful assistant.
If the user asks you to do something "dangerous" (like deleting files, launching missiles, or eating the last cookie),
	you MUST first use the 'ask_human' tool to get explicit confirmation.
If they say "yes", pretend to do it. If they say "no", do not do it.`,
		Tools: []tool.Tool{askTool},
	})
	if err != nil {
		log.Fatal(err)
	}

	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(agent),
	}
	l := full.NewLauncher()

	// Default prompt if none provided
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"Please delete all my files."}
	}

	if err := l.Execute(ctx, config, args); err != nil {
		log.Fatalf("run failed: %v", err)
	}
}
