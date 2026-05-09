package frontend

import (
	"errors"
	"html/template"
	"strings"
	"testing"
)

type fakeRenderNode struct {
	html template.HTML
	err  error
}

func (n fakeRenderNode) Render() (template.HTML, error) {
	return n.html, n.err
}

func TestMarkdownRendersGoldmarkHTML(t *testing.T) {
	got := renderFrontendNodeForTest(t, Markdown(MarkdownProps{
		Content: "# Title\n\n**bold**",
		Props:   ComponentProps{Class: "doc-body"},
	}))

	for _, want := range []string{
		`<article class="prose doc-body">`,
		`<h1>Title</h1>`,
		`<p><strong>bold</strong></p>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestContainerRendersLayoutChildren(t *testing.T) {
	got := renderFrontendNodeForTest(t, Container(
		ContainerProps{MaxWidth: "sm", Padding: "sm", Centered: true, Props: ComponentProps{Class: "dashboard"}},
		Text("Hello"),
		Element("strong", ElementProps{}, Text("World")),
	))

	for _, want := range []string{
		`<div class="max-w-3xl p-3 mx-auto dashboard">`,
		`<span>Hello</span>`,
		`<strong><span>World</span></strong>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestRenderNodesHandlesNilChild(t *testing.T) {
	rendered, err := renderNodes([]Node{Text("before"), nil, Text("after")})
	if err != nil {
		t.Fatalf("render nodes failed: %v", err)
	}
	if len(rendered) != 3 {
		t.Fatalf("expected 3 rendered children, got %d", len(rendered))
	}
	if rendered[1] != "" {
		t.Fatalf("expected nil child to render as empty HTML, got %q", rendered[1])
	}
}

func TestRenderNodesReturnsChildError(t *testing.T) {
	wantErr := errors.New("child render failed")
	_, err := renderNodes([]Node{Text("ok"), fakeRenderNode{err: wantErr}})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestSwitchWithVariantsAndRadioGroupWithVariants(t *testing.T) {
	switchHTML := renderFrontendNodeForTest(t, SwitchWithVariants("notify", "1", "Notify me", true, SwitchVariantProps{
		Variant:  VariantPrimary,
		Size:     SizeSM,
		Class:    "custom-switch",
		Disabled: true,
	}))
	for _, want := range []string{
		`name="notify"`,
		`value="1"`,
		`class="toggle toggle-primary toggle-sm custom-switch"`,
		`checked`,
		`disabled`,
		`<span class="label-text">Notify me</span>`,
	} {
		if !strings.Contains(switchHTML, want) {
			t.Fatalf("expected switch HTML to contain %q in %q", want, switchHTML)
		}
	}

	radioHTML := renderFrontendNodeForTest(t, RadioGroupWithVariants("tier", "Plan tier", []RadioItem{
		{Label: "Team", Value: "team", Checked: true},
		{Label: "Enterprise", Value: "enterprise", Disabled: true},
	}, RadioVariantProps{
		Variant:  VariantSecondary,
		Size:     SizeLG,
		Class:    "custom-radio",
		Disabled: true,
	}))
	for _, want := range []string{
		`aria-label="Plan tier"`,
		`name="tier"`,
		`value="team"`,
		`value="enterprise"`,
		`class="radio radio-secondary radio-lg custom-radio"`,
		`checked`,
		`disabled`,
	} {
		if !strings.Contains(radioHTML, want) {
			t.Fatalf("expected radio group HTML to contain %q in %q", want, radioHTML)
		}
	}
}
