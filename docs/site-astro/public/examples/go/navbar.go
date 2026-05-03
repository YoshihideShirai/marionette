package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterNavbarExample(app *mb.App) {
	app.Page("/navbar", func(ctx *mb.Context) mf.Node {
		return mf.Navbar(mf.TextComponent(mf.TextProps{Text: "Marionette"}), mf.Button("Docs", mf.ComponentProps{Variant: "ghost", Size: "sm"}), mf.Button("Login", mf.ComponentProps{Variant: "primary", Size: "sm"}))
	})
}
