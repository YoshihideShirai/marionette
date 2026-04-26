package main

import (
	"log"
	"strconv"

	"github.com/example/marionette/internal/marionette"
)

func main() {
	app := marionette.New()
	app.Set("count", 0)

	app.Render(func(ctx *marionette.Context) marionette.Node {
		return renderCounter(app.GetInt("count"))
	})

	app.Handle("counter/increment", func(ctx *marionette.Context) marionette.Node {
		count := app.GetInt("count") + 1
		app.Set("count", count)
		return renderCounter(count)
	})

	if err := app.Run("127.0.0.1:8080"); err != nil {
		log.Fatal(err)
	}
}

func renderCounter(count int) marionette.Node {
	return marionette.DivClass("app", "card bg-base-100 shadow-xl max-w-md",
		marionette.DivClass("", "card-body",
			marionette.Column(
				marionette.DivClass("", "text-2xl font-bold", marionette.Text("Counter Demo")),
				marionette.DivClass("", "text-base", marionette.Text("Count: "+strconv.Itoa(count))),
				marionette.Button("Increment").OnClick("counter/increment").TargetSelector("#app"),
			),
		),
	)
}
