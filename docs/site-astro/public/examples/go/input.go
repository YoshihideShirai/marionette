package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterInputExample(app *mb.App) {
	app.Page("/input", func(ctx *mb.Context) mf.Node {
		return mf.InputWithOptions("start_date", "2030-01-01", mf.InputOptions{
			Type:        "date",
			Placeholder: "Start date",
			Required:    true,
			Props:       mf.ComponentProps{Size: "sm"},
		})
	})
}
