package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterBadgeExample(app *mb.App) {
	app.Page("/badge", func(ctx *mb.Context) mf.Node {
		return mf.Actions(mf.ActionsProps{Gap: "sm", Wrap: true},
			mf.Badge(mf.BadgeProps{Label: "New", Variant: "primary"}),
			mf.Badge(mf.BadgeProps{Label: "Beta", Variant: "outline"}),
			mf.Badge(mf.BadgeProps{Label: "Stable", Variant: "success"}),
		)
	})
}
