package main

import (
	"fmt"
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
	return mf.ComponentContainer(mf.ContainerProps{MaxWidth: "4xl", Centered: true},
		mf.ComponentStack(mf.StackProps{Direction: "column", Gap: "6"},
			mf.ComponentPageHeader(mf.PageHeaderProps{
				Title:       "Simple Tasks",
				Description: "Marionette end-to-end sample",
			}),
			mf.ComponentForm(mf.FormProps{
				Class: "space-y-3",
				Attrs: mf.Attrs{
					"hx-post":   "/tasks/create",
					"hx-target": "#task-list",
					"hx-swap":   "innerHTML",
				},
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
				mf.ComponentSubmitButton("Add Task", mf.ComponentProps{}),
			),
			mf.DivID("task-list", taskList(tasks)),
		),
	)
}

func taskList(tasks []task) mf.Node {
	if len(tasks) == 0 {
		return mf.Element("p", mf.ElementProps{Class: "opacity-70"}, mf.Text("No tasks yet"))
	}

	rows := make([]mf.TableRowData, 0, len(tasks))
	for _, t := range tasks {
		rows = append(rows, mf.TableRow(
			mf.Text(fmt.Sprintf("%d", t.ID)),
			mf.Text(t.Name),
		))
	}
	return mf.ComponentTable(mf.TableProps{
		Columns: []mf.TableColumn{{Label: "ID"}, {Label: "Name"}},
		Rows:    toComponentRows(rows),
	})
}

func toComponentRows(rows []mf.TableRowData) []mf.TableComponentRow {
	converted := make([]mf.TableComponentRow, 0, len(rows))
	for _, row := range rows {
		converted = append(converted, mf.TableComponentRow{Cells: row.Cells})
	}
	return converted
}
