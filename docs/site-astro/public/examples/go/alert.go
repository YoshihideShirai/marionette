package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterAlertExample wires a Marionette page used in docs snippets.
func RegisterAlertExample(app *mb.App) {
	app.Page("/alert", func(ctx *mb.Context) mf.Node {
		return mf.Div(
			mf.H1(mf.Text("Alert example")),
			mf.P(mf.Text("Implement this UI with Marionette components.")),
		)
	})
}
