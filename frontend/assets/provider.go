package assets

import (
	"path"
	"strings"
)

// AssetName identifies a framework/library asset that can be resolved by an AssetProvider.
type AssetName string

const (
	DaisyUI            AssetName = "daisyui"
	TailwindCSSBrowser AssetName = "tailwindcss-browser"
	ChartJS            AssetName = "chartjs"
	HTMX               AssetName = "htmx"
)

const (
	// DaisyUICSSFile is the default self-hosted DaisyUI stylesheet file name.
	DaisyUICSSFile = "daisyui.css"
	// TailwindCSSBrowserJSFile is the default self-hosted Tailwind browser runtime file name.
	TailwindCSSBrowserJSFile = "tailwindcss-browser.js"
	// HTMXJSFile is the default self-hosted HTMX runtime file name.
	HTMXJSFile = "htmx.min.js"
	// ChartJSFile is the default self-hosted Chart.js UMD bundle file name.
	ChartJSFile = "chart.umd.js"
)

// AssetProvider resolves known Marionette framework/library names to CSS and JS URLs.
type AssetProvider interface {
	StylesheetURL(name AssetName) (string, bool)
	ScriptURL(name AssetName) (string, bool)
}

// AssetResolver is kept as a descriptive alias for providers used by renderers.
type AssetResolver = AssetProvider

// CDNAssetProvider resolves framework/library assets to Marionette's default CDN URLs.
type CDNAssetProvider struct{}

// DefaultProvider is the provider used when an app or shell does not specify one.
var DefaultProvider AssetProvider = CDNAssetProvider{}

func (CDNAssetProvider) StylesheetURL(name AssetName) (string, bool) {
	switch name {
	case DaisyUI:
		return DaisyUICSSURL, true
	default:
		return "", false
	}
}

func (CDNAssetProvider) ScriptURL(name AssetName) (string, bool) {
	switch name {
	case TailwindCSSBrowser:
		return TailwindBrowserURL, true
	case ChartJS:
		return ChartJSURL, true
	case HTMX:
		return HTMXURL, true
	default:
		return "", false
	}
}

// LocalAssetProvider resolves assets to paths below a local/static base URL.
type LocalAssetProvider struct {
	BasePath    string
	Stylesheets map[AssetName]string
	Scripts     map[AssetName]string
}

// EmbeddedAssetProvider resolves embedded static assets exposed under a URL prefix.
type EmbeddedAssetProvider = LocalAssetProvider

// NewLocalAssetProvider creates a local provider with Marionette's default file names.
func NewLocalAssetProvider(basePath string) LocalAssetProvider {
	return LocalAssetProvider{
		BasePath:    basePath,
		Stylesheets: DefaultLocalStylesheets(),
		Scripts:     DefaultLocalScripts(),
	}
}

// NewEmbeddedAssetProvider creates an embedded provider with Marionette's default file names.
func NewEmbeddedAssetProvider(basePath string) EmbeddedAssetProvider {
	return EmbeddedAssetProvider(NewLocalAssetProvider(basePath))
}

func DefaultLocalStylesheets() map[AssetName]string {
	return map[AssetName]string{
		DaisyUI: DaisyUICSSFile,
	}
}

func DefaultLocalScripts() map[AssetName]string {
	return map[AssetName]string{
		TailwindCSSBrowser: TailwindCSSBrowserJSFile,
		ChartJS:            ChartJSFile,
		HTMX:               HTMXJSFile,
	}
}

func (p LocalAssetProvider) StylesheetURL(name AssetName) (string, bool) {
	file := strings.TrimSpace(p.Stylesheets[name])
	if file == "" {
		return "", false
	}
	return assetURL(p.BasePath, file), true
}

func (p LocalAssetProvider) ScriptURL(name AssetName) (string, bool) {
	file := strings.TrimSpace(p.Scripts[name])
	if file == "" {
		return "", false
	}
	return assetURL(p.BasePath, file), true
}

func assetURL(basePath, name string) string {
	basePath = strings.TrimRight(strings.TrimSpace(basePath), "/")
	name = strings.TrimLeft(strings.TrimSpace(name), "/")
	if basePath == "" {
		return name
	}
	if strings.HasPrefix(basePath, "http://") || strings.HasPrefix(basePath, "https://") {
		return basePath + "/" + name
	}
	if strings.HasPrefix(basePath, "/") {
		return path.Join(basePath, name)
	}
	return path.Join("/", basePath, name)
}
