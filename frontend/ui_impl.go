package frontend

import (
	"context"
	"fmt"
	"io"

	lowhtml "github.com/YoshihideShirai/marionette/frontend/html"
	shared "github.com/YoshihideShirai/marionette/frontend/shared"
	dataframeimports "github.com/rocketlaunchr/dataframe-go/imports"
)

type FlashLevel string

const (
	FlashSuccess FlashLevel = "success"
	FlashError   FlashLevel = "error"
	FlashInfo    FlashLevel = "info"
	FlashWarn    FlashLevel = "warn"
)

type FlashMessage struct {
	Level   FlashLevel `json:"level"`
	Message string     `json:"message"`
}

// Node is a declarative UI element that can render itself as safe HTML.
type Node = lowhtml.Node
type element = lowhtml.ElementNode

// Attrs defines HTML attributes for low-level element constructors.
type Attrs = lowhtml.Attrs

// ElementProps defines common HTML element attributes while keeping class and
// id easy to scan at call sites.
type ElementProps = lowhtml.ElementProps

// Raw allows trusted HTML snippets (e.g. full page shell).
type Raw = lowhtml.Raw

func textNode(v string) Node {
	return lowhtml.Text(v)
}

func htmlElement(tag string, props ElementProps, children ...Node) Node {
	return lowhtml.Element(tag, props, children...)
}

func Text(v string) Node { return lowhtml.Text(v) }

func Element(tag string, props ElementProps, children ...Node) Node {
	return lowhtml.Element(tag, props, children...)
}

func Div(children ...Node) Node { return lowhtml.Div(children...) }

func DivProps(props ElementProps, children ...Node) Node {
	return lowhtml.DivProps(props, children...)
}

func Span(children ...Node) Node { return lowhtml.Span(children...) }

func SpanProps(props ElementProps, children ...Node) Node {
	return lowhtml.SpanProps(props, children...)
}

func P(children ...Node) Node { return lowhtml.P(children...) }

func AnchorProps(props ElementProps, children ...Node) Node {
	return lowhtml.AnchorProps(props, children...)
}

func AsideProps(props ElementProps, children ...Node) Node {
	return lowhtml.AsideProps(props, children...)
}

func DescriptionListProps(props ElementProps, children ...Node) Node {
	return lowhtml.DescriptionListProps(props, children...)
}

func DescriptionTermProps(props ElementProps, children ...Node) Node {
	return lowhtml.DescriptionTermProps(props, children...)
}

func DescriptionDetailsProps(props ElementProps, children ...Node) Node {
	return lowhtml.DescriptionDetailsProps(props, children...)
}

func PProps(props ElementProps, children ...Node) Node { return lowhtml.PProps(props, children...) }

func LabelElement(children ...Node) Node { return lowhtml.LabelElement(children...) }

func LabelElementProps(props ElementProps, children ...Node) Node {
	return lowhtml.LabelElementProps(props, children...)
}

func InputElement(props ElementProps) Node { return lowhtml.InputElement(props) }

func Ul(children ...Node) Node { return lowhtml.Ul(children...) }

func UlProps(props ElementProps, children ...Node) Node { return lowhtml.UlProps(props, children...) }

func Li(children ...Node) Node { return lowhtml.Li(children...) }

func LiProps(props ElementProps, children ...Node) Node { return lowhtml.LiProps(props, children...) }

func H1(children ...Node) Node { return lowhtml.H1(children...) }

func H1Props(props ElementProps, children ...Node) Node { return lowhtml.H1Props(props, children...) }

func H2(children ...Node) Node { return lowhtml.H2(children...) }

func H2Props(props ElementProps, children ...Node) Node { return lowhtml.H2Props(props, children...) }

func H3(children ...Node) Node { return lowhtml.H3(children...) }

func H3Props(props ElementProps, children ...Node) Node { return lowhtml.H3Props(props, children...) }

func H4(children ...Node) Node { return lowhtml.H4(children...) }

func H4Props(props ElementProps, children ...Node) Node { return lowhtml.H4Props(props, children...) }

type TableRowData = shared.TableRowData

func HTMXTable(headers []string, rows ...TableRowData) Node {
	return lowhtml.HTMXTable(headers, rows...)
}

func TableRow(cells ...Node) TableRowData {
	return lowhtml.TableRow(cells...)
}

func DataFrameFromCSV(r io.ReadSeeker, props TableProps, opts ...dataframeimports.CSVLoadOptions) (Node, error) {
	if r == nil {
		return nil, fmt.Errorf("csv reader is nil")
	}
	df, err := dataframeimports.LoadFromCSV(context.Background(), r, opts...)
	if err != nil {
		return nil, err
	}
	return DataFrame(df, props), nil
}

func DataFrameFromTSV(r io.ReadSeeker, props TableProps, opts ...dataframeimports.CSVLoadOptions) (Node, error) {
	tsvOpts := make([]dataframeimports.CSVLoadOptions, len(opts))
	copy(tsvOpts, opts)
	if len(tsvOpts) == 0 {
		tsvOpts = append(tsvOpts, dataframeimports.CSVLoadOptions{Comma: '\t'})
	} else if tsvOpts[0].Comma == 0 {
		tsvOpts[0].Comma = '\t'
	}
	return DataFrameFromCSV(r, props, tsvOpts...)
}

type SidebarItem = shared.SidebarItem

func Sidebar(brand, title string, items ...SidebarItem) *lowhtml.Sidebar {
	return lowhtml.NewSidebar(brand, title, items...)
}

func SidebarLink(label, href string) SidebarItem {
	return lowhtml.SidebarLink(label, href)
}

func Form(action string, children ...Node) *lowhtml.Form { return lowhtml.NewForm(action, children...) }

func Input(name, value string, props ...ComponentProps) Node {
	if len(props) > 0 {
		return lowhtml.Input(name, value, props[0].Class)
	}
	return lowhtml.Input(name, value)
}

func FileUpload(name string, required bool, props ...ComponentProps) Node {
	if len(props) > 0 {
		return lowhtml.FileUpload(name, required, props[0].Class)
	}
	return lowhtml.FileUpload(name, required)
}

func HiddenInput(name, value string) Node { return lowhtml.HiddenInput(name, value) }

func Submit(label string) Node { return lowhtml.Submit(label) }

func HTMXButton(label string) *lowhtml.Button { return lowhtml.HTMXButton(label) }

func actionPath(action string) string { return lowhtml.ActionPath(action) }

func FlashAlerts(flashes []FlashMessage) Node {
	if len(flashes) == 0 {
		return htmlElement("div", ElementProps{ID: "flash-alerts", Class: "hidden"})
	}

	children := make([]Node, 0, len(flashes))
	for _, flash := range flashes {
		children = append(children, htmlElement("div", ElementProps{Class: "alert " + flashLevelClass(flash.Level)}, textNode(flash.Message)))
	}

	return htmlElement("div", ElementProps{ID: "flash-alerts", Class: "space-y-2"}, children...)
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
