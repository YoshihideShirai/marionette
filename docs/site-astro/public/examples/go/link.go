package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterLinkExample(app *mb.App) {
	app.Page("/link", func(ctx *mb.Context) mf.Node {
		return mf.Actions(mf.ActionsProps{Gap: "sm", Wrap: true},
			mf.Link(mf.LinkProps{
				Label: "View dashboard",
				Href:  "/dashboard",
			}),
			mf.ExternalLink("Open docs", "https://example.com/docs", mf.ComponentProps{}),
			mf.ExternalIconLink("↗", "Open docs in a new tab", "https://example.com/docs", mf.ComponentProps{
				Variant: "ghost",
				Size:    "sm",
			}),
			mf.DownloadLink("Download CSV", ctx.Asset("reports/users.csv"), "users.csv", mf.ComponentProps{
				Variant: "primary",
				Size:    "sm",
			}),
		)
	})
}
