package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterDataFrameChartExample wires a Marionette page used in docs snippets.
func RegisterDataFrameChartExample(app *mb.App) {
	app.Page("/dataframe-chart", func(ctx *mb.Context) mf.Node {
		return mf.DataFrameChart(nil, mf.DataFrameChartProps{
			Chart: mf.ChartProps{
				Type:   mf.ChartTypeBar,
				Title:  "DataFrame chart",
				Labels: []string{"North", "South", "West"},
				Datasets: []mf.ChartDataset{
					{Label: "Revenue", Data: []float64{70, 56, 82}},
				},
			},
		})
	})
}
