package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterAppShellExample(app *mb.App) {
	app.Page("/app-shell", func(ctx *mb.Context) mf.Node {
		return mf.AppShell(mf.AppShellProps{
			Header:  mf.H2(mf.Text("Workspace")),
			Content: mf.TextComponent(mf.TextProps{Text: "Main content rendered inside AppShell.", Tone: "muted"}),
		})
	})
}
