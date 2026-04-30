package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	mh "github.com/YoshihideShirai/marionette/frontend/html"
)

// RegisterDataFrameChartExample wires a Marionette page used in docs snippets.
func RegisterDataFrameChartExample(app *mb.App) {
	app.Page("/dataframe-chart", func(ctx *mb.Context) mf.Node {
		return mh.Div(
			mh.H1(mh.Text("DataFrameChart example")),
			mh.P(mh.Text("Implement this UI with Marionette components.")),
		)
	})
}
