package frontend

import (
	"github.com/YoshihideShirai/marionette/frontend/assets"
	daisyuipresets "github.com/YoshihideShirai/marionette/frontend/daisyui/presets"
	"github.com/YoshihideShirai/marionette/frontend/twailwindcss"
)

type StyleTemplate struct {
	Name                 string
	FrameworkStylesheets []string
	FrameworkScripts     []string

	// FrameworkStylesheetAssets and FrameworkScriptAssets name provider-resolved assets.
	// FrameworkStylesheets/FrameworkScripts remain supported for direct URL compatibility.
	FrameworkStylesheetAssets []assets.AssetName
	FrameworkScriptAssets     []assets.AssetName
}

var DaisyUITemplate = StyleTemplate{
	Name:                      "daisyui",
	FrameworkStylesheetAssets: daisyuipresets.FrameworkStylesheetAssets(),
	FrameworkScriptAssets:     daisyuipresets.FrameworkScriptAssets(),
	FrameworkStylesheets:      daisyuipresets.FrameworkStylesheets(),
	FrameworkScripts:          daisyuipresets.FrameworkScripts(),
}

var TailwindCSSTemplate = StyleTemplate{
	Name:                      twailwindcss.TemplateName,
	FrameworkStylesheetAssets: twailwindcss.FrameworkStylesheetAssets(),
	FrameworkScriptAssets:     twailwindcss.FrameworkScriptAssets(),
	FrameworkStylesheets:      twailwindcss.FrameworkStylesheets(),
	FrameworkScripts:          twailwindcss.FrameworkScripts(),
}

func DefaultStyleTemplate() StyleTemplate {
	return DaisyUITemplate
}

func StyleTemplateByName(name string) (StyleTemplate, bool) {
	switch name {
	case "daisyui", "tailadmin":
		return DaisyUITemplate, true
	case twailwindcss.TemplateName:
		return TailwindCSSTemplate, true
	default:
		return StyleTemplate{}, false
	}
}
