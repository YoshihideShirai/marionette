package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterChatBubbleExample(app *mb.App) {
	app.Page("/chat-bubble", func(ctx *mb.Context) mf.Node {
		return mf.ChatBubble(mf.TextComponent(mf.TextProps{Text: "Hello"}), false)
	})
}
