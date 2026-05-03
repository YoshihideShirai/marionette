package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterDrawerExample(app *mb.App) {
	app.Page("/drawer", func(ctx *mb.Context) mf.Node {
		return mf.Drawer("d1", mf.Menu(mf.Anchor("#", mf.TextComponent(mf.TextProps{Text: "Sidebar"}))), mf.Button("open", mf.ComponentProps{}))
	})
}
