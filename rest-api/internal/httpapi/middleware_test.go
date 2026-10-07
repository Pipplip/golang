package httpapi

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoveryMiddlewareReturns500OnPanic(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))

	panicHandler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})

	handler := RecoveryMiddleware(logger)(panicHandler)
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "internal server error") {
		t.Fatalf("expected error response body, got %q", rec.Body.String())
	}
}

func TestLoggingMiddlewareLogsRequest(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler := LoggingMiddleware(logger)(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}

	logText, _ := io.ReadAll(&logs)
	logLine := string(logText)
	if !strings.Contains(logLine, "method=GET") {
		t.Fatalf("expected method in log line: %q", logLine)
	}
	if !strings.Contains(logLine, "path=/books") {
		t.Fatalf("expected path in log line: %q", logLine)
	}
	if !strings.Contains(logLine, "status=204") {
		t.Fatalf("expected status in log line: %q", logLine)
	}
}
