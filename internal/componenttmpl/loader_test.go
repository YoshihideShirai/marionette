package componenttmpl

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemplateFile(t *testing.T, dir, name, contents string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write template %s: %v", path, err)
	}
}

func TestLoadReadsTmplAndHTMLTemplates(t *testing.T) {
	dir := t.TempDir()
	writeTemplateFile(t, dir, "button.tmpl", `{{define "components/button"}}<button>{{.Label}}</button>{{end}}`)
	writeTemplateFile(t, dir, "card.html", `{{define "components/card"}}<section>{{.Title}}</section>{{end}}`)

	tmpl, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	var button bytes.Buffer
	if err := tmpl.ExecuteTemplate(&button, "components/button", map[string]string{"Label": "Save"}); err != nil {
		t.Fatalf("ExecuteTemplate(components/button) error = %v", err)
	}
	if got, want := button.String(), "<button>Save</button>"; got != want {
		t.Fatalf("button template = %q, want %q", got, want)
	}

	var card bytes.Buffer
	if err := tmpl.ExecuteTemplate(&card, "components/card", map[string]string{"Title": "Profile"}); err != nil {
		t.Fatalf("ExecuteTemplate(components/card) error = %v", err)
	}
	if got, want := card.String(), "<section>Profile</section>"; got != want {
		t.Fatalf("card template = %q, want %q", got, want)
	}
}

func TestSourceLoadAndNodeRender(t *testing.T) {
	dir := t.TempDir()
	writeTemplateFile(t, dir, "greeting.tmpl", `{{define "components/greeting"}}<p>Hello, {{.Name}}!</p>{{end}}`)

	source := NewSource(dir)
	if _, err := source.Load(); err != nil {
		t.Fatalf("Source.Load() error = %v", err)
	}

	html, err := (Node{
		Source: source,
		Name:   "components/greeting",
		Data:   map[string]string{"Name": "Aiko <Admin>"},
	}).Render()
	if err != nil {
		t.Fatalf("Node.Render() error = %v", err)
	}
	if got, want := string(html), "<p>Hello, Aiko &lt;Admin&gt;!</p>"; got != want {
		t.Fatalf("Node.Render() = %q, want %q", got, want)
	}
}

func TestLoadEmptyDirectoryReturnsError(t *testing.T) {
	_, err := Load(t.TempDir())
	if err == nil {
		t.Fatal("Load() error = nil, want no templates error")
	}
	if !strings.Contains(err.Error(), "no component templates found") {
		t.Fatalf("Load() error = %q, want no component templates found", err)
	}
}

func TestLoadBrokenTemplateReturnsParseError(t *testing.T) {
	dir := t.TempDir()
	writeTemplateFile(t, dir, "broken.tmpl", `{{define "components/broken"}}<p>{{if .Open}}</p>{{end}}`)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("Load() error = nil, want parse error")
	}
	if !strings.Contains(err.Error(), "failed to parse component templates") {
		t.Fatalf("Load() error = %q, want failed to parse component templates", err)
	}
}

func TestNilSourceLoadReturnsError(t *testing.T) {
	var source *Source
	_, err := source.Load()
	if err == nil {
		t.Fatal("(*Source)(nil).Load() error = nil, want nil source error")
	}
	if got, want := err.Error(), "component template source is nil"; got != want {
		t.Fatalf("(*Source)(nil).Load() error = %q, want %q", got, want)
	}
}

func TestNodeRenderMissingTemplateReturnsError(t *testing.T) {
	dir := t.TempDir()
	writeTemplateFile(t, dir, "existing.tmpl", `{{define "components/existing"}}<p>Existing</p>{{end}}`)

	_, err := (Node{
		Source: NewSource(dir),
		Name:   "components/missing",
	}).Render()
	if err == nil {
		t.Fatal("Node.Render() error = nil, want missing template error")
	}
}
