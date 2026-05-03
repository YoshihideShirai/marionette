package main

import m "github.com/YoshihideShirai/marionette"

func example() m.Node {
	return m.Container(
		m.Divider(m.DividerProps{}),
		m.Hero("DaisyUI aliases", "Call daisyui components from frontend package directly."),
		m.Timeline(
			m.TimelineItem("Step 1", "", m.Code("import frontend")),
			m.TimelineItem("Step 2", "", m.Code("frontend.Hero(...)")),
		),
	)
}
