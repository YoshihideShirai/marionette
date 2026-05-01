package marionette

import (
	"io"

	rdf "github.com/rocketlaunchr/dataframe-go"
	dataframeimports "github.com/rocketlaunchr/dataframe-go/imports"
)

func SubmitButton(label string, props ComponentProps) Node    { return UISubmitButton(label, props) }
func InputWithOptions(name, value string, options InputOptions) Node {
	return UIInputWithOptions(name, value, options)
}
func ActionForm(props ActionFormProps, children ...Node) Node {
	return UIActionForm(props, children...)
}
func FormField(control Node, props FormFieldProps) Node           { return UIFormField(control, props) }
func Modal(props ModalProps) Node                                 { return UIModal(props) }
func Toast(props ToastProps) Node                                 { return UIToast(props) }
func Alert(props AlertProps) Node                                 { return UIAlert(props) }
func Skeleton(props SkeletonProps) Node                           { return UISkeleton(props) }
func Progress(props ProgressProps) Node                           { return UIProgress(props) }
func EmptyState(props EmptyStateProps) Node                       { return UIEmptyState(props) }
func Chart(props ChartProps) Node                                 { return UIChart(props) }
func Image(props ImageProps) Node                                 { return UIImage(props) }
func DataFrameComponent(df *rdf.DataFrame, props TableProps) Node { return DataFrame(df, props) }
func DataFrameChart(df *rdf.DataFrame, props DataFrameChartProps) Node {
	return UIDataFrameChart(df, props)
}
func DataFrameFromCSV(r io.ReadSeeker, props TableProps, opts ...dataframeimports.CSVLoadOptions) (Node, error) {
	return UIDataFrameFromCSV(r, props, opts...)
}
func DataFrameFromTSV(r io.ReadSeeker, props TableProps, opts ...dataframeimports.CSVLoadOptions) (Node, error) {
	return UIDataFrameFromTSV(r, props, opts...)
}
func Pagination(props PaginationProps) Node             { return UIPagination(props) }
func Tabs(props TabsProps) Node                         { return UITabs(props) }
func Breadcrumb(props BreadcrumbProps) Node             { return UIBreadcrumb(props) }
func Badge(props BadgeProps) Node                       { return UIBadge(props) }
func Actions(props ActionsProps, children ...Node) Node { return UIActions(props, children...) }
func Divider(props DividerProps) Node                   { return UIDivider(props) }
func TextComponent(props TextProps) Node                { return UIText(props) }
func HiddenField(name, value string) Node               { return UIHiddenField(name, value) }
func Stack(props StackProps, children ...Node) Node     { return UIStack(props, children...) }
func Grid(props GridProps, children ...Node) Node       { return UIGrid(props, children...) }
func Split(props SplitProps) Node                       { return UISplit(props) }
func PageHeader(props PageHeaderProps) Node             { return UIPageHeader(props) }
func Region(props RegionProps, children ...Node) Node   { return UIRegion(props, children...) }
func Box(props BoxProps, children ...Node) Node         { return UIBox(props, children...) }
func AppShell(props AppShellProps) Node                 { return UIAppShell(props) }
func Card(props CardProps, children ...Node) Node       { return UICard(props, children...) }
func Section(props SectionProps, children ...Node) Node { return UISection(props, children...) }
func Markdown(props MarkdownProps) Node                 { return UIMarkdown(props) }
