package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterTimelineExample(app *mb.App) {
	app.Page("/timeline", func(ctx *mb.Context) mf.Node {
		return mf.Timeline(mf.TimelineItem("1", "", mf.TextComponent(mf.TextProps{Text: "Import frontend"})), mf.TimelineItem("2", "", mf.TextComponent(mf.TextProps{Text: "Call frontend.Timeline()"})))
	})
}
