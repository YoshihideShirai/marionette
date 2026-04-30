package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterAlertExample wires a Marionette page used in docs snippets.
func RegisterAlertExample(app *mb.App) {
	app.Page("/alert", func(ctx *mb.Context) mf.Node {
		return mf.PageHeader(mf.PageHeaderProps{Title: "Alert example", Description: "Implement this UI with Marionette components."})
	})
}
