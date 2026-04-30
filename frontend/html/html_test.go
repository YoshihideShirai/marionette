package html

import (
	"strings"
	"testing"
)

func TestElementEscapesTextAndAttrs(t *testing.T) {
	rendered, err := DivProps(
		ElementProps{
			ID:    "panel",
			Class: "card",
			Attrs: Attrs{
				"data-value": `<script>alert("x")</script>`,
			},
		},
		Text(`<b>Hello</b>`),
	).Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	got := string(rendered)
	if !strings.Contains(got, `data-value="&lt;script&gt;alert(&#34;x&#34;)&lt;/script&gt;"`) {
		t.Fatalf("rendered attrs were not escaped: %s", got)
	}
	if !strings.Contains(got, `&lt;b&gt;Hello&lt;/b&gt;`) {
		t.Fatalf("rendered text was not escaped: %s", got)
	}
}

func TestElementRejectsInvalidTag(t *testing.T) {
	_, err := Element("div onclick=alert(1)", ElementProps{}).Render()
	if err == nil {
		t.Fatal("Render() error = nil, want invalid tag error")
	}
}
