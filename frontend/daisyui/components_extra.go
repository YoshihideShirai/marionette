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
	props.Class = strings.TrimSpace("btn btn-primary " + props.Class)
	return Button(label, props)
}

func SecondaryButton(label string, props shared.ComponentProps) shared.Node {
	if props.Variant == "" {
		props.Variant = "secondary"
	}
	props.Class = strings.TrimSpace("btn btn-secondary " + props.Class)
	return Button(label, props)
}

func GhostButton(label string, props shared.ComponentProps) shared.Node {
	if props.Variant == "" {
		props.Variant = "ghost"
	}
	props.Class = strings.TrimSpace("btn btn-ghost " + props.Class)
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
	return lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": "menu bg-base-200 rounded-box"}, Children: items}
}

func Footer(children ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "footer", Attrs: map[string]string{"class": "footer sm:footer-horizontal bg-base-200 text-base-content p-10"}, Children: children}
}

func Drawer(id string, side, content shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "drawer", "id": id}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "input", Attrs: map[string]string{"id": id + "-toggle", "type": "checkbox", "class": "drawer-toggle"}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "drawer-content"}, Children: []shared.Node{content}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "drawer-side"}, Children: []shared.Node{side}},
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
	className := "collapse collapse-arrow bg-base-100 border border-base-300"
	attrs := map[string]string{"class": className}
	if open {
		attrs["open"] = "open"
	}
	return lowhtml.ElementNode{Tag: "details", Attrs: attrs, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "summary", Attrs: map[string]string{"class": "collapse-title font-semibold"}, Text: title},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "collapse-content text-sm"}, Children: []shared.Node{content}},
	}}
}

func MockupWindow(title string, content shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-window border border-base-300"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "px-4 py-16 bg-base-200"}, Children: []shared.Node{content}},
		lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": "sr-only"}, Text: title},
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
		lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": "indicator-item badge badge-secondary"}, Children: []shared.Node{item}},
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
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"tabindex": "0", "class": "dropdown-content z-1 card card-sm bg-base-100 shadow-md"}, Children: []shared.Node{menu}},
	}}
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
	attrs := map[string]string{"class": "radial-progress " + sizeClass, "style": "--value:" + strconv.Itoa(value) + ";", "role": "progressbar"}
	return lowhtml.ElementNode{Tag: "div", Attrs: attrs, Text: strconv.Itoa(value) + "%"}
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

func Join(children ...shared.Node) shared.Node {
	wrapped := make([]shared.Node, 0, len(children))
	for _, child := range children {
		wrapped = append(wrapped, lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "join-item"}, Children: []shared.Node{child}})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "join"}, Children: wrapped}
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
	return lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": "countdown font-mono text-2xl"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"style": "--value:" + strconv.Itoa(value) + ";"}},
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
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-browser border border-base-300"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-browser-toolbar"}, Children: []shared.Node{lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "input"}, Text: "https://example.com"}}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "grid place-content-center h-80"}, Children: []shared.Node{content}},
	}}
}

func PhoneMockup(content shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-phone border-primary"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-phone-camera"}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-phone-display"}, Children: []shared.Node{content}},
	}}
}

func CodeMockup(lines ...string) shared.Node {
	children := make([]shared.Node, 0, len(lines))
	for _, line := range lines {
		children = append(children, lowhtml.ElementNode{Tag: "pre", Children: []shared.Node{lowhtml.ElementNode{Tag: "code", Text: line}}})
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

func List(items ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": "list bg-base-100 rounded-box shadow-md"}, Children: items}
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
	return lowhtml.ElementNode{Tag: "table", Attrs: map[string]string{"class": "table"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "thead", Children: []shared.Node{lowhtml.ElementNode{Tag: "tr", Children: headersNode}}},
		lowhtml.ElementNode{Tag: "tbody", Children: tbodyRows},
	}}
}

func TextRotate(words []string, animationClass string) shared.Node {
	if animationClass == "" {
		animationClass = "animate-pulse"
	}
	items := make([]shared.Node, 0, len(words))
	for _, w := range words {
		items = append(items, lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": animationClass}, Text: w})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "flex gap-2"}, Children: items}
}

func Hover3DCard(content shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "card bg-base-100 shadow-xl transition-transform duration-300 hover:-translate-y-1 hover:shadow-2xl"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "card-body"}, Children: []shared.Node{content}},
	}}
}

func HoverGallery(items ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "grid grid-cols-2 md:grid-cols-3 gap-4"}, Children: items}
}

func Accordion(title string, content shared.Node, open bool) shared.Node {
	return Collapse(title, content, open)
}

func FAB(icon shared.Node, label string) shared.Node {
	return lowhtml.ElementNode{Tag: "button", Attrs: map[string]string{"class": "btn btn-primary btn-circle fixed bottom-6 right-6"}, Children: []shared.Node{
		icon,
		lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": "sr-only"}, Text: label},
	}}
}

func SpeedDial(trigger shared.Node, items ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "dropdown dropdown-top dropdown-end fixed bottom-6 right-6"}, Children: []shared.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"tabindex": "0", "role": "button"}, Children: []shared.Node{trigger}},
		lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"tabindex": "0", "class": "dropdown-content menu bg-base-100 rounded-box z-1 w-52 p-2 shadow"}, Children: items},
	}}
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
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "join"}, Children: options}
}

// ThemeControllerOption renders a daisyUI theme-controller radio input.
// See: https://daisyui.com/components/theme-controller/
func ThemeControllerOption(theme string, checked bool, className string) shared.Node {
	attrs := map[string]string{
		"type":  "radio",
		"name":  "theme-buttons",
		"class": strings.TrimSpace("theme-controller " + className),
		"value": theme,
	}
	if checked {
		attrs["checked"] = "checked"
	}
	return lowhtml.ElementNode{Tag: "input", Attrs: attrs, Text: theme}
}

func DockItem(child shared.Node, active bool) shared.Node {
	className := "dock-label"
	if active {
		className = "dock-active"
	}
	return lowhtml.ElementNode{Tag: "button", Attrs: map[string]string{"class": className}, Children: []shared.Node{child}}
}

func FilterItem(label string, active bool) shared.Node {
	className := "btn btn-sm"
	if active {
		className += " btn-active"
	}
	return lowhtml.ElementNode{Tag: "button", Attrs: map[string]string{"class": className}, Text: label}
}

func CalendarGrid(days ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "grid grid-cols-7 gap-1"}, Children: days}
}
