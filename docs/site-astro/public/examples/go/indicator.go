package main
import m "github.com/YoshihideShirai/marionette"
func example() m.Node { return m.Indicator(m.Text("new"), m.Button("Inbox", m.ComponentProps{})) }
