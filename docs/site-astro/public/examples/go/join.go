package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterJoinExample(app *mb.App) {
	app.Page("/join", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "join example"})
	})
}
