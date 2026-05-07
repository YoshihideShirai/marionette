package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

func RegisterDashwindSettingsSectionExample(app *mb.App) {
	app.Page("/dashwind-settings-section", func(ctx *mb.Context) mf.Node {
		return dw.SettingsSection(dw.SettingsSectionProps{
			Title:       "Workspace settings",
			Description: "Card-backed form section with labels, help text, and submit actions.",
			SubmitLabel: "Save changes",
			Fields: []dw.Field{
				{Label: "Workspace name", Value: "Growth Ops", Required: true},
				{Label: "Notification email", Type: "email", Value: "ops@example.com", Help: "Weekly summary recipient."},
			},
		})
	})
}
