package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterStatExample(app *mb.App) {
	app.Page("/stat", func(ctx *mb.Context) mf.Node {
		return mf.Stat("Sales", "$12K", "+12%")
	})
}
