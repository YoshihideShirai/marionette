package daisyui

import (
	"strconv"
	"strings"

	lowhtml "github.com/YoshihideShirai/marionette/frontend/html"
	shared "github.com/YoshihideShirai/marionette/frontend/shared"
)

type DrawerProps struct {
	ID           string
	Class        string
	ContentClass string
	SideClass    string
	OverlayClass string
	Content      shared.Node
	Side         shared.Node
}

func DrawerWithProps(props DrawerProps) shared.Node {
	drawerID := strings.TrimSpace(props.ID)
	if drawerID == "" {
		drawerID = "drawer"
	}
	overlayClass := strings.TrimSpace(props.OverlayClass)
	if overlayClass == "" {
		overlayClass = "drawer-overlay"
	}
	return node("div", map[string]string{"class": strings.TrimSpace("drawer " + props.Class)},
		node("input", map[string]string{"id": drawerID, "type": "checkbox", "class": "drawer-toggle"}),
		node("div", map[string]string{"class": strings.TrimSpace("drawer-content " + props.ContentClass)}, props.Content),
		node("div", map[string]string{"class": strings.TrimSpace("drawer-side " + props.SideClass)},
			node("label", map[string]string{"for": drawerID, "aria-label": "close sidebar", "class": overlayClass}),
			props.Side,
		),
	)
}

type NavbarProps struct {
	Class string
}

func NavbarWithProps(props NavbarProps, children ...shared.Node) shared.Node {
	return node("div", map[string]string{"class": strings.TrimSpace("navbar " + props.Class)}, children...)
}

type MenuProps struct {
	Class string
}

func MenuWithProps(props MenuProps, items ...shared.Node) shared.Node {
	return node("ul", map[string]string{"class": strings.TrimSpace("menu " + props.Class)}, items...)
}

type MenuLinkProps struct {
	Label  string
	Href   string
	Icon   string
	Active bool
	Class  string
}

func MenuTitle(label string) shared.Node {
	return node("li", map[string]string{"class": "menu-title"}, textNode("span", nil, label))
}

func MenuLink(props MenuLinkProps) shared.Node {
	className := strings.TrimSpace(props.Class)
	if props.Active {
		className = strings.TrimSpace(className + " active")
	}
	children := make([]shared.Node, 0, 2)
	if props.Icon != "" {
		children = append(children, textNode("span", map[string]string{"class": "w-6 text-center"}, props.Icon))
	}
	children = append(children, lowhtml.ElementNode{Tag: "span", Text: props.Label})
	return node("li", nil, node("a", map[string]string{"class": className, "href": props.Href}, children...))
}

type StatsProps struct {
	Class string
}

type StatProps struct {
	Title            string
	Value            string
	Description      string
	Figure           shared.Node
	Class            string
	TitleClass       string
	ValueClass       string
	DescriptionClass string
	FigureClass      string
}

func StatsWithProps(props StatsProps, items ...shared.Node) shared.Node {
	return node("div", map[string]string{"class": strings.TrimSpace("stats " + props.Class)}, items...)
}

func StatItem(props StatProps) shared.Node {
	children := make([]shared.Node, 0, 4)
	if props.Figure != nil {
		children = append(children, node("div", map[string]string{"class": strings.TrimSpace("stat-figure " + props.FigureClass)}, props.Figure))
	}
	children = append(children,
		textNode("div", map[string]string{"class": strings.TrimSpace("stat-title " + props.TitleClass)}, props.Title),
		textNode("div", map[string]string{"class": strings.TrimSpace("stat-value " + props.ValueClass)}, props.Value),
		textNode("div", map[string]string{"class": strings.TrimSpace("stat-desc " + props.DescriptionClass)}, props.Description),
	)
	return node("div", map[string]string{"class": strings.TrimSpace("stat " + props.Class)}, children...)
}

type CardPanelProps struct {
	Title       string
	Description string
	Class       string
	BodyClass   string
	Actions     shared.Node
}

func CardPanel(props CardPanelProps, children ...shared.Node) shared.Node {
	bodyChildren := make([]shared.Node, 0, len(children)+3)
	if props.Title != "" {
		bodyChildren = append(bodyChildren, textNode("h2", map[string]string{"class": "card-title"}, props.Title))
	}
	if props.Description != "" {
		bodyChildren = append(bodyChildren, textNode("p", map[string]string{"class": "text-sm text-base-content/60"}, props.Description))
	}
	bodyChildren = append(bodyChildren, children...)
	if props.Actions != nil {
		bodyChildren = append(bodyChildren, node("div", map[string]string{"class": "card-actions justify-end"}, props.Actions))
	}
	return node("div", map[string]string{"class": strings.TrimSpace("card " + props.Class)},
		node("div", map[string]string{"class": strings.TrimSpace("card-body " + props.BodyClass)}, bodyChildren...),
	)
}

type TableProps struct {
	Headers      []string
	Rows         [][]shared.Node
	Class        string
	WrapperClass string
}

func TableWithProps(props TableProps) shared.Node {
	headers := make([]shared.Node, 0, len(props.Headers))
	for _, header := range props.Headers {
		headers = append(headers, textNode("th", nil, header))
	}
	bodyRows := make([]shared.Node, 0, len(props.Rows))
	for _, row := range props.Rows {
		cells := make([]shared.Node, 0, len(row))
		for _, cell := range row {
			cells = append(cells, node("td", nil, cell))
		}
		bodyRows = append(bodyRows, node("tr", nil, cells...))
	}
	return node("div", map[string]string{"class": strings.TrimSpace("overflow-x-auto " + props.WrapperClass)},
		node("table", map[string]string{"class": strings.TrimSpace("table " + props.Class)},
			node("thead", nil, node("tr", nil, headers...)),
			node("tbody", nil, bodyRows...),
		),
	)
}

func ProgressWithClass(value, max float64, className string) shared.Node {
	return node("progress", map[string]string{
		"class": strings.TrimSpace("progress " + className),
		"value": strconv.FormatFloat(value, 'f', -1, 64),
		"max":   strconv.FormatFloat(max, 'f', -1, 64),
	})
}

func ButtonWithAttrs(label string, props shared.ComponentProps, attrs map[string]string) shared.Node {
	buttonAttrs := map[string]string{}
	for key, value := range attrs {
		buttonAttrs[key] = value
	}
	buttonAttrs["class"] = strings.TrimSpace("btn " + props.Class + " " + buttonAttrs["class"])
	if props.Disabled {
		buttonAttrs["disabled"] = "disabled"
	}
	return textNode("button", buttonAttrs, label)
}

type ActionFormOptions struct {
	Action string
	Target string
	Swap   string
	Method string
	Class  string
}

func ActionFormWithOptions(props ActionFormOptions, children ...shared.Node) shared.Node {
	method := strings.ToLower(strings.TrimSpace(props.Method))
	if method == "" {
		method = "post"
	}
	attrs := map[string]string{"method": method, "class": strings.TrimSpace(props.Class)}
	if props.Action != "" {
		attrs["action"] = props.Action
		attrs["hx-"+method] = props.Action
	}
	if props.Target != "" {
		attrs["hx-target"] = props.Target
	}
	if props.Swap != "" {
		attrs["hx-swap"] = props.Swap
	}
	return node("form", attrs, children...)
}

func ButtonContentWithAttrs(props shared.ComponentProps, attrs map[string]string, children ...shared.Node) shared.Node {
	buttonAttrs := map[string]string{}
	for key, value := range attrs {
		buttonAttrs[key] = value
	}
	buttonAttrs["class"] = strings.TrimSpace("btn " + props.Class + " " + buttonAttrs["class"])
	if props.Disabled {
		buttonAttrs["disabled"] = "disabled"
	}
	return node("button", buttonAttrs, children...)
}

func AvatarPlaceholder(label string, outerClass string, innerClass string) shared.Node {
	return node("div", map[string]string{"class": strings.TrimSpace("avatar placeholder " + outerClass)},
		node("div", map[string]string{"class": innerClass}, textNode("span", nil, label)),
	)
}
