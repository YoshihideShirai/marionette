package marionette

import (
	"net/url"

	mf "github.com/YoshihideShirai/marionette/frontend"
	rdf "github.com/rocketlaunchr/dataframe-go"
)

type DataFrameFilterOp = mf.DataFrameFilterOp

const (
	DataFrameFilterEq       = mf.DataFrameFilterEq
	DataFrameFilterContains = mf.DataFrameFilterContains
	DataFrameFilterGT       = mf.DataFrameFilterGT
	DataFrameFilterGTE      = mf.DataFrameFilterGTE
	DataFrameFilterLT       = mf.DataFrameFilterLT
	DataFrameFilterLTE      = mf.DataFrameFilterLTE
)

type DataFrameFilter = mf.DataFrameFilter
type DataFrameSort = mf.DataFrameSort
type DataFrameComputedColumn = mf.DataFrameComputedColumn
type DataFrameViewProps = mf.DataFrameViewProps
type DataQueryState = mf.DataQueryState
type DataFrameChartSeries = mf.DataFrameChartSeries
type DataFrameChartProps = mf.DataFrameChartProps

func ApplyDataFrameView(df *rdf.DataFrame, view DataFrameViewProps) *rdf.DataFrame {
	return mf.ApplyDataFrameView(df, view)
}

func DataFrame(df *rdf.DataFrame, props TableProps) Node {
	return mf.DataFrame(df, props)
}

func DataFrameChart(df *rdf.DataFrame, props DataFrameChartProps) Node {
	return mf.DataFrameChart(df, props)
}

func DataQueryStateFromView(view DataFrameViewProps) DataQueryState {
	return mf.DataQueryStateFromView(view)
}

func encodeDataQueryState(state DataQueryState, values url.Values) {
	state.Encode(values)
}
