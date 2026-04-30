package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterSwitchExample(app *mb.App) {
	app.Page("/switch", func(ctx *mb.Context) mf.Node {
		return mf.SwitchComponent(mf.SwitchComponentProps{
			Name:    "notifications",
			Value:   "enabled",
			Label:   "Enable notifications",
			Checked: true,
		})
	})
}
