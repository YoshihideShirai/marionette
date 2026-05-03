package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterJoinExample(app *mb.App) {
	app.Page("/join", func(ctx *mb.Context) mf.Node {
		return mf.Join(mf.Button("1", mf.ComponentProps{}), mf.Button("2", mf.ComponentProps{}))
	})
}
