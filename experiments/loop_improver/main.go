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
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/loopagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/exitlooptool"
	"google.golang.org/genai"
)

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

	// 1. The Writer
	writer, err := llmagent.New(llmagent.Config{
		Name:        "writer",
		Model:       model,
		Instruction: "You are a comedy writer. Write a short joke about the user's topic. If you receive feedback, improve your joke.",
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. The Critic with the Exit Tool
	exitTool, err := exitlooptool.New()
	if err != nil {
		log.Fatal(err)
	}

	critic, err := llmagent.New(llmagent.Config{
		Name:  "critic",
		Model: model,
		Instruction: "You are a harsh comedy critic. Rate the previous joke on a scale of 1-10. " +
			"If the rating is 8 or higher, call the exit_loop tool. " +
			"If it's lower, provide specific, constructive feedback on how to make it funnier.",
		Tools: []tool.Tool{exitTool},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 3. The Loop
	// It will run [writer -> critic] repeatedly.
	loop, err := loopagent.New(loopagent.Config{
		AgentConfig: agent.Config{
			Name:        "writers_room",
			SubAgents:   []agent.Agent{writer, critic},
		},
		MaxIterations: 3, // Safety limit so we don't burn tokens forever if the critic is never satisfied.
	})
	if err != nil {
		log.Fatal(err)
	}

	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(loop),
	}
	l := full.NewLauncher()

	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"Recursion"}
	}

	if err := l.Execute(ctx, config, args); err != nil {
		log.Fatalf("run failed: %v", err)
	}
}
