package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterThemeControllerExample(app *mb.App) {
	app.Page("/theme-controller", func(ctx *mb.Context) mf.Node {
		return mf.ThemeController(
			mf.ThemeControllerOption("light", false, "btn join-item"),
			mf.ThemeControllerOption("dark", false, "btn join-item"),
			mf.ThemeControllerOption("cupcake", true, "btn join-item"),
		)
	})
}
