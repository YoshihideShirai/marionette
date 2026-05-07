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
				{Label: "Users", Data: []float64{3, 7, 12}, BackgroundColors: []string{"#2563eb", "#14b8a6", "#f59e0b"}},
			},
			Options: mf.ChartOptions{BeginAtZero: true, HideLegend: true},
		})
	})
}
