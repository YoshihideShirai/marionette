package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterDropdownExample(app *mb.App) {
	app.Page("/dropdown", func(ctx *mb.Context) mf.Node {
		return mf.Dropdown(mf.Button("Menu", mf.ComponentProps{}), mf.Menu(mf.Anchor("#", mf.TextComponent(mf.TextProps{Text: "Item 1"})), mf.Anchor("#", mf.TextComponent(mf.TextProps{Text: "Item 2"}))))
	})
}
