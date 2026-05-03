package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterCarouselItemExample(app *mb.App) {
	app.Page("/carousel-item", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "carousel-item example"})
	})
}
