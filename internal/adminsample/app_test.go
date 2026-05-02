package adminsample

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestOrdersFilterFragmentKeepsMainContentTarget(t *testing.T) {
	app := BuildApp()
	app.Set("loggedIn", true)
	app.Set("currentPage", "overview")
	handler := app.Handler()

	first := postFilter(t, handler, "Active")
	if strings.Contains(first, "<!doctype html>") {
		t.Fatalf("expected fragment response, got full document")
	}
	for _, want := range []string{`id="main-content"`, `Filter applied: Active`, `selected="selected" value="Active"`} {
		if !strings.Contains(first, want) {
			t.Fatalf("expected first filter response to contain %q, got %q", want, first)
		}
	}

	second := postFilter(t, handler, "Review")
	for _, want := range []string{`id="main-content"`, `Filter applied: Review`, `selected="selected" value="Review"`} {
		if !strings.Contains(second, want) {
			t.Fatalf("expected second filter response to contain %q, got %q", want, second)
		}
	}
}

func postFilter(t *testing.T, handler http.Handler, status string) string {
	t.Helper()
	form := url.Values{"status": {status}}
	req := httptest.NewRequest(http.MethodPost, "/orders/filter", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	return rr.Body.String()
}
