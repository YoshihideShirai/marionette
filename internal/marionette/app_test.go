package marionette

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestPageRendersFullHTML(t *testing.T) {
	app := New()
	app.Page("/", func(ctx *Context) Node {
		return DivClass("app", "card", Text("Dashboard"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "<!doctype html>") {
		t.Fatalf("expected full html shell, got %q", body)
	}
	if !strings.Contains(body, "Dashboard") {
		t.Fatalf("expected page content, got %q", body)
	}
}

func TestUnregisteredPageReturnsNotFound(t *testing.T) {
	app := New()
	app.Page("/", func(ctx *Context) Node {
		return Div("app", Text("Dashboard"))
	})

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestActionRendersFragmentAndParsesForm(t *testing.T) {
	app := New()
	app.Action("users/create", func(ctx *Context) Node {
		ctx.Set("name", ctx.FormValue("name"))
		return DivClass("users", "card", Text(ctx.Get("name").(string)))
	})

	form := url.Values{"name": {"Aiko"}}
	req := httptest.NewRequest(http.MethodPost, "/users/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if strings.Contains(body, "<!doctype html>") {
		t.Fatalf("expected fragment response, got %q", body)
	}
	if !strings.Contains(body, `id="users"`) || !strings.Contains(body, "Aiko") {
		t.Fatalf("expected rendered fragment with form value, got %q", body)
	}
}

func TestActionRejectsGet(t *testing.T) {
	app := New()
	app.Action("users/create", func(ctx *Context) Node {
		return Div("users")
	})

	req := httptest.NewRequest(http.MethodGet, "/users/create", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}

func TestContextQueryAndStateHelpers(t *testing.T) {
	app := New()
	app.Page("/", func(ctx *Context) Node {
		ctx.Set("count", 2)
		return Div("app", Text(ctx.Query("filter")+":"+ctx.Param("id")+":"+strings.Repeat("x", ctx.GetInt("count"))))
	})

	req := httptest.NewRequest(http.MethodGet, "/?filter=active", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if !strings.Contains(rr.Body.String(), "active::xx") {
		t.Fatalf("expected query and state helper output, got %q", rr.Body.String())
	}
}
