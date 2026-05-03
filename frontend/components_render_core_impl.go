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

func UISubmitButton(label string, props ComponentProps) Node {
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

func UIThemeToggleButton(props ComponentProps) Node {
	return templateNode{
		name: "components/theme_toggle_button",
		data: struct {
			Class    string
			Disabled bool
		}{
			Class:    buttonClass(props),
			Disabled: props.Disabled,
		},
	}
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

func UIInput(name, value string, props ComponentProps) Node {
	return UIInputWithOptions(name, value, InputOptions{
		Type:        "text",
		Placeholder: strings.TrimSpace(name),
		Props:       props,
	})
}

func UIInputWithOptions(name, value string, options InputOptions) Node {
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

func UITextarea(name, value string, options TextareaOptions) Node {
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

func UIForm(props FormProps, children ...Node) Node {
	attrs := make(Attrs, len(props.Attrs)+2)
	for key, value := range props.Attrs {
		attrs[key] = value
	}
	if method := strings.TrimSpace(props.Method); method != "" {
		attrs["method"] = method
	}
	if action := strings.TrimSpace(props.Action); action != "" {
		attrs["action"] = action
	}
	return htmlElement("form", ElementProps{
		ID:    strings.TrimSpace(props.ID),
		Class: strings.TrimSpace(props.Class),
		Attrs: attrs,
	}, children...)
}

func UIActionForm(props ActionFormProps, children ...Node) Node {
	action := actionPath(strings.TrimSpace(props.Action))
	if action == "/" {
		return renderErrorNode{err: fmt.Errorf("action form action is required")}
	}

	method := strings.ToLower(strings.TrimSpace(props.Method))
	if method == "" {
		method = "post"
	}
	if method != "post" && method != "get" {
		return renderErrorNode{err: fmt.Errorf("unsupported action form method: %s", method)}
	}

	attrs := Attrs{
		"action": action,
		"method": method,
	}
	attrs["hx-"+method] = action
	if target := strings.TrimSpace(props.Target); target != "" {
		attrs["hx-target"] = target
	}
	if swap := strings.TrimSpace(props.Swap); swap != "" {
		attrs["hx-swap"] = swap
	}

	return htmlElement("form", ElementProps{
		ID:    strings.TrimSpace(props.ID),
		Class: strings.TrimSpace(props.Props.Class),
		Attrs: attrs,
	}, children...)
}

func UIFormField(control Node, props FormFieldProps) Node {
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

func UISelect(name string, options []SelectOption, props ComponentProps) Node {
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

func UIModal(props ModalProps) Node {
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

func UIToast(props ToastProps) Node {
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

func UIAlert(props AlertProps) Node {
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

func UISkeleton(props SkeletonProps) Node {
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

func UIProgress(props ProgressProps) Node {
	maxValue := props.Max
	if maxValue <= 0 {
		maxValue = 100
	}
	value := props.Value
	if value < 0 {
		value = 0
	}
	if value > maxValue {
		value = maxValue
	}
	percent := 0.0
	if maxValue > 0 {
		percent = value / maxValue * 100
	}
	ariaLabel := strings.TrimSpace(props.AriaLabel)
	if ariaLabel == "" {
		ariaLabel = strings.TrimSpace(props.Label)
	}
	if ariaLabel == "" {
		ariaLabel = "progress"
	}
	return templateNode{
		name: "components/progress",
		data: struct {
			Class         string
			Label         string
			AriaLabel     string
			Value         float64
			Max           float64
			Percent       string
			ShowValue     bool
			Indeterminate bool
		}{
			Class:         progressClass(props.Props),
			Label:         strings.TrimSpace(props.Label),
			AriaLabel:     ariaLabel,
			Value:         value,
			Max:           maxValue,
			Percent:       fmt.Sprintf("%.0f%%", percent),
			ShowValue:     props.ShowValue,
			Indeterminate: props.Indeterminate,
		},
	}
}

func UIEmptyState(props EmptyStateProps) Node {
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

func UIChart(props ChartProps) Node {
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

func UIImage(props ImageProps) Node {
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

func UIPagination(props PaginationProps) Node {
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

func UITabs(props TabsProps) Node {
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

func UIBreadcrumb(props BreadcrumbProps) Node {
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

func UICheckbox(props CheckboxComponentProps) Node {
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

func UIRadioGroup(props RadioGroupComponentProps) Node {
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

func UISwitch(props SwitchComponentProps) Node {
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

func UIBadge(props BadgeProps) Node {
	return element{Tag: "span", Attrs: map[string]string{"class": badgeClass(props.Props)}, Text: strings.TrimSpace(props.Label)}
}

func UIActions(props ActionsProps, children ...Node) Node {
	return htmlElement("div", ElementProps{
		Class: actionsClass(props),
	}, children...)
}

func UIDivider(props DividerProps) Node {
	return htmlElement("div", ElementProps{
		Class: dividerClass(props),
	})
}

func UIText(props TextProps) Node {
	return element{Tag: "span", Attrs: map[string]string{"class": textClass(props)}, Text: strings.TrimSpace(props.Text)}
}

func UIFontIcon(props FontIconProps) Node {
	name := strings.TrimSpace(props.Name)
	if name == "" {
		return renderErrorNode{err: fmt.Errorf("font icon name is required")}
	}
	library := strings.ToLower(strings.TrimSpace(props.Library))
	if library == "" {
		library = "material-icons"
	}
	tag := "i"
	className := props.Props.Class
	text := ""
	switch library {
	case "material", "material-icons":
		tag = "span"
		className = joinClass("material-icons", className)
		text = name
	case "fi", "uicons", "flaticon":
		className = joinClass("fi fi-"+name, className)
	default:
		return renderErrorNode{err: fmt.Errorf("unsupported font icon library: %q", props.Library)}
	}
	attrs := map[string]string{"class": className}
	if strings.TrimSpace(props.AriaLabel) != "" {
		attrs["aria-label"] = strings.TrimSpace(props.AriaLabel)
		attrs["role"] = "img"
	} else if props.Decorative {
		attrs["aria-hidden"] = "true"
	}
	return element{Tag: tag, Attrs: attrs, Text: text}
}

func UIHiddenField(name, value string) Node {
	return HiddenInput(name, value)
}

func UIStack(props StackProps, children ...Node) Node {
	return layoutChildrenNode("components/stack", stackClass(props), children)
}

func UIGrid(props GridProps, children ...Node) Node {
	return layoutChildrenNode("components/grid", gridClass(props), children)
}

func UISplit(props SplitProps) Node {
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

func UIPageHeader(props PageHeaderProps) Node {
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

func Container(props ContainerProps, children ...Node) Node {
	return layoutChildrenNode("components/container", containerClass(props), children)
}

func UIRegion(props RegionProps, children ...Node) Node {
	id := strings.TrimSpace(props.ID)
	if id == "" {
		return renderErrorNode{err: fmt.Errorf("region id is required")}
	}
	return htmlElement("div", ElementProps{
		ID:    id,
		Class: strings.TrimSpace(props.Props.Class),
	}, children...)
}

func UIBox(props BoxProps, children ...Node) Node {
	return htmlElement("div", ElementProps{
		Class: boxClass(props),
	}, children...)
}

func UIAppShell(props AppShellProps) Node {
	id := strings.TrimSpace(props.ID)
	if id == "" {
		id = "app"
	}
	mainID := strings.TrimSpace(props.MainID)
	if mainID == "" {
		mainID = "main-content"
	}
	return Region(RegionProps{ID: id, Props: ComponentProps{Class: appShellClass(props.Props)}},
		props.Sidebar,
		htmlElement("div", ElementProps{Class: "min-w-0 space-y-6"},
			props.Flashes,
			props.Header,
			Region(RegionProps{ID: mainID, Props: ComponentProps{Class: "space-y-6"}}, props.Content),
		),
	)
}

func UICard(props CardProps, children ...Node) Node {
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
			BodyClass   string
			Title       string
			Description string
			Actions     template.HTML
			Children    []template.HTML
		}{
			Class:       cardClass(props.Props),
			BodyClass:   cardBodyClass(props),
			Title:       strings.TrimSpace(props.Title),
			Description: strings.TrimSpace(props.Description),
			Actions:     actionsHTML,
			Children:    childHTML,
		},
	}
}

func UISection(props SectionProps, children ...Node) Node {
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
