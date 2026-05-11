package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

type dashwindAccountRow struct{ Name, Plan string }

func RegisterDashwindResourcePageExample(app *mb.App) {
	app.Page("/dashwind-resource-page", func(ctx *mb.Context) mf.Node {
		rows := []dashwindAccountRow{{"Acme Co", "Enterprise"}, {"Globex", "Team"}}
		return dw.ResourcePage(dw.ResourcePageProps[dashwindAccountRow]{
			Title:         "Accounts",
			Description:   "Conventional resource page with header actions, filters, table, and row actions.",
			Rows:          rows,
			PrimaryAction: dw.Action{Label: "New account", Href: "#", Class: "btn-primary btn-sm"},
			Search:        mf.InputElement(mf.ElementProps{Attrs: mf.Attrs{"name": "q", "type": "search", "placeholder": "Search accounts", "class": "input w-full"}}),
			Columns: []dw.Column[dashwindAccountRow]{
				{Header: "Account", Cell: func(row dashwindAccountRow) mf.Node { return mf.Text(row.Name) }},
				{Header: "Plan", Cell: func(row dashwindAccountRow) mf.Node { return mf.Badge(mf.BadgeProps{Label: row.Plan}) }},
			},
			RowActions: func(row dashwindAccountRow) []dw.Action {
				return []dw.Action{{Label: "Open", Href: "#", Class: "btn-ghost btn-xs"}}
			},
		})
	})
}
