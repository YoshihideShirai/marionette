package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterCollapseExample(app *mb.App) {
	app.Page("/collapse", func(ctx *mb.Context) mf.Node {
		return mf.Collapse("Title", mf.Paragraph(mf.ParagraphProps{Text: "Content"}), false)
	})
}
