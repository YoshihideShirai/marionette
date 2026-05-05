package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterBoxExample(app *mb.App) {
	app.Page("/box", func(ctx *mb.Context) mf.Node {
		return mf.Box(mf.BoxProps{Tone: "muted", Border: true, Padding: "lg"},
			mf.H3("Box surface"),
			mf.P("Use Box to group related content blocks."),
		)
	})
}
