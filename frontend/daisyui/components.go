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
