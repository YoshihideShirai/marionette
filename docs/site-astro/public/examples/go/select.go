package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterSelectExample(app *mb.App) {
	app.Page("/select", func(ctx *mb.Context) mf.Node {
		return mf.Select("role", []mf.SelectOption{
			{Label: "Admin", Value: "admin"},
			{Label: "Viewer", Value: "viewer", Selected: true},
		}, mf.ComponentProps{})
	})
}
