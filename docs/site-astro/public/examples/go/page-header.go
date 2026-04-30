package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterPageHeaderExample(app *mb.App) {
	app.Page("/page-header", func(ctx *mb.Context) mf.Node {
		return mf.PageHeader(mf.PageHeaderProps{
			Title:       "Users",
			Description: "Manage account records.",
			Actions:     mf.Button("Create", mf.ComponentProps{Size: "sm"}),
		})
	})
}
