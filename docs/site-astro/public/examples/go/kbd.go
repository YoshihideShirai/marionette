package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterKbdExample(app *mb.App) {
	app.Page("/kbd", func(ctx *mb.Context) mf.Node {
		return mf.Div(mf.TextComponent(mf.TextProps{Text: "Press "}), mf.Kbd("⌘"), mf.TextComponent(mf.TextProps{Text: " + "}), mf.Kbd("K"))
	})
}
