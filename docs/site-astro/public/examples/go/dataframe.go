package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterDataFrameExample wires a Marionette page used in docs snippets.
func RegisterDataFrameExample(app *mb.App) {
	app.Page("/dataframe", func(ctx *mb.Context) mf.Node {
		return mf.PageHeader(mf.PageHeaderProps{Title: "DataFrame example", Description: "Implement this UI with Marionette components."})
	})
}
