package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterHeroExample(app *mb.App) {
	app.Page("/hero", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "hero example"})
	})
}
