package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterSectionExample(app *mb.App) {
	app.Page("/section", func(ctx *mb.Context) mf.Node {
		return mf.Section(mf.SectionProps{
			Title:       "Details",
			Description: "Supporting information for the current workflow.",
		}, mf.TextComponent(mf.TextProps{Text: "Section body"}))
	})
}
