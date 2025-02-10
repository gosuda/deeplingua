package main

import (
	"context"
	"sync/atomic"

	"github.com/lemon-mint/coord/llm"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

type LoadBalancingModel struct {
	models []llm.Model
	idx    atomic.Int64
}

func NewLoadBalancingModel(models ...llm.Model) *LoadBalancingModel {
	return &LoadBalancingModel{
		models: models,
	}
}

func (g *LoadBalancingModel) GenerateStream(ctx context.Context, chat *llm.ChatContext, input *llm.Content) *llm.StreamContent {
	idx := g.idx.Add(1) % int64(len(g.models))
	return g.models[idx].GenerateStream(ctx, chat, input)
}

func (g *LoadBalancingModel) Close() error {
	return nil
}

func (g *LoadBalancingModel) Name() string {
	return "LoadBalancingModel"
}

type RateLimitingModel struct {
	llm.Model
	rlim *rate.Limiter
}

func NewRateLimitingModel(model llm.Model, tps rate.Limit) *RateLimitingModel {
	return &RateLimitingModel{
		Model: model,
		rlim:  rate.NewLimiter(tps, 0),
	}
}

func (g *RateLimitingModel) GenerateStream(ctx context.Context, chat *llm.ChatContext, input *llm.Content) *llm.StreamContent {
	g.rlim.Wait(ctx)
	log.Info().Msg("calling llm model")
	return g.Model.GenerateStream(ctx, chat, input)
}

func Ptr[T any](v T) *T {
	return &v
}
