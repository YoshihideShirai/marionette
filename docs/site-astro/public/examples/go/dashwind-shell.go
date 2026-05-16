package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

func RegisterDashwindShellExample(app *mb.App) {
	app.Page("/dashwind-shell", func(ctx *mb.Context) mf.Node {
		return dw.Shell(dw.ShellProps{
			Brand:       dw.Brand{Title: "DashWind", Subtitle: "Admin kit", Mark: "D", Href: "#"},
			CurrentPath: "/dashboard",
			Navigation: dw.Navigation{{Label: "Workspace", Items: []dw.NavItem{
				{Path: "/dashboard", Label: "Dashboard", Icon: "▦"},
				{Path: "/customers", Label: "Customers", Icon: "☷", Badge: "12"},
			}}},
			SearchPlaceholder: "Search dashboard",
		}, dw.CardPanel(dw.CardPanelProps{Title: "Shell content", Description: "Responsive drawer, topbar, search, and actions."}, mf.Text("Ready for htmx fragments.")))
	})
}
