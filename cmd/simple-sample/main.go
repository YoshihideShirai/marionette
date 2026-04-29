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
	return mf.DivClass("space-y-4",
		mf.Element("h1", mf.ElementProps{Class: "text-2xl font-bold"}, mf.Text("Simple Tasks")),
		mf.Element("p", mf.ElementProps{Class: "text-sm opacity-80"}, mf.Text("Marionette end-to-end sample")),
		mf.Element("form", mf.ElementProps{Attrs: mf.Attrs{
			"hx-post":   "/tasks/create",
			"hx-target": "#task-list",
			"hx-swap":   "innerHTML",
			"class":     "flex gap-2",
		}},
			mf.Element("input", mf.ElementProps{Attrs: mf.Attrs{"name": "name", "placeholder": "Task name", "class": "input input-bordered w-full"}}),
			mf.Element("button", mf.ElementProps{Attrs: mf.Attrs{"type": "submit", "class": "btn btn-primary"}}, mf.Text("Add Task")),
		),
		mf.DivID("task-list", taskList(tasks)),
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
	return mf.Table([]string{"ID", "Name"}, rows...)
}
