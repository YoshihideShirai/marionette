package main
import m "github.com/YoshihideShirai/marionette"
func example() m.Node { return m.Timeline(m.TimelineItem("1", "", m.Text("Import frontend")), m.TimelineItem("2", "", m.Text("Call frontend.Timeline()"))) }
