package chartjs

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestConfigJSONDefaultsToLineChart(t *testing.T) {
	config := mustConfigMap(t, ChartProps{
		Labels:   []string{"Jan", "Feb"},
		Datasets: []ChartDataset{{Label: "Revenue", Data: []float64{10, 20}}},
	})

	if got := config["type"]; got != string(ChartTypeLine) {
		t.Fatalf("type = %v, want %q", got, ChartTypeLine)
	}

	data := mustMap(t, config["data"])
	if got, want := data["labels"], []any{"Jan", "Feb"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("labels = %#v, want %#v", got, want)
	}

	datasets := mustSlice(t, data["datasets"])
	if len(datasets) != 1 {
		t.Fatalf("datasets length = %d, want 1", len(datasets))
	}
	dataset := mustMap(t, datasets[0])
	if got := dataset["label"]; got != "Revenue" {
		t.Fatalf("dataset label = %v, want Revenue", got)
	}
	if got, want := dataset["data"], []any{float64(10), float64(20)}; !reflect.DeepEqual(got, want) {
		t.Fatalf("dataset data = %#v, want %#v", got, want)
	}
}

func TestConfigJSONAppliesDefaultColorsToMultiValueCircularAndBarCharts(t *testing.T) {
	tests := []struct {
		name            string
		chartType       ChartType
		wantBorderColor []any
	}{
		{name: "bar", chartType: ChartTypeBar, wantBorderColor: []any{"#2563eb", "#14b8a6", "#f59e0b"}},
		{name: "pie", chartType: ChartTypePie, wantBorderColor: []any{"#ffffff", "#ffffff", "#ffffff"}},
		{name: "doughnut", chartType: ChartTypeDoughnut, wantBorderColor: []any{"#ffffff", "#ffffff", "#ffffff"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := mustConfigMap(t, ChartProps{
				Type:   tt.chartType,
				Labels: []string{"A", "B", "C"},
				Datasets: []ChartDataset{{
					Label: "Series",
					Data:  []float64{1, 2, 3},
				}},
			})

			dataset := firstDataset(t, config)
			wantBackgroundColor := []any{"#2563eb", "#14b8a6", "#f59e0b"}
			if got := dataset["backgroundColor"]; !reflect.DeepEqual(got, wantBackgroundColor) {
				t.Fatalf("backgroundColor = %#v, want %#v", got, wantBackgroundColor)
			}
			if got := dataset["borderColor"]; !reflect.DeepEqual(got, tt.wantBorderColor) {
				t.Fatalf("borderColor = %#v, want %#v", got, tt.wantBorderColor)
			}
		})
	}
}

func TestConfigJSONReflectsOptionsAndScales(t *testing.T) {
	config := mustConfigMap(t, ChartProps{
		Type:     ChartTypeBar,
		Labels:   []string{"Q1"},
		Datasets: []ChartDataset{{Label: "Revenue", Data: []float64{42}}},
		Options: ChartOptions{
			HideLegend:  true,
			AspectRatio: 1.5,
			BeginAtZero: true,
			Stacked:     true,
			XAxisLabel:  " Quarter ",
			YAxisLabel:  " Revenue ",
		},
	})

	options := mustMap(t, config["options"])
	if got := options["maintainAspectRatio"]; got != true {
		t.Fatalf("maintainAspectRatio = %v, want true", got)
	}
	if got := options["aspectRatio"]; got != 1.5 {
		t.Fatalf("aspectRatio = %v, want 1.5", got)
	}

	plugins := mustMap(t, options["plugins"])
	legend := mustMap(t, plugins["legend"])
	if got := legend["display"]; got != false {
		t.Fatalf("legend display = %v, want false", got)
	}

	assertScales(t, mustMap(t, options["scales"]))
	assertScales(t, Scales(ChartOptions{
		BeginAtZero: true,
		Stacked:     true,
		XAxisLabel:  " Quarter ",
		YAxisLabel:  " Revenue ",
	}))
}

func TestConfigJSONUsesPointsBeforeDatasetData(t *testing.T) {
	config := mustConfigMap(t, ChartProps{
		Type:   ChartTypeScatter,
		Labels: []string{"first", "second"},
		Datasets: []ChartDataset{{
			Label:  "Coordinates",
			Data:   []float64{100, 200},
			Points: []ChartPoint{{X: 1.5, Y: 2.5}, {X: 3.5, Y: 4.5}},
		}},
	})

	points := mustSlice(t, firstDataset(t, config)["data"])
	if len(points) != 2 {
		t.Fatalf("points length = %d, want 2", len(points))
	}
	first := mustMap(t, points[0])
	if got := first["x"]; got != 1.5 {
		t.Fatalf("first point x = %v, want 1.5", got)
	}
	if got := first["y"]; got != 2.5 {
		t.Fatalf("first point y = %v, want 2.5", got)
	}
}

func TestFallbackText(t *testing.T) {
	if got, want := FallbackText(ChartProps{}), "Chart data is available in the fallback table below."; got != want {
		t.Fatalf("FallbackText() = %q, want %q", got, want)
	}
	if got, want := FallbackText(ChartProps{Title: " Sales "}), "Sales data is available in the fallback table below."; got != want {
		t.Fatalf("FallbackText() = %q, want %q", got, want)
	}
}

func TestFallbackRowsFormatsMissingDataAndPoints(t *testing.T) {
	rows := FallbackRows(ChartProps{
		Labels: []string{"A", "B", "C"},
		Datasets: []ChartDataset{
			{Label: "Data", Data: []float64{10, 20}},
			{Label: "Points", Data: []float64{100, 200, 300}, Points: []ChartPoint{{X: 1, Y: 2}, {X: 3.5, Y: 4.25}}},
		},
	})

	want := []FallbackRow{
		{Label: "A", Values: []string{"10", "1, 2"}},
		{Label: "B", Values: []string{"20", "3.5, 4.25"}},
		{Label: "C", Values: []string{"", "300"}},
	}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("FallbackRows() = %#v, want %#v", rows, want)
	}
}

func mustConfigMap(t *testing.T, props ChartProps) map[string]any {
	t.Helper()
	configJSON, err := ConfigJSON(props)
	if err != nil {
		t.Fatalf("ConfigJSON() error = %v", err)
	}
	var config map[string]any
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		t.Fatalf("json.Unmarshal(ConfigJSON()) error = %v; json = %s", err, configJSON)
	}
	return config
}

func firstDataset(t *testing.T, config map[string]any) map[string]any {
	t.Helper()
	data := mustMap(t, config["data"])
	datasets := mustSlice(t, data["datasets"])
	if len(datasets) == 0 {
		t.Fatal("datasets length = 0, want at least 1")
	}
	return mustMap(t, datasets[0])
}

func assertScales(t *testing.T, scales map[string]any) {
	t.Helper()
	x := mustMap(t, scales["x"])
	y := mustMap(t, scales["y"])

	if got := x["stacked"]; got != true {
		t.Fatalf("x.stacked = %v, want true", got)
	}
	if got := y["stacked"]; got != true {
		t.Fatalf("y.stacked = %v, want true", got)
	}
	if got := y["beginAtZero"]; got != true {
		t.Fatalf("y.beginAtZero = %v, want true", got)
	}
	assertTitle(t, x["title"], "Quarter")
	assertTitle(t, y["title"], "Revenue")
}

func assertTitle(t *testing.T, raw any, text string) {
	t.Helper()
	title := mustMap(t, raw)
	if got := title["display"]; got != true {
		t.Fatalf("title.display = %v, want true", got)
	}
	if got := title["text"]; got != text {
		t.Fatalf("title.text = %v, want %q", got, text)
	}
}

func mustMap(t *testing.T, raw any) map[string]any {
	t.Helper()
	value, ok := raw.(map[string]any)
	if !ok {
		t.Fatalf("value = %#v (%T), want map[string]any", raw, raw)
	}
	return value
}

func mustSlice(t *testing.T, raw any) []any {
	t.Helper()
	value, ok := raw.([]any)
	if !ok {
		t.Fatalf("value = %#v (%T), want []any", raw, raw)
	}
	return value
}
