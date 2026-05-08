package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterMarkdownExample(app *mb.App) {
	app.Page("/markdown", func(ctx *mb.Context) mf.Node {
		return mf.Markdown(mf.MarkdownProps{Content: "# Overview\n\n- Fast rendering\n- Safe HTML output"})
	})
}
