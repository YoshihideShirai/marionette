package frontend

import (
	"fmt"
	"strings"
)

type FormRowProps struct {
	ID          string
	Label       string
	Description string
	Error       string
	Required    bool
	Control     Node
}

type FieldErrorProps struct {
	ID      string
	Message string
}

type TextFieldProps struct {
	ID          string
	Name        string
	Value       string
	Placeholder string
	Type        string
	Description string
	Error       string
	Required    bool
	Disabled    bool
	ReadOnly    bool
	Ref         string
}

type TextareaProps struct {
	ID          string
	Name        string
	Value       string
	Placeholder string
	Rows        int
	Description string
	Error       string
	Required    bool
	Disabled    bool
	ReadOnly    bool
	Ref         string
}

type SelectFieldProps struct {
	ID          string
	Name        string
	Options     []SelectOption
	Description string
	Error       string
	Required    bool
	Disabled    bool
	ReadOnly    bool
	Ref         string
}

type CheckboxProps struct {
	ID          string
	Name        string
	Value       string
	Checked     bool
	Label       string
	Description string
	Error       string
	Disabled    bool
	ReadOnly    bool
	Ref         string
}

type RadioOption struct {
	Label string
	Value string
}

type RadioGroupProps struct {
	ID          string
	Name        string
	Value       string
	Options     []RadioOption
	Description string
	Error       string
	Disabled    bool
	ReadOnly    bool
	Ref         string
}

type SwitchProps struct {
	ID          string
	Name        string
	Value       string
	Checked     bool
	Label       string
	Description string
	Error       string
	Disabled    bool
	ReadOnly    bool
	Ref         string
}

func FormRow(props FormRowProps) Node {
	id := strings.TrimSpace(props.ID)
	if id == "" {
		return renderErrorNode{err: fmt.Errorf("form row id is required")}
	}
	if props.Control == nil {
		return renderErrorNode{err: fmt.Errorf("form row control is required")}
	}
	descID := id + "-description"
	errorID := id + "-error"
	children := []Node{}
	if strings.TrimSpace(props.Label) != "" {
		labelChildren := []Node{element{Tag: "span", Attrs: map[string]string{"class": "text-sm font-medium"}, Text: props.Label}}
		if props.Required {
			labelChildren = append(labelChildren, element{Tag: "span", Attrs: map[string]string{"class": "text-error"}, Text: "*"})
		}
		children = append(children, element{Tag: "label", Attrs: map[string]string{"for": id, "class": "mb-1 inline-flex items-center gap-1"}, Children: labelChildren})
	}
	children = append(children, props.Control)
	if strings.TrimSpace(props.Description) != "" {
		children = append(children, element{Tag: "p", Attrs: map[string]string{"id": descID, "class": "text-xs text-base-content/70"}, Text: props.Description})
	}
	if strings.TrimSpace(props.Error) != "" {
		children = append(children, FieldError(FieldErrorProps{ID: errorID, Message: props.Error}))
	}
	return element{Tag: "div", Attrs: map[string]string{"class": "ui-form-row form-control gap-1.5"}, Children: children}
}

func FieldError(props FieldErrorProps) Node {
	if strings.TrimSpace(props.Message) == "" {
		return Raw("")
	}
	id := strings.TrimSpace(props.ID)
	if id == "" {
		return renderErrorNode{err: fmt.Errorf("field error id is required")}
	}
	return element{Tag: "p", Attrs: map[string]string{"id": id, "class": "ui-field-error text-xs font-medium text-error"}, Text: props.Message}
}

func TextField(props TextFieldProps) Node {
	typeValue := strings.TrimSpace(props.Type)
	if typeValue == "" {
		typeValue = "text"
	}
	return element{Tag: "input", Attrs: inputControlAttrs(controlAttrConfig{
		ID:          props.ID,
		Name:        props.Name,
		Value:       props.Value,
		Description: props.Description,
		Error:       props.Error,
		Disabled:    props.Disabled,
		ReadOnly:    props.ReadOnly,
		Required:    props.Required,
		Ref:         props.Ref,
		Class:       formInputClass(props.Error),
		Type:        typeValue,
		Placeholder: props.Placeholder,
	})}
}

func Textarea(props TextareaProps) Node {
	attrs := inputControlAttrs(controlAttrConfig{
		ID:          props.ID,
		Name:        props.Name,
		Value:       "",
		Description: props.Description,
		Error:       props.Error,
		Disabled:    props.Disabled,
		ReadOnly:    props.ReadOnly,
		Required:    props.Required,
		Ref:         props.Ref,
		Class:       formInputClass(props.Error) + " min-h-24",
		Placeholder: props.Placeholder,
	})
	if props.Rows > 0 {
		attrs["rows"] = fmt.Sprintf("%d", props.Rows)
	}
	return element{Tag: "textarea", Attrs: attrs, Text: props.Value}
}

func Select(props SelectFieldProps) Node {
	attrs := inputControlAttrs(controlAttrConfig{
		ID:          props.ID,
		Name:        props.Name,
		Description: props.Description,
		Error:       props.Error,
		Disabled:    props.Disabled,
		ReadOnly:    props.ReadOnly,
		Required:    props.Required,
		Ref:         props.Ref,
		Class:       formSelectClass(props.Error),
	})
	children := make([]Node, 0, len(props.Options))
	for _, opt := range props.Options {
		o := element{Tag: "option", Attrs: map[string]string{"value": opt.Value}, Text: opt.Label}
		if opt.Selected {
			o.Attrs["selected"] = "selected"
		}
		children = append(children, o)
	}
	return element{Tag: "select", Attrs: attrs, Children: children}
}

func Checkbox(props CheckboxProps) Node {
	attrs := checkableAttrs(props.ID, props.Name, props.Value, props.Description, props.Error, props.Ref, props.Disabled, props.ReadOnly, "checkbox")
	if props.Checked {
		attrs["checked"] = "checked"
	}
	control := element{Tag: "input", Attrs: attrs}
	content := []Node{control, element{Tag: "span", Attrs: map[string]string{"class": "text-sm"}, Text: props.Label}}
	return element{Tag: "label", Attrs: map[string]string{"for": strings.TrimSpace(props.ID), "class": "inline-flex items-center gap-2"}, Children: content}
}

func RadioGroup(props RadioGroupProps) Node {
	groupID := strings.TrimSpace(props.ID)
	if groupID == "" {
		return renderErrorNode{err: fmt.Errorf("radio group id is required")}
	}
	items := make([]Node, 0, len(props.Options))
	for i, opt := range props.Options {
		itemID := fmt.Sprintf("%s-%d", groupID, i)
		attrs := checkableAttrs(itemID, props.Name, opt.Value, props.Description, props.Error, props.Ref, props.Disabled, props.ReadOnly, "radio")
		if props.Value == opt.Value {
			attrs["checked"] = "checked"
		}
		items = append(items, element{Tag: "label", Attrs: map[string]string{"for": itemID, "class": "inline-flex items-center gap-2"}, Children: []Node{
			element{Tag: "input", Attrs: attrs},
			element{Tag: "span", Attrs: map[string]string{"class": "text-sm"}, Text: opt.Label},
		}})
	}
	return element{Tag: "div", Attrs: map[string]string{"id": groupID, "class": "flex flex-wrap gap-4"}, Children: items}
}

func Switch(props SwitchProps) Node {
	attrs := checkableAttrs(props.ID, props.Name, props.Value, props.Description, props.Error, props.Ref, props.Disabled, props.ReadOnly, "checkbox")
	attrs["class"] = formSwitchClass(props.Error)
	if props.Checked {
		attrs["checked"] = "checked"
	}
	return element{Tag: "label", Attrs: map[string]string{"for": strings.TrimSpace(props.ID), "class": "inline-flex items-center gap-2"}, Children: []Node{
		element{Tag: "input", Attrs: attrs},
		element{Tag: "span", Attrs: map[string]string{"class": "text-sm"}, Text: props.Label},
	}}
}

type controlAttrConfig struct {
	ID          string
	Name        string
	Value       string
	Description string
	Error       string
	Disabled    bool
	ReadOnly    bool
	Required    bool
	Ref         string
	Class       string
	Type        string
	Placeholder string
}

func inputControlAttrs(cfg controlAttrConfig) map[string]string {
	id := strings.TrimSpace(cfg.ID)
	if id == "" {
		return map[string]string{"data-render-error": "id is required"}
	}
	attrs := map[string]string{
		"id":               id,
		"name":             strings.TrimSpace(cfg.Name),
		"class":            cfg.Class,
		"aria-describedby": describedBy(id, cfg.Description, cfg.Error),
	}
	if cfg.Type != "" {
		attrs["type"] = cfg.Type
	}
	if cfg.Value != "" {
		attrs["value"] = cfg.Value
	}
	if cfg.Placeholder != "" {
		attrs["placeholder"] = cfg.Placeholder
	}
	if cfg.Required {
		attrs["required"] = "required"
	}
	if cfg.Disabled {
		attrs["disabled"] = "disabled"
	}
	if cfg.ReadOnly {
		attrs["readonly"] = "readonly"
	}
	if strings.TrimSpace(cfg.Error) != "" {
		attrs["aria-invalid"] = "true"
	}
	if strings.TrimSpace(cfg.Ref) != "" {
		attrs["data-ref"] = strings.TrimSpace(cfg.Ref)
	}
	return attrs
}

func checkableAttrs(id, name, value, description, errMsg, ref string, disabled, readOnly bool, typ string) map[string]string {
	attrs := map[string]string{
		"id":               strings.TrimSpace(id),
		"name":             strings.TrimSpace(name),
		"type":             typ,
		"class":            formCheckableClass(errMsg),
		"aria-describedby": describedBy(strings.TrimSpace(id), description, errMsg),
	}
	if strings.TrimSpace(value) != "" {
		attrs["value"] = strings.TrimSpace(value)
	}
	if disabled {
		attrs["disabled"] = "disabled"
	}
	if readOnly {
		attrs["readonly"] = "readonly"
	}
	if strings.TrimSpace(errMsg) != "" {
		attrs["aria-invalid"] = "true"
	}
	if strings.TrimSpace(ref) != "" {
		attrs["data-ref"] = strings.TrimSpace(ref)
	}
	return attrs
}

func describedBy(id, description, errMsg string) string {
	ids := []string{}
	if strings.TrimSpace(description) != "" {
		ids = append(ids, id+"-description")
	}
	if strings.TrimSpace(errMsg) != "" {
		ids = append(ids, id+"-error")
	}
	return strings.Join(ids, " ")
}

func formInputClass(errMsg string) string {
	base := "input input-bordered w-full text-sm focus:outline-none focus-visible:ring-2 focus-visible:ring-primary/40 disabled:opacity-60 disabled:cursor-not-allowed read-only:bg-base-200"
	if strings.TrimSpace(errMsg) != "" {
		return base + " input-error"
	}
	return base
}

func formSelectClass(errMsg string) string {
	base := "select select-bordered w-full text-sm focus:outline-none focus-visible:ring-2 focus-visible:ring-primary/40 disabled:opacity-60 disabled:cursor-not-allowed"
	if strings.TrimSpace(errMsg) != "" {
		return base + " select-error"
	}
	return base
}

func formCheckableClass(errMsg string) string {
	base := "checkbox focus-visible:ring-2 focus-visible:ring-primary/40 disabled:opacity-60 disabled:cursor-not-allowed"
	if strings.TrimSpace(errMsg) != "" {
		return base + " border-error"
	}
	return base
}

func formSwitchClass(errMsg string) string {
	base := "toggle focus-visible:ring-2 focus-visible:ring-primary/40 disabled:opacity-60 disabled:cursor-not-allowed"
	if strings.TrimSpace(errMsg) != "" {
		return base + " border-error"
	}
	return base
}
