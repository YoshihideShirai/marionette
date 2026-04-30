// Package html provides low-level UI node constructors for advanced Marionette
// users who need custom markup below the component API.
package html

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"sort"
)

// Node is a declarative UI element that can render itself as safe HTML.
type Node interface {
	Render() (template.HTML, error)
}

// ElementNode is the rendered representation used by low-level element
// constructors.
type ElementNode struct {
	Tag      string
	Attrs    map[string]string
	Children []Node
	Text     string
}

// Attrs defines HTML attributes for low-level element constructors.
type Attrs map[string]string

// ElementProps defines common HTML element attributes while keeping class and
// id easy to scan at call sites.
type ElementProps struct {
	ID    string
	Class string
	Attrs Attrs
}

var tagPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]*$`)

func (e ElementNode) Render() (template.HTML, error) {
	if !tagPattern.MatchString(e.Tag) {
		return "", fmt.Errorf("invalid tag: %q", e.Tag)
	}

	children := make([]template.HTML, 0, len(e.Children))
	for _, child := range e.Children {
		r, err := child.Render()
		if err != nil {
			return "", err
		}
		children = append(children, r)
	}

	var b bytes.Buffer
	b.WriteString("<")
	b.WriteString(e.Tag)

	keys := make([]string, 0, len(e.Attrs))
	for k := range e.Attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(" ")
		b.WriteString(template.HTMLEscapeString(k))
		b.WriteString(`="`)
		b.WriteString(template.HTMLEscapeString(e.Attrs[k]))
		b.WriteString(`"`)
	}
	b.WriteString(">")
	b.WriteString(template.HTMLEscapeString(e.Text))
	b.WriteString(string(joinHTML(children)))
	b.WriteString("</")
	b.WriteString(e.Tag)
	b.WriteString(">")

	return template.HTML(b.String()), nil
}

// Raw allows trusted HTML snippets (e.g. full page shell).
type Raw string

func (r Raw) Render() (template.HTML, error) { return template.HTML(r), nil }

func Text(v string) Node {
	return ElementNode{Tag: "span", Text: v}
}

func Element(tag string, props ElementProps, children ...Node) Node {
	return ElementNode{Tag: tag, Attrs: elementAttrs(props), Children: children}
}

func Div(children ...Node) Node {
	return DivProps(ElementProps{}, children...)
}

func DivID(id string, children ...Node) Node {
	return DivProps(ElementProps{ID: id}, children...)
}

func DivClass(className string, children ...Node) Node {
	return DivProps(ElementProps{Class: className}, children...)
}

func DivAttrs(attrs Attrs, children ...Node) Node {
	return DivProps(ElementProps{Attrs: attrs}, children...)
}

func DivProps(props ElementProps, children ...Node) Node {
	return Element("div", props, children...)
}

func Span(children ...Node) Node                          { return Element("span", ElementProps{}, children...) }
func SpanProps(props ElementProps, children ...Node) Node { return Element("span", props, children...) }

func P(children ...Node) Node                          { return Element("p", ElementProps{}, children...) }
func PProps(props ElementProps, children ...Node) Node { return Element("p", props, children...) }

func H1(children ...Node) Node                          { return Element("h1", ElementProps{}, children...) }
func H1Props(props ElementProps, children ...Node) Node { return Element("h1", props, children...) }

func H2(children ...Node) Node                          { return Element("h2", ElementProps{}, children...) }
func H2Props(props ElementProps, children ...Node) Node { return Element("h2", props, children...) }

func H3(children ...Node) Node                          { return Element("h3", ElementProps{}, children...) }
func H3Props(props ElementProps, children ...Node) Node { return Element("h3", props, children...) }

func H4(children ...Node) Node                          { return Element("h4", ElementProps{}, children...) }
func H4Props(props ElementProps, children ...Node) Node { return Element("h4", props, children...) }

func Column(children ...Node) Node {
	return ElementNode{Tag: "div", Attrs: map[string]string{"class": "flex flex-col gap-3"}, Children: children}
}

func elementAttrs(props ElementProps) map[string]string {
	attrs := make(map[string]string, len(props.Attrs)+2)
	for key, value := range props.Attrs {
		attrs[key] = value
	}
	if props.ID != "" {
		attrs["id"] = props.ID
	}
	if props.Class != "" {
		attrs["class"] = joinClass(attrs["class"], props.Class)
	}
	return attrs
}

func joinHTML(parts []template.HTML) template.HTML {
	var b bytes.Buffer
	for _, p := range parts {
		b.WriteString(string(p))
	}
	return template.HTML(b.String())
}

func joinClass(parts ...string) string {
	var b bytes.Buffer
	for _, part := range parts {
		if part == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(part)
	}
	return b.String()
}
