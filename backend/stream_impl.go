package backend

import "strings"

const textStreamStatePrefix = "__marionette_text_stream:"

// TextStreamOptions configures a server-side text stream owned by Context.
type TextStreamOptions struct {
	// Name identifies the stream within the app state.
	Name string
	// Text is the complete text that will be revealed chunk by chunk.
	Text string
	// ChunkSize controls how many word chunks are revealed by each AdvanceTextStream call.
	// Values less than 1 use the default chunk size.
	ChunkSize int
}

// TextStreamStep is the result of advancing a server-side text stream.
type TextStreamStep struct {
	Content string
	Done    bool
	Active  bool
	Cursor  int
	Total   int
}

type textStreamState struct {
	Text      string
	Cursor    int
	ChunkSize int
}

const defaultTextStreamChunkSize = 5

// StartTextStream stores a text stream in app state so later actions can advance it.
func (c *Context) StartTextStream(options TextStreamOptions) {
	name := strings.TrimSpace(options.Name)
	if c == nil || name == "" {
		return
	}
	chunkSize := options.ChunkSize
	if chunkSize <= 0 {
		chunkSize = defaultTextStreamChunkSize
	}
	c.SetGlobal(textStreamKey(name), textStreamState{Text: options.Text, ChunkSize: chunkSize})
}

// AdvanceTextStream reveals the next chunk of a server-side text stream.
func (c *Context) AdvanceTextStream(name string) TextStreamStep {
	key := textStreamKey(name)
	if key == "" || c == nil {
		return TextStreamStep{Done: true}
	}

	var step TextStreamStep
	c.UpdateGlobal(key, func(old any) any {
		state, ok := old.(textStreamState)
		if !ok {
			step = TextStreamStep{Done: true}
			return old
		}

		chunks := splitTextStreamChunks(state.Text)
		chunkSize := state.ChunkSize
		if chunkSize <= 0 {
			chunkSize = defaultTextStreamChunkSize
		}
		next := state.Cursor + chunkSize
		if next > len(chunks) {
			next = len(chunks)
		}
		if next <= 0 {
			step = TextStreamStep{Done: true, Active: true, Total: len(chunks)}
			return nil
		}

		done := next >= len(chunks)
		step = TextStreamStep{
			Content: strings.Join(chunks[:next], " "),
			Done:    done,
			Active:  true,
			Cursor:  next,
			Total:   len(chunks),
		}
		if done {
			return nil
		}
		state.Cursor = next
		state.ChunkSize = chunkSize
		return state
	})
	return step
}

// ResetTextStream clears a server-side text stream from app state.
func (c *Context) ResetTextStream(name string) {
	key := textStreamKey(name)
	if key == "" || c == nil {
		return
	}
	c.SetGlobal(key, nil)
}

// TextStreamActive reports whether a named server-side text stream can still advance.
func (c *Context) TextStreamActive(name string) bool {
	key := textStreamKey(name)
	if key == "" || c == nil {
		return false
	}
	_, ok := c.GetGlobal(key).(textStreamState)
	return ok
}

func textStreamKey(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	return textStreamStatePrefix + name
}

func splitTextStreamChunks(text string) []string {
	chunks := strings.Fields(text)
	if len(chunks) == 0 {
		return []string{text}
	}
	return chunks
}
