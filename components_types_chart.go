package marionette

import mf "github.com/YoshihideShirai/marionette/frontend"

// このファイルはChartコンポーネントのProps/DTO型を定義する。
// グラフ描画に関する型をここに集約する。

type ChartType = mf.ChartType

const (
	ChartTypeBar      = mf.ChartTypeBar
	ChartTypeLine     = mf.ChartTypeLine
	ChartTypePie      = mf.ChartTypePie
	ChartTypeDoughnut = mf.ChartTypeDoughnut
	ChartTypeScatter  = mf.ChartTypeScatter
)

type ChartDataset = mf.ChartDataset
type ChartPoint = mf.ChartPoint
type ChartOptions = mf.ChartOptions
type ChartProps = mf.ChartProps

type chartFallbackRow struct {
	Label  string
	Values []string
}
