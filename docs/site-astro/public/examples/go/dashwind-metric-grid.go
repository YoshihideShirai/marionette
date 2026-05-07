package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

func RegisterDashwindMetricGridExample(app *mb.App) {
	app.Page("/dashwind-metric-grid", func(ctx *mb.Context) mf.Node {
		return dw.MetricGrid(dw.MetricGridProps{Items: []dw.Metric{
			{Title: "Revenue", Value: "$128.4k", Description: "This month", Trend: "+12.5% vs last month", TrendTone: dw.ToneSuccess, Icon: mf.Text("💸")},
			{Title: "Open tickets", Value: "37", Description: "Across support queues", Trend: "5 need review", TrendTone: dw.ToneWarning, Icon: mf.Text("🎫")},
			{Title: "Churn risk", Value: "4", Description: "Enterprise accounts", Trend: "2 escalated", TrendTone: dw.ToneError, Icon: mf.Text("⚠")},
		}})
	})
}
