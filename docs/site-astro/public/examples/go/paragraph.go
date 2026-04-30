package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	mh "github.com/YoshihideShirai/marionette/frontend/html"
)

func RegisterParagraphExample(app *mb.App) {
	app.Page("/paragraph", func(ctx *mb.Context) mf.Node {
		return mh.Div(
			mh.P(mh.Text("Paragraph helper for readable long-form text.")),
		)
	})
}
