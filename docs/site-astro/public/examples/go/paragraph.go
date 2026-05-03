package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterParagraphExample(app *mb.App) {
	app.Page("/paragraph", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "Paragraph helper for readable long-form text."})
	})
}
