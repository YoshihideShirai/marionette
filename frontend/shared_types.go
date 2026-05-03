package frontend

import lowhtml "github.com/YoshihideShirai/marionette/frontend/html"

// Shared type definitions extracted for future decoupling between
// frontend and frontend/daisyui packages.
// NOTE: kept in package frontend for backward compatibility.

type SharedNode = lowhtml.Node

type SharedComponentProps struct {
	Class    string
	Variant  string
	Size     string
	Disabled bool
}

type SharedSelectOption struct {
	Label    string
	Value    string
	Selected bool
}

type SharedTabsItem struct {
	Label, Href      string
	Active, Disabled bool
}

type SharedTabsProps struct {
	Items     []SharedTabsItem
	AriaLabel string
	Props     SharedComponentProps
}

