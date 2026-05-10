package main

import (
	"strconv"
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

func main() {
	app := buildApp("Simple Tasks", "Marionette end-to-end sample")
	if err := app.Run("127.0.0.1:8081"); err != nil {
		panic(err)
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
		dw.PageHeader(dw.PageHeaderProps{
			Title:       title,
			Description: description,
			Actions:     mf.Button("Export", mf.ComponentProps{Class: "btn-outline btn-sm"}),
		}),
		dw.MetricGrid(dw.MetricGridProps{Items: []dw.Metric{
			{Title: "Tasks", Value: strconv.Itoa(len(tasks)), Description: "Stored in Go state", Trend: "Updated by POST action", TrendTone: dw.ToneSuccess, Icon: metricIcon("T")},
			{Title: "Routes", Value: "2", Description: "Page + action", Trend: "Small surface area", TrendTone: dw.ToneNeutral, Icon: metricIcon("R")},
			{Title: "Fragments", Value: "1", Description: "#task-list region", Trend: "htmx swap target", TrendTone: dw.ToneSuccess, Icon: metricIcon("F")},
			{Title: "Stack", Value: "Go", Description: "No SPA build", Trend: "daisyUI included", TrendTone: dw.ToneNeutral, Icon: metricIcon("G")},
		}}),
		mf.Grid(mf.GridProps{Columns: "3", Gap: "lg"},
			dw.CardPanel(dw.CardPanelProps{Title: "Add task", Description: "Submit refreshes only the task table.", Class: "lg:col-span-1"},
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

	return dw.CardPanel(dw.CardPanelProps{Title: "Tasks", Description: "This panel is returned as an htmx fragment."},
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
		Brand:             dw.Brand{Title: "Simple", Subtitle: "Task sample", Mark: "S"},
		CurrentPath:       "/",
		SearchPlaceholder: "Search tasks",
		Navigation: []dw.NavGroup{{
			Label: "Sample",
			Items: []dw.NavItem{{Path: "/", Label: title, Icon: "T"}},
		}},
		User: dw.UserMenu{Name: "Demo User", Email: "demo@example.com", Initials: "DU"},
	}, body)
}

func metricIcon(label string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "grid h-10 w-10 place-items-center rounded-box bg-primary/10 font-bold text-primary"}, mf.Text(label))
}
