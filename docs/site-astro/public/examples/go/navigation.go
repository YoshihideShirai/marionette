package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterNavigationExample(app *mb.App) {
	app.Page("/navigation", func(ctx *mb.Context) mf.Node {
		return mf.Stack(mf.StackProps{Gap: "lg"},
			mf.Breadcrumb(mf.BreadcrumbProps{Items: []mf.BreadcrumbItem{
				{Label: "Home", Href: "/"},
				{Label: "Users", Active: true},
			}}),
			mf.Tabs(mf.TabsProps{Items: []mf.TabsItem{
				{Label: "Profile", Href: "/profile", Active: true},
				{Label: "Settings", Href: "/settings"},
			}}),
			mf.Pagination(mf.PaginationProps{Page: 2, TotalPages: 5, PrevHref: "?page=1", NextHref: "?page=3"}),
		)
	})
}
