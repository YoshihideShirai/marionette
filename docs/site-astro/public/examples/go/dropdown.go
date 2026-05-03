package main
import m "github.com/YoshihideShirai/marionette"
func example() m.Node { return m.Dropdown(m.Button("Menu", m.ComponentProps{}), m.Menu(m.Text("Item 1"), m.Text("Item 2"))) }
