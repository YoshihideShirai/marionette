package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterLayoutExample(app *mb.App) {
	app.Page("/layout", func(ctx *mb.Context) mf.Node {
		return mf.ContainerComponent(mf.ContainerProps{MaxWidth: "md", Centered: true},
			mf.Grid(mf.GridProps{Columns: "2", Gap: "lg"},
				mf.Card(mf.CardProps{Title: "Main"}, mf.UIText(mf.TextProps{Text: "Primary content"})),
				mf.Card(mf.CardProps{Title: "Aside"}, mf.UIText(mf.TextProps{Text: "Supporting content"})),
			),
		)
	})
}
