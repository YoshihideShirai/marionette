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
