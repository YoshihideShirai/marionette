package goexamples

import (
	"strconv"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	rdf "github.com/rocketlaunchr/dataframe-go"
)

func RegisterDataFrameExample(app *mb.App) {
	app.Page("/dataframe", func(ctx *mb.Context) mf.Node {
		q := ctx.Query("q")
		sort := ctx.Query("sort")
		page, _ := strconv.Atoi(ctx.Query("page"))
		df := rdf.NewDataFrame(
			rdf.NewSeriesString("Name", nil, "Aiko", "Ren", "Mina", "Sora", "Yui"),
			rdf.NewSeriesString("Role", nil, "Admin", "Editor", "Viewer", "Admin", "Viewer"),
			rdf.NewSeriesInt64("Score", nil, 88, 72, 65, 95, 81),
		)
		view := mf.DataFrameViewProps{Page: page, PageSize: 3, ComputedColumns: []mf.DataFrameComputedColumn{{Name: "Tier", Compute: func(row map[string]any) any {
			if row["Score"].(int64) >= 85 {
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
		return mf.DataFrame(df, mf.TableProps{View: view, Columns: []mf.TableColumn{{Label: "Name", SortKey: "Name", SortHref: "/dataframe?sort=Name&q=" + q}, {Label: "Role", SortKey: "Role", SortHref: "/dataframe?sort=Role&q=" + q}, {Label: "Score", SortKey: "Score", SortHref: "/dataframe?sort=Score&q=" + q}}})
	})
}
