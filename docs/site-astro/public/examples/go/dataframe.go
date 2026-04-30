package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	mh "github.com/YoshihideShirai/marionette/frontend/html"
)

// RegisterDataFrameExample wires a Marionette page used in docs snippets.
func RegisterDataFrameExample(app *mb.App) {
	app.Page("/dataframe", func(ctx *mb.Context) mf.Node {
		return mh.Div(
			mh.H1(mh.Text("DataFrame example")),
			mh.P(mh.Text("Implement this UI with Marionette components.")),
		)
	})
}
