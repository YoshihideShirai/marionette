package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterDataFrameExample wires a Marionette page used in docs snippets.
func RegisterDataFrameExample(app *mb.App) {
	app.Page("/dataframe", func(ctx *mb.Context) mf.Node {
		return mf.Div(
			mf.Element("h1", mf.ElementProps{}, mf.Text("DataFrame example")),
			mf.Element("p", mf.ElementProps{}, mf.Text("Implement this UI with Marionette components.")),
		)
	})
}
