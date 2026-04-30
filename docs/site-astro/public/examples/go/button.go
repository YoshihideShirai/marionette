package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterButtonExample(app *mb.App) {
	app.Page("/button", func(ctx *mb.Context) mf.Node {
		return mf.Actions(mf.ActionsProps{Gap: "sm"},
			mf.Button("Cancel", mf.ComponentProps{Variant: "ghost"}),
			mf.Button("Save", mf.ComponentProps{Variant: "primary"}),
		)
	})
}
