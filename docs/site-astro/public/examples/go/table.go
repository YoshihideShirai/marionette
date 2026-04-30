package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterTableExample(app *mb.App) {
	app.Page("/table", func(ctx *mb.Context) mf.Node {
		return mf.Table(mf.TableProps{
			Columns: []mf.TableColumn{{Label: "Name"}, {Label: "Role"}},
			Rows: []mf.TableComponentRow{
				mf.TableRowValues("Aiko Tanaka", mf.Badge(mf.BadgeProps{Label: "Admin"})),
				mf.TableRowValues("Ren Sato", mf.Badge(mf.BadgeProps{Label: "Editor"})),
			},
		})
	})
}
