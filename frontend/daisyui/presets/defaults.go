package presets

import (
	"github.com/YoshihideShirai/marionette/frontend/assets"
	"github.com/YoshihideShirai/marionette/frontend/twailwindcss"
)

func FrameworkStylesheetAssets() []assets.AssetName {
	return []assets.AssetName{assets.DaisyUI}
}

func FrameworkScriptAssets() []assets.AssetName {
	return twailwindcss.FrameworkScriptAssets()
}

func FrameworkStylesheets() []string {
	return resolveStylesheets(assets.DefaultProvider, FrameworkStylesheetAssets())
}

func FrameworkScripts() []string {
	return resolveScripts(assets.DefaultProvider, FrameworkScriptAssets())
}

func resolveStylesheets(provider assets.AssetProvider, names []assets.AssetName) []string {
	out := make([]string, 0, len(names))
	for _, name := range names {
		if url, ok := provider.StylesheetURL(name); ok && url != "" {
			out = append(out, url)
		}
	}
	return out
}

func resolveScripts(provider assets.AssetProvider, names []assets.AssetName) []string {
	out := make([]string, 0, len(names))
	for _, name := range names {
		if url, ok := provider.ScriptURL(name); ok && url != "" {
			out = append(out, url)
		}
	}
	return out
}
