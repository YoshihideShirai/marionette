package frontend

import shared "github.com/YoshihideShirai/marionette/frontend/shared"

// このファイルはImage以降の補助コンポーネントProps/DTOを定義する。
// ナビゲーション・レイアウト・装飾系の型はここに配置する。

type ImageProps = shared.ImageProps

type PaginationProps = shared.PaginationProps

type TabsItem = shared.TabsItem
type TabsProps = shared.TabsProps

type BreadcrumbItem = shared.BreadcrumbItem
type BreadcrumbProps = shared.BreadcrumbProps

type CheckboxComponentProps = shared.CheckboxComponentProps
type RadioItem = shared.RadioItem
type RadioGroupComponentProps = shared.RadioGroupComponentProps
type SwitchComponentProps = shared.SwitchComponentProps

type BadgeProps = shared.BadgeProps
type ActionsProps = shared.ActionsProps
type DividerProps = shared.DividerProps
type TextProps = shared.TextProps
type FontIconProps = shared.FontIconProps

type StackProps = shared.StackProps
type GridProps = shared.GridProps
type SplitProps = shared.SplitProps
type PageHeaderProps = shared.PageHeaderProps
type ContainerProps = shared.ContainerProps
type RegionProps = shared.RegionProps

type CardProps struct {
	Title, Description string
	Actions            Node
	Gap                string
	Props              ComponentProps
}

type SectionProps = shared.SectionProps
type BoxProps = shared.BoxProps
type AppShellProps = shared.AppShellProps
