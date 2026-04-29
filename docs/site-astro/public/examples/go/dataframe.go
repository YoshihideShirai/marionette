//go:build ignore
// +build ignore

package main

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func main() {
	app := mb.New()

	app.Page("/", func(ctx *mb.Context) mf.Node {
		return mf.Div(
			mf.H1(mf.Text("UIDataFrame example")),
			mf.P(mf.Text("Render this component with Marionette frontend nodes.")),
		)
	})

	app.Run(":8080")
}
