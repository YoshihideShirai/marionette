package main
import m "github.com/YoshihideShirai/marionette"
func example() m.Node { return m.Drawer("drawer", m.Menu(m.Text("Sidebar")), m.Button("Open", m.ComponentProps{})) }
