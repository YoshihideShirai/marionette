package marionette

import (
	"fmt"
	"html"
	"strings"
)

// Node is a declarative UI element that can render itself as HTML.
type Node interface {
	Render() string
}

type element struct {
	tag      string
	attrs    map[string]string
	children []Node
	text     string
}

func (e element) Render() string {
	var b strings.Builder
	b.WriteString("<")
	b.WriteString(e.tag)
	for k, v := range e.attrs {
		b.WriteString(" ")
		b.WriteString(k)
		b.WriteString(`="`)
		b.WriteString(html.EscapeString(v))
		b.WriteString(`"`)
	}
	b.WriteString(">")
	if e.text != "" {
		b.WriteString(html.EscapeString(e.text))
	}
	for _, c := range e.children {
		b.WriteString(c.Render())
	}
	b.WriteString("</")
	b.WriteString(e.tag)
	b.WriteString(">")
	return b.String()
}

// Raw allows trusted HTML snippets (e.g. full page shell).
type Raw string

func (r Raw) Render() string { return string(r) }

func Text(v string) Node {
	return element{tag: "span", text: v}
}

func Div(id string, children ...Node) Node {
	attrs := map[string]string{}
	if id != "" {
		attrs["id"] = id
	}
	return element{tag: "div", attrs: attrs, children: children}
}

func Column(children ...Node) Node {
	return element{tag: "div", attrs: map[string]string{"style": "display:flex;flex-direction:column;gap:12px;"}, children: children}
}

type button struct {
	label  string
	action string
	target string
}

func Button(label string) *button {
	return &button{label: label, target: "#app"}
}

func (b *button) OnClick(action string) *button {
	b.action = action
	return b
}

func (b *button) Target(selector string) *button {
	b.target = selector
	return b
}

func (b *button) Render() string {
	if b.action == "" {
		panic("button action is required; call OnClick")
	}
	return fmt.Sprintf(
		`<button hx-post="/%s" hx-target="%s" hx-swap="outerHTML">%s</button>`,
		html.EscapeString(b.action),
		html.EscapeString(b.target),
		html.EscapeString(b.label),
	)
}
