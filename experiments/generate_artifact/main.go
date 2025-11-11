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
	"google.golang.org/adk/artifact"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

type SaveReportInput struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type SaveReportOutput struct {
	Success bool `json:"success"`
}

func saveReportHandler(ctx tool.Context, input SaveReportInput) SaveReportOutput {
	// We use the Artifacts service from the context.
	// It automatically handles AppName, UserID, and SessionID.
	_, err := ctx.Artifacts().Save(context.Background(), input.Filename, genai.NewPartFromText(input.Content))
	if err != nil {
		log.Printf("Error saving artifact: %v", err)
		return SaveReportOutput{Success: false}
	}
	fmt.Printf("\n[SYSTEM] Saved artifact '%s' (Content length: %d)\n", input.Filename, len(input.Content))
	return SaveReportOutput{Success: true}
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

	saveTool, err := functiontool.New(functiontool.Config{
		Name:        "save_report",
		Description: "Saves a text report to the user's session artifacts.",
	}, saveReportHandler)
	if err != nil {
		log.Fatal(err)
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:        "reporter",
		Model:       model,
		Instruction: "You are a researcher. When asked to write a report, generate the content and then ALWAYS save it using the save_report tool.",
		Tools:       []tool.Tool{saveTool},
	})
	if err != nil {
		log.Fatal(err)
	}

	// We MUST provide an ArtifactService to the launcher.
	// InMemoryService is good for testing.
	config := &adk.Config{
		AgentLoader:     services.NewSingleAgentLoader(agent),
		ArtifactService: artifact.InMemoryService(),
	}
	l := full.NewLauncher()

	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"Write a very short report about goldfish."}
	}

	if err := l.Execute(ctx, config, args); err != nil {
		log.Fatalf("run failed: %v", err)
	}
}
