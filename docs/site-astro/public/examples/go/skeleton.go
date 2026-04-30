package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterSkeletonExample wires a Marionette page used in docs snippets.
func RegisterSkeletonExample(app *mb.App) {
	app.Page("/skeleton", func(ctx *mb.Context) mf.Node {
		return mf.PageHeader(mf.PageHeaderProps{Title: "Skeleton example", Description: "Implement this UI with Marionette components."})
	})
}
