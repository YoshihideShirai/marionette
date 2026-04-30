package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterChartExample(app *mb.App) {
	app.Page("/chart", func(ctx *mb.Context) mf.Node {
		return mf.Chart(mf.ChartProps{
			Type:        mf.ChartTypeLine,
			Title:       "Weekly signups",
			Description: "New accounts by weekday.",
			Labels:      []string{"Mon", "Tue", "Wed"},
			Datasets: []mf.ChartDataset{
				{
					Label:           "Signups",
					Data:            []float64{12, 19, 14},
					BorderColor:     "#2563eb",
					BackgroundColor: "rgba(37, 99, 235, 0.16)",
					Fill:            true,
					Tension:         0.3,
				},
			},
			Options: mf.ChartOptions{BeginAtZero: true, YAxisLabel: "Users"},
			Height:  260,
		})
	})
}
