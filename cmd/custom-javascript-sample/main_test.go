package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootPageIncludesCustomAssetsAndFormulaPanel(t *testing.T) {
	app := buildApp()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Custom JavaScript Sample",
		"External MathJax plus app-level custom JavaScript.",
		mathJaxCHTMLURL,
		"custom-js-status",
		"Waiting for custom JavaScript",
		"htmx:afterSwap",
		".math-sample-equation",
		"Euler identity",
		`data-mathjax`,
		`hx-post="/formula/next"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected root page to contain %q, got %q", want, body)
		}
	}
}

func TestNextFormulaActionCyclesFormulaPanelFragment(t *testing.T) {
	app := buildApp()

	req := httptest.NewRequest(http.MethodPost, "/formula/next", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{"Quadratic formula", `data-mathjax`, `Next formula`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected formula fragment to contain %q, got %q", want, body)
		}
	}
	if strings.Contains(body, "<!doctype html>") {
		t.Fatalf("expected action to return a fragment, got full document %q", body)
	}
}
