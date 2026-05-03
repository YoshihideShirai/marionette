package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterCodeExample(app *mb.App) {
	app.Page("/code", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "code example"})
	})
}
