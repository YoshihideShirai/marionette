package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDataAppRootSmoke(t *testing.T) {
	app := buildApp()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Data App Tutorial",
		"CSV load",
		`id="sales-dashboard"`,
		"Regional sales table",
		"Regional sales chart",
		"East",
		"West",
		"North",
		"Sales by region",
		`hx-post="/filters/apply"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected root page to contain %q, got %q", want, body)
		}
	}
	if strings.Contains(body, "South") {
		t.Fatalf("expected default sales filter to hide South, got %q", body)
	}
}

func TestFilterActionRendersDashboardFragment(t *testing.T) {
	app := buildApp()

	req := httptest.NewRequest(http.MethodPost, "/filters/apply", strings.NewReader("min_sales=1300"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{`id="sales-dashboard"`, "Regional sales table", "East", "West", "Sales by region"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected filter fragment to contain %q, got %q", want, body)
		}
	}
	if strings.Contains(body, "North") || strings.Contains(body, "South") {
		t.Fatalf("expected min sales filter to hide lower sales regions, got %q", body)
	}
	if strings.Contains(body, "<!doctype html>") {
		t.Fatalf("expected action to return a fragment, got full document %q", body)
	}
}
