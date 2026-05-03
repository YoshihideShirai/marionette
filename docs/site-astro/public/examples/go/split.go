package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterSplitExample(app *mb.App) {
	app.Page("/split", func(ctx *mb.Context) mf.Node {
		return mf.Split(mf.SplitProps{
			Main:  mf.Card(mf.CardProps{Title: "Main"}, mf.TextComponent(mf.TextProps{Text: "Primary workflow"})),
			Aside: mf.Card(mf.CardProps{Title: "Aside"}, mf.TextComponent(mf.TextProps{Text: "Supporting details"})),
			Gap:   "lg",
		})
	})
}
