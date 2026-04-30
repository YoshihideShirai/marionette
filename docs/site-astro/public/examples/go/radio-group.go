package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterRadioGroupExample(app *mb.App) {
	app.Page("/radio-group", func(ctx *mb.Context) mf.Node {
		return mf.RadioGroupComponent(mf.RadioGroupComponentProps{
			Name:      "role",
			AriaLabel: "Role",
			Items: []mf.RadioItem{
				{Label: "Admin", Value: "admin"},
				{Label: "Viewer", Value: "viewer", Checked: true},
			},
		})
	})
}
