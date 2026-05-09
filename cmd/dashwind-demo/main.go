package main

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	"github.com/YoshihideShirai/marionette/internal/dashwinddemo"
)

func main() {
	app := buildApp()
	if err := app.Run("127.0.0.1:8083"); err != nil {
		panic(err)
	}
}

func buildApp() *mb.App {
	return dashwinddemo.BuildApp()
}
