package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterDropdownExample(app *mb.App) {
	app.Page("/dropdown", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "dropdown example"})
	})
}
