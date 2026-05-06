package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterThemeControllerExample(app *mb.App) {
	app.Page("/theme-controller", func(ctx *mb.Context) mf.Node {
		return mf.DivProps(mf.ElementProps{Class: "space-y-8"},
			mf.ThemeController(
				mf.ThemeControllerOption("light", false, "btn"),
				mf.ThemeControllerOption("dark", false, "btn"),
				mf.ThemeControllerOption("cupcake", true, "btn"),
				mf.ThemeControllerOption("dracula", false, "btn"),
			),
			mf.ThemeController(
				mf.ThemeControllerOption("corporate", true, "radio radio-xs"),
				mf.ThemeControllerOption("forest", false, "radio radio-xs"),
				mf.ThemeControllerOption("luxury", false, "radio radio-xs"),
				mf.ThemeControllerOption("night", false, "radio radio-xs"),
			),
			mf.DivProps(mf.ElementProps{Class: "space-y-2"},
				themeToggle("Dark", "toggle", "dark"),
				themeToggle("Retro", "checkbox", "retro"),
			),
			mf.DivProps(mf.ElementProps{Class: "grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-3"},
				themeSwatch("light", "Light"),
				themeSwatch("cupcake", "Cupcake"),
				themeSwatch("emerald", "Emerald"),
				themeSwatch("night", "Night"),
				themeSwatch("dracula", "Dracula"),
			),
			themeSelect(),
		)
	})
}

func themeToggle(label, inputClass, value string) mf.Node {
	return mf.Element("label", mf.ElementProps{Class: "label cursor-pointer gap-2 w-fit"},
		mf.SpanProps(mf.ElementProps{Class: "label"}, mf.Text(label)),
		mf.Element("input", mf.ElementProps{Attrs: mf.Attrs{
			"type":  "checkbox",
			"class": inputClass + " theme-controller",
			"value": value,
		}}),
	)
}

func themeSwatch(value, label string) mf.Node {
	return mf.Element("label", mf.ElementProps{Class: "cursor-pointer"},
		mf.Element("input", mf.ElementProps{Attrs: mf.Attrs{
			"type":  "radio",
			"name":  "theme-picker",
			"value": value,
			"class": "theme-controller sr-only",
		}}),
		mf.SpanProps(mf.ElementProps{Class: "block rounded-box border border-base-300 p-2"},
			mf.SpanProps(mf.ElementProps{Class: "font-semibold text-xs"}, mf.Text(label)),
			mf.DivProps(mf.ElementProps{Class: "mt-2 grid grid-cols-4 gap-1"},
				colorDot("bg-primary"),
				colorDot("bg-secondary"),
				colorDot("bg-accent"),
				colorDot("bg-neutral"),
			),
		),
	)
}

func colorDot(className string) mf.Node {
	return mf.SpanProps(mf.ElementProps{Class: "h-3 rounded " + className})
}

func themeSelect() mf.Node {
	return mf.Element("fieldset", mf.ElementProps{Class: "fieldset w-full max-w-xs"},
		mf.Element("legend", mf.ElementProps{Class: "fieldset-legend"}, mf.Text("Choose a theme")),
		mf.Element("select", mf.ElementProps{Class: "select theme-controller", Attrs: mf.Attrs{"name": "theme-picker"}},
			themeOption("light", "Light"),
			themeOption("cupcake", "Cupcake"),
			themeOption("bumblebee", "Bumblebee"),
			themeOption("synthwave", "Synthwave"),
			themeOption("forest", "Forest"),
			themeOption("black", "Black"),
		),
	)
}

func themeOption(value, label string) mf.Node {
	return mf.Element("option", mf.ElementProps{Attrs: mf.Attrs{"value": value}}, mf.Text(label))
}
