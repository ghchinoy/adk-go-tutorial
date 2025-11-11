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
	"google.golang.org/adk/agent/workflowagents/parallelagent"
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

	optimist, err := llmagent.New(llmagent.Config{
		Name:        "optimist",
		Model:       model,
		Instruction: "You are an eternal optimist. Give a short, positive take on the user's topic.",
	})
	if err != nil {
		log.Fatal(err)
	}

	pessimist, err := llmagent.New(llmagent.Config{
		Name:        "pessimist",
		Model:       model,
		Instruction: "You are a grumpy pessimist. Give a short, negative take on the user's topic.",
	})
	if err != nil {
		log.Fatal(err)
	}

	// The Orchestrator: Parallel Agent
	// It will run both agents at the same time.
	orchestrator, err := parallelagent.New(parallelagent.Config{
		AgentConfig: agent.Config{
			Name:        "debate_team",
			Description: "Gets two opposing viewpoints on a topic.",
			SubAgents:   []agent.Agent{optimist, pessimist},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(orchestrator),
	}
	l := full.NewLauncher()

	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"Artificial Intelligence"}
	}

	if err := l.Execute(ctx, config, args); err != nil {
		log.Fatalf("run failed: %v", err)
	}
}
