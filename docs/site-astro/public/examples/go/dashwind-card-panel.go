package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

func RegisterDashwindCardPanelExample(app *mb.App) {
	app.Page("/dashwind-card-panel", func(ctx *mb.Context) mf.Node {
		return dw.CardPanel(dw.CardPanelProps{
			Title:       "Revenue health",
			Description: "A reusable admin card surface with DaisyUI styling.",
			Actions:     mf.Button("Export", mf.ComponentProps{Variant: "ghost", Size: "sm"}),
		}, mf.TextComponent(mf.TextProps{Text: "$128.4k", Size: "3xl", Weight: "bold"}))
	})
}
