package presets

import (
	"slices"
	"testing"

	"github.com/YoshihideShirai/marionette/frontend/assets"
)

type fakeAssetProvider struct {
	stylesheets map[assets.AssetName]assetResult
	scripts     map[assets.AssetName]assetResult
}

type assetResult struct {
	url string
	ok  bool
}

func (p fakeAssetProvider) StylesheetURL(name assets.AssetName) (string, bool) {
	result, ok := p.stylesheets[name]
	if !ok {
		return "", false
	}
	return result.url, result.ok
}

func (p fakeAssetProvider) ScriptURL(name assets.AssetName) (string, bool) {
	result, ok := p.scripts[name]
	if !ok {
		return "", false
	}
	return result.url, result.ok
}

func TestFrameworkStylesheetAssetsReturnsDaisyUI(t *testing.T) {
	got := FrameworkStylesheetAssets()
	want := []assets.AssetName{assets.DaisyUI}
	if !slices.Equal(got, want) {
		t.Fatalf("FrameworkStylesheetAssets() = %v, want %v", got, want)
	}
}

func TestFrameworkScriptAssetsDelegatesTailwindBrowser(t *testing.T) {
	got := FrameworkScriptAssets()
	if !slices.Contains(got, assets.TailwindCSSBrowser) {
		t.Fatalf("FrameworkScriptAssets() = %v, want to include %q", got, assets.TailwindCSSBrowser)
	}
}

func TestFrameworkStylesheetsResolvesDefaultProvider(t *testing.T) {
	got := FrameworkStylesheets()
	if len(got) == 0 {
		t.Fatalf("FrameworkStylesheets() is empty, want DaisyUI provider URL")
	}
	if got[0] != assets.DaisyUICSSURL {
		t.Fatalf("FrameworkStylesheets()[0] = %q, want %q", got[0], assets.DaisyUICSSURL)
	}
}

func TestFrameworkScriptsResolvesDefaultProvider(t *testing.T) {
	got := FrameworkScripts()
	if len(got) == 0 {
		t.Fatalf("FrameworkScripts() is empty, want Tailwind browser URL")
	}
	if !slices.Contains(got, assets.TailwindBrowserURL) {
		t.Fatalf("FrameworkScripts() = %v, want to include %q", got, assets.TailwindBrowserURL)
	}
}

func TestResolveStylesheetsSkipsUnknownOrBlankProviderResults(t *testing.T) {
	provider := fakeAssetProvider{
		stylesheets: map[assets.AssetName]assetResult{
			assets.DaisyUI:            {url: "/assets/daisyui.css", ok: true},
			assets.AssetName("blank"): {url: "", ok: true},
		},
	}

	got := resolveStylesheets(provider, []assets.AssetName{
		assets.AssetName("unknown"),
		assets.AssetName("blank"),
		assets.DaisyUI,
	})
	want := []string{"/assets/daisyui.css"}
	if !slices.Equal(got, want) {
		t.Fatalf("resolveStylesheets() = %v, want %v", got, want)
	}
}

func TestResolveScriptsSkipsUnknownOrBlankProviderResults(t *testing.T) {
	provider := fakeAssetProvider{
		scripts: map[assets.AssetName]assetResult{
			assets.TailwindCSSBrowser: {url: "/assets/tailwindcss-browser.js", ok: true},
			assets.AssetName("blank"): {url: "", ok: true},
		},
	}

	got := resolveScripts(provider, []assets.AssetName{
		assets.AssetName("unknown"),
		assets.AssetName("blank"),
		assets.TailwindCSSBrowser,
	})
	want := []string{"/assets/tailwindcss-browser.js"}
	if !slices.Equal(got, want) {
		t.Fatalf("resolveScripts() = %v, want %v", got, want)
	}
}
