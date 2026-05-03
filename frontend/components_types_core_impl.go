package frontend

import shared "github.com/YoshihideShirai/marionette/frontend/shared"

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

type SelectOption = shared.SelectOption

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

type EmptyStateProps struct {
	Title       string
	Description string
	Skeleton    bool
	Rows        int
	Icon        string
	Props       ComponentProps
}

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
