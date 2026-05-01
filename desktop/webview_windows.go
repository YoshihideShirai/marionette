//go:build windows && marionette_desktop

package desktop

import webview2 "github.com/jchv/go-webview2"

func openWebView(url string, options Options) error {
	window := webview2.New(options.Debug)
	defer window.Destroy()

	window.SetTitle(options.Title)
	window.SetSize(options.Width, options.Height, webview2.HintNone)
	window.Navigate(url)
	window.Run()
	return nil
}
