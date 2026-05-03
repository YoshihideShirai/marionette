package frontend

import daisy "github.com/YoshihideShirai/marionette/frontend/daisyui"

func Avatar(src, alt, class string) Node  { return daisy.Avatar(src, alt, class) }
func Navbar(start, center, end Node) Node { return daisy.Navbar(start, center, end) }
func Hero(title, description string, actions ...Node) Node {
	return daisy.Hero(title, description, actions...)
}
func Menu(items ...Node) Node                   { return daisy.Menu(items...) }
func Footer(children ...Node) Node              { return daisy.Footer(children...) }
func Drawer(id string, side, content Node) Node { return daisy.Drawer(id, side, content) }
func Stat(title, value, desc string) Node       { return daisy.Stat(title, value, desc) }
func Steps(items ...Node) Node                  { return daisy.Steps(items...) }
func Step(label string, active bool) Node       { return daisy.Step(label, active) }
func Timeline(items ...Node) Node               { return daisy.Timeline(items...) }
func TimelineItem(startLabel, endLabel string, content Node) Node {
	return daisy.TimelineItem(startLabel, endLabel, content)
}
func Collapse(title string, content Node, open bool) Node {
	return daisy.Collapse(title, content, open)
}
func MockupWindow(title string, content Node) Node        { return daisy.MockupWindow(title, content) }
func Kbd(text string) Node                                { return daisy.Kbd(text) }
func Code(text string) Node                               { return daisy.Code(text) }
func Indicator(item, target Node) Node                    { return daisy.Indicator(item, target) }
func Dropdown(trigger, menu Node) Node                    { return daisy.Dropdown(trigger, menu) }
func Tooltip(text string, child Node) Node                { return daisy.Tooltip(text, child) }
func Loading(sizeClass string) Node                       { return daisy.Loading(sizeClass) }
func RadialProgress(value int, sizeClass string) Node     { return daisy.RadialProgress(value, sizeClass) }
func Rating(name string, max int, checked int) Node       { return daisy.Rating(name, max, checked) }
func Range(name string, value int, min int, max int) Node { return daisy.Range(name, value, min, max) }
func Toggle(name string, checked bool) Node               { return daisy.Toggle(name, checked) }
func Join(children ...Node) Node                          { return daisy.Join(children...) }
func Mask(shapeClass string, child Node) Node             { return daisy.Mask(shapeClass, child) }
func Carousel(items ...Node) Node                         { return daisy.Carousel(items...) }
func CarouselItem(id string, child Node) Node             { return daisy.CarouselItem(id, child) }
func ChatBubble(content Node, end bool) Node              { return daisy.ChatBubble(content, end) }
func ThemeController(options ...Node) Node                { return daisy.ThemeController(options...) }
func ThemeControllerOption(theme string, checked bool, className string) Node {
	return daisy.ThemeControllerOption(theme, checked, className)
}
