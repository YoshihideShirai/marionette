package main
import m "github.com/YoshihideShirai/marionette"
func example() m.Node { return m.Navbar(m.Text("Marionette"), m.Button("Docs", m.ComponentProps{}), m.Button("Login", m.ComponentProps{Class:"btn-primary"})) }
