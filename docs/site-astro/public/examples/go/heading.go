package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterHeadingExample(app *mb.App) {
	app.Page("/heading", func(ctx *mb.Context) mf.Node {
		return mf.Div(
			mf.H1(mf.Text("Heading helper example")),
			mf.H2(mf.Text("H2 subtitle")),
			mf.H3(mf.Text("H3 section")),
			mf.H4(mf.Text("H4 caption")),
		)
	})
}
