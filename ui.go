package marionette

import (
	"io"

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

func Sidebar(brand, title string, items ...SidebarItem) *mh.Sidebar {
	return mh.NewSidebar(brand, title, items...)
}

func Form(action string, children ...Node) *mh.Form { return mh.NewForm(action, children...) }

func Input(name, value string, props ...ComponentProps) Node { return mf.Input(name, value, props...) }
func FileUpload(name string, required bool, props ...ComponentProps) Node {
	return mf.FileUpload(name, required, props...)
}
func HiddenInput(name, value string) Node { return mf.HiddenInput(name, value) }
func Submit(label string) Node            { return mf.Submit(label) }

func HTMXButton(label string) *mh.Button { return mh.HTMXButton(label) }

func FlashAlerts(flashes []FlashMessage) Node { return mf.FlashAlerts(flashes) }

func actionPath(action string) string { return mh.ActionPath(action) }
