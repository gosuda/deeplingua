package main

import (
	"context"

	"github.com/lemon-mint/coord"
	"github.com/lemon-mint/coord/llm"
	"github.com/lemon-mint/coord/pconf"
	"github.com/lemon-mint/coord/provider"
	_ "github.com/lemon-mint/coord/provider/aistudio"
	_ "github.com/lemon-mint/coord/provider/anthropic"
	_ "github.com/lemon-mint/coord/provider/openai"
	_ "github.com/lemon-mint/coord/provider/vertexai"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

type Configs struct {
	Models       []Model `json:"models,omitempty"`
	StartIndex   int     `json:"start_index,omitempty"`
	CustomPrompt *string `json:"custom_prompt,omitempty"`
}

type Model struct {
	Provider    string   `json:"provider,omitempty"`
	ModelID     string   `json:"model_id,omitempty"`
	Location    string   `json:"location,omitempty"`
	Project     string   `json:"project,omitempty"`
	Temperature *float32 `json:"temperature,omitempty"`
	RateLimit   *float64 `json:"rate_limit,omitempty"`
	MaxTokens   int      `json:"max_tokens,omitempty"`
	APIKey      string   `json:"api_key,omitempty"`
	BaseURL     string   `json:"base_url,omitempty"`
}

func ApplyConfig(c *Configs) {
	if c == nil {
		return
	}

	models := make([]llm.Model, 0, len(c.Models))
	for i, m := range c.Models {
		var client provider.LLMClient
		var options []pconf.Config
		var err error
		switch m.Provider {
		case "openai":
			if m.APIKey != "" {
				options = append(options, pconf.WithAPIKey(m.APIKey))
			}
			if m.BaseURL != "" {
				options = append(options, pconf.WithBaseURL(m.BaseURL))
			}
			client, err = coord.NewLLMClient(context.Background(), "openai", options...)
			if err != nil {
				log.Fatal().Err(err).Int("index", i).Msg("failed to create openai client")
			}
		case "vertexai":
			if m.Project != "" {
				options = append(options, pconf.WithProjectID(m.Project))
			}
			if m.Location != "" {
				options = append(options, pconf.WithLocation(m.Location))
			}
			client, err = coord.NewLLMClient(context.Background(), "vertexai", options...)
			if err != nil {
				log.Fatal().Err(err).Int("index", i).Msg("failed to create vertexai client")
			}
		case "aistudio":
			if m.APIKey != "" {
				options = append(options, pconf.WithAPIKey(m.APIKey))
			}
			client, err = coord.NewLLMClient(context.Background(), "aistudio", options...)
			if err != nil {
				log.Fatal().Err(err).Int("index", i).Msg("failed to create aistudio client")
			}
		case "anthropic":
			if m.APIKey != "" {
				options = append(options, pconf.WithAPIKey(m.APIKey))
			}
			client, err = coord.NewLLMClient(context.Background(), "anthropic", options...)
			if err != nil {
				log.Fatal().Err(err).Int("index", i).Msg("failed to create anthropic client")
			}
		}

		output_tokens := 8192
		if m.MaxTokens > 0 {
			output_tokens = m.MaxTokens
		}

		model, err := client.NewLLM(m.ModelID, &llm.Config{
			Temperature:           m.Temperature,
			MaxOutputTokens:       Ptr(output_tokens),
			SafetyFilterThreshold: llm.BlockOff,
		})
		if err != nil {
			log.Fatal().Err(err).Int("index", i).Msg("failed to create model")
		}

		if m.RateLimit != nil {
			model = NewRateLimitingModel(model, rate.Limit(*m.RateLimit))
		}
		models = append(models, model)
	}
	model := NewLoadBalancingModel(models...)
	translationModel = model

	startIndex = c.StartIndex
	if c.CustomPrompt != nil {
		customPrompt = *c.CustomPrompt
	}
}
