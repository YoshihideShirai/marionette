package desktop

import "strings"

// Options configures the native desktop window used to host a Marionette app.
type Options struct {
	Title  string
	Width  int
	Height int
	Debug  bool
}

func normalizeOptions(options Options) Options {
	options.Title = strings.TrimSpace(options.Title)
	if options.Title == "" {
		options.Title = "Marionette"
	}
	if options.Width <= 0 {
		options.Width = 1200
	}
	if options.Height <= 0 {
		options.Height = 800
	}
	return options
}
