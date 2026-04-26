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

func loadComponentTemplates() (*template.Template, error) {
	componentTemplatesOnce.Do(func() {
		_, currentFile, _, ok := runtime.Caller(0)
		if !ok {
			cachedTemplatesErr = fmt.Errorf("failed to resolve component template path")
			return
		}
		componentsDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "templates", "components")
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
