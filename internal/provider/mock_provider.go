package provider

import (
	"aegis-ai-gateway/internal/domain"
	"context"
	"time"
)


type MockProvider struct {
	name string
	delay time.Duration
}

func NewMockProvider(name string, delay time.Duration) *MockProvider {
	return &MockProvider{
		name:  name,
		delay: delay,
	}
}

func (p *MockProvider) Generate(ctx context.Context, request domain.GenerateRequest) (response domain.GenerateResponse, error error) {
	for {
	select{  
		case <- time.After(p.delay):
		  continue
				case <- ctx.Done():
				return domain.GenerateResponse{}, ctx.Err()
	}
	}
}