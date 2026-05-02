package frontend

import daisyuipresets "github.com/YoshihideShirai/marionette/frontend/daisyui/presets"

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
	Name: "tailwindcss",
	FrameworkScripts: []string{
		"https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4",
	},
}

func DefaultStyleTemplate() StyleTemplate {
	return DaisyUITemplate
}

func StyleTemplateByName(name string) (StyleTemplate, bool) {
	switch name {
	case "daisyui", "tailadmin":
		return DaisyUITemplate, true
	case "tailwindcss":
		return TailwindCSSTemplate, true
	default:
		return StyleTemplate{}, false
	}
}
