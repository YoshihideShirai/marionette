package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterBreadcrumbExample(app *mb.App) {
	app.Page("/breadcrumb", func(ctx *mb.Context) mf.Node {
		return mf.Breadcrumb(mf.BreadcrumbProps{
			Items: []mf.BreadcrumbItem{
				{Label: "Home", Href: "/"},
				{Label: "Users", Href: "/users"},
				{Label: "Aiko", Active: true},
			},
		})
	})
}
