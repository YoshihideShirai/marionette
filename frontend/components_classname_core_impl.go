package frontend

import "strings"

// このファイルはコンポーネント用クラス名の組み立て関数を定義する。
// classNameに関するロジック追加時はここを起点に配置する。

func buttonClass(props ComponentProps) string {
	base := []string{"btn", "w-fit", buttonVariantClass(props.Variant), buttonSizeClass(props.Size)}
	if props.Class != "" {
		base = append(base, props.Class)
	}
	return joinClass(base...)
}

func linkClass(props ComponentProps, hasIcon, iconOnly bool) string {
	var base []string
	if props.Variant != "" || props.Size != "" {
		base = []string{buttonClass(props)}
		if iconOnly {
			base = append(base, "btn-square")
		}
	} else {
		base = []string{"link", "link-hover", "w-fit"}
		if hasIcon {
			base = append(base, "inline-flex", "items-center", "gap-1")
		}
		if props.Class != "" {
			base = append(base, props.Class)
		}
	}
	if props.Disabled {
		base = append(base, "pointer-events-none", "cursor-not-allowed", "opacity-50")
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
	tokens := strings.Fields(strings.TrimSpace(variant))
	if len(tokens) == 0 {
		return "btn-primary"
	}
	classes := make([]string, 0, len(tokens))
	for _, token := range tokens {
		switch token {
		case "default", "base":
			// plain DaisyUI button (no tone class)
		case "primary":
			classes = append(classes, "btn-primary")
		case "secondary":
			classes = append(classes, "btn-secondary")
		case "accent":
			classes = append(classes, "btn-accent")
		case "neutral":
			classes = append(classes, "btn-neutral")
		case "info":
			classes = append(classes, "btn-info")
		case "success":
			classes = append(classes, "btn-success")
		case "warning":
			classes = append(classes, "btn-warning")
		case "danger", "error":
			classes = append(classes, "btn-error")
		case "ghost":
			classes = append(classes, "btn-ghost")
		case "link":
			classes = append(classes, "btn-link")
		case "outline":
			classes = append(classes, "btn-outline")
		case "dash", "dashed":
			classes = append(classes, "btn-dash")
		case "soft":
			classes = append(classes, "btn-soft")
		case "glass":
			classes = append(classes, "btn-glass")
		case "active":
			classes = append(classes, "btn-active")
		case "disabled":
			classes = append(classes, "btn-disabled")
		case "wide":
			classes = append(classes, "btn-wide")
		case "block":
			classes = append(classes, "btn-block")
		case "square":
			classes = append(classes, "btn-square")
		case "circle":
			classes = append(classes, "btn-circle")
		}
	}
	if len(classes) == 0 {
		return "btn-primary"
	}
	return joinClass(classes...)
}

func buttonSizeClass(size string) string {
	switch strings.TrimSpace(size) {
	case "xs":
		return "btn-xs"
	case "sm":
		return "btn-sm"
	case "md", "":
		return ""
	case "lg":
		return "btn-lg"
	case "xl":
		return "btn-xl"
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

func badgeClass(props ComponentProps) string {
	base := []string{"badge", badgeVariantClass(props.Variant), badgeSizeClass(props.Size)}
	if props.Class != "" {
		base = append(base, props.Class)
	}
	return joinClass(base...)
}

func badgeVariantClass(variant string) string {
	switch strings.TrimSpace(variant) {
	case "primary":
		return "badge-primary"
	case "secondary":
		return "badge-secondary"
	case "accent":
		return "badge-accent"
	case "danger", "error":
		return "badge-error"
	case "outline":
		return "badge-outline"
	case "ghost":
		return "badge-ghost"
	default:
		return ""
	}
}

func badgeSizeClass(size string) string {
	switch strings.TrimSpace(size) {
	case "sm":
		return "badge-sm"
	case "lg":
		return "badge-lg"
	default:
		return ""
	}
}

func actionsClass(props ActionsProps) string {
	base := []string{"flex", "items-center", gapClass(props.Gap), actionsAlignClass(props.Align)}
	if props.Wrap {
		base = append(base, "flex-wrap")
	}
	if props.Props.Class != "" {
		base = append(base, props.Props.Class)
	}
	return joinClass(base...)
}

func actionsAlignClass(align string) string {
	switch strings.TrimSpace(align) {
	case "center":
		return "justify-center"
	case "end", "right":
		return "justify-end"
	case "between":
		return "justify-between"
	default:
		return "justify-start"
	}
}

func dividerClass(props DividerProps) string {
	base := []string{"divider", dividerSpacingClass(props.Spacing)}
	if props.Props.Class != "" {
		base = append(base, props.Props.Class)
	}
	return joinClass(base...)
}

func dividerSpacingClass(spacing string) string {
	switch strings.TrimSpace(spacing) {
	case "none":
		return "my-0"
	case "xs":
		return "my-1"
	case "sm":
		return "my-2"
	case "lg":
		return "my-6"
	default:
		return ""
	}
}

func textClass(props TextProps) string {
	base := []string{textSizeClass(props.Size), textWeightClass(props.Weight), textToneClass(props.Tone)}
	if props.Props.Class != "" {
		base = append(base, props.Props.Class)
	}
	return joinClass(base...)
}

func textSizeClass(size string) string {
	switch strings.TrimSpace(size) {
	case "xs":
		return "text-xs"
	case "sm":
		return "text-sm"
	case "lg":
		return "text-lg"
	case "xl":
		return "text-xl"
	case "2xl":
		return "text-2xl"
	case "3xl":
		return "text-3xl"
	default:
		return ""
	}
}

func textWeightClass(weight string) string {
	switch strings.TrimSpace(weight) {
	case "medium":
		return "font-medium"
	case "semibold":
		return "font-semibold"
	case "bold":
		return "font-bold"
	default:
		return ""
	}
}

func textToneClass(tone string) string {
	switch strings.TrimSpace(tone) {
	case "muted":
		return "text-base-content/60"
	case "subtle":
		return "text-base-content/70"
	default:
		return ""
	}
}

func boxClass(props BoxProps) string {
	base := []string{boxToneClass(props.Tone), boxPaddingClass(props.Padding)}
	if props.Border {
		base = append(base, "border border-base-300")
	}
	if props.Props.Class != "" {
		base = append(base, props.Props.Class)
	}
	return joinClass(base...)
}

func boxToneClass(tone string) string {
	switch strings.TrimSpace(tone) {
	case "base":
		return "bg-base-100"
	case "muted":
		return "bg-base-200"
	default:
		return ""
	}
}

func boxPaddingClass(padding string) string {
	switch strings.TrimSpace(padding) {
	case "none":
		return "p-0"
	case "sm":
		return "p-3"
	case "lg":
		return "p-6"
	default:
		return "p-4"
	}
}

func appShellClass(props ComponentProps) string {
	return joinClass("grid gap-6 lg:grid-cols-[16rem_minmax(0,1fr)]", props.Class)
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

func cardBodyClass(props CardProps) string {
	return joinClass("card-body", gapClass(props.Gap))
}

func sectionClass(props ComponentProps) string {
	return joinClass("space-y-4", props.Class)
}

func chartClass(props ComponentProps) string {
	return joinClass("card bg-base-100 shadow-sm", props.Class)
}

func imageClass(props ComponentProps) string {
	return joinClass("space-y-2", props.Class)
}

func imageFrameClass(props ImageProps) string {
	return joinClass("overflow-hidden rounded-box bg-base-200", imageAspectClass(props.AspectRatio))
}

func imageElementClass(props ImageProps) string {
	base := []string{"block", "w-full", imageObjectFitClass(props.ObjectFit)}
	if strings.TrimSpace(props.AspectRatio) != "" {
		base = append(base, "h-full")
	} else {
		base = append(base, "h-auto")
	}
	return joinClass(base...)
}

func imageAspectClass(aspectRatio string) string {
	switch strings.TrimSpace(aspectRatio) {
	case "square":
		return "aspect-square"
	case "video":
		return "aspect-video"
	case "wide":
		return "aspect-[16/9]"
	case "portrait":
		return "aspect-[3/4]"
	default:
		return ""
	}
}

func imageObjectFitClass(objectFit string) string {
	switch strings.TrimSpace(objectFit) {
	case "contain":
		return "object-contain"
	case "fill":
		return "object-fill"
	case "none":
		return "object-none"
	case "scale-down":
		return "object-scale-down"
	default:
		return "object-cover"
	}
}

func progressClass(props ComponentProps) string {
	base := []string{"progress", "w-full", progressVariantClass(props.Variant), progressSizeClass(props.Size)}
	if props.Class != "" {
		base = append(base, props.Class)
	}
	return joinClass(base...)
}

func progressVariantClass(variant string) string {
	switch strings.TrimSpace(variant) {
	case "primary":
		return "progress-primary"
	case "secondary":
		return "progress-secondary"
	case "accent":
		return "progress-accent"
	case "success":
		return "progress-success"
	case "info":
		return "progress-info"
	case "warning":
		return "progress-warning"
	case "danger", "error":
		return "progress-error"
	default:
		return ""
	}
}

func progressSizeClass(size string) string {
	switch strings.TrimSpace(size) {
	case "sm":
		return "h-1"
	case "lg":
		return "h-4"
	default:
		return "h-2"
	}
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
