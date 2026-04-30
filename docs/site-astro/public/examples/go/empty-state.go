package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterEmptyStateExample(app *mb.App) {
	app.Page("/empty-state", func(ctx *mb.Context) mf.Node {
		return mf.EmptyState(mf.EmptyStateProps{
			Title:       "No users",
			Description: "Create a user to populate this table.",
			Icon:        "0",
		})
	})
}
