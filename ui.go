package marionette

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"regexp"
	"sort"
	"strings"

	dataframeimports "github.com/rocketlaunchr/dataframe-go/imports"
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

func Element(tag string, props ElementProps, children ...Node) Node {
	return element{Tag: tag, Attrs: elementAttrs(props), Children: children}
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

func Column(children ...Node) Node {
	return element{Tag: "div", Attrs: map[string]string{"class": "flex flex-col gap-3"}, Children: children}
}

type table struct {
	Headers []string
	Rows    []TableRowData
}

type TableRowData struct {
	Cells []Node
}

func HTMXTable(headers []string, rows ...TableRowData) Node {
	return table{Headers: headers, Rows: rows}
}

func TableRow(cells ...Node) TableRowData {
	return TableRowData{Cells: cells}
}

func UIDataFrameFromCSV(r io.ReadSeeker, props TableProps, opts ...dataframeimports.CSVLoadOptions) (Node, error) {
	if r == nil {
		return nil, fmt.Errorf("csv reader is nil")
	}
	df, err := dataframeimports.LoadFromCSV(context.Background(), r, opts...)
	if err != nil {
		return nil, err
	}
	return DataFrame(df, props), nil
}

func UIDataFrameFromTSV(r io.ReadSeeker, props TableProps, opts ...dataframeimports.CSVLoadOptions) (Node, error) {
	tsvOpts := make([]dataframeimports.CSVLoadOptions, len(opts))
	copy(tsvOpts, opts)
	if len(tsvOpts) == 0 {
		tsvOpts = append(tsvOpts, dataframeimports.CSVLoadOptions{Comma: '\t'})
	} else if tsvOpts[0].Comma == 0 {
		tsvOpts[0].Comma = '\t'
	}
	return UIDataFrameFromCSV(r, props, tsvOpts...)
}

func (t table) Render() (template.HTML, error) {
	headerCells := make([]Node, 0, len(t.Headers))
	for _, header := range t.Headers {
		headerCells = append(headerCells, element{Tag: "th", Text: header})
	}

	bodyRows := make([]Node, 0, len(t.Rows))
	for _, row := range t.Rows {
		cells := make([]Node, 0, len(row.Cells))
		for _, cell := range row.Cells {
			cells = append(cells, element{Tag: "td", Children: []Node{cell}})
		}
		bodyRows = append(bodyRows, element{Tag: "tr", Children: cells})
	}

	return element{
		Tag:   "table",
		Attrs: map[string]string{"class": "table"},
		Children: []Node{
			element{
				Tag: "thead",
				Children: []Node{
					element{Tag: "tr", Children: headerCells},
				},
			},
			element{Tag: "tbody", Children: bodyRows},
		},
	}.Render()
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

func HiddenInput(name, value string) Node {
	return element{
		Tag: "input",
		Attrs: map[string]string{
			"name":  name,
			"type":  "hidden",
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

func HTMXButton(label string) *button {
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

func FlashAlerts(flashes []FlashMessage) Node {
	if len(flashes) == 0 {
		return DivProps(ElementProps{ID: "flash-alerts", Class: "hidden"})
	}

	children := make([]Node, 0, len(flashes))
	for _, flash := range flashes {
		children = append(children, DivClass("alert "+flashLevelClass(flash.Level), Text(flash.Message)))
	}

	return DivProps(ElementProps{ID: "flash-alerts", Class: "space-y-2"}, children...)
}

func flashLevelClass(level FlashLevel) string {
	switch level {
	case FlashSuccess:
		return "alert-success"
	case FlashError:
		return "alert-error"
	case FlashWarn:
		return "alert-warning"
	default:
		return "alert-info"
	}
}
