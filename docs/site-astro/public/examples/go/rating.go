package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterRatingExample(app *mb.App) {
	app.Page("/rating", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "rating example"})
	})
}
