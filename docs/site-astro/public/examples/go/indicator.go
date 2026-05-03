package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterIndicatorExample(app *mb.App) {
	app.Page("/indicator", func(ctx *mb.Context) mf.Node {
		return mf.Indicator(mf.Badge(mf.BadgeProps{Label: "new", Props: mf.ComponentProps{Variant: "secondary"}}), mf.Button("Inbox", mf.ComponentProps{}))
	})
}
