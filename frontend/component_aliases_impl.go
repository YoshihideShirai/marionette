package frontend

import (
	"io"

	rdf "github.com/rocketlaunchr/dataframe-go"
	dataframeimports "github.com/rocketlaunchr/dataframe-go/imports"
)

func ButtonComponent(label string, props ComponentProps) Node { return UIButton(label, props) }
func SubmitButton(label string, props ComponentProps) Node    { return UISubmitButton(label, props) }
func ThemeToggleButton(props ComponentProps) Node             { return UIThemeToggleButton(props) }
func InputComponent(name, value string, props ComponentProps) Node {
	return UIInput(name, value, props)
}
func InputWithOptions(name, value string, options InputOptions) Node {
	return UIInputWithOptions(name, value, options)
}
func TextareaComponent(name, value string, options TextareaOptions) Node {
	return UITextarea(name, value, options)
}
func FormComponent(props FormProps, children ...Node) Node       { return UIForm(props, children...) }
func FormFieldComponent(control Node, props FormFieldProps) Node { return UIFormField(control, props) }
func SelectComponent(name string, options []SelectOption, props ComponentProps) Node {
	return UISelect(name, options, props)
}
func Modal(props ModalProps) Node                                 { return UIModal(props) }
func Toast(props ToastProps) Node                                 { return UIToast(props) }
func Alert(props AlertProps) Node                                 { return UIAlert(props) }
func Skeleton(props SkeletonProps) Node                           { return UISkeleton(props) }
func EmptyState(props EmptyStateProps) Node                       { return UIEmptyState(props) }
func TableComponent(props TableProps) Node                        { return UITable(props) }
func Chart(props ChartProps) Node                                 { return UIChart(props) }
func DataFrameComponent(df *rdf.DataFrame, props TableProps) Node { return UIDataFrame(df, props) }
func DataFrameChart(df *rdf.DataFrame, props DataFrameChartProps) Node {
	return UIDataFrameChart(df, props)
}
func DataFrameFromCSV(r io.ReadSeeker, props TableProps, opts ...dataframeimports.CSVLoadOptions) (Node, error) {
	return UIDataFrameFromCSV(r, props, opts...)
}
func DataFrameFromTSV(r io.ReadSeeker, props TableProps, opts ...dataframeimports.CSVLoadOptions) (Node, error) {
	return UIDataFrameFromTSV(r, props, opts...)
}
func Pagination(props PaginationProps) Node                   { return UIPagination(props) }
func Tabs(props TabsProps) Node                               { return UITabs(props) }
func Breadcrumb(props BreadcrumbProps) Node                   { return UIBreadcrumb(props) }
func CheckboxComponent(props CheckboxComponentProps) Node     { return UICheckbox(props) }
func RadioGroupComponent(props RadioGroupComponentProps) Node { return UIRadioGroup(props) }
func SwitchComponent(props SwitchComponentProps) Node         { return UISwitch(props) }
func Stack(props StackProps, children ...Node) Node           { return UIStack(props, children...) }
func Grid(props GridProps, children ...Node) Node             { return UIGrid(props, children...) }
func Split(props SplitProps) Node                             { return UISplit(props) }
func PageHeader(props PageHeaderProps) Node                   { return UIPageHeader(props) }
func ContainerComponent(props ContainerProps, children ...Node) Node {
	return UIContainer(props, children...)
}
func Region(props RegionProps, children ...Node) Node   { return UIRegion(props, children...) }
func Card(props CardProps, children ...Node) Node       { return UICard(props, children...) }
func Section(props SectionProps, children ...Node) Node { return UISection(props, children...) }
