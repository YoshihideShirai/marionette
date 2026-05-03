package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterTooltipExample(app *mb.App) {
	app.Page("/tooltip", func(ctx *mb.Context) mf.Node {
		return mf.Tooltip("hello", mf.Button("Hover", mf.ComponentProps{}))
	})
}
