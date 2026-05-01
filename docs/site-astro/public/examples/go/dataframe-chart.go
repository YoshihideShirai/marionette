package goexamples

import (
	"strconv"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	rdf "github.com/rocketlaunchr/dataframe-go"
)

func RegisterDataFrameChartExample(app *mb.App) {
	app.Page("/dataframe-chart", func(ctx *mb.Context) mf.Node {
		minScore, _ := strconv.Atoi(ctx.Query("min_score"))
		df := rdf.NewDataFrame(
			rdf.NewSeriesString("Name", nil, "Aiko", "Ren", "Mina", "Sora", "Yui"),
			rdf.NewSeriesInt64("Score", nil, 88, 72, 65, 95, 81),
			rdf.NewSeriesInt64("Tasks", nil, 12, 8, 6, 15, 10),
		)
		view := mf.DataFrameViewProps{Sort: []mf.DataFrameSort{{Column: "Score", Desc: true}}}
		if minScore > 0 {
			view.Filters = append(view.Filters, mf.DataFrameFilter{Column: "Score", Op: mf.DataFrameFilterGTE, Value: minScore})
		}
		return mf.DataFrameChart(df, mf.DataFrameChartProps{View: view, Chart: mf.ChartProps{Type: mf.ChartTypeBar, Title: "DataFrame chart"}, LabelColumn: "Name", Series: []mf.DataFrameChartSeries{{Column: "Score"}, {Column: "Tasks"}}})
	})
}
