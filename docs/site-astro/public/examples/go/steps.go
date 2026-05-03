package main
import m "github.com/YoshihideShirai/marionette"
func example() m.Node { return m.Steps(m.Step("A", true), m.Step("B", false)) }
