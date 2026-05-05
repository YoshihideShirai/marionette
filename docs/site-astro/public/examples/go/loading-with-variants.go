package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterLoadingWithVariantsExample(app *mb.App) {
	app.Page("/loading-with-variants", func(ctx *mb.Context) mf.Node {
		return mf.JoinWithDirection("horizontal",
			mf.LoadingWithVariants("spinner", "sm"),
			mf.LoadingWithVariants("dots", "md"),
			mf.LoadingWithVariants("ring", "lg"),
		)
	})
}
