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
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
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

	// Agent 1: The Idea Generator
	ideaAgent, err := llmagent.New(llmagent.Config{
		Name:        "idea_generator",
		Model:       model,
		Description: "Generates a random, funny topic.",
		Instruction: "You are a creative assistant. When asked, generate ONE random, funny, and specific topic for a joke. Output ONLY the topic, nothing else.",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Agent 2: The Joke Writer
	jokeAgent, err := llmagent.New(llmagent.Config{
		Name:        "joke_writer",
		Model:       model,
		Description: "Writes a joke about a given topic.",
		Instruction: "You are a professional comedian. Write a short, punchy joke about the topic provided by the previous agent.",
	})
	if err != nil {
		log.Fatal(err)
	}

	// The Orchestrator: Sequential Agent
	// It will run ideaAgent, wait for it to finish, then run jokeAgent.
	orchestrator, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "joke_machine",
			Description: "Generates a topic and then writes a joke about it.",
			SubAgents:   []agent.Agent{ideaAgent, jokeAgent},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(orchestrator),
	}
	l := full.NewLauncher()

	// If no args are provided, we supply a default prompt to kick off the sequence.
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"Go!"}
	}

	if err := l.Execute(ctx, config, args); err != nil {
		log.Fatalf("run failed: %v", err)
	}
}
