package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterFileUploadExample(app *mb.App) {
	app.Page("/file-upload", func(ctx *mb.Context) mf.Node {
		return mf.FileUpload("attachment", true, mf.ComponentProps{Variant: "bordered"})
	})
}
