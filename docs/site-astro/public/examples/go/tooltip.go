package main
import m "github.com/YoshihideShirai/marionette"
func example() m.Node { return m.Tooltip("hello", m.Button("Hover", m.ComponentProps{})) }
