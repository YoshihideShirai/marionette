package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterContainerExample(app *mb.App) {
	app.Page("/container", func(ctx *mb.Context) mf.Node {
		return mf.Container(
			mf.ContainerProps{MaxWidth: "md", Padding: "lg", Centered: true},
			mf.Card(mf.CardProps{Title: "Contained content"}, mf.UIText(mf.TextProps{Text: "Centered page section."})),
		)
	})
}
