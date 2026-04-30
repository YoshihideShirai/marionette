package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	mh "github.com/YoshihideShirai/marionette/frontend/html"
)

// RegisterOverlaySystemExample wires a Marionette page used in docs snippets.
func RegisterOverlaySystemExample(app *mb.App) {
	app.Page("/overlay-system", func(ctx *mb.Context) mf.Node {
		return mh.Div(
			mh.H1(mh.Text("OverlaySystem example")),
			mh.P(mh.Text("Implement this UI with Marionette components.")),
		)
	})
}
