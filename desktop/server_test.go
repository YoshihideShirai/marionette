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

func TestStartLocalServerServesPagesActionsAndAssets(t *testing.T) {
	app := backend.New()
	app.Set("name", "Aiko")
	app.Assets("/assets", fstest.MapFS{
		"app.css": {Data: []byte("body { color: red; }")},
	})
	app.Page("/", func(ctx *backend.Context) frontend.Node {
		return frontend.Container(frontend.ContainerProps{}, frontend.Text("Hello "+ctx.Get("name").(string)))
	})
	app.Action("rename", func(ctx *backend.Context) frontend.Node {
		ctx.Set("name", ctx.FormValue("name"))
		return frontend.Text("Hello " + ctx.Get("name").(string))
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
