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
