package main
import m "github.com/YoshihideShirai/marionette"
func example() m.Node { return m.Join(m.Button("1", m.ComponentProps{}), m.Button("2", m.ComponentProps{})) }
