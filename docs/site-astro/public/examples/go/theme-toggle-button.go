package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterThemeToggleButtonExample(app *mb.App) {
	app.Page("/theme-toggle-button", func(ctx *mb.Context) mf.Node {
		return mf.ThemeToggleButton(mf.ComponentProps{Variant: "outline", Size: "sm"})
	})
}
