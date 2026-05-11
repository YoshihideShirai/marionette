package daisyui

import (
	"strconv"
	"strings"

	lowhtml "github.com/YoshihideShirai/marionette/frontend/html"
	shared "github.com/YoshihideShirai/marionette/frontend/shared"
)

func H1(children ...shared.Node) shared.Node {
	return node("h1", map[string]string{"class": "text-4xl font-bold"}, children...)
}
func H2(children ...shared.Node) shared.Node {
	return node("h2", map[string]string{"class": "text-3xl font-bold"}, children...)
}
func H3(children ...shared.Node) shared.Node {
	return node("h3", map[string]string{"class": "text-2xl font-semibold"}, children...)
}
func H4(children ...shared.Node) shared.Node {
	return node("h4", map[string]string{"class": "text-xl font-semibold"}, children...)
}
func TextNode(text string) shared.Node { return textNode("span", nil, text) }

func PrimaryButton(label string, props shared.ComponentProps) shared.Node {
	if props.Variant == "" {
		props.Variant = "primary"
	}
	props.Class = strings.TrimSpace("btn-primary " + props.Class)
	return Button(label, props)
}

func SecondaryButton(label string, props shared.ComponentProps) shared.Node {
	if props.Variant == "" {
		props.Variant = "secondary"
	}
	props.Class = strings.TrimSpace("btn-secondary " + props.Class)
	return Button(label, props)
}

func GhostButton(label string, props shared.ComponentProps) shared.Node {
	if props.Variant == "" {
		props.Variant = "ghost"
	}
	props.Class = strings.TrimSpace("btn-ghost " + props.Class)
	return Button(label, props)
}

// Avatar follows daisyUI's avatar markup: .avatar > .w-*/mask wrapper > img
func Avatar(src, alt, class string) shared.Node {
	return lowhtml.ElementNode{
		Tag:   "div",
		Attrs: map[string]string{"class": "avatar"},
		Children: []shared.Node{
			lowhtml.ElementNode{
				Tag:      "div",
				Attrs:    map[string]string{"class": class},
				Children: []shared.Node{lowhtml.ElementNode{Tag: "img", Attrs: map[string]string{"src": src, "alt": alt}}},
			},
		},
	}
}

func Navbar(start, center, end shared.Node) shared.Node {
	return node("div", map[string]string{"class": "navbar bg-base-100 shadow-sm"},
		node("div", map[string]string{"class": "navbar-start"}, start),
		node("div", map[string]string{"class": "navbar-center"}, center),
		node("div", map[string]string{"class": "navbar-end"}, end),
	)
}

func Hero(title, description string, actions ...shared.Node) shared.Node {
	children := []shared.Node{
		textNode("h1", map[string]string{"class": "text-5xl font-bold"}, title),
		textNode("p", map[string]string{"class": "py-6"}, description),
	}
	if len(actions) > 0 {
		children = append(children, node("div", map[string]string{"class": "flex gap-2"}, actions...))
	}
	return node("div", map[string]string{"class": "hero bg-base-200 rounded-box"},
		node("div", map[string]string{"class": "hero-content text-center"},
			node("div", map[string]string{"class": "max-w-md"}, children...),
		),
	)
}

func Menu(items ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": "menu bg-base-200 rounded-box"}, Children: listItemChildren("", items)}
}

func Footer(children ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "footer", Attrs: map[string]string{"class": "footer sm:footer-horizontal bg-base-200 text-base-content p-10"}, Children: children}
}

func Drawer(id string, side, content shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "drawer"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "input", Attrs: map[string]string{"id": id, "type": "checkbox", "class": "drawer-toggle"}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "drawer-content"}, Children: []shared.Node{content}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "drawer-side"}, Children: []shared.Node{
			lowhtml.ElementNode{Tag: "label", Attrs: map[string]string{"for": id, "aria-label": "close sidebar", "class": "drawer-overlay"}},
			side,
		}},
	}}
}

func Stat(title, value, desc string) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "stats shadow"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "stat"}, Children: []shared.Node{
			lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "stat-title"}, Text: title},
			lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "stat-value"}, Text: value},
			lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "stat-desc"}, Text: desc},
		}},
	}}
}

func Steps(items ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": "steps steps-vertical lg:steps-horizontal"}, Children: items}
}

func Step(label string, active bool) shared.Node {
	className := "step"
	if active {
		className = "step step-primary"
	}
	return lowhtml.ElementNode{Tag: "li", Attrs: map[string]string{"class": className}, Text: label}
}

func Timeline(items ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": "timeline timeline-vertical"}, Children: items}
}

func TimelineItem(startLabel, endLabel string, content shared.Node) shared.Node {
	children := []shared.Node{}
	if startLabel != "" {
		children = append(children, lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "timeline-start"}, Text: startLabel})
	}
	children = append(children,
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "timeline-middle"}, Text: "●"},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "timeline-end timeline-box"}, Children: []shared.Node{content}},
	)
	if endLabel != "" {
		children = append(children, lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "timeline-end"}, Text: endLabel})
	}
	return lowhtml.ElementNode{Tag: "li", Children: children}
}

func Collapse(title string, content shared.Node, open bool) shared.Node {
	inputAttrs := map[string]string{"type": "checkbox"}
	if open {
		inputAttrs["checked"] = "checked"
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "collapse collapse-arrow bg-base-100 border border-base-300"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "input", Attrs: inputAttrs},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "collapse-title font-semibold"}, Text: title},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "collapse-content text-sm"}, Children: []shared.Node{content}},
	}}
}

func MockupWindow(title string, content shared.Node) shared.Node {
	_ = title
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-window"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "px-4 py-16 bg-base-200"}, Children: []shared.Node{content}},
	}}
}

func Kbd(text string) shared.Node {
	return lowhtml.ElementNode{Tag: "kbd", Attrs: map[string]string{"class": "kbd"}, Text: text}
}

func Code(text string) shared.Node {
	return lowhtml.ElementNode{Tag: "code", Attrs: map[string]string{"class": "bg-base-200 rounded px-1 py-0.5"}, Text: text}
}

func Indicator(item, target shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "indicator"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": "indicator-item"}, Children: []shared.Node{item}},
		target,
	}}
}

func Link(label, href string, props shared.ComponentProps) shared.Node {
	className := "link"
	if props.Class != "" {
		className += " " + props.Class
	}
	return lowhtml.ElementNode{Tag: "a", Attrs: map[string]string{"class": className, "href": href}, Text: label}
}

func Dropdown(trigger, menu shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "dropdown"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"tabindex": "0", "role": "button"}, Children: []shared.Node{trigger}},
		dropdownContent(menu),
	}}
}

func dropdownContent(content shared.Node) shared.Node {
	if n, ok := content.(lowhtml.ElementNode); ok {
		n.Attrs = cloneAttrs(n.Attrs)
		n.Attrs["class"] = appendClass(n.Attrs["class"], "dropdown-content")
		n.Attrs["tabindex"] = "-1"
		return n
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"tabindex": "-1", "class": "dropdown-content"}, Children: []shared.Node{content}}
}

func Tooltip(text string, child shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "tooltip", "data-tip": text}, Children: []shared.Node{child}}
}

func Loading(sizeClass string) shared.Node {
	className := "loading loading-spinner"
	if sizeClass != "" {
		className += " " + sizeClass
	}
	return lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": className}}
}

func RadialProgress(value int, sizeClass string) shared.Node {
	valueText := strconv.Itoa(value)
	attrs := map[string]string{"class": strings.TrimSpace("radial-progress " + sizeClass), "style": "--value:" + valueText + ";", "role": "progressbar", "aria-valuenow": valueText}
	return lowhtml.ElementNode{Tag: "div", Attrs: attrs, Text: valueText + "%"}
}

func Rating(name string, max int, checked int) shared.Node {
	stars := make([]shared.Node, 0, max)
	for i := 1; i <= max; i++ {
		attrs := map[string]string{"type": "radio", "name": name, "class": "mask mask-star-2 bg-orange-400", "value": strconv.Itoa(i)}
		if i == checked {
			attrs["checked"] = "checked"
		}
		stars = append(stars, lowhtml.ElementNode{Tag: "input", Attrs: attrs})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "rating"}, Children: stars}
}

func Range(name string, value int, min int, max int) shared.Node {
	return lowhtml.ElementNode{Tag: "input", Attrs: map[string]string{"type": "range", "name": name, "value": strconv.Itoa(value), "min": strconv.Itoa(min), "max": strconv.Itoa(max), "class": "range"}}
}

func Toggle(name string, checked bool) shared.Node {
	attrs := map[string]string{"type": "checkbox", "name": name, "class": "toggle"}
	if checked {
		attrs["checked"] = "checked"
	}
	return lowhtml.ElementNode{Tag: "input", Attrs: attrs}
}

func ToggleVariant(name string, checked bool, variant string) shared.Node {
	className := "toggle"
	switch strings.ToLower(strings.TrimSpace(variant)) {
	case "primary":
		className += " toggle-primary"
	case "secondary":
		className += " toggle-secondary"
	case "accent":
		className += " toggle-accent"
	case "neutral":
		className += " toggle-neutral"
	case "info":
		className += " toggle-info"
	case "success":
		className += " toggle-success"
	case "warning":
		className += " toggle-warning"
	case "danger", "error":
		className += " toggle-error"
	}
	attrs := map[string]string{"type": "checkbox", "name": name, "class": className}
	if checked {
		attrs["checked"] = "checked"
	}
	return lowhtml.ElementNode{Tag: "input", Attrs: attrs}
}

func ToggleWithIcons(name string, checked bool, className string) shared.Node {
	inputAttrs := map[string]string{"type": "checkbox", "name": name}
	if checked {
		inputAttrs["checked"] = "checked"
	}
	return lowhtml.ElementNode{Tag: "label", Attrs: map[string]string{"class": strings.TrimSpace("toggle text-base-content " + className)}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "input", Attrs: inputAttrs},
		lowhtml.ElementNode{Tag: "svg", Attrs: map[string]string{"aria-label": "enabled", "xmlns": "http://www.w3.org/2000/svg", "viewBox": "0 0 24 24"}, Children: []shared.Node{
			lowhtml.ElementNode{Tag: "g", Attrs: map[string]string{"stroke-linejoin": "round", "stroke-linecap": "round", "stroke-width": "4", "fill": "none", "stroke": "currentColor"}, Children: []shared.Node{
				lowhtml.ElementNode{Tag: "path", Attrs: map[string]string{"d": "M20 6 9 17l-5-5"}},
			}},
		}},
		lowhtml.ElementNode{Tag: "svg", Attrs: map[string]string{"aria-label": "disabled", "xmlns": "http://www.w3.org/2000/svg", "viewBox": "0 0 24 24", "fill": "none", "stroke": "currentColor", "stroke-width": "4", "stroke-linecap": "round", "stroke-linejoin": "round"}, Children: []shared.Node{
			lowhtml.ElementNode{Tag: "path", Attrs: map[string]string{"d": "M18 6 6 18"}},
			lowhtml.ElementNode{Tag: "path", Attrs: map[string]string{"d": "m6 6 12 12"}},
		}},
	}}
}

func Join(children ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "join"}, Children: children}
}

func Mask(shapeClass string, child shared.Node) shared.Node {
	className := "mask " + shapeClass
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": className}, Children: []shared.Node{child}}
}

func Carousel(items ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "carousel w-full"}, Children: items}
}

func CarouselItem(id string, child shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"id": id, "class": "carousel-item w-full"}, Children: []shared.Node{child}}
}

func ChatBubble(content shared.Node, end bool) shared.Node {
	position := "chat-start"
	if end {
		position = "chat-end"
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "chat " + position}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "chat-bubble"}, Children: []shared.Node{content}},
	}}
}

func Countdown(value int) shared.Node {
	valueText := strconv.Itoa(value)
	return lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": "countdown font-mono text-2xl", "aria-live": "polite", "aria-label": valueText}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"style": "--value:" + valueText + ";"}, Text: valueText},
	}}
}

func Status(colorClass string) shared.Node {
	return lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": "status " + colorClass}}
}

func Dock(items ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "dock"}, Children: items}
}

func Fieldset(legend string, fields ...shared.Node) shared.Node {
	children := []shared.Node{lowhtml.ElementNode{Tag: "legend", Attrs: map[string]string{"class": "fieldset-legend"}, Text: legend}}
	children = append(children, fields...)
	return lowhtml.ElementNode{Tag: "fieldset", Attrs: map[string]string{"class": "fieldset"}, Children: children}
}

func Label(text string) shared.Node {
	return lowhtml.ElementNode{Tag: "label", Attrs: map[string]string{"class": "label"}, Text: text}
}

func Validator(message string) shared.Node {
	return lowhtml.ElementNode{Tag: "p", Attrs: map[string]string{"class": "validator-hint"}, Text: message}
}

func BrowserMockup(content shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-browser"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-browser-toolbar"}, Children: []shared.Node{lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "input"}, Text: "https://example.com"}}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "grid place-content-center h-80"}, Children: []shared.Node{content}},
	}}
}

func PhoneMockup(content shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-phone"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-phone-camera"}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-phone-display"}, Children: []shared.Node{content}},
	}}
}

func CodeMockup(lines ...string) shared.Node {
	children := make([]shared.Node, 0, len(lines))
	for _, line := range lines {
		children = append(children, lowhtml.ElementNode{Tag: "pre", Attrs: map[string]string{"data-prefix": "$"}, Children: []shared.Node{lowhtml.ElementNode{Tag: "code", Text: line}}})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-code"}, Children: children}
}

func Calendar(content shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "card bg-base-100 border border-base-300"}, Children: []shared.Node{content}}
}

func Filter(items ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "filter"}, Children: items}
}

func Diff(before, after shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "figure", Attrs: map[string]string{"class": "diff aspect-16/9"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "diff-item-1"}, Children: []shared.Node{before}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "diff-item-2"}, Children: []shared.Node{after}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "diff-resizer"}},
	}}
}

func listItemChildren(itemClass string, items []shared.Node) []shared.Node {
	children := make([]shared.Node, 0, len(items))
	for _, item := range items {
		if n, ok := item.(lowhtml.ElementNode); ok && n.Tag == "li" {
			if itemClass != "" {
				n.Attrs = cloneAttrs(n.Attrs)
				n.Attrs["class"] = appendClass(n.Attrs["class"], itemClass)
			}
			children = append(children, n)
			continue
		}
		attrs := map[string]string{}
		if itemClass != "" {
			attrs["class"] = itemClass
		}
		children = append(children, lowhtml.ElementNode{Tag: "li", Attrs: attrs, Children: []shared.Node{item}})
	}
	return children
}

func List(items ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": "list bg-base-100 rounded-box shadow-md"}, Children: listItemChildren("list-row", items)}
}

func Table(headers []string, rows ...[]shared.Node) shared.Node {
	headersNode := make([]shared.Node, 0, len(headers))
	for _, h := range headers {
		headersNode = append(headersNode, lowhtml.ElementNode{Tag: "th", Text: h})
	}
	tbodyRows := make([]shared.Node, 0, len(rows))
	for _, row := range rows {
		cells := make([]shared.Node, 0, len(row))
		for _, cell := range row {
			cells = append(cells, lowhtml.ElementNode{Tag: "td", Children: []shared.Node{cell}})
		}
		tbodyRows = append(tbodyRows, lowhtml.ElementNode{Tag: "tr", Children: cells})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "overflow-x-auto"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "table", Attrs: map[string]string{"class": "table"}, Children: []shared.Node{
			lowhtml.ElementNode{Tag: "thead", Children: []shared.Node{lowhtml.ElementNode{Tag: "tr", Children: headersNode}}},
			lowhtml.ElementNode{Tag: "tbody", Children: tbodyRows},
		}},
	}}
}

func TextRotate(words []string, animationClass string) shared.Node {
	className := strings.TrimSpace("text-rotate " + animationClass)
	items := make([]shared.Node, 0, len(words))
	for _, w := range words {
		items = append(items, lowhtml.ElementNode{Tag: "span", Text: w})
	}
	return lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": className}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "span", Children: items},
	}}
}

func Hover3DCard(content shared.Node) shared.Node {
	children := make([]shared.Node, 0, 9)
	children = append(children, content)
	for i := 0; i < 8; i++ {
		children = append(children, lowhtml.ElementNode{Tag: "div"})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "hover-3d"}, Children: children}
}

func HoverGallery(items ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "figure", Attrs: map[string]string{"class": "hover-gallery"}, Children: items}
}

func Accordion(title string, content shared.Node, open bool) shared.Node {
	inputAttrs := map[string]string{"type": "radio", "name": "accordion"}
	if open {
		inputAttrs["checked"] = "checked"
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "collapse collapse-arrow bg-base-100 border border-base-300"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "input", Attrs: inputAttrs},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "collapse-title font-semibold"}, Text: title},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "collapse-content text-sm"}, Children: []shared.Node{content}},
	}}
}

func FAB(icon shared.Node, label string) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "fab"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "button", Attrs: map[string]string{"class": "btn btn-lg btn-circle btn-primary", "aria-label": label}, Children: []shared.Node{icon}},
	}}
}

func SpeedDial(trigger shared.Node, items ...shared.Node) shared.Node {
	children := make([]shared.Node, 0, len(items)+1)
	children = append(children, lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"tabindex": "0", "role": "button", "class": "btn btn-lg btn-circle btn-primary"}, Children: []shared.Node{trigger}})
	children = append(children, items...)
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "fab"}, Children: children}
}

func Swap(onNode, offNode shared.Node, active bool) shared.Node {
	attrs := map[string]string{"class": "swap"}
	if active {
		attrs["class"] = "swap swap-active"
	}
	return lowhtml.ElementNode{Tag: "label", Attrs: attrs, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "input", Attrs: map[string]string{"type": "checkbox"}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "swap-on"}, Children: []shared.Node{onNode}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "swap-off"}, Children: []shared.Node{offNode}},
	}}
}

func ThemeController(options ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Children: options}
}

// ThemeControllerOption renders a daisyUI theme-controller radio input.
// See: https://daisyui.com/components/theme-controller/
func ThemeControllerOption(theme string, checked bool, className string) shared.Node {
	attrs := map[string]string{
		"type":       "radio",
		"name":       "theme-buttons",
		"class":      strings.TrimSpace("theme-controller " + className),
		"value":      theme,
		"aria-label": theme,
	}
	if checked {
		attrs["checked"] = "checked"
	}
	return lowhtml.ElementNode{Tag: "input", Attrs: attrs}
}

func DockItem(child shared.Node, active bool) shared.Node {
	attrs := map[string]string{}
	if active {
		attrs["class"] = "dock-active"
	}
	return lowhtml.ElementNode{Tag: "button", Attrs: attrs, Children: []shared.Node{child}}
}

func FilterItem(label string, active bool) shared.Node {
	attrs := map[string]string{"class": "btn btn-sm", "type": "radio", "name": "filter", "aria-label": label}
	if active {
		attrs["checked"] = "checked"
	}
	return lowhtml.ElementNode{Tag: "input", Attrs: attrs}
}

func CalendarGrid(days ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "grid grid-cols-7 gap-1"}, Children: days}
}

func ButtonWithVariants(label string, variant string, size string, style string, props shared.ComponentProps) shared.Node {
	classes := []string{}
	if variant != "" {
		classes = append(classes, "btn-"+variant)
	}
	if size != "" {
		classes = append(classes, "btn-"+size)
	}
	if style != "" {
		classes = append(classes, "btn-"+style)
	}
	props.Class = strings.TrimSpace(strings.Join(append(classes, props.Class), " "))
	return Button(label, props)
}

func InputWithVariants(name, value, color, size, style string, props shared.ComponentProps) shared.Node {
	classes := []string{}
	if color != "" {
		classes = append(classes, "input-"+color)
	}
	if sizeClass := daisySizeClass("input", size); sizeClass != "" {
		classes = append(classes, sizeClass)
	}
	if style != "" && style != "bordered" {
		classes = append(classes, "input-"+style)
	}
	props.Class = strings.TrimSpace(strings.Join(append(classes, props.Class), " "))
	return Input(name, value, props)
}

func SelectWithVariants(name string, options []shared.SelectOption, color, size, style string, props shared.ComponentProps) shared.Node {
	classes := []string{}
	if color != "" {
		classes = append(classes, "select-"+color)
	}
	if sizeClass := daisySizeClass("select", size); sizeClass != "" {
		classes = append(classes, sizeClass)
	}
	if style != "" && style != "bordered" {
		classes = append(classes, "select-"+style)
	}
	props.Class = strings.TrimSpace(strings.Join(append(classes, props.Class), " "))
	return Select(name, options, props)
}

func TextareaWithVariants(name, value, color, size, style string, options shared.TextareaOptions) shared.Node {
	classes := []string{}
	if color != "" {
		classes = append(classes, "textarea-"+color)
	}
	if sizeClass := daisySizeClass("textarea", size); sizeClass != "" {
		classes = append(classes, sizeClass)
	}
	if style != "" && style != "bordered" {
		classes = append(classes, "textarea-"+style)
	}
	options.Props.Class = strings.TrimSpace(strings.Join(append(classes, options.Props.Class), " "))
	return Textarea(name, value, options)
}

func ProgressWithVariant(value, max float64, label, color string, props shared.ComponentProps) shared.Node {
	if color != "" {
		props.Class = strings.TrimSpace("progress-" + color + " " + props.Class)
	}
	return Progress(shared.ProgressProps{Value: value, Max: max, Label: label, Props: props})
}

func BadgeWithVariant(label, color, size, style string, props shared.ComponentProps) shared.Node {
	classes := []string{"badge"}
	if color != "" {
		classes = append(classes, "badge-"+color)
	}
	if size != "" {
		classes = append(classes, "badge-"+size)
	}
	if style != "" {
		classes = append(classes, "badge-"+style)
	}
	props.Class = strings.TrimSpace(strings.Join(append(classes, props.Class), " "))
	return Badge(shared.BadgeProps{Label: label, Props: props})
}

func CheckboxWithVariants(name, value, label, color, size string, checked bool, props shared.ComponentProps) shared.Node {
	classes := []string{}
	if color != "" {
		classes = append(classes, "checkbox-"+color)
	}
	if sizeClass := daisySizeClass("checkbox", size); sizeClass != "" {
		classes = append(classes, sizeClass)
	}
	props.Class = strings.TrimSpace(strings.Join(append(classes, props.Class), " "))
	return Checkbox(shared.CheckboxComponentProps{Name: name, Value: value, Label: label, Checked: checked, Props: props})
}

func RadioGroupWithVariants(name, color, size string, items []shared.RadioItem, props shared.ComponentProps) shared.Node {
	inputClass := "radio"
	if color != "" {
		inputClass += " radio-" + color
	}
	if sizeClass := daisySizeClass("radio", size); sizeClass != "" {
		inputClass += " " + sizeClass
	}
	children := make([]shared.Node, 0, len(items))
	for _, item := range items {
		attrs := map[string]string{"type": "radio", "name": name, "value": item.Value, "class": inputClass}
		if item.Checked {
			attrs["checked"] = "checked"
		}
		if item.Disabled {
			attrs["disabled"] = "disabled"
		}
		children = append(children,
			lowhtml.ElementNode{Tag: "label", Attrs: map[string]string{"class": "label cursor-pointer gap-2"}, Children: []shared.Node{
				lowhtml.ElementNode{Tag: "input", Attrs: attrs},
				textNode("span", map[string]string{"class": "label"}, item.Label),
			}},
		)
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": strings.TrimSpace("space-y-2 " + props.Class)}, Children: children}
}

func RangeWithVariants(name string, value int, min int, max int, color string, size string) shared.Node {
	classes := []string{"range"}
	if color != "" {
		classes = append(classes, "range-"+color)
	}
	if size != "" {
		classes = append(classes, "range-"+size)
	}
	return lowhtml.ElementNode{Tag: "input", Attrs: map[string]string{"type": "range", "name": name, "value": strconv.Itoa(value), "min": strconv.Itoa(min), "max": strconv.Itoa(max), "class": strings.Join(classes, " ")}}
}

func RatingWithVariants(name string, max int, checked int, size string, half bool, allowClear bool) shared.Node {
	classes := []string{"rating"}
	if size != "" {
		classes = append(classes, "rating-"+size)
	}
	if half {
		classes = append(classes, "rating-half")
	}
	stars := make([]shared.Node, 0, max+1)
	if allowClear {
		stars = append(stars, lowhtml.ElementNode{Tag: "input", Attrs: map[string]string{"type": "radio", "name": name, "class": "rating-hidden", "value": "0"}})
	}
	for i := 1; i <= max; i++ {
		attrs := map[string]string{"type": "radio", "name": name, "class": "mask mask-star-2 bg-orange-400", "value": strconv.Itoa(i)}
		if i == checked {
			attrs["checked"] = "checked"
		}
		stars = append(stars, lowhtml.ElementNode{Tag: "input", Attrs: attrs})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": strings.Join(classes, " ")}, Children: stars}
}

func ToastWithPlacement(children []shared.Node, horizontal, vertical string, className string) shared.Node {
	classes := []string{"toast"}
	if horizontal != "" {
		classes = append(classes, "toast-"+horizontal)
	}
	if vertical != "" {
		classes = append(classes, "toast-"+vertical)
	}
	if className != "" {
		classes = append(classes, className)
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": strings.Join(classes, " ")}, Children: children}
}

func TooltipWithVariants(text string, child shared.Node, placement string, color string, open bool) shared.Node {
	classes := []string{"tooltip"}
	if placement != "" {
		classes = append(classes, "tooltip-"+placement)
	}
	if color != "" {
		classes = append(classes, "tooltip-"+color)
	}
	if open {
		classes = append(classes, "tooltip-open")
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": strings.Join(classes, " "), "data-tip": text}, Children: []shared.Node{child}}
}

func TableWithVariants(headers []string, rows [][]shared.Node, zebra bool, pinRows bool, pinCols bool, size string) shared.Node {
	classes := []string{"table"}
	if zebra {
		classes = append(classes, "table-zebra")
	}
	if pinRows {
		classes = append(classes, "table-pin-rows")
	}
	if pinCols {
		classes = append(classes, "table-pin-cols")
	}
	if size != "" {
		classes = append(classes, "table-"+size)
	}
	headersNode := make([]shared.Node, 0, len(headers))
	for _, h := range headers {
		headersNode = append(headersNode, lowhtml.ElementNode{Tag: "th", Text: h})
	}
	tbodyRows := make([]shared.Node, 0, len(rows))
	for _, row := range rows {
		cells := make([]shared.Node, 0, len(row))
		for _, cell := range row {
			cells = append(cells, lowhtml.ElementNode{Tag: "td", Children: []shared.Node{cell}})
		}
		tbodyRows = append(tbodyRows, lowhtml.ElementNode{Tag: "tr", Children: cells})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "overflow-x-auto"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "table", Attrs: map[string]string{"class": strings.Join(classes, " ")}, Children: []shared.Node{
			lowhtml.ElementNode{Tag: "thead", Children: []shared.Node{lowhtml.ElementNode{Tag: "tr", Children: headersNode}}},
			lowhtml.ElementNode{Tag: "tbody", Children: tbodyRows},
		}},
	}}
}

func ModalWithPlacement(props shared.ModalProps, placement string) shared.Node {
	attrs := map[string]string{"class": "modal"}
	if props.Open {
		attrs["open"] = "open"
	}
	if placement != "" {
		attrs["class"] += " modal-" + placement
	}
	return node("dialog", attrs,
		node("div", map[string]string{"class": "modal-box"},
			textNode("h3", map[string]string{"class": "font-bold text-lg"}, props.Title),
			props.Body,
			node("div", map[string]string{"class": "modal-action"}, props.Actions),
		),
		node("form", map[string]string{"method": "dialog", "class": "modal-backdrop"}, textNode("button", nil, "close")),
	)
}

func TabsWithVariants(items []shared.TabsItem, style string, placement string, size string, className string) shared.Node {
	classes := []string{"tabs"}
	if style != "" {
		classes = append(classes, "tabs-"+style)
	}
	if placement != "" {
		classes = append(classes, "tabs-"+placement)
	}
	if size != "" {
		classes = append(classes, "tabs-"+size)
	}
	if className != "" {
		classes = append(classes, className)
	}
	tabNodes := make([]shared.Node, 0, len(items))
	for _, item := range items {
		attrs := map[string]string{"class": "tab", "role": "tab", "type": "button", "aria-selected": "false"}
		if item.Active {
			attrs["class"] += " tab-active"
			attrs["aria-selected"] = "true"
		}
		if item.Disabled {
			attrs["class"] += " tab-disabled"
			attrs["disabled"] = "disabled"
		}
		tabNodes = append(tabNodes, textNode("button", attrs, item.Label))
	}
	return node("div", map[string]string{"class": strings.Join(classes, " "), "role": "tablist"}, tabNodes...)
}

func StepsWithVariants(items []shared.Node, direction string, color string, className string) shared.Node {
	classes := []string{"steps"}
	if direction != "" {
		classes = append(classes, "steps-"+direction)
	}
	if className != "" {
		classes = append(classes, className)
	}
	if color != "" {
		items = withStepColor(items, "step-"+color)
	}
	return lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": strings.Join(classes, " ")}, Children: items}
}

func withStepColor(items []shared.Node, colorClass string) []shared.Node {
	colored := make([]shared.Node, 0, len(items))
	for _, item := range items {
		switch n := item.(type) {
		case lowhtml.ElementNode:
			if n.Tag == "li" && hasClass(n.Attrs["class"], "step") {
				n.Attrs = cloneAttrs(n.Attrs)
				n.Attrs["class"] = appendClass(n.Attrs["class"], colorClass)
			}
			colored = append(colored, n)
		default:
			colored = append(colored, item)
		}
	}
	return colored
}

func hasClass(className, target string) bool {
	for _, class := range strings.Fields(className) {
		if class == target {
			return true
		}
	}
	return false
}

func cloneAttrs(attrs map[string]string) map[string]string {
	cloned := make(map[string]string, len(attrs)+1)
	for key, value := range attrs {
		cloned[key] = value
	}
	return cloned
}

func appendClass(className string, parts ...string) string {
	classes := strings.Fields(className)
	seen := make(map[string]bool, len(classes)+len(parts))
	for _, class := range classes {
		seen[class] = true
	}
	for _, part := range parts {
		if part != "" && !seen[part] {
			classes = append(classes, part)
			seen[part] = true
		}
	}
	return strings.Join(classes, " ")
}

func TimelineWithDirection(items []shared.Node, direction string, compact bool, snapIcon bool, className string) shared.Node {
	classes := []string{"timeline"}
	if direction != "" {
		classes = append(classes, "timeline-"+direction)
	}
	if compact {
		classes = append(classes, "timeline-compact")
	}
	if snapIcon {
		classes = append(classes, "timeline-snap-icon")
	}
	if className != "" {
		classes = append(classes, className)
	}
	return lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": strings.Join(classes, " ")}, Children: items}
}

func LoadingWithVariants(kind string, size string) shared.Node {
	classes := []string{"loading"}
	if kind == "" {
		kind = "spinner"
	}
	classes = append(classes, "loading-"+kind)
	if size != "" {
		classes = append(classes, "loading-"+size)
	}
	return lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": strings.Join(classes, " ")}}
}

func StatusWithVariants(color string, size string) shared.Node {
	classes := []string{"status"}
	if color != "" {
		classes = append(classes, "status-"+color)
	}
	if size != "" {
		classes = append(classes, "status-"+size)
	}
	return lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": strings.Join(classes, " ")}}
}

func ToggleWithVariants(name string, checked bool, color string, size string) shared.Node {
	className := "toggle"
	if color != "" {
		className += " toggle-" + color
	}
	if size != "" {
		className += " toggle-" + size
	}
	attrs := map[string]string{"type": "checkbox", "name": name, "class": className}
	if checked {
		attrs["checked"] = "checked"
	}
	return lowhtml.ElementNode{Tag: "input", Attrs: attrs}
}

func SwapWithVariants(onNode, offNode shared.Node, active bool, rotate bool, flip bool) shared.Node {
	className := "swap"
	if active {
		className += " swap-active"
	}
	if rotate {
		className += " swap-rotate"
	}
	if flip {
		className += " swap-flip"
	}
	return lowhtml.ElementNode{Tag: "label", Attrs: map[string]string{"class": className}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "input", Attrs: map[string]string{"type": "checkbox"}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "swap-on"}, Children: []shared.Node{onNode}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "swap-off"}, Children: []shared.Node{offNode}},
	}}
}

func JoinWithDirection(direction string, children ...shared.Node) shared.Node {
	className := "join"
	if direction != "" {
		className += " join-" + direction
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": className}, Children: children}
}

func DropdownWithPlacement(trigger, menu shared.Node, placement string) shared.Node {
	className := "dropdown"
	if placement != "" {
		className += " dropdown-" + placement
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": className}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"tabindex": "0", "role": "button"}, Children: []shared.Node{trigger}},
		dropdownContent(menu),
	}}
}
