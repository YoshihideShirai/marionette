package presets

import "github.com/YoshihideShirai/marionette/frontend/assets"

func FrameworkStylesheets() []string {
	return []string{assets.DaisyUICSSURL}
}

func FrameworkScripts() []string {
	return []string{assets.TailwindBrowserURL}
}
