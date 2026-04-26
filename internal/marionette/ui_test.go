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
