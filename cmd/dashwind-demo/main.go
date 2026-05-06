package main

import "github.com/YoshihideShirai/marionette/internal/dashwinddemo"

func main() {
	app := dashwinddemo.BuildApp()
	if err := app.Run("127.0.0.1:8083"); err != nil {
		panic(err)
	}
}
