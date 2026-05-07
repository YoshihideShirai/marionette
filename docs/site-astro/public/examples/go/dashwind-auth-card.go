package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

func RegisterDashwindAuthCardExample(app *mb.App) {
	app.Page("/dashwind-auth-card", func(ctx *mb.Context) mf.Node {
		return dw.AuthCard(dw.AuthCardProps{
			Title:       "Sign in",
			Description: "Use the preset for split-screen authentication flows.",
			SubmitLabel: "Continue",
			Fields: []dw.Field{
				{Label: "Email", Type: "email", Value: "ada@example.com", Required: true},
				{Label: "Password", Type: "password", Required: true, Help: "At least 12 characters."},
			},
			Footer: mf.Link(mf.LinkProps{Label: "Forgot password?", Href: "#"}),
		})
	})
}
