// Package html provides low-level HTML/HTMX primitives and tag-building
// utilities used by higher-level frontend packages.
//
// This package is primarily an internal implementation layer. Most Marionette
// applications should depend on user-facing APIs from package frontend (and
// package frontend/daisyui when using daisyUI components) instead of consuming
// these primitives directly.
package html

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

// SidebarItem is the low-level representation for sidebar links used by
// primitive and sample navigation builders.
type SidebarItem struct {
	Label   string
	Href    string
	Current bool
}

// Active returns a copy of the item marked as the current navigation entry.
func (i SidebarItem) Active() SidebarItem {
	i.Current = true
	return i
}

// TableRowData stores primitive table cells.
type TableRowData struct {
	Cells []Node
}

// Table is a low-level table node for small HTMX samples and starter UIs.
type Table struct {
	Headers []string
	Rows    []TableRowData
}

// HTMXTable returns a primitive table node used by HTMX-focused examples.
func HTMXTable(headers []string, rows ...TableRowData) Node {
	return Table{Headers: headers, Rows: rows}
}

// TableRow groups cell nodes for HTMXTable.
func TableRow(cells ...Node) TableRowData {
	return TableRowData{Cells: cells}
}

func (t Table) Render() (template.HTML, error) {
	headerCells := make([]Node, 0, len(t.Headers))
	for _, header := range t.Headers {
		headerCells = append(headerCells, ElementNode{Tag: "th", Text: header})
	}

	bodyRows := make([]Node, 0, len(t.Rows))
	for _, row := range t.Rows {
		cells := make([]Node, 0, len(row.Cells))
		for _, cell := range row.Cells {
			cells = append(cells, ElementNode{Tag: "td", Children: []Node{cell}})
		}
		bodyRows = append(bodyRows, ElementNode{Tag: "tr", Children: cells})
	}

	return ElementNode{
		Tag:   "table",
		Attrs: map[string]string{"class": "table"},
		Children: []Node{
			ElementNode{
				Tag: "thead",
				Children: []Node{
					ElementNode{Tag: "tr", Children: headerCells},
				},
			},
			ElementNode{Tag: "tbody", Children: bodyRows},
		},
	}.Render()
}

// Sidebar is a low-level sidebar node for samples and starter layouts.
type Sidebar struct {
	Brand     string
	Title     string
	Items     []SidebarItem
	NoteTitle string
	NoteText  string
}

// NewSidebar creates a low-level sidebar node.
func NewSidebar(brand, title string, items ...SidebarItem) *Sidebar {
	return &Sidebar{Brand: brand, Title: title, Items: items}
}

// SidebarLink creates a low-level sidebar link item.
func SidebarLink(label, href string) SidebarItem {
	return SidebarItem{Label: label, Href: href}
}

// Note returns the sidebar with an optional note block configured.
func (s *Sidebar) Note(title, text string) *Sidebar {
	s.NoteTitle = title
	s.NoteText = text
	return s
}

func (s *Sidebar) Render() (template.HTML, error) {
	children := []Node{
		ElementNode{
			Tag:   "div",
			Attrs: map[string]string{"class": "mb-6"},
			Children: []Node{
				ElementNode{Tag: "div", Attrs: map[string]string{"class": "text-sm font-semibold uppercase tracking-wide text-base-content/50"}, Text: s.Brand},
				ElementNode{Tag: "div", Attrs: map[string]string{"class": "text-lg font-bold"}, Text: s.Title},
			},
		},
		s.renderNav(),
	}
	if s.NoteTitle != "" || s.NoteText != "" {
		children = append(children, ElementNode{
			Tag:   "div",
			Attrs: map[string]string{"class": "mt-6 rounded-box bg-base-200 p-3 text-sm text-base-content/70"},
			Children: []Node{
				ElementNode{Tag: "div", Attrs: map[string]string{"class": "font-medium text-base-content"}, Text: s.NoteTitle},
				ElementNode{Tag: "div", Text: s.NoteText},
			},
		})
	}

	return ElementNode{
		Tag:      "aside",
		Attrs:    map[string]string{"class": "rounded-box border border-base-300 bg-base-100 p-4 shadow-sm lg:min-h-[calc(100vh-3rem)]"},
		Children: children,
	}.Render()
}

func (s *Sidebar) renderNav() Node {
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
		items = append(items, ElementNode{Tag: "a", Attrs: map[string]string{"class": className, "href": href}, Text: item.Label})
	}
	return ElementNode{Tag: "nav", Attrs: map[string]string{"class": "flex flex-col gap-1"}, Children: items}
}

// Form is a primitive HTMX form node.
type Form struct {
	Action   string
	TargetQ  string
	Children []Node
}

// NewForm creates a primitive HTMX form that posts to action and targets #app by default.
func NewForm(action string, children ...Node) *Form {
	return &Form{Action: action, TargetQ: "#app", Children: children}
}

// Target configures the HTMX target selector.
func (f *Form) Target(selector string) *Form {
	f.TargetQ = selector
	return f
}

func (f *Form) Render() (template.HTML, error) {
	return ElementNode{
		Tag: "form",
		Attrs: map[string]string{
			"class":     "flex flex-col gap-3",
			"hx-post":   ActionPath(f.Action),
			"hx-target": f.TargetQ,
			"hx-swap":   "outerHTML",
		},
		Children: f.Children,
	}.Render()
}

// Input creates a primitive text input node.
func Input(name, value string, className ...string) Node {
	return ElementNode{Tag: "input", Attrs: map[string]string{
		"class": joinClass(append([]string{"input input-bordered w-full"}, className...)...),
		"name":  name,
		"type":  "text",
		"value": value,
	}}
}

// FileUpload creates a primitive file input node.
func FileUpload(name string, required bool, className ...string) Node {
	attrs := map[string]string{
		"class": joinClass(append([]string{"input input-bordered w-full"}, className...)...),
		"name":  name,
		"type":  "file",
		"value": "",
	}
	if required {
		attrs["required"] = "required"
	}
	return ElementNode{Tag: "input", Attrs: attrs}
}

// HiddenInput creates a primitive hidden input node.
func HiddenInput(name, value string) Node {
	return ElementNode{Tag: "input", Attrs: map[string]string{"name": name, "type": "hidden", "value": value}}
}

// Submit creates a primitive submit button node.
func Submit(label string) Node {
	return ElementNode{Tag: "button", Attrs: map[string]string{"class": "btn btn-primary w-fit", "type": "submit"}, Text: label}
}

// Button is a primitive HTMX button node.
type Button struct {
	Label   string
	Action  string
	TargetQ string
}

// HTMXButton creates a primitive HTMX button that targets #app by default.
func HTMXButton(label string) *Button {
	return &Button{Label: label, TargetQ: "#app"}
}

// OnClick configures the button POST action.
func (b *Button) OnClick(action string) *Button { return b.Post(action) }

// Post configures the button POST action.
func (b *Button) Post(action string) *Button {
	b.Action = action
	return b
}

// TargetSelector configures the button HTMX target selector.
func (b *Button) TargetSelector(selector string) *Button {
	b.TargetQ = selector
	return b
}

// Target configures the button HTMX target selector.
func (b *Button) Target(selector string) *Button { return b.TargetSelector(selector) }

func (b *Button) Render() (template.HTML, error) {
	return ElementNode{Tag: "button", Attrs: map[string]string{
		"class":     "btn btn-primary w-fit",
		"hx-post":   ActionPath(b.Action),
		"hx-target": b.TargetQ,
		"hx-swap":   "outerHTML",
	}, Text: b.Label}.Render()
}

// ActionPath normalizes HTMX action paths by ensuring one leading slash.
func ActionPath(action string) string {
	trimmed := strings.TrimSpace(action)
	if strings.HasPrefix(trimmed, "/") {
		return trimmed
	}
	return "/" + strings.TrimLeft(trimmed, "/")
}

var tagPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]*$`)

func (e ElementNode) Render() (template.HTML, error) {
	if !tagPattern.MatchString(e.Tag) {
		return "", fmt.Errorf("invalid tag: %q", e.Tag)
	}

	children := make([]template.HTML, 0, len(e.Children))
	for _, child := range e.Children {
		if child == nil {
			continue
		}
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

func P(children ...Node) Node                               { return Element("p", ElementProps{}, children...) }
func AnchorProps(props ElementProps, children ...Node) Node { return Element("a", props, children...) }
func AsideProps(props ElementProps, children ...Node) Node {
	return Element("aside", props, children...)
}
func DescriptionListProps(props ElementProps, children ...Node) Node {
	return Element("dl", props, children...)
}
func DescriptionTermProps(props ElementProps, children ...Node) Node {
	return Element("dt", props, children...)
}
func DescriptionDetailsProps(props ElementProps, children ...Node) Node {
	return Element("dd", props, children...)
}

func PProps(props ElementProps, children ...Node) Node { return Element("p", props, children...) }

func LabelElement(children ...Node) Node {
	return LabelElementProps(ElementProps{}, children...)
}
func LabelElementProps(props ElementProps, children ...Node) Node {
	return Element("label", props, children...)
}

func InputElement(props ElementProps) Node { return Element("input", props) }

func Ul(children ...Node) Node                          { return UlProps(ElementProps{}, children...) }
func UlProps(props ElementProps, children ...Node) Node { return Element("ul", props, children...) }

func Li(children ...Node) Node                          { return LiProps(ElementProps{}, children...) }
func LiProps(props ElementProps, children ...Node) Node { return Element("li", props, children...) }

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
