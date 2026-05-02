package main

import "github.com/YoshihideShirai/marionette/internal/adminsample"

func main() {
	app := adminsample.BuildApp()
	if err := app.Run("127.0.0.1:8082"); err != nil {
		panic(err)
	}
}
