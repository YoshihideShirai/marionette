package marionette

import shared "github.com/YoshihideShirai/marionette/frontend/shared"

// このファイルはChartコンポーネントのProps/DTO型を定義する。
// グラフ描画に関する型はfrontend/sharedを正として公開する。

type ChartType = shared.ChartType

const (
	ChartTypeBar      = shared.ChartTypeBar
	ChartTypeLine     = shared.ChartTypeLine
	ChartTypePie      = shared.ChartTypePie
	ChartTypeDoughnut = shared.ChartTypeDoughnut
	ChartTypeScatter  = shared.ChartTypeScatter
)

type ChartDataset = shared.ChartDataset
type ChartPoint = shared.ChartPoint
type ChartOptions = shared.ChartOptions
type ChartProps = shared.ChartProps

type chartFallbackRow struct {
	Label  string
	Values []string
}
