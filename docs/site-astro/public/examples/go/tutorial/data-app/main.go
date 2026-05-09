package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	rdf "github.com/rocketlaunchr/dataframe-go"
	dataframeimports "github.com/rocketlaunchr/dataframe-go/imports"
)

const salesCSV = `Region,Sales,Orders
North,1200,18
South,950,12
East,1750,23
West,1330,16
`

func main() {
	app := buildApp()
	if err := app.Run("127.0.0.1:8084"); err != nil {
		log.Fatal(err)
	}
}

func buildApp() *mb.App {
	app := mb.New()
	app.Page("/", func(ctx *mb.Context) mf.Node {
		return renderDashboard(ctx)
	}, mb.WithTitle("Data App Tutorial"))
	app.Action("filters/apply", func(ctx *mb.Context) mf.Node {
		return dashboardPanel(ctx)
	})
	return app
}

func renderDashboard(ctx *mb.Context) mf.Node {
	return mf.Container(mf.ContainerProps{MaxWidth: "5xl", Centered: true},
		mf.Stack(mf.StackProps{Direction: "column", Gap: "6"},
			mf.PageHeader(mf.PageHeaderProps{
				Title:       "Data App Tutorial",
				Description: "CSV load, shared DataFrame view, synchronized table and chart, and action rerender.",
			}),
			dashboardPanel(ctx),
		),
	)
}

func dashboardPanel(ctx *mb.Context) mf.Node {
	minSales := minSalesFromContext(ctx)
	view := salesView(minSales)

	table, err := mf.DataFrameFromCSV(bytes.NewReader([]byte(salesCSV)), mf.TableProps{
		View: view,
		Columns: []mf.TableColumn{
			{Label: "Region", SortKey: "Region"},
			{Label: "Sales", SortKey: "Sales", SortActive: true},
			{Label: "Orders", SortKey: "Orders"},
		},
	})
	if err != nil {
		return mf.Alert(mf.AlertProps{Title: "Failed to load CSV", Description: err.Error()})
	}

	df, err := loadSalesDataFrame()
	if err != nil {
		return mf.Alert(mf.AlertProps{Title: "Failed to load chart data", Description: err.Error()})
	}

	return mf.Region(mf.RegionProps{ID: "sales-dashboard"},
		mf.Stack(mf.StackProps{Direction: "column", Gap: "4"},
			mf.ActionForm(mf.ActionFormProps{
				Action: "/filters/apply",
				Target: "#sales-dashboard",
				Swap:   "outerHTML",
				Props:  mf.ComponentProps{Class: "flex items-end gap-3"},
			},
				mf.FormRow(mf.FormRowProps{
					ID:    "min-sales",
					Label: "Minimum sales",
					Control: mf.TextField(mf.TextFieldProps{
						ID:    "min-sales",
						Name:  "min_sales",
						Value: fmt.Sprint(minSales),
					}),
				}),
				mf.SubmitButton("Apply filters", mf.ComponentProps{Variant: "primary"}),
			),
			mf.Grid(mf.GridProps{Columns: "2", Gap: "4"},
				mf.Card(mf.CardProps{Title: "Regional sales table"}, table),
				mf.Card(mf.CardProps{Title: "Regional sales chart"},
					mf.DataFrameChart(df, mf.DataFrameChartProps{
						View:        view,
						Chart:       mf.ChartProps{Type: mf.ChartTypeBar, Title: "Sales by region"},
						LabelColumn: "Region",
						Series:      []mf.DataFrameChartSeries{{Column: "Sales", Label: "Sales"}},
					}),
				),
			),
		),
	)
}

func minSalesFromContext(ctx *mb.Context) int {
	if ctx == nil {
		return 1000
	}
	value := ctx.FormValue("min_sales")
	if value == "" {
		value = ctx.Query("min_sales")
	}
	if value == "" {
		return 1000
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 1000
	}
	return parsed
}

func salesView(minSales int) mf.DataFrameViewProps {
	return mf.DataFrameViewProps{
		Filters: []mf.DataFrameFilter{{Column: "Sales", Op: mf.DataFrameFilterGTE, Value: float64(minSales)}},
		Sort:    []mf.DataFrameSort{{Column: "Sales", Desc: true}},
	}
}

func loadSalesDataFrame() (*rdf.DataFrame, error) {
	return dataframeimports.LoadFromCSV(context.Background(), bytes.NewReader([]byte(salesCSV)))
}
