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
	"os"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// 1. Define Inputs/Outputs for our tools

type SaveColorInput struct {
	Color string `json:"color"`
}

type SaveColorOutput struct {
	Success bool `json:"success"`
}

type GetColorInput struct{}

type GetColorOutput struct {
	Color string `json:"color"`
}

// 2. Define Handlers that use tool.Context to access Session State

func saveColorHandler(ctx tool.Context, input SaveColorInput) SaveColorOutput {
	// We use KeyPrefixUser so this preference is tied to the user,
	// not just this specific conversation session.
	key := session.KeyPrefixUser + "fav_color"
	err := ctx.State().Set(key, input.Color)
	if err != nil {
		log.Printf("Error saving state: %v", err)
		return SaveColorOutput{Success: false}
	}
	return SaveColorOutput{Success: true}
}

func getColorHandler(ctx tool.Context, _ GetColorInput) GetColorOutput {
	key := session.KeyPrefixUser + "fav_color"
	val, err := ctx.State().Get(key)
	if err != nil {
		// Key doesn't exist yet
		return GetColorOutput{Color: "unknown"}
	}
	color, ok := val.(string)
	if !ok {
		return GetColorOutput{Color: "unknown"}
	}
	return GetColorOutput{Color: color}
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		Backend:  genai.BackendVertexAI,
		Project:  os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	// 3. Create the Tools
	saveTool, err := functiontool.New(functiontool.Config{
		Name:        "save_favorite_color",
		Description: "Saves the user's favorite color to their profile.",
	}, saveColorHandler)
	if err != nil {
		log.Fatal(err)
	}

	getTool, err := functiontool.New(functiontool.Config{
		Name:        "get_favorite_color",
		Description: "Retrieves the user's favorite color if known.",
	}, getColorHandler)
	if err != nil {
		log.Fatal(err)
	}

	// 4. Create the Agent
	agent, err := llmagent.New(llmagent.Config{
		Name:  "memory_agent",
		Model: model,
		Instruction: "You are a helpful assistant that remembers user preferences. " +
			"If the user tells you their favorite color, save it. " +
			"If you need to know their favorite color, look it up.",
		Tools: []tool.Tool{saveTool, getTool},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 5. Launch
	// Note: We are using the default InMemory session service here.
	// In a real app, you'd use VertexAIService or a custom one for persistence.
	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(agent),
	}
	l := full.NewLauncher()

	// We use console mode to have a multi-turn conversation.
	// Run with no arguments to enter interactive mode.
	if len(os.Args) > 1 {
		fmt.Println("NOTE: To test memory, run without arguments to enter interactive console mode.")
	}

	if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("run failed: %v", err)
	}
}
