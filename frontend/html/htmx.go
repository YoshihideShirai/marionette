package html

import (
	"html/template"
	"strings"
)

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
