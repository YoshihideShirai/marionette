package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterToastExample wires a Marionette page used in docs snippets.
func RegisterToastExample(app *mb.App) {
	app.Page("/toast", func(ctx *mb.Context) mf.Node {
		return mf.Toast(mf.ToastProps{
			Title:       "Settings saved",
			Description: "Your workspace preferences were updated.",
			Icon:        "OK",
			Props:       mf.ComponentProps{Variant: "success"},
		})
	})
}
