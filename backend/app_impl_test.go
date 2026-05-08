package backend

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	frontend "github.com/YoshihideShirai/marionette/frontend"
	"github.com/YoshihideShirai/marionette/frontend/assets"
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
	if strings.Contains(body, assets.DaisyUICSSURL) {
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
	if !strings.Contains(body, `src="`+assets.TailwindBrowserURL+`"`) {
		t.Fatalf("expected tailwind framework script import, got %q", body)
	}
	if strings.Contains(body, assets.DaisyUICSSURL) {
		t.Fatalf("did not expect daisyui import for tailwind template, got %q", body)
	}
}

func TestAppUseAssetsOverridesBuiltInAssetProvider(t *testing.T) {
	app := New()
	app.UseAssets(assets.NewLocalAssetProvider("/vendor"))
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Dashboard"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	for _, want := range []string{
		`href="/vendor/daisyui.css"`,
		`src="/vendor/tailwindcss-browser.js"`,
		`src="/vendor/htmx.min.js"`,
		`src="/vendor/chart.umd.js"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected provider-resolved asset %q, got %q", want, body)
		}
	}
	if strings.Contains(body, assets.DaisyUICSSURL) || strings.Contains(body, assets.HTMXURL) {
		t.Fatalf("did not expect CDN assets when provider is overridden, got %q", body)
	}
}

func TestUseOfflineAssetsKeepsGeneratedHTMLLocal(t *testing.T) {
	app := New()
	app.Assets("/vendor", fstest.MapFS{
		assets.DaisyUICSSFile:           {Data: []byte("/* daisyui */")},
		assets.TailwindCSSBrowserJSFile: {Data: []byte("// tailwind")},
		assets.HTMXJSFile:               {Data: []byte("// htmx")},
		assets.ChartJSFile:              {Data: []byte("// chart")},
	}, WithAssetCache(time.Hour))
	app.UseOfflineAssets("/vendor")
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Offline"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	if strings.Contains(body, "http://") || strings.Contains(body, "https://") {
		t.Fatalf("expected offline HTML without absolute network URLs, got %q", body)
	}
	for _, want := range []string{
		`href="/vendor/daisyui.css"`,
		`src="/vendor/tailwindcss-browser.js"`,
		`src="/vendor/htmx.min.js"`,
		`src="/vendor/chart.umd.js"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected offline asset %q, got %q", want, body)
		}
	}
}

func TestUseStyleTemplateByNameRejectsUnknown(t *testing.T) {
	app := New()
	err := app.UseStyleTemplateByName("unknown-template")
	if err == nil {
		t.Fatal("expected error for unknown template")
	}
}

func TestPageCanDisableChartsWhenNoChartComponentsAreUsed(t *testing.T) {
	app := New()
	app.DisableCharts()
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Dashboard"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	if strings.Contains(body, assets.ChartJSURL) || strings.Contains(body, "chart.umd.js") {
		t.Fatalf("did not expect Chart.js asset on a page without charts, got %q", body)
	}
	if strings.Contains(body, "window.mrnInitCharts") {
		t.Fatalf("did not expect chart bootstrap on a page without charts, got %q", body)
	}
	if !strings.Contains(body, assets.HTMXURL) {
		t.Fatalf("expected HTMX to remain enabled by default, got %q", body)
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

func TestAppGlobalStateHelpers(t *testing.T) {
	app := New()
	app.SetGlobal("count", 3)
	app.SetGlobal("name", "Aiko")

	if got := app.GetGlobalInt("count"); got != 3 {
		t.Fatalf("expected global int helper to return 3, got %d", got)
	}
	if got := app.GetGlobal("name"); got != "Aiko" {
		t.Fatalf("expected global value %q, got %v", "Aiko", got)
	}
	if got := app.GetGlobalInt("name"); got != 0 {
		t.Fatalf("expected non-int global value to return 0, got %d", got)
	}
}

func TestGetGlobalSnapshotClonesWhileLocked(t *testing.T) {
	app := New()
	app.SetGlobal("names", []string{"Aiko", "Ren"})

	snapshot, ok := app.GetGlobalSnapshot("names", func(value any) any {
		names, _ := value.([]string)
		return append([]string(nil), names...)
	}).([]string)
	if !ok {
		t.Fatalf("expected cloned names snapshot to be []string")
	}
	snapshot[0] = "changed"

	stored := app.GetGlobalSnapshot("names", func(value any) any {
		names, _ := value.([]string)
		return append([]string(nil), names...)
	}).([]string)
	if got := stored[0]; got != "Aiko" {
		t.Fatalf("expected snapshot mutation not to change global state, got %q", got)
	}

	ctx := &Context{}
	ctx.SetGlobal("labels", map[string]string{"role": "Admin"})
	labels := ctx.GetGlobalSnapshot("labels", func(value any) any {
		old, _ := value.(map[string]string)
		next := make(map[string]string, len(old))
		for key, value := range old {
			next[key] = value
		}
		return next
	}).(map[string]string)
	labels["role"] = "Viewer"
	if got := ctx.GetGlobal("labels").(map[string]string)["role"]; got != "Admin" {
		t.Fatalf("expected fallback snapshot mutation not to change global state, got %q", got)
	}
}

func TestUpdateGlobalHoldsLockAcrossTransform(t *testing.T) {
	app := New()
	app.SetGlobal("count", 0)

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				app.UpdateGlobal("count", func(old any) any {
					oldInt, _ := old.(int)
					return oldInt + 1
				})
			}
		}()
	}
	wg.Wait()

	if got := app.GetGlobalInt("count"); got != 10000 {
		t.Fatalf("expected atomic global update to reach 10000, got %d", got)
	}
	if got := app.IncrementGlobalInt("count", 5); got != 10005 {
		t.Fatalf("expected increment helper to return 10005, got %d", got)
	}
}

func TestContextUpdateGlobalDelegatesToAppAndSupportsFallbackState(t *testing.T) {
	app := New()
	app.SetGlobal("count", 1)
	app.Page("/", func(ctx *Context) frontend.Node {
		next := ctx.UpdateGlobal("count", func(old any) any {
			return old.(int) + 2
		})
		return frontend.DivProps(frontend.ElementProps{ID: "app"}, frontend.Text(strconv.Itoa(next.(int))))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if !strings.Contains(rr.Body.String(), ">3<") {
		t.Fatalf("expected delegated context update output, got %q", rr.Body.String())
	}
	if got := app.GetGlobalInt("count"); got != 3 {
		t.Fatalf("expected app state to be updated, got %d", got)
	}

	ctx := &Context{}
	if got := ctx.IncrementGlobalInt("count", 4); got != 4 {
		t.Fatalf("expected fallback increment to return 4, got %d", got)
	}
	if got := ctx.State["count"]; got != 4 {
		t.Fatalf("expected fallback state to be updated, got %v", got)
	}
}

func TestContextLocalIsRequestScopedAndSharedStateUsesHelpers(t *testing.T) {
	app := New()
	app.SetGlobal("shared", "app")
	app.Page("/", func(ctx *Context) frontend.Node {
		if ctx.Local == nil {
			t.Fatalf("expected Context.Local to be initialized")
		}
		ctx.Local["request"] = "local"
		return frontend.DivProps(frontend.ElementProps{ID: "app"}, frontend.Text(ctx.Local["request"].(string)+":"+ctx.GetGlobal("shared").(string)))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if !strings.Contains(rr.Body.String(), "local:app") {
		t.Fatalf("expected local and shared state output, got %q", rr.Body.String())
	}
	if got := app.state["request"]; got != nil {
		t.Fatalf("expected Context.Local writes to stay request-local, got %v", got)
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

func TestAssetPolicyRejectsExternalAddScriptInOfflineMode(t *testing.T) {
	app := New()
	app.UseOfflineAssets("/vendor")
	app.AddScript("https://cdn.example.com/widget.js")
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Offline"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected offline policy failure status 500, got %d with body %q", rr.Code, rr.Body.String())
	}
	if body := rr.Body.String(); !strings.Contains(body, "asset policy forbids external URL") || !strings.Contains(body, "https://cdn.example.com/widget.js") {
		t.Fatalf("expected clear external URL policy error, got %q", body)
	}
}

func TestAssetPolicyRejectsExternalStyleTemplateInOfflineMode(t *testing.T) {
	app := New()
	app.UseOfflineAssets("/vendor")
	app.UseStyleTemplate(frontend.StyleTemplate{
		Name:                 "external-template",
		FrameworkStylesheets: []string{"https://cdn.example.com/template.css"},
	})
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Offline"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected offline policy failure status 500, got %d with body %q", rr.Code, rr.Body.String())
	}
	if body := rr.Body.String(); !strings.Contains(body, "asset policy forbids external URL") || !strings.Contains(body, "https://cdn.example.com/template.css") {
		t.Fatalf("expected clear external URL policy error, got %q", body)
	}
}

func TestAssetModeOfflineRejectsDefaultCDNFrameworkAssets(t *testing.T) {
	app := New()
	if err := app.SetAssetMode(AssetModeOffline); err != nil {
		t.Fatalf("SetAssetMode failed: %v", err)
	}
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Offline"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected offline policy failure status 500, got %d with body %q", rr.Code, rr.Body.String())
	}
	if body := rr.Body.String(); !strings.Contains(body, "asset policy forbids external URL") {
		t.Fatalf("expected clear external URL policy error, got %q", body)
	}
}
