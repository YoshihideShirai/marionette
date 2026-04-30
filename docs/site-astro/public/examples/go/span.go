package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterSpanExample(app *mb.App) {
	app.Page("/span", func(ctx *mb.Context) mf.Node {
		return mf.Stack(mf.StackProps{Direction: "horizontal", Gap: "none"},
			mf.TextComponent(mf.TextProps{Text: "Inline"}),
			mf.TextComponent(mf.TextProps{Text: " text helper"}),
		)
	})
}
