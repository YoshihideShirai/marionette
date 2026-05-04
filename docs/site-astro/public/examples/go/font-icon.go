package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterFontIconExample(app *mb.App) {
	app.Page("/font-icon", func(ctx *mb.Context) mf.Node {
		return mf.Stack(
			mf.StackProps{Direction: "horizontal", Gap: "md", Align: "center"},
			mf.FontIcon(mf.FontIconProps{Library: "material-icons", Name: "check_circle", Decorative: true}),
			mf.FontIcon(mf.FontIconProps{Library: "material-icons", Name: "warning", Decorative: true}),
			mf.FontIcon(mf.FontIconProps{Library: "material-icons", Name: "info", Decorative: true}),
		)
	})
}
