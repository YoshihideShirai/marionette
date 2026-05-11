package frontend

import (
	"strings"

	rdf "github.com/rocketlaunchr/dataframe-go"

	components "github.com/YoshihideShirai/marionette/frontend/components"
	daisy "github.com/YoshihideShirai/marionette/frontend/daisyui"
)

func ThemeToggleButton(props ComponentProps) Node { return daisy.ThemeToggleButton(props) }
func InputWithOptions(name, value string, options InputOptions) Node {
	return daisy.InputWithOptions(name, value, options)
}
func ActionForm(props ActionFormProps, children ...Node) Node {
	return daisy.ActionForm(props, children...)
}
func FormField(control Node, props FormFieldProps) Node { return daisy.FormField(control, props) }
func Modal(props ModalProps) Node                       { return daisy.Modal(props) }
func ModalWithVariants(props ModalProps, variant ModalVariantProps) Node {
	return daisy.ModalWithVariants(props, variant)
}
func Toast(props ToastProps) Node {
	alertProps := AlertContentProps{
		Title:       props.Title,
		Description: props.Description,
		Icon:        textIcon(props.Icon),
		Props:       ComponentProps{Variant: props.Props.Variant},
	}
	return daisy.ToastWithContent(props.Props, daisy.AlertWithContent(alertProps))
}
func ToastWithContent(props ComponentProps, children ...Node) Node {
	return daisy.ToastWithContent(props, children...)
}
func Alert(props AlertProps) Node {
	return daisy.AlertWithContent(AlertContentProps{
		Title:       props.Title,
		Description: props.Description,
		Icon:        textIcon(props.Icon),
		Props:       props.Props,
	})
}
func AlertWithContent(props AlertContentProps, children ...Node) Node {
	return daisy.AlertWithContent(props, children...)
}

func textIcon(icon string) Node {
	if strings.TrimSpace(icon) == "" {
		return nil
	}
	return daisy.TextNode(strings.TrimSpace(icon))
}
func Skeleton(props SkeletonProps) Node { return daisy.Skeleton(props.Rows, props.Props) }
func Progress(props ProgressProps) Node { return daisy.Progress(props) }
func EmptyState(props EmptyStateProps) Node {
	return daisy.EmptyState(props)
}
func DataFrameComponent(df *rdf.DataFrame, props TableProps) Node { return DataFrame(df, props) }
func Badge(props BadgeProps) Node                                 { return daisy.Badge(props) }
func Actions(props ActionsProps, children ...Node) Node           { return daisy.Actions(props, children...) }
func Divider(props DividerProps) Node                             { return daisy.Divider(props) }
func TextComponent(props TextProps) Node                          { return daisy.Text(props) }
func FontIcon(props FontIconProps) Node                           { return daisy.FontIcon(props) }
func HiddenField(name, value string) Node                         { return daisy.HiddenField(name, value) }
func Stack(props StackProps, children ...Node) Node               { return daisy.Stack(props, children...) }
func Grid(props GridProps, children ...Node) Node                 { return components.Grid(props, children...) }
func Split(props SplitProps) Node                                 { return daisy.Split(props) }
func PageHeader(props PageHeaderProps) Node                       { return daisy.PageHeader(props) }
func Region(props RegionProps, children ...Node) Node             { return components.Region(props, children...) }
func Box(props BoxProps, children ...Node) Node                   { return daisy.Box(props, children...) }
func AppShell(props AppShellProps) Node                           { return daisy.AppShell(props) }
func Card(props CardProps, children ...Node) Node {
	return daisy.Card(props.Title, props.Description, props.Actions, children, props.Props)
}
func CardWithVariants(props CardProps, variant CardVariantProps, children ...Node) Node {
	return daisy.CardWithVariants(props, variant, children...)
}
func Section(props SectionProps, children ...Node) Node { return daisy.Section(props, children...) }
