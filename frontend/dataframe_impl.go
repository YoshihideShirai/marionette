package frontend

import (
	"fmt"
	"strconv"
	"strings"

	rdf "github.com/rocketlaunchr/dataframe-go"
)

type DataFrameChartSeries struct {
	Column          string
	Label           string
	BackgroundColor string
	BorderColor     string
	Fill            bool
	Tension         float64
}

type DataFrameChartProps struct {
	Chart       ChartProps
	LabelColumn string
	Series      []DataFrameChartSeries
}

// UIDataFrame renders a github.com/rocketlaunchr/dataframe-go DataFrame
// as a table component.
//
// Values are converted as follows:
//   - nil => empty text cell
//   - Node => rendered directly
//   - all others => fmt.Sprint(value) wrapped in Text(...)
func UIDataFrame(df *rdf.DataFrame, props TableProps) Node {
	tableProps := props
	if df == nil {
		return UITable(tableProps)
	}

	columnNames := df.Names()
	if len(columnNames) > 0 {
		columns := make([]TableColumn, 0, len(columnNames))
		for _, name := range columnNames {
			columns = append(columns, TableColumn{Label: strings.TrimSpace(name)})
		}
		tableProps.Columns = columns
	}

	rows := make([]TableComponentRow, 0, df.NRows())
	for i := 0; i < df.NRows(); i++ {
		rowData := df.Row(i, true, rdf.SeriesName)
		cells := make([]Node, 0, len(columnNames))
		for _, name := range columnNames {
			value := rowData[name]
			switch v := value.(type) {
			case nil:
				cells = append(cells, Text(""))
			case Node:
				cells = append(cells, v)
			default:
				cells = append(cells, Text(fmt.Sprint(v)))
			}
		}
		rows = append(rows, TableComponentRow{Cells: cells})
	}

	tableProps.Rows = rows
	return UITable(tableProps)
}

// UIDataFrameChart renders dataframe columns through UIChart.
//
// LabelColumn selects the x-axis labels. If blank, the first dataframe column is
// used. Series selects numeric columns. If blank, every column after LabelColumn
// is rendered as a dataset.
func UIDataFrameChart(df *rdf.DataFrame, props DataFrameChartProps) Node {
	chartProps := props.Chart
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
		datasets[i] = ChartDataset{
			Label:           label,
			BackgroundColor: strings.TrimSpace(item.BackgroundColor),
			BorderColor:     strings.TrimSpace(item.BorderColor),
			Fill:            item.Fill,
			Tension:         item.Tension,
			Data:            make([]float64, 0, df.NRows()),
		}
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

func chartFloat(value any) float64 {
	switch v := value.(type) {
	case nil:
		return 0
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	case uint:
		return float64(v)
	case uint64:
		return float64(v)
	case uint32:
		return float64(v)
	default:
		n, err := strconv.ParseFloat(strings.TrimSpace(fmt.Sprint(v)), 64)
		if err != nil {
			return 0
		}
		return n
	}
}
