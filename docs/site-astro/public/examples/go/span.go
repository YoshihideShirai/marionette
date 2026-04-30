package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	mh "github.com/YoshihideShirai/marionette/frontend/html"
)

func RegisterSpanExample(app *mb.App) {
	app.Page("/span", func(ctx *mb.Context) mf.Node {
		return mh.Div(
			mh.Span(mh.Text("Inline")),
			mh.Span(mh.Text(" text helper")),
		)
	})
}
