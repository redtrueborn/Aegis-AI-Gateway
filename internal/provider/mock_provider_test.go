package provider

import (
	"context"
	"errors"
	"testing"
	"time"

	"aegis-ai-gateway/internal/domain"
)

var _ LLMProvider = (*MockProvider)(nil)

func TestNewMockProvider_SetsFields(t *testing.T) {
	p := NewMockProvider("mock", 25*time.Millisecond)

	if p == nil {
		t.Fatal("expected provider, got nil")
	}
	if p.name != "mock" {
		t.Fatalf("name = %q, want %q", p.name, "mock")
	}
	if p.delay != 25*time.Millisecond {
		t.Fatalf("delay = %s, want %s", p.delay, 25*time.Millisecond)
	}
}

func TestMockProvider_Generate_ReturnsCanceledContextError(t *testing.T) {
	p := NewMockProvider("mock", 10*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := p.Generate(ctx, domain.GenerateRequest{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want %v", err, context.Canceled)
	}
}

func TestMockProvider_Generate_ReturnsDeadlineExceeded(t *testing.T) {
	p := NewMockProvider("mock", 20*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	_, err := p.Generate(ctx, domain.GenerateRequest{})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want %v", err, context.DeadlineExceeded)
	}
}
