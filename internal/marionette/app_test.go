package marionette

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexRendersHTMLNotEscaped(t *testing.T) {
	app := New()
	app.Render(func(ctx *Context) Node {
		return DivClass("app", "card", Text("Hello"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	if strings.Contains(body, "&lt;div") {
		t.Fatalf("expected unescaped HTML, got %q", body)
	}
	if !strings.Contains(body, `<div class="card" id="app">`) {
		t.Fatalf("expected rendered card fragment, got %q", body)
	}
}
