package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterAvatarExample(app *mb.App) {
	app.Page("/avatar", func(ctx *mb.Context) mf.Node {
		return mf.Avatar("https://placehold.co/96x96", "avatar", "w-24 rounded")
	})
}
