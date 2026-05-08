package assets

import "testing"

func TestCDNAssetProviderResolvesKnownAssets(t *testing.T) {
	provider := CDNAssetProvider{}
	if got, ok := provider.StylesheetURL(DaisyUI); !ok || got != DaisyUICSSURL {
		t.Fatalf("expected DaisyUI CDN URL, got %q ok=%v", got, ok)
	}
	if got, ok := provider.ScriptURL(HTMX); !ok || got != HTMXURL {
		t.Fatalf("expected HTMX CDN URL, got %q ok=%v", got, ok)
	}
}

func TestLocalAssetProviderResolvesBasePath(t *testing.T) {
	provider := NewLocalAssetProvider("/assets/vendor")
	if got, ok := provider.StylesheetURL(DaisyUI); !ok || got != "/assets/vendor/daisyui.css" {
		t.Fatalf("expected local DaisyUI URL, got %q ok=%v", got, ok)
	}
	if got, ok := provider.ScriptURL(ChartJS); !ok || got != "/assets/vendor/chart.js" {
		t.Fatalf("expected local Chart.js URL, got %q ok=%v", got, ok)
	}
}
