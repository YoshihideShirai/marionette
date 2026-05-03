package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterChatBubbleExample(app *mb.App) {
	app.Page("/chat-bubble", func(ctx *mb.Context) mf.Node {
		return mf.Div(mf.Class("chat chat-start"), mf.Div(mf.Class("chat-bubble"), mf.TextComponent(mf.TextProps{Text: "Hello"})))
	})
}
