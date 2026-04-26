package marionette

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
	"strings"
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

type InputOptions struct {
	Type        string
	Placeholder string
	Min         string
	Max         string
	Required    bool
	Error       string
	Props       ComponentProps
}

type EmptyStateProps struct {
	Title       string
	Description string
	Skeleton    bool
	Rows        int
}

type templateNode struct {
	name string
	data any
}

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
			Error       string
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
			Error:       strings.TrimSpace(options.Error),
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

func ComponentEmptyState(props EmptyStateProps) Node {
	rows := props.Rows
	if rows <= 0 {
		rows = 3
	}
	return templateNode{
		name: "components/empty_state",
		data: struct {
			Title       string
			Description string
			Skeleton    bool
			Rows        []int
		}{
			Title:       strings.TrimSpace(props.Title),
			Description: strings.TrimSpace(props.Description),
			Skeleton:    props.Skeleton,
			Rows:        make([]int, rows),
		},
	}
}

func loadComponentTemplates() (*template.Template, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to resolve component template path")
	}
	componentsDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "templates", "components")
	tmplFiles, err := filepath.Glob(filepath.Join(componentsDir, "*.tmpl"))
	if err != nil {
		return nil, err
	}
	htmlFiles, err := filepath.Glob(filepath.Join(componentsDir, "*.html"))
	if err != nil {
		return nil, err
	}
	files := append(tmplFiles, htmlFiles...)
	if len(files) == 0 {
		return nil, fmt.Errorf("no component templates found in %s", componentsDir)
	}
	return template.ParseFiles(files...)
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
