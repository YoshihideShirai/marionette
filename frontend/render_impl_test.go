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
	for _, want := range []string{"system", "mrnSetTheme", "data-mrn-theme-mode", "data-mrn-theme-icon", "prefers-color-scheme: dark"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected theme bootstrap to contain %q, got %q", want, out)
		}
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
		`src="/vendor/htmx.min.js"`,
		`src="/vendor/chart.umd.js"`,
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

func TestShellIncludesDefaultFeatureScripts(t *testing.T) {
	out, err := shell(template.HTML(`<div id="app"></div>`))
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}
	for _, want := range []string{assets.HTMXURL, assets.ChartJSURL, "window.mrnInitCharts", "htmx:afterSwap"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected default feature asset %q in shell output, got %q", want, out)
		}
	}
}

func TestShellDisableHTMXKeepsCharts(t *testing.T) {
	out, err := shellWithOptions(template.HTML(`<div id="app"></div>`), shellOptions{DisableHTMX: true})
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}
	for _, notWant := range []string{assets.HTMXURL} {
		if strings.Contains(out, notWant) {
			t.Fatalf("did not expect disabled HTMX asset %q in shell output, got %q", notWant, out)
		}
	}
	for _, want := range []string{assets.ChartJSURL, "window.mrnInitCharts", "htmx:afterSwap"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected chart asset/bootstrap %q to remain in shell output, got %q", want, out)
		}
	}
}

func TestShellDisableChartsKeepsHTMX(t *testing.T) {
	out, err := shellWithOptions(template.HTML(`<div id="app"></div>`), shellOptions{DisableCharts: true})
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}
	for _, notWant := range []string{assets.ChartJSURL, "window.mrnInitCharts", "htmx:afterSwap"} {
		if strings.Contains(out, notWant) {
			t.Fatalf("did not expect disabled chart asset/bootstrap %q in shell output, got %q", notWant, out)
		}
	}
	if !strings.Contains(out, assets.HTMXURL) {
		t.Fatalf("expected HTMX asset to remain in shell output, got %q", out)
	}
}

func TestShellDisableHTMXAndChartsOmitsFeatureScripts(t *testing.T) {
	out, err := shellWithOptions(template.HTML(`<div id="app"></div>`), shellOptions{
		DisableHTMX:   true,
		DisableCharts: true,
	})
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}
	for _, notWant := range []string{assets.HTMXURL, assets.ChartJSURL, "window.mrnInitCharts", "htmx:afterSwap"} {
		if strings.Contains(out, notWant) {
			t.Fatalf("did not expect disabled feature asset %q in shell output, got %q", notWant, out)
		}
	}
	if !strings.Contains(out, "mrnToggleTheme") {
		t.Fatalf("expected theme bootstrap to remain enabled, got %q", out)
	}
}

func TestShellAssetProviderSilentlySkipsMissingAssets(t *testing.T) {
	provider := assets.LocalAssetProvider{
		BasePath: "/vendor",
		Stylesheets: map[assets.AssetName]string{
			assets.DaisyUI: assets.DaisyUICSSFile,
		},
		Scripts: map[assets.AssetName]string{
			assets.HTMX: assets.HTMXJSFile,
		},
	}
	out, err := shellWithOptions(template.HTML(`<div id="app"></div>`), shellOptions{AssetProvider: provider})
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}
	for _, want := range []string{`href="/vendor/daisyui.css"`, `src="/vendor/htmx.min.js"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected provider-resolved asset %q in shell output, got %q", want, out)
		}
	}
	for _, notWant := range []string{
		`src="/vendor/tailwindcss-browser.js"`,
		`src="/vendor/chart.umd.js"`,
		assets.TailwindBrowserURL,
		assets.ChartJSURL,
	} {
		if strings.Contains(out, notWant) {
			t.Fatalf("did not expect missing/skipped asset %q in shell output, got %q", notWant, out)
		}
	}
}

func TestShellAssetPolicyForbidsExternalURLsByKind(t *testing.T) {
	policy := assets.AssetPolicy{ForbidExternalURLs: true}
	tests := []struct {
		name    string
		options shellOptions
		want    string
	}{
		{
			name: "framework stylesheet",
			options: shellOptions{
				AssetPolicy:          policy,
				FrameworkStylesheets: []string{"https://cdn.example.com/framework.css"},
				DisableHTMX:          true,
				DisableCharts:        true,
			},
			want: "asset policy forbids external URL for framework stylesheet",
		},
		{
			name: "framework script",
			options: shellOptions{
				AssetPolicy:      policy,
				FrameworkScripts: []string{"https://cdn.example.com/framework.js"},
				DisableHTMX:      true,
				DisableCharts:    true,
			},
			want: "asset policy forbids external URL for framework script",
		},
		{
			name: "stylesheet",
			options: shellOptions{
				AssetPolicy:          policy,
				FrameworkStylesheets: []string{"/assets/framework.css"},
				Stylesheets:          []string{"https://cdn.example.com/app.css"},
				DisableHTMX:          true,
				DisableCharts:        true,
			},
			want: "asset policy forbids external URL for stylesheet",
		},
		{
			name: "script",
			options: shellOptions{
				AssetPolicy:          policy,
				FrameworkStylesheets: []string{"/assets/framework.css"},
				Scripts:              []string{"https://cdn.example.com/app.js"},
				DisableHTMX:          true,
				DisableCharts:        true,
			},
			want: "asset policy forbids external URL for script",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := shellWithOptions(template.HTML(`<div id="app"></div>`), tt.options)
			if err == nil {
				t.Fatalf("expected asset policy error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error to contain %q, got %q", tt.want, err.Error())
			}
		})
	}
}

func TestShellAssetsUseStyleTemplateDeepCopiesTemplate(t *testing.T) {
	tpl := StyleTemplate{
		Name:                      "copy-test",
		FrameworkStylesheets:      []string{"/assets/original-framework.css"},
		FrameworkScripts:          []string{"/assets/original-framework.js"},
		FrameworkStylesheetAssets: []assets.AssetName{assets.DaisyUI},
		FrameworkScriptAssets:     []assets.AssetName{assets.TailwindCSSBrowser},
	}
	var shellAssets ShellAssets
	shellAssets.UseStyleTemplate(tpl)

	tpl.FrameworkStylesheets[0] = "/assets/mutated-framework.css"
	tpl.FrameworkScripts[0] = "/assets/mutated-framework.js"
	tpl.FrameworkStylesheetAssets[0] = assets.AssetName("mutated-stylesheet")
	tpl.FrameworkScriptAssets[0] = assets.AssetName("mutated-script")

	if got := shellAssets.StyleTemplate.FrameworkStylesheetAssets[0]; got != assets.DaisyUI {
		t.Fatalf("expected stylesheet asset slice to be deep-copied, got %q", got)
	}
	if got := shellAssets.StyleTemplate.FrameworkScriptAssets[0]; got != assets.TailwindCSSBrowser {
		t.Fatalf("expected script asset slice to be deep-copied, got %q", got)
	}

	// Force direct framework URLs so the rendered shell proves the copied URL slices,
	// rather than the mutated source template, are used after UseStyleTemplate returns.
	shellAssets.StyleTemplate.FrameworkStylesheetAssets = nil
	shellAssets.StyleTemplate.FrameworkScriptAssets = nil
	out, err := shellWithOptions(template.HTML(`<div id="app"></div>`), shellOptions{
		StyleTemplate: shellAssets.StyleTemplate,
		DisableHTMX:   true,
		DisableCharts: true,
	})
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}
	for _, want := range []string{`href="/assets/original-framework.css"`, `src="/assets/original-framework.js"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected copied template asset %q in shell output, got %q", want, out)
		}
	}
	for _, notWant := range []string{"/assets/mutated-framework.css", "/assets/mutated-framework.js"} {
		if strings.Contains(out, notWant) {
			t.Fatalf("did not expect mutated template asset %q in shell output, got %q", notWant, out)
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
