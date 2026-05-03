package frontend

// このファイルはChartコンポーネントのProps/DTO型を定義する。
// グラフ描画に関する型をここに集約する。

type ChartType string

const (
	ChartTypeBar      ChartType = "bar"
	ChartTypeLine     ChartType = "line"
	ChartTypePie      ChartType = "pie"
	ChartTypeDoughnut ChartType = "doughnut"
	ChartTypeScatter  ChartType = "scatter"
)

type ChartDataset struct {
	Label           string
	Data            []float64
	Points          []ChartPoint
	BackgroundColor string
	BorderColor     string
	Fill            bool
	Tension         float64
}

type ChartPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ChartOptions struct {
	BeginAtZero bool
	Stacked     bool
	HideLegend  bool
	AspectRatio float64
	XAxisLabel  string
	YAxisLabel  string
}

type ChartProps struct {
	Type            ChartType
	Title           string
	Description     string
	Labels          []string
	Datasets        []ChartDataset
	Options         ChartOptions
	AriaLabel       string
	Height          int
	Props           ComponentProps
	QueryStateName  string
	QueryStateLabel string
}

type chartFallbackRow struct {
	Label  string
	Values []string
}
