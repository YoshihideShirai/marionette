package frontend

import (
	"encoding/json"
	"html"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"testing"

	rdf "github.com/rocketlaunchr/dataframe-go"
)

func smallDataFrame() *rdf.DataFrame {
	return rdf.NewDataFrame(
		rdf.NewSeriesMixed("name", nil, "alpha", "beta", "gamma", "delta"),
		rdf.NewSeriesMixed("group", nil, "b", "a", "a", "b"),
		rdf.NewSeriesMixed("score", nil, 10, 2, 30, 4),
		rdf.NewSeriesMixed("score_text", nil, "10", "2", "30", "4"),
		rdf.NewSeriesMixed("note", nil, "Hello red", "blue", "red alert", "green"),
	)
}

func frameColumnValues(t *testing.T, df *rdf.DataFrame, column string) []any {
	t.Helper()
	values := make([]any, 0, df.NRows())
	for i := 0; i < df.NRows(); i++ {
		values = append(values, rowValue(df.Row(i, true, rdf.SeriesName), column))
	}
	return values
}

func frameColumnStrings(t *testing.T, df *rdf.DataFrame, column string) []string {
	t.Helper()
	values := frameColumnValues(t, df, column)
	strings := make([]string, 0, len(values))
	for _, value := range values {
		strings = append(strings, value.(string))
	}
	return strings
}

func TestApplyDataFrameView(t *testing.T) {
	tests := []struct {
		name      string
		view      DataFrameViewProps
		wantNames []string
	}{
		{
			name:      "eq filters exact values",
			view:      DataFrameViewProps{Filters: []DataFrameFilter{{Column: "group", Op: DataFrameFilterEq, Value: "a"}}},
			wantNames: []string{"beta", "gamma"},
		},
		{
			name:      "neq excludes exact values",
			view:      DataFrameViewProps{Filters: []DataFrameFilter{{Column: "group", Op: DataFrameFilterNotEq, Value: "a"}}},
			wantNames: []string{"alpha", "delta"},
		},
		{
			name:      "contains matches case insensitively",
			view:      DataFrameViewProps{Filters: []DataFrameFilter{{Column: "note", Op: DataFrameFilterContains, Value: "RED"}}},
			wantNames: []string{"alpha", "gamma"},
		},
		{
			name:      "gt compares numeric values with numeric strings",
			view:      DataFrameViewProps{Filters: []DataFrameFilter{{Column: "score", Op: DataFrameFilterGT, Value: "9"}}},
			wantNames: []string{"alpha", "gamma"},
		},
		{
			name:      "gte compares numeric-string values with numbers",
			view:      DataFrameViewProps{Filters: []DataFrameFilter{{Column: "score_text", Op: DataFrameFilterGTE, Value: 10}}},
			wantNames: []string{"alpha", "gamma"},
		},
		{
			name:      "lt compares numeric values with numeric strings",
			view:      DataFrameViewProps{Filters: []DataFrameFilter{{Column: "score", Op: DataFrameFilterLT, Value: "10"}}},
			wantNames: []string{"beta", "delta"},
		},
		{
			name:      "lte compares numeric-string values with numbers",
			view:      DataFrameViewProps{Filters: []DataFrameFilter{{Column: "score_text", Op: DataFrameFilterLTE, Value: 4}}},
			wantNames: []string{"beta", "delta"},
		},
		{
			name:      "multiple column sort honors desc for later keys",
			view:      DataFrameViewProps{Sort: []DataFrameSort{{Column: "group"}, {Column: "score", Desc: true}}},
			wantNames: []string{"gamma", "beta", "alpha", "delta"},
		},
		{
			name:      "page below one is normalized to first page",
			view:      DataFrameViewProps{Page: 0, PageSize: 2},
			wantNames: []string{"alpha", "beta"},
		},
		{
			name:      "out of range page returns zero rows",
			view:      DataFrameViewProps{Page: 3, PageSize: 2},
			wantNames: []string{},
		},
		{
			name: "computed columns are added and invalid definitions ignored",
			view: DataFrameViewProps{ComputedColumns: []DataFrameComputedColumn{
				{Name: " doubled ", Compute: func(row map[string]any) any {
					value, _ := numericValue(row["score"])
					return int(value * 2)
				}},
				{Name: "   ", Compute: func(row map[string]any) any { return "ignored" }},
				{Name: "nil_compute"},
			}},
			wantNames: []string{"alpha", "beta", "gamma", "delta"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyDataFrameView(smallDataFrame(), tt.view)
			if got == nil {
				t.Fatal("ApplyDataFrameView() = nil")
			}
			if names := frameColumnStrings(t, got, "name"); !reflect.DeepEqual(names, tt.wantNames) {
				t.Fatalf("name column = %#v, want %#v", names, tt.wantNames)
			}
			if tt.name == "computed columns are added and invalid definitions ignored" {
				if !slices.Contains(got.Names(), "doubled") {
					t.Fatalf("computed column names = %#v, want doubled", got.Names())
				}
				if slices.Contains(got.Names(), "") || slices.Contains(got.Names(), "nil_compute") {
					t.Fatalf("computed column names = %#v, want invalid columns ignored", got.Names())
				}
				if values := frameColumnValues(t, got, "doubled"); !reflect.DeepEqual(values, []any{int64(20), int64(4), int64(60), int64(8)}) {
					t.Fatalf("doubled column = %#v, want %#v", values, []any{int64(20), int64(4), int64(60), int64(8)})
				}
			}
		})
	}
}

func TestDataFrameChartDefaultsAndCoercion(t *testing.T) {
	node := DataFrameChart(smallDataFrame(), DataFrameChartProps{})
	htmlText, err := node.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	config := extractChartConfig(t, string(htmlText))

	if !reflect.DeepEqual(config.Data.Labels, []string{"alpha", "beta", "gamma", "delta"}) {
		t.Fatalf("labels = %#v, want first column values", config.Data.Labels)
	}
	gotLabels := make([]string, 0, len(config.Data.Datasets))
	for _, dataset := range config.Data.Datasets {
		gotLabels = append(gotLabels, dataset.Label)
	}
	if !reflect.DeepEqual(gotLabels, []string{"group", "score", "score_text", "note"}) {
		t.Fatalf("dataset labels = %#v, want non-label columns", gotLabels)
	}
	if !reflect.DeepEqual(config.Data.Datasets[0].Data, []float64{0, 0, 0, 0}) {
		t.Fatalf("non-numeric group data = %#v, want zeros", config.Data.Datasets[0].Data)
	}
	if !reflect.DeepEqual(config.Data.Datasets[1].Data, []float64{10, 2, 30, 4}) {
		t.Fatalf("numeric score data = %#v", config.Data.Datasets[1].Data)
	}
	if !reflect.DeepEqual(config.Data.Datasets[2].Data, []float64{10, 2, 30, 4}) {
		t.Fatalf("numeric-string score_text data = %#v", config.Data.Datasets[2].Data)
	}
	if !reflect.DeepEqual(config.Data.Datasets[3].Data, []float64{0, 0, 0, 0}) {
		t.Fatalf("non-numeric note data = %#v, want zeros", config.Data.Datasets[3].Data)
	}
}

func TestDataFrameChartQueryStateLabelDefaultsFromLabelColumn(t *testing.T) {
	node := DataFrameChart(smallDataFrame(), DataFrameChartProps{
		LabelColumn: " name ",
		Chart:       ChartProps{QueryStateName: "chart-filter"},
	})
	htmlText, err := node.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !regexp.MustCompile(`data-mrn-filter-column="name"`).MatchString(string(htmlText)) {
		t.Fatalf("rendered chart = %s, want QueryStateLabel defaulted from LabelColumn", htmlText)
	}
}

func TestDataQueryStateFromViewAndToView(t *testing.T) {
	view := DataFrameViewProps{
		Filters:  []DataFrameFilter{{Column: "group", Op: DataFrameFilterEq, Value: "a"}},
		Page:     3,
		PageSize: 10,
	}
	state := DataQueryStateFromView(view)
	if !reflect.DeepEqual(state.Filters, view.Filters) {
		t.Fatalf("DataQueryStateFromView().Filters = %#v, want %#v", state.Filters, view.Filters)
	}
	state.Filters[0].Value = "mutated"
	if view.Filters[0].Value == "mutated" {
		t.Fatal("DataQueryStateFromView() did not copy filters")
	}

	base := DataFrameViewProps{Page: 2, PageSize: 25}
	got := state.ToView(base)
	if got.Page != base.Page || got.PageSize != base.PageSize {
		t.Fatalf("ToView() = %#v, want non-filter fields preserved from %#v", got, base)
	}
	if !reflect.DeepEqual(got.Filters, state.Filters) {
		t.Fatalf("ToView().Filters = %#v, want %#v", got.Filters, state.Filters)
	}
	state.Filters[0].Value = "changed again"
	if got.Filters[0].Value == "changed again" {
		t.Fatal("ToView() did not copy filters")
	}
}

func TestDataQueryStateEncode(t *testing.T) {
	values := url.Values{"existing": []string{"kept"}}
	DataQueryState{Filters: []DataFrameFilter{
		{Column: " group ", Op: DataFrameFilterEq, Value: "a"},
		{Column: "score", Op: DataFrameFilterGT, Value: 10},
		{Column: " ", Op: DataFrameFilterLT, Value: 99},
		{Column: "note", Value: "red"},
	}}.Encode(values)

	want := url.Values{
		"existing":  []string{"kept"},
		"df.filter": []string{"group|eq|a", "score|gt|10", "note|eq|red"},
	}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("values = %#v, want %#v", values, want)
	}
}

type chartConfig struct {
	Data struct {
		Labels   []string `json:"labels"`
		Datasets []struct {
			Label string    `json:"label"`
			Data  []float64 `json:"data"`
		} `json:"datasets"`
	} `json:"data"`
}

func extractChartConfig(t *testing.T, rendered string) chartConfig {
	t.Helper()
	matches := regexp.MustCompile(`(?s)<script type="application/json" data-mrn-chart-config>(.*?)</script>`).FindStringSubmatch(rendered)
	if len(matches) != 2 {
		t.Fatalf("rendered chart = %s, want chart config script", rendered)
	}
	var config chartConfig
	if err := json.Unmarshal([]byte(html.UnescapeString(matches[1])), &config); err != nil {
		t.Fatalf("json.Unmarshal(%q) error = %v", matches[1], err)
	}
	return config
}
