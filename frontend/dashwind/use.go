package dashwind

import (
	"encoding/json"
	"strings"

	"github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	"github.com/YoshihideShirai/marionette/frontend/assets"
)

// Options configures DashWind assets registered by Use.
type Options struct {
	// Theme sets the initial DaisyUI theme for pages that use DashWind.
	// When empty, Marionette's default theme bootstrap remains in control.
	Theme string
	// CustomCSS adds trusted inline CSS after the DashWind default CSS.
	CustomCSS string
	// DisableDefaultCSS skips DefaultCSS registration when the app provides its own shell CSS.
	DisableDefaultCSS bool
	// AssetsBasePath rewrites the DaisyUI/Tailwind framework imports to a self-hosted base path.
	// For example, "/assets/dashwind" resolves to "/assets/dashwind/daisyui.css" and
	// "/assets/dashwind/tailwindcss-browser.js". Leave empty to use the DaisyUI style template CDN defaults.
	AssetsBasePath string
}

// Use registers the DaisyUI style template, DashWind CSS, and optional DashWind settings on app.
func Use(app *backend.App, options Options) {
	if app == nil {
		return
	}
	app.UseStyleTemplate(styleTemplate())
	if base := strings.TrimSpace(options.AssetsBasePath); base != "" {
		app.UseAssets(assets.NewLocalAssetProvider(base))
	}
	if !options.DisableDefaultCSS {
		app.AddStyle(DefaultCSS)
	}
	if css := strings.TrimSpace(options.CustomCSS); css != "" {
		app.AddStyle(css)
	}
	if theme := strings.TrimSpace(options.Theme); theme != "" {
		app.AddJavaScript(themeBootstrapJS(theme))
	}
}

func styleTemplate() mf.StyleTemplate {
	return mf.DaisyUITemplate
}

func themeBootstrapJS(theme string) string {
	encoded, err := json.Marshal(theme)
	if err != nil {
		return ""
	}
	return `document.documentElement.setAttribute("data-theme", ` + string(encoded) + `);`
}
