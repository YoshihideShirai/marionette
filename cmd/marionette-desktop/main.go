package main

import (
	"log"
	"strconv"
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	"github.com/YoshihideShirai/marionette/desktop"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

func main() {
	app := buildApp(desktopTitle, desktopDescription)
	if err := desktop.Run(app, desktopOptions()); err != nil {
		log.Fatal(err)
	}
}

const (
	desktopTitle       = "Marionette Desktop"
	desktopDescription = "The same server-side Marionette app running inside a desktop WebView shell."
)

func desktopOptions() desktop.Options {
	return desktop.Options{
		Title:  desktopTitle,
		Width:  1200,
		Height: 800,
	}
}

type task struct {
	ID   int
	Name string
}

func buildApp(title, description string) *mb.App {
	app := mb.New()
	dw.Use(app, dw.Options{Theme: "corporate"})
	app.SetGlobal("tasks", []task{})
	app.SetGlobal("nextID", 1)

	app.Page("/", func(ctx *mb.Context) mf.Node {
		return page(ctx, title, description)
	}, mb.WithTitle(title))

	app.Action("tasks/create", func(ctx *mb.Context) mf.Node {
		name := strings.TrimSpace(ctx.FormValue("name"))
		if name != "" {
			nextID := ctx.IncrementGlobalInt("nextID", 1) - 1
			ctx.UpdateGlobal("tasks", func(old any) any {
				tasks := append([]task(nil), old.([]task)...)
				return append(tasks, task{ID: nextID, Name: name})
			})
		}
		return taskList(ctx.GetGlobal("tasks").([]task))
	})

	return app
}

func page(ctx *mb.Context, title, description string) mf.Node {
	tasks := ctx.GetGlobal("tasks").([]task)
	body := mf.DivProps(mf.ElementProps{Class: "space-y-6"},
		dw.PageHeader(dw.PageHeaderProps{Title: title, Description: description}),
		dw.MetricGrid(dw.MetricGridProps{Items: []dw.Metric{
			{Title: "Tasks", Value: strconv.Itoa(len(tasks)), Description: "Local desktop session", Trend: "Same Go app", TrendTone: dw.ToneSuccess, Icon: metricIcon("T")},
			{Title: "Shell", Value: "WebView", Description: "Desktop wrapper", Trend: "1200 x 800", TrendTone: dw.ToneNeutral, Icon: metricIcon("W")},
			{Title: "Updates", Value: "htmx", Description: "Fragment swap", Trend: "#task-list", TrendTone: dw.ToneSuccess, Icon: metricIcon("H")},
			{Title: "Runtime", Value: "Go", Description: "Server-side UI", Trend: "No extra build", TrendTone: dw.ToneNeutral, Icon: metricIcon("G")},
		}}),
		mf.Grid(mf.GridProps{Columns: "3", Gap: "lg"},
			dw.CardPanel(dw.CardPanelProps{Title: "Add task", Description: "Works the same in browser and desktop shells.", Class: "lg:col-span-1"},
				mf.ActionForm(mf.ActionFormProps{
					Action: "/tasks/create",
					Target: "#task-list",
					Swap:   "innerHTML",
					Props:  mf.ComponentProps{Class: "space-y-4"},
				},
					mf.FormRow(mf.FormRowProps{
						ID:       "task-name",
						Label:    "Task",
						Required: true,
						Control: mf.TextField(mf.TextFieldProps{
							ID:          "task-name",
							Name:        "name",
							Placeholder: "Task name",
							Required:    true,
						}),
					}),
					mf.SubmitButton("Add Task", mf.ComponentProps{Class: "btn-primary w-full"}),
				),
			),
			mf.DivProps(mf.ElementProps{Class: "lg:col-span-2"},
				mf.Region(mf.RegionProps{ID: "task-list"}, taskList(tasks)),
			),
		),
	)
	return appShell(title, body)
}

func taskList(tasks []task) mf.Node {
	if len(tasks) == 0 {
		return dw.CardPanel(dw.CardPanelProps{Title: "Tasks"},
			mf.EmptyState(mf.EmptyStateProps{Title: "No tasks yet", Description: "Create a task to populate this table."}),
		)
	}

	return dw.CardPanel(dw.CardPanelProps{Title: "Tasks", Description: "The action returns this panel as a fragment."},
		dw.DataTable(dw.DataTableProps[task]{
			Columns: []dw.Column[task]{
				{Header: "ID", Cell: func(t task) mf.Node { return mf.Text(strconv.Itoa(t.ID)) }, Class: "w-20"},
				{Header: "Name", Cell: func(t task) mf.Node { return mf.Text(t.Name) }},
			},
			Rows:    tasks,
			Zebra:   true,
			Compact: true,
		}),
	)
}

func appShell(title string, body mf.Node) mf.Node {
	return dw.Shell(dw.ShellProps{
		Brand:             dw.Brand{Title: "Desktop", Subtitle: "WebView sample", Mark: "D"},
		CurrentPath:       "/",
		SearchPlaceholder: "Search tasks",
		Navigation: []dw.NavGroup{{
			Label: "Sample",
			Items: []dw.NavItem{{Path: "/", Label: title, Icon: "D"}},
		}},
		User: dw.UserMenu{Name: "Desktop User", Email: "desktop@example.com", Initials: "DU"},
	}, body)
}

func metricIcon(label string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "grid h-10 w-10 place-items-center rounded-box bg-primary/10 font-bold text-primary"}, mf.Text(label))
}
