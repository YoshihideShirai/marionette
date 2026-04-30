package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterDataFrameChartExample wires a Marionette page used in docs snippets.
func RegisterDataFrameChartExample(app *mb.App) {
	app.Page("/dataframe-chart", func(ctx *mb.Context) mf.Node {
		return mf.PageHeader(mf.PageHeaderProps{Title: "DataFrameChart example", Description: "Implement this UI with Marionette components."})
	})
}
