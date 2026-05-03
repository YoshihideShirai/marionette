package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterToggleExample(app *mb.App) {
	app.Page("/toggle", func(ctx *mb.Context) mf.Node {
		return mf.Stack(mf.StackProps{Gap: "sm"},
			mf.Toggle("demo-toggle", true),
			mf.ToggleVariant("toggle-primary", true, "primary"),
			mf.ToggleVariant("toggle-secondary", true, "secondary"),
			mf.ToggleVariant("toggle-accent", true, "accent"),
			mf.ToggleVariant("toggle-neutral", true, "neutral"),
			mf.ToggleVariant("toggle-info", true, "info"),
			mf.ToggleVariant("toggle-success", true, "success"),
			mf.ToggleVariant("toggle-warning", true, "warning"),
			mf.ToggleVariant("toggle-error", true, "error"),
			mf.ToggleWithIcons("demo-toggle-icons", true, ""),
		)
	})
}
