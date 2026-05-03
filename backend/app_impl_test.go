package backend

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

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

func TestPageCanSwitchStyleTemplateImports(t *testing.T) {
	app := New()
	app.UseStyleTemplate(frontend.StyleTemplate{
		Name:                 "tailadmin-custom",
		FrameworkStylesheets: []string{"https://cdn.example.com/tailadmin.css"},
		FrameworkScripts:     []string{"https://cdn.example.com/tailwind.js"},
	})
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Dashboard"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, `href="https://cdn.example.com/tailadmin.css"`) {
		t.Fatalf("expected template stylesheet import, got %q", body)
	}
	if !strings.Contains(body, `src="https://cdn.example.com/tailwind.js"`) {
		t.Fatalf("expected template script import, got %q", body)
	}
	if strings.Contains(body, `cdn.jsdelivr.net/npm/daisyui@5`) {
		t.Fatalf("expected default framework import to be replaced, got %q", body)
	}
}

func TestPageCanUseTailwindTemplatePreset(t *testing.T) {
	app := New()
	app.UseTailwindCSSTemplate()
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Dashboard"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, `src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"`) {
		t.Fatalf("expected tailwind framework script import, got %q", body)
	}
	if strings.Contains(body, `cdn.jsdelivr.net/npm/daisyui@5`) {
		t.Fatalf("did not expect daisyui import for tailwind template, got %q", body)
	}
}

func TestUseStyleTemplateByNameRejectsUnknown(t *testing.T) {
	app := New()
	err := app.UseStyleTemplateByName("unknown-template")
	if err == nil {
		t.Fatal("expected error for unknown template")
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

func TestAssetsServeFilesAndApplyHeaders(t *testing.T) {
	app := New()
	app.Assets("/assets", fstest.MapFS{
		"app.css": {Data: []byte("body { color: red; }")},
	}, WithAssetCache(time.Hour), WithAssetImmutable())

	req := httptest.NewRequest(http.MethodGet, "/assets/app.css", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if got := strings.TrimSpace(rr.Body.String()); got != "body { color: red; }" {
		t.Fatalf("expected asset body, got %q", got)
	}
	if got := rr.Header().Get("Cache-Control"); got != "public, max-age=3600, immutable" {
		t.Fatalf("expected cache header, got %q", got)
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/css") {
		t.Fatalf("expected css content type, got %q", got)
	}
}

func TestAssetsCanServeDownloads(t *testing.T) {
	app := New()
	app.Assets("/assets", fstest.MapFS{
		"reports/users report.csv": {Data: []byte("name\nAiko\n")},
	}, WithAssetDownload(), WithAssetCache(time.Hour), WithAssetContentTypes(map[string]string{".csv": "text/csv; charset=utf-8"}))

	req := httptest.NewRequest(http.MethodGet, app.Asset("reports/users report.csv"), nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if got := rr.Body.String(); got != "name\nAiko\n" {
		t.Fatalf("expected download body, got %q", got)
	}
	if got := rr.Header().Get("Content-Disposition"); got != `attachment; filename="users report.csv"` {
		t.Fatalf("expected content disposition, got %q", got)
	}
	if got := rr.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("expected csv content type, got %q", got)
	}
	if got := rr.Header().Get("Cache-Control"); got != "public, max-age=3600" {
		t.Fatalf("expected cache header, got %q", got)
	}
}

func TestDownloadsConvenienceRouteServesHeadRequests(t *testing.T) {
	app := New()
	app.Downloads("/downloads", fstest.MapFS{
		"report.csv": {Data: []byte("name\nAiko\n")},
	})

	req := httptest.NewRequest(http.MethodHead, "/downloads/report.csv", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if got := rr.Body.String(); got != "" {
		t.Fatalf("expected empty body for HEAD, got %q", got)
	}
	if got := rr.Header().Get("Content-Disposition"); got != `attachment; filename=report.csv` {
		t.Fatalf("expected content disposition, got %q", got)
	}
}

func TestAssetsBlockDirectoryIndexByDefault(t *testing.T) {
	app := New()
	app.Assets("/assets", fstest.MapFS{
		"icons":           {Mode: 0o755 | fs.ModeDir},
		"icons/check.svg": {Data: []byte("<svg></svg>")},
	})

	req := httptest.NewRequest(http.MethodGet, "/assets/icons/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for directory index, got %d", rr.Code)
	}
}

func TestAssetBuildsURLFromFirstRegisteredAssetPrefix(t *testing.T) {
	app := New()
	app.Assets("assets", fstest.MapFS{})

	if got := app.Asset("images/hero image.png"); got != "/assets/images/hero%20image.png" {
		t.Fatalf("expected escaped asset URL, got %q", got)
	}
	if got := app.Asset("https://cdn.example.com/app.css"); got != "https://cdn.example.com/app.css" {
		t.Fatalf("expected absolute asset URL to pass through, got %q", got)
	}
}

func TestAssetsServeEscapedAssetPaths(t *testing.T) {
	app := New()
	app.Assets("/assets", fstest.MapFS{
		"images/hero image.png": {Data: []byte("png")},
	})

	req := httptest.NewRequest(http.MethodGet, app.Asset("images/hero image.png"), nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if got := rr.Body.String(); got != "png" {
		t.Fatalf("expected escaped asset body, got %q", got)
	}
}
