package marionette

import (
	"html/template"
	"io"
	"strings"

	mf "github.com/YoshihideShirai/marionette/frontend"
	mh "github.com/YoshihideShirai/marionette/frontend/html"
	shared "github.com/YoshihideShirai/marionette/frontend/shared"
	dataframeimports "github.com/rocketlaunchr/dataframe-go/imports"
)

// Node is a declarative UI element that can render itself as safe HTML.
type Node = shared.Node
type element = mh.ElementNode

type Attrs = shared.Attrs
type ElementProps = mf.ElementProps
type Raw = mf.Raw

type TableRowData = shared.TableRowData

type SidebarItem = shared.SidebarItem

func Text(v string) Node { return mf.Text(v) }
func Element(tag string, props ElementProps, children ...Node) Node {
	return mf.Element(tag, props, children...)
}
func Div(children ...Node) Node                           { return mf.Div(children...) }
func DivID(id string, children ...Node) Node              { return mh.DivID(id, children...) }
func DivClass(className string, children ...Node) Node    { return mh.DivClass(className, children...) }
func DivAttrs(attrs Attrs, children ...Node) Node         { return mh.DivAttrs(attrs, children...) }
func DivProps(props ElementProps, children ...Node) Node  { return mf.DivProps(props, children...) }
func Span(children ...Node) Node                          { return mf.Span(children...) }
func SpanProps(props ElementProps, children ...Node) Node { return mf.SpanProps(props, children...) }
func P(children ...Node) Node                             { return mf.P(children...) }
func AnchorProps(props ElementProps, children ...Node) Node {
	return mf.AnchorProps(props, children...)
}
func AsideProps(props ElementProps, children ...Node) Node { return mf.AsideProps(props, children...) }
func DescriptionListProps(props ElementProps, children ...Node) Node {
	return mf.DescriptionListProps(props, children...)
}
func DescriptionTermProps(props ElementProps, children ...Node) Node {
	return mf.DescriptionTermProps(props, children...)
}
func DescriptionDetailsProps(props ElementProps, children ...Node) Node {
	return mf.DescriptionDetailsProps(props, children...)
}
func PProps(props ElementProps, children ...Node) Node { return mf.PProps(props, children...) }
func LabelElement(children ...Node) Node               { return mf.LabelElement(children...) }
func LabelElementProps(props ElementProps, children ...Node) Node {
	return mf.LabelElementProps(props, children...)
}
func InputElement(props ElementProps) Node                  { return mf.InputElement(props) }
func Ul(children ...Node) Node                              { return mf.Ul(children...) }
func UlProps(props ElementProps, children ...Node) Node     { return mf.UlProps(props, children...) }
func Li(children ...Node) Node                              { return mf.Li(children...) }
func LiProps(props ElementProps, children ...Node) Node     { return mf.LiProps(props, children...) }
func H1(children ...Node) Node                              { return mf.H1(children...) }
func H1Props(props ElementProps, children ...Node) Node     { return mf.H1Props(props, children...) }
func H2(children ...Node) Node                              { return mf.H2(children...) }
func H2Props(props ElementProps, children ...Node) Node     { return mf.H2Props(props, children...) }
func H3(children ...Node) Node                              { return mf.H3(children...) }
func H3Props(props ElementProps, children ...Node) Node     { return mf.H3Props(props, children...) }
func H4(children ...Node) Node                              { return mf.H4(children...) }
func H4Props(props ElementProps, children ...Node) Node     { return mf.H4Props(props, children...) }
func Column(children ...Node) Node                          { return mh.Column(children...) }
func HTMXTable(headers []string, rows ...TableRowData) Node { return mf.HTMXTable(headers, rows...) }
func TableRow(cells ...Node) TableRowData                   { return mf.TableRow(cells...) }

func DataFrameFromCSV(r io.ReadSeeker, props TableProps, opts ...dataframeimports.CSVLoadOptions) (Node, error) {
	return mf.DataFrameFromCSV(r, props, opts...)
}
func DataFrameFromTSV(r io.ReadSeeker, props TableProps, opts ...dataframeimports.CSVLoadOptions) (Node, error) {
	return mf.DataFrameFromTSV(r, props, opts...)
}

func SidebarLink(label, href string) SidebarItem { return mf.SidebarLink(label, href) }

type sidebar struct {
	brand     string
	title     string
	items     []SidebarItem
	noteTitle string
	noteText  string
}

func Sidebar(brand, title string, items ...SidebarItem) *sidebar {
	return &sidebar{brand: brand, title: title, items: items}
}
func (s *sidebar) Note(title, text string) *sidebar {
	s.noteTitle = title
	s.noteText = text
	return s
}
func (s *sidebar) Render() (template.HTML, error) {
	inner := mf.Sidebar(s.brand, s.title, s.items...)
	if strings.TrimSpace(s.noteTitle) != "" || strings.TrimSpace(s.noteText) != "" {
		inner = inner.Note(s.noteTitle, s.noteText)
	}
	return inner.Render()
}

type form struct {
	action string
	target string
	kids   []Node
}

func Form(action string, children ...Node) *form { return &form{action: action, kids: children} }
func (f *form) Target(selector string) *form     { f.target = selector; return f }
func (f *form) Render() (template.HTML, error) {
	inner := mf.Form(f.action, f.kids...)
	if f.target != "" {
		inner = inner.Target(f.target)
	}
	return inner.Render()
}

func Input(name, value string, props ...ComponentProps) Node { return mf.Input(name, value, props...) }
func FileUpload(name string, required bool, props ...ComponentProps) Node {
	return mf.FileUpload(name, required, props...)
}
func HiddenInput(name, value string) Node { return mf.HiddenInput(name, value) }
func Submit(label string) Node            { return mf.Submit(label) }

type button struct {
	label  string
	action string
	target string
}

func HTMXButton(label string) *button                    { return &button{label: label} }
func (b *button) OnClick(action string) *button          { return b.Post(action) }
func (b *button) Post(action string) *button             { b.action = action; return b }
func (b *button) TargetSelector(selector string) *button { return b.Target(selector) }
func (b *button) Target(selector string) *button         { b.target = selector; return b }
func (b *button) Render() (template.HTML, error) {
	inner := mf.HTMXButton(b.label)
	if b.action != "" {
		inner = inner.Post(b.action)
	}
	if b.target != "" {
		inner = inner.Target(b.target)
	}
	return inner.Render()
}

func FlashAlerts(flashes []FlashMessage) Node { return mf.FlashAlerts(flashes) }

func actionPath(action string) string {
	trimmed := strings.TrimSpace(action)
	if strings.HasPrefix(trimmed, "/") {
		return trimmed
	}
	return "/" + strings.TrimLeft(trimmed, "/")
}
