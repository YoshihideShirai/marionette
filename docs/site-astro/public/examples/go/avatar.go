package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterAvatarExample(app *mb.App) {
	app.Page("/avatar", func(ctx *mb.Context) mf.Node {
		return mf.TextComponent(mf.TextProps{Text: "avatar example"})
	})
}
