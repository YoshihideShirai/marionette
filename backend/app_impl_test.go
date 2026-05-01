package backend

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	frontend "github.com/YoshihideShirai/marionette/frontend"
)

func TestPageIncludesCustomStyles(t *testing.T) {
	app := New()
	app.AddStylesheet("/assets/app.css")
	app.AddStyle(`
		#marionette-root {
			max-width: 48rem;
		}
	`)
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Dashboard"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, `<link href="/assets/app.css" rel="stylesheet" type="text/css" />`) {
		t.Fatalf("expected custom stylesheet link, got %q", body)
	}
	if !strings.Contains(body, `<style>#marionette-root`) {
		t.Fatalf("expected custom inline CSS, got %q", body)
	}
	defaultIndex := strings.Index(body, "--mrn-page-max-width")
	customIndex := strings.Index(body, "#marionette-root {\n\t\t\tmax-width: 48rem;")
	if defaultIndex == -1 {
		t.Fatalf("expected default Marionette CSS in response, got %q", body)
	}
	if customIndex == -1 {
		t.Fatalf("expected custom CSS in response, got %q", body)
	}
	if customIndex < defaultIndex {
		t.Fatalf("expected custom CSS after default CSS, got %q", body)
	}
}

func TestPageCanSetHTMLTitle(t *testing.T) {
	app := New()
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Dashboard"))
	}, WithTitle(`Users & Teams`))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, `<title>Users &amp; Teams</title>`) {
		t.Fatalf("expected escaped custom title, got %q", body)
	}
}

func TestPageIncludesCustomScripts(t *testing.T) {
	app := New()
	app.AddScript("https://cdn.example.com/widget.js")
	app.AddJavaScript(`
		window.marionetteWidgetReady = true;
	`)
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Dashboard"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	scriptIndex := strings.Index(body, `<script src="https://cdn.example.com/widget.js"></script>`)
	inlineIndex := strings.Index(body, `<script>window.marionetteWidgetReady = true;</script>`)
	if scriptIndex == -1 {
		t.Fatalf("expected custom external script, got %q", body)
	}
	if inlineIndex == -1 {
		t.Fatalf("expected custom inline JavaScript, got %q", body)
	}
	if inlineIndex < scriptIndex {
		t.Fatalf("expected inline JavaScript after external scripts, got %q", body)
	}
}
