package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterChartBarExample(app *mb.App) {
	app.Page("/chart-bar", func(ctx *mb.Context) mf.Node {
		return mf.Chart(mf.ChartProps{
			Type:   mf.ChartTypeBar,
			Title:  "Role distribution",
			Labels: []string{"Admin", "Editor", "Viewer"},
			Datasets: []mf.ChartDataset{
				{Label: "Users", Data: []float64{3, 7, 12}, BackgroundColor: "#93c5fd"},
			},
			Options: mf.ChartOptions{BeginAtZero: true, HideLegend: true},
		})
	})
}
