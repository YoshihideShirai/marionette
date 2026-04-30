package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterTabsExample(app *mb.App) {
	app.Page("/tabs", func(ctx *mb.Context) mf.Node {
		return mf.Tabs(mf.TabsProps{
			AriaLabel: "User sections",
			Items: []mf.TabsItem{
				{Label: "Profile", Href: "/users/1/profile", Active: true},
				{Label: "Permissions", Href: "/users/1/permissions"},
				{Label: "Audit", Disabled: true},
			},
		})
	})
}
