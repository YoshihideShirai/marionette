package marionette

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
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

func TestFlashPersistsForNextRequestAndAutoClears(t *testing.T) {
	app := New()
	app.Action("save", func(ctx *Context) Node {
		ctx.FlashSuccess("saved")
		return Div("app", Text("ok"))
	})
	app.Page("/", func(ctx *Context) Node {
		return FlashAlerts(ctx.Flashes())
	})

	actionReq := httptest.NewRequest(http.MethodPost, "/save", strings.NewReader(""))
	actionReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	actionRes := httptest.NewRecorder()
	app.Handler().ServeHTTP(actionRes, actionReq)

	var flashCookie *http.Cookie
	for _, c := range actionRes.Result().Cookies() {
		if c.Name == flashCookieName {
			flashCookie = c
			break
		}
	}
	if flashCookie == nil {
		t.Fatalf("expected %q cookie to be set", flashCookieName)
	}

	firstGet := httptest.NewRequest(http.MethodGet, "/", nil)
	firstGet.AddCookie(flashCookie)
	firstRes := httptest.NewRecorder()
	app.Handler().ServeHTTP(firstRes, firstGet)
	if !strings.Contains(firstRes.Body.String(), "saved") {
		t.Fatalf("expected flash message in first GET response, got %q", firstRes.Body.String())
	}

	secondGet := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, c := range firstRes.Result().Cookies() {
		if c.Name == flashCookieName {
			secondGet.AddCookie(c)
		}
	}
	secondRes := httptest.NewRecorder()
	app.Handler().ServeHTTP(secondRes, secondGet)
	if strings.Contains(secondRes.Body.String(), "saved") {
		t.Fatalf("expected flash message to be auto-cleared, got %q", secondRes.Body.String())
	}
}

func TestAppStateConcurrentSetViaContext(t *testing.T) {
	app := New()
	app.Action("set", func(ctx *Context) Node {
		ctx.Set("name", ctx.FormValue("name"))
		return Div("ok")
	})

	handler := app.Handler()
	const workers = 32

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		i := i
		go func() {
			defer wg.Done()
			form := url.Values{"name": {strconv.Itoa(i)}}
			req := httptest.NewRequest(http.MethodPost, "/set", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Errorf("expected 200, got %d", rr.Code)
			}
		}()
	}
	wg.Wait()
}

func TestContextSetIsVisibleFromAppGetInt(t *testing.T) {
	app := New()
	app.Page("/", func(ctx *Context) Node {
		ctx.Set("count", 7)
		return Div("app")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if got := app.GetInt("count"); got != 7 {
		t.Fatalf("expected shared app/context state to be 7, got %d", got)
	}
}
