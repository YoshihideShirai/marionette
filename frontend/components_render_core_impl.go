package frontend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
)

// このファイルはNode生成ロジックを定義する。
// コンポーネントの描画処理はここに追加する。

func Button(label string, props ComponentProps) Node {
	return componentButton(label, "button", props)
}

func SubmitButton(label string, props ComponentProps) Node {
	return componentButton(label, "submit", props)
}

func Link(props LinkProps) Node {
	href := strings.TrimSpace(props.Href)
	if href == "" || props.Props.Disabled {
		href = "#"
	}

	target := strings.TrimSpace(props.Target)
	if target == "" && props.External {
		target = "_blank"
	}

	rel := strings.TrimSpace(props.Rel)
	if rel == "" && (props.External || target == "_blank") {
		rel = "noopener noreferrer"
	}

	filename := strings.TrimSpace(props.Filename)
	download := props.Download || filename != ""
	label := strings.TrimSpace(props.Label)
	icon := strings.TrimSpace(props.Icon)
	ariaLabel := strings.TrimSpace(props.AriaLabel)
	if ariaLabel == "" && icon != "" && label != "" {
		ariaLabel = label
	}

	return templateNode{
		name: "components/link",
		data: struct {
			Class     string
			Label     string
			Icon      string
			Href      string
			Target    string
			Rel       string
			Download  bool
			Filename  string
			AriaLabel string
			Disabled  bool
		}{
			Class:     linkClass(props.Props, icon != "", label == ""),
			Label:     label,
			Icon:      icon,
			Href:      href,
			Target:    target,
			Rel:       rel,
			Download:  download,
			Filename:  filename,
			AriaLabel: ariaLabel,
			Disabled:  props.Props.Disabled,
		},
	}
}

func ExternalLink(label, href string, props ComponentProps) Node {
	return Link(LinkProps{
		Label:    label,
		Href:     href,
		External: true,
		Props:    props,
	})
}

func ExternalIconLink(icon, ariaLabel, href string, props ComponentProps) Node {
	return Link(LinkProps{
		Icon:      icon,
		AriaLabel: ariaLabel,
		Href:      href,
		External:  true,
		Props:     props,
	})
}

func DownloadLink(label, href, filename string, props ComponentProps) Node {
	return Link(LinkProps{
		Label:    label,
		Href:     href,
		Download: true,
		Filename: filename,
		Props:    props,
	})
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

func Table(props TableProps) Node {
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
			QueryStateName   string
			SelectedFilters  []DataFrameFilter
		}{
			Columns:          props.Columns,
			Rows:             rows,
			EmptyTitle:       strings.TrimSpace(props.EmptyTitle),
			EmptyDescription: strings.TrimSpace(props.EmptyDescription),
			QueryStateName:   strings.TrimSpace(props.QueryStateName),
			SelectedFilters:  props.SelectedFilters,
		},
	}
}

func Chart(props ChartProps) Node {
	config, err := chartConfigJSON(props)
	if err != nil {
		return renderErrorNode{err: err}
	}
	height := props.Height
	if height <= 0 {
		height = 320
	}
	ariaLabel := strings.TrimSpace(props.AriaLabel)
	if ariaLabel == "" {
		ariaLabel = strings.TrimSpace(props.Title)
	}
	if ariaLabel == "" {
		ariaLabel = "Chart"
	}
	return templateNode{
		name: "components/chart",
		data: struct {
			Class           string
			Title           string
			Description     string
			AriaLabel       string
			Height          int
			Config          template.JS
			Labels          []string
			Datasets        []ChartDataset
			Rows            []chartFallbackRow
			FallbackText    string
			QueryStateName  string
			QueryStateLabel string
		}{
			Class:           chartClass(props.Props),
			Title:           strings.TrimSpace(props.Title),
			Description:     strings.TrimSpace(props.Description),
			AriaLabel:       ariaLabel,
			Height:          height,
			Config:          template.JS(config),
			Labels:          props.Labels,
			Datasets:        props.Datasets,
			Rows:            chartFallbackRows(props),
			FallbackText:    chartFallbackText(props),
			QueryStateName:  strings.TrimSpace(props.QueryStateName),
			QueryStateLabel: strings.TrimSpace(props.QueryStateLabel),
		},
	}
}

func Image(props ImageProps) Node {
	src := strings.TrimSpace(props.Src)
	if src == "" {
		return renderErrorNode{err: fmt.Errorf("image src is required")}
	}
	loading := strings.TrimSpace(props.Loading)
	if loading == "" {
		loading = "lazy"
	}
	decoding := strings.TrimSpace(props.Decoding)
	if decoding == "" {
		decoding = "async"
	}
	return templateNode{
		name: "components/image",
		data: struct {
			Class      string
			FrameClass string
			ImageClass string
			Src        string
			Alt        string
			Caption    string
			Width      int
			Height     int
			Loading    string
			Decoding   string
		}{
			Class:      imageClass(props.Props),
			FrameClass: imageFrameClass(props),
			ImageClass: imageElementClass(props),
			Src:        src,
			Alt:        props.Alt,
			Caption:    strings.TrimSpace(props.Caption),
			Width:      props.Width,
			Height:     props.Height,
			Loading:    loading,
			Decoding:   decoding,
		},
	}
}

func Pagination(props PaginationProps) Node {
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

func chartConfigJSON(props ChartProps) (string, error) {
	chartType := strings.TrimSpace(string(props.Type))
	if chartType == "" {
		chartType = string(ChartTypeLine)
	}

	datasets := make([]map[string]any, 0, len(props.Datasets))
	for _, dataset := range props.Datasets {
		data := any(dataset.Data)
		if len(dataset.Points) > 0 {
			data = dataset.Points
		}
		item := map[string]any{
			"label": strings.TrimSpace(dataset.Label),
			"data":  data,
		}
		if color := strings.TrimSpace(dataset.BackgroundColor); color != "" {
			item["backgroundColor"] = color
		}
		if color := strings.TrimSpace(dataset.BorderColor); color != "" {
			item["borderColor"] = color
		}
		if dataset.Fill {
			item["fill"] = true
		}
		if dataset.Tension > 0 {
			item["tension"] = dataset.Tension
		}
		datasets = append(datasets, item)
	}

	options := map[string]any{
		"responsive":          true,
		"maintainAspectRatio": false,
		"plugins": map[string]any{
			"legend": map[string]any{"display": !props.Options.HideLegend},
		},
	}
	if props.Options.AspectRatio > 0 {
		options["maintainAspectRatio"] = true
		options["aspectRatio"] = props.Options.AspectRatio
	}

	if chartType != string(ChartTypePie) && chartType != string(ChartTypeDoughnut) {
		options["scales"] = chartScales(props.Options)
	}

	payload := map[string]any{
		"type": chartType,
		"data": map[string]any{
			"labels":   props.Labels,
			"datasets": datasets,
		},
		"options": options,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func chartScales(options ChartOptions) map[string]any {
	x := map[string]any{}
	y := map[string]any{
		"beginAtZero": options.BeginAtZero,
	}
	if options.Stacked {
		x["stacked"] = true
		y["stacked"] = true
	}
	if label := strings.TrimSpace(options.XAxisLabel); label != "" {
		x["title"] = map[string]any{"display": true, "text": label}
	}
	if label := strings.TrimSpace(options.YAxisLabel); label != "" {
		y["title"] = map[string]any{"display": true, "text": label}
	}
	return map[string]any{"x": x, "y": y}
}

func chartFallbackText(props ChartProps) string {
	title := strings.TrimSpace(props.Title)
	if title == "" {
		title = "Chart"
	}
	return title + " data is available in the fallback table below."
}

func chartFallbackRows(props ChartProps) []chartFallbackRow {
	rows := make([]chartFallbackRow, 0, len(props.Labels))
	for i, label := range props.Labels {
		values := make([]string, 0, len(props.Datasets))
		for _, dataset := range props.Datasets {
			if i < len(dataset.Points) {
				values = append(values, fmt.Sprintf("%g, %g", dataset.Points[i].X, dataset.Points[i].Y))
				continue
			}
			if i >= len(dataset.Data) {
				values = append(values, "")
				continue
			}
			values = append(values, fmt.Sprint(dataset.Data[i]))
		}
		rows = append(rows, chartFallbackRow{Label: label, Values: values})
	}
	return rows
}

func Tabs(props TabsProps) Node {
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

func Breadcrumb(props BreadcrumbProps) Node {
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

func checkboxComponent(props CheckboxComponentProps) Node {
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

func radioGroupComponent(props RadioGroupComponentProps) Node {
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

func switchComponent(props SwitchComponentProps) Node {
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

func Container(props ContainerProps, children ...Node) Node {
	return layoutChildrenNode("components/container", containerClass(props), children)
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
