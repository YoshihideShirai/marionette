package main

import (
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
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
	return mf.Container(mf.ContainerProps{MaxWidth: "4xl", Centered: true},
		mf.Stack(mf.StackProps{Direction: "column", Gap: "6"},
			mf.PageHeader(mf.PageHeaderProps{
				Title:       title,
				Description: description,
			}),
			mf.ActionForm(mf.ActionFormProps{
				Action: "/tasks/create",
				Target: "#task-list",
				Swap:   "innerHTML",
				Props:  mf.ComponentProps{Class: "space-y-3"},
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
				mf.SubmitButton("Add Task", mf.ComponentProps{}),
			),
			mf.Region(mf.RegionProps{ID: "task-list"}, taskList(tasks)),
		),
	)
}

func taskList(tasks []task) mf.Node {
	if len(tasks) == 0 {
		return mf.EmptyState(mf.EmptyStateProps{Title: "No tasks yet"})
	}

	rows := make([]mf.TableComponentRow, 0, len(tasks))
	for _, t := range tasks {
		rows = append(rows, mf.TableRowValues(t.ID, t.Name))
	}
	return mf.Table(mf.TableProps{
		Columns: []mf.TableColumn{{Label: "ID"}, {Label: "Name"}},
		Rows:    rows,
	})
}
