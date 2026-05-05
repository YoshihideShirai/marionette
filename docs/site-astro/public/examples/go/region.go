package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterRegionExample(app *mb.App) {
	app.Page("/region", func(ctx *mb.Context) mf.Node {
		return mf.Region(mf.RegionProps{ID: "summary-region"},
			mf.Card(mf.CardProps{Title: "Summary"}, mf.TextComponent(mf.TextProps{Text: "Region content can be updated independently."})),
		)
	})
}
