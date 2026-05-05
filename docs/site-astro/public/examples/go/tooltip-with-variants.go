package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterTooltipWithVariantsExample(app *mb.App) {
	app.Page("/tooltip-with-variants", func(ctx *mb.Context) mf.Node {
		return mf.TooltipWithVariants("Saved", mf.Button("Hover", mf.ComponentProps{}), "right", "success", false)
	})
}
