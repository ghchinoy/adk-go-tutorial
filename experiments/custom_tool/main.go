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
	"fmt"
	"log"
	"math/rand"
	"os"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// 1. Define Input/Output structs for your tool.
// ADK uses these to automatically generate the JSON schema for the LLM.
type RollDiceInput struct {
	NumDice int `json:"num_dice"`
	Sides   int `json:"sides"`
}

type RollDiceOutput struct {
	Rolls []int `json:"rolls"`
	Total int   `json:"total"`
}

// 2. Define the handler function.
// It must match the signature: func(tool.Context, Input) Output
func rollDiceHandler(ctx tool.Context, input RollDiceInput) RollDiceOutput {
	// Set defaults if zero values are passed
	if input.NumDice <= 0 {
		input.NumDice = 1
	}
	if input.Sides <= 0 {
		input.Sides = 6
	}

	rolls := make([]int, input.NumDice)
	total := 0
	for i := 0; i < input.NumDice; i++ {
		roll := rand.Intn(input.Sides) + 1
		rolls[i] = roll
		total += roll
	}

	fmt.Printf("DEBUG: Rolling %d d%d: %v (Total: %d)\n", input.NumDice, input.Sides, rolls, total)
	return RollDiceOutput{
		Rolls: rolls,
		Total: total,
	}
}

func main() {
	ctx := context.Background()

	// Initialize Model
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		Backend:  genai.BackendVertexAI,
		Project:  os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	// 3. Create the Tool
	// We use functiontool.New with our handler. Go's generics handle the rest.
	diceTool, err := functiontool.New(functiontool.Config{
		Name:        "roll_dice",
		Description: "Rolls one or more dice and returns the results.",
	}, rollDiceHandler)
	if err != nil {
		log.Fatalf("Failed to create dice tool: %v", err)
	}

	// 4. Create Agent with the Tool
	agent, err := llmagent.New(llmagent.Config{
		Name:        "gambler_agent",
		Model:       model,
		Description: "An agent that can roll dice.",
		Instruction: "You are a helpful assistant that can roll dice for the user. " +
			"When asked to roll dice, call the roll_dice tool and report the results.",
		Tools: []tool.Tool{
			diceTool,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// 5. Launch
	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(agent),
	}
	l := full.NewLauncher()
	if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("run failed: %v", err)
	}
}
