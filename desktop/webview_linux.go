//go:build linux && marionette_desktop

package desktop

import webkitgtk "github.com/malivvan/webkitgtk"

func openWebView(url string, options Options) error {
	app := webkitgtk.New(webkitgtk.AppOptions{
		Name: "Marionette",
	})
	app.Open(webkitgtk.WindowOptions{
		Title:           options.Title,
		URL:             url,
		Width:           options.Width,
		Height:          options.Height,
		DevToolsEnabled: options.Debug,
	})
	return app.Run()
}
