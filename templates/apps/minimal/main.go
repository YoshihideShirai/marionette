package main

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

func main() {
	app := mb.New()
	dw.Use(app, dw.Options{Theme: "corporate"})
	app.Page("/", homePage, mb.WithTitle("Minimal Marionette App"))

	if err := app.Run("127.0.0.1:8090"); err != nil {
		panic(err)
	}
}

func homePage(ctx *mb.Context) mf.Node {
	body := mf.DivProps(mf.ElementProps{Class: "space-y-6"},
		dw.PageHeader(dw.PageHeaderProps{
			Title:       "Minimal App",
			Description: "A polished single-page starting point with shell, metrics, and content panels.",
		}),
		dw.MetricGrid(dw.MetricGridProps{Items: []dw.Metric{
			{Title: "Pages", Value: "1", Description: "Root route", Trend: "Ready to extend", TrendTone: dw.ToneSuccess, Icon: metricIcon("P")},
			{Title: "Actions", Value: "0", Description: "Add POST handlers", Trend: "Use htmx regions", TrendTone: dw.ToneNeutral, Icon: metricIcon("A")},
			{Title: "State", Value: "Go", Description: "Server-owned", Trend: "Simple by default", TrendTone: dw.ToneSuccess, Icon: metricIcon("S")},
			{Title: "Theme", Value: "DW", Description: "DashWind shell", Trend: "Production shape", TrendTone: dw.ToneNeutral, Icon: metricIcon("D")},
		}}),
		mf.Grid(mf.GridProps{Columns: "2", Gap: "lg"},
			dw.CardPanel(dw.CardPanelProps{Title: "Main workflow", Description: "Replace this panel with the first job your app performs."},
				mf.DivProps(mf.ElementProps{Class: "space-y-4"},
					mf.PProps(mf.ElementProps{Class: "text-sm text-base-content/70"}, mf.Text("Start from one page handler, then add action endpoints only where the UI needs partial updates.")),
					mf.Actions(mf.ActionsProps{Align: "start", Wrap: true},
						mf.Button("Primary action", mf.ComponentProps{Class: "btn-primary"}),
						mf.Button("Secondary", mf.ComponentProps{Class: "btn-outline"}),
					),
				),
			),
			dw.CardPanel(dw.CardPanelProps{Title: "Implementation notes", Description: "Useful defaults for a new Marionette app."},
				mf.UlProps(mf.ElementProps{Class: "list-disc space-y-2 pl-5 text-sm text-base-content/70"},
					mf.Li(mf.Text("Keep page rendering, state transitions, and htmx fragments in Go.")),
					mf.Li(mf.Text("Use regions when only part of the page should refresh.")),
					mf.Li(mf.Text("Move repeated screen patterns into small helper functions.")),
				),
			),
		),
	)
	return appShell("/", body)
}

func appShell(currentPath string, body mf.Node) mf.Node {
	return dw.Shell(dw.ShellProps{
		Brand:             dw.Brand{Title: "Starter", Subtitle: "Minimal template", Mark: "M"},
		CurrentPath:       currentPath,
		SearchPlaceholder: "Search workspace",
		Navigation: []dw.NavGroup{{
			Label: "Template",
			Items: []dw.NavItem{{Path: "/", Label: "Overview", Icon: "O"}},
		}},
		User: dw.UserMenu{Name: "Template User", Email: "user@example.com", Initials: "TU"},
	}, body)
}

func metricIcon(label string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "grid h-10 w-10 place-items-center rounded-box bg-primary/10 font-bold text-primary"}, mf.Text(label))
}
