package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterFormFieldExample(app *mb.App) {
	app.Page("/form-field", func(ctx *mb.Context) mf.Node {
		return mf.FormFieldComponent(
			mf.InputWithOptions("email", "", mf.InputOptions{Placeholder: "team@example.com"}),
			mf.FormFieldProps{Label: "Email", Required: true, Hint: "Use your work email."},
		)
	})
}
