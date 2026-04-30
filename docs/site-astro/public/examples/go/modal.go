package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterModalExample(app *mb.App) {
	app.Page("/modal", func(ctx *mb.Context) mf.Node {
		return mf.Modal(mf.ModalProps{
			Title: "Delete user",
			Body:  mf.UIText(mf.TextProps{Text: "Confirm deletion before continuing."}),
			Actions: mf.Actions(mf.ActionsProps{Align: "end", Gap: "sm"},
				mf.Button("Cancel", mf.ComponentProps{Variant: "ghost", Size: "sm"}),
				mf.Button("Delete", mf.ComponentProps{Variant: "danger", Size: "sm"}),
			),
			Open: true,
		})
	})
}
