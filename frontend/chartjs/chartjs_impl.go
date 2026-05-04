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
		} else if colors := defaultColors(chartType, len(dataset.Data)); len(colors) > 0 {
			item["backgroundColor"] = colors
		}
		if color := strings.TrimSpace(dataset.BorderColor); color != "" {
			item["borderColor"] = color
		} else if colors := defaultBorderColors(chartType, len(dataset.Data)); len(colors) > 0 {
			item["borderColor"] = colors
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

func defaultColors(chartType string, count int) []string {
	if count <= 1 {
		return nil
	}
	switch chartType {
	case string(ChartTypeBar), string(ChartTypePie), string(ChartTypeDoughnut):
		palette := []string{"#2563eb", "#14b8a6", "#f59e0b", "#8b5cf6", "#ef4444", "#22c55e"}
		if count > len(palette) {
			count = len(palette)
		}
		return palette[:count]
	default:
		return nil
	}
}

func defaultBorderColors(chartType string, count int) []string {
	if count <= 1 {
		return nil
	}
	switch chartType {
	case string(ChartTypePie), string(ChartTypeDoughnut):
		colors := make([]string, count)
		for i := range colors {
			colors[i] = "#ffffff"
		}
		return colors
	default:
		return nil
	}
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
