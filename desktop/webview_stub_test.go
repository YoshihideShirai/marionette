//go:build !marionette_desktop

package desktop

import (
	"strings"
	"testing"
)

func TestOpenWebViewStubReportsBuildTag(t *testing.T) {
	err := openWebView("http://127.0.0.1:1/", Options{})
	if err == nil || !strings.Contains(err.Error(), "-tags marionette_desktop") {
		t.Fatalf("expected build tag guidance, got %v", err)
	}
}
