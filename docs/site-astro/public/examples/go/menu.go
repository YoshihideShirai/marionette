package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterMenuExample(app *mb.App) {
	app.Page("/menu", func(ctx *mb.Context) mf.Node {
		return mf.Menu(mf.TextComponent(mf.TextProps{Text: "Item 1"}), mf.TextComponent(mf.TextProps{Text: "Item 2"}))
	})
}
