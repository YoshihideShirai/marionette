package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterCardExample(app *mb.App) {
	app.Page("/card", func(ctx *mb.Context) mf.Node {
		return mf.Card(mf.CardProps{
			Title:       "Workspace summary",
			Description: "Header, description, actions, then body content.",
			Actions:     mf.Button("Edit", mf.ComponentProps{Variant: "ghost", Size: "sm"}),
		}, mf.UIText(mf.TextProps{Text: "Ready"}))
	})
}
