package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterTextareaExample(app *mb.App) {
	app.Page("/textarea", func(ctx *mb.Context) mf.Node {
		return mf.Textarea("notes", "Initial onboarding notes", mf.TextareaOptions{
			Rows:        4,
			Placeholder: "Add notes",
			Props:       mf.ComponentProps{Size: "sm"},
		})
	})
}
