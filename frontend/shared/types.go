package shared

import lowhtml "github.com/YoshihideShirai/marionette/frontend/html"

type Node = lowhtml.Node
type Attrs map[string]string

type ComponentProps struct {
	Class    string
	Variant  string
	Size     string
	Disabled bool
}

type SelectOption struct {
	Label    string
	Value    string
	Selected bool
}

type TabsItem struct {
	Label, Href      string
	Active, Disabled bool
}

type TabsProps struct {
	Items     []TabsItem
	AriaLabel string
	Props     ComponentProps
}

type ModalProps struct {
	Title   string
	Body    Node
	Actions Node
	Open    bool
}

type FormFieldProps struct {
	Label    string
	Required bool
	Hint     string
	Error    string
}

type FormProps struct {
	ID     string
	Class  string
	Method string
	Action string
	Attrs  Attrs
}

type ActionFormProps struct {
	ID     string
	Action string
	Target string
	Swap   string
	Method string
	Props  ComponentProps
}

type InputOptions struct {
	Type        string
	Placeholder string
	Min         string
	Max         string
	Required    bool
	Props       ComponentProps
}

type TextareaOptions struct {
	Placeholder string
	Rows        int
	Required    bool
	Props       ComponentProps
}

type PaginationProps struct {
	Page, TotalPages   int
	PrevHref, NextHref string
}

type BreadcrumbItem struct {
	Label, Href string
	Active      bool
}
type BreadcrumbProps struct {
	Items     []BreadcrumbItem
	AriaLabel string
	Props     ComponentProps
}

type CheckboxComponentProps struct {
	Name, Value, Label string
	Checked            bool
	Props              ComponentProps
}

type RadioItem struct {
	Label, Value      string
	Checked, Disabled bool
}
type RadioGroupComponentProps struct {
	Name      string
	Items     []RadioItem
	AriaLabel string
	Props     ComponentProps
}

type SwitchComponentProps struct {
	Name, Value, Label string
	Checked            bool
	Props              ComponentProps
}

type BadgeProps struct {
	Label string
	Props ComponentProps
}
type ActionsProps struct {
	Align, Gap string
	Wrap       bool
	Props      ComponentProps
}
type DividerProps struct {
	Spacing string
	Props   ComponentProps
}
type TextProps struct {
	Text, Size, Weight, Tone string
	Props                    ComponentProps
}
type FontIconProps struct {
	Name, Library, AriaLabel string
	Decorative               bool
	Props                    ComponentProps
}

type StackProps struct {
	Direction, Gap, Align, Justify string
	Wrap                           bool
	Props                          ComponentProps
}
type GridProps struct {
	Columns, Gap, MinColumnWidth string
	Props                        ComponentProps
}
type SplitProps struct {
	Main, Aside     Node
	AsideWidth      string
	ReverseOnMobile bool
	Gap             string
	Props           ComponentProps
}
type PageHeaderProps struct {
	Title, Description string
	Actions            Node
	Props              ComponentProps
}
type RegionProps struct {
	ID    string
	Props ComponentProps
}
type SectionProps struct {
	Title, Description string
	Actions            Node
	Props              ComponentProps
}
type BoxProps struct {
	Padding string
	Border  bool
	Tone    string
	Props   ComponentProps
}
type AppShellProps struct {
	ID, MainID                        string
	Sidebar, Flashes, Header, Content Node
	Props                             ComponentProps
}
type ContainerProps struct {
	MaxWidth, Padding string
	Centered          bool
	Props             ComponentProps
}

type EmptyStateProps struct {
	Title       string
	Description string
	Skeleton    bool
	Rows        int
	Icon        string
	Props       ComponentProps
}

type ImageProps struct {
	Src         string
	Alt         string
	Caption     string
	Width       int
	Height      int
	Loading     string
	Decoding    string
	AspectRatio string
	ObjectFit   string
	Props       ComponentProps
}

type SidebarItem struct {
	Label   string
	Href    string
	Current bool
}

type TableRowData struct {
	Cells []Node
}

type ChartType string

const (
	ChartTypeBar      ChartType = "bar"
	ChartTypeLine     ChartType = "line"
	ChartTypePie      ChartType = "pie"
	ChartTypeDoughnut ChartType = "doughnut"
	ChartTypeScatter  ChartType = "scatter"
)

type ChartDataset struct {
	Label           string
	Data            []float64
	Points          []ChartPoint
	BackgroundColor string
	BorderColor     string
	Fill            bool
	Tension         float64
}
type ChartPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type ChartOptions struct {
	BeginAtZero bool
	Stacked     bool
	HideLegend  bool
	AspectRatio float64
	XAxisLabel  string
	YAxisLabel  string
}
type ChartProps struct {
	Type            ChartType
	Title           string
	Description     string
	Labels          []string
	Datasets        []ChartDataset
	Options         ChartOptions
	AriaLabel       string
	Height          int
	Props           ComponentProps
	QueryStateName  string
	QueryStateLabel string
}
