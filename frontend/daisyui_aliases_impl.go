package frontend

import daisy "github.com/YoshihideShirai/marionette/frontend/daisyui"

func TextNode(text string) Node { return daisy.TextNode(text) }
func PrimaryButton(label string, props ComponentProps) Node {
	return daisy.PrimaryButton(label, props)
}
func SecondaryButton(label string, props ComponentProps) Node {
	return daisy.SecondaryButton(label, props)
}
func GhostButton(label string, props ComponentProps) Node { return daisy.GhostButton(label, props) }
func Avatar(src, alt, class string) Node                  { return daisy.Avatar(src, alt, class) }
func Navbar(start, center, end Node) Node                 { return daisy.Navbar(start, center, end) }
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
func ToggleVariant(name string, checked bool, variant string) Node {
	return daisy.ToggleVariant(name, checked, variant)
}
func ToggleWithIcons(name string, checked bool, className string) Node {
	return daisy.ToggleWithIcons(name, checked, className)
}
func Join(children ...Node) Node              { return daisy.Join(children...) }
func Mask(shapeClass string, child Node) Node { return daisy.Mask(shapeClass, child) }
func Carousel(items ...Node) Node             { return daisy.Carousel(items...) }
func CarouselItem(id string, child Node) Node { return daisy.CarouselItem(id, child) }
func ChatBubble(content Node, end bool) Node  { return daisy.ChatBubble(content, end) }
func ThemeController(options ...Node) Node    { return daisy.ThemeController(options...) }
func ThemeControllerOption(theme string, checked bool, className string) Node {
	return daisy.ThemeControllerOption(theme, checked, className)
}
func Countdown(value int) Node                    { return daisy.Countdown(value) }
func Status(colorClass string) Node               { return daisy.Status(colorClass) }
func Dock(items ...Node) Node                     { return daisy.Dock(items...) }
func Fieldset(legend string, fields ...Node) Node { return daisy.Fieldset(legend, fields...) }
func Label(text string) Node                      { return daisy.Label(text) }
func Validator(message string) Node               { return daisy.Validator(message) }
func BrowserMockup(content Node) Node             { return daisy.BrowserMockup(content) }
func PhoneMockup(content Node) Node               { return daisy.PhoneMockup(content) }
func CodeMockup(lines ...string) Node             { return daisy.CodeMockup(lines...) }
func Calendar(content Node) Node                  { return daisy.Calendar(content) }
func Swap(onNode, offNode Node, active bool) Node { return daisy.Swap(onNode, offNode, active) }
func Filter(items ...Node) Node                   { return daisy.Filter(items...) }
func Diff(before, after Node) Node                { return daisy.Diff(before, after) }
func List(items ...Node) Node                     { return daisy.List(items...) }
func Accordion(title string, content Node, open bool) Node {
	return daisy.Accordion(title, content, open)
}
func FAB(icon Node, label string) Node           { return daisy.FAB(icon, label) }
func SpeedDial(trigger Node, items ...Node) Node { return daisy.SpeedDial(trigger, items...) }
func DockItem(child Node, active bool) Node      { return daisy.DockItem(child, active) }
func FilterItem(label string, active bool) Node  { return daisy.FilterItem(label, active) }
func CalendarGrid(days ...Node) Node             { return daisy.CalendarGrid(days...) }
func TextRotate(words []string, animationClass string) Node {
	return daisy.TextRotate(words, animationClass)
}
func Hover3DCard(content Node) Node   { return daisy.Hover3DCard(content) }
func HoverGallery(items ...Node) Node { return daisy.HoverGallery(items...) }

func SelectWithVariants(name string, options []SelectOption, color, size, style string, props ComponentProps) Node {
	return daisy.SelectWithVariants(name, options, color, size, style, props)
}
func ProgressWithVariant(value, max float64, label, color string, props ComponentProps) Node {
	return daisy.ProgressWithVariant(value, max, label, color, props)
}
func BadgeWithVariant(label, color, size, style string, props ComponentProps) Node {
	return daisy.BadgeWithVariant(label, color, size, style, props)
}
func RangeWithVariants(name string, value int, min int, max int, color string, size string) Node {
	return daisy.RangeWithVariants(name, value, min, max, color, size)
}
func RatingWithVariants(name string, max int, checked int, size string, half bool, allowClear bool) Node {
	return daisy.RatingWithVariants(name, max, checked, size, half, allowClear)
}
func ToastWithPlacement(children []Node, horizontal, vertical, className string) Node {
	return daisy.ToastWithPlacement(children, horizontal, vertical, className)
}
func TooltipWithVariants(text string, child Node, placement string, color string, open bool) Node {
	return daisy.TooltipWithVariants(text, child, placement, color, open)
}
func TableWithVariants(headers []string, rows [][]Node, zebra bool, pinRows bool, pinCols bool, size string) Node {
	return daisy.TableWithVariants(headers, rows, zebra, pinRows, pinCols, size)
}
func ModalWithPlacement(props ModalProps, placement string) Node {
	return daisy.ModalWithPlacement(props, placement)
}
func StepsWithVariants(items []Node, direction string, color string, className string) Node {
	return daisy.StepsWithVariants(items, direction, color, className)
}
func TimelineWithDirection(items []Node, direction string, compact bool, snapIcon bool, className string) Node {
	return daisy.TimelineWithDirection(items, direction, compact, snapIcon, className)
}
func LoadingWithVariants(kind string, size string) Node { return daisy.LoadingWithVariants(kind, size) }
func StatusWithVariants(color string, size string) Node { return daisy.StatusWithVariants(color, size) }
func ToggleWithVariants(name string, checked bool, color string, size string) Node {
	return daisy.ToggleWithVariants(name, checked, color, size)
}
func SwapWithVariants(onNode, offNode Node, active bool, rotate bool, flip bool) Node {
	return daisy.SwapWithVariants(onNode, offNode, active, rotate, flip)
}
func JoinWithDirection(direction string, children ...Node) Node {
	return daisy.JoinWithDirection(direction, children...)
}
func DropdownWithPlacement(trigger, menu Node, placement string) Node {
	return daisy.DropdownWithPlacement(trigger, menu, placement)
}
