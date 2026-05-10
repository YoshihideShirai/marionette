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
		"server-driven state, htmx partial updates, and token-by-token SSE rendering",
		"Demo conversation",
		"Marionette AI",
		`id="chat-panel"`,
		`hx-post="/chat/send"`,
		`hx-target="#chat-panel"`,
		"SSE mock replies append token by token",
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
		"SSE streaming",
		`id="message-content-3"`,
		`id="message-cursor-3"`,
		`data-marionette-sse-url="/chat/stream"`,
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

	body := getStream(t, handler)
	for _, want := range []string{
		"event: html",
		"event: done",
		"In",
		"Marionette",
		"beforeend:#message-content-3",
		"Complete",
		" StreamAction",
		" SSE",
		" endpoint",
		" token",
		" client-side",
		" state.",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected SSE stream to contain %q, got %q", want, body)
		}
	}
	if strings.Contains(body, `hx-post="/chat/stream"`) {
		t.Fatalf("expected SSE stream not to use polling trigger, got %q", body)
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
	if strings.Contains(body, "quarterly forecast") || strings.Contains(body, "Streaming response") || strings.Contains(body, `data-marionette-sse-url="/chat/stream"`) {
		t.Fatalf("expected reset fragment to remove previous conversation and stream trigger, got %q", body)
	}
}

func getStream(t *testing.T, handler http.Handler) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/chat/stream", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected stream status 200, got %d", rr.Code)
	}
	return rr.Body.String()
}
