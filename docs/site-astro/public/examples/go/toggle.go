package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterToggleExample(app *mb.App) {
	app.Page("/toggle", func(ctx *mb.Context) mf.Node {
		return mf.Toggle("demo-toggle", true)
	})
}
