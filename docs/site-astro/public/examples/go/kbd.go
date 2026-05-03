package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterKbdExample(app *mb.App) {
	app.Page("/kbd", func(ctx *mb.Context) mf.Node {
		return mf.Stack(mf.StackProps{Direction: "horizontal", Gap: "xs", Align: "center"},
			mf.TextComponent(mf.TextProps{Text: "Press"}),
			mf.Kbd("⌘"),
			mf.TextComponent(mf.TextProps{Text: "+"}),
			mf.Kbd("K"),
		)
	})
}
