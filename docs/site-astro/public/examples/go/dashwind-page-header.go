package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

func RegisterDashwindPageHeaderExample(app *mb.App) {
	app.Page("/dashwind-page-header", func(ctx *mb.Context) mf.Node {
		return dw.PageHeader(dw.PageHeaderProps{
			Title:       "Customers",
			Description: "Track accounts, lifecycle stage, and owner activity.",
			Actions:     mf.Button("New customer", mf.ComponentProps{Variant: "primary", Size: "sm"}),
		})
	})
}
