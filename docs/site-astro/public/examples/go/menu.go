package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterMenuExample(app *mb.App) {
	app.Page("/menu", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "menu example"})
	})
}
