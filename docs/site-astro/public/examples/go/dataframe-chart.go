package goexamples

import (
	"bytes"
	"context"
	"strconv"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	rdf "github.com/rocketlaunchr/dataframe-go"
	dataframeimports "github.com/rocketlaunchr/dataframe-go/imports"
)

const dataframeChartTSV = `Name	Score	Tasks
Aiko	88	12
Ren	72	8
Mina	65	6
Sora	95	15
Yui	81	10
`

func RegisterDataFrameChartExample(app *mb.App) {
	app.Page("/dataframe-chart", func(ctx *mb.Context) mf.Node {
		minScore, _ := strconv.Atoi(ctx.Query("min_score"))
		df, err := dataframeimports.LoadFromCSV(context.Background(), bytes.NewReader([]byte(dataframeChartTSV)), dataframeimports.CSVLoadOptions{Comma: '\t'})
		if err != nil {
			return mf.Alert(mf.AlertProps{Title: "Failed to load TSV", Description: err.Error()})
		}

		return renderDataFrameChart(df, minScore)
	})
}

func renderDataFrameChart(df *rdf.DataFrame, minScore int) mf.Node {
	view := mf.DataFrameViewProps{Sort: []mf.DataFrameSort{{Column: "Score", Desc: true}}}
	if minScore > 0 {
		view.Filters = append(view.Filters, mf.DataFrameFilter{Column: "Score", Op: mf.DataFrameFilterGTE, Value: float64(minScore)})
	}
	return mf.DataFrameChart(df, mf.DataFrameChartProps{View: view, Chart: mf.ChartProps{Type: mf.ChartTypeBar, Title: "DataFrame chart"}, LabelColumn: "Name", Series: []mf.DataFrameChartSeries{{Column: "Score"}, {Column: "Tasks"}}})
}
