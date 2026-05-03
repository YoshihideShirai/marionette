package shared

import lowhtml "github.com/YoshihideShirai/marionette/frontend/html"

type Node = lowhtml.Node

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
