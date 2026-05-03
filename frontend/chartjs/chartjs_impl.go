package chartjs

import (
	"encoding/json"
	"fmt"
	"strings"

	shared "github.com/YoshihideShirai/marionette/frontend/shared"
)

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

type FallbackRow struct {
	Label  string
	Values []string
}

func ConfigJSON(props ChartProps) (string, error) {
	chartType := strings.TrimSpace(string(props.Type))
	if chartType == "" {
		chartType = string(ChartTypeLine)
	}
	datasets := make([]map[string]any, 0, len(props.Datasets))
	for _, dataset := range props.Datasets {
		data := any(dataset.Data)
		if len(dataset.Points) > 0 {
			data = dataset.Points
		}
		item := map[string]any{"label": strings.TrimSpace(dataset.Label), "data": data}
		if color := strings.TrimSpace(dataset.BackgroundColor); color != "" {
			item["backgroundColor"] = color
		}
		if color := strings.TrimSpace(dataset.BorderColor); color != "" {
			item["borderColor"] = color
		}
		if dataset.Fill {
			item["fill"] = true
		}
		if dataset.Tension > 0 {
			item["tension"] = dataset.Tension
		}
		datasets = append(datasets, item)
	}
	options := map[string]any{
		"responsive":          true,
		"maintainAspectRatio": false,
		"plugins":             map[string]any{"legend": map[string]any{"display": !props.Options.HideLegend}},
	}
	if props.Options.AspectRatio > 0 {
		options["maintainAspectRatio"] = true
		options["aspectRatio"] = props.Options.AspectRatio
	}
	if chartType != string(ChartTypePie) && chartType != string(ChartTypeDoughnut) {
		options["scales"] = Scales(props.Options)
	}
	payload := map[string]any{"type": chartType, "data": map[string]any{"labels": props.Labels, "datasets": datasets}, "options": options}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Scales(options ChartOptions) map[string]any {
	x := map[string]any{}
	y := map[string]any{"beginAtZero": options.BeginAtZero}
	if options.Stacked {
		x["stacked"] = true
		y["stacked"] = true
	}
	if label := strings.TrimSpace(options.XAxisLabel); label != "" {
		x["title"] = map[string]any{"display": true, "text": label}
	}
	if label := strings.TrimSpace(options.YAxisLabel); label != "" {
		y["title"] = map[string]any{"display": true, "text": label}
	}
	return map[string]any{"x": x, "y": y}
}

func FallbackText(props ChartProps) string {
	title := strings.TrimSpace(props.Title)
	if title == "" {
		title = "Chart"
	}
	return title + " data is available in the fallback table below."
}

func FallbackRows(props ChartProps) []FallbackRow {
	rows := make([]FallbackRow, 0, len(props.Labels))
	for i, label := range props.Labels {
		values := make([]string, 0, len(props.Datasets))
		for _, dataset := range props.Datasets {
			if i < len(dataset.Points) {
				values = append(values, fmt.Sprintf("%g, %g", dataset.Points[i].X, dataset.Points[i].Y))
				continue
			}
			if i >= len(dataset.Data) {
				values = append(values, "")
				continue
			}
			values = append(values, fmt.Sprint(dataset.Data[i]))
		}
		rows = append(rows, FallbackRow{Label: label, Values: values})
	}
	return rows
}
