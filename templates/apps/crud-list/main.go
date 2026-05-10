package main

import (
	"strconv"
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

type task struct {
	ID     int
	Name   string
	Owner  string
	Status string
}

func main() {
	app := buildApp()
	if err := app.Run("127.0.0.1:8091"); err != nil {
		panic(err)
	}
}

func buildApp() *mb.App {
	app := mb.New()
	dw.Use(app, dw.Options{Theme: "corporate"})
	app.SetGlobal("tasks", []task{
		{ID: 1, Name: "Map first workflow", Owner: "Product", Status: "Active"},
		{ID: 2, Name: "Invite first user", Owner: "Success", Status: "Queued"},
		{ID: 3, Name: "Review launch checklist", Owner: "Ops", Status: "Done"},
	})
	app.SetGlobal("nextID", 4)

	app.Page("/", page, mb.WithTitle("CRUD List Template"))
	app.Action("tasks/create", createTask)
	app.Action("tasks/complete", completeTask)

	return app
}

func page(ctx *mb.Context) mf.Node {
	tasks := ctx.GetGlobal("tasks").([]task)
	body := mf.DivProps(mf.ElementProps{Class: "space-y-6"},
		dw.PageHeader(dw.PageHeaderProps{
			Title:       "Task Operations",
			Description: "A CRUD-ready layout with form actions, table fragments, and status metrics.",
			Actions:     mf.Button("Import CSV", mf.ComponentProps{Class: "btn-outline btn-sm"}),
		}),
		metrics(tasks),
		mf.Grid(mf.GridProps{Columns: "3", Gap: "lg"},
			dw.CardPanel(dw.CardPanelProps{Title: "Create task", Description: "The POST action refreshes only the table region.", Class: "lg:col-span-1"},
				createForm(),
			),
			mf.DivProps(mf.ElementProps{Class: "lg:col-span-2"},
				mf.Region(mf.RegionProps{ID: "task-list"}, taskList(tasks)),
			),
		),
	)
	return appShell("/", body)
}

func createForm() mf.Node {
	return mf.ActionForm(mf.ActionFormProps{
		Action: "/tasks/create",
		Target: "#task-list",
		Swap:   "innerHTML",
		Props:  mf.ComponentProps{Class: "space-y-4"},
	},
		mf.FormRow(mf.FormRowProps{
			ID:       "task-name",
			Label:    "Task name",
			Required: true,
			Control: mf.TextField(mf.TextFieldProps{
				ID:          "task-name",
				Name:        "name",
				Placeholder: "Review invoices",
				Required:    true,
			}),
		}),
		mf.FormRow(mf.FormRowProps{
			ID:    "task-owner",
			Label: "Owner",
			Control: mf.Select(mf.SelectFieldProps{
				ID:   "task-owner",
				Name: "owner",
				Options: []mf.SelectOption{
					{Label: "Ops", Value: "Ops", Selected: true},
					{Label: "Product", Value: "Product"},
					{Label: "Success", Value: "Success"},
				},
			}),
		}),
		mf.SubmitButton("Add task", mf.ComponentProps{Class: "btn-primary w-full"}),
	)
}

func createTask(ctx *mb.Context) mf.Node {
	name := strings.TrimSpace(ctx.FormValue("name"))
	owner := strings.TrimSpace(ctx.FormValue("owner"))
	if owner == "" {
		owner = "Ops"
	}
	if name != "" {
		nextID := ctx.IncrementGlobalInt("nextID", 1) - 1
		ctx.UpdateGlobal("tasks", func(old any) any {
			tasks := append([]task(nil), old.([]task)...)
			return append(tasks, task{ID: nextID, Name: name, Owner: owner, Status: "Queued"})
		})
	}
	return taskList(ctx.GetGlobal("tasks").([]task))
}

func completeTask(ctx *mb.Context) mf.Node {
	id := strings.TrimSpace(ctx.FormValue("id"))
	ctx.UpdateGlobal("tasks", func(old any) any {
		tasks := append([]task(nil), old.([]task)...)
		for i := range tasks {
			if id == intString(tasks[i].ID) {
				tasks[i].Status = "Done"
			}
		}
		return tasks
	})
	return taskList(ctx.GetGlobal("tasks").([]task))
}

func metrics(tasks []task) mf.Node {
	done := 0
	for _, t := range tasks {
		if t.Status == "Done" {
			done++
		}
	}
	return dw.MetricGrid(dw.MetricGridProps{Items: []dw.Metric{
		{Title: "Total tasks", Value: intString(len(tasks)), Description: "Current list", Trend: "Stored in app state", TrendTone: dw.ToneNeutral, Icon: metricIcon("T")},
		{Title: "Completed", Value: intString(done), Description: "Marked done", Trend: "POST fragment updates", TrendTone: dw.ToneSuccess, Icon: metricIcon("C")},
		{Title: "Open", Value: intString(len(tasks) - done), Description: "Needs action", Trend: "Queue visible", TrendTone: dw.ToneWarning, Icon: metricIcon("O")},
		{Title: "Owners", Value: "3", Description: "Demo teams", Trend: "Select field example", TrendTone: dw.ToneNeutral, Icon: metricIcon("U")},
	}})
}

func taskList(tasks []task) mf.Node {
	if len(tasks) == 0 {
		return dw.CardPanel(dw.CardPanelProps{Title: "Tasks"},
			mf.EmptyState(mf.EmptyStateProps{Title: "No tasks yet", Description: "Create the first task to populate this region."}),
		)
	}
	return dw.CardPanel(dw.CardPanelProps{Title: "Tasks", Description: "Rows can trigger POST actions without leaving the page."},
		dw.DataTable(dw.DataTableProps[task]{
			Columns: []dw.Column[task]{
				{Header: "ID", Cell: func(t task) mf.Node { return mf.Text(intString(t.ID)) }, Class: "w-16"},
				{Header: "Name", Cell: func(t task) mf.Node { return mf.Text(t.Name) }},
				{Header: "Owner", Cell: func(t task) mf.Node { return mf.Text(t.Owner) }},
				{Header: "Status", Cell: func(t task) mf.Node { return statusBadge(t.Status) }},
				{Header: "", Cell: completeButton, Class: "text-right"},
			},
			Rows:    tasks,
			Zebra:   true,
			Compact: true,
		}),
	)
}

func completeButton(t task) mf.Node {
	if t.Status == "Done" {
		return mf.Text("")
	}
	return mf.ActionForm(mf.ActionFormProps{Action: "/tasks/complete", Target: "#task-list", Swap: "innerHTML"},
		mf.HiddenField("id", intString(t.ID)),
		mf.SubmitButton("Complete", mf.ComponentProps{Class: "btn-sm btn-outline"}),
	)
}

func statusBadge(status string) mf.Node {
	className := "badge-info"
	if status == "Done" {
		className = "badge-success"
	}
	if status == "Active" {
		className = "badge-warning"
	}
	return mf.Badge(mf.BadgeProps{Label: status, Props: mf.ComponentProps{Class: className}})
}

func appShell(currentPath string, body mf.Node) mf.Node {
	return dw.Shell(dw.ShellProps{
		Brand:             dw.Brand{Title: "TaskOps", Subtitle: "CRUD template", Mark: "T"},
		CurrentPath:       currentPath,
		SearchPlaceholder: "Search tasks",
		Navigation: []dw.NavGroup{{
			Label: "Workspace",
			Items: []dw.NavItem{{Path: "/", Label: "Tasks", Icon: "T"}, {Path: "#", Label: "Archive", Icon: "A", Badge: "Soon"}},
		}},
		User: dw.UserMenu{Name: "Template User", Email: "tasks@example.com", Initials: "TU"},
	}, body)
}

func metricIcon(label string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "grid h-10 w-10 place-items-center rounded-box bg-primary/10 font-bold text-primary"}, mf.Text(label))
}

func intString(v int) string {
	return strconv.Itoa(v)
}
