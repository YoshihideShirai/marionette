//go:build !marionette_desktop

package desktop

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/YoshihideShirai/marionette/backend"
	"github.com/YoshihideShirai/marionette/frontend"
)

func TestOpenWebViewStubReportsBuildTag(t *testing.T) {
	err := openWebView("http://127.0.0.1:1/", Options{})
	if err == nil || !strings.Contains(err.Error(), "-tags marionette_desktop") {
		t.Fatalf("expected build tag guidance, got %v", err)
	}
}

func TestRunCleansUpServerAfterWebViewReturns(t *testing.T) {
	app := backend.New()
	app.Page("/", func(ctx *backend.Context) frontend.Node {
		return frontend.Text("ready")
	})

	originalOpenWebView := openWebView
	t.Cleanup(func() { openWebView = originalOpenWebView })

	var openedURL string
	openWebView = func(url string, options Options) error {
		openedURL = url
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return err
		}
		if !strings.Contains(string(body), "ready") {
			return errors.New("server did not serve app while webview was open")
		}
		return nil
	}

	if err := Run(app, Options{}); err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if openedURL == "" {
		t.Fatal("openWebView was not called")
	}

	client := http.Client{Timeout: 100 * time.Millisecond}
	deadline := time.Now().Add(time.Second)
	for {
		resp, err := client.Get(openedURL)
		if err != nil {
			return
		}
		_ = resp.Body.Close()
		if time.Now().After(deadline) {
			t.Fatalf("server at %s was still accepting connections after Run returned", openedURL)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
