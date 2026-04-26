package marionette

import (
	"strings"
	"testing"
)

func TestButtonRenderUsesHTMXMarkup(t *testing.T) {
	html, err := Button("Increment").OnClick("counter/increment").Target("#app").Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	if got == "" {
		t.Fatal("expected non-empty button html")
	}
	if want := `hx-post="/counter/increment"`; !strings.Contains(got, want) {
		t.Fatalf("expected %q in %q", want, got)
	}
}

func TestElementRenderEscapesText(t *testing.T) {
	html, err := Text(`<script>alert(1)</script>`).Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	if strings.Contains(got, `<script>alert(1)</script>`) {
		t.Fatalf("expected escaped content, got %q", got)
	}
}
