package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// RegisterFormExample wires a Marionette page used in docs snippets.
func RegisterFormExample(app *mb.App) {
	app.Page("/form", func(ctx *mb.Context) mf.Node {
		return mf.FormComponent(mf.FormProps{Method: "post", Action: "/users"},
			mf.FormFieldComponent(
				mf.InputComponent("name", "", mf.ComponentProps{}),
				mf.FormFieldProps{Label: "Name", Required: true, Hint: "Enter a display name."},
			),
			mf.SubmitButton("Create user", mf.ComponentProps{Variant: "primary"}),
		)
	})
}
