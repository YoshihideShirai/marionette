package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterTableWithVariantsExample(app *mb.App) {
	app.Page("/table-with-variants", func(ctx *mb.Context) mf.Node {
		return mf.TableWithVariants(
			[]string{"Name", "Role"},
			[][]mf.Node{{mf.Text("Aiko"), mf.Text("Admin")}, {mf.Text("Ren"), mf.Text("Editor")}},
			true, false, false, "sm",
		)
	})
}
