package goexamples

import (
	"bytes"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

const dataframeCSV = `Name,Role,Score
Aiko,Admin,88
Ren,Editor,72
Mina,Viewer,65
Sora,Admin,95
Yui,Viewer,81
`

func RegisterDataFrameExample(app *mb.App) {
	app.Page("/dataframe", func(ctx *mb.Context) mf.Node {
		q := ctx.Query("q")
		sort := ctx.Query("sort")

		view := mf.DataFrameViewProps{PageSize: 3, ComputedColumns: []mf.DataFrameComputedColumn{{Name: "Tier", Compute: func(row map[string]any) any {
			if row["Score"].(float64) >= 85 {
				return "Gold"
			}
			return "Silver"
		}}}}
		if q != "" {
			view.Filters = append(view.Filters, mf.DataFrameFilter{Column: "Name", Op: mf.DataFrameFilterContains, Value: q})
		}
		if sort != "" {
			view.Sort = append(view.Sort, mf.DataFrameSort{Column: sort})
		}

		node, err := mf.UIDataFrameFromCSV(bytes.NewReader([]byte(dataframeCSV)), mf.TableProps{
			View: view,
			Columns: []mf.TableColumn{
				{Label: "Name", SortKey: "Name", SortHref: "/dataframe?sort=Name&q=" + q},
				{Label: "Role", SortKey: "Role", SortHref: "/dataframe?sort=Role&q=" + q},
				{Label: "Score", SortKey: "Score", SortHref: "/dataframe?sort=Score&q=" + q},
			},
		})
		if err != nil {
			return mf.Alert(mf.AlertProps{Title: "Failed to load CSV", Description: err.Error()})
		}
		return node
	})
}
