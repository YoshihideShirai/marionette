package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDesktopOptionsDoNotRequireWebView(t *testing.T) {
	opts := desktopOptions()
	if opts.Title != desktopTitle {
		t.Fatalf("expected title %q, got %q", desktopTitle, opts.Title)
	}
	if opts.Width != 1200 || opts.Height != 800 {
		t.Fatalf("expected 1200x800 window, got %dx%d", opts.Width, opts.Height)
	}
}

func TestDesktopAppConstructionRendersRootPage(t *testing.T) {
	app := buildApp(desktopTitle, desktopDescription)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{
		desktopTitle,
		desktopDescription,
		"No tasks yet",
		`hx-post="/tasks/create"`,
		`id="task-list"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected desktop root to contain %q, got %q", want, body)
		}
	}
}
