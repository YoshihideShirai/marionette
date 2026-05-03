package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterHeroExample(app *mb.App) {
	app.Page("/hero", func(ctx *mb.Context) mf.Node {
		return mf.Hero(mf.HeroProps{},
			mf.Div(mf.Class("hero-content text-center"),
				mf.Div(mf.Class("max-w-md"),
					mf.Heading(mf.HeadingProps{Level: 1, Text: "Hero alias"}),
					mf.Paragraph(mf.ParagraphProps{Text: "Use frontend.Hero directly."}),
				),
			),
		)
	})
}
