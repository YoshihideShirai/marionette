package frontend

import (
	"fmt"
	"strings"

	rdf "github.com/rocketlaunchr/dataframe-go"
)

// ComponentDataFrame renders a github.com/rocketlaunchr/dataframe-go DataFrame
// as a table component.
//
// Values are converted as follows:
//   - nil => empty text cell
//   - Node => rendered directly
//   - all others => fmt.Sprint(value) wrapped in Text(...)
func ComponentDataFrame(df *rdf.DataFrame, props TableProps) Node {
	tableProps := props
	if df == nil {
		return ComponentTable(tableProps)
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
	return ComponentTable(tableProps)
}
