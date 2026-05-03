package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterMenuExample(app *mb.App) {
	app.Page("/menu", func(ctx *mb.Context) mf.Node {
		return mf.Menu(mf.Anchor("#", mf.TextComponent(mf.TextProps{Text: "Item 1"})), mf.Anchor("#", mf.TextComponent(mf.TextProps{Text: "Item 2"})))
	})
}
