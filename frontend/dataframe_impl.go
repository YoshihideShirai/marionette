package frontend

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"

	rdf "github.com/rocketlaunchr/dataframe-go"
)

type DataFrameFilterOp string

const (
	DataFrameFilterEq       DataFrameFilterOp = "eq"
	DataFrameFilterNotEq    DataFrameFilterOp = "neq"
	DataFrameFilterContains DataFrameFilterOp = "contains"
	DataFrameFilterGT       DataFrameFilterOp = "gt"
	DataFrameFilterGTE      DataFrameFilterOp = "gte"
	DataFrameFilterLT       DataFrameFilterOp = "lt"
	DataFrameFilterLTE      DataFrameFilterOp = "lte"
)

type DataFrameFilter struct {
	Column string
	Op     DataFrameFilterOp
	Value  any
}
type DataFrameSort struct {
	Column string
	Desc   bool
}
type DataFrameComputedColumn struct {
	Name    string
	Compute func(row map[string]any) any
}
type DataFrameViewProps struct {
	Filters         []DataFrameFilter
	Sort            []DataFrameSort
	Page            int
	PageSize        int
	ComputedColumns []DataFrameComputedColumn
}

type DataFrameChartSeries struct {
	Column, Label, BackgroundColor, BorderColor string
	Fill                                        bool
	Tension                                     float64
}
type DataFrameChartProps struct {
	Chart       ChartProps
	LabelColumn string
	Series      []DataFrameChartSeries
	View        DataFrameViewProps
}

func ApplyDataFrameView(df *rdf.DataFrame, view DataFrameViewProps) *rdf.DataFrame {
	if df == nil {
		return nil
	}
	columnNames := df.Names()
	rows := make([]map[any]any, 0, df.NRows())
	for i := 0; i < df.NRows(); i++ {
		rows = append(rows, df.Row(i, true, rdf.SeriesName))
	}
	for _, cc := range view.ComputedColumns {
		if cc.Compute == nil || strings.TrimSpace(cc.Name) == "" {
			continue
		}
		name := strings.TrimSpace(cc.Name)
		for _, row := range rows {
			input := map[string]any{}
			for k, v := range row {
				if key, ok := k.(string); ok {
					input[key] = v
				}
			}
			row[name] = cc.Compute(input)
		}
		columnNames = append(columnNames, name)
	}
	if len(view.Filters) > 0 {
		filtered := make([]map[any]any, 0, len(rows))
		for _, row := range rows {
			ok := true
			for _, f := range view.Filters {
				if !matchesFilter(rowValue(row, strings.TrimSpace(f.Column)), f) {
					ok = false
					break
				}
			}
			if ok {
				filtered = append(filtered, row)
			}
		}
		rows = filtered
	}
	if len(view.Sort) > 0 {
		slices.SortStableFunc(rows, func(a, b map[any]any) int {
			for _, s := range view.Sort {
				av, bv := rowValue(a, strings.TrimSpace(s.Column)), rowValue(b, strings.TrimSpace(s.Column))
				c := compareAny(av, bv)
				if c != 0 {
					if s.Desc {
						return -c
					}
					return c
				}
			}
			return 0
		})
	}
	if view.PageSize > 0 {
		page := view.Page
		if page < 1 {
			page = 1
		}
		start := (page - 1) * view.PageSize
		if start >= len(rows) {
			rows = []map[any]any{}
		} else {
			end := min(start+view.PageSize, len(rows))
			rows = rows[start:end]
		}
	}
	series := make([]rdf.Series, 0, len(columnNames))
	for _, name := range columnNames {
		vals := make([]any, 0, len(rows))
		for _, row := range rows {
			vals = append(vals, rowValue(row, name))
		}
		series = append(series, rdf.NewSeriesMixed(name, nil, vals...))
	}
	return rdf.NewDataFrame(series...)
}

// UIDataFrame renders ...
func DataFrame(df *rdf.DataFrame, props TableProps) Node {
	tableProps := props
	df = ApplyDataFrameView(df, tableProps.View)
	if df == nil {
		return Table(tableProps)
	}
	columnNames := df.Names()
	if len(columnNames) > 0 {
		cols := make([]TableColumn, 0, len(columnNames))
		for _, name := range columnNames {
			col := TableColumn{Label: strings.TrimSpace(name)}
			for _, in := range props.Columns {
				if strings.TrimSpace(in.SortKey) == strings.TrimSpace(name) || strings.TrimSpace(in.Label) == strings.TrimSpace(name) {
					col.SortKey = in.SortKey
					col.SortHref = in.SortHref
					col.SortActive = in.SortActive
					break
				}
			}
			cols = append(cols, col)
		}
		tableProps.Columns = cols
	}
	rows := make([]TableComponentRow, 0, df.NRows())
	for i := 0; i < df.NRows(); i++ {
		rowData := df.Row(i, true, rdf.SeriesName)
		cells := make([]Node, 0, len(columnNames))
		for _, name := range columnNames {
			switch v := rowData[name].(type) {
			case nil:
				cells = append(cells, textNode(""))
			case Node:
				cells = append(cells, v)
			default:
				cells = append(cells, textNode(fmt.Sprint(v)))
			}
		}
		rows = append(rows, TableComponentRow{Cells: cells})
	}
	tableProps.Rows = rows
	return Table(tableProps)
}

func DataFrameChart(df *rdf.DataFrame, props DataFrameChartProps) Node {
	chartProps := props.Chart
	df = ApplyDataFrameView(df, props.View)
	if df == nil {
		return UIChart(chartProps)
	}
	columnNames := df.Names()
	if len(columnNames) == 0 {
		return UIChart(chartProps)
	}
	labelColumn := strings.TrimSpace(props.LabelColumn)
	if labelColumn == "" {
		labelColumn = columnNames[0]
	}
	series := props.Series
	if len(series) == 0 {
		for _, name := range columnNames {
			if name == labelColumn {
				continue
			}
			series = append(series, DataFrameChartSeries{Column: name, Label: name})
		}
	}
	labels := make([]string, 0, df.NRows())
	datasets := make([]ChartDataset, len(series))
	for i, item := range series {
		label := strings.TrimSpace(item.Label)
		if label == "" {
			label = strings.TrimSpace(item.Column)
		}
		datasets[i] = ChartDataset{Label: label, BackgroundColor: strings.TrimSpace(item.BackgroundColor), BorderColor: strings.TrimSpace(item.BorderColor), Fill: item.Fill, Tension: item.Tension, Data: make([]float64, 0, df.NRows())}
	}
	for row := 0; row < df.NRows(); row++ {
		rowData := df.Row(row, true, rdf.SeriesName)
		labels = append(labels, fmt.Sprint(rowData[labelColumn]))
		for i, item := range series {
			datasets[i].Data = append(datasets[i].Data, chartFloat(rowData[strings.TrimSpace(item.Column)]))
		}
	}
	chartProps.Labels = labels
	chartProps.Datasets = datasets
	return UIChart(chartProps)
}

func rowValue(row map[any]any, key string) any {
	if value, ok := row[key]; ok {
		return value
	}
	for k, value := range row {
		if ks, ok := k.(string); ok && ks == key {
			return value
		}
	}
	return nil
}

func matchesFilter(value any, f DataFrameFilter) bool {
	op := DataFrameFilterOp(strings.TrimSpace(string(f.Op)))
	if op == "" {
		op = DataFrameFilterEq
	}
	left, right := strings.TrimSpace(fmt.Sprint(value)), strings.TrimSpace(fmt.Sprint(f.Value))
	switch op {
	case DataFrameFilterEq:
		return left == right
	case DataFrameFilterNotEq:
		return left != right
	case DataFrameFilterContains:
		return strings.Contains(strings.ToLower(left), strings.ToLower(right))
	case DataFrameFilterGT, DataFrameFilterGTE, DataFrameFilterLT, DataFrameFilterLTE:
		c := compareAny(value, f.Value)
		if op == DataFrameFilterGT {
			return c > 0
		}
		if op == DataFrameFilterGTE {
			return c >= 0
		}
		if op == DataFrameFilterLT {
			return c < 0
		}
		return c <= 0
	default:
		return left == right
	}
}
func compareAny(a, b any) int {
	af, aok := numericValue(a)
	bf, bok := numericValue(b)
	if aok && bok {
		return cmp.Compare(af, bf)
	}
	return cmp.Compare(strings.TrimSpace(fmt.Sprint(a)), strings.TrimSpace(fmt.Sprint(b)))
}
func numericValue(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint64:
		return float64(n), true
	case uint32:
		return float64(n), true
	default:
		f, err := strconv.ParseFloat(strings.TrimSpace(fmt.Sprint(v)), 64)
		return f, err == nil
	}
}
func chartFloat(value any) float64 { f, _ := numericValue(value); return f }

// UIDataFrameChart is kept for backward compatibility.
func UIDataFrameChart(df *rdf.DataFrame, props DataFrameChartProps) Node {
	return DataFrameChart(df, props)
}
