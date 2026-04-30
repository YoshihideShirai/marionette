package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterSidebarExample(app *mb.App) {
	app.Page("/sidebar", func(ctx *mb.Context) mf.Node {
		return mf.Sidebar("Marionette", "Admin Console",
			mf.SidebarLink("Dashboard", "/").Active(),
			mf.SidebarLink("Users", "/users"),
			mf.SidebarLink("Analytics", "/analytics"),
			mf.SidebarLink("Settings", "/settings"),
		).Note("Demo workspace", "In-memory data for admin UI prototyping.")
	})
}
