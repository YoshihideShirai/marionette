package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterTextComponentExample(app *mb.App) {
	app.Page("/text-component", func(ctx *mb.Context) mf.Node {
		return mf.Stack(mf.StackProps{Gap: "sm"},
			mf.TextComponent(mf.TextProps{Text: "Display text", Size: "xl", Weight: "semibold"}),
			mf.TextComponent(mf.TextProps{Text: "Supportive body copy", Tone: "muted"}),
		)
	})
}
