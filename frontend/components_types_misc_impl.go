package frontend

// このファイルはImage以降の補助コンポーネントProps/DTOを定義する。
// ナビゲーション・レイアウト・装飾系の型はここに配置する。

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

type PaginationProps struct {
	Page, TotalPages   int
	PrevHref, NextHref string
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

type ContainerProps struct {
	MaxWidth, Padding string
	Centered          bool
	Props             ComponentProps
}

type RegionProps struct {
	ID    string
	Props ComponentProps
}

type CardProps struct {
	Title, Description string
	Actions            Node
	Gap                string
	Props              ComponentProps
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
