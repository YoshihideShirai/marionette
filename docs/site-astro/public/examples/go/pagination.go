package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterPaginationExample(app *mb.App) {
	app.Page("/pagination", func(ctx *mb.Context) mf.Node {
		return mf.Pagination(mf.PaginationProps{
			Page:       2,
			TotalPages: 4,
			PrevHref:   "/users?page=1",
			NextHref:   "/users?page=3",
		})
	})
}
