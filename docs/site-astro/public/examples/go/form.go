package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterFormExample wires a Marionette page used in docs snippets.
func RegisterFormExample(app *mb.App) {
	app.Page("/form", func(ctx *mb.Context) mf.Node {
		return mf.Div(
			mf.Element("h1", mf.ElementProps{}, mf.Text("Form example")),
			mf.Element("p", mf.ElementProps{}, mf.Text("Implement this UI with Marionette components.")),
		)
	})
}
