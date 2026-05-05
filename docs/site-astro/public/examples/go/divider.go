package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterDividerExample(app *mb.App) {
	app.Page("/divider", func(ctx *mb.Context) mf.Node {
		return mf.Stack(mf.StackProps{Gap: "md"},
			mf.TextComponent(mf.TextProps{Text: "Section A"}),
			mf.Divider(mf.DividerProps{Spacing: "md"}),
			mf.TextComponent(mf.TextProps{Text: "Section B"}),
		)
	})
}
