package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDashwindDemoRootSmoke(t *testing.T) {
	app := buildApp()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{
		`id="dashwind-main"`,
		"Dashboard",
		"Revenue",
		"Leads",
		`hx-post="/dashboard/period"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected dashboard root to contain %q, got %q", want, body)
		}
	}
}
