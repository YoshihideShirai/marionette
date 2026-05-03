package frontend

import (
	"html/template"

	shared "github.com/YoshihideShirai/marionette/frontend/shared"
)

// このファイルは基本的なコンポーネントProps/DTO型を定義する。
// 新しい汎用UIコンポーネントの型はここに追加する。

// ComponentProps defines shared style knobs for template components.
type ComponentProps = shared.ComponentProps

type LinkProps struct {
	Label     string
	Icon      string
	Href      string
	Target    string
	Rel       string
	External  bool
	Download  bool
	Filename  string
	AriaLabel string
	Props     ComponentProps
}

type IconButtonProps struct {
	Label        string
	IconSVG      template.HTML
	IconPosition string
	Type         string
	Props        ComponentProps
}

type LoginButtonProps struct {
	Label   string
	IconSVG template.HTML
	Type    string
	Props   ComponentProps
}

type SelectOption = shared.SelectOption

type ModalProps = shared.ModalProps

type FormFieldProps = shared.FormFieldProps

type FormProps = shared.FormProps

type ActionFormProps = shared.ActionFormProps

type InputOptions = shared.InputOptions

type TextareaOptions = shared.TextareaOptions

type EmptyStateProps = shared.EmptyStateProps

type AlertProps struct {
	Title       string
	Description string
	Icon        string
	Props       ComponentProps
}

type ToastProps struct {
	Title       string
	Description string
	Icon        string
	Props       ComponentProps
	Live        string
}

type SkeletonProps struct {
	Rows  int
	Props ComponentProps
}

type ProgressProps struct {
	Value         float64
	Max           float64
	Label         string
	AriaLabel     string
	ShowValue     bool
	Indeterminate bool
	Props         ComponentProps
}
