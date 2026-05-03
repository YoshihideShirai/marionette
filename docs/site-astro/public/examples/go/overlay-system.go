package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterOverlaySystemExample wires a Marionette page used in docs snippets.
func RegisterOverlaySystemExample(app *mb.App) {
	app.Page("/overlay-system", func(ctx *mb.Context) mf.Node {
		return mf.Modal(mf.ModalProps{
			Title:   "Overlay demo",
			Body:    mf.TextComponent(mf.TextProps{Text: "Modal uses the shared overlay layer in the runtime shell."}),
			Actions: mf.Button("Close", mf.ComponentProps{Variant: "ghost", Size: "sm"}),
			Open:    true,
		})
	})
}
