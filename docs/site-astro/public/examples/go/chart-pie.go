package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterChartPieExample(app *mb.App) {
	app.Page("/chart-pie", func(ctx *mb.Context) mf.Node {
		return mf.Chart(mf.ChartProps{
			Type:   mf.ChartTypePie,
			Title:  "Traffic sources",
			Labels: []string{"Search", "Direct", "Referral"},
			Datasets: []mf.ChartDataset{
				{Label: "Sessions", Data: []float64{48, 32, 20}, BackgroundColor: "#2563eb"},
			},
		})
	})
}
