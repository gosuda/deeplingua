package main

import (
	"context"
	"os"

	"github.com/lemon-mint/coord"
	"github.com/lemon-mint/coord/llm"
	"github.com/lemon-mint/coord/pconf"
	_ "github.com/lemon-mint/coord/provider/vertexai"
)

func init() {
	// copy this to models.go and remove "return" below
	return

	project_id := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project_id == "" {
		panic("GOOGLE_CLOUD_PROJECT environment variable is not set")
	}

	client, err := coord.NewLLMClient(context.Background(), "vertexai", pconf.WithProjectID(project_id), pconf.WithLocation("us-central1"))
	if err != nil {
		panic(err)
	}

	translationModel, err = client.NewLLM("gemini-2.0-flash-001", &llm.Config{
		MaxOutputTokens:       Ptr(8192),
		SafetyFilterThreshold: llm.BlockOff,
	})
}
