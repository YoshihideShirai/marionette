package frontend

type StyleTemplate struct {
	Name                 string
	FrameworkStylesheets []string
	FrameworkScripts     []string
}

var TailAdminTemplate = StyleTemplate{
	Name: "tailadmin",
	FrameworkStylesheets: []string{
		"https://cdn.jsdelivr.net/npm/daisyui@5",
	},
	FrameworkScripts: []string{
		"https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4",
	},
}

var TailwindCSSTemplate = StyleTemplate{
	Name: "tailwindcss",
	FrameworkScripts: []string{
		"https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4",
	},
}

func DefaultStyleTemplate() StyleTemplate {
	return TailAdminTemplate
}

func StyleTemplateByName(name string) (StyleTemplate, bool) {
	switch name {
	case "tailadmin":
		return TailAdminTemplate, true
	case "tailwindcss":
		return TailwindCSSTemplate, true
	default:
		return StyleTemplate{}, false
	}
}
