package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

type dashwindLeadRow struct{ Name, Company, Stage string }

func RegisterDashwindDataTableExample(app *mb.App) {
	app.Page("/dashwind-data-table", func(ctx *mb.Context) mf.Node {
		rows := []dashwindLeadRow{{"Ada Lovelace", "Analytical Engines", "Qualified"}, {"Grace Hopper", "Compiler Labs", "Proposal"}}
		return dw.DataTable(dw.DataTableProps[dashwindLeadRow]{
			Rows:  rows,
			Zebra: true,
			Columns: []dw.Column[dashwindLeadRow]{
				{Header: "Name", Sortable: true, Cell: func(row dashwindLeadRow) mf.Node { return mf.Text(row.Name) }},
				{Header: "Company", Cell: func(row dashwindLeadRow) mf.Node { return mf.Text(row.Company) }},
				{Header: "Stage", Class: "text-right", Cell: func(row dashwindLeadRow) mf.Node {
					return mf.Badge(mf.BadgeProps{Label: row.Stage, Props: mf.ComponentProps{Variant: "primary"}})
				}},
			},
		})
	})
}
