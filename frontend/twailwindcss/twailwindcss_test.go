package twailwindcss

import (
	"slices"
	"testing"

	"github.com/YoshihideShirai/marionette/frontend/assets"
)

func TestTemplateNameIsTailwindCSS(t *testing.T) {
	if TemplateName != "tailwindcss" {
		t.Fatalf("TemplateName = %q, want %q", TemplateName, "tailwindcss")
	}
}

func TestBrowserURLMatchesTailwindBrowserAssetURL(t *testing.T) {
	want, ok := assets.DefaultProvider.ScriptURL(assets.TailwindCSSBrowser)
	if !ok {
		t.Fatal("DefaultProvider.ScriptURL(TailwindCSSBrowser) returned ok=false")
	}
	if BrowserURL != want {
		t.Fatalf("BrowserURL = %q, want TailwindCSSBrowser asset URL %q", BrowserURL, want)
	}
}

func TestFrameworkStylesheetAssetsIsNil(t *testing.T) {
	if got := FrameworkStylesheetAssets(); got != nil {
		t.Fatalf("FrameworkStylesheetAssets() = %#v, want nil slice", got)
	}
}

func TestFrameworkScriptAssetsReturnsTailwindBrowserAsset(t *testing.T) {
	got := FrameworkScriptAssets()
	want := []assets.AssetName{assets.TailwindCSSBrowser}
	if got == nil {
		t.Fatalf("FrameworkScriptAssets() = nil, want non-nil slice %v", want)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("FrameworkScriptAssets() = %v, want %v", got, want)
	}
}

func TestFrameworkStylesheetsIsNil(t *testing.T) {
	if got := FrameworkStylesheets(); got != nil {
		t.Fatalf("FrameworkStylesheets() = %#v, want nil slice", got)
	}
}

func TestFrameworkScriptsReturnsTailwindBrowserURL(t *testing.T) {
	scripts := FrameworkScripts()
	if scripts == nil {
		t.Fatalf("FrameworkScripts() = nil, want non-nil slice containing %q", BrowserURL)
	}
	if len(scripts) != 1 {
		t.Fatalf("FrameworkScripts() length = %d, want 1 (%v)", len(scripts), scripts)
	}
	if scripts[0] != BrowserURL {
		t.Fatalf("FrameworkScripts()[0] = %q, want BrowserURL %q", scripts[0], BrowserURL)
	}
}
