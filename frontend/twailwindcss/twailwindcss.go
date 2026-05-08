package twailwindcss

import "github.com/YoshihideShirai/marionette/frontend/assets"

const (
	TemplateName = "tailwindcss"
	BrowserURL   = assets.TailwindBrowserURL
)

func FrameworkStylesheetAssets() []assets.AssetName {
	return nil
}

func FrameworkScriptAssets() []assets.AssetName {
	return []assets.AssetName{assets.TailwindCSSBrowser}
}

func FrameworkStylesheets() []string {
	return nil
}

func FrameworkScripts() []string {
	if url, ok := assets.DefaultProvider.ScriptURL(assets.TailwindCSSBrowser); ok && url != "" {
		return []string{url}
	}
	return nil
}
