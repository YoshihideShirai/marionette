package main

import (
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

type task struct {
	ID   int
	Name string
}

func main() {
	app := mb.New()
	app.Set("tasks", []task{})
	app.Set("nextID", 1)

	app.Page("/", func(ctx *mb.Context) mf.Node {
		return page(ctx)
	})

	app.Action("tasks/create", func(ctx *mb.Context) mf.Node {
		name := strings.TrimSpace(ctx.FormValue("name"))
		if name != "" {
			tasks := ctx.Get("tasks").([]task)
			nextID := ctx.Get("nextID").(int)
			tasks = append(tasks, task{ID: nextID, Name: name})
			ctx.Set("tasks", tasks)
			ctx.Set("nextID", nextID+1)
		}
		return taskList(ctx.Get("tasks").([]task))
	})

	if err := app.Run("127.0.0.1:8081"); err != nil {
		panic(err)
	}
}

func page(ctx *mb.Context) mf.Node {
	tasks := ctx.Get("tasks").([]task)
	return mf.ContainerComponent(mf.ContainerProps{MaxWidth: "4xl", Centered: true},
		mf.Stack(mf.StackProps{Direction: "column", Gap: "6"},
			mf.PageHeader(mf.PageHeaderProps{
				Title:       "Simple Tasks",
				Description: "Marionette end-to-end sample",
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
	return mf.TableComponent(mf.TableProps{
		Columns: []mf.TableColumn{{Label: "ID"}, {Label: "Name"}},
		Rows:    rows,
	})
}
