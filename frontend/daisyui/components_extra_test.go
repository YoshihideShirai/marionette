package daisyui

import (
	"strings"
	"testing"

	lowhtml "github.com/YoshihideShirai/marionette/frontend/html"
	shared "github.com/YoshihideShirai/marionette/frontend/shared"
)

func TestSmallExtraComponentsRenderExpectedMarkup(t *testing.T) {
	tests := []struct {
		name string
		node shared.Node
		want []string
	}{
		{
			name: "kbd",
			node: Kbd("⌘K"),
			want: []string{`<kbd class="kbd">⌘K</kbd>`},
		},
		{
			name: "code",
			node: Code("go test"),
			want: []string{`<code class="bg-base-200 rounded px-1 py-0.5">go test</code>`},
		},
		{
			name: "indicator",
			node: Indicator(TextNode("3"), TextNode("Inbox")),
			want: []string{`<div class="indicator">`, `class="indicator-item"`, `>3</span>`, `>Inbox</span>`},
		},
		{
			name: "link",
			node: Link("Docs", "/docs", shared.ComponentProps{Class: "link-primary"}),
			want: []string{`<a class="link link-primary" href="/docs">Docs</a>`},
		},
		{
			name: "tooltip",
			node: Tooltip("Helpful", TextNode("?")),
			want: []string{`<div class="tooltip" data-tip="Helpful">`, `>?</span>`},
		},
		{
			name: "loading",
			node: Loading("loading-lg"),
			want: []string{`<span class="loading loading-spinner loading-lg"></span>`},
		},
		{
			name: "radial progress",
			node: RadialProgress(73, "text-primary"),
			want: []string{`aria-valuenow="73"`, `class="radial-progress text-primary"`, `role="progressbar"`, `style="--value:73;"`, `>73%</div>`},
		},
		{
			name: "range",
			node: Range("volume", 7, 0, 10),
			want: []string{`class="range"`, `max="10"`, `min="0"`, `name="volume"`, `type="range"`, `value="7"`},
		},
		{
			name: "toggle",
			node: Toggle("enabled", true),
			want: []string{`checked="checked"`, `class="toggle"`, `name="enabled"`, `type="checkbox"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderComponentForTest(t, tt.node)
			assertContainsAll(t, got, tt.want)
		})
	}
}

func TestCompositeExtraComponentsRenderExpectedContainers(t *testing.T) {
	tests := []struct {
		name string
		node shared.Node
		want []string
	}{
		{name: "join", node: Join(TextNode("Joined child")), want: []string{`<div class="join">`, `>Joined child</span>`}},
		{name: "mask", node: Mask("mask-squircle", TextNode("Masked child")), want: []string{`<div class="mask mask-squircle">`, `>Masked child</span>`}},
		{name: "carousel", node: Carousel(TextNode("Slide child")), want: []string{`<div class="carousel w-full">`, `>Slide child</span>`}},
		{name: "carousel item", node: CarouselItem("slide-1", TextNode("Slide 1")), want: []string{`class="carousel-item w-full"`, `id="slide-1"`, `>Slide 1</span>`}},
		{name: "chat bubble", node: ChatBubble(TextNode("Hello"), true), want: []string{`<div class="chat chat-end">`, `class="chat-bubble"`, `>Hello</span>`}},
		{name: "dock", node: Dock(TextNode("Dock child")), want: []string{`<div class="dock">`, `>Dock child</span>`}},
		{name: "fieldset", node: Fieldset("Profile", TextNode("Field child")), want: []string{`<fieldset class="fieldset">`, `<legend class="fieldset-legend">Profile</legend>`, `>Field child</span>`}},
		{name: "browser mockup", node: BrowserMockup(TextNode("Browser child")), want: []string{`<div class="mockup-browser">`, `class="mockup-browser-toolbar"`, `https://example.com`, `>Browser child</span>`}},
		{name: "phone mockup", node: PhoneMockup(TextNode("Phone child")), want: []string{`<div class="mockup-phone">`, `class="mockup-phone-camera"`, `class="mockup-phone-display"`, `>Phone child</span>`}},
		{name: "code mockup", node: CodeMockup("npm test"), want: []string{`<div class="mockup-code">`, `<pre data-prefix="$"><code>npm test</code></pre>`}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderComponentForTest(t, tt.node)
			assertContainsAll(t, got, tt.want)
		})
	}
}

func TestDropdownAddsDropdownContentAttrs(t *testing.T) {
	t.Run("wraps non element content", func(t *testing.T) {
		got := renderComponentForTest(t, Dropdown(TextNode("Open"), lowhtml.Raw("Menu")))
		assertContainsAll(t, got, []string{`<div class="dropdown">`, `role="button"`, `tabindex="0"`, `<div class="dropdown-content" tabindex="-1">Menu</div>`})
	})

	t.Run("appends class to element content", func(t *testing.T) {
		menu := lowhtml.ElementNode{Tag: "ul", Attrs: map[string]string{"class": "menu shadow", "id": "actions"}, Children: []shared.Node{TextNode("Edit")}}
		got := renderComponentForTest(t, Dropdown(TextNode("Open"), menu))
		assertContainsAll(t, got, []string{`<ul class="menu shadow dropdown-content" id="actions" tabindex="-1">`, `>Edit</span>`})
	})
}

func TestRatingAndRangeInputs(t *testing.T) {
	rating := renderComponentForTest(t, Rating("score", 5, 3))
	assertContainsAll(t, rating, []string{`<div class="rating">`, `name="score"`, `type="radio"`, `value="3"`})
	if got := strings.Count(rating, `type="radio"`); got != 5 {
		t.Fatalf("expected 5 radio inputs, got %d in %q", got, rating)
	}
	if got := strings.Count(rating, `checked="checked"`); got != 1 {
		t.Fatalf("expected 1 checked input, got %d in %q", got, rating)
	}
	assertContainsAll(t, rating, []string{`<input checked="checked" class="mask mask-star-2 bg-orange-400" name="score" type="radio" value="3"></input>`})

	rangeHTML := renderComponentForTest(t, Range("score", 8, 1, 10))
	assertContainsAll(t, rangeHTML, []string{`class="range"`, `max="10"`, `min="1"`, `name="score"`, `type="range"`, `value="8"`})
}

func TestVariantComponentsRenderExpectedClasses(t *testing.T) {
	tests := []struct {
		name       string
		node       shared.Node
		want       []string
		wantCounts map[string]int
	}{
		{
			name: "progress color and size",
			node: ProgressWithVariant(45, 90, "Uploading", "success", shared.ComponentProps{Size: "lg", Class: "custom-progress"}),
			want: []string{`<progress class="progress w-full  h-4 progress-success custom-progress" max="90" value="45">`, `<span>Uploading</span>`},
		},
		{
			name: "badge color size style",
			node: BadgeWithVariant("New", "primary", "lg", "outline", shared.ComponentProps{Class: "tracking-wide"}),
			want: []string{`<span class="badge badge badge-primary badge-lg badge-outline tracking-wide">New</span>`},
		},
		{
			name: "range color size and value",
			node: RangeWithVariants("volume", 4, 0, 11, "secondary", "xs"),
			want: []string{`class="range range-secondary range-xs"`, `max="11"`, `min="0"`, `name="volume"`, `type="range"`, `value="4"`},
		},
		{
			name:       "rating size half allow clear and checked",
			node:       RatingWithVariants("rating", 4, 2, "sm", true, true),
			want:       []string{`<div class="rating rating-sm rating-half">`, `class="rating-hidden"`, `value="0"`, `value="2"`},
			wantCounts: map[string]int{`type="radio"`: 5, `checked="checked"`: 1},
		},
		{
			name: "toast placement and custom class",
			node: ToastWithPlacement([]shared.Node{TextNode("Saved")}, "end", "bottom", "z-50"),
			want: []string{`<div class="toast toast-end toast-bottom z-50">`, `>Saved</span>`},
		},
		{
			name: "tooltip placement color and open",
			node: TooltipWithVariants("Copied", TextNode("Copy"), "left", "info", true),
			want: []string{`<div class="tooltip tooltip-left tooltip-info tooltip-open" data-tip="Copied">`, `>Copy</span>`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderComponentForTest(t, tt.node)
			assertContainsAll(t, got, tt.want)
			for needle, wantCount := range tt.wantCounts {
				if gotCount := strings.Count(got, needle); gotCount != wantCount {
					t.Fatalf("expected %q count %d, got %d in %q", needle, wantCount, gotCount, got)
				}
			}
		})
	}
}

func TestToggleVariantRendersDaisyUIColorClasses(t *testing.T) {
	tests := []struct {
		name    string
		variant string
		want    string
	}{
		{name: "primary", variant: "primary", want: `class="toggle toggle-primary"`},
		{name: "secondary", variant: "secondary", want: `class="toggle toggle-secondary"`},
		{name: "accent", variant: "accent", want: `class="toggle toggle-accent"`},
		{name: "neutral", variant: "neutral", want: `class="toggle toggle-neutral"`},
		{name: "info", variant: "info", want: `class="toggle toggle-info"`},
		{name: "success", variant: "success", want: `class="toggle toggle-success"`},
		{name: "warning", variant: "warning", want: `class="toggle toggle-warning"`},
		{name: "error", variant: "error", want: `class="toggle toggle-error"`},
		{name: "danger alias", variant: "danger", want: `class="toggle toggle-error"`},
		{name: "normalized", variant: " Primary ", want: `class="toggle toggle-primary"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := ToggleVariant("demo", true, tt.variant).Render()
			if err != nil {
				t.Fatalf("render ToggleVariant: %v", err)
			}
			got := string(html)
			if !strings.Contains(got, tt.want) {
				t.Fatalf("expected %q in %q", tt.want, got)
			}
			if !strings.Contains(got, `checked="checked"`) {
				t.Fatalf("expected checked attribute in %q", got)
			}
		})
	}
}

func TestDrawerWithPropsRendersConfigurableShell(t *testing.T) {
	html, err := DrawerWithProps(DrawerProps{
		ID:           "demo-drawer",
		Class:        "lg:drawer-open app-shell",
		ContentClass: "flex min-h-screen",
		SideClass:    "z-40",
		Content:      TextNode("content"),
		Side:         TextNode("side"),
	}).Render()
	if err != nil {
		t.Fatalf("render DrawerWithProps: %v", err)
	}
	got := string(html)
	for _, want := range []string{`class="drawer lg:drawer-open app-shell"`, `id="demo-drawer"`, `class="drawer-content flex min-h-screen"`, `class="drawer-side z-40"`, `aria-label="close sidebar"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestDashwindPrimitiveComponentsRenderDaisyUIMarkup(t *testing.T) {
	table := TableWithProps(TableProps{
		Headers: []string{"Name", "Status"},
		Rows:    [][]shared.Node{{TextNode("Acme"), Badge(shared.BadgeProps{Label: "Paid", Props: shared.ComponentProps{Class: "badge-success"}})}},
		Class:   "table-zebra",
	})
	card := CardPanel(CardPanelProps{Title: "Recent", Description: "Transactions", Class: "bg-base-100 shadow"}, table)
	html, err := card.Render()
	if err != nil {
		t.Fatalf("render card/table: %v", err)
	}
	got := string(html)
	for _, want := range []string{`class="card bg-base-100 shadow"`, `class="card-body"`, `class="table table-zebra"`, `<th>Name</th>`, `class="badge badge-success"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestThemeToggleButtonIncludesSystemModeState(t *testing.T) {
	html, err := ThemeToggleButton(shared.ComponentProps{Class: "btn-sm"}).Render()
	if err != nil {
		t.Fatalf("render ThemeToggleButton: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`aria-label="Toggle theme: system, light, dark"`,
		`data-mrn-theme-toggle="true"`,
		`data-mrn-theme-mode="system"`,
		`data-mrn-theme-icon="true"`,
		`◐`,
		`globalThis.mrnToggleTheme()`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
	for _, unwanted := range []string{`data-mrn-theme-label`, `System`, `Light`, `Dark`, `Theme:`} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("expected no visible theme text %q in %q", unwanted, got)
		}
	}
}

func renderComponentForTest(t *testing.T, node shared.Node) string {
	t.Helper()
	html, err := node.Render()
	if err != nil {
		t.Fatalf("render component: %v", err)
	}
	return string(html)
}

func assertContainsAll(t *testing.T, got string, wants []string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func assertContainsNone(t *testing.T, got string, unwants []string) {
	t.Helper()
	for _, unwant := range unwants {
		if strings.Contains(got, unwant) {
			t.Fatalf("expected %q to be absent from %q", unwant, got)
		}
	}
}

func TestCoreFormControlPropsRenderBoundaryAttributes(t *testing.T) {
	tests := []struct {
		name      string
		node      shared.Node
		want      []string
		wantNever []string
	}{
		{
			name: "button disabled emits disabled attribute",
			node: Button("Save", shared.ComponentProps{Class: "btn-primary", Disabled: true}),
			want: []string{`<button class="btn btn-primary" disabled="disabled">Save</button>`},
		},
		{
			name:      "button enabled omits disabled attribute",
			node:      Button("Save", shared.ComponentProps{}),
			want:      []string{`<button class="btn">Save</button>`},
			wantNever: []string{`disabled="disabled"`},
		},
		{
			name: "input disabled carries name value and class",
			node: Input("email", "alice@example.test", shared.ComponentProps{Disabled: true}),
			want: []string{`class="input w-full"`, `disabled="disabled"`, `name="email"`, `value="alice@example.test"`},
		},
		{
			name: "input size prop renders xl class",
			node: Input("email", "alice@example.test", shared.ComponentProps{Size: "xl"}),
			want: []string{`class="input w-full input-xl"`},
		},
		{
			name: "select marks only selected option",
			node: Select("status", []shared.SelectOption{{Label: "Draft", Value: "draft"}, {Label: "Published", Value: "published", Selected: true}}, shared.ComponentProps{Class: "select-sm"}),
			want: []string{`<select class="select select-sm" name="status">`, `<option value="draft">Draft</option>`, `<option selected="selected" value="published">Published</option>`},
		},
		{
			name: "select size prop renders md class",
			node: Select("status", []shared.SelectOption{{Label: "Draft", Value: "draft"}}, shared.ComponentProps{Size: "md"}),
			want: []string{`<select class="select select-md" name="status">`},
		},
		{
			name: "textarea positive rows and required are rendered",
			node: Textarea("bio", "hello", shared.TextareaOptions{Placeholder: "About you", Rows: 4, Required: true, Props: shared.ComponentProps{Class: "textarea-primary"}}),
			want: []string{`class="textarea w-full textarea-primary"`, `name="bio"`, `placeholder="About you"`, `required="required"`, `rows="4"`, `>hello</textarea>`},
		},
		{
			name: "textarea size prop renders xs class",
			node: Textarea("bio", "hello", shared.TextareaOptions{Props: shared.ComponentProps{Size: "xs"}}),
			want: []string{`class="textarea w-full textarea-xs"`},
		},
		{
			name:      "textarea default rows omits rows attribute",
			node:      Textarea("bio", "hello", shared.TextareaOptions{Rows: 0}),
			want:      []string{`<textarea class="textarea w-full" name="bio">hello</textarea>`},
			wantNever: []string{`rows=`},
		},
		{
			name: "checkbox checked emits type and checked",
			node: Checkbox(shared.CheckboxComponentProps{Name: "agree", Value: "yes", Label: "Agree", Checked: true, Props: shared.ComponentProps{Class: "checkbox-success"}}),
			want: []string{`<input checked="checked" class="checkbox checkbox-success" name="agree" type="checkbox" value="yes"></input>`, `<span class="label">Agree</span>`},
		},
		{
			name: "checkbox size prop renders md class",
			node: Checkbox(shared.CheckboxComponentProps{Name: "agree", Value: "yes", Label: "Agree", Props: shared.ComponentProps{Size: "md"}}),
			want: []string{`class="checkbox checkbox-md"`},
		},
		{
			name: "checkbox disabled prop renders disabled attribute",
			node: Checkbox(shared.CheckboxComponentProps{Name: "agree", Value: "yes", Label: "Agree", Props: shared.ComponentProps{Disabled: true}}),
			want: []string{`class="checkbox" disabled="disabled" name="agree" type="checkbox" value="yes"`},
		},
		{
			name:      "checkbox unchecked omits checked",
			node:      Checkbox(shared.CheckboxComponentProps{Name: "agree", Value: "yes", Label: "Agree"}),
			want:      []string{`class="checkbox"`, `name="agree"`, `type="checkbox"`, `value="yes"`},
			wantNever: []string{`checked="checked"`},
		},
		{
			name: "radio group checked item emits type and checked",
			node: RadioGroup(shared.RadioGroupComponentProps{Name: "plan", Items: []shared.RadioItem{{Label: "Free", Value: "free"}, {Label: "Pro", Value: "pro", Checked: true}}, Props: shared.ComponentProps{Class: "rounded"}}),
			want: []string{`<div class="space-y-2 rounded">`, `class="radio" name="plan" type="radio" value="free"`, `checked="checked" class="radio" name="plan" type="radio" value="pro"`},
		},
		{
			name: "radio group size prop renders xl class",
			node: RadioGroup(shared.RadioGroupComponentProps{Name: "plan", Items: []shared.RadioItem{{Label: "Free", Value: "free"}}, Props: shared.ComponentProps{Size: "xl"}}),
			want: []string{`class="radio radio-xl" name="plan" type="radio" value="free"`},
		},
		{
			name: "radio group disabled prop disables all items",
			node: RadioGroup(shared.RadioGroupComponentProps{Name: "plan", Items: []shared.RadioItem{{Label: "Free", Value: "free"}, {Label: "Pro", Value: "pro"}}, Props: shared.ComponentProps{Disabled: true}}),
			want: []string{`class="radio" disabled="disabled" name="plan" type="radio" value="free"`, `class="radio" disabled="disabled" name="plan" type="radio" value="pro"`},
		},
		{
			name: "radio group item disabled prop disables one item",
			node: RadioGroup(shared.RadioGroupComponentProps{Name: "plan", Items: []shared.RadioItem{{Label: "Free", Value: "free"}, {Label: "Pro", Value: "pro", Disabled: true}}}),
			want: []string{`class="radio" disabled="disabled" name="plan" type="radio" value="pro"`},
		},
		{
			name: "switch checked emits checkbox type and checked",
			node: Switch(shared.SwitchComponentProps{Name: "enabled", Value: "1", Label: "Enabled", Checked: true, Props: shared.ComponentProps{Class: "toggle-primary"}}),
			want: []string{`<input checked="checked" class="toggle toggle-primary" name="enabled" type="checkbox" value="1"></input>`, `<span class="label">Enabled</span>`},
		},
		{
			name: "switch size prop renders xs class",
			node: Switch(shared.SwitchComponentProps{Name: "enabled", Value: "1", Label: "Enabled", Props: shared.ComponentProps{Size: "xs"}}),
			want: []string{`class="toggle toggle-xs"`},
		},
		{
			name: "switch disabled prop renders disabled attribute",
			node: Switch(shared.SwitchComponentProps{Name: "enabled", Value: "1", Label: "Enabled", Props: shared.ComponentProps{Disabled: true}}),
			want: []string{`class="toggle" disabled="disabled" name="enabled" type="checkbox" value="1"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderComponentForTest(t, tt.node)
			assertContainsAll(t, got, tt.want)
			assertContainsNone(t, got, tt.wantNever)
		})
	}
}

func TestCoreNavigationAndFeedbackPropsRenderBoundaryAttributes(t *testing.T) {
	tests := []struct {
		name       string
		node       shared.Node
		want       []string
		wantNever  []string
		wantCounts map[string]int
	}{
		{
			name: "progress defaults max and includes value when determinate",
			node: Progress(shared.ProgressProps{Value: 25, Max: 0, Label: "Loading", Props: shared.ComponentProps{Class: "progress-info"}}),
			want: []string{`<progress class="progress w-full  h-2 progress-info" max="100" value="25">`, `<span>Loading</span>`},
		},
		{
			name:      "progress indeterminate omits value but keeps explicit max",
			node:      Progress(shared.ProgressProps{Value: 25, Max: 80, Indeterminate: true}),
			want:      []string{`<progress class="progress w-full  h-2" max="80">`},
			wantNever: []string{`value="25"`},
		},
		{
			name: "tabs render tablist role aria label and active selected state",
			node: Tabs(shared.TabsProps{AriaLabel: "Sections", Items: []shared.TabsItem{{Label: "One", Active: true}, {Label: "Two", Disabled: true}}, Props: shared.ComponentProps{Class: "tabs-boxed"}}),
			want: []string{`<div aria-label="Sections" class="tabs tabs-boxed" role="tablist">`, `<button aria-selected="true" class="tab tab-active" role="tab" type="button">One</button>`, `<button aria-selected="false" class="tab tab-disabled" disabled="disabled" role="tab" type="button">Two</button>`},
		},
		{
			name: "pagination renders hrefs and current page button",
			node: Pagination(shared.PaginationProps{Page: 2, TotalPages: 5, PrevHref: "/p/1", NextHref: "/p/3"}),
			want: []string{`<div class="join">`, `<a class="join-item btn" href="/p/1">«</a>`, `<button class="join-item btn">2</button>`, `<a class="join-item btn" href="/p/3">»</a>`},
		},
		{
			name:       "empty state skeleton defaults rows and accessibility attributes",
			node:       EmptyState(shared.EmptyStateProps{Skeleton: true, Rows: 0}),
			want:       []string{`<div aria-busy="true" aria-live="polite" class="space-y-2">`, `class="skeleton h-4 w-full"`},
			wantCounts: map[string]int{`class="skeleton h-4 w-full"`: 3},
		},
		{
			name:      "empty state content uses props class and text",
			node:      EmptyState(shared.EmptyStateProps{Title: "No records", Description: "Try another filter", Props: shared.ComponentProps{Class: "min-h-64"}}),
			want:      []string{`<div class="hero bg-base-200 rounded-box min-h-64">`, `<h2 class="text-2xl font-bold">No records</h2>`, `<p>Try another filter</p>`},
			wantNever: []string{`aria-busy="true"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderComponentForTest(t, tt.node)
			assertContainsAll(t, got, tt.want)
			assertContainsNone(t, got, tt.wantNever)
			for needle, wantCount := range tt.wantCounts {
				if gotCount := strings.Count(got, needle); gotCount != wantCount {
					t.Fatalf("expected %q count %d, got %d in %q", needle, wantCount, gotCount, got)
				}
			}
		})
	}
}

func TestTableWithPropsClassIndexBoundsDoNotPanic(t *testing.T) {
	got := renderComponentForTest(t, TableWithProps(TableProps{
		Headers:         []string{"Name", "Status"},
		Rows:            [][]shared.Node{{TextNode("Acme"), TextNode("Paid")}, {TextNode("Beta"), TextNode("Due")}},
		HeaderClasses:   []string{"first-col"},
		HeaderSortables: []bool{false, true},
		CellClasses:     [][]string{{"name-cell"}},
	}))

	assertContainsAll(t, got, []string{
		`<th class="first-col">Name</th>`,
		`<th aria-sort="none"><span class="inline-flex items-center gap-1"><span>Status</span><span aria-hidden="true" class="text-base-content/40">↕</span></span></th>`,
		`<td class="name-cell"><span>Acme</span></td>`,
		`<td><span>Paid</span></td>`,
		`<td><span>Beta</span></td>`,
	})
}

func TestButtonWithAttrsMergePriority(t *testing.T) {
	tests := []struct {
		name  string
		props shared.ComponentProps
		attrs map[string]string
		want  []string
	}{
		{
			name:  "caller attrs preserved and class appended after props class",
			props: shared.ComponentProps{Class: "btn-primary"},
			attrs: map[string]string{"class": "w-full", "type": "submit", "data-action": "save"},
			want:  []string{`class="btn btn-primary w-full"`, `data-action="save"`, `type="submit"`},
		},
		{
			name:  "component disabled prop overrides caller disabled value",
			props: shared.ComponentProps{Disabled: true},
			attrs: map[string]string{"disabled": "false", "class": "btn-ghost"},
			want:  []string{`class="btn  btn-ghost"`, `disabled="disabled"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderComponentForTest(t, ButtonWithAttrs("Run", tt.props, tt.attrs))
			assertContainsAll(t, got, tt.want)
		})
	}
}

func TestActionFormWithOptionsMergesDefaultsAndOptions(t *testing.T) {
	tests := []struct {
		name      string
		props     ActionFormOptions
		want      []string
		wantNever []string
	}{
		{
			name:  "default method is post and action adds matching htmx verb",
			props: ActionFormOptions{Action: "/items", Target: "#list", Swap: "outerHTML", Class: "inline"},
			want:  []string{`action="/items"`, `class="inline"`, `hx-post="/items"`, `hx-swap="outerHTML"`, `hx-target="#list"`, `method="post"`},
		},
		{
			name:      "explicit method is normalized for method and htmx attribute",
			props:     ActionFormOptions{Action: "/items/1", Method: " PUT "},
			want:      []string{`action="/items/1"`, `hx-put="/items/1"`, `method="put"`},
			wantNever: []string{`hx-post=`, `method=" PUT "`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderComponentForTest(t, ActionFormWithOptions(tt.props, TextNode("child")))
			assertContainsAll(t, got, tt.want)
			assertContainsNone(t, got, tt.wantNever)
		})
	}
}

func TestVariantComponentsRenderBoundaryAttributes(t *testing.T) {
	tests := []struct {
		name      string
		node      shared.Node
		want      []string
		wantNever []string
	}{
		{
			name: "button variants merge classes and disabled",
			node: ButtonWithVariants("Save", "primary", "sm", "wide", shared.ComponentProps{Class: "extra", Disabled: true}),
			want: []string{`class="btn btn-primary btn-sm btn-wide extra"`, `disabled="disabled"`},
		},
		{
			name: "input variants render disabled input classes",
			node: InputWithVariants("email", "a@example.test", "success", "lg", "ghost", shared.ComponentProps{Disabled: true}),
			want: []string{`class="input w-full input-success input-lg input-ghost"`, `disabled="disabled"`, `name="email"`, `value="a@example.test"`},
		},
		{
			name: "input variants accept xl size",
			node: InputWithVariants("email", "a@example.test", "success", "xl", "", shared.ComponentProps{}),
			want: []string{`class="input w-full input-success input-xl"`},
		},
		{
			name:      "input variants omit unknown size token",
			node:      InputWithVariants("email", "a@example.test", "success", "huge", "", shared.ComponentProps{}),
			want:      []string{`class="input w-full input-success"`},
			wantNever: []string{`input-huge`, `huge`},
		},
		{
			name: "select variants preserve selected option",
			node: SelectWithVariants("role", []shared.SelectOption{{Label: "Admin", Value: "admin", Selected: true}}, "primary", "sm", "ghost", shared.ComponentProps{}),
			want: []string{`class="select select-primary select-sm select-ghost"`, `<option selected="selected" value="admin">Admin</option>`},
		},
		{
			name: "select variants accept md size",
			node: SelectWithVariants("role", []shared.SelectOption{{Label: "Admin", Value: "admin"}}, "primary", "md", "", shared.ComponentProps{}),
			want: []string{`class="select select-primary select-md"`},
		},
		{
			name:      "textarea variants omit default rows",
			node:      TextareaWithVariants("notes", "body", "warning", "lg", "ghost", shared.TextareaOptions{}),
			want:      []string{`class="textarea w-full textarea-warning textarea-lg textarea-ghost"`, `name="notes"`, `>body</textarea>`},
			wantNever: []string{`rows=`},
		},
		{
			name: "textarea variants accept xs size",
			node: TextareaWithVariants("notes", "body", "warning", "xs", "", shared.TextareaOptions{}),
			want: []string{`class="textarea w-full textarea-warning textarea-xs"`},
		},
		{
			name:      "progress variant indeterminate omits value and defaults max",
			node:      Progress(shared.ProgressProps{Indeterminate: true, Props: shared.ComponentProps{Class: "progress-error"}}),
			want:      []string{`<progress class="progress w-full  h-2 progress-error" max="100">`},
			wantNever: []string{`value=`},
		},
		{
			name: "checkbox variants carry checked state",
			node: CheckboxWithVariants("terms", "yes", "Terms", "success", "sm", true, shared.ComponentProps{Class: "extra"}),
			want: []string{`checked="checked"`, `class="checkbox checkbox-success checkbox-sm extra"`, `type="checkbox"`},
		},
		{
			name: "checkbox variants accept xl size",
			node: CheckboxWithVariants("terms", "yes", "Terms", "success", "xl", true, shared.ComponentProps{}),
			want: []string{`class="checkbox checkbox-success checkbox-xl"`},
		},
		{
			name: "checkbox variants preserve disabled prop",
			node: CheckboxWithVariants("terms", "yes", "Terms", "success", "sm", false, shared.ComponentProps{Disabled: true}),
			want: []string{`class="checkbox checkbox-success checkbox-sm" disabled="disabled" name="terms" type="checkbox" value="yes"`},
		},
		{
			name: "radio variants carry disabled and checked states",
			node: RadioGroupWithVariants("tier", "primary", "lg", []shared.RadioItem{{Label: "Team", Value: "team", Checked: true}, {Label: "Enterprise", Value: "enterprise", Disabled: true}}, shared.ComponentProps{}),
			want: []string{`checked="checked" class="radio radio-primary radio-lg" name="tier" type="radio" value="team"`, `class="radio radio-primary radio-lg" disabled="disabled" name="tier" type="radio" value="enterprise"`},
		},
		{
			name: "radio variants accept md size",
			node: RadioGroupWithVariants("tier", "primary", "md", []shared.RadioItem{{Label: "Team", Value: "team"}}, shared.ComponentProps{}),
			want: []string{`class="radio radio-primary radio-md" name="tier" type="radio" value="team"`},
		},
		{
			name: "radio variants disabled prop disables all items",
			node: RadioGroupWithVariants("tier", "primary", "md", []shared.RadioItem{{Label: "Team", Value: "team"}, {Label: "Enterprise", Value: "enterprise"}}, shared.ComponentProps{Disabled: true}),
			want: []string{`class="radio radio-primary radio-md" disabled="disabled" name="tier" type="radio" value="team"`, `class="radio radio-primary radio-md" disabled="disabled" name="tier" type="radio" value="enterprise"`},
		},
		{
			name: "tabs variants emit roles and active selected state",
			node: TabsWithVariants([]shared.TabsItem{{Label: "Open", Active: true}, {Label: "Closed", Disabled: true}}, "boxed", "lifted", "sm", "custom"),
			want: []string{`<div class="tabs tabs-boxed tabs-lifted tabs-sm custom" role="tablist">`, `<button aria-selected="true" class="tab tab-active" role="tab" type="button">Open</button>`, `<button aria-selected="false" class="tab tab-disabled" disabled="disabled" role="tab" type="button">Closed</button>`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderComponentForTest(t, tt.node)
			assertContainsAll(t, got, tt.want)
			assertContainsNone(t, got, tt.wantNever)
		})
	}
}

func TestCoreComponentsRenderExpectedRolesAndClasses(t *testing.T) {
	alert := renderComponentForTest(t, Alert("Heads up", "Check settings", shared.ComponentProps{Class: "alert-warning"}))
	assertContainsAll(t, alert, []string{`role="alert"`, `class="alert alert-warning"`, `Heads up Check settings`})

	toast := renderComponentForTest(t, Toast("Saved", "Profile updated", shared.ComponentProps{Class: "toast-end"}))
	assertContainsAll(t, toast, []string{`class="toast toast-end"`, `class="alert"`, `role="alert"`, `Saved Profile updated`})

	modal := renderComponentForTest(t, Modal(shared.ModalProps{
		Title:   "Confirm delete",
		Body:    TextNode("Are you sure?"),
		Actions: Button("Delete", shared.ComponentProps{Class: "btn-error"}),
		Open:    true,
	}))
	assertContainsAll(t, modal, []string{`<dialog class="modal" open="open">`, `class="modal-box"`, `class="modal-backdrop"`, `method="dialog"`})
}

func TestCardRendersOptionalHeaderActionsAndChildren(t *testing.T) {
	withHeader := renderComponentForTest(t, Card(
		"Project",
		"Recent activity",
		Button("Archive", shared.ComponentProps{Class: "btn-sm"}),
		[]shared.Node{TextNode("card child")},
		shared.ComponentProps{Class: "border"},
	))
	assertContainsAll(t, withHeader, []string{`class="card bg-base-100 shadow-sm border"`, `class="card-body"`, `class="card-title"`, `Project`, `Recent activity`, `class="card-actions justify-end"`, `Archive`, `card child`})

	withoutHeader := renderComponentForTest(t, Card(
		"",
		"",
		nil,
		[]shared.Node{TextNode("only child")},
		shared.ComponentProps{},
	))
	assertContainsAll(t, withoutHeader, []string{`class="card bg-base-100 shadow-sm"`, `only child`})
	assertContainsNone(t, withoutHeader, []string{`class="card-body"`, `class="card-title"`, `class="card-actions justify-end"`})
}

func TestSkeletonDefaultsRowsToThree(t *testing.T) {
	for _, rows := range []int{0, -2} {
		t.Run("rows", func(t *testing.T) {
			got := renderComponentForTest(t, Skeleton(rows, shared.ComponentProps{Class: "animate-pulse"}))
			assertContainsAll(t, got, []string{`class="space-y-2 animate-pulse"`})
			if count := strings.Count(got, `class="skeleton h-4 w-full"`); count != 3 {
				t.Fatalf("expected 3 default skeleton rows, got %d in %q", count, got)
			}
		})
	}
}

func TestFormComponentsRenderHTMXAttributes(t *testing.T) {
	got := renderComponentForTest(t, ActionForm(shared.ActionFormProps{
		Action: "/contacts",
		Target: "#contacts",
		Swap:   "outerHTML",
		Props:  shared.ComponentProps{Class: "gap-4"},
	}, HiddenField("csrf", "token")))
	assertContainsAll(t, got, []string{`method="post"`, `action="/contacts"`, `hx-post="/contacts"`, `hx-target="#contacts"`, `hx-swap="outerHTML"`, `class="space-y-4 gap-4"`, `type="hidden"`})
}

func TestFormFieldRendersRequiredHintAndError(t *testing.T) {
	got := renderComponentForTest(t, FormField(
		Input("email", "", shared.ComponentProps{}),
		shared.FormFieldProps{Label: "Email", Required: true, Hint: "We never share it.", Error: "Email is required"},
	))
	assertContainsAll(t, got, []string{`class="fieldset w-full"`, `class="fieldset-legend"`, `Email *`, `We never share it.`, `class="label text-error"`, `Email is required`})
}

func TestLayoutComponentsRenderLandmarks(t *testing.T) {
	appShell := renderComponentForTest(t, AppShell(shared.AppShellProps{
		ID:      "app-shell",
		MainID:  "main-content",
		Header:  textNode("header", map[string]string{"class": "site-header"}, "Header"),
		Content: TextNode("Main content"),
		Props:   shared.ComponentProps{Class: "theme-shell"},
	}))
	assertContainsAll(t, appShell, []string{`<div class="min-h-screen bg-base-100 theme-shell" id="app-shell">`, `<main class="mx-auto w-full max-w-7xl p-4 md:p-6" id="main-content">`, `Main content`})

	region := renderComponentForTest(t, Region(shared.RegionProps{ID: "reports", Props: shared.ComponentProps{Class: "space-y-6"}}, TextNode("Reports")))
	assertContainsAll(t, region, []string{`<section class="space-y-3 space-y-6" id="reports">`, `Reports`})

	split := renderComponentForTest(t, Split(shared.SplitProps{Main: TextNode("Primary"), Aside: TextNode("Aside"), ReverseOnMobile: true, Props: shared.ComponentProps{Class: "items-start"}}))
	assertContainsAll(t, split, []string{`class="flex flex-col-reverse gap-4 items-start"`, `class="flex-1"`, `<aside class="w-full lg:w-80">`, `Primary`, `Aside`})

	container := renderComponentForTest(t, Container(shared.ContainerProps{MaxWidth: "max-w-5xl", Padding: "px-4", Centered: true, Props: shared.ComponentProps{Class: "dashboard"}}, TextNode("Contained")))
	assertContainsAll(t, container, []string{`class="w-full max-w-5xl mx-auto px-4 dashboard"`, `Contained`})
}

func TestAdminNavigationAndStatsComponents(t *testing.T) {
	navbar := renderComponentForTest(t, NavbarWithProps(NavbarProps{Class: "bg-base-200"}, MenuLink(MenuLinkProps{Label: "Dashboard", Href: "/admin", Icon: "🏠", Active: true, Class: "font-medium"})))
	assertContainsAll(t, navbar, []string{`class="navbar bg-base-200"`, `href="/admin"`, `class="font-medium active"`, `class="w-6 text-center"`, `Dashboard`, `🏠`})

	stats := renderComponentForTest(t, StatsWithProps(StatsProps{Class: "shadow"},
		StatItem(StatProps{Title: "Revenue", Value: "$1.2k", Description: "Today", Figure: TextNode("↗"), FigureClass: "text-success", Class: "place-items-center", ValueClass: "text-primary"}),
		StatItem(StatProps{Title: "Users", Value: "42", Description: "Active"}),
	))
	assertContainsAll(t, stats, []string{`class="stats shadow"`, `class="stat place-items-center"`, `class="stat-figure text-success"`, `class="stat-value text-primary"`, `Revenue`, `$1.2k`, `Users`, `42`})
	if count := strings.Count(stats, `class="stat-figure`); count != 1 {
		t.Fatalf("expected only the stat with a figure to render stat-figure, got %d in %q", count, stats)
	}
}

func TestAdminButtonAttrsMergeAndDisabled(t *testing.T) {
	button := renderComponentForTest(t, ButtonWithAttrs("Save", shared.ComponentProps{Class: "btn-primary", Disabled: true}, map[string]string{"class": "gap-2", "type": "submit", "data-action": "save"}))
	assertContainsAll(t, button, []string{`class="btn btn-primary gap-2"`, `data-action="save"`, `disabled="disabled"`, `type="submit"`, `Save`})

	contentButton := renderComponentForTest(t, ButtonContentWithAttrs(shared.ComponentProps{Class: "btn-ghost", Disabled: true}, map[string]string{"class": "w-full", "aria-label": "Open menu"}, TextNode("Menu")))
	assertContainsAll(t, contentButton, []string{`aria-label="Open menu"`, `class="btn btn-ghost w-full"`, `disabled="disabled"`, `Menu`})
}
