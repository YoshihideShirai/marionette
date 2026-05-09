package daisyui

import (
	"strings"
	"testing"

	shared "github.com/YoshihideShirai/marionette/frontend/shared"
)

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
			node: Input("email", "alice@example.test", shared.ComponentProps{Class: "input-bordered", Disabled: true}),
			want: []string{`class="input w-full input-bordered"`, `disabled="disabled"`, `name="email"`, `value="alice@example.test"`},
		},
		{
			name: "select marks only selected option",
			node: Select("status", []shared.SelectOption{{Label: "Draft", Value: "draft"}, {Label: "Published", Value: "published", Selected: true}}, shared.ComponentProps{Class: "select-sm"}),
			want: []string{`<select class="select select-sm" name="status">`, `<option value="draft">Draft</option>`, `<option selected="selected" value="published">Published</option>`},
		},
		{
			name: "textarea positive rows and required are rendered",
			node: Textarea("bio", "hello", shared.TextareaOptions{Placeholder: "About you", Rows: 4, Required: true, Props: shared.ComponentProps{Class: "textarea-primary"}}),
			want: []string{`class="textarea w-full textarea-primary"`, `name="bio"`, `placeholder="About you"`, `required="required"`, `rows="4"`, `>hello</textarea>`},
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
			name: "switch checked emits checkbox type and checked",
			node: Switch(shared.SwitchComponentProps{Name: "enabled", Value: "1", Label: "Enabled", Checked: true, Props: shared.ComponentProps{Class: "toggle-primary"}}),
			want: []string{`<input checked="checked" class="toggle toggle-primary" name="enabled" type="checkbox" value="1"></input>`, `<span class="label">Enabled</span>`},
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
			name: "select variants preserve selected option",
			node: SelectWithVariants("role", []shared.SelectOption{{Label: "Admin", Value: "admin", Selected: true}}, "primary", "sm", "ghost", shared.ComponentProps{}),
			want: []string{`class="select select-primary select-sm select-ghost"`, `<option selected="selected" value="admin">Admin</option>`},
		},
		{
			name:      "textarea variants omit default rows",
			node:      TextareaWithVariants("notes", "body", "warning", "lg", "ghost", shared.TextareaOptions{}),
			want:      []string{`class="textarea w-full textarea-warning textarea-lg textarea-ghost"`, `name="notes"`, `>body</textarea>`},
			wantNever: []string{`rows=`},
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
			name: "radio variants carry disabled and checked states",
			node: RadioGroupWithVariants("tier", "primary", "lg", []shared.RadioItem{{Label: "Team", Value: "team", Checked: true}, {Label: "Enterprise", Value: "enterprise", Disabled: true}}, shared.ComponentProps{}),
			want: []string{`checked="checked" class="radio radio-primary radio-lg" name="tier" type="radio" value="team"`, `class="radio radio-primary radio-lg" disabled="disabled" name="tier" type="radio" value="enterprise"`},
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
