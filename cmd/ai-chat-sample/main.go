package main

import (
	"fmt"
	stdhtml "html"
	"strings"
	"time"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

const thinkingDelay = 650 * time.Millisecond

type chatMessage struct {
	ID        int
	Role      string
	Name      string
	Content   string
	Streaming bool
	Thinking  bool
}

func main() {
	app := buildApp()
	if err := app.Run("127.0.0.1:8084"); err != nil {
		panic(err)
	}
}

func buildApp() *mb.App {
	app := mb.New()
	app.SetGlobal("messages", []chatMessage{welcomeMessage()})
	app.SetGlobal("nextMessageID", 2)
	app.SetGlobal("chatError", "")
	app.DisableCharts()
	app.EnableSSE()
	app.AddStyle(`
		#marionette-root { min-height: 100vh; width: 100%; }
		.ai-chat-message { scroll-margin-block: 1rem; }
		.ai-chat-token { white-space: pre-wrap; }
		.ai-chat-composer textarea:focus { outline: none; box-shadow: none; }
	`)
	app.Page("/", func(ctx *mb.Context) mf.Node {
		return page(ctx)
	}, mb.WithTitle("AI Chat"))

	app.Action("chat/send", func(ctx *mb.Context) mf.Node {
		prompt := strings.TrimSpace(ctx.FormValue("prompt"))
		if prompt == "" {
			ctx.SetGlobal("chatError", "Please enter a message.")
			return chatPanel(ctx)
		}

		ctx.SetGlobal("chatError", "")
		userID := ctx.IncrementGlobalInt("nextMessageID", 1) - 1
		assistantID := ctx.IncrementGlobalInt("nextMessageID", 1) - 1
		ctx.StartTextStream(mb.TextStreamOptions{
			Name:      "chat-reply",
			Text:      demoReply(prompt),
			ChunkSize: 1,
		})
		ctx.UpdateGlobal("messages", func(old any) any {
			messages := cloneMessages(old).([]chatMessage)
			messages = append(messages,
				chatMessage{ID: userID, Role: "user", Name: "You", Content: prompt},
				chatMessage{ID: assistantID, Role: "assistant", Name: "Marionette AI", Streaming: true, Thinking: true},
			)
			return messages
		})
		return chatPanel(ctx)
	})

	app.StreamAction("chat/stream", func(ctx *mb.Context) mb.Stream {
		return func(yield func(mf.Node) bool) {
			for {
				messageID, step, wasThinking := advanceStream(ctx)
				if !step.Active || messageID == 0 {
					return
				}
				if wasThinking {
					time.Sleep(thinkingDelay)
				}
				if !yield(streamDelta(messageID, step, wasThinking)) || step.Done {
					return
				}
			}
		}
	})

	app.Action("chat/reset", func(ctx *mb.Context) mf.Node {
		ctx.SetGlobal("messages", []chatMessage{welcomeMessage()})
		ctx.SetGlobal("nextMessageID", 2)
		ctx.SetGlobal("chatError", "")
		ctx.ResetTextStream("chat-reply")
		return chatPanel(ctx)
	})

	return app
}

func page(ctx *mb.Context) mf.Node {
	return chatPanel(ctx)
}

func chatPanel(ctx *mb.Context) mf.Node {
	messages := ctx.GetGlobalSnapshot("messages", cloneMessages).([]chatMessage)
	errorMessage, _ := ctx.GetGlobal("chatError").(string)

	children := []mf.Node{
		topBar(),
		conversation(messages),
	}
	if hasStreamingMessage(messages) {
		children = append(children, streamConnector())
	}
	if strings.TrimSpace(errorMessage) != "" {
		children = append(children, mf.Alert(mf.AlertProps{
			Title:       "Input error",
			Description: errorMessage,
			Props:       mf.ComponentProps{Class: "alert-warning"},
		}))
	}
	children = append(children, promptForm())

	return mf.Region(mf.RegionProps{ID: "chat-panel", Props: mf.ComponentProps{Class: "flex min-h-screen flex-col bg-base-100"}}, children...)
}

func topBar() mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "sticky top-0 z-10 border-b border-base-200 bg-base-100/95 px-4 py-3 backdrop-blur"},
		mf.DivProps(mf.ElementProps{Class: "mx-auto flex w-full max-w-5xl items-center justify-between gap-3"},
			mf.DivProps(mf.ElementProps{Class: "min-w-0"},
				mf.H1Props(mf.ElementProps{Class: "truncate text-base font-semibold"}, mf.Text("Marionette AI")),
				mf.PProps(mf.ElementProps{Class: "text-xs text-base-content/60"}, mf.Text("AI Chat")),
			),
			resetForm(),
		),
	)
}

func conversation(messages []chatMessage) mf.Node {
	items := make([]mf.Node, 0, len(messages))
	for _, msg := range messages {
		items = append(items, messageBubble(msg))
	}
	return mf.DivProps(mf.ElementProps{Class: "mx-auto flex w-full max-w-3xl flex-1 flex-col gap-6 px-4 py-8"}, items...)
}

func messageBubble(msg chatMessage) mf.Node {
	contentID := fmt.Sprintf("message-content-%d", msg.ID)
	cursorID := fmt.Sprintf("message-cursor-%d", msg.ID)
	statusID := fmt.Sprintf("message-status-%d", msg.ID)
	statusText := ""
	if msg.Streaming {
		statusText = "SSE streaming"
		if msg.Thinking {
			statusText = "Thinking"
		}
	}

	bubbleChildren := []mf.Node{
		mf.SpanProps(mf.ElementProps{ID: contentID, Class: "ai-chat-token"}, mf.Text(msg.Content)),
	}
	if msg.Thinking {
		thinkingID := fmt.Sprintf("message-thinking-%d", msg.ID)
		bubbleChildren = append(bubbleChildren,
			mf.SpanProps(mf.ElementProps{ID: thinkingID, Class: "inline-flex items-center gap-2 opacity-70"},
				mf.Text("Thinking"),
				mf.SpanProps(mf.ElementProps{Class: "loading loading-dots loading-xs"}, mf.Text("")),
			),
		)
	}
	if msg.Streaming {
		cursorClass := "opacity-70"
		cursorText := "▌"
		if msg.Thinking {
			cursorClass = "hidden"
			cursorText = ""
		}
		bubbleChildren = append(bubbleChildren, mf.SpanProps(mf.ElementProps{ID: cursorID, Class: cursorClass}, mf.Text(cursorText)))
	}

	statusNode := mf.SpanProps(mf.ElementProps{ID: statusID, Class: "sr-only"}, mf.Text(statusText))
	if msg.Role == "user" {
		return mf.DivProps(mf.ElementProps{ID: fmt.Sprintf("message-%d", msg.ID), Class: "ai-chat-message flex justify-end"},
			mf.DivProps(mf.ElementProps{Class: "max-w-[80%] rounded-3xl bg-base-200 px-4 py-2.5 text-sm leading-relaxed text-base-content sm:max-w-[70%]"},
				bubbleChildren...,
			),
			statusNode,
		)
	}

	return mf.DivProps(mf.ElementProps{ID: fmt.Sprintf("message-%d", msg.ID), Class: "ai-chat-message flex gap-4"},
		mf.DivProps(mf.ElementProps{Class: "mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-neutral text-xs font-semibold text-neutral-content"}, mf.Text("AI")),
		mf.DivProps(mf.ElementProps{Class: "min-w-0 flex-1 space-y-1"},
			mf.DivProps(mf.ElementProps{Class: "text-sm font-medium text-base-content"}, mf.Text(msg.Name)),
			mf.DivProps(mf.ElementProps{Class: "text-sm leading-7 text-base-content"}, bubbleChildren...),
			statusNode,
		),
	)
}

func streamConnector() mf.Node {
	return mf.DivProps(mf.ElementProps{
		ID:    "chat-stream-connector",
		Class: "hidden",
		Attrs: mf.Attrs{
			"aria-hidden":               "true",
			"data-marionette-sse-url":   "/chat/stream",
			"data-marionette-sse-scope": "#chat-panel",
		},
	})
}

func promptForm() mf.Node {
	return mf.ActionForm(mf.ActionFormProps{
		Action: "/chat/send",
		Target: "#chat-panel",
		Swap:   "outerHTML",
		Props:  mf.ComponentProps{Class: "ai-chat-composer sticky bottom-0 mx-auto w-full max-w-3xl bg-base-100 px-4 pb-5 pt-2"},
	},
		mf.DivProps(mf.ElementProps{Class: "rounded-[1.75rem] border border-base-300 bg-base-100 p-2 shadow-sm"},
			mf.Element("textarea", mf.ElementProps{
				ID:    "chat-prompt",
				Class: "min-h-16 w-full resize-none border-0 bg-transparent px-3 py-2 text-sm leading-6",
				Attrs: mf.Attrs{
					"name":        "prompt",
					"placeholder": "Message Marionette AI",
					"rows":        "2",
					"required":    "required",
					"aria-label":  "Message",
				},
			}),
			mf.DivProps(mf.ElementProps{Class: "flex items-center justify-end px-1"},
				mf.Element("button", mf.ElementProps{
					Class: "btn btn-neutral btn-sm rounded-full px-5",
					Attrs: mf.Attrs{"type": "submit"},
				}, mf.Text("Send")),
			),
		),
	)
}

func resetForm() mf.Node {
	return mf.ActionForm(mf.ActionFormProps{
		Action: "/chat/reset",
		Target: "#chat-panel",
		Swap:   "outerHTML",
		Props:  mf.ComponentProps{Class: "m-0"},
	},
		mf.Element("button", mf.ElementProps{
			Class: "btn btn-ghost btn-sm rounded-full",
			Attrs: mf.Attrs{"type": "submit"},
		}, mf.Text("Reset conversation")),
	)
}

func welcomeMessage() chatMessage {
	return chatMessage{
		ID:      1,
		Role:    "assistant",
		Name:    "Marionette AI",
		Content: "Hello. Ask me anything.",
	}
}

func demoReply(prompt string) string {
	lower := strings.ToLower(prompt)
	switch {
	case strings.Contains(lower, "htmx") || strings.Contains(lower, "stream") || strings.Contains(lower, "swap"):
		return "In Marionette, ActionForm Target and Swap let a POST response update only the selected region. This sample combines those partial updates with a StreamAction SSE endpoint so the assistant reply is appended token by token without client-side state."
	case strings.Contains(lower, "sales"):
		return "For sales data, consider showing KPIs in cards, details in a table, and trends in a chart. You could also extract conditions from the chat and apply them to DataQueryState."
	case strings.Contains(lower, "api") || strings.Contains(lower, "llm") || strings.Contains(lower, "ai"):
		return "To connect an external LLM, replace this demoReply function with a streaming API client and feed chunks through the text stream APIs. The UI can keep the same Node structure while the reply chunks come from the model."
	default:
		return fmt.Sprintf("I reviewed %q. With Marionette, input handling, state updates, partial updates, and streaming-style UI feedback can stay in one Go flow.", prompt)
	}
}

func advanceStream(ctx *mb.Context) (int, mb.TextStreamStep, bool) {
	step := ctx.AdvanceTextStream("chat-reply")
	if !step.Active {
		return 0, step, false
	}

	messageID := 0
	wasThinking := false
	ctx.UpdateGlobal("messages", func(old any) any {
		messages := cloneMessages(old).([]chatMessage)
		for i := range messages {
			if messages[i].Streaming {
				messageID = messages[i].ID
				wasThinking = messages[i].Thinking
				messages[i].Content = step.Content
				messages[i].Streaming = !step.Done
				messages[i].Thinking = false
				break
			}
		}
		return messages
	})
	return messageID, step, wasThinking
}

func streamDelta(messageID int, step mb.TextStreamStep, wasThinking bool) mf.Node {
	chunk := stdhtml.EscapeString(step.Delta)
	contentID := fmt.Sprintf("message-content-%d", messageID)
	cursorID := fmt.Sprintf("message-cursor-%d", messageID)
	statusID := fmt.Sprintf("message-status-%d", messageID)
	updates := ""
	if wasThinking {
		thinkingID := fmt.Sprintf("message-thinking-%d", messageID)
		updates += fmt.Sprintf(`<span id="%s" hx-swap-oob="outerHTML"></span><span id="%s" class="opacity-70" hx-swap-oob="outerHTML">▌</span><span id="%s" class="badge badge-info badge-xs ml-2" hx-swap-oob="outerHTML">SSE streaming</span>`, thinkingID, cursorID, statusID)
	}
	if step.Done {
		updates += fmt.Sprintf(`<span id="%s" class="badge badge-success badge-xs ml-2" hx-swap-oob="outerHTML">Complete</span><span id="%s" hx-swap-oob="outerHTML"></span>`, statusID, cursorID)
	}
	return mf.Raw(fmt.Sprintf(`<span hx-swap-oob="beforeend:#%s">%s</span>%s`, contentID, chunk, updates))
}

func hasStreamingMessage(messages []chatMessage) bool {
	for _, msg := range messages {
		if msg.Streaming {
			return true
		}
	}
	return false
}

func cloneMessages(old any) any {
	messages, _ := old.([]chatMessage)
	cloned := make([]chatMessage, len(messages))
	copy(cloned, messages)
	return cloned
}
