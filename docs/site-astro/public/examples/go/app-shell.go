package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterAppShellExample(app *mb.App) {
	app.Page("/app-shell", func(ctx *mb.Context) mf.Node {
		return mf.AppShell(
			mf.AppShellProps{Title: "Workspace"},
			mf.P("Main content rendered inside AppShell."),
		)
	})
}
