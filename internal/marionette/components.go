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
	return templateNode{
		name: "components/input",
		data: struct {
			Class       string
			Name        string
			Value       string
			Placeholder string
			Disabled    bool
		}{
			Class:       inputClass(props),
			Name:        name,
			Value:       value,
			Placeholder: strings.TrimSpace(name),
			Disabled:    props.Disabled,
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

func loadComponentTemplates() (*template.Template, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to resolve component template path")
	}
	tmplDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "templates", "components", "*.tmpl")
	return template.ParseGlob(tmplDir)
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
