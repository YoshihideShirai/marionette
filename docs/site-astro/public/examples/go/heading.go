package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	mh "github.com/YoshihideShirai/marionette/frontend/html"
)

func RegisterHeadingExample(app *mb.App) {
	app.Page("/heading", func(ctx *mb.Context) mf.Node {
		return mf.Stack(mf.StackProps{Gap: "sm"},
			mh.H1(mh.Text("Heading helper example")),
			mh.H2(mh.Text("H2 subtitle")),
			mh.H3(mh.Text("H3 section")),
			mh.H4(mh.Text("H4 caption")),
		)
	})
}
