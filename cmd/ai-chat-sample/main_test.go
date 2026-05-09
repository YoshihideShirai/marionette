package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootPageRendersAIChatSample(t *testing.T) {
	app := buildApp()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{
		"AI Chat Sample",
		"server-driven state, htmx partial updates, and simulated streaming replies",
		"Demo conversation",
		"Marionette AI",
		`id="chat-panel"`,
		`hx-post="/chat/send"`,
		`hx-target="#chat-panel"`,
		"Streaming mock replies render chunk by chunk",
		"Reset conversation",
		"Integration note",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected root page to contain %q, got %q", want, body)
		}
	}
}

func TestChatSendActionAppendsUserAndStartsStreamingAssistantMessage(t *testing.T) {
	app := buildApp()
	req := httptest.NewRequest(http.MethodPost, "/chat/send", strings.NewReader("prompt=Explain+htmx+streaming"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{
		`id="chat-panel"`,
		"You",
		"Explain htmx streaming",
		"Streaming response…",
		"Streaming",
		`hx-post="/chat/stream"`,
		`hx-trigger="load delay:350ms"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected send fragment to contain %q, got %q", want, body)
		}
	}
	if strings.Contains(body, "<!doctype html>") {
		t.Fatalf("expected action to return a fragment, got full document %q", body)
	}
}

func TestChatStreamActionProgressivelyRevealsAssistantMessage(t *testing.T) {
	app := buildApp()
	handler := app.Handler()

	sendReq := httptest.NewRequest(http.MethodPost, "/chat/send", strings.NewReader("prompt=Explain+htmx+streaming"))
	sendReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler.ServeHTTP(httptest.NewRecorder(), sendReq)

	first := postStream(t, handler)
	for _, want := range []string{
		"In Marionette, ActionForm Target and",
		"▌",
		`hx-post="/chat/stream"`,
	} {
		if !strings.Contains(first, want) {
			t.Fatalf("expected first stream fragment to contain %q, got %q", want, first)
		}
	}

	latest := first
	for i := 0; i < 12 && strings.Contains(latest, `hx-post="/chat/stream"`); i++ {
		latest = postStream(t, handler)
	}
	for _, want := range []string{
		"This sample combines those partial updates with a polling trigger",
		"assistant reply appears chunk by chunk.",
	} {
		if !strings.Contains(latest, want) {
			t.Fatalf("expected final stream fragment to contain %q, got %q", want, latest)
		}
	}
	if strings.Contains(latest, "▌") || strings.Contains(latest, `hx-post="/chat/stream"`) {
		t.Fatalf("expected completed stream without cursor or stream trigger, got %q", latest)
	}
}

func TestChatSendActionShowsValidationError(t *testing.T) {
	app := buildApp()
	req := httptest.NewRequest(http.MethodPost, "/chat/send", strings.NewReader("prompt=+"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{"Input error", "Please enter a message."} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected validation fragment to contain %q, got %q", want, body)
		}
	}
}

func TestChatResetActionRestoresWelcomeMessage(t *testing.T) {
	app := buildApp()
	handler := app.Handler()

	sendReq := httptest.NewRequest(http.MethodPost, "/chat/send", strings.NewReader("prompt=quarterly+forecast"))
	sendReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler.ServeHTTP(httptest.NewRecorder(), sendReq)

	resetReq := httptest.NewRequest(http.MethodPost, "/chat/reset", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, resetReq)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Hello. This is an AI chat sample demo.") {
		t.Fatalf("expected reset fragment to contain welcome message, got %q", body)
	}
	if strings.Contains(body, "quarterly forecast") || strings.Contains(body, "Streaming response") || strings.Contains(body, `hx-post="/chat/stream"`) {
		t.Fatalf("expected reset fragment to remove previous conversation and stream trigger, got %q", body)
	}
}

func postStream(t *testing.T, handler http.Handler) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/chat/stream", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected stream status 200, got %d", rr.Code)
	}
	return rr.Body.String()
}
