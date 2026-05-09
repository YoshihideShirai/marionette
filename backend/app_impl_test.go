package backend

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestContextUpdateGlobalDelegatesToApp(t *testing.T) {
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

func TestFlashCookiesAreIssuedReadAndCleared(t *testing.T) {
	app := New()
	app.Page("/set-flash", func(ctx *Context) frontend.Node {
		ctx.FlashSuccess(" saved ")
		return frontend.Text("set")
	})
	app.Page("/read-flash", func(ctx *Context) frontend.Node {
		flashes := ctx.Flashes()
		if len(flashes) != 1 {
			t.Fatalf("expected one flash, got %#v", flashes)
		}
		if flashes[0].Level != FlashSuccess {
			t.Fatalf("expected success flash level, got %q", flashes[0].Level)
		}
		if flashes[0].Message != "saved" {
			t.Fatalf("expected trimmed flash message, got %q", flashes[0].Message)
		}
		return frontend.Text(flashes[0].Message)
	})

	handler := app.Handler()
	setReq := httptest.NewRequest(http.MethodGet, "/set-flash", nil)
	setRR := httptest.NewRecorder()
	handler.ServeHTTP(setRR, setReq)

	flashCookie := responseCookie(setRR, flashCookieName)
	if flashCookie == nil {
		t.Fatalf("expected %s cookie to be issued", flashCookieName)
	}
	if flashCookie.Value == "" {
		t.Fatalf("expected %s cookie value to be populated", flashCookieName)
	}

	readReq := httptest.NewRequest(http.MethodGet, "/read-flash", nil)
	readReq.AddCookie(flashCookie)
	readRR := httptest.NewRecorder()
	handler.ServeHTTP(readRR, readReq)

	if !strings.Contains(readRR.Body.String(), "saved") {
		t.Fatalf("expected rendered flash message, got %q", readRR.Body.String())
	}
	clearCookie := responseCookie(readRR, flashCookieName)
	if clearCookie == nil {
		t.Fatalf("expected %s clear cookie to be issued", flashCookieName)
	}
	if clearCookie.Value != "" || clearCookie.MaxAge != -1 {
		t.Fatalf("expected clear cookie with empty value and MaxAge -1, got %#v", clearCookie)
	}
}

func TestEmptyFlashDoesNotIssueCookie(t *testing.T) {
	app := New()
	app.Page("/empty-flash", func(ctx *Context) frontend.Node {
		ctx.FlashSuccess("   ")
		return frontend.Text("empty")
	})

	req := httptest.NewRequest(http.MethodGet, "/empty-flash", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if cookie := responseCookie(rr, flashCookieName); cookie != nil {
		t.Fatalf("did not expect %s cookie for empty flash, got %#v", flashCookieName, cookie)
	}
}

func TestSessionCookiePersistsAndCanBeCleared(t *testing.T) {
	app := New()
	app.Page("/set-session", func(ctx *Context) frontend.Node {
		ctx.SetSession("user", "aiko")
		return frontend.Text("set")
	})
	app.Page("/read-session", func(ctx *Context) frontend.Node {
		if got := ctx.Session("user"); got != "aiko" {
			t.Fatalf("expected session value to be restored, got %q", got)
		}
		return frontend.Text(ctx.Session("user"))
	})
	app.Page("/clear-session", func(ctx *Context) frontend.Node {
		if got := ctx.Session("user"); got != "aiko" {
			t.Fatalf("expected session value before clear, got %q", got)
		}
		ctx.ClearSession()
		if got := ctx.Session("user"); got != "" {
			t.Fatalf("expected session value to be empty after clear, got %q", got)
		}
		return frontend.Text("cleared")
	})

	handler := app.Handler()
	setReq := httptest.NewRequest(http.MethodGet, "/set-session", nil)
	setRR := httptest.NewRecorder()
	handler.ServeHTTP(setRR, setReq)

	sessionCookie := responseCookie(setRR, sessionCookieName)
	if sessionCookie == nil {
		t.Fatalf("expected %s cookie to be issued", sessionCookieName)
	}
	if sessionCookie.Value == "" {
		t.Fatalf("expected %s cookie value to be populated", sessionCookieName)
	}

	readReq := httptest.NewRequest(http.MethodGet, "/read-session", nil)
	readReq.AddCookie(sessionCookie)
	readRR := httptest.NewRecorder()
	handler.ServeHTTP(readRR, readReq)
	if !strings.Contains(readRR.Body.String(), "aiko") {
		t.Fatalf("expected rendered session value, got %q", readRR.Body.String())
	}

	clearReq := httptest.NewRequest(http.MethodGet, "/clear-session", nil)
	clearReq.AddCookie(sessionCookie)
	clearRR := httptest.NewRecorder()
	handler.ServeHTTP(clearRR, clearReq)

	clearedCookie := responseCookie(clearRR, sessionCookieName)
	if clearedCookie == nil {
		t.Fatalf("expected %s cookie after ClearSession", sessionCookieName)
	}
	if decoded := decodeSessionFromCookie(t, clearedCookie); len(decoded) != 0 {
		t.Fatalf("expected cleared session cookie to decode to empty map, got %#v", decoded)
	}
}

func TestSecureCookiesUseAppCookieSecureSetting(t *testing.T) {
	app := New()
	app.SetCookieSecure(true)
	app.Page("/secure-cookies", func(ctx *Context) frontend.Node {
		ctx.FlashSuccess("saved")
		ctx.SetSession("user", "aiko")
		return frontend.Text("secure")
	})

	req := httptest.NewRequest(http.MethodGet, "/secure-cookies", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	flashCookie := responseCookie(rr, flashCookieName)
	if flashCookie == nil {
		t.Fatalf("expected %s cookie", flashCookieName)
	}
	if !flashCookie.Secure {
		t.Fatalf("expected %s cookie to be Secure", flashCookieName)
	}
	sessionCookie := responseCookie(rr, sessionCookieName)
	if sessionCookie == nil {
		t.Fatalf("expected %s cookie", sessionCookieName)
	}
	if !sessionCookie.Secure {
		t.Fatalf("expected %s cookie to be Secure", sessionCookieName)
	}
}

func TestRouteBehavior(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(*App)
		method       string
		path         string
		body         string
		contentType  string
		wantStatus   int
		wantContains string
	}{
		{
			name: "page path without leading slash is available with slash",
			setup: func(app *App) {
				app.Page("users", func(ctx *Context) frontend.Node {
					return frontend.Text("users page")
				})
			},
			method:       http.MethodGet,
			path:         "/users",
			wantStatus:   http.StatusOK,
			wantContains: "users page",
		},
		{
			name: "empty page path registers root",
			setup: func(app *App) {
				app.Page("", func(ctx *Context) frontend.Node {
					return frontend.Text("root page")
				})
			},
			method:       http.MethodGet,
			path:         "/",
			wantStatus:   http.StatusOK,
			wantContains: "root page",
		},
		{
			name: "action path with leading slash is available",
			setup: func(app *App) {
				app.Action("/save", func(ctx *Context) frontend.Node {
					return frontend.Text(ctx.FormValue("name"))
				})
			},
			method:       http.MethodPost,
			path:         "/save",
			body:         "name=slash",
			contentType:  "application/x-www-form-urlencoded",
			wantStatus:   http.StatusOK,
			wantContains: "slash",
		},
		{
			name: "action path without leading slash is normalized",
			setup: func(app *App) {
				app.Action("save", func(ctx *Context) frontend.Node {
					return frontend.Text(ctx.FormValue("name"))
				})
			},
			method:       http.MethodPost,
			path:         "/save",
			body:         "name=plain",
			contentType:  "application/x-www-form-urlencoded",
			wantStatus:   http.StatusOK,
			wantContains: "plain",
		},
		{
			name: "post to page route is rejected",
			setup: func(app *App) {
				app.Page("/page", func(ctx *Context) frontend.Node {
					return frontend.Text("page")
				})
			},
			method:       http.MethodPost,
			path:         "/page",
			wantStatus:   http.StatusMethodNotAllowed,
			wantContains: "method not allowed",
		},
		{
			name: "get to action route is rejected",
			setup: func(app *App) {
				app.Action("submit", func(ctx *Context) frontend.Node {
					return frontend.Text("submitted")
				})
			},
			method:       http.MethodGet,
			path:         "/submit",
			wantStatus:   http.StatusMethodNotAllowed,
			wantContains: "method not allowed",
		},
		{
			name: "action rejects invalid form body",
			setup: func(app *App) {
				app.Action("submit", func(ctx *Context) frontend.Node {
					return frontend.Text("submitted")
				})
			},
			method:       http.MethodPost,
			path:         "/submit",
			body:         "name=%zz",
			contentType:  "application/x-www-form-urlencoded",
			wantStatus:   http.StatusBadRequest,
			wantContains: "invalid URL escape",
		},
		{
			name: "action rejects invalid content type combination",
			setup: func(app *App) {
				app.Action("submit", func(ctx *Context) frontend.Node {
					return frontend.Text("submitted")
				})
			},
			method:       http.MethodPost,
			path:         "/submit",
			body:         "name=Aiko",
			contentType:  "application/x-www-form-urlencoded; %",
			wantStatus:   http.StatusBadRequest,
			wantContains: "mime",
		},
		{
			name:         "unregistered root reports missing page registration",
			setup:        func(app *App) {},
			method:       http.MethodGet,
			path:         "/",
			wantStatus:   http.StatusInternalServerError,
			wantContains: "missing app.Page or app.Render registration for /",
		},
		{
			name:         "unregistered non-root returns not found",
			setup:        func(app *App) {},
			method:       http.MethodGet,
			path:         "/missing",
			wantStatus:   http.StatusNotFound,
			wantContains: "404 page not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New()
			tt.setup(app)

			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			rr := httptest.NewRecorder()
			app.Handler().ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d with body %q", tt.wantStatus, rr.Code, rr.Body.String())
			}
			if tt.wantContains != "" && !strings.Contains(rr.Body.String(), tt.wantContains) {
				t.Fatalf("expected body to contain %q, got %q", tt.wantContains, rr.Body.String())
			}
		})
	}
}

func responseCookie(rr *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func decodeSessionFromCookie(t *testing.T, cookie *http.Cookie) map[string]string {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)
	return decodeSession(req)
}

func TestAssetsIgnoreEmptyAndRootPrefixes(t *testing.T) {
	for _, prefix := range []string{"", "/"} {
		t.Run(prefix, func(t *testing.T) {
			app := New()
			app.Assets(prefix, fstest.MapFS{
				"app.css": {Data: []byte("body {}")},
			})

			if len(app.assets) != 0 {
				t.Fatalf("expected no asset routes to be registered, got %d", len(app.assets))
			}

			req := httptest.NewRequest(http.MethodGet, "/app.css", nil)
			rr := httptest.NewRecorder()
			app.Handler().ServeHTTP(rr, req)

			if rr.Code != http.StatusNotFound {
				t.Fatalf("expected 404 for unregistered asset route, got %d", rr.Code)
			}
		})
	}
}

func TestAssetNormalizesAndEscapesURLs(t *testing.T) {
	app := New()
	app.Assets("/assets", fstest.MapFS{})

	tests := map[string]string{
		"../secret.txt":                   "/assets/secret.txt",
		"images/a b.png":                  "/assets/images/a%20b.png",
		"https://cdn.example.com/app.css": "https://cdn.example.com/app.css",
		"data:text/plain,hi":              "data:text/plain,hi",
	}
	for name, want := range tests {
		if got := app.Asset(name); got != want {
			t.Fatalf("Asset(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestAssetsRejectInvalidEscapes(t *testing.T) {
	app := New()
	app.Assets("/assets", fstest.MapFS{
		"app.css": {Data: []byte("body {}")},
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.URL = &url.URL{Path: "/assets/%zz"}
	req.RequestURI = "/assets/%zz"
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid escape, got %d", rr.Code)
	}
}

func TestAssetsDoNotResolveDotDotOutsidePrefix(t *testing.T) {
	app := New()
	app.Assets("/assets", fstest.MapFS{
		"file": {Data: []byte("root file")},
	})

	req := httptest.NewRequest(http.MethodGet, "/assets/%2e%2e/file", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for dot-dot asset request, got %d with body %q", rr.Code, rr.Body.String())
	}
}

func TestAssetsRejectNonGetAndHeadMethods(t *testing.T) {
	app := New()
	app.Assets("/assets", fstest.MapFS{
		"app.css": {Data: []byte("body {}")},
	})

	req := httptest.NewRequest(http.MethodPost, "/assets/app.css", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 for POST asset request, got %d", rr.Code)
	}
}

func TestAssetsAllowDirectoryIndexOnlyWhenEnabled(t *testing.T) {
	fsys := fstest.MapFS{
		"icons":           {Mode: 0o755 | fs.ModeDir},
		"icons/check.svg": {Data: []byte("<svg></svg>")},
	}

	t.Run("disabled", func(t *testing.T) {
		app := New()
		app.Assets("/assets", fsys)

		req := httptest.NewRequest(http.MethodGet, "/assets/icons/", nil)
		rr := httptest.NewRecorder()
		app.Handler().ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("expected 404 when directory index is disabled, got %d", rr.Code)
		}
	})

	t.Run("enabled", func(t *testing.T) {
		app := New()
		app.Assets("/assets", fsys, WithAssetIndex(true))

		req := httptest.NewRequest(http.MethodGet, "/assets/icons/", nil)
		rr := httptest.NewRecorder()
		app.Handler().ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 when directory index is enabled, got %d", rr.Code)
		}
		if body := rr.Body.String(); !strings.Contains(body, "check.svg") {
			t.Fatalf("expected directory index to include child file, got %q", body)
		}
	})
}

func TestAssetContentTypeExtensionNormalization(t *testing.T) {
	app := New()
	app.Assets("/assets", fstest.MapFS{
		"report.csv": {Data: []byte("name\nAiko\n")},
	}, WithAssetContentTypes(map[string]string{"CSV": "text/csv; charset=utf-8"}))

	req := httptest.NewRequest(http.MethodGet, "/assets/report.csv", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("expected normalized CSV content type, got %q", got)
	}
}

func TestContextRequestHelpers(t *testing.T) {
	app := New()
	app.Action("/requests/{id}", func(ctx *Context) frontend.Node {
		got := strings.Join([]string{
			ctx.Param("id"),
			ctx.Query("filter"),
			ctx.FormValue("name"),
		}, ":")
		return frontend.DivProps(frontend.ElementProps{ID: "request-helpers"}, frontend.Text(got))
	})

	form := url.Values{"name": {"Aiko"}}
	req := httptest.NewRequest(http.MethodPost, "/requests/42?filter=active", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %q", rr.Code, rr.Body.String())
	}
	if body := rr.Body.String(); !strings.Contains(body, ">42:active:Aiko<") {
		t.Fatalf("expected route param, query, and form values in response, got %q", body)
	}
}

func TestContextGlobalHelpers(t *testing.T) {
	app := New()
	app.Page("/", func(ctx *Context) frontend.Node {
		ctx.SetGlobal("count", 2)
		ctx.SetGlobal("names", []string{"Aiko", "Ren"})

		snapshot := ctx.GetGlobalSnapshot("names", func(value any) any {
			names, _ := value.([]string)
			return append([]string(nil), names...)
		}).([]string)
		snapshot[0] = "changed"

		next := ctx.IncrementGlobalInt("count", 3)
		stored := ctx.GetGlobalSnapshot("names", func(value any) any {
			names, _ := value.([]string)
			return append([]string(nil), names...)
		}).([]string)

		got := stored[0] + ":" + strconv.Itoa(ctx.GetGlobalInt("count")) + ":" + strconv.Itoa(next)
		return frontend.DivProps(frontend.ElementProps{ID: "global-helpers"}, frontend.Text(got))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if body := rr.Body.String(); !strings.Contains(body, ">Aiko:5:5<") {
		t.Fatalf("expected context global helper output, got %q", body)
	}
	if got := app.GetGlobalInt("count"); got != 5 {
		t.Fatalf("expected app global count to be 5, got %d", got)
	}
}

func TestRenderAndHandleConvenienceMethods(t *testing.T) {
	app := New()
	app.Render(func(ctx *Context) frontend.Node {
		return frontend.DivProps(frontend.ElementProps{ID: "root"}, frontend.Text("root page"))
	})
	app.Handle("save", func(ctx *Context) frontend.Node {
		return frontend.DivProps(frontend.ElementProps{ID: "save"}, frontend.Text("saved"))
	})

	handler := app.Handler()
	rootRR := httptest.NewRecorder()
	handler.ServeHTTP(rootRR, httptest.NewRequest(http.MethodGet, "/", nil))
	if body := rootRR.Body.String(); !strings.Contains(body, "root page") {
		t.Fatalf("expected Render to register / page, got %q", body)
	}

	actionRR := httptest.NewRecorder()
	handler.ServeHTTP(actionRR, httptest.NewRequest(http.MethodPost, "/save", nil))
	if actionRR.Code != http.StatusOK {
		t.Fatalf("expected Handle action status 200, got %d with body %q", actionRR.Code, actionRR.Body.String())
	}
	if body := actionRR.Body.String(); !strings.Contains(body, "saved") {
		t.Fatalf("expected Handle to register POST action, got %q", body)
	}
}

func TestEnableDisableHTMX(t *testing.T) {
	app := New()
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Dashboard"))
	})

	app.DisableHTMX()
	disabledRR := httptest.NewRecorder()
	app.Handler().ServeHTTP(disabledRR, httptest.NewRequest(http.MethodGet, "/", nil))
	if body := disabledRR.Body.String(); strings.Contains(body, assets.HTMXURL) {
		t.Fatalf("did not expect HTMX script after DisableHTMX, got %q", body)
	}

	app.EnableHTMX(true)
	enabledRR := httptest.NewRecorder()
	app.Handler().ServeHTTP(enabledRR, httptest.NewRequest(http.MethodGet, "/", nil))
	if body := enabledRR.Body.String(); !strings.Contains(body, assets.HTMXURL) {
		t.Fatalf("expected HTMX script after EnableHTMX(true), got %q", body)
	}
}

func TestUseAssetProviderAlias(t *testing.T) {
	app := New()
	app.UseAssetProvider(assets.NewLocalAssetProvider("/vendor"))
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
		t.Fatalf("did not expect CDN assets when provider alias is used, got %q", body)
	}
}

func TestUseDaisyUITemplate(t *testing.T) {
	app := New()
	app.UseTailwindCSSTemplate()
	app.UseDaisyUITemplate()
	app.Page("/", func(ctx *Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Dashboard"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, `href="`+assets.DaisyUICSSURL+`"`) {
		t.Fatalf("expected DaisyUI stylesheet in shell, got %q", body)
	}
	if !strings.Contains(body, `src="`+assets.TailwindBrowserURL+`"`) {
		t.Fatalf("expected Tailwind browser script from DaisyUI template in shell, got %q", body)
	}
}
