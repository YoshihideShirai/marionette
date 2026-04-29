package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterToastExample wires a Marionette page used in docs snippets.
func RegisterToastExample(app *mb.App) {
	app.Page("/toast", func(ctx *mb.Context) mf.Node {
		return mf.Div(
			mf.H1(mf.Text("Toast example")),
			mf.P(mf.Text("Implement this UI with Marionette components.")),
		)
	})
}
