package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterLoadingExample(app *mb.App) {
	app.Page("/loading", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "loading example"})
	})
}
