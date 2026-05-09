package frontend

import "testing"

func TestBoxClassMatrix(t *testing.T) {
	t.Run("tone", func(t *testing.T) {
		tests := []struct {
			name      string
			tone      string
			want      []string
			wantNever []string
		}{
			{name: "base", tone: "base", want: []string{"bg-base-100"}, wantNever: []string{"bg-base-200"}},
			{name: "muted", tone: "muted", want: []string{"bg-base-200"}, wantNever: []string{"bg-base-100"}},
			{name: "unknown", tone: "unknown", wantNever: []string{"bg-base-100", "bg-base-200", "unknown"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := boxClass(BoxProps{Tone: tt.tone})
				assertHasClassTokens(t, className, tt.want...)
				assertMissingClassTokens(t, className, tt.wantNever...)
			})
		}
	})

	t.Run("padding", func(t *testing.T) {
		tests := []struct {
			name      string
			padding   string
			want      []string
			wantNever []string
		}{
			{name: "none", padding: "none", want: []string{"p-0"}, wantNever: []string{"p-3", "p-4", "p-6"}},
			{name: "sm", padding: "sm", want: []string{"p-3"}, wantNever: []string{"p-0", "p-4", "p-6"}},
			{name: "lg", padding: "lg", want: []string{"p-6"}, wantNever: []string{"p-0", "p-3", "p-4"}},
			{name: "default", padding: "", want: []string{"p-4"}, wantNever: []string{"p-0", "p-3", "p-6"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := boxClass(BoxProps{Padding: tt.padding})
				assertHasClassTokens(t, className, tt.want...)
				assertMissingClassTokens(t, className, tt.wantNever...)
			})
		}
	})

	t.Run("border", func(t *testing.T) {
		tests := []struct {
			name   string
			border bool
			want   []string
			not    []string
		}{
			{name: "true", border: true, want: []string{"border", "border-base-300"}},
			{name: "false", border: false, not: []string{"border", "border-base-300"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := boxClass(BoxProps{Border: tt.border, Props: ComponentProps{Class: "custom-box"}})
				assertHasClassTokens(t, className, append(tt.want, "p-4", "custom-box")...)
				assertMissingClassTokens(t, className, tt.not...)
			})
		}
	})
}

func TestFeedbackClassMatrix(t *testing.T) {
	t.Run("components", func(t *testing.T) {
		for _, component := range []string{"alert", "toast", "banner"} {
			t.Run(component, func(t *testing.T) {
				className := feedbackClass(component, ComponentProps{Class: "custom-feedback"})
				assertHasClassTokens(t, className, "ui-feedback", "ui-feedback-"+component, "ui-feedback-info", "ui-feedback-md", "custom-feedback")
			})
		}
	})

	t.Run("variants", func(t *testing.T) {
		tests := []struct {
			name      string
			variant   string
			want      []string
			wantNever []string
		}{
			{name: "success", variant: "success", want: []string{"ui-feedback-success"}, wantNever: []string{"ui-feedback-info", "ui-feedback-warning", "ui-feedback-error"}},
			{name: "info", variant: "info", want: []string{"ui-feedback-info"}, wantNever: []string{"ui-feedback-success", "ui-feedback-warning", "ui-feedback-error"}},
			{name: "warning", variant: "warning", want: []string{"ui-feedback-warning"}, wantNever: []string{"ui-feedback-success", "ui-feedback-info", "ui-feedback-error"}},
			{name: "error", variant: "error", want: []string{"ui-feedback-error"}, wantNever: []string{"ui-feedback-success", "ui-feedback-info", "ui-feedback-warning"}},
			{name: "unknown defaults info", variant: "unknown", want: []string{"ui-feedback-info"}, wantNever: []string{"ui-feedback-unknown", "unknown"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := feedbackClass("alert", ComponentProps{Variant: tt.variant})
				assertHasClassTokens(t, className, append([]string{"ui-feedback", "ui-feedback-alert"}, tt.want...)...)
				assertMissingClassTokens(t, className, tt.wantNever...)
			})
		}
	})

	t.Run("sizes", func(t *testing.T) {
		tests := []struct {
			name      string
			size      string
			want      []string
			wantNever []string
		}{
			{name: "sm", size: "sm", want: []string{"ui-feedback-sm"}, wantNever: []string{"ui-feedback-md", "ui-feedback-lg"}},
			{name: "lg", size: "lg", want: []string{"ui-feedback-lg"}, wantNever: []string{"ui-feedback-sm", "ui-feedback-md"}},
			{name: "default", size: "", want: []string{"ui-feedback-md"}, wantNever: []string{"ui-feedback-sm", "ui-feedback-lg"}},
			{name: "unknown defaults md", size: "huge", want: []string{"ui-feedback-md"}, wantNever: []string{"ui-feedback-sm", "ui-feedback-lg", "huge"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := feedbackClass("toast", ComponentProps{Size: tt.size})
				assertHasClassTokens(t, className, append([]string{"ui-feedback", "ui-feedback-toast"}, tt.want...)...)
				assertMissingClassTokens(t, className, tt.wantNever...)
			})
		}
	})
}

func TestLayoutClassMatrix(t *testing.T) {
	t.Run("stackClass", func(t *testing.T) {
		tests := []struct {
			name  string
			props StackProps
			want  []string
		}{
			{name: "default vertical gap", props: StackProps{}, want: []string{"flex", "flex-col", "gap-4", "items-stretch", "justify-start"}},
			{name: "horizontal custom", props: StackProps{Direction: "horizontal", Gap: "lg", Align: "center", Justify: "between", Props: ComponentProps{Class: "custom-stack"}}, want: []string{"flex", "flex-row", "gap-6", "items-center", "justify-between", "custom-stack"}},
			{name: "row wrap", props: StackProps{Direction: "row", Gap: "none", Align: "end", Justify: "end", Wrap: true}, want: []string{"flex", "flex-row", "gap-0", "items-end", "justify-end", "flex-wrap"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assertHasClassTokens(t, stackClass(tt.props), tt.want...)
			})
		}
	})

	t.Run("gridClass", func(t *testing.T) {
		tests := []struct {
			name      string
			props     GridProps
			want      []string
			wantNever []string
		}{
			{name: "default columns", props: GridProps{}, want: []string{"grid", "gap-4", "grid-cols-1", "md:grid-cols-2", "xl:grid-cols-3"}},
			{name: "fixed two columns custom", props: GridProps{Columns: "2", Gap: "sm", Props: ComponentProps{Class: "custom-grid"}}, want: []string{"grid", "gap-2", "grid-cols-1", "md:grid-cols-2", "custom-grid"}, wantNever: []string{"xl:grid-cols-3"}},
			{name: "min width wins", props: GridProps{Columns: "4", MinColumnWidth: "lg", Gap: "xl"}, want: []string{"grid", "gap-8", "grid-cols-[repeat(auto-fit,minmax(22rem,1fr))]"}, wantNever: []string{"sm:grid-cols-2", "xl:grid-cols-4"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := gridClass(tt.props)
				assertHasClassTokens(t, className, tt.want...)
				assertMissingClassTokens(t, className, tt.wantNever...)
			})
		}
	})

	t.Run("splitClass", func(t *testing.T) {
		tests := []struct {
			name  string
			props SplitProps
			want  []string
		}{
			{name: "default", props: SplitProps{}, want: []string{"grid", "items-start", "gap-4", "lg:grid-cols-[minmax(0,1fr)_22rem]"}},
			{name: "small aside custom", props: SplitProps{AsideWidth: "sm", Gap: "xs", Props: ComponentProps{Class: "custom-split"}}, want: []string{"grid", "items-start", "gap-1", "lg:grid-cols-[minmax(0,1fr)_16rem]", "custom-split"}},
			{name: "large aside", props: SplitProps{AsideWidth: "lg", Gap: "lg"}, want: []string{"grid", "items-start", "gap-6", "lg:grid-cols-[minmax(0,1fr)_28rem]"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assertHasClassTokens(t, splitClass(tt.props), tt.want...)
			})
		}
	})

	t.Run("containerClass", func(t *testing.T) {
		tests := []struct {
			name      string
			props     ContainerProps
			want      []string
			wantNever []string
		}{
			{name: "default", props: ContainerProps{}, want: []string{"max-w-7xl", "p-6"}, wantNever: []string{"mx-auto"}},
			{name: "small centered custom", props: ContainerProps{MaxWidth: "sm", Padding: "sm", Centered: true, Props: ComponentProps{Class: "custom-container"}}, want: []string{"max-w-3xl", "p-3", "mx-auto", "custom-container"}},
			{name: "full no padding", props: ContainerProps{MaxWidth: "full", Padding: "none"}, want: []string{"max-w-none", "p-0"}, wantNever: []string{"mx-auto"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := containerClass(tt.props)
				assertHasClassTokens(t, className, tt.want...)
				assertMissingClassTokens(t, className, tt.wantNever...)
			})
		}
	})
}

func TestSurfaceClassMatrix(t *testing.T) {
	t.Run("cardClass", func(t *testing.T) {
		tests := []struct {
			name  string
			props ComponentProps
			want  []string
			not   []string
		}{
			{name: "base", props: ComponentProps{}, want: []string{"card", "bg-base-100", "shadow-sm"}},
			{name: "variant and size do not leak", props: ComponentProps{Variant: "primary", Size: "lg"}, want: []string{"card", "bg-base-100", "shadow-sm"}, not: []string{"primary", "lg", "card-primary", "card-lg"}},
			{name: "custom class", props: ComponentProps{Class: "custom-card"}, want: []string{"card", "bg-base-100", "shadow-sm", "custom-card"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := cardClass(tt.props)
				assertHasClassTokens(t, className, tt.want...)
				assertMissingClassTokens(t, className, tt.not...)
			})
		}
	})

	t.Run("cardBodyClass", func(t *testing.T) {
		tests := []struct {
			name string
			gap  string
			want []string
		}{
			{name: "default size", want: []string{"card-body", "gap-4"}},
			{name: "sm size", gap: "sm", want: []string{"card-body", "gap-2"}},
			{name: "lg size", gap: "lg", want: []string{"card-body", "gap-6"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assertHasClassTokens(t, cardBodyClass(CardProps{Gap: tt.gap}), tt.want...)
			})
		}
	})

	t.Run("sectionClass", func(t *testing.T) {
		className := sectionClass(ComponentProps{Variant: "muted", Size: "lg", Class: "custom-section"})
		assertHasClassTokens(t, className, "space-y-4", "custom-section")
		assertMissingClassTokens(t, className, "muted", "lg", "section-muted", "section-lg")
	})

	t.Run("imageClass", func(t *testing.T) {
		className := imageClass(ComponentProps{Variant: "rounded", Size: "sm", Class: "custom-image"})
		assertHasClassTokens(t, className, "space-y-2", "custom-image")
		assertMissingClassTokens(t, className, "rounded", "sm", "image-rounded", "image-sm")
	})
}
