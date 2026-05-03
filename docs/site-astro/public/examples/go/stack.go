package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterStackExample(app *mb.App) {
	app.Page("/stack", func(ctx *mb.Context) mf.Node {
		return mf.Stack(mf.StackProps{Direction: "horizontal", Gap: "sm", Align: "center"},
			mf.Badge(mf.BadgeProps{Label: "Admin", Props: mf.ComponentProps{Variant: "primary"}}),
			mf.TextComponent(mf.TextProps{Text: "Aiko Tanaka", Weight: "medium"}),
		)
	})
}
