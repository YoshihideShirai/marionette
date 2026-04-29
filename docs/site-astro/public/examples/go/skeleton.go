package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterSkeletonExample wires a Marionette page used in docs snippets.
func RegisterSkeletonExample(app *mb.App) {
	app.Page("/skeleton", func(ctx *mb.Context) mf.Node {
		return mf.Div(
			mf.H1(mf.Text("Skeleton example")),
			mf.P(mf.Text("Implement this UI with Marionette components.")),
		)
	})
}
