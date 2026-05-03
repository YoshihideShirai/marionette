package daisyui

import (
	"strconv"
	"strings"

	lowhtml "github.com/YoshihideShirai/marionette/frontend/html"
	shared "github.com/YoshihideShirai/marionette/frontend/shared"
)

func node(tag string, attrs map[string]string, children ...shared.Node) shared.Node {
	return lowhtml.ElementNode{Tag: tag, Attrs: attrs, Children: children}
}

func textNode(tag string, attrs map[string]string, text string) shared.Node {
	return lowhtml.ElementNode{Tag: tag, Attrs: attrs, Text: text}
}

func Button(label string, props shared.ComponentProps) shared.Node {
	className := strings.TrimSpace("btn " + props.Class)
	attrs := map[string]string{"class": className}
	if props.Disabled {
		attrs["disabled"] = "disabled"
	}
	return textNode("button", attrs, label)
}

func Alert(title, description string, props shared.ComponentProps) shared.Node {
	return node("div", map[string]string{"class": strings.TrimSpace("alert " + props.Class)},
		textNode("span", nil, strings.TrimSpace(title+" "+description)),
	)
}

func Card(title, description string, actions shared.Node, children []shared.Node, props shared.ComponentProps) shared.Node {
	cardChildren := make([]shared.Node, 0, len(children)+1)
	if title != "" || description != "" || actions != nil {
		headerChildren := []shared.Node{}
		if title != "" {
			headerChildren = append(headerChildren, textNode("h2", map[string]string{"class": "card-title"}, title))
		}
		if description != "" {
			headerChildren = append(headerChildren, textNode("p", nil, description))
		}
		if actions != nil {
			headerChildren = append(headerChildren, node("div", map[string]string{"class": "card-actions justify-end"}, actions))
		}
		cardChildren = append(cardChildren, node("div", map[string]string{"class": "card-body"}, headerChildren...))
	}
	cardChildren = append(cardChildren, children...)
	return node("div", map[string]string{"class": strings.TrimSpace("card bg-base-100 shadow-sm " + props.Class)}, cardChildren...)
}

func Input(name, value string, props shared.ComponentProps) shared.Node {
	attrs := map[string]string{
		"name":  name,
		"value": value,
		"class": strings.TrimSpace("input input-bordered w-full " + props.Class),
	}
	if props.Disabled {
		attrs["disabled"] = "disabled"
	}
	return node("input", attrs)
}

func Toast(title, description string, props shared.ComponentProps) shared.Node {
	return node("div", map[string]string{"class": strings.TrimSpace("toast " + props.Class)},
		node("div", map[string]string{"class": "alert"}, textNode("span", nil, strings.TrimSpace(title+" "+description))),
	)
}

func Modal(props shared.ModalProps) shared.Node {
	className := "modal"
	if props.Open {
		className += " modal-open"
	}
	return node("div", map[string]string{"class": className},
		node("div", map[string]string{"class": "modal-box"},
			textNode("h3", map[string]string{"class": "font-bold text-lg"}, props.Title),
			props.Body,
			node("div", map[string]string{"class": "modal-action"}, props.Actions),
		),
	)
}

func Select(name string, options []shared.SelectOption, props shared.ComponentProps) shared.Node {
	children := make([]shared.Node, 0, len(options))
	for _, opt := range options {
		attrs := map[string]string{"value": opt.Value}
		if opt.Selected {
			attrs["selected"] = "selected"
		}
		children = append(children, textNode("option", attrs, opt.Label))
	}
	return node("select", map[string]string{
		"name":  name,
		"class": strings.TrimSpace("select select-bordered " + props.Class),
	}, children...)
}

func Tabs(props shared.TabsProps) shared.Node {
	tabNodes := make([]shared.Node, 0, len(props.Items))
	for _, item := range props.Items {
		className := "tab"
		if item.Active {
			className += " tab-active"
		}
		tabNodes = append(tabNodes, textNode("a", map[string]string{"class": className, "href": item.Href}, item.Label))
	}
	return node("div", map[string]string{"class": strings.TrimSpace("tabs " + props.Props.Class)}, tabNodes...)
}

func Badge(props shared.BadgeProps) shared.Node {
	return textNode("span", map[string]string{"class": strings.TrimSpace("badge " + props.Props.Class)}, props.Label)
}

func Skeleton(rows int, props shared.ComponentProps) shared.Node {
	if rows <= 0 {
		rows = 3
	}
	items := make([]shared.Node, 0, rows)
	for i := 0; i < rows; i++ {
		items = append(items, node("div", map[string]string{"class": "skeleton h-4 w-full"}))
	}
	return node("div", map[string]string{"class": strings.TrimSpace("space-y-2 " + props.Class)}, items...)
}

func Progress(value, max float64, label string, props shared.ComponentProps) shared.Node {
	return node("progress", map[string]string{"class": strings.TrimSpace("progress w-full " + props.Class), "value": strconv.FormatFloat(value, 'f', -1, 64), "max": strconv.FormatFloat(max, 'f', -1, 64)}, textNode("span", nil, label))
}

func Checkbox(props shared.CheckboxComponentProps) shared.Node {
	inputAttrs := map[string]string{"type": "checkbox", "class": strings.TrimSpace("checkbox " + props.Props.Class), "name": props.Name, "value": props.Value}
	if props.Checked {
		inputAttrs["checked"] = "checked"
	}
	return node("label", map[string]string{"class": "label cursor-pointer gap-2"}, node("input", inputAttrs), textNode("span", map[string]string{"class": "label-text"}, props.Label))
}

func RadioGroup(props shared.RadioGroupComponentProps) shared.Node {
	items := make([]shared.Node, 0, len(props.Items))
	for _, item := range props.Items {
		attrs := map[string]string{"type": "radio", "name": props.Name, "value": item.Value, "class": "radio"}
		if item.Checked {
			attrs["checked"] = "checked"
		}
		items = append(items, node("label", map[string]string{"class": "label cursor-pointer gap-2"}, node("input", attrs), textNode("span", map[string]string{"class": "label-text"}, item.Label)))
	}
	return node("div", map[string]string{"class": strings.TrimSpace("space-y-2 " + props.Props.Class)}, items...)
}

func Switch(props shared.SwitchComponentProps) shared.Node {
	attrs := map[string]string{"type": "checkbox", "class": strings.TrimSpace("toggle " + props.Props.Class), "name": props.Name, "value": props.Value}
	if props.Checked {
		attrs["checked"] = "checked"
	}
	return node("label", map[string]string{"class": "label cursor-pointer gap-2"}, node("input", attrs), textNode("span", map[string]string{"class": "label-text"}, props.Label))
}

func Pagination(props shared.PaginationProps) shared.Node {
	return node("div", map[string]string{"class": "join"},
		textNode("a", map[string]string{"class": "join-item btn", "href": props.PrevHref}, "«"),
		textNode("button", map[string]string{"class": "join-item btn"}, strconv.Itoa(props.Page)),
		textNode("a", map[string]string{"class": "join-item btn", "href": props.NextHref}, "»"),
	)
}

func EmptyState(props shared.EmptyStateProps) shared.Node {
	return node("div", map[string]string{"class": strings.TrimSpace("hero bg-base-200 rounded-box " + props.Props.Class)},
		node("div", map[string]string{"class": "hero-content text-center"},
			node("div", nil, textNode("h2", map[string]string{"class": "text-2xl font-bold"}, props.Title), textNode("p", nil, props.Description)),
		),
	)
}

func PageHeader(props shared.PageHeaderProps) shared.Node {
	return node("header", map[string]string{"class": strings.TrimSpace("mb-6 space-y-2 " + props.Props.Class)},
		textNode("h1", map[string]string{"class": "text-3xl font-bold"}, props.Title),
		textNode("p", map[string]string{"class": "text-base-content/70"}, props.Description),
		props.Actions,
	)
}

func Section(props shared.SectionProps, children ...shared.Node) shared.Node {
	nodes := make([]shared.Node, 0, len(children)+2)
	if props.Title != "" {
		nodes = append(nodes, textNode("h2", map[string]string{"class": "text-xl font-semibold"}, props.Title))
	}
	if props.Description != "" {
		nodes = append(nodes, textNode("p", map[string]string{"class": "text-base-content/70"}, props.Description))
	}
	nodes = append(nodes, children...)
	return node("section", map[string]string{"class": strings.TrimSpace("space-y-4 " + props.Props.Class)}, nodes...)
}

func Grid(props shared.GridProps, children ...shared.Node) shared.Node {
	className := "grid"
	if props.Columns != "" {
		className += " " + props.Columns
	}
	if props.Gap != "" {
		className += " " + props.Gap
	} else {
		className += " gap-4"
	}
	if props.Props.Class != "" {
		className += " " + props.Props.Class
	}
	return node("div", map[string]string{"class": className}, children...)
}

func Stack(props shared.StackProps, children ...shared.Node) shared.Node {
	className := "flex"
	if props.Direction != "" {
		className += " " + props.Direction
	} else {
		className += " flex-col"
	}
	if props.Gap != "" {
		className += " " + props.Gap
	} else {
		className += " gap-2"
	}
	if props.Props.Class != "" {
		className += " " + props.Props.Class
	}
	return node("div", map[string]string{"class": className}, children...)
}

func Breadcrumb(props shared.BreadcrumbProps) shared.Node {
	items := make([]shared.Node, 0, len(props.Items))
	for _, item := range props.Items {
		items = append(items, node("li", nil, textNode("a", map[string]string{"href": item.Href}, item.Label)))
	}
	return node("div", map[string]string{"class": strings.TrimSpace("breadcrumbs text-sm " + props.Props.Class)},
		node("ul", nil, items...),
	)
}

func Divider(props shared.DividerProps) shared.Node {
	className := "divider"
	if props.Props.Class != "" {
		className += " " + props.Props.Class
	}
	if props.Spacing != "" {
		className += " " + props.Spacing
	}
	return node("div", map[string]string{"class": className})
}

func Actions(props shared.ActionsProps, children ...shared.Node) shared.Node {
	return node("div", map[string]string{"class": strings.TrimSpace("flex items-center gap-2 " + props.Props.Class)}, children...)
}

func HiddenField(name, value string) shared.Node {
	return node("input", map[string]string{"type": "hidden", "name": name, "value": value})
}

func Box(props shared.BoxProps, children ...shared.Node) shared.Node {
	return node("div", map[string]string{"class": strings.TrimSpace("rounded-box border border-base-300 p-4 " + props.Props.Class)}, children...)
}

func AppShell(props shared.AppShellProps) shared.Node {
	attrs := map[string]string{"class": strings.TrimSpace("min-h-screen bg-base-100 " + props.Props.Class)}
	if props.ID != "" {
		attrs["id"] = props.ID
	}
	mainAttrs := map[string]string{"class": "mx-auto w-full max-w-7xl p-4 md:p-6"}
	if props.MainID != "" {
		mainAttrs["id"] = props.MainID
	}
	return node("div", attrs,
		props.Sidebar,
		props.Flashes,
		props.Header,
		node("main", mainAttrs, props.Content),
	)
}

func Image(props shared.ImageProps) shared.Node {
	attrs := map[string]string{"src": props.Src, "alt": props.Alt, "class": strings.TrimSpace("rounded-lg " + props.Props.Class)}
	if props.Width > 0 {
		attrs["width"] = strconv.Itoa(props.Width)
	}
	if props.Height > 0 {
		attrs["height"] = strconv.Itoa(props.Height)
	}
	return node("figure", nil, node("img", attrs), textNode("figcaption", map[string]string{"class": "text-sm mt-2"}, props.Caption))
}

func Chart(props shared.ChartProps) shared.Node {
	return node("div", map[string]string{"class": strings.TrimSpace("card bg-base-100 border border-base-300 " + props.Props.Class)},
		node("div", map[string]string{"class": "card-body"},
			textNode("h3", map[string]string{"class": "card-title"}, props.Title),
			textNode("p", map[string]string{"class": "text-sm opacity-70"}, props.Description),
		),
	)
}

func Form(props shared.FormProps, children ...shared.Node) shared.Node {
	attrs := map[string]string{
		"method": props.Method,
		"action": props.Action,
		"class":  strings.TrimSpace("space-y-4 " + props.Class),
	}
	if props.ID != "" {
		attrs["id"] = props.ID
	}
	return node("form", attrs, children...)
}

func ActionForm(props shared.ActionFormProps, children ...shared.Node) shared.Node {
	attrs := map[string]string{
		"method": props.Method,
		"action": props.Action,
		"class":  strings.TrimSpace("space-y-4 " + props.Props.Class),
	}
	if props.Target != "" {
		attrs["hx-target"] = props.Target
	}
	if props.Swap != "" {
		attrs["hx-swap"] = props.Swap
	}
	return node("form", attrs, children...)
}

func FormField(control shared.Node, props shared.FormFieldProps) shared.Node {
	children := []shared.Node{textNode("span", map[string]string{"class": "label-text"}, props.Label), control}
	if props.Hint != "" {
		children = append(children, textNode("span", map[string]string{"class": "label-text-alt"}, props.Hint))
	}
	if props.Error != "" {
		children = append(children, textNode("span", map[string]string{"class": "label-text-alt text-error"}, props.Error))
	}
	return node("label", map[string]string{"class": "form-control w-full gap-1"}, children...)
}

func Textarea(name, value string, options shared.TextareaOptions) shared.Node {
	attrs := map[string]string{
		"name":  name,
		"class": strings.TrimSpace("textarea textarea-bordered w-full " + options.Props.Class),
	}
	if options.Rows > 0 {
		attrs["rows"] = strconv.Itoa(options.Rows)
	}
	if options.Placeholder != "" {
		attrs["placeholder"] = options.Placeholder
	}
	if options.Required {
		attrs["required"] = "required"
	}
	return textNode("textarea", attrs, value)
}

func Region(props shared.RegionProps, children ...shared.Node) shared.Node {
	attrs := map[string]string{
		"class": strings.TrimSpace("space-y-3 " + props.Props.Class),
	}
	if props.ID != "" {
		attrs["id"] = props.ID
	}
	return node("section", attrs, children...)
}

func Split(props shared.SplitProps) shared.Node {
	wrapClass := "flex flex-col gap-4"
	if props.ReverseOnMobile {
		wrapClass = "flex flex-col-reverse gap-4"
	}
	if props.Props.Class != "" {
		wrapClass += " " + props.Props.Class
	}
	return node("div", map[string]string{"class": wrapClass},
		node("div", map[string]string{"class": "flex-1"}, props.Main),
		node("aside", map[string]string{"class": "w-full lg:w-80"}, props.Aside),
	)
}

func Container(props shared.ContainerProps, children ...shared.Node) shared.Node {
	className := "w-full"
	if props.MaxWidth != "" {
		className += " " + props.MaxWidth
	} else {
		className += " max-w-7xl"
	}
	if props.Centered {
		className += " mx-auto"
	}
	if props.Padding != "" {
		className += " " + props.Padding
	}
	if props.Props.Class != "" {
		className += " " + props.Props.Class
	}
	return node("div", map[string]string{"class": className}, children...)
}

func ThemeToggleButton(props shared.ComponentProps) shared.Node {
	className := strings.TrimSpace("btn btn-ghost " + props.Class)
	return node("button", map[string]string{"class": className, "type": "button", "aria-label": "Toggle theme"},
		textNode("span", nil, "🌓"),
	)
}

func Text(props shared.TextProps) shared.Node {
	return textNode("p", map[string]string{"class": strings.TrimSpace(props.Size + " " + props.Weight + " " + props.Props.Class)}, props.Text)
}

func FontIcon(props shared.FontIconProps) shared.Node {
	attrs := map[string]string{"class": strings.TrimSpace(props.Library + " " + props.Name + " " + props.Props.Class)}
	if props.AriaLabel != "" {
		attrs["aria-label"] = props.AriaLabel
	}
	if props.Decorative {
		attrs["aria-hidden"] = "true"
	}
	return node("i", attrs)
}

func ThemeToggle(props shared.ComponentProps) shared.Node {
	return ThemeToggleButton(props)
}

func HTMXTable(headers []string, rows ...shared.TableRowData) shared.Node {
	headerNodes := make([]shared.Node, 0, len(headers))
	for _, h := range headers {
		headerNodes = append(headerNodes, textNode("th", nil, h))
	}
	rowNodes := make([]shared.Node, 0, len(rows))
	for _, r := range rows {
		cells := make([]shared.Node, 0, len(r.Cells))
		for _, c := range r.Cells {
			cells = append(cells, node("td", nil, c))
		}
		rowNodes = append(rowNodes, node("tr", nil, cells...))
	}
	return node("div", map[string]string{"class": "overflow-x-auto"},
		node("table", map[string]string{"class": "table table-zebra w-full"},
			node("thead", nil, node("tr", nil, headerNodes...)),
			node("tbody", nil, rowNodes...),
		),
	)
}

func TableRow(cells ...shared.Node) shared.TableRowData {
	return shared.TableRowData{Cells: cells}
}

func SubmitButton(label string, props shared.ComponentProps) shared.Node {
	btn := Button(label, props)
	return btn
}

func InputWithOptions(name, value string, options shared.InputOptions) shared.Node {
	attrs := map[string]string{
		"name":  name,
		"value": value,
		"type":  options.Type,
		"class": strings.TrimSpace("input input-bordered w-full " + options.Props.Class),
	}
	if options.Placeholder != "" {
		attrs["placeholder"] = options.Placeholder
	}
	return node("input", attrs)
}

func FileUpload(name string, required bool, props ...shared.ComponentProps) shared.Node {
	p := shared.ComponentProps{}
	if len(props) > 0 {
		p = props[0]
	}
	attrs := map[string]string{"type": "file", "name": name, "class": strings.TrimSpace("file-input file-input-bordered w-full " + p.Class)}
	if required {
		attrs["required"] = "required"
	}
	return node("input", attrs)
}

func Sidebar(brand, title string, items ...shared.SidebarItem) shared.Node {
	nodes := make([]shared.Node, 0, len(items))
	for _, item := range items {
		nodes = append(nodes, node("li", nil, textNode("a", map[string]string{"href": item.Href}, item.Label)))
	}
	return node("aside", map[string]string{"class": "w-80 bg-base-200 p-4"},
		textNode("div", map[string]string{"class": "text-lg font-bold mb-1"}, brand),
		textNode("div", map[string]string{"class": "text-sm opacity-70 mb-4"}, title),
		node("ul", map[string]string{"class": "menu"}, nodes...),
	)
}

func SidebarLink(label, href string) shared.SidebarItem {
	return shared.SidebarItem{Label: label, Href: href}
}

func DownloadLink(label, href, filename string, props shared.ComponentProps) shared.Node {
	attrs := map[string]string{"href": href, "download": filename, "class": strings.TrimSpace("link link-primary " + props.Class)}
	return textNode("a", attrs, label)
}

func DrawerLayout(drawerID string, navbar, content shared.Node, sidebarItems []shared.Node) shared.Node {
	if drawerID == "" {
		drawerID = "drawer"
	}
	items := make([]shared.Node, 0, len(sidebarItems))
	for _, item := range sidebarItems {
		items = append(items, node("li", nil, item))
	}
	return node("div", map[string]string{"class": "drawer lg:drawer-open"},
		node("input", map[string]string{"id": drawerID, "type": "checkbox", "class": "drawer-toggle"}),
		node("div", map[string]string{"class": "drawer-content flex flex-col"}, navbar, content),
		node("div", map[string]string{"class": "drawer-side"},
			node("label", map[string]string{"for": drawerID, "aria-label": "close sidebar", "class": "drawer-overlay"}),
			node("ul", map[string]string{"class": "menu bg-base-200 min-h-full w-80 p-4"}, items...),
		),
	)
}

func DrawerNavbar(drawerID, title string, desktopItems []shared.Node) shared.Node {
	if drawerID == "" {
		drawerID = "drawer"
	}
	if title == "" {
		title = "Navbar Title"
	}
	menuItems := make([]shared.Node, 0, len(desktopItems))
	for _, item := range desktopItems {
		menuItems = append(menuItems, node("li", nil, item))
	}
	return node("div", map[string]string{"class": "navbar bg-base-300 w-full"},
		node("div", map[string]string{"class": "flex-none lg:hidden"},
			node("label", map[string]string{"for": drawerID, "aria-label": "open sidebar", "class": "btn btn-square btn-ghost"},
				node("svg", map[string]string{"xmlns": "http://www.w3.org/2000/svg", "fill": "none", "viewBox": "0 0 24 24", "class": "inline-block h-6 w-6 stroke-current"},
					node("path", map[string]string{"stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", "d": "M4 6h16M4 12h16M4 18h16"}),
				),
			),
		),
		textNode("div", map[string]string{"class": "mx-2 flex-1 px-2"}, title),
		node("div", map[string]string{"class": "hidden flex-none lg:block"},
			node("ul", map[string]string{"class": "menu menu-horizontal"}, menuItems...),
		),
	)
}

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
