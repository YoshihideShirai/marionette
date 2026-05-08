package frontend

import (
	"html/template"
	"strings"
	"testing"

	"github.com/YoshihideShirai/marionette/frontend/assets"
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
	if !strings.Contains(out, "mrnToggleTheme") {
		t.Fatalf("expected theme toggle helper in shell output, got %q", out)
	}
	if !strings.Contains(out, "marionette-theme") {
		t.Fatalf("expected localStorage theme key in shell output, got %q", out)
	}
}

func TestShellAppliesDefaultStyleTemplate(t *testing.T) {
	out, err := shell(template.HTML(`<div id="app"></div>`))
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}
	defaults := DefaultStyleTemplate()
	for _, href := range defaults.FrameworkStylesheets {
		if !strings.Contains(out, `href="`+href+`"`) {
			t.Fatalf("expected default stylesheet %q in shell output", href)
		}
	}
	for _, src := range defaults.FrameworkScripts {
		if !strings.Contains(out, `src="`+src+`"`) {
			t.Fatalf("expected default script %q in shell output", src)
		}
	}
}

func TestShellResolvesBuiltInAssetsThroughProvider(t *testing.T) {
	provider := assets.NewLocalAssetProvider("/vendor")
	out, err := shellWithOptions(template.HTML(`<div id="app"></div>`), shellOptions{AssetProvider: provider})
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}
	for _, want := range []string{
		`href="/vendor/daisyui.css"`,
		`src="/vendor/tailwindcss-browser.js"`,
		`src="/vendor/htmx.js"`,
		`src="/vendor/chart.js"`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected provider-resolved asset %q in shell output, got %q", want, out)
		}
	}
	for _, legacyCDN := range []string{assets.DaisyUICSSURL, assets.TailwindBrowserURL, assets.HTMXURL, assets.ChartJSURL} {
		if strings.Contains(out, legacyCDN) {
			t.Fatalf("did not expect CDN URL %q when provider is overridden, got %q", legacyCDN, out)
		}
	}
}

func TestShellOptionsStyleTemplateOverridesDefault(t *testing.T) {
	out, err := shellWithOptions(template.HTML(`<div id="app"></div>`), shellOptions{
		StyleTemplate: StyleTemplate{
			Name:                 "custom",
			FrameworkStylesheets: []string{"https://cdn.example.com/custom.css"},
			FrameworkScripts:     []string{"https://cdn.example.com/custom.js"},
		},
	})
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}
	if !strings.Contains(out, `href="https://cdn.example.com/custom.css"`) {
		t.Fatalf("expected custom style template stylesheet in shell output")
	}
	if !strings.Contains(out, `src="https://cdn.example.com/custom.js"`) {
		t.Fatalf("expected custom style template script in shell output")
	}
	if strings.Contains(out, DefaultStyleTemplate().FrameworkStylesheets[0]) {
		t.Fatalf("did not expect default style template stylesheet after override")
	}
}

func TestShellIncludesDefaultStyleBeforeCustomStyles(t *testing.T) {
	out, err := shellWithOptions(template.HTML(`<div id="app"></div>`), shellOptions{Styles: []template.CSS{`#marionette-root { max-width: 48rem; }`}})
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
