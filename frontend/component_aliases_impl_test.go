package frontend

import (
	"strings"
	"testing"

	componentspkg "github.com/YoshihideShirai/marionette/frontend/components"
	daisyui "github.com/YoshihideShirai/marionette/frontend/daisyui"
)

func TestComponentAliasesMatchCanonicalRenderOutput(t *testing.T) {
	child := TextNode("child")
	actions := TextNode("actions")
	modalBody := TextNode("modal body")
	tests := []struct {
		name      string
		alias     Node
		canonical Node
	}{
		{
			name:      "Grid aliases components.Grid",
			alias:     Grid(GridProps{Columns: "2", Gap: "sm", Props: ComponentProps{Class: "custom-grid"}}, child),
			canonical: componentspkg.Grid(GridProps{Columns: "2", Gap: "sm", Props: ComponentProps{Class: "custom-grid"}}, child),
		},
		{
			name:      "Region aliases components.Region",
			alias:     Region(RegionProps{ID: "demo-region", Props: ComponentProps{Class: "custom-region"}}, child),
			canonical: componentspkg.Region(RegionProps{ID: "demo-region", Props: ComponentProps{Class: "custom-region"}}, child),
		},
		{
			name:      "Badge aliases daisyui.Badge",
			alias:     Badge(BadgeProps{Label: "Ready", Props: ComponentProps{Class: "badge-success"}}),
			canonical: daisyui.Badge(BadgeProps{Label: "Ready", Props: ComponentProps{Class: "badge-success"}}),
		},
		{
			name: "Card aliases daisyui.Card",
			alias: Card(CardProps{
				Title:       "Card title",
				Description: "Card description",
				Actions:     actions,
				Props:       ComponentProps{Class: "custom-card"},
			}, child),
			canonical: daisyui.Card("Card title", "Card description", actions, []Node{child}, ComponentProps{Class: "custom-card"}),
		},
		{
			name: "Modal aliases daisyui.Modal",
			alias: Modal(ModalProps{
				Title:   "Modal title",
				Body:    modalBody,
				Actions: actions,
				Open:    true,
			}),
			canonical: daisyui.Modal(ModalProps{
				Title:   "Modal title",
				Body:    modalBody,
				Actions: actions,
				Open:    true,
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderFrontendNodeForTest(t, tt.alias)
			want := renderFrontendNodeForTest(t, tt.canonical)
			if got != want {
				t.Fatalf("alias render output mismatch\nwant: %q\n got: %q", want, got)
			}
		})
	}
}

func TestSimpleAliasesRenderExpectedMarkup(t *testing.T) {
	child := TextNode("child")
	tests := []struct {
		name      string
		alias     Node
		canonical Node
		want      []string
	}{
		{
			name:      "Alert aliases daisyui.Alert",
			alias:     Alert(AlertProps{Title: "Heads up", Description: "Saved", Props: ComponentProps{Class: "alert-info"}}),
			canonical: daisyui.Alert("Heads up", "Saved", ComponentProps{Class: "alert-info"}),
			want:      []string{`class="alert alert-info"`, `role="alert"`, `<span>Heads up Saved</span>`},
		},
		{
			name:      "Toast aliases daisyui.Toast",
			alias:     Toast(ToastProps{Title: "Notice", Description: "Queued", Props: ComponentProps{Class: "toast-end"}}),
			canonical: daisyui.Toast("Notice", "Queued", ComponentProps{Class: "toast-end"}),
			want:      []string{`class="toast toast-end"`, `class="alert"`, `role="alert"`, `<span>Notice Queued</span>`},
		},
		{
			name:      "Skeleton aliases daisyui.Skeleton",
			alias:     Skeleton(SkeletonProps{Rows: 2, Props: ComponentProps{Class: "w-48"}}),
			canonical: daisyui.Skeleton(2, ComponentProps{Class: "w-48"}),
			want:      []string{`class="space-y-2 w-48"`, `class="skeleton h-4 w-full"`},
		},
		{
			name:      "Progress aliases daisyui.Progress",
			alias:     Progress(ProgressProps{Value: 25, Max: 50, Label: "Loading", Props: ComponentProps{Variant: "primary", Size: "sm", Class: "custom-progress"}}),
			canonical: daisyui.Progress(ProgressProps{Value: 25, Max: 50, Label: "Loading", Props: ComponentProps{Variant: "primary", Size: "sm", Class: "custom-progress"}}),
			want:      []string{`class="progress w-full progress-primary h-1 custom-progress"`, `max="50"`, `value="25"`, `<span>Loading</span>`},
		},
		{
			name:      "EmptyState aliases daisyui.EmptyState",
			alias:     EmptyState(EmptyStateProps{Title: "Nothing here", Description: "Create one", Props: ComponentProps{Class: "min-h-64"}}),
			canonical: daisyui.EmptyState(EmptyStateProps{Title: "Nothing here", Description: "Create one", Props: ComponentProps{Class: "min-h-64"}}),
			want:      []string{`class="hero bg-base-200 rounded-box min-h-64"`, `class="hero-content text-center"`, `<h2 class="text-2xl font-bold">Nothing here</h2>`, `<p>Create one</p>`},
		},
		{
			name:      "Actions aliases daisyui.Actions",
			alias:     Actions(ActionsProps{Props: ComponentProps{Class: "justify-end"}}, child),
			canonical: daisyui.Actions(ActionsProps{Props: ComponentProps{Class: "justify-end"}}, child),
			want:      []string{`class="flex items-center gap-2 justify-end"`, `<span>child</span>`},
		},
		{
			name:      "Divider aliases daisyui.Divider",
			alias:     Divider(DividerProps{Spacing: "primary", Props: ComponentProps{Class: "my-8"}}),
			canonical: daisyui.Divider(DividerProps{Spacing: "primary", Props: ComponentProps{Class: "my-8"}}),
			want:      []string{`class="divider my-8 divider-primary"`},
		},
		{
			name:      "HiddenField aliases daisyui.HiddenField",
			alias:     HiddenField("csrf", "token-123"),
			canonical: daisyui.HiddenField("csrf", "token-123"),
			want:      []string{`name="csrf"`, `type="hidden"`, `value="token-123"`},
		},
		{
			name:      "Box aliases daisyui.Box",
			alias:     Box(BoxProps{Props: ComponentProps{Class: "bg-base-200"}}, child),
			canonical: daisyui.Box(BoxProps{Props: ComponentProps{Class: "bg-base-200"}}, child),
			want:      []string{`class="rounded-box border border-base-300 p-4 bg-base-200"`, `<span>child</span>`},
		},
		{
			name:      "Section aliases daisyui.Section",
			alias:     Section(SectionProps{Title: "Overview", Description: "Recent activity", Props: ComponentProps{Class: "pt-4"}}, child),
			canonical: daisyui.Section(SectionProps{Title: "Overview", Description: "Recent activity", Props: ComponentProps{Class: "pt-4"}}, child),
			want:      []string{`<section class="space-y-4 pt-4">`, `<h2 class="text-xl font-semibold">Overview</h2>`, `<p class="text-base-content/70">Recent activity</p>`, `<span>child</span>`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderFrontendNodeForTest(t, tt.alias)
			want := renderFrontendNodeForTest(t, tt.canonical)
			if got != want {
				t.Fatalf("alias render output mismatch\nwant: %q\n got: %q", want, got)
			}
			for _, fragment := range tt.want {
				if !strings.Contains(got, fragment) {
					t.Fatalf("expected %q in %q", fragment, got)
				}
			}
		})
	}
}
