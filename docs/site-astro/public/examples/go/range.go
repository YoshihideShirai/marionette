package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterRangeExample(app *mb.App) {
	app.Page("/range", func(ctx *mb.Context) mf.Node {
		return mf.Range("volume", 40, 0, 100)
	})
}
