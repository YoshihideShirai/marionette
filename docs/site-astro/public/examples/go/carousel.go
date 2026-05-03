package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterCarouselExample(app *mb.App) {
	app.Page("/carousel", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "carousel example"})
	})
}
