package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterSpanExample(app *mb.App) {
	app.Page("/span", func(ctx *mb.Context) mf.Node {
		return mf.Div(
			mf.Span(mf.Text("Inline")),
			mf.Span(mf.Text(" text helper")),
		)
	})
}
