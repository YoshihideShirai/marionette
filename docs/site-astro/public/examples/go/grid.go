package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterGridExample(app *mb.App) {
	app.Page("/grid", func(ctx *mb.Context) mf.Node {
		return mf.Grid(mf.GridProps{Columns: "3", Gap: "lg"},
			mf.Card(mf.CardProps{Title: "Users"}, mf.UIText(mf.TextProps{Text: "24"})),
			mf.Card(mf.CardProps{Title: "Teams"}, mf.UIText(mf.TextProps{Text: "6"})),
			mf.Card(mf.CardProps{Title: "Alerts"}, mf.UIText(mf.TextProps{Text: "1"})),
		)
	})
}
