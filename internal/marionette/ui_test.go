package marionette

import (
	"strings"
	"testing"
)

func TestButtonRenderUsesHTMXMarkup(t *testing.T) {
	html, err := Button("Increment").OnClick("counter/increment").Target("#app").Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	if got == "" {
		t.Fatal("expected non-empty button html")
	}
	if want := `hx-post="/counter/increment"`; !strings.Contains(got, want) {
		t.Fatalf("expected %q in %q", want, got)
	}
}

func TestButtonPostAcceptsLeadingSlash(t *testing.T) {
	html, err := Button("Save").Post("/users/create").Target("#users").Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	if strings.Contains(got, `hx-post="//users/create"`) {
		t.Fatalf("expected normalized post path, got %q", got)
	}
	if want := `hx-target="#users"`; !strings.Contains(got, want) {
		t.Fatalf("expected %q in %q", want, got)
	}
}

func TestFormInputAndSubmitRenderHTMXMarkup(t *testing.T) {
	html, err := Form("users/create", Input("name", `<Aiko>`), Submit("Create")).Target("#users").Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`hx-post="/users/create"`,
		`hx-target="#users"`,
		`name="name"`,
		`value="&lt;Aiko&gt;"`,
		`type="submit"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestSidebarRendersNavigationAndEscapesText(t *testing.T) {
	html, err := Sidebar("Marionette", "Admin <Console>",
		SidebarLink("Users", "/").Active(),
		SidebarLink("Settings", "/settings"),
	).Note("Demo", `<unsafe>`).Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`<aside`,
		`href="/"`,
		`href="/settings"`,
		`btn btn-primary justify-start`,
		`Admin &lt;Console&gt;`,
		`&lt;unsafe&gt;`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
	if strings.Contains(got, "Admin <Console>") || strings.Contains(got, `<unsafe>`) {
		t.Fatalf("expected sidebar text to be escaped, got %q", got)
	}
}

func TestTableRendersHeadersRowsAndEscapesCells(t *testing.T) {
	html, err := Table([]string{"Name", "Role"},
		TableRow(Text(`<Aiko>`), DivClass("", "badge", Text("Admin"))),
	).Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`<table class="table">`,
		`<th>Name</th>`,
		`<td><span>&lt;Aiko&gt;</span></td>`,
		`class="badge"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
	if strings.Contains(got, `<Aiko>`) {
		t.Fatalf("expected table cell text to be escaped, got %q", got)
	}
}

func TestHiddenInputRenderEscapesValue(t *testing.T) {
	html, err := HiddenInput("id", `"42"`).Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`name="id"`,
		`type="hidden"`,
		`value="&#34;42&#34;"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestElementRenderEscapesText(t *testing.T) {
	html, err := Text(`<script>alert(1)</script>`).Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	if strings.Contains(got, `<script>alert(1)</script>`) {
		t.Fatalf("expected escaped content, got %q", got)
	}
}

func TestTemplatePartialsRenderSharedProps(t *testing.T) {
	buttonHTML, err := ComponentButton("Send", ComponentProps{Class: "tracking-wide", Variant: "secondary", Size: "sm", Disabled: true}).Render()
	if err != nil {
		t.Fatalf("button render failed: %v", err)
	}
	inputHTML, err := ComponentInput("email", "demo@example.com", ComponentProps{Variant: "ghost", Size: "sm", Disabled: true}).Render()
	if err != nil {
		t.Fatalf("input render failed: %v", err)
	}
	selectHTML, err := ComponentSelect("role", []SelectOption{{Label: "Viewer", Value: "viewer", Selected: true}}, ComponentProps{Variant: "ghost", Size: "sm", Disabled: true}).Render()
	if err != nil {
		t.Fatalf("select render failed: %v", err)
	}

	for _, tc := range []struct {
		name string
		html string
		want []string
	}{
		{name: "button", html: string(buttonHTML), want: []string{`btn-secondary`, `btn-sm`, `tracking-wide`, `disabled`}},
		{name: "input", html: string(inputHTML), want: []string{`input-ghost`, `input-sm`, `name="email"`, `disabled`}},
		{name: "select", html: string(selectHTML), want: []string{`select-ghost`, `select-sm`, `name="role"`, `selected`, `disabled`}},
	} {
		for _, want := range tc.want {
			if !strings.Contains(tc.html, want) {
				t.Fatalf("%s expected %q in %q", tc.name, want, tc.html)
			}
		}
	}
}

func TestComponentInputWithOptionsRendersDateConstraints(t *testing.T) {
	html, err := ComponentInputWithOptions("start_date", "2030-01-01", InputOptions{
		Type:     "date",
		Min:      "2024-01-01",
		Max:      "2026-12-31",
		Required: true,
		Props:    ComponentProps{Variant: "default", Size: "sm"},
	}).Render()
	if err != nil {
		t.Fatalf("input render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{`type="date"`, `min="2024-01-01"`, `max="2026-12-31"`, `required`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestComponentFormFieldRendersLabelHintAndError(t *testing.T) {
	html, err := ComponentFormField(
		ComponentInputWithOptions("name", "", InputOptions{Required: true, Props: ComponentProps{Variant: "default", Size: "sm"}}),
		FormFieldProps{
			Label:    "Name",
			Required: true,
			Hint:     "Enter a display name.",
			Error:    "Name is required.",
		},
	).Render()
	if err != nil {
		t.Fatalf("form field render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{`label-text`, `Name`, `*`, `Enter a display name.`, `Name is required.`, `name="name"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestComponentModalRendersSSRState(t *testing.T) {
	closedHTML, err := ComponentModal(ModalProps{
		Title:   "Delete user",
		Body:    Text("Confirm deletion"),
		Actions: ComponentButton("Cancel", ComponentProps{Variant: "ghost", Size: "sm"}),
		Open:    false,
	}).Render()
	if err != nil {
		t.Fatalf("closed modal render failed: %v", err)
	}
	openHTML, err := ComponentModal(ModalProps{
		Title:   "Delete user",
		Body:    Text("Confirm deletion"),
		Actions: ComponentButton("Delete", ComponentProps{Variant: "danger", Size: "sm"}),
		Open:    true,
	}).Render()
	if err != nil {
		t.Fatalf("open modal render failed: %v", err)
	}

	if strings.Contains(string(closedHTML), "modal-open") {
		t.Fatalf("expected closed modal without modal-open class, got %q", closedHTML)
	}
	for _, want := range []string{`modal-open`, `Delete user`, `Confirm deletion`, `btn-error`} {
		if !strings.Contains(string(openHTML), want) {
			t.Fatalf("expected %q in %q", want, openHTML)
		}
	}
}

func TestComponentEmptyStateRendersSkeletonAndCopy(t *testing.T) {
	skeletonHTML, err := ComponentEmptyState(EmptyStateProps{Skeleton: true, Rows: 2}).Render()
	if err != nil {
		t.Fatalf("skeleton render failed: %v", err)
	}
	if !strings.Contains(string(skeletonHTML), `aria-busy="true"`) {
		t.Fatalf("expected skeleton aria-busy state, got %q", skeletonHTML)
	}

	emptyHTML, err := ComponentEmptyState(EmptyStateProps{Title: "No users", Description: "Create one first."}).Render()
	if err != nil {
		t.Fatalf("empty render failed: %v", err)
	}
	for _, want := range []string{"No users", "Create one first."} {
		if !strings.Contains(string(emptyHTML), want) {
			t.Fatalf("expected %q in %q", want, emptyHTML)
		}
	}
}

func TestComponentTableRendersSortHeadersAndEmptyState(t *testing.T) {
	emptyHTML, err := ComponentTable(TableProps{
		Columns: []TableColumn{
			{Label: "Name", SortKey: "name", SortHref: "/?sort=name", SortActive: true},
			{Label: "Role"},
		},
		EmptyTitle:       "No users",
		EmptyDescription: "Create a user to get started.",
	}).Render()
	if err != nil {
		t.Fatalf("empty table render failed: %v", err)
	}
	got := string(emptyHTML)
	for _, want := range []string{`href="/?sort=name"`, `No users`, `Create a user to get started.`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestComponentPaginationRendersState(t *testing.T) {
	html, err := ComponentPagination(PaginationProps{
		Page:       2,
		TotalPages: 4,
		PrevHref:   "/?page=1&per_page=10",
		NextHref:   "/?page=3&per_page=10",
	}).Render()
	if err != nil {
		t.Fatalf("pagination render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{"Page 2 / 4", `href="/?page=1&amp;per_page=10"`, `href="/?page=3&amp;per_page=10"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}
