package server

import (
	"aegis-ai-gateway/internal/domain"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type stubProvider struct {
	response domain.GenerateResponse
	err      error
	called   bool
	req      domain.GenerateRequest
}

func (s *stubProvider) Generate(_ context.Context, req domain.GenerateRequest) (domain.GenerateResponse, error) {
	s.called = true
	s.req = req
	if s.err != nil {
		return domain.GenerateResponse{}, s.err
	}
	return s.response, nil
}

func TestChatHandler_MethodNotAllowed(t *testing.T) {
	provider := &stubProvider{}
	h := &ChatHandler{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		provider:       provider,
		requestTimeout: time.Second,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/chat", nil)
	rr := httptest.NewRecorder()

	h.ServeHttp(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
	if provider.called {
		t.Fatal("provider Generate should not be called for invalid method")
	}
}

func TestChatHandler_EmptyBody(t *testing.T) {
	provider := &stubProvider{}
	h := &ChatHandler{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		provider:       provider,
		requestTimeout: time.Second,
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/chat", nil)
	req.Body = nil
	rr := httptest.NewRecorder()

	h.ServeHttp(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
	if provider.called {
		t.Fatal("provider Generate should not be called for empty body")
	}
}

func TestChatHandler_InvalidJSON(t *testing.T) {
	provider := &stubProvider{}
	h := &ChatHandler{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		provider:       provider,
		requestTimeout: time.Second,
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/chat", strings.NewReader(`{"model":`))
	rr := httptest.NewRecorder()

	h.ServeHttp(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
	if provider.called {
		t.Fatal("provider Generate should not be called for invalid JSON")
	}
}

func TestChatHandler_ProviderError(t *testing.T) {
	provider := &stubProvider{err: errors.New("provider failed")}
	h := &ChatHandler{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		provider:       provider,
		requestTimeout: time.Second,
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/chat", strings.NewReader(`{"model":"gpt-4o","message":{"content":"hello","role":"user"}}`))
	rr := httptest.NewRecorder()

	h.ServeHttp(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}
	if !provider.called {
		t.Fatal("provider Generate should be called on valid request")
	}
}

func TestChatHandler_Success(t *testing.T) {
	provider := &stubProvider{
		response: domain.GenerateResponse{
			Provider: domain.ProviderMock,
			Model:    "gpt-4o-mini",
			Output: map[string]any{
				"text": "hello back",
			},
		},
	}
	h := &ChatHandler{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		provider:       provider,
		requestTimeout: time.Second,
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/chat", strings.NewReader(`{"model":"gpt-4o","message":{"content":"hello","role":"user"}}`))
	rr := httptest.NewRecorder()

	h.ServeHttp(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", got, "application/json")
	}
	if !provider.called {
		t.Fatal("provider Generate should be called for valid request")
	}
	if provider.req.Model != "gpt-4o" {
		t.Fatalf("provider request model = %q, want %q", provider.req.Model, "gpt-4o")
	}
	if len(provider.req.Messages) != 1 {
		t.Fatalf("provider request messages length = %d, want %d", len(provider.req.Messages), 1)
	}
	if provider.req.Messages[0].Content != "hello" {
		t.Fatalf("provider request message content = %q, want %q", provider.req.Messages[0].Content, "hello")
	}
	if !strings.Contains(rr.Body.String(), `"status":"completed"`) {
		t.Fatalf("response body missing status field: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"provider":"mock"`) {
		t.Fatalf("response body missing provider field: %s", rr.Body.String())
	}
}
