package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterChartLineExample(app *mb.App) {
	app.Page("/chart-line", func(ctx *mb.Context) mf.Node {
		return mf.Chart(mf.ChartProps{
			Type:   mf.ChartTypeLine,
			Title:  "Onboarding trend",
			Labels: []string{"Jan", "Feb", "Mar"},
			Datasets: []mf.ChartDataset{
				{Label: "Starts", Data: []float64{8, 14, 18}, BorderColor: "#2563eb", Fill: true, Tension: 0.35},
			},
		})
	})
}
