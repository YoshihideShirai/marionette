package frontend

import "fmt"

// このファイルはTableコンポーネントのProps/DTO型と補助関数を定義する。
// テーブル表示に関する型・変換ロジックをここに集約する。

type TableColumn struct {
	Label      string
	SortKey    string
	SortHref   string
	SortActive bool
}

type TableComponentRow struct {
	Cells []Node
}

func TableRowValues(values ...any) TableComponentRow {
	cells := make([]Node, 0, len(values))
	for _, value := range values {
		switch v := value.(type) {
		case nil:
			cells = append(cells, textNode(""))
		case Node:
			cells = append(cells, v)
		default:
			cells = append(cells, textNode(fmt.Sprint(v)))
		}
	}
	return TableComponentRow{Cells: cells}
}

type TableProps struct {
	Columns          []TableColumn
	Rows             []TableComponentRow
	EmptyTitle       string
	EmptyDescription string
	View             DataFrameViewProps
	QueryStateName   string
	SelectedFilters  []DataFrameFilter
}
