package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterImageExample(app *mb.App) {
	app.Page("/image", func(ctx *mb.Context) mf.Node {
		return mf.Image(mf.ImageProps{
			Src:         "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=80",
			Alt:         "Desk with laptop and notebook",
			Caption:     "Workspace preview",
			Width:       1200,
			Height:      800,
			AspectRatio: "video",
			ObjectFit:   "cover",
		})
	})
}
