package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterFeedbackExample(app *mb.App) {
	app.Page("/feedback", func(ctx *mb.Context) mf.Node {
		return mf.Stack(mf.StackProps{Gap: "sm"},
			mf.Toast(mf.ToastProps{Title: "Saved", Description: "All changes were synced.", Props: mf.ComponentProps{Variant: "success"}}),
			mf.Alert(mf.AlertProps{Title: "Needs review", Description: "Please check the highlighted fields.", Props: mf.ComponentProps{Variant: "warning"}}),
			mf.EmptyState(mf.EmptyStateProps{Title: "No data", Description: "Create a record to get started."}),
			mf.Skeleton(mf.SkeletonProps{Rows: 2}),
		)
	})
}
