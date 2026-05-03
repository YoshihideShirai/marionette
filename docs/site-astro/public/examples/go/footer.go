package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterFooterExample(app *mb.App) {
	app.Page("/footer", func(ctx *mb.Context) mf.Node {
		return mf.Footer(mf.TextComponent(mf.TextProps{Text: "Footer alias"}))
	})
}
