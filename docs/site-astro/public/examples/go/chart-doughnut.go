package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterChartDoughnutExample(app *mb.App) {
	app.Page("/chart-doughnut", func(ctx *mb.Context) mf.Node {
		return mf.Chart(mf.ChartProps{
			Type:   mf.ChartTypeDoughnut,
			Title:  "Plan mix",
			Labels: []string{"Free", "Team", "Enterprise"},
			Datasets: []mf.ChartDataset{
				{Label: "Accounts", Data: []float64{64, 28, 8}},
			},
		})
	})
}
