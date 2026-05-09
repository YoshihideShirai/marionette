//go:build darwin && marionette_desktop

package desktop

// macOS desktop runtime is intentionally unsupported.
var openWebView = defaultOpenWebView

func defaultOpenWebView(_ string, _ Options) error {
	return marionetteDesktopMacOSIsUnsupported
}
