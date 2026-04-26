package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewServer_RegistersHealthRoute(t *testing.T) {
	s := NewServer(Config{Addr: ":8080"}, slog.New(slog.NewTextHandler(io.Discard, nil)))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	s.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestHttpServer_ConfiguresFields(t *testing.T) {
	s := NewServer(Config{Addr: ":9090"}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	h := s.HttpServer()

	if h.Addr != ":9090" {
		t.Fatalf("Addr = %q, want %q", h.Addr, ":9090")
	}
	if h.Handler == nil {
		t.Fatal("Handler is nil, want non-nil")
	}
	if h.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("ReadHeaderTimeout = %s, want %s", h.ReadHeaderTimeout, 5*time.Second)
	}
}

func TestHeathHandler_GetReturnsExpectedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	HeathHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", got, "application/json")
	}

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body as JSON: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("status field = %#v, want %q", body["status"], "ok")
	}
	if body["service"] != "aegis-ai-gateway" {
		t.Fatalf("service field = %#v, want %q", body["service"], "aegis-ai-gateway")
	}
}

func TestHeathHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	rr := httptest.NewRecorder()

	HeathHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}
