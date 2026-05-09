package main

import (
	"fmt"
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

type chatMessage struct {
	ID        int
	Role      string
	Name      string
	Content   string
	Streaming bool
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
	app.AddStyle(`
		#marionette-root { width: min(100%, 72rem); }
		.ai-chat-message { scroll-margin-block: 1rem; }
	`)

	app.Page("/", func(ctx *mb.Context) mf.Node {
		return page(ctx)
	}, mb.WithTitle("AI Chat Sample"))

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
			ChunkSize: 5,
		})
		ctx.UpdateGlobal("messages", func(old any) any {
			messages := cloneMessages(old).([]chatMessage)
			messages = append(messages,
				chatMessage{ID: userID, Role: "user", Name: "You", Content: prompt},
				chatMessage{ID: assistantID, Role: "assistant", Name: "Marionette AI", Content: "Streaming response…", Streaming: true},
			)
			return messages
		})
		return chatPanel(ctx)
	})

	app.Action("chat/stream", func(ctx *mb.Context) mf.Node {
		advanceStream(ctx)
		return chatPanel(ctx)
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
	return mf.Container(mf.ContainerProps{MaxWidth: "6xl", Centered: true},
		mf.Stack(mf.StackProps{Direction: "column", Gap: "6"},
			mf.PageHeader(mf.PageHeaderProps{
				Title:       "AI Chat Sample",
				Description: "A Go-only demo for building a chat UI with server-driven state, htmx partial updates, and simulated streaming replies. It uses a sample reply generator instead of calling an external AI API.",
			}),
			mf.Grid(mf.GridProps{Columns: "3", Gap: "lg"},
				mf.DivProps(mf.ElementProps{Class: "lg:col-span-2"}, chatPanel(ctx)),
				sidebar(),
			),
		),
	)
}

func chatPanel(ctx *mb.Context) mf.Node {
	messages := ctx.GetGlobalSnapshot("messages", cloneMessages).([]chatMessage)
	errorMessage, _ := ctx.GetGlobal("chatError").(string)

	children := []mf.Node{
		conversation(messages),
	}
	if hasStreamingMessage(messages) {
		children = append(children, streamTrigger())
	}
	if strings.TrimSpace(errorMessage) != "" {
		children = append(children, mf.Alert(mf.AlertProps{
			Title:       "Input error",
			Description: errorMessage,
			Props:       mf.ComponentProps{Class: "alert-warning"},
		}))
	}
	children = append(children, promptForm(), resetForm())

	return mf.Region(mf.RegionProps{ID: "chat-panel"},
		mf.Card(mf.CardProps{
			Title:       "Demo conversation",
			Description: "The server owns the conversation state, swaps this card after each action, and streams mock assistant replies in small chunks.",
			Props:       mf.ComponentProps{Class: "border border-base-300"},
		}, children...),
	)
}

func conversation(messages []chatMessage) mf.Node {
	items := make([]mf.Node, 0, len(messages))
	for _, msg := range messages {
		items = append(items, messageBubble(msg))
	}
	return mf.DivProps(mf.ElementProps{Class: "card-body max-h-[34rem] overflow-y-auto space-y-4 bg-base-200/60"}, items...)
}

func messageBubble(msg chatMessage) mf.Node {
	alignClass := "chat-start"
	bubbleClass := "chat-bubble chat-bubble-secondary"
	if msg.Role == "user" {
		alignClass = "chat-end"
		bubbleClass = "chat-bubble chat-bubble-primary"
	}

	headerChildren := []mf.Node{mf.Text(msg.Name)}
	if msg.Streaming {
		headerChildren = append(headerChildren,
			mf.SpanProps(mf.ElementProps{Class: "badge badge-info badge-xs ml-2"}, mf.Text("Streaming")),
		)
	}

	return mf.DivProps(mf.ElementProps{Class: "ai-chat-message chat " + alignClass},
		mf.DivProps(mf.ElementProps{Class: "chat-header text-xs opacity-70"}, headerChildren...),
		mf.DivProps(mf.ElementProps{Class: bubbleClass}, mf.Text(msg.Content)),
	)
}

func streamTrigger() mf.Node {
	return mf.StreamTrigger(mf.StreamTriggerProps{
		Action: "/chat/stream",
		Target: "#chat-panel",
		Swap:   "outerHTML",
		Delay:  "350ms",
	})
}

func promptForm() mf.Node {
	return mf.ActionForm(mf.ActionFormProps{
		Action: "/chat/send",
		Target: "#chat-panel",
		Swap:   "outerHTML",
		Props:  mf.ComponentProps{Class: "card-body gap-4 border-t border-base-300"},
	},
		mf.FormRow(mf.FormRowProps{
			ID:          "chat-prompt",
			Label:       "Message",
			Description: "Examples: Build a sales summary UI / Explain htmx streaming",
			Required:    true,
			Control: mf.Textarea(mf.TextareaProps{
				ID:          "chat-prompt",
				Name:        "prompt",
				Placeholder: "Type what you want to ask the AI",
				Rows:        4,
				Required:    true,
			}),
		}),
		mf.Actions(mf.ActionsProps{Props: mf.ComponentProps{Class: "justify-end"}},
			mf.SubmitButton("Send message", mf.ComponentProps{Class: "btn-primary"}),
		),
	)
}

func resetForm() mf.Node {
	return mf.ActionForm(mf.ActionFormProps{
		Action: "/chat/reset",
		Target: "#chat-panel",
		Swap:   "outerHTML",
		Props:  mf.ComponentProps{Class: "card-body pt-0 items-end"},
	},
		mf.SubmitButton("Reset conversation", mf.ComponentProps{Class: "btn-ghost btn-sm"}),
	)
}

func sidebar() mf.Node {
	return mf.Stack(mf.StackProps{Direction: "column", Gap: "4"},
		mf.Card(mf.CardProps{
			Title:       "What this demonstrates",
			Description: "Build an AI-chat-style operations UI with standard Marionette components.",
			Props:       mf.ComponentProps{Class: "border border-base-300"},
		},
			mf.DivProps(mf.ElementProps{Class: "card-body pt-0"},
				mf.UlProps(mf.ElementProps{Class: "list-disc space-y-2 pl-5 text-sm text-base-content/80"},
					mf.Li(mf.Text("Streaming mock replies render chunk by chunk")),
					mf.Li(mf.Text("Conversation history is stored in Go app state")),
					mf.Li(mf.Text("Form submissions update only the card with htmx")),
					mf.Li(mf.Text("Runs locally without an external API key")),
				),
			),
		),
		mf.Alert(mf.AlertProps{
			Title:       "Integration note",
			Description: "For production, feed real LLM chunks through the text stream APIs and load API keys from environment variables or secret management.",
			Props:       mf.ComponentProps{Class: "alert-info"},
		}),
	)
}

func welcomeMessage() chatMessage {
	return chatMessage{
		ID:      1,
		Role:    "assistant",
		Name:    "Marionette AI",
		Content: "Hello. This is an AI chat sample demo. Send a message to update server-side state, stream a mock reply, and re-render only the conversation card.",
	}
}

func demoReply(prompt string) string {
	lower := strings.ToLower(prompt)
	switch {
	case strings.Contains(lower, "htmx") || strings.Contains(lower, "stream") || strings.Contains(lower, "swap"):
		return "In Marionette, ActionForm Target and Swap let a POST response update only the selected region. This sample combines those partial updates with a polling trigger so the assistant reply appears chunk by chunk."
	case strings.Contains(lower, "sales"):
		return "For sales data, consider showing KPIs in cards, details in a table, and trends in a chart. You could also extract conditions from the chat and apply them to DataQueryState."
	case strings.Contains(lower, "api") || strings.Contains(lower, "llm") || strings.Contains(lower, "ai"):
		return "To connect an external LLM, replace this demoReply function with a streaming API client and feed chunks through the text stream APIs. The UI can keep the same Node structure while the reply chunks come from the model."
	default:
		return fmt.Sprintf("I reviewed %q. With Marionette, input handling, state updates, partial updates, and streaming-style UI feedback can stay in one Go flow.", prompt)
	}
}

func advanceStream(ctx *mb.Context) {
	step := ctx.AdvanceTextStream("chat-reply")
	if !step.Active {
		return
	}

	content := step.Content
	if !step.Done {
		content += " ▌"
	}

	ctx.UpdateGlobal("messages", func(old any) any {
		messages := cloneMessages(old).([]chatMessage)
		for i := range messages {
			if messages[i].Streaming {
				messages[i].Content = content
				messages[i].Streaming = !step.Done
				break
			}
		}
		return messages
	})
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
