package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterRadialProgressExample(app *mb.App) {
	app.Page("/radial-progress", func(ctx *mb.Context) mf.Node {
		return mf.RadialProgress(70, "")
	})
}
