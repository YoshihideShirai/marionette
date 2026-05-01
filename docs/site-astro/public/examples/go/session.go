package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterSessionExample wires a small login/logout session sample for docs snippets.
func RegisterSessionExample(app *mb.App) {
	app.Page("/session", func(ctx *mb.Context) mf.Node {
		user := ctx.Session("user")
		if user == "" {
			return mf.Card(mf.CardProps{
				Title:       "Session sample",
				Description: "You are not signed in.",
			},
				mf.Form("/session/login",
					mf.SubmitButton("Sign in as Aiko", mf.ComponentProps{Variant: "primary"}),
				),
			)
		}

		return mf.Card(mf.CardProps{
			Title:       "Session sample",
			Description: "Signed in as " + user,
		},
			mf.Form("/session/logout",
				mf.SubmitButton("Sign out", mf.ComponentProps{Variant: "ghost"}),
			),
		)
	})

	app.Action("session/login", func(ctx *mb.Context) mf.Node {
		ctx.SetSession("user", "Aiko")
		return mf.Toast(mf.ToastProps{Title: "Signed in", Description: "Signed in as Aiko", Props: mf.ComponentProps{Variant: "success"}})
	})

	app.Action("session/logout", func(ctx *mb.Context) mf.Node {
		ctx.ClearSession()
		return mf.Toast(mf.ToastProps{Title: "Signed out", Description: "Session cleared", Props: mf.ComponentProps{Variant: "info"}})
	})
}
