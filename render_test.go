package marionette

import (
	"html/template"
	"strings"
	"testing"
)

func TestShellUsesStableRootID(t *testing.T) {
	out, err := shell(template.HTML(`<div id="app"></div>`))
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}

	if !strings.Contains(out, `id="marionette-root"`) {
		t.Fatalf("expected marionette-root id in shell, got %q", out)
	}
	if strings.Count(out, `id="app"`) != 1 {
		t.Fatalf("expected exactly one app id in shell output, got %q", out)
	}
}

func TestShellIncludesThemeBootstrapScript(t *testing.T) {
	out, err := shell(template.HTML(`<div id="app"></div>`))
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}

	if !strings.Contains(out, "window.mrnToggleTheme") {
		t.Fatalf("expected theme toggle helper in shell output, got %q", out)
	}
	if !strings.Contains(out, "marionette-theme") {
		t.Fatalf("expected localStorage theme key in shell output, got %q", out)
	}
}

func TestShellIncludesDefaultStyleBeforeCustomStyles(t *testing.T) {
	out, err := shellWithOptions(template.HTML(`<div id="app"></div>`), shellOptions{
		Styles: []template.CSS{`#marionette-root { max-width: 48rem; }`},
	})
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}

	defaultIndex := strings.Index(out, "--mrn-page-max-width")
	if defaultIndex == -1 {
		t.Fatalf("expected default Marionette CSS in shell output, got %q", out)
	}
	customIndex := strings.Index(out, "#marionette-root { max-width: 48rem; }")
	if customIndex == -1 {
		t.Fatalf("expected custom CSS in shell output, got %q", out)
	}
	if customIndex < defaultIndex {
		t.Fatalf("expected custom CSS after default CSS, got %q", out)
	}
}
