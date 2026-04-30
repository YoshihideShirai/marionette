package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterAlertExample wires a Marionette page used in docs snippets.
func RegisterAlertExample(app *mb.App) {
	app.Page("/alert", func(ctx *mb.Context) mf.Node {
		return mf.Alert(mf.AlertProps{
			Title:       "Payment failed",
			Description: "Please verify your card and try again.",
			Icon:        "!",
			Props:       mf.ComponentProps{Variant: "error"},
		})
	})
}
