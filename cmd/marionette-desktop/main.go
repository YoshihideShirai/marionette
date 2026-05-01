package main

import (
	"log"

	"github.com/YoshihideShirai/marionette/desktop"
	"github.com/YoshihideShirai/marionette/internal/tasksdemo"
)

func main() {
	app := tasksdemo.BuildApp(
		"Marionette Desktop",
		"The same server-side Marionette app running inside a desktop WebView shell.",
	)
	if err := desktop.Run(app, desktop.Options{
		Title:  "Marionette Desktop",
		Width:  1200,
		Height: 800,
	}); err != nil {
		log.Fatal(err)
	}
}
