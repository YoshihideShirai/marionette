package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterStepsExample(app *mb.App) {
	app.Page("/steps", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "steps example"})
	})
}
