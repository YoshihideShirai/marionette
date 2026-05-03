package daisyui

import (
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
