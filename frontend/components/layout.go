package components

import (
	"bytes"
	"html/template"
	"strings"

	lowhtml "github.com/YoshihideShirai/marionette/frontend/html"
	shared "github.com/YoshihideShirai/marionette/frontend/shared"
)

var gridTmpl = template.Must(template.New("components/grid").Parse(`{{define "components/grid" -}}
<div class="{{.Class}}">
  {{range .Children}}{{.}}{{end}}
</div>
{{- end}}`))

type gridNode struct {
	className string
	children  []shared.Node
}

func Grid(props shared.GridProps, children ...shared.Node) shared.Node {
	return gridNode{className: gridClass(props), children: children}
}

func (n gridNode) Render() (template.HTML, error) {
	children := make([]template.HTML, 0, len(n.children))
	for _, child := range n.children {
		if child == nil {
			children = append(children, "")
			continue
		}
		rendered, err := child.Render()
		if err != nil {
			return "", err
		}
		children = append(children, rendered)
	}
	var b bytes.Buffer
	err := gridTmpl.ExecuteTemplate(&b, "components/grid", struct {
		Class    string
		Children []template.HTML
	}{Class: n.className, Children: children})
	return template.HTML(b.String()), err
}

func Region(props shared.RegionProps, children ...shared.Node) shared.Node {
	attrs := map[string]string{"id": strings.TrimSpace(props.ID)}
	if props.Props.Class != "" {
		attrs["class"] = strings.TrimSpace(props.Props.Class)
	}
	return lowhtml.ElementNode{Tag: "div", Attrs: attrs, Children: children}
}

func gridClass(props shared.GridProps) string {
	base := []string{"grid", gapClass(props.Gap), gridColumnsClass(props.Columns, props.MinColumnWidth)}
	if props.Props.Class != "" {
		base = append(base, props.Props.Class)
	}
	return joinClass(base...)
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

func joinClass(parts ...string) string {
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return strings.Join(out, " ")
}
