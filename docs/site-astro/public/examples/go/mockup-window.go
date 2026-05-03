package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterMockupWindowExample(app *mb.App) {
	app.Page("/mockup-window", func(ctx *mb.Context) mf.Node {
		return mf.MockupWindow("mockup-window", mf.TextComponent(mf.TextProps{Text: "Mock content"}))
	})
}
