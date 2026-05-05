package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterBadgeExample(app *mb.App) {
	app.Page("/badge", func(ctx *mb.Context) mf.Node {
		return mf.Actions(mf.ActionsProps{Gap: "sm", Wrap: true},
			mf.Badge(mf.BadgeProps{Label: "New", Props: mf.ComponentProps{Variant: "primary"}}),
			mf.Badge(mf.BadgeProps{Label: "Beta", Props: mf.ComponentProps{Variant: "outline"}}),
			mf.Badge(mf.BadgeProps{Label: "Stable", Props: mf.ComponentProps{Variant: "success"}}),
		)
	})
}
