package daisyui

import frontend "github.com/YoshihideShirai/marionette/frontend"

func Button(label string, props frontend.ComponentProps) frontend.Node {
	return frontend.Button(label, props)
}

func Alert(title, description string, props frontend.ComponentProps) frontend.Node {
	return frontend.UIAlert(frontend.AlertProps{Title: title, Description: description, Props: props})
}

func Card(title, description string, actions frontend.Node, children []frontend.Node, props frontend.ComponentProps) frontend.Node {
	return frontend.UICard(frontend.CardProps{Title: title, Description: description, Actions: actions, Props: props}, children...)
}

func Input(name, value string, props frontend.ComponentProps) frontend.Node {
	return frontend.Input(name, value, props)
}

func Toast(title, description string, props frontend.ComponentProps) frontend.Node {
	return frontend.UIToast(frontend.ToastProps{Title: title, Description: description, Props: props})
}

func Modal(props frontend.ModalProps) frontend.Node {
	return frontend.UIModal(props)
}

func Select(name string, options []frontend.SelectOption, props frontend.ComponentProps) frontend.Node {
	return frontend.UISelect(name, options, props)
}

func Tabs(props frontend.TabsProps) frontend.Node {
	return frontend.UITabs(props)
}

func Badge(props frontend.BadgeProps) frontend.Node {
	return frontend.UIBadge(props)
}

func Skeleton(rows int, props frontend.ComponentProps) frontend.Node {
	return frontend.UISkeleton(frontend.SkeletonProps{Rows: rows, Props: props})
}

func Progress(value, max float64, label string, props frontend.ComponentProps) frontend.Node {
	return frontend.UIProgress(frontend.ProgressProps{Value: value, Max: max, Label: label, Props: props})
}
