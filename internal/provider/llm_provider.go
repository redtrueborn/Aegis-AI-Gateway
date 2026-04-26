package provider

import (
	"aegis-ai-gateway/internal/domain"
	"context"
)


type LLMProvider interface {
	Generate(ctx context.Context, request domain.GenerateRequest)  (response domain.GenerateResponse, error error)
}