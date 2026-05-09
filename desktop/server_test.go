package desktop

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/YoshihideShirai/marionette/backend"
	"github.com/YoshihideShirai/marionette/frontend"
)

func TestNormalizeOptionsDefaults(t *testing.T) {
	got := normalizeOptions(Options{})
	want := Options{Title: "Marionette", Width: 1200, Height: 800}
	if got != want {
		t.Fatalf("normalizeOptions(Options{}) = %+v, want %+v", got, want)
	}
}

func TestNormalizeOptionsTrimsBlankTitleToDefault(t *testing.T) {
	got := normalizeOptions(Options{Title: " \t\n ", Width: 640, Height: 480})
	if got.Title != "Marionette" {
		t.Fatalf("blank title normalized to %q, want %q", got.Title, "Marionette")
	}
}

func TestNormalizeOptionsDefaultInvalidSizes(t *testing.T) {
	tests := []struct {
		name    string
		options Options
	}{
		{name: "negative width", options: Options{Title: "App", Width: -1, Height: 480}},
		{name: "zero width", options: Options{Title: "App", Width: 0, Height: 480}},
		{name: "negative height", options: Options{Title: "App", Width: 640, Height: -1}},
		{name: "zero height", options: Options{Title: "App", Width: 640, Height: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeOptions(tt.options)
			if tt.options.Width <= 0 && got.Width != 1200 {
				t.Fatalf("width normalized to %d, want default 1200", got.Width)
			}
			if tt.options.Height <= 0 && got.Height != 800 {
				t.Fatalf("height normalized to %d, want default 800", got.Height)
			}
		})
	}
}

func TestNormalizeOptionsPreservesValidValues(t *testing.T) {
	got := normalizeOptions(Options{Title: "  Custom App  ", Width: 1024, Height: 768, Debug: true})
	want := Options{Title: "Custom App", Width: 1024, Height: 768, Debug: true}
	if got != want {
		t.Fatalf("normalizeOptions preserved %+v, want %+v", got, want)
	}
}

func TestStartLocalServerRejectsNilHandler(t *testing.T) {
	server, err := startLocalServer(nil)
	if err == nil {
		t.Fatal("expected nil handler error")
	}
	if server != nil {
		t.Fatalf("server = %#v, want nil", server)
	}
	if got, want := err.Error(), "desktop: handler is nil"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
}

func TestStartLocalServerServesPagesActionsAndAssets(t *testing.T) {
	app := backend.New()
	app.SetGlobal("name", "Aiko")
	app.Assets("/assets", fstest.MapFS{
		"app.css": {Data: []byte("body { color: red; }")},
	})
	app.Page("/", func(ctx *backend.Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Hello "+ctx.GetGlobal("name").(string)))
	})
	app.Action("rename", func(ctx *backend.Context) frontend.Node {
		ctx.SetGlobal("name", ctx.FormValue("name"))
		return frontend.Text("Hello " + ctx.GetGlobal("name").(string))
	})

	server, err := startLocalServer(app.Handler())
	if err != nil {
		t.Fatalf("startLocalServer failed: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			t.Fatalf("server shutdown failed: %v", err)
		}
	}()

	body := getBody(t, server.URL)
	if !strings.Contains(body, "Hello Aiko") {
		t.Fatalf("expected page response from local server, got %q", body)
	}

	form := url.Values{"name": {"Ren"}}
	resp, err := http.PostForm(server.URL+"rename", form)
	if err != nil {
		t.Fatalf("POST action failed: %v", err)
	}
	actionBody, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		t.Fatalf("read action body failed: %v", err)
	}
	if got := strings.TrimSpace(string(actionBody)); got != "<span>Hello Ren</span>" {
		t.Fatalf("expected action response, got %q", got)
	}

	assetBody := getBody(t, server.URL+"assets/app.css")
	if strings.TrimSpace(assetBody) != "body { color: red; }" {
		t.Fatalf("expected asset response, got %q", assetBody)
	}
}

func TestLocalServerShutdownNilContext(t *testing.T) {
	server, err := startLocalServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	if err != nil {
		t.Fatalf("startLocalServer failed: %v", err)
	}

	if err := server.Shutdown(nil); err != nil {
		t.Fatalf("Shutdown(nil) failed: %v", err)
	}
}

func TestLocalServerShutdownIsIdempotent(t *testing.T) {
	server, err := startLocalServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	if err != nil {
		t.Fatalf("startLocalServer failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		t.Fatalf("first Shutdown failed: %v", err)
	}
	if err := server.Shutdown(ctx); err != nil {
		t.Fatalf("second Shutdown failed: %v", err)
	}
}

func TestRunRejectsNilApp(t *testing.T) {
	if err := Run(nil, Options{}); err == nil {
		t.Fatal("expected nil app error")
	}
}

func getBody(t *testing.T, rawURL string) string {
	t.Helper()
	resp, err := http.Get(rawURL)
	if err != nil {
		t.Fatalf("GET %s failed: %v", rawURL, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read %s failed: %v", rawURL, err)
	}
	return string(body)
}
