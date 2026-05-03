package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterTimelineExample(app *mb.App) {
	app.Page("/timeline", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "timeline example"})
	})
}
