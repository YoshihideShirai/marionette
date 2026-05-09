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
	if got, ok := provider.ScriptURL(ChartJS); !ok || got != "/assets/vendor/chart.umd.js" {
		t.Fatalf("expected local Chart.js URL, got %q ok=%v", got, ok)
	}
}

func TestLocalAssetProviderResolvesURLVariants(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		file     string
		want     string
	}{
		{name: "empty base path", basePath: "", file: "daisyui.css", want: "daisyui.css"},
		{name: "relative base path is rooted", basePath: "assets/vendor", file: "daisyui.css", want: "/assets/vendor/daisyui.css"},
		{name: "trims slashes and whitespace", basePath: " /assets/vendor/ ", file: " /daisyui.css ", want: "/assets/vendor/daisyui.css"},
		{name: "keeps https base URL", basePath: "https://cdn.example.com/vendor/", file: "/daisyui.css", want: "https://cdn.example.com/vendor/daisyui.css"},
		{name: "keeps http base URL", basePath: "http://cdn.example.com/vendor/", file: "/daisyui.css", want: "http://cdn.example.com/vendor/daisyui.css"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := LocalAssetProvider{
				BasePath: tt.basePath,
				Stylesheets: map[AssetName]string{
					DaisyUI: tt.file,
				},
			}
			got, ok := provider.StylesheetURL(DaisyUI)
			if !ok {
				t.Fatalf("StylesheetURL() ok = false, want true")
			}
			if got != tt.want {
				t.Fatalf("StylesheetURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLocalAssetProviderIgnoresBlankOrUnknownAssets(t *testing.T) {
	provider := LocalAssetProvider{
		BasePath: "/assets",
		Stylesheets: map[AssetName]string{
			DaisyUI: "  ",
		},
		Scripts: map[AssetName]string{
			HTMX: "\t",
		},
	}

	if got, ok := provider.StylesheetURL(DaisyUI); ok || got != "" {
		t.Fatalf("StylesheetURL(DaisyUI) = %q ok=%v, want empty and false", got, ok)
	}
	if got, ok := provider.StylesheetURL(AssetName("missing")); ok || got != "" {
		t.Fatalf("StylesheetURL(missing) = %q ok=%v, want empty and false", got, ok)
	}
	if got, ok := provider.ScriptURL(HTMX); ok || got != "" {
		t.Fatalf("ScriptURL(HTMX) = %q ok=%v, want empty and false", got, ok)
	}
	if got, ok := provider.ScriptURL(AssetName("missing")); ok || got != "" {
		t.Fatalf("ScriptURL(missing) = %q ok=%v, want empty and false", got, ok)
	}
}
