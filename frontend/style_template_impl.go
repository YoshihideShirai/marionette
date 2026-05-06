package frontend

import (
	daisyuipresets "github.com/YoshihideShirai/marionette/frontend/daisyui/presets"
	"github.com/YoshihideShirai/marionette/frontend/twailwindcss"
)

type StyleTemplate struct {
	Name                 string
	FrameworkStylesheets []string
	FrameworkScripts     []string
}

var DaisyUITemplate = StyleTemplate{
	Name:                 "daisyui",
	FrameworkStylesheets: daisyuipresets.FrameworkStylesheets(),
	FrameworkScripts:     daisyuipresets.FrameworkScripts(),
}

var TailwindCSSTemplate = StyleTemplate{
	Name:                 twailwindcss.TemplateName,
	FrameworkStylesheets: twailwindcss.FrameworkStylesheets(),
	FrameworkScripts:     twailwindcss.FrameworkScripts(),
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
