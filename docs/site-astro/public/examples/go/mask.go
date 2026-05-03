package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterMaskExample(app *mb.App) {
	app.Page("/mask", func(ctx *mb.Context) mf.Node {
		return mf.Mask("mask-squircle w-24", mf.Image(mf.ImageProps{Src: "https://placehold.co/96x96", Alt: "mask"}))
	})
}
