package marionette

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// ComponentProps defines shared style knobs for template components.
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

type ModalProps struct {
	Title   string
	Body    Node
	Actions Node
	Open    bool
}

type FormFieldProps struct {
	Label    string
	Required bool
	Hint     string
	Error    string
}

type InputOptions struct {
	Type        string
	Placeholder string
	Min         string
	Max         string
	Required    bool
	Props       ComponentProps
}

type TextareaOptions struct {
	Placeholder string
	Rows        int
	Required    bool
	Props       ComponentProps
}

type EmptyStateProps struct {
	Title       string
	Description string
	Skeleton    bool
	Rows        int
	Icon        string
	Props       ComponentProps
}

type AlertProps struct {
	Title       string
	Description string
	Icon        string
	Props       ComponentProps
}

type ToastProps struct {
	Title       string
	Description string
	Icon        string
	Props       ComponentProps
	Live        string
}

type SkeletonProps struct {
	Rows  int
	Props ComponentProps
}

type TableColumn struct {
	Label      string
	SortKey    string
	SortHref   string
	SortActive bool
}

type TableComponentRow struct {
	Cells []Node
}

type TableProps struct {
	Columns          []TableColumn
	Rows             []TableComponentRow
	EmptyTitle       string
	EmptyDescription string
}

type PaginationProps struct {
	Page       int
	TotalPages int
	PrevHref   string
	NextHref   string
}

type TabsItem struct {
	Label    string
	Href     string
	Active   bool
	Disabled bool
}

type TabsProps struct {
	Items     []TabsItem
	AriaLabel string
	Props     ComponentProps
}

type BreadcrumbItem struct {
	Label  string
	Href   string
	Active bool
}

type BreadcrumbProps struct {
	Items     []BreadcrumbItem
	AriaLabel string
	Props     ComponentProps
}

type CheckboxComponentProps struct {
	Name    string
	Value   string
	Label   string
	Checked bool
	Props   ComponentProps
}

type RadioItem struct {
	Label    string
	Value    string
	Checked  bool
	Disabled bool
}

type RadioGroupComponentProps struct {
	Name      string
	Items     []RadioItem
	AriaLabel string
	Props     ComponentProps
}

type SwitchComponentProps struct {
	Name    string
	Value   string
	Label   string
	Checked bool
	Props   ComponentProps
}

type StackProps struct {
	Direction string
	Gap       string
	Align     string
	Justify   string
	Wrap      bool
	Props     ComponentProps
}

type GridProps struct {
	Columns        string
	Gap            string
	MinColumnWidth string
	Props          ComponentProps
}

type SplitProps struct {
	Main            Node
	Aside           Node
	AsideWidth      string
	ReverseOnMobile bool
	Gap             string
	Props           ComponentProps
}

type PageHeaderProps struct {
	Title       string
	Description string
	Actions     Node
	Props       ComponentProps
}

type ContainerProps struct {
	MaxWidth string
	Padding  string
	Centered bool
	Props    ComponentProps
}

type CardProps struct {
	Title       string
	Description string
	Actions     Node
	Props       ComponentProps
}

type SectionProps struct {
	Title       string
	Description string
	Actions     Node
	Props       ComponentProps
}

type templateNode struct {
	name string
	data any
}

var (
	cachedTemplates        *template.Template
	cachedTemplatesErr     error
	componentTemplatesOnce sync.Once
)

func (n templateNode) Render() (template.HTML, error) {
	tmpl, err := loadComponentTemplates()
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	if err := tmpl.ExecuteTemplate(&out, n.name, n.data); err != nil {
		return "", err
	}
	return template.HTML(out.String()), nil
}

func ComponentButton(label string, props ComponentProps) Node {
	return componentButton(label, "button", props)
}

func ComponentSubmitButton(label string, props ComponentProps) Node {
	return componentButton(label, "submit", props)
}

func componentButton(label, buttonType string, props ComponentProps) Node {
	return templateNode{
		name: "components/button",
		data: struct {
			Class    string
			Type     string
			Label    string
			Disabled bool
		}{
			Class:    buttonClass(props),
			Type:     buttonType,
			Label:    label,
			Disabled: props.Disabled,
		},
	}
}

func ComponentInput(name, value string, props ComponentProps) Node {
	return ComponentInputWithOptions(name, value, InputOptions{
		Type:        "text",
		Placeholder: strings.TrimSpace(name),
		Props:       props,
	})
}

func ComponentInputWithOptions(name, value string, options InputOptions) Node {
	inputType := strings.TrimSpace(options.Type)
	if inputType == "" {
		inputType = "text"
	}
	return templateNode{
		name: "components/input",
		data: struct {
			Class       string
			Name        string
			Type        string
			Value       string
			Placeholder string
			Min         string
			Max         string
			Required    bool
			Disabled    bool
		}{
			Class:       inputClass(options.Props),
			Name:        name,
			Type:        inputType,
			Value:       value,
			Placeholder: options.Placeholder,
			Min:         strings.TrimSpace(options.Min),
			Max:         strings.TrimSpace(options.Max),
			Required:    options.Required,
			Disabled:    options.Props.Disabled,
		},
	}
}

func ComponentTextarea(name, value string, options TextareaOptions) Node {
	rows := options.Rows
	if rows <= 0 {
		rows = 3
	}
	return templateNode{
		name: "components/textarea",
		data: struct {
			Class       string
			Name        string
			Value       string
			Placeholder string
			Rows        int
			Required    bool
			Disabled    bool
		}{
			Class:       textareaClass(options.Props),
			Name:        strings.TrimSpace(name),
			Value:       value,
			Placeholder: strings.TrimSpace(options.Placeholder),
			Rows:        rows,
			Required:    options.Required,
			Disabled:    options.Props.Disabled,
		},
	}
}

func ComponentFormField(control Node, props FormFieldProps) Node {
	controlHTML, err := renderNode(control)
	if err != nil {
		return renderErrorNode{err: err}
	}
	return templateNode{
		name: "components/form_field",
		data: struct {
			Label    string
			Required bool
			Hint     string
			Error    string
			Control  template.HTML
		}{
			Label:    strings.TrimSpace(props.Label),
			Required: props.Required,
			Hint:     strings.TrimSpace(props.Hint),
			Error:    strings.TrimSpace(props.Error),
			Control:  controlHTML,
		},
	}
}

func ComponentSelect(name string, options []SelectOption, props ComponentProps) Node {
	return templateNode{
		name: "components/select",
		data: struct {
			Class    string
			Name     string
			Options  []SelectOption
			Disabled bool
		}{
			Class:    selectClass(props),
			Name:     name,
			Options:  options,
			Disabled: props.Disabled,
		},
	}
}

func ComponentModal(props ModalProps) Node {
	bodyHTML, err := renderNode(props.Body)
	if err != nil {
		return renderErrorNode{err: err}
	}
	actionsHTML, err := renderNode(props.Actions)
	if err != nil {
		return renderErrorNode{err: err}
	}
	return templateNode{
		name: "components/modal",
		data: struct {
			Title   string
			Body    template.HTML
			Actions template.HTML
			Open    bool
		}{
			Title:   props.Title,
			Body:    bodyHTML,
			Actions: actionsHTML,
			Open:    props.Open,
		},
	}
}

func ComponentToast(props ToastProps) Node {
	live := strings.TrimSpace(props.Live)
	if live == "" {
		live = "polite"
	}
	return templateNode{
		name: "components/toast",
		data: struct {
			Class       string
			Title       string
			Description string
			Icon        string
			Live        string
		}{
			Class:       feedbackClass("toast", props.Props),
			Title:       strings.TrimSpace(props.Title),
			Description: strings.TrimSpace(props.Description),
			Icon:        strings.TrimSpace(props.Icon),
			Live:        live,
		},
	}
}

func ComponentAlert(props AlertProps) Node {
	return templateNode{
		name: "components/alert",
		data: struct {
			Class       string
			Title       string
			Description string
			Icon        string
		}{
			Class:       feedbackClass("alert", props.Props),
			Title:       strings.TrimSpace(props.Title),
			Description: strings.TrimSpace(props.Description),
			Icon:        strings.TrimSpace(props.Icon),
		},
	}
}

func ComponentSkeleton(props SkeletonProps) Node {
	rows := props.Rows
	if rows <= 0 {
		rows = 3
	}
	return templateNode{
		name: "components/skeleton",
		data: struct {
			Class string
			Rows  []int
		}{
			Class: feedbackClass("skeleton", props.Props),
			Rows:  make([]int, rows),
		},
	}
}

func ComponentEmptyState(props EmptyStateProps) Node {
	rows := props.Rows
	if rows <= 0 {
		rows = 3
	}
	return templateNode{
		name: "components/empty_state",
		data: struct {
			Class       string
			Title       string
			Description string
			Skeleton    bool
			Rows        []int
			Icon        string
		}{
			Class:       feedbackClass("empty-state", props.Props),
			Title:       strings.TrimSpace(props.Title),
			Description: strings.TrimSpace(props.Description),
			Skeleton:    props.Skeleton,
			Rows:        make([]int, rows),
			Icon:        strings.TrimSpace(props.Icon),
		},
	}
}

func ComponentTable(props TableProps) Node {
	rows := make([]struct {
		Cells []template.HTML
	}, 0, len(props.Rows))
	for _, row := range props.Rows {
		cells := make([]template.HTML, 0, len(row.Cells))
		for _, cell := range row.Cells {
			cellHTML, err := renderNode(cell)
			if err != nil {
				return renderErrorNode{err: err}
			}
			cells = append(cells, cellHTML)
		}
		rows = append(rows, struct {
			Cells []template.HTML
		}{Cells: cells})
	}

	return templateNode{
		name: "components/table",
		data: struct {
			Columns          []TableColumn
			Rows             []struct{ Cells []template.HTML }
			EmptyTitle       string
			EmptyDescription string
		}{
			Columns:          props.Columns,
			Rows:             rows,
			EmptyTitle:       strings.TrimSpace(props.EmptyTitle),
			EmptyDescription: strings.TrimSpace(props.EmptyDescription),
		},
	}
}

func ComponentPagination(props PaginationProps) Node {
	page := props.Page
	if page < 1 {
		page = 1
	}
	totalPages := props.TotalPages
	if totalPages < 1 {
		totalPages = 1
	}
	return templateNode{
		name: "components/pagination",
		data: struct {
			Page       int
			TotalPages int
			PrevHref   string
			NextHref   string
		}{
			Page:       page,
			TotalPages: totalPages,
			PrevHref:   strings.TrimSpace(props.PrevHref),
			NextHref:   strings.TrimSpace(props.NextHref),
		},
	}
}

func ComponentTabs(props TabsProps) Node {
	items := make([]TabsItem, 0, len(props.Items))
	for _, item := range props.Items {
		items = append(items, TabsItem{
			Label:    strings.TrimSpace(item.Label),
			Href:     strings.TrimSpace(item.Href),
			Active:   item.Active,
			Disabled: item.Disabled,
		})
	}
	ariaLabel := strings.TrimSpace(props.AriaLabel)
	if ariaLabel == "" {
		ariaLabel = "tabs"
	}
	return templateNode{
		name: "components/tabs",
		data: struct {
			Class     string
			AriaLabel string
			Items     []TabsItem
		}{
			Class:     joinClass("tabs tabs-boxed", props.Props.Class),
			AriaLabel: ariaLabel,
			Items:     items,
		},
	}
}

func ComponentBreadcrumb(props BreadcrumbProps) Node {
	items := make([]BreadcrumbItem, 0, len(props.Items))
	for _, item := range props.Items {
		items = append(items, BreadcrumbItem{
			Label:  strings.TrimSpace(item.Label),
			Href:   strings.TrimSpace(item.Href),
			Active: item.Active,
		})
	}
	ariaLabel := strings.TrimSpace(props.AriaLabel)
	if ariaLabel == "" {
		ariaLabel = "breadcrumb"
	}
	return templateNode{
		name: "components/breadcrumb",
		data: struct {
			Class     string
			AriaLabel string
			Items     []BreadcrumbItem
		}{
			Class:     joinClass("breadcrumbs text-sm", props.Props.Class),
			AriaLabel: ariaLabel,
			Items:     items,
		},
	}
}

func ComponentCheckbox(props CheckboxComponentProps) Node {
	return templateNode{
		name: "components/checkbox",
		data: struct {
			Label    string
			Name     string
			Value    string
			Class    string
			Checked  bool
			Disabled bool
		}{
			Label:    strings.TrimSpace(props.Label),
			Name:     strings.TrimSpace(props.Name),
			Value:    strings.TrimSpace(props.Value),
			Class:    checkboxClass(props.Props),
			Checked:  props.Checked,
			Disabled: props.Props.Disabled,
		},
	}
}

func ComponentRadioGroup(props RadioGroupComponentProps) Node {
	items := make([]RadioItem, 0, len(props.Items))
	for _, item := range props.Items {
		items = append(items, RadioItem{
			Label:    strings.TrimSpace(item.Label),
			Value:    strings.TrimSpace(item.Value),
			Checked:  item.Checked,
			Disabled: item.Disabled,
		})
	}
	ariaLabel := strings.TrimSpace(props.AriaLabel)
	if ariaLabel == "" {
		ariaLabel = "radio group"
	}
	return templateNode{
		name: "components/radio_group",
		data: struct {
			Name      string
			Class     string
			AriaLabel string
			Items     []RadioItem
			Disabled  bool
		}{
			Name:      strings.TrimSpace(props.Name),
			Class:     radioClass(props.Props),
			AriaLabel: ariaLabel,
			Items:     items,
			Disabled:  props.Props.Disabled,
		},
	}
}

func ComponentSwitch(props SwitchComponentProps) Node {
	return templateNode{
		name: "components/switch",
		data: struct {
			Label    string
			Name     string
			Value    string
			Class    string
			Checked  bool
			Disabled bool
		}{
			Label:    strings.TrimSpace(props.Label),
			Name:     strings.TrimSpace(props.Name),
			Value:    strings.TrimSpace(props.Value),
			Class:    switchClass(props.Props),
			Checked:  props.Checked,
			Disabled: props.Props.Disabled,
		},
	}
}

func ComponentStack(props StackProps, children ...Node) Node {
	return layoutChildrenNode("components/stack", stackClass(props), children)
}

func ComponentGrid(props GridProps, children ...Node) Node {
	return layoutChildrenNode("components/grid", gridClass(props), children)
}

func ComponentSplit(props SplitProps) Node {
	mainHTML, err := renderNode(props.Main)
	if err != nil {
		return renderErrorNode{err: err}
	}
	asideHTML, err := renderNode(props.Aside)
	if err != nil {
		return renderErrorNode{err: err}
	}
	return templateNode{
		name: "components/split",
		data: struct {
			Class           string
			MainClass       string
			AsideClass      string
			Main            template.HTML
			Aside           template.HTML
			ReverseOnMobile bool
		}{
			Class:           splitClass(props),
			MainClass:       splitPaneClass("main", props.ReverseOnMobile),
			AsideClass:      splitPaneClass("aside", props.ReverseOnMobile),
			Main:            mainHTML,
			Aside:           asideHTML,
			ReverseOnMobile: props.ReverseOnMobile,
		},
	}
}

func ComponentPageHeader(props PageHeaderProps) Node {
	actionsHTML, err := renderNode(props.Actions)
	if err != nil {
		return renderErrorNode{err: err}
	}
	return templateNode{
		name: "components/page_header",
		data: struct {
			Class       string
			Title       string
			Description string
			Actions     template.HTML
		}{
			Class:       pageHeaderClass(props.Props),
			Title:       strings.TrimSpace(props.Title),
			Description: strings.TrimSpace(props.Description),
			Actions:     actionsHTML,
		},
	}
}

func ComponentContainer(props ContainerProps, children ...Node) Node {
	return layoutChildrenNode("components/container", containerClass(props), children)
}

func ComponentCard(props CardProps, children ...Node) Node {
	childHTML, err := renderNodes(children)
	if err != nil {
		return renderErrorNode{err: err}
	}
	actionsHTML, err := renderNode(props.Actions)
	if err != nil {
		return renderErrorNode{err: err}
	}
	return templateNode{
		name: "components/card",
		data: struct {
			Class       string
			Title       string
			Description string
			Actions     template.HTML
			Children    []template.HTML
		}{
			Class:       cardClass(props.Props),
			Title:       strings.TrimSpace(props.Title),
			Description: strings.TrimSpace(props.Description),
			Actions:     actionsHTML,
			Children:    childHTML,
		},
	}
}

func ComponentSection(props SectionProps, children ...Node) Node {
	childHTML, err := renderNodes(children)
	if err != nil {
		return renderErrorNode{err: err}
	}
	actionsHTML, err := renderNode(props.Actions)
	if err != nil {
		return renderErrorNode{err: err}
	}
	return templateNode{
		name: "components/section",
		data: struct {
			Class       string
			Title       string
			Description string
			Actions     template.HTML
			Children    []template.HTML
		}{
			Class:       sectionClass(props.Props),
			Title:       strings.TrimSpace(props.Title),
			Description: strings.TrimSpace(props.Description),
			Actions:     actionsHTML,
			Children:    childHTML,
		},
	}
}

func layoutChildrenNode(name, className string, children []Node) Node {
	childHTML, err := renderNodes(children)
	if err != nil {
		return renderErrorNode{err: err}
	}
	return templateNode{
		name: name,
		data: struct {
			Class    string
			Children []template.HTML
		}{
			Class:    className,
			Children: childHTML,
		},
	}
}

func loadComponentTemplates() (*template.Template, error) {
	componentTemplatesOnce.Do(func() {
		_, currentFile, _, ok := runtime.Caller(0)
		if !ok {
			cachedTemplatesErr = fmt.Errorf("failed to resolve component template path")
			return
		}
		componentsDir := filepath.Join(filepath.Dir(currentFile), "templates", "components")
		tmplFiles, err := filepath.Glob(filepath.Join(componentsDir, "*.tmpl"))
		if err != nil {
			cachedTemplatesErr = err
			return
		}
		htmlFiles, err := filepath.Glob(filepath.Join(componentsDir, "*.html"))
		if err != nil {
			cachedTemplatesErr = err
			return
		}
		files := append(tmplFiles, htmlFiles...)
		if len(files) == 0 {
			cachedTemplatesErr = fmt.Errorf("no component templates found in %s", componentsDir)
			return
		}
		cachedTemplates, cachedTemplatesErr = template.ParseFiles(files...)
	})
	return cachedTemplates, cachedTemplatesErr
}

func renderNode(node Node) (template.HTML, error) {
	if node == nil {
		return "", nil
	}
	return node.Render()
}

func renderNodes(nodes []Node) ([]template.HTML, error) {
	rendered := make([]template.HTML, 0, len(nodes))
	for _, node := range nodes {
		html, err := renderNode(node)
		if err != nil {
			return nil, err
		}
		rendered = append(rendered, html)
	}
	return rendered, nil
}

type renderErrorNode struct {
	err error
}

func (n renderErrorNode) Render() (template.HTML, error) {
	return "", n.err
}

func buttonClass(props ComponentProps) string {
	base := []string{"btn", "w-fit", buttonVariantClass(props.Variant), buttonSizeClass(props.Size)}
	if props.Class != "" {
		base = append(base, props.Class)
	}
	return joinClass(base...)
}

func inputClass(props ComponentProps) string {
	variantClass := "input-bordered"
	if props.Variant == "ghost" {
		variantClass = "input-ghost"
	}
	base := []string{"input", "w-full", variantClass, inputSizeClass(props.Size)}
	if props.Class != "" {
		base = append(base, props.Class)
	}
	return joinClass(base...)
}

func selectClass(props ComponentProps) string {
	variantClass := "select-bordered"
	if props.Variant == "ghost" {
		variantClass = "select-ghost"
	}
	base := []string{"select", "w-full", variantClass, selectSizeClass(props.Size)}
	if props.Class != "" {
		base = append(base, props.Class)
	}
	return joinClass(base...)
}

func textareaClass(props ComponentProps) string {
	variantClass := "textarea-bordered"
	if props.Variant == "ghost" {
		variantClass = "textarea-ghost"
	}
	base := []string{"textarea", "w-full", variantClass, textareaSizeClass(props.Size)}
	if props.Class != "" {
		base = append(base, props.Class)
	}
	return joinClass(base...)
}

func buttonVariantClass(variant string) string {
	switch variant {
	case "secondary":
		return "btn-secondary"
	case "ghost":
		return "btn-ghost"
	case "danger":
		return "btn-error"
	default:
		return "btn-primary"
	}
}

func buttonSizeClass(size string) string {
	switch size {
	case "sm":
		return "btn-sm"
	case "lg":
		return "btn-lg"
	default:
		return ""
	}
}

func inputSizeClass(size string) string {
	switch size {
	case "sm":
		return "input-sm"
	case "lg":
		return "input-lg"
	default:
		return ""
	}
}

func selectSizeClass(size string) string {
	switch size {
	case "sm":
		return "select-sm"
	case "lg":
		return "select-lg"
	default:
		return ""
	}
}

func textareaSizeClass(size string) string {
	switch size {
	case "sm":
		return "textarea-sm"
	case "lg":
		return "textarea-lg"
	default:
		return ""
	}
}

func checkboxClass(props ComponentProps) string {
	base := []string{"checkbox", checkboxSizeClass(props.Size)}
	if props.Class != "" {
		base = append(base, props.Class)
	}
	return joinClass(base...)
}

func checkboxSizeClass(size string) string {
	switch size {
	case "sm":
		return "checkbox-sm"
	case "lg":
		return "checkbox-lg"
	default:
		return ""
	}
}

func radioClass(props ComponentProps) string {
	base := []string{"radio", radioSizeClass(props.Size)}
	if props.Class != "" {
		base = append(base, props.Class)
	}
	return joinClass(base...)
}

func radioSizeClass(size string) string {
	switch size {
	case "sm":
		return "radio-sm"
	case "lg":
		return "radio-lg"
	default:
		return ""
	}
}

func switchClass(props ComponentProps) string {
	base := []string{"toggle", toggleSizeClass(props.Size)}
	if props.Class != "" {
		base = append(base, props.Class)
	}
	return joinClass(base...)
}

func toggleSizeClass(size string) string {
	switch size {
	case "sm":
		return "toggle-sm"
	case "lg":
		return "toggle-lg"
	default:
		return ""
	}
}

func feedbackClass(component string, props ComponentProps) string {
	base := []string{"ui-feedback", "ui-feedback-" + component, feedbackVariantClass(props.Variant), feedbackSizeClass(props.Size)}
	if props.Class != "" {
		base = append(base, props.Class)
	}
	return joinClass(base...)
}

func feedbackVariantClass(variant string) string {
	switch variant {
	case "success", "info", "warning", "error":
		return "ui-feedback-" + variant
	default:
		return "ui-feedback-info"
	}
}

func feedbackSizeClass(size string) string {
	switch size {
	case "sm", "lg":
		return "ui-feedback-" + size
	default:
		return "ui-feedback-md"
	}
}

func stackClass(props StackProps) string {
	base := []string{"flex", stackDirectionClass(props.Direction), gapClass(props.Gap), alignClass(props.Align), justifyClass(props.Justify)}
	if props.Wrap {
		base = append(base, "flex-wrap")
	}
	if props.Props.Class != "" {
		base = append(base, props.Props.Class)
	}
	return joinClass(base...)
}

func stackDirectionClass(direction string) string {
	switch strings.TrimSpace(direction) {
	case "horizontal", "row":
		return "flex-row"
	default:
		return "flex-col"
	}
}

func gridClass(props GridProps) string {
	base := []string{"grid", gapClass(props.Gap), gridColumnsClass(props.Columns, props.MinColumnWidth)}
	if props.Props.Class != "" {
		base = append(base, props.Props.Class)
	}
	return joinClass(base...)
}

func splitClass(props SplitProps) string {
	base := []string{"grid", "items-start", gapClass(props.Gap), splitColumnsClass(props.AsideWidth)}
	if props.Props.Class != "" {
		base = append(base, props.Props.Class)
	}
	return joinClass(base...)
}

func splitPaneClass(pane string, reverseOnMobile bool) string {
	base := []string{"min-w-0"}
	if !reverseOnMobile {
		return joinClass(base...)
	}
	if pane == "main" {
		base = append(base, "order-2", "lg:order-1")
	} else {
		base = append(base, "order-1", "lg:order-2")
	}
	return joinClass(base...)
}

func pageHeaderClass(props ComponentProps) string {
	return joinClass("flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between", props.Class)
}

func containerClass(props ContainerProps) string {
	base := []string{containerMaxWidthClass(props.MaxWidth), containerPaddingClass(props.Padding)}
	if props.Centered {
		base = append(base, "mx-auto")
	}
	if props.Props.Class != "" {
		base = append(base, props.Props.Class)
	}
	return joinClass(base...)
}

func cardClass(props ComponentProps) string {
	return joinClass("card bg-base-100 shadow-sm", props.Class)
}

func sectionClass(props ComponentProps) string {
	return joinClass("space-y-4", props.Class)
}

func gapClass(gap string) string {
	switch strings.TrimSpace(gap) {
	case "none", "0":
		return "gap-0"
	case "xs":
		return "gap-1"
	case "sm":
		return "gap-2"
	case "lg":
		return "gap-6"
	case "xl":
		return "gap-8"
	default:
		return "gap-4"
	}
}

func alignClass(align string) string {
	switch strings.TrimSpace(align) {
	case "start":
		return "items-start"
	case "center":
		return "items-center"
	case "end":
		return "items-end"
	default:
		return "items-stretch"
	}
}

func justifyClass(justify string) string {
	switch strings.TrimSpace(justify) {
	case "center":
		return "justify-center"
	case "end":
		return "justify-end"
	case "between":
		return "justify-between"
	default:
		return "justify-start"
	}
}

func gridColumnsClass(columns, minColumnWidth string) string {
	switch strings.TrimSpace(minColumnWidth) {
	case "sm":
		return "grid-cols-[repeat(auto-fit,minmax(14rem,1fr))]"
	case "md":
		return "grid-cols-[repeat(auto-fit,minmax(18rem,1fr))]"
	case "lg":
		return "grid-cols-[repeat(auto-fit,minmax(22rem,1fr))]"
	}

	switch strings.TrimSpace(columns) {
	case "1":
		return "grid-cols-1"
	case "2":
		return "grid-cols-1 md:grid-cols-2"
	case "4":
		return "grid-cols-1 sm:grid-cols-2 xl:grid-cols-4"
	default:
		return "grid-cols-1 md:grid-cols-2 xl:grid-cols-3"
	}
}

func splitColumnsClass(asideWidth string) string {
	switch strings.TrimSpace(asideWidth) {
	case "sm":
		return "lg:grid-cols-[minmax(0,1fr)_16rem]"
	case "lg":
		return "lg:grid-cols-[minmax(0,1fr)_28rem]"
	default:
		return "lg:grid-cols-[minmax(0,1fr)_22rem]"
	}
}

func containerMaxWidthClass(maxWidth string) string {
	switch strings.TrimSpace(maxWidth) {
	case "sm":
		return "max-w-3xl"
	case "md":
		return "max-w-5xl"
	case "full":
		return "max-w-none"
	default:
		return "max-w-7xl"
	}
}

func containerPaddingClass(padding string) string {
	switch strings.TrimSpace(padding) {
	case "none", "0":
		return "p-0"
	case "sm":
		return "p-3"
	case "lg":
		return "p-8"
	default:
		return "p-6"
	}
}

func joinClass(parts ...string) string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			continue
		}
		filtered = append(filtered, strings.TrimSpace(part))
	}
	return strings.Join(filtered, " ")
}
