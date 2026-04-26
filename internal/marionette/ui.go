package marionette

import (
	"bytes"
	"html/template"
)

// Node is a declarative UI element that can render itself as safe HTML.
type Node interface {
	Render() (template.HTML, error)
}

type element struct {
	Tag      string
	Attrs    map[string]string
	Children []Node
	Text     string
}

var elementTmpl = template.Must(template.New("element").Parse(`<{{.Tag}}{{range $k, $v := .Attrs}} {{$k}}="{{$v}}"{{end}}>{{.Text}}{{.ChildrenHTML}}</{{.Tag}}>`))

func (e element) Render() (template.HTML, error) {
	children := make([]template.HTML, 0, len(e.Children))
	for _, child := range e.Children {
		r, err := child.Render()
		if err != nil {
			return "", err
		}
		children = append(children, r)
	}

	view := struct {
		Tag          string
		Attrs        map[string]string
		Text         string
		ChildrenHTML template.HTML
	}{
		Tag:          e.Tag,
		Attrs:        e.Attrs,
		Text:         e.Text,
		ChildrenHTML: joinHTML(children),
	}

	var b bytes.Buffer
	if err := elementTmpl.Execute(&b, view); err != nil {
		return "", err
	}
	return template.HTML(b.String()), nil
}

// Raw allows trusted HTML snippets (e.g. full page shell).
type Raw string

func (r Raw) Render() (template.HTML, error) { return template.HTML(r), nil }

func Text(v string) Node {
	return element{Tag: "span", Text: v}
}

func Div(id string, children ...Node) Node {
	return DivClass(id, "", children...)
}

func DivClass(id, className string, children ...Node) Node {
	attrs := map[string]string{}
	if id != "" {
		attrs["id"] = id
	}
	if className != "" {
		attrs["class"] = className
	}
	return element{Tag: "div", Attrs: attrs, Children: children}
}

func Column(children ...Node) Node {
	return element{Tag: "div", Attrs: map[string]string{"class": "flex flex-col gap-3"}, Children: children}
}

type button struct {
	Label   string
	Action  string
	TargetQ string
}

var buttonTmpl = template.Must(template.New("button").Parse(`<button class="btn btn-primary w-fit" hx-post="/{{.Action}}" hx-target="{{.TargetQ}}" hx-swap="outerHTML">{{.Label}}</button>`))

func Button(label string) *button {
	return &button{Label: label, TargetQ: "#app"}
}

func (b *button) OnClick(action string) *button {
	b.Action = action
	return b
}

func (b *button) TargetSelector(selector string) *button {
	b.TargetQ = selector
	return b
}

func (b *button) Target(selector string) *button {
	return b.TargetSelector(selector)
}

func (b *button) Render() (template.HTML, error) {
	var out bytes.Buffer
	if err := buttonTmpl.Execute(&out, b); err != nil {
		return "", err
	}
	return template.HTML(out.String()), nil
}

func joinHTML(parts []template.HTML) template.HTML {
	var b bytes.Buffer
	for _, p := range parts {
		b.WriteString(string(p))
	}
	return template.HTML(b.String())
}
