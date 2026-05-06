package presets

import (
	"github.com/YoshihideShirai/marionette/frontend/assets"
	"github.com/YoshihideShirai/marionette/frontend/twailwindcss"
)

func FrameworkStylesheets() []string {
	return []string{assets.DaisyUICSSURL}
}

func FrameworkScripts() []string {
	return twailwindcss.FrameworkScripts()
}
