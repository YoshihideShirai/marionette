//go:build darwin && marionette_desktop

package desktop

// macOS desktop runtime is intentionally unsupported.
func openWebView(_ string, _ Options) error {
	return marionetteDesktopMacOSIsUnsupported
}
