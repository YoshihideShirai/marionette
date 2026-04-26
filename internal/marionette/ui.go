package marionette

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"sort"
	"strings"
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

var tagPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]*$`)

func (e element) Render() (template.HTML, error) {
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

type sidebar struct {
	Brand     string
	Title     string
	Items     []SidebarItem
	NoteTitle string
	NoteText  string
}

type SidebarItem struct {
	Label   string
	Href    string
	Current bool
}

func Sidebar(brand, title string, items ...SidebarItem) *sidebar {
	return &sidebar{Brand: brand, Title: title, Items: items}
}

func SidebarLink(label, href string) SidebarItem {
	return SidebarItem{Label: label, Href: href}
}

func (i SidebarItem) Active() SidebarItem {
	i.Current = true
	return i
}

func (s *sidebar) Note(title, text string) *sidebar {
	s.NoteTitle = title
	s.NoteText = text
	return s
}

func (s *sidebar) Render() (template.HTML, error) {
	children := []Node{
		element{
			Tag: "div",
			Attrs: map[string]string{
				"class": "mb-6",
			},
			Children: []Node{
				element{
					Tag:   "div",
					Attrs: map[string]string{"class": "text-sm font-semibold uppercase tracking-wide text-base-content/50"},
					Text:  s.Brand,
				},
				element{
					Tag:   "div",
					Attrs: map[string]string{"class": "text-lg font-bold"},
					Text:  s.Title,
				},
			},
		},
		s.renderNav(),
	}
	if s.NoteTitle != "" || s.NoteText != "" {
		children = append(children, element{
			Tag:   "div",
			Attrs: map[string]string{"class": "mt-6 rounded-box bg-base-200 p-3 text-sm text-base-content/70"},
			Children: []Node{
				element{Tag: "div", Attrs: map[string]string{"class": "font-medium text-base-content"}, Text: s.NoteTitle},
				element{Tag: "div", Text: s.NoteText},
			},
		})
	}

	return element{
		Tag:      "aside",
		Attrs:    map[string]string{"class": "rounded-box border border-base-300 bg-base-100 p-4 shadow-sm lg:min-h-[calc(100vh-3rem)]"},
		Children: children,
	}.Render()
}

func (s *sidebar) renderNav() Node {
	items := make([]Node, 0, len(s.Items))
	for _, item := range s.Items {
		href := item.Href
		if href == "" {
			href = "#"
		}
		className := "btn btn-ghost justify-start text-base-content/70"
		if item.Current {
			className = "btn btn-primary justify-start"
		}
		items = append(items, element{
			Tag:   "a",
			Attrs: map[string]string{"class": className, "href": href},
			Text:  item.Label,
		})
	}
	return element{
		Tag:      "nav",
		Attrs:    map[string]string{"class": "flex flex-col gap-1"},
		Children: items,
	}
}

type form struct {
	Action   string
	TargetQ  string
	Children []Node
}

func Form(action string, children ...Node) *form {
	return &form{Action: action, TargetQ: "#app", Children: children}
}

func (f *form) Target(selector string) *form {
	f.TargetQ = selector
	return f
}

func (f *form) Render() (template.HTML, error) {
	return element{
		Tag: "form",
		Attrs: map[string]string{
			"class":     "flex flex-col gap-3",
			"hx-post":   actionPath(f.Action),
			"hx-target": f.TargetQ,
			"hx-swap":   "outerHTML",
		},
		Children: f.Children,
	}.Render()
}

func Input(name, value string) Node {
	return element{
		Tag: "input",
		Attrs: map[string]string{
			"class": "input input-bordered w-full",
			"name":  name,
			"type":  "text",
			"value": value,
		},
	}
}

func Submit(label string) Node {
	return element{
		Tag: "button",
		Attrs: map[string]string{
			"class": "btn btn-primary w-fit",
			"type":  "submit",
		},
		Text: label,
	}
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
	return b.Post(action)
}

func (b *button) Post(action string) *button {
	b.Action = strings.TrimPrefix(action, "/")
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

func actionPath(action string) string {
	if strings.HasPrefix(action, "/") {
		return action
	}
	return "/" + action
}

func joinHTML(parts []template.HTML) template.HTML {
	var b bytes.Buffer
	for _, p := range parts {
		b.WriteString(string(p))
	}
	return template.HTML(b.String())
}
