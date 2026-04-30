package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterChartScatterExample(app *mb.App) {
	app.Page("/chart-scatter", func(ctx *mb.Context) mf.Node {
		return mf.Chart(mf.ChartProps{
			Type:  mf.ChartTypeScatter,
			Title: "Response time by load",
			Datasets: []mf.ChartDataset{
				{
					Label:       "Requests",
					Points:      []mf.ChartPoint{{X: 10, Y: 80}, {X: 40, Y: 130}, {X: 80, Y: 220}},
					BorderColor: "#0f766e",
				},
			},
			Options: mf.ChartOptions{XAxisLabel: "Requests/sec", YAxisLabel: "Latency (ms)"},
		})
	})
}
