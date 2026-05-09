package frontend

import (
	"strings"
	"sync"
	"testing"
)

func TestHeadingHelpersRenderExpectedTags(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want string
	}{
		{name: "h1", node: H1(Text("Title 1")), want: `<h1><span>Title 1</span></h1>`},
		{name: "h2", node: H2(Text("Title 2")), want: `<h2><span>Title 2</span></h2>`},
		{name: "h3", node: H3(Text("Title 3")), want: `<h3><span>Title 3</span></h3>`},
		{name: "h4", node: H4(Text("Title 4")), want: `<h4><span>Title 4</span></h4>`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := tt.node.Render()
			if err != nil {
				t.Fatalf("render failed: %v", err)
			}
			if got := string(html); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestFileUploadRendersFileInput(t *testing.T) {
	html, err := FileUpload("attachment", true).Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{`name="attachment"`, `type="file"`, `required`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestMenuIconRendersHamburgerSVG(t *testing.T) {
	html, err := MenuIcon(ComponentProps{Class: "text-primary"}).Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`<svg`,
		`class="inline-block h-6 w-6 stroke-current text-primary"`,
		`fill="none"`,
		`viewBox="0 0 24 24"`,
		`<path`,
		`d="M4 6h16M4 12h16M4 18h16"`,
		`stroke-linecap="round"`,
		`stroke-linejoin="round"`,
		`stroke-width="2"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestLinkRendersExpectedAttributes(t *testing.T) {
	tests := []struct {
		name    string
		node    Node
		want    []string
		notWant []string
	}{
		{
			name: "escapes label and href",
			node: Link(LinkProps{Label: `Docs <now>`, Href: `/docs?team=a&role=admin`}),
			want: []string{`class="link link-hover w-fit"`, `href="/docs?team=a&amp;role=admin"`, `<span>Docs &lt;now&gt;</span>`},
		},
		{
			name: "external defaults target and rel",
			node: ExternalLink("Docs", "https://example.com/docs", ComponentProps{}),
			want: []string{`href="https://example.com/docs"`, `target="_blank"`, `rel="noopener noreferrer"`},
		},
		{
			name: "external icon link",
			node: ExternalIconLink("↗", "Open docs", "https://example.com/docs", ComponentProps{Variant: "ghost", Size: "sm"}),
			want: []string{`class="btn w-fit btn-ghost btn-sm btn-square"`, `href="https://example.com/docs"`, `target="_blank"`, `rel="noopener noreferrer"`, `aria-label="Open docs"`, `<span class="ui-icon material-icons" aria-hidden="true">↗</span>`},
		},
		{
			name: "download filename implies download",
			node: DownloadLink("CSV", "/assets/users.csv", "users report.csv", ComponentProps{Variant: "primary", Size: "sm"}),
			want: []string{`class="btn w-fit btn-primary btn-sm"`, `href="/assets/users.csv"`, `download="users report.csv"`, `<span>CSV</span>`},
		},
		{
			name:    "disabled link is inert",
			node:    Link(LinkProps{Label: "Disabled", Href: "/danger", Props: ComponentProps{Disabled: true}}),
			want:    []string{`href="#"`, `aria-disabled="true"`, `tabindex="-1"`, `pointer-events-none`, `cursor-not-allowed`, `opacity-50`},
			notWant: []string{`href="/danger"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := tt.node.Render()
			if err != nil {
				t.Fatalf("render failed: %v", err)
			}
			got := string(html)
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Fatalf("expected %q in %q", want, got)
				}
			}
			for _, notWant := range tt.notWant {
				if strings.Contains(got, notWant) {
					t.Fatalf("did not expect %q in %q", notWant, got)
				}
			}
		})
	}
}

func TestLoadComponentTemplatesCachesParsedTemplates(t *testing.T) {
	componentTemplateSource = nil
	componentTemplateSourceErr = nil
	componentTemplateSourceOnce = sync.Once{}

	source := componentTemplateSourceForPackage()
	if componentTemplateSourceErr != nil {
		t.Fatalf("source setup failed: %v", componentTemplateSourceErr)
	}
	first, err := source.Load()
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}
	second, err := source.Load()
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}
	if first == nil || second == nil {
		t.Fatal("expected non-nil template sets")
	}
	if first != second {
		t.Fatalf("expected cached template pointer reuse, got %p and %p", first, second)
	}
}

func TestActionFormRendersHTMXActionAttributes(t *testing.T) {
	html, err := ActionForm(ActionFormProps{
		Action: "tasks/create",
		Target: "#task-list",
		Swap:   "innerHTML",
		Props:  ComponentProps{Class: "space-y-3"},
	}, TextField(TextFieldProps{ID: "task-name", Name: "name"})).Render()
	if err != nil {
		t.Fatalf("action form render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`action="tasks/create"`,
		`method="post"`,
		`hx-post="tasks/create"`,
		`hx-target="#task-list"`,
		`hx-swap="innerHTML"`,
		`class="space-y-4 space-y-3"`,
		`name="name"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestButtonClassVariantMatrix(t *testing.T) {
	tests := []struct {
		name      string
		variant   string
		want      []string
		wantNever []string
	}{
		{name: "primary", variant: "primary", want: []string{"btn", "w-fit", "btn-primary"}},
		{name: "secondary", variant: "secondary", want: []string{"btn-secondary"}},
		{name: "danger maps to error", variant: "danger", want: []string{"btn-error"}, wantNever: []string{"btn-danger"}},
		{name: "error", variant: "error", want: []string{"btn-error"}},
		{name: "outline", variant: "outline", want: []string{"btn-outline"}},
		{name: "dash", variant: "dash", want: []string{"btn-dash"}},
		{name: "dashed", variant: "dashed", want: []string{"btn-dash"}},
		{name: "soft", variant: "soft", want: []string{"btn-soft"}},
		{name: "glass", variant: "glass", want: []string{"btn-glass"}},
		{name: "active", variant: "active", want: []string{"btn-active"}},
		{name: "wide", variant: "wide", want: []string{"btn-wide"}},
		{name: "block", variant: "block", want: []string{"btn-block"}},
		{name: "square", variant: "square", want: []string{"btn-square"}},
		{name: "circle", variant: "circle", want: []string{"btn-circle"}},
		{name: "unknown falls back to primary", variant: "unknown", want: []string{"btn-primary"}, wantNever: []string{"unknown", "btn-unknown"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			className := buttonClass(ComponentProps{Variant: tt.variant})
			assertHasClassTokens(t, className, tt.want...)
			assertMissingClassTokens(t, className, tt.wantNever...)
		})
	}
}

func TestInputSelectTextareaClassVariants(t *testing.T) {
	tests := []struct {
		name      string
		className string
		want      []string
		wantNever []string
	}{
		{name: "input default is bordered", className: inputClass(ComponentProps{}), want: []string{"input", "w-full", "input-bordered"}, wantNever: []string{"input-ghost"}},
		{name: "input ghost variant", className: inputClass(ComponentProps{Variant: "ghost"}), want: []string{"input-ghost"}, wantNever: []string{"input-bordered"}},
		{name: "input sm size", className: inputClass(ComponentProps{Size: "sm"}), want: []string{"input-sm"}},
		{name: "input lg size", className: inputClass(ComponentProps{Size: "lg"}), want: []string{"input-lg"}},
		{name: "input unknown size", className: inputClass(ComponentProps{Size: "huge"}), want: []string{"input-bordered"}, wantNever: []string{"input-sm", "input-lg", "huge"}},
		{name: "select default is bordered", className: selectClass(ComponentProps{}), want: []string{"select", "w-full", "select-bordered"}, wantNever: []string{"select-ghost"}},
		{name: "select ghost variant", className: selectClass(ComponentProps{Variant: "ghost"}), want: []string{"select-ghost"}, wantNever: []string{"select-bordered"}},
		{name: "select sm size", className: selectClass(ComponentProps{Size: "sm"}), want: []string{"select-sm"}},
		{name: "select lg size", className: selectClass(ComponentProps{Size: "lg"}), want: []string{"select-lg"}},
		{name: "select unknown size", className: selectClass(ComponentProps{Size: "huge"}), want: []string{"select-bordered"}, wantNever: []string{"select-sm", "select-lg", "huge"}},
		{name: "textarea default is bordered", className: textareaClass(ComponentProps{}), want: []string{"textarea", "w-full", "textarea-bordered"}, wantNever: []string{"textarea-ghost"}},
		{name: "textarea ghost variant", className: textareaClass(ComponentProps{Variant: "ghost"}), want: []string{"textarea-ghost"}, wantNever: []string{"textarea-bordered"}},
		{name: "textarea sm size", className: textareaClass(ComponentProps{Size: "sm"}), want: []string{"textarea-sm"}},
		{name: "textarea lg size", className: textareaClass(ComponentProps{Size: "lg"}), want: []string{"textarea-lg"}},
		{name: "textarea unknown size", className: textareaClass(ComponentProps{Size: "huge"}), want: []string{"textarea-bordered"}, wantNever: []string{"textarea-sm", "textarea-lg", "huge"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertHasClassTokens(t, tt.className, tt.want...)
			assertMissingClassTokens(t, tt.className, tt.wantNever...)
		})
	}
}

func TestCheckboxRadioSwitchSizeClasses(t *testing.T) {
	tests := []struct {
		name      string
		className string
		want      []string
		wantNever []string
	}{
		{name: "checkbox sm", className: checkboxClass(ComponentProps{Size: "sm"}), want: []string{"checkbox", "checkbox-sm"}},
		{name: "checkbox lg", className: checkboxClass(ComponentProps{Size: "lg"}), want: []string{"checkbox", "checkbox-lg"}},
		{name: "checkbox unknown", className: checkboxClass(ComponentProps{Size: "huge"}), want: []string{"checkbox"}, wantNever: []string{"checkbox-sm", "checkbox-lg", "huge"}},
		{name: "radio sm", className: radioClass(ComponentProps{Size: "sm"}), want: []string{"radio", "radio-sm"}},
		{name: "radio lg", className: radioClass(ComponentProps{Size: "lg"}), want: []string{"radio", "radio-lg"}},
		{name: "radio unknown", className: radioClass(ComponentProps{Size: "huge"}), want: []string{"radio"}, wantNever: []string{"radio-sm", "radio-lg", "huge"}},
		{name: "switch sm", className: switchClass(ComponentProps{Size: "sm"}), want: []string{"toggle", "toggle-sm"}},
		{name: "switch lg", className: switchClass(ComponentProps{Size: "lg"}), want: []string{"toggle", "toggle-lg"}},
		{name: "switch unknown", className: switchClass(ComponentProps{Size: "huge"}), want: []string{"toggle"}, wantNever: []string{"toggle-sm", "toggle-lg", "huge"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertHasClassTokens(t, tt.className, tt.want...)
			assertMissingClassTokens(t, tt.className, tt.wantNever...)
		})
	}
}

func TestBadgeActionsDividerTextClassMatrices(t *testing.T) {
	t.Run("badge variants", func(t *testing.T) {
		tests := []struct {
			name      string
			variant   string
			want      []string
			wantNever []string
		}{
			{name: "primary", variant: "primary", want: []string{"badge", "badge-primary"}},
			{name: "secondary", variant: "secondary", want: []string{"badge-secondary"}},
			{name: "accent", variant: "accent", want: []string{"badge-accent"}},
			{name: "danger maps to error", variant: "danger", want: []string{"badge-error"}, wantNever: []string{"badge-danger"}},
			{name: "error", variant: "error", want: []string{"badge-error"}},
			{name: "outline", variant: "outline", want: []string{"badge-outline"}},
			{name: "ghost", variant: "ghost", want: []string{"badge-ghost"}},
			{name: "unknown", variant: "unknown", want: []string{"badge"}, wantNever: []string{"unknown", "badge-unknown"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := badgeClass(ComponentProps{Variant: tt.variant})
				assertHasClassTokens(t, className, tt.want...)
				assertMissingClassTokens(t, className, tt.wantNever...)
			})
		}
	})

	t.Run("badge sizes", func(t *testing.T) {
		tests := []struct {
			name      string
			size      string
			want      []string
			wantNever []string
		}{
			{name: "sm", size: "sm", want: []string{"badge-sm"}},
			{name: "lg", size: "lg", want: []string{"badge-lg"}},
			{name: "unknown", size: "huge", want: []string{"badge"}, wantNever: []string{"badge-sm", "badge-lg", "huge"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := badgeClass(ComponentProps{Size: tt.size})
				assertHasClassTokens(t, className, tt.want...)
				assertMissingClassTokens(t, className, tt.wantNever...)
			})
		}
	})

	t.Run("actions align", func(t *testing.T) {
		tests := []struct {
			name  string
			align string
			want  string
		}{
			{name: "center", align: "center", want: "justify-center"},
			{name: "end", align: "end", want: "justify-end"},
			{name: "right", align: "right", want: "justify-end"},
			{name: "between", align: "between", want: "justify-between"},
			{name: "unknown defaults start", align: "unknown", want: "justify-start"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := actionsClass(ActionsProps{Align: tt.align})
				assertHasClassTokens(t, className, "flex", "items-center", tt.want)
			})
		}
	})

	t.Run("divider spacing", func(t *testing.T) {
		tests := []struct {
			name      string
			spacing   string
			want      []string
			wantNever []string
		}{
			{name: "none", spacing: "none", want: []string{"divider", "my-0"}},
			{name: "xs", spacing: "xs", want: []string{"my-1"}},
			{name: "sm", spacing: "sm", want: []string{"my-2"}},
			{name: "lg", spacing: "lg", want: []string{"my-6"}},
			{name: "unknown", spacing: "huge", want: []string{"divider"}, wantNever: []string{"my-0", "my-1", "my-2", "my-6", "huge"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := dividerClass(DividerProps{Spacing: tt.spacing})
				assertHasClassTokens(t, className, tt.want...)
				assertMissingClassTokens(t, className, tt.wantNever...)
			})
		}
	})

	t.Run("text size", func(t *testing.T) {
		tests := []struct {
			name      string
			size      string
			want      []string
			wantNever []string
		}{
			{name: "xs", size: "xs", want: []string{"text-xs"}},
			{name: "sm", size: "sm", want: []string{"text-sm"}},
			{name: "lg", size: "lg", want: []string{"text-lg"}},
			{name: "xl", size: "xl", want: []string{"text-xl"}},
			{name: "2xl", size: "2xl", want: []string{"text-2xl"}},
			{name: "3xl", size: "3xl", want: []string{"text-3xl"}},
			{name: "unknown", size: "huge", wantNever: []string{"huge", "text-huge"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := textClass(TextProps{Size: tt.size})
				assertHasClassTokens(t, className, tt.want...)
				assertMissingClassTokens(t, className, tt.wantNever...)
			})
		}
	})

	t.Run("text tone", func(t *testing.T) {
		tests := []struct {
			name      string
			tone      string
			want      []string
			wantNever []string
		}{
			{name: "muted", tone: "muted", want: []string{"text-base-content/60"}},
			{name: "subtle", tone: "subtle", want: []string{"text-base-content/70"}},
			{name: "unknown", tone: "loud", wantNever: []string{"loud", "text-loud", "text-base-content/60", "text-base-content/70"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				className := textClass(TextProps{Tone: tt.tone})
				assertHasClassTokens(t, className, tt.want...)
				assertMissingClassTokens(t, className, tt.wantNever...)
			})
		}
	})
}

func TestLinkClassSecurityStates(t *testing.T) {
	className := linkClass(ComponentProps{Disabled: true}, false, false)
	assertHasClassTokens(t, className, "link", "link-hover", "w-fit", "pointer-events-none", "cursor-not-allowed", "opacity-50")
}

func renderFrontendNodeForTest(t *testing.T, node Node) string {
	t.Helper()
	html, err := node.Render()
	if err != nil {
		t.Fatalf("render node: %v", err)
	}
	return string(html)
}

func assertHasClassTokens(t *testing.T, className string, tokens ...string) {
	t.Helper()
	classes := classTokenSet(className)
	for _, token := range tokens {
		if !classes[token] {
			t.Fatalf("expected class %q in %q", token, className)
		}
	}
}

func assertMissingClassTokens(t *testing.T, className string, tokens ...string) {
	t.Helper()
	classes := classTokenSet(className)
	for _, token := range tokens {
		if classes[token] {
			t.Fatalf("did not expect class %q in %q", token, className)
		}
	}
}

func classTokenSet(className string) map[string]bool {
	set := map[string]bool{}
	for _, token := range strings.Fields(className) {
		set[token] = true
	}
	return set
}

func TestStreamTriggerRendersHTMXPollingAttributes(t *testing.T) {
	html, err := StreamTrigger(StreamTriggerProps{
		ID:     "reply-stream",
		Action: "/chat/stream",
		Target: "#chat-panel",
		Delay:  "250ms",
		Props:  ComponentProps{Class: "custom-stream"},
	}).Render()
	if err != nil {
		t.Fatalf("stream trigger render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`id="reply-stream"`,
		`class="hidden custom-stream"`,
		`aria-hidden="true"`,
		`hx-post="/chat/stream"`,
		`hx-trigger="load delay:250ms"`,
		`hx-target="#chat-panel"`,
		`hx-swap="outerHTML"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}
