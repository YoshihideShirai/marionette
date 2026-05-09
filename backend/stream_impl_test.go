package backend

import "testing"

func TestTextStreamAdvancesByConfiguredChunks(t *testing.T) {
	app := New()
	ctx := &Context{app: app}

	ctx.StartTextStream(TextStreamOptions{Name: "reply", Text: "one two three four five", ChunkSize: 2})
	if !ctx.TextStreamActive("reply") {
		t.Fatalf("expected text stream to be active after start")
	}

	first := ctx.AdvanceTextStream("reply")
	if !first.Active || first.Done || first.Content != "one two" || first.Cursor != 2 || first.Total != 5 {
		t.Fatalf("unexpected first step: %+v", first)
	}

	second := ctx.AdvanceTextStream("reply")
	if !second.Active || second.Done || second.Content != "one two three four" || second.Cursor != 4 || second.Total != 5 {
		t.Fatalf("unexpected second step: %+v", second)
	}

	final := ctx.AdvanceTextStream("reply")
	if !final.Active || !final.Done || final.Content != "one two three four five" || final.Cursor != 5 || final.Total != 5 {
		t.Fatalf("unexpected final step: %+v", final)
	}
	if ctx.TextStreamActive("reply") {
		t.Fatalf("expected stream to be inactive after final step")
	}
}

func TestTextStreamResetClearsState(t *testing.T) {
	app := New()
	ctx := &Context{app: app}

	ctx.StartTextStream(TextStreamOptions{Name: "reply", Text: "one two three", ChunkSize: 1})
	ctx.ResetTextStream("reply")

	if ctx.TextStreamActive("reply") {
		t.Fatalf("expected stream to be inactive after reset")
	}
	step := ctx.AdvanceTextStream("reply")
	if step.Active || !step.Done || step.Content != "" {
		t.Fatalf("unexpected step after reset: %+v", step)
	}
}

func TestTextStreamIgnoresBlankNames(t *testing.T) {
	app := New()
	ctx := &Context{app: app}

	ctx.StartTextStream(TextStreamOptions{Name: " ", Text: "ignored"})
	if ctx.TextStreamActive(" ") {
		t.Fatalf("expected blank stream name to stay inactive")
	}
	step := ctx.AdvanceTextStream(" ")
	if step.Active || !step.Done {
		t.Fatalf("unexpected blank-name step: %+v", step)
	}
}
