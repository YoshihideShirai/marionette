package main

import (
	"fmt"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

type activity struct {
	Time   string
	Event  string
	Owner  string
	Status string
}

func main() {
	app := mb.New()
	dw.Use(app, dw.Options{Theme: "corporate"})
	app.Page("/", dashboardPage, mb.WithTitle("Dashboard Template"))

	if err := app.Run("127.0.0.1:8092"); err != nil {
		panic(err)
	}
}

func dashboardPage(ctx *mb.Context) mf.Node {
	body := mf.DivProps(mf.ElementProps{Class: "space-y-6"},
		dw.PageHeader(dw.PageHeaderProps{
			Title:       "Operations Dashboard",
			Description: "A scan-friendly layout for KPIs, queues, and recent activity.",
			Actions: mf.Actions(mf.ActionsProps{Align: "end", Wrap: true},
				mf.Button("Export", mf.ComponentProps{Class: "btn-outline btn-sm"}),
				mf.Button("Refresh", mf.ComponentProps{Class: "btn-primary btn-sm"}),
			),
		}),
		dw.MetricGrid(dw.MetricGridProps{Items: []dw.Metric{
			{Title: "Open tickets", Value: "128", Description: "18 need review", Trend: "12% below limit", TrendTone: dw.ToneSuccess, Icon: metricIcon("T")},
			{Title: "SLA", Value: "97.4%", Description: "Current period", Trend: "On target", TrendTone: dw.ToneSuccess, Icon: metricIcon("S")},
			{Title: "Deploys", Value: "12", Description: "Last 24 hours", Trend: "3 scheduled", TrendTone: dw.ToneNeutral, Icon: metricIcon("D")},
			{Title: "Incidents", Value: "1", Description: "Active watch", Trend: "Owner assigned", TrendTone: dw.ToneWarning, Icon: metricIcon("I")},
		}}),
		mf.Split(mf.SplitProps{
			Main: dw.CardPanel(dw.CardPanelProps{Title: "Recent activity", Description: "Most recent operational events."}, activityTable()),
			Aside: mf.Stack(mf.StackProps{Direction: "column", Gap: "4"},
				dw.CardPanel(dw.CardPanelProps{Title: "Queue health", Description: "Live backlog pressure."},
					progressRow("Ingest queue", 76, "progress-primary"),
					progressRow("Review queue", 42, "progress-warning"),
					progressRow("Escalations", 18, "progress-error"),
				),
				dw.CardPanel(dw.CardPanelProps{Title: "Next actions", Description: "Short operational checklist."},
					mf.UlProps(mf.ElementProps{Class: "space-y-3 text-sm"},
						actionItem("Review high-priority tickets", "18 waiting"),
						actionItem("Confirm deploy window", "Today 16:00"),
						actionItem("Assign follow-up owner", "1 incident"),
					),
				),
			),
			AsideWidth: "sm",
			Gap:        "6",
		}),
	)
	return appShell("/", body)
}

func activityTable() mf.Node {
	items := []activity{
		{Time: "09:20", Event: "Payment sync completed", Owner: "Finance", Status: "Done"},
		{Time: "10:05", Event: "Customer import queued", Owner: "Growth", Status: "Queued"},
		{Time: "10:40", Event: "API latency alert", Owner: "Platform", Status: "Watching"},
		{Time: "11:15", Event: "Deploy checklist opened", Owner: "Release", Status: "Active"},
	}
	return dw.DataTable(dw.DataTableProps[activity]{
		Columns: []dw.Column[activity]{
			{Header: "Time", Cell: func(a activity) mf.Node { return mf.Text(a.Time) }, Class: "w-20"},
			{Header: "Event", Cell: func(a activity) mf.Node { return mf.Text(a.Event) }},
			{Header: "Owner", Cell: func(a activity) mf.Node { return mf.Text(a.Owner) }},
			{Header: "Status", Cell: func(a activity) mf.Node { return statusBadge(a.Status) }},
		},
		Rows:    items,
		Zebra:   true,
		Compact: true,
	})
}

func progressRow(label string, value float64, className string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "space-y-2"},
		mf.DivProps(mf.ElementProps{Class: "flex items-center justify-between text-sm"},
			mf.SpanProps(mf.ElementProps{Class: "font-medium"}, mf.Text(label)),
			mf.SpanProps(mf.ElementProps{Class: "text-base-content/60"}, mf.Text(percent(value))),
		),
		mf.Progress(mf.ProgressProps{Value: value, Max: 100, Props: mf.ComponentProps{Class: className}}),
	)
}

func actionItem(label, meta string) mf.Node {
	return mf.LiProps(mf.ElementProps{Class: "flex items-center justify-between gap-3 rounded-box bg-base-200 px-3 py-2"},
		mf.SpanProps(mf.ElementProps{}, mf.Text(label)),
		mf.Badge(mf.BadgeProps{Label: meta, Props: mf.ComponentProps{Class: "badge-outline"}}),
	)
}

func statusBadge(status string) mf.Node {
	className := "badge-info"
	if status == "Done" {
		className = "badge-success"
	}
	if status == "Watching" {
		className = "badge-warning"
	}
	return mf.Badge(mf.BadgeProps{Label: status, Props: mf.ComponentProps{Class: className}})
}

func appShell(currentPath string, body mf.Node) mf.Node {
	return dw.Shell(dw.ShellProps{
		Brand:             dw.Brand{Title: "OpsDesk", Subtitle: "Dashboard template", Mark: "O"},
		CurrentPath:       currentPath,
		SearchPlaceholder: "Search operations",
		Navigation: []dw.NavGroup{{
			Label: "Workspace",
			Items: []dw.NavItem{{Path: "/", Label: "Dashboard", Icon: "D"}, {Path: "#", Label: "Activity", Icon: "A", Badge: "4"}},
		}},
		User: dw.UserMenu{Name: "Template User", Email: "ops@example.com", Initials: "TU"},
	}, body)
}

func metricIcon(label string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "grid h-10 w-10 place-items-center rounded-box bg-primary/10 font-bold text-primary"}, mf.Text(label))
}

func percent(value float64) string {
	return fmt.Sprintf("%.0f%%", value)
}
