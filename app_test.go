package marionette

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

	"github.com/YoshihideShirai/marionette/frontend/assets"
)

func TestIndexRendersHTMLNotEscaped(t *testing.T) {
	app := New()
	app.Render(func(ctx *Context) Node {
		return DivProps(ElementProps{ID: "app", Class: "card"}, Text("Hello"))
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
		return DivProps(ElementProps{ID: "app", Class: "card"}, Text("Dashboard"))
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

func TestPageCanSetHTMLTitle(t *testing.T) {
	app := New()
	app.Page("/", func(ctx *Context) Node {
		return DivID("app", Text("Dashboard"))
	}, WithTitle(`Users & Teams`))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, `<title>Users &amp; Teams</title>`) {
		t.Fatalf("expected escaped custom title, got %q", body)
	}
}

func TestPageIncludesCustomStyles(t *testing.T) {
	app := New()
	app.AddStylesheet("/assets/app.css")
	app.AddStyle(`
		#marionette-root {
			max-width: 48rem;
		}
	`)
	app.Page("/", func(ctx *Context) Node {
		return DivID("app", Text("Dashboard"))
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
}

func TestPageCanDisableChartsWhenNoChartComponentsAreUsed(t *testing.T) {
	app := New()
	app.DisableCharts()
	app.Page("/", func(ctx *Context) Node {
		return DivID("app", Text("Dashboard"))
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
	app.Page("/", func(ctx *Context) Node {
		return DivID("app", Text("Dashboard"))
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

func TestUnregisteredPageReturnsNotFound(t *testing.T) {
	app := New()
	app.Page("/", func(ctx *Context) Node {
		return DivID("app", Text("Dashboard"))
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
		ctx.SetGlobal("name", ctx.FormValue("name"))
		return DivProps(ElementProps{ID: "users", Class: "card"}, Text(ctx.GetGlobal("name").(string)))
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
		return DivID("users")
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
		ctx.SetGlobal("count", 2)
		return DivID("app", Text(ctx.Query("filter")+":"+ctx.Param("id")+":"+strings.Repeat("x", ctx.GetGlobalInt("count"))))
	})

	req := httptest.NewRequest(http.MethodGet, "/?filter=active", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if !strings.Contains(rr.Body.String(), "active::xx") {
		t.Fatalf("expected query and state helper output, got %q", rr.Body.String())
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
	app.Page("/", func(ctx *Context) Node {
		next := ctx.UpdateGlobal("count", func(old any) any {
			return old.(int) + 2
		})
		return DivID("app", Text(strconv.Itoa(next.(int))))
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
	app.Page("/", func(ctx *Context) Node {
		if ctx.Local == nil {
			t.Fatalf("expected Context.Local to be initialized")
		}
		ctx.Local["request"] = "local"
		return DivID("app", Text(ctx.Local["request"].(string)+":"+ctx.GetGlobal("shared").(string)))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if !strings.Contains(rr.Body.String(), "local:app") {
		t.Fatalf("expected local and shared state output, got %q", rr.Body.String())
	}
	if got := app.GetGlobal("request"); got != nil {
		t.Fatalf("expected Context.Local writes to stay request-local, got %v", got)
	}
}

func TestFlashPersistsForNextRequestAndAutoClears(t *testing.T) {
	app := New()
	app.Action("save", func(ctx *Context) Node {
		ctx.FlashSuccess("saved")
		return DivID("app", Text("ok"))
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
		ctx.SetGlobal("name", ctx.FormValue("name"))
		return DivID("ok")
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
		ctx.SetGlobal("count", 7)
		return DivID("app")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if got := app.GetGlobalInt("count"); got != 7 {
		t.Fatalf("expected shared app/context state to be 7, got %d", got)
	}
}

func TestFlashCookieSecureDefaultsToFalse(t *testing.T) {
	app := New()
	app.Action("save", func(ctx *Context) Node {
		ctx.FlashSuccess("saved")
		return DivID("app", Text("ok"))
	})

	req := httptest.NewRequest(http.MethodPost, "/save", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	var flashCookie *http.Cookie
	for _, c := range rr.Result().Cookies() {
		if c.Name == flashCookieName {
			flashCookie = c
			break
		}
	}
	if flashCookie == nil {
		t.Fatalf("expected %q cookie to be set", flashCookieName)
	}
	if flashCookie.Secure {
		t.Fatalf("expected flash cookie secure to default to false")
	}
}

func TestFlashCookieSecureCanBeEnabled(t *testing.T) {
	app := New()
	app.SetCookieSecure(true)
	app.Action("save", func(ctx *Context) Node {
		ctx.FlashSuccess("saved")
		return DivID("app", Text("ok"))
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
	if !flashCookie.Secure {
		t.Fatalf("expected flash cookie secure to be true when enabled")
	}

	getReq := httptest.NewRequest(http.MethodGet, "/", nil)
	getReq.AddCookie(flashCookie)
	getRes := httptest.NewRecorder()
	app.Handler().ServeHTTP(getRes, getReq)

	var clearCookie *http.Cookie
	for _, c := range getRes.Result().Cookies() {
		if c.Name == flashCookieName {
			clearCookie = c
			break
		}
	}
	if clearCookie == nil {
		t.Fatalf("expected %q cookie to be cleared", flashCookieName)
	}
	if !clearCookie.Secure {
		t.Fatalf("expected clear flash cookie secure to be true when enabled")
	}
}

func TestSessionSetAndReadFromCookie(t *testing.T) {
	app := New()
	app.Action("set-session", func(ctx *Context) Node {
		ctx.SetSession("user_id", "u_123")
		return Div(Text("ok"))
	})
	app.Page("/profile", func(ctx *Context) Node {
		return Div(Text(ctx.Session("user_id")))
	})

	setReq := httptest.NewRequest(http.MethodPost, "/set-session", strings.NewReader(""))
	setReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setRes := httptest.NewRecorder()
	app.Handler().ServeHTTP(setRes, setReq)

	if setRes.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", setRes.Code)
	}
	cookie := setRes.Result().Cookies()[0]
	if cookie.Name != sessionCookieName {
		t.Fatalf("expected %q cookie to be set", sessionCookieName)
	}

	readReq := httptest.NewRequest(http.MethodGet, "/profile", nil)
	readReq.AddCookie(cookie)
	readRes := httptest.NewRecorder()
	app.Handler().ServeHTTP(readRes, readReq)

	if !strings.Contains(readRes.Body.String(), "u_123") {
		t.Fatalf("expected session value in response, got %q", readRes.Body.String())
	}
}

func TestSessionCanBeCleared(t *testing.T) {
	app := New()
	app.Action("logout", func(ctx *Context) Node {
		ctx.ClearSession()
		return Div(Text("bye"))
	})

	encoded, err := encodeSession(map[string]string{"user_id": "u_123"})
	if err != nil {
		t.Fatalf("unexpected encode error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(""))
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: encoded})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	app.Handler().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}
	cookies := res.Result().Cookies()
	if len(cookies) == 0 || cookies[0].Name != sessionCookieName {
		t.Fatalf("expected %q cookie to be set", sessionCookieName)
	}
	outReq := httptest.NewRequest(http.MethodGet, "/", nil)
	outReq.AddCookie(cookies[0])
	decoded := decodeSession(outReq)
	if len(decoded) != 0 {
		t.Fatalf("expected cleared session, got %v", decoded)
	}
}
