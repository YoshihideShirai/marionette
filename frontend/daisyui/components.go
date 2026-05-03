package daisyui

import (
	"strconv"

	frontend "github.com/YoshihideShirai/marionette/frontend"
	lowhtml "github.com/YoshihideShirai/marionette/frontend/html"
)

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

func Checkbox(props frontend.CheckboxComponentProps) frontend.Node {
	return frontend.UICheckbox(props)
}

func RadioGroup(props frontend.RadioGroupComponentProps) frontend.Node {
	return frontend.UIRadioGroup(props)
}

func Switch(props frontend.SwitchComponentProps) frontend.Node {
	return frontend.UISwitch(props)
}

func Pagination(props frontend.PaginationProps) frontend.Node {
	return frontend.UIPagination(props)
}

func EmptyState(props frontend.EmptyStateProps) frontend.Node {
	return frontend.UIEmptyState(props)
}

func PageHeader(props frontend.PageHeaderProps) frontend.Node {
	return frontend.UIPageHeader(props)
}

func Section(props frontend.SectionProps, children ...frontend.Node) frontend.Node {
	return frontend.UISection(props, children...)
}

func Grid(props frontend.GridProps, children ...frontend.Node) frontend.Node {
	return frontend.UIGrid(props, children...)
}

func Stack(props frontend.StackProps, children ...frontend.Node) frontend.Node {
	return frontend.UIStack(props, children...)
}

func Breadcrumb(props frontend.BreadcrumbProps) frontend.Node {
	return frontend.UIBreadcrumb(props)
}

func Divider(props frontend.DividerProps) frontend.Node {
	return frontend.UIDivider(props)
}

func Actions(props frontend.ActionsProps, children ...frontend.Node) frontend.Node {
	return frontend.UIActions(props, children...)
}

func HiddenField(name, value string) frontend.Node {
	return frontend.UIHiddenField(name, value)
}

func Box(props frontend.BoxProps, children ...frontend.Node) frontend.Node {
	return frontend.UIBox(props, children...)
}

func AppShell(props frontend.AppShellProps) frontend.Node {
	return frontend.UIAppShell(props)
}

func Image(props frontend.ImageProps) frontend.Node {
	return frontend.UIImage(props)
}

func Chart(props frontend.ChartProps) frontend.Node {
	return frontend.UIChart(props)
}

func Form(props frontend.FormProps, children ...frontend.Node) frontend.Node {
	return frontend.UIForm(props, children...)
}

func ActionForm(props frontend.ActionFormProps, children ...frontend.Node) frontend.Node {
	return frontend.UIActionForm(props, children...)
}

func FormField(control frontend.Node, props frontend.FormFieldProps) frontend.Node {
	return frontend.UIFormField(control, props)
}

func Textarea(name, value string, options frontend.TextareaOptions) frontend.Node {
	return frontend.UITextarea(name, value, options)
}

func Region(props frontend.RegionProps, children ...frontend.Node) frontend.Node {
	return frontend.UIRegion(props, children...)
}

func Split(props frontend.SplitProps) frontend.Node {
	return frontend.UISplit(props)
}

func Container(props frontend.ContainerProps, children ...frontend.Node) frontend.Node {
	return frontend.Container(props, children...)
}

func ThemeToggleButton(props frontend.ComponentProps) frontend.Node {
	return frontend.UIThemeToggleButton(props)
}

func Text(props frontend.TextProps) frontend.Node {
	return frontend.UIText(props)
}

func FontIcon(props frontend.FontIconProps) frontend.Node {
	return frontend.UIFontIcon(props)
}

func ThemeToggle(props frontend.ComponentProps) frontend.Node {
	return frontend.UIThemeToggleButton(props)
}

func HTMXTable(headers []string, rows ...frontend.TableRowData) frontend.Node {
	return frontend.HTMXTable(headers, rows...)
}

func TableRow(cells ...frontend.Node) frontend.TableRowData {
	return frontend.TableRow(cells...)
}

func SubmitButton(label string, props frontend.ComponentProps) frontend.Node {
	return frontend.UISubmitButton(label, props)
}

func InputWithOptions(name, value string, options frontend.InputOptions) frontend.Node {
	return frontend.UIInputWithOptions(name, value, options)
}

func FileUpload(name string, required bool, props ...frontend.ComponentProps) frontend.Node {
	return frontend.FileUpload(name, required, props...)
}

func Sidebar(brand, title string, items ...frontend.SidebarItem) frontend.Node {
	return frontend.Sidebar(brand, title, items...)
}

func SidebarLink(label, href string) frontend.SidebarItem {
	return frontend.SidebarLink(label, href)
}

func DownloadLink(label, href, filename string, props frontend.ComponentProps) frontend.Node {
	return frontend.DownloadLink(label, href, filename, props)
}

func H1(children ...frontend.Node) frontend.Node { return frontend.H1(children...) }
func H2(children ...frontend.Node) frontend.Node { return frontend.H2(children...) }
func H3(children ...frontend.Node) frontend.Node { return frontend.H3(children...) }
func H4(children ...frontend.Node) frontend.Node { return frontend.H4(children...) }
func TextNode(text string) frontend.Node         { return frontend.Text(text) }

func PrimaryButton(label string, props frontend.ComponentProps) frontend.Node {
	if props.Variant == "" {
		props.Variant = "primary"
	}
	return frontend.Button(label, props)
}

func SecondaryButton(label string, props frontend.ComponentProps) frontend.Node {
	if props.Variant == "" {
		props.Variant = "secondary"
	}
	return frontend.Button(label, props)
}

func GhostButton(label string, props frontend.ComponentProps) frontend.Node {
	if props.Variant == "" {
		props.Variant = "ghost"
	}
	return frontend.Button(label, props)
}

// Avatar follows daisyUI's avatar markup: .avatar > .w-*/mask wrapper > img
func Avatar(src, alt, class string) frontend.Node {
	return lowhtml.ElementNode{
		Tag:   "div",
		Attrs: map[string]string{"class": "avatar"},
		Children: []frontend.Node{
			lowhtml.ElementNode{
				Tag:      "div",
				Attrs:    map[string]string{"class": class},
				Children: []frontend.Node{lowhtml.ElementNode{Tag: "img", Attrs: map[string]string{"src": src, "alt": alt}}},
			},
		},
	}
}

func Navbar(start, center, end frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "navbar bg-base-100 shadow-sm"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "navbar-start"}, Children: []frontend.Node{start}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "navbar-center"}, Children: []frontend.Node{center}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "navbar-end"}, Children: []frontend.Node{end}},
	}}
}

func Hero(title, description string, actions ...frontend.Node) frontend.Node {
	children := []frontend.Node{
		lowhtml.ElementNode{Tag: "h1", Attrs: map[string]string{"class": "text-5xl font-bold"}, Text: title},
		lowhtml.ElementNode{Tag: "p", Attrs: map[string]string{"class": "py-6"}, Text: description},
	}
	if len(actions) > 0 {
		children = append(children, lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "flex gap-2"}, Children: actions})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "hero bg-base-200 rounded-box"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "hero-content text-center"}, Children: []frontend.Node{
			lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "max-w-md"}, Children: children},
		}},
	}}
}

func Menu(items ...frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": "menu bg-base-200 rounded-box"}, Children: items}
}

func Footer(children ...frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "footer", Attrs: map[string]string{"class": "footer sm:footer-horizontal bg-base-200 text-base-content p-10"}, Children: children}
}

func Drawer(id string, side, content frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "drawer", "id": id}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "input", Attrs: map[string]string{"id": id + "-toggle", "type": "checkbox", "class": "drawer-toggle"}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "drawer-content"}, Children: []frontend.Node{content}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "drawer-side"}, Children: []frontend.Node{side}},
	}}
}

func Stat(title, value, desc string) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "stats shadow"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "stat"}, Children: []frontend.Node{
			lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "stat-title"}, Text: title},
			lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "stat-value"}, Text: value},
			lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "stat-desc"}, Text: desc},
		}},
	}}
}

func Steps(items ...frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": "steps steps-vertical lg:steps-horizontal"}, Children: items}
}

func Step(label string, active bool) frontend.Node {
	className := "step"
	if active {
		className = "step step-primary"
	}
	return lowhtml.ElementNode{Tag: "li", Attrs: map[string]string{"class": className}, Text: label}
}

func Timeline(items ...frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": "timeline timeline-vertical"}, Children: items}
}

func TimelineItem(startLabel, endLabel string, content frontend.Node) frontend.Node {
	children := []frontend.Node{}
	if startLabel != "" {
		children = append(children, lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "timeline-start"}, Text: startLabel})
	}
	children = append(children,
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "timeline-middle"}, Text: "●"},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "timeline-end timeline-box"}, Children: []frontend.Node{content}},
	)
	if endLabel != "" {
		children = append(children, lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "timeline-end"}, Text: endLabel})
	}
	return lowhtml.ElementNode{Tag: "li", Children: children}
}

func Collapse(title string, content frontend.Node, open bool) frontend.Node {
	className := "collapse collapse-arrow bg-base-100 border border-base-300"
	attrs := map[string]string{"class": className}
	if open {
		attrs["open"] = "open"
	}
	return lowhtml.ElementNode{Tag: "details", Attrs: attrs, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "summary", Attrs: map[string]string{"class": "collapse-title font-semibold"}, Text: title},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "collapse-content text-sm"}, Children: []frontend.Node{content}},
	}}
}

func MockupWindow(title string, content frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-window border border-base-300"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "px-4 py-16 bg-base-200"}, Children: []frontend.Node{content}},
		lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": "sr-only"}, Text: title},
	}}
}

func Kbd(text string) frontend.Node {
	return lowhtml.ElementNode{Tag: "kbd", Attrs: map[string]string{"class": "kbd"}, Text: text}
}

func Code(text string) frontend.Node {
	return lowhtml.ElementNode{Tag: "code", Attrs: map[string]string{"class": "bg-base-200 rounded px-1 py-0.5"}, Text: text}
}

func Indicator(item, target frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "indicator"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": "indicator-item badge badge-secondary"}, Children: []frontend.Node{item}},
		target,
	}}
}

func Link(label, href string, props frontend.ComponentProps) frontend.Node {
	className := "link"
	if props.Class != "" {
		className += " " + props.Class
	}
	return lowhtml.ElementNode{Tag: "a", Attrs: map[string]string{"class": className, "href": href}, Text: label}
}

func Dropdown(trigger, menu frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "dropdown"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"tabindex": "0", "role": "button"}, Children: []frontend.Node{trigger}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"tabindex": "0", "class": "dropdown-content z-1 card card-sm bg-base-100 shadow-md"}, Children: []frontend.Node{menu}},
	}}
}

func Tooltip(text string, child frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "tooltip", "data-tip": text}, Children: []frontend.Node{child}}
}

func Loading(sizeClass string) frontend.Node {
	className := "loading loading-spinner"
	if sizeClass != "" {
		className += " " + sizeClass
	}
	return lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": className}}
}

func RadialProgress(value int, sizeClass string) frontend.Node {
	attrs := map[string]string{"class": "radial-progress " + sizeClass, "style": "--value:" + strconv.Itoa(value) + ";", "role": "progressbar"}
	return lowhtml.ElementNode{Tag: "div", Attrs: attrs, Text: strconv.Itoa(value) + "%"}
}

func Rating(name string, max int, checked int) frontend.Node {
	stars := make([]frontend.Node, 0, max)
	for i := 1; i <= max; i++ {
		attrs := map[string]string{"type": "radio", "name": name, "class": "mask mask-star-2 bg-orange-400", "value": strconv.Itoa(i)}
		if i == checked {
			attrs["checked"] = "checked"
		}
		stars = append(stars, lowhtml.ElementNode{Tag: "input", Attrs: attrs})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "rating"}, Children: stars}
}

func Range(name string, value int, min int, max int) frontend.Node {
	return lowhtml.ElementNode{Tag: "input", Attrs: map[string]string{"type": "range", "name": name, "value": strconv.Itoa(value), "min": strconv.Itoa(min), "max": strconv.Itoa(max), "class": "range"}}
}

func Toggle(name string, checked bool) frontend.Node {
	attrs := map[string]string{"type": "checkbox", "name": name, "class": "toggle"}
	if checked {
		attrs["checked"] = "checked"
	}
	return lowhtml.ElementNode{Tag: "input", Attrs: attrs}
}

func Join(children ...frontend.Node) frontend.Node {
	wrapped := make([]frontend.Node, 0, len(children))
	for _, child := range children {
		wrapped = append(wrapped, lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "join-item"}, Children: []frontend.Node{child}})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "join"}, Children: wrapped}
}

func Mask(shapeClass string, child frontend.Node) frontend.Node {
	className := "mask " + shapeClass
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": className}, Children: []frontend.Node{child}}
}

func Carousel(items ...frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "carousel w-full"}, Children: items}
}

func CarouselItem(id string, child frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"id": id, "class": "carousel-item w-full"}, Children: []frontend.Node{child}}
}

func ChatBubble(content frontend.Node, end bool) frontend.Node {
	position := "chat-start"
	if end {
		position = "chat-end"
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "chat " + position}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "chat-bubble"}, Children: []frontend.Node{content}},
	}}
}

func Countdown(value int) frontend.Node {
	return lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": "countdown font-mono text-2xl"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"style": "--value:" + strconv.Itoa(value) + ";"}},
	}}
}

func Status(colorClass string) frontend.Node {
	return lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": "status " + colorClass}}
}

func Dock(items ...frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "dock"}, Children: items}
}

func Fieldset(legend string, fields ...frontend.Node) frontend.Node {
	children := []frontend.Node{lowhtml.ElementNode{Tag: "legend", Attrs: map[string]string{"class": "fieldset-legend"}, Text: legend}}
	children = append(children, fields...)
	return lowhtml.ElementNode{Tag: "fieldset", Attrs: map[string]string{"class": "fieldset"}, Children: children}
}

func Label(text string) frontend.Node {
	return lowhtml.ElementNode{Tag: "label", Attrs: map[string]string{"class": "label"}, Text: text}
}

func Validator(message string) frontend.Node {
	return lowhtml.ElementNode{Tag: "p", Attrs: map[string]string{"class": "validator-hint"}, Text: message}
}

func BrowserMockup(content frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-browser border border-base-300"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-browser-toolbar"}, Children: []frontend.Node{lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "input"}, Text: "https://example.com"}}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "grid place-content-center h-80"}, Children: []frontend.Node{content}},
	}}
}

func PhoneMockup(content frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-phone border-primary"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-phone-camera"}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-phone-display"}, Children: []frontend.Node{content}},
	}}
}

func CodeMockup(lines ...string) frontend.Node {
	children := make([]frontend.Node, 0, len(lines))
	for _, line := range lines {
		children = append(children, lowhtml.ElementNode{Tag: "pre", Children: []frontend.Node{lowhtml.ElementNode{Tag: "code", Text: line}}})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "mockup-code"}, Children: children}
}

func Calendar(content frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "card bg-base-100 border border-base-300"}, Children: []frontend.Node{content}}
}

func Filter(items ...frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "filter"}, Children: items}
}

func Diff(before, after frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "figure", Attrs: map[string]string{"class": "diff aspect-16/9"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "diff-item-1"}, Children: []frontend.Node{before}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "diff-item-2"}, Children: []frontend.Node{after}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "diff-resizer"}},
	}}
}

func List(items ...frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": "list bg-base-100 rounded-box shadow-md"}, Children: items}
}

func Table(headers []string, rows ...[]frontend.Node) frontend.Node {
	headersNode := make([]frontend.Node, 0, len(headers))
	for _, h := range headers {
		headersNode = append(headersNode, lowhtml.ElementNode{Tag: "th", Text: h})
	}
	tbodyRows := make([]frontend.Node, 0, len(rows))
	for _, row := range rows {
		cells := make([]frontend.Node, 0, len(row))
		for _, cell := range row {
			cells = append(cells, lowhtml.ElementNode{Tag: "td", Children: []frontend.Node{cell}})
		}
		tbodyRows = append(tbodyRows, lowhtml.ElementNode{Tag: "tr", Children: cells})
	}
	return lowhtml.ElementNode{Tag: "table", Attrs: map[string]string{"class": "table"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "thead", Children: []frontend.Node{lowhtml.ElementNode{Tag: "tr", Children: headersNode}}},
		lowhtml.ElementNode{Tag: "tbody", Children: tbodyRows},
	}}
}

func TextRotate(words []string, animationClass string) frontend.Node {
	if animationClass == "" {
		animationClass = "animate-pulse"
	}
	items := make([]frontend.Node, 0, len(words))
	for _, w := range words {
		items = append(items, lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": animationClass}, Text: w})
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "flex gap-2"}, Children: items}
}

func Hover3DCard(content frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "card bg-base-100 shadow-xl transition-transform duration-300 hover:-translate-y-1 hover:shadow-2xl"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "card-body"}, Children: []frontend.Node{content}},
	}}
}

func HoverGallery(items ...frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "grid grid-cols-2 md:grid-cols-3 gap-4"}, Children: items}
}

func Accordion(title string, content frontend.Node, open bool) frontend.Node {
	return Collapse(title, content, open)
}

func FAB(icon frontend.Node, label string) frontend.Node {
	return lowhtml.ElementNode{Tag: "button", Attrs: map[string]string{"class": "btn btn-primary btn-circle fixed bottom-6 right-6"}, Children: []frontend.Node{
		icon,
		lowhtml.ElementNode{Tag: "span", Attrs: map[string]string{"class": "sr-only"}, Text: label},
	}}
}

func SpeedDial(trigger frontend.Node, items ...frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "dropdown dropdown-top dropdown-end fixed bottom-6 right-6"}, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"tabindex": "0", "role": "button"}, Children: []frontend.Node{trigger}},
		lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"tabindex": "0", "class": "dropdown-content menu bg-base-100 rounded-box z-1 w-52 p-2 shadow"}, Children: items},
	}}
}

func Swap(onNode, offNode frontend.Node, active bool) frontend.Node {
	attrs := map[string]string{"class": "swap"}
	if active {
		attrs["class"] = "swap swap-active"
	}
	return lowhtml.ElementNode{Tag: "label", Attrs: attrs, Children: []frontend.Node{
		lowhtml.ElementNode{Tag: "input", Attrs: map[string]string{"type": "checkbox"}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "swap-on"}, Children: []frontend.Node{onNode}},
		lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "swap-off"}, Children: []frontend.Node{offNode}},
	}}
}

func ThemeController(options ...frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "join"}, Children: options}
}

func DockItem(child frontend.Node, active bool) frontend.Node {
	className := "dock-label"
	if active {
		className = "dock-active"
	}
	return lowhtml.ElementNode{Tag: "button", Attrs: map[string]string{"class": className}, Children: []frontend.Node{child}}
}

func FilterItem(label string, active bool) frontend.Node {
	className := "btn btn-sm"
	if active {
		className += " btn-active"
	}
	return lowhtml.ElementNode{Tag: "button", Attrs: map[string]string{"class": className}, Text: label}
}

func CalendarGrid(days ...frontend.Node) frontend.Node {
	return lowhtml.ElementNode{Tag: "div", Attrs: map[string]string{"class": "grid grid-cols-7 gap-1"}, Children: days}
}
