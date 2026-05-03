package frontend

import chartjs "github.com/YoshihideShirai/marionette/frontend/chartjs"

// このファイルはChartコンポーネントのProps/DTO型を定義する。
// グラフ描画に関する型をここに集約する。

type ChartType = chartjs.ChartType

const (
	ChartTypeBar      = chartjs.ChartTypeBar
	ChartTypeLine     = chartjs.ChartTypeLine
	ChartTypePie      = chartjs.ChartTypePie
	ChartTypeDoughnut = chartjs.ChartTypeDoughnut
	ChartTypeScatter  = chartjs.ChartTypeScatter
)

type ChartDataset = chartjs.ChartDataset

type ChartPoint = chartjs.ChartPoint

type ChartOptions = chartjs.ChartOptions

type ChartProps = chartjs.ChartProps

type chartFallbackRow struct {
	Label  string
	Values []string
}
