package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	mh "github.com/YoshihideShirai/marionette/frontend/html"
)

// RegisterToastExample wires a Marionette page used in docs snippets.
func RegisterToastExample(app *mb.App) {
	app.Page("/toast", func(ctx *mb.Context) mf.Node {
		return mh.Div(
			mh.H1(mh.Text("Toast example")),
			mh.P(mh.Text("Implement this UI with Marionette components.")),
		)
	})
}
