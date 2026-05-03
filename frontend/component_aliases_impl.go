package frontend

import (
	rdf "github.com/rocketlaunchr/dataframe-go"
	"strconv"

	daisy "github.com/YoshihideShirai/marionette/frontend/daisyui"
	lowhtml "github.com/YoshihideShirai/marionette/frontend/html"
)

func SubmitButton(label string, props ComponentProps) Node { return daisy.SubmitButton(label, props) }
func ThemeToggleButton(props ComponentProps) Node {
	className := "btn btn-ghost"
	if props.Class != "" {
		className += " " + props.Class
	}
	return lowhtml.ElementNode{Tag: "button", Attrs: map[string]string{"type": "button", "class": className, "aria-label": "Toggle theme", "onclick": "window.mrnToggleTheme()"}, Children: []Node{lowhtml.ElementNode{Tag: "span", Text: "🌓 Theme"}}}
}
func InputWithOptions(name, value string, options InputOptions) Node {
	return daisy.InputWithOptions(name, value, options)
}
func ActionForm(props ActionFormProps, children ...Node) Node {
	return daisy.ActionForm(props, children...)
}
func FormField(control Node, props FormFieldProps) Node { return daisy.FormField(control, props) }
func Modal(props ModalProps) Node                       { return daisy.Modal(props) }
func Toast(props ToastProps) Node                       { return daisy.Toast(props.Title, props.Description, props.Props) }
func Alert(props AlertProps) Node                       { return daisy.Alert(props.Title, props.Description, props.Props) }
func Skeleton(props SkeletonProps) Node                 { return daisy.Skeleton(props.Rows, props.Props) }
func Progress(props ProgressProps) Node {
	max := props.Max
	if max <= 0 {
		max = 100
	}
	attrs := map[string]string{"class": progressClass(props.Props), "value": strconv.FormatFloat(props.Value, 'f', -1, 64), "max": strconv.FormatFloat(max, 'f', -1, 64)}
	if props.Indeterminate {
		delete(attrs, "value")
	}
	return lowhtml.ElementNode{Tag: "progress", Attrs: attrs, Children: []Node{lowhtml.ElementNode{Tag: "span", Text: props.Label}}}
}
func EmptyState(props EmptyStateProps) Node {
	if props.Skeleton {
		rows := props.Rows
		if rows <= 0 {
			rows = 3
		}
		children := make([]Node, 0, rows)
		for range rows {
			children = append(children, lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "skeleton h-4 w-full"}})
		}
		return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "space-y-2", "aria-busy": "true", "aria-live": "polite"}, Children: children}
	}
	return daisy.EmptyState(props)
}
func Chart(props ChartProps) Node                                 { return UIChart(props) }
func Image(props ImageProps) Node                                 { return daisy.Image(props) }
func DataFrameComponent(df *rdf.DataFrame, props TableProps) Node { return DataFrame(df, props) }
func Pagination(props PaginationProps) Node                       { return UIPagination(props) }
func Tabs(props TabsProps) Node                                   { return daisy.Tabs(props) }
func Breadcrumb(props BreadcrumbProps) Node                       { return daisy.Breadcrumb(props) }
func Badge(props BadgeProps) Node                                 { return daisy.Badge(props) }
func Actions(props ActionsProps, children ...Node) Node           { return daisy.Actions(props, children...) }
func Divider(props DividerProps) Node                             { return daisy.Divider(props) }
func TextComponent(props TextProps) Node                          { return daisy.Text(props) }
func UIText(props TextProps) Node                                 { return TextComponent(props) }
func FontIcon(props FontIconProps) Node                           { return daisy.FontIcon(props) }
func HiddenField(name, value string) Node                         { return daisy.HiddenField(name, value) }
func Stack(props StackProps, children ...Node) Node               { return daisy.Stack(props, children...) }
func Grid(props GridProps, children ...Node) Node {
	return layoutChildrenNode("components/grid", gridClass(props), children)
}
func Split(props SplitProps) Node           { return daisy.Split(props) }
func PageHeader(props PageHeaderProps) Node { return daisy.PageHeader(props) }
func Region(props RegionProps, children ...Node) Node {
	attrs := map[string]string{"id": props.ID}
	if props.Props.Class != "" {
		attrs["class"] = props.Props.Class
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: attrs, Children: children}
}
func Box(props BoxProps, children ...Node) Node { return daisy.Box(props, children...) }
func AppShell(props AppShellProps) Node         { return daisy.AppShell(props) }
func Card(props CardProps, children ...Node) Node {
	return daisy.Card(props.Title, props.Description, props.Actions, children, props.Props)
}
func Section(props SectionProps, children ...Node) Node { return daisy.Section(props, children...) }
