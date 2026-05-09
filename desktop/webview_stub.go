//go:build !marionette_desktop

package desktop

import "fmt"

var openWebView = defaultOpenWebView

func defaultOpenWebView(_ string, _ Options) error {
	return fmt.Errorf("desktop: WebView support is not enabled; rebuild with -tags marionette_desktop")
}
