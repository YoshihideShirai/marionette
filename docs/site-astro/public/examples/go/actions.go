package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterActionsExample(app *mb.App) {
	app.Page("/actions", func(ctx *mb.Context) mf.Node {
		return mf.Actions(mf.ActionsProps{Gap: "sm", Justify: "between", Wrap: true},
			mf.Button("Back", mf.ComponentProps{Variant: "ghost"}),
			mf.Button("Save", mf.ComponentProps{Variant: "primary"}),
		)
	})
}
