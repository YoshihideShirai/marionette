package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterCheckboxExample(app *mb.App) {
	app.Page("/checkbox", func(ctx *mb.Context) mf.Node {
		return mf.CheckboxComponent(mf.CheckboxComponentProps{
			Name:    "send_invite",
			Value:   "yes",
			Label:   "Send invite email",
			Checked: true,
		})
	})
}
