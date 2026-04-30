package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterOverlaySystemExample wires a Marionette page used in docs snippets.
func RegisterOverlaySystemExample(app *mb.App) {
	app.Page("/overlay-system", func(ctx *mb.Context) mf.Node {
		return mf.PageHeader(mf.PageHeaderProps{Title: "OverlaySystem example", Description: "Implement this UI with Marionette components."})
	})
}
