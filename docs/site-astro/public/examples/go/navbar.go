package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterNavbarExample(app *mb.App) {
	app.Page("/navbar", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "navbar example"})
	})
}
