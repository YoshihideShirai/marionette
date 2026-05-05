package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterBoxExample(app *mb.App) {
	app.Page("/box", func(ctx *mb.Context) mf.Node {
		return mf.Box(mf.BoxProps{Tone: "muted", Border: true, Padding: "lg"},
			mf.H3(mf.Text("Box surface")),
			mf.TextComponent(mf.TextProps{Text: "Use Box to group related content blocks."}),
		)
	})
}
