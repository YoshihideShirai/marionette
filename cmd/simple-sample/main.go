package main

import (
	"github.com/YoshihideShirai/marionette/internal/tasksdemo"
)

func main() {
	app := tasksdemo.BuildApp("Simple Tasks", "Marionette end-to-end sample")
	if err := app.Run("127.0.0.1:8081"); err != nil {
		panic(err)
	}
}
