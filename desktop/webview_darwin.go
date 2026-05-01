//go:build darwin && marionette_desktop

package desktop

import webview "github.com/webview/webview_go"

func openWebView(url string, options Options) error {
	window := webview.New(options.Debug)
	defer window.Destroy()

	window.SetTitle(options.Title)
	window.SetSize(options.Width, options.Height, webview.HintNone)
	window.Navigate(url)
	window.Run()
	return nil
}
