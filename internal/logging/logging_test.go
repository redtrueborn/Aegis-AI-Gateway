package logging

import (
	"io"
	"log/slog"
	"testing"
)

func TestInitLogger_ReturnsDefaultLogger(t *testing.T) {
	logger := InitLogger()
	if logger == nil {
		t.Fatal("expected logger, got nil")
	}

	if logger != slog.Default() {
		t.Fatal("expected InitLogger to return slog.Default logger")
	}
}

func TestInitLogger_UsesCurrentGlobalDefault(t *testing.T) {
	originalDefault := slog.Default()
	t.Cleanup(func() {
		slog.SetDefault(originalDefault)
	})

	customDefault := slog.New(slog.NewTextHandler(io.Discard, nil))
	slog.SetDefault(customDefault)

	logger := InitLogger()
	if logger != customDefault {
		t.Fatal("expected InitLogger to return current global default logger")
	}
}
