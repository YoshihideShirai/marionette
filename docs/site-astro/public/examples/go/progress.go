package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterProgressExample(app *mb.App) {
	app.Page("/progress", func(ctx *mb.Context) mf.Node {
		return mf.Stack(mf.StackProps{Gap: "lg"},
			mf.Progress(mf.ProgressProps{
				Value:     72,
				Max:       100,
				Label:     "Upload progress",
				ShowValue: true,
				Props:     mf.ComponentProps{Variant: "success", Size: "lg"},
			}),
			mf.Progress(mf.ProgressProps{
				Label:         "Preparing import",
				Indeterminate: true,
				Props:         mf.ComponentProps{Variant: "info", Size: "sm"},
			}),
		)
	})
}
