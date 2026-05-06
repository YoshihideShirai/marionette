## 6.5 daisyUI convenience components

These helpers are re-exported from `frontend` and also available from
`frontend/daisyui`. They map directly to common daisyUI component structures and
usually accept already-renderable `Node` children.

### Typography and buttons
- `TextNode(text string) Node`
  - Renders text in a `span`.
- `PrimaryButton(label string, props ComponentProps) Node`
- `SecondaryButton(label string, props ComponentProps) Node`
- `GhostButton(label string, props ComponentProps) Node`
  - Convenience wrappers around `Button` that add `btn-primary`,
    `btn-secondary`, or `btn-ghost` classes.
- `ButtonWithVariants(label, variant, size, style string, props ComponentProps) Node`
  - Adds `btn-<variant>`, `btn-<size>`, and `btn-<style>` when non-empty.

### Navigation, menus, and layout primitives
- `Navbar(start, center, end Node) Node`
  - Emits `navbar-start`, `navbar-center`, and `navbar-end` slots.
- `Menu(items ...Node) Node`
  - Wraps non-`li` items in `li` elements.
- `Footer(children ...Node) Node`
- `Drawer(id string, side, content Node) Node`
  - Uses `id` for the drawer toggle input and matching overlay label.
- `Dropdown(trigger, menu Node) Node`
- `DropdownWithPlacement(trigger, menu Node, placement string) Node`
  - Adds `dropdown-<placement>` when placement is non-empty.
- `Join(children ...Node) Node`
- `JoinWithDirection(direction string, children ...Node) Node`
  - Adds `join-<direction>` when direction is non-empty.

### Hero, avatar, card-style, and mockup helpers
- `Hero(title, description string, actions ...Node) Node`
- `Avatar(src, alt, class string) Node`
  - Uses daisyUI avatar markup with the provided wrapper class around `img`.
- `Stat(title, value, desc string) Node`
- `BrowserMockup(content Node) Node`
- `PhoneMockup(content Node) Node`
- `MockupWindow(title string, content Node) Node`
  - `title` is currently accepted for API symmetry but is not rendered.
- `CodeMockup(lines ...string) Node`
  - Renders each line as a mock terminal row with `data-prefix="$"`.

### Disclosure, steps, timeline, tabs, and carousel
- `Accordion(title string, content Node, open bool) Node`
  - Uses radio-backed collapse markup with `name="accordion"`.
- `Collapse(title string, content Node, open bool) Node`
  - Uses checkbox-backed collapse markup.
- `Steps(items ...Node) Node`
- `Step(label string, active bool) Node`
  - `active` adds `step-primary`.
- `StepsWithVariants(items []Node, direction, color, className string) Node`
  - Adds `steps-<direction>` and `className`; when `color` is non-empty,
    existing `step` children receive `step-<color>`.
- `Timeline(items ...Node) Node`
- `TimelineItem(startLabel, endLabel string, content Node) Node`
- `TimelineWithDirection(items []Node, direction string, compact bool, snapIcon bool, className string) Node`
  - Adds `timeline-<direction>`, `timeline-compact`, `timeline-snap-icon`, and
    `className` when requested.
- `TabsWithVariants(items []TabsItem, style, placement, size, className string) Node`
  - Adds `tabs-<style>`, `tabs-<placement>`, `tabs-<size>`, and `className`.
  - `TabsItem.Active` sets `tab-active` and `aria-selected="true"`.
  - `TabsItem.Disabled` sets `tab-disabled` and `disabled`.
- `Carousel(items ...Node) Node`
- `CarouselItem(id string, child Node) Node`

### Inputs and form adjuncts
- `Range(name string, value int, min int, max int) Node`
- `RangeWithVariants(name string, value int, min int, max int, color, size string) Node`
  - Adds `range-<color>` and `range-<size>` when non-empty.
- `Rating(name string, max int, checked int) Node`
- `RatingWithVariants(name string, max int, checked int, size string, half bool, allowClear bool) Node`
  - Adds `rating-<size>` and `rating-half` when requested.
  - `allowClear` prepends a hidden value `0` radio option.
- `Toggle(name string, checked bool) Node`
- `ToggleVariant(name string, checked bool, variant string) Node`
  - Recognizes `primary`, `secondary`, `accent`, `neutral`, `info`,
    `success`, `warning`, and `danger`/`error`.
- `ToggleWithVariants(name string, checked bool, color string, size string) Node`
  - Adds `toggle-<color>` and `toggle-<size>` when non-empty.
- `ToggleWithIcons(name string, checked bool, className string) Node`
- `ThemeController(options ...Node) Node`
- `ThemeControllerOption(theme string, checked bool, className string) Node`
  - Emits a radio input named `theme-buttons` with `theme-controller` and the
    provided theme value.
- `Fieldset(legend string, fields ...Node) Node`
- `Label(text string) Node`
- `Validator(message string) Node`

### Feedback, overlays, and visual effects
- `Tooltip(text string, child Node) Node`
- `TooltipWithVariants(text string, child Node, placement, color string, open bool) Node`
  - Adds `tooltip-<placement>`, `tooltip-<color>`, and `tooltip-open`.
- `ToastWithPlacement(children []Node, horizontal, vertical, className string) Node`
  - Adds `toast-<horizontal>`, `toast-<vertical>`, and `className`.
- `Loading(sizeClass string) Node`
  - Defaults to `loading loading-spinner`; appends `sizeClass` when non-empty.
- `LoadingWithVariants(kind string, size string) Node`
  - Blank `kind` defaults to `spinner`; adds `loading-<size>` when non-empty.
- `RadialProgress(value int, sizeClass string) Node`
- `Status(colorClass string) Node`
- `StatusWithVariants(color string, size string) Node`
  - Adds `status-<color>` and `status-<size>` when non-empty.
- `Swap(onNode, offNode Node, active bool) Node`
- `SwapWithVariants(onNode, offNode Node, active bool, rotate bool, flip bool) Node`
  - Adds `swap-active`, `swap-rotate`, and `swap-flip` according to flags.
- `FAB(icon Node, label string) Node`
  - Uses `label` as the button `aria-label`.
- `SpeedDial(trigger Node, items ...Node) Node`

### Data display, media, and utility helpers
- `List(items ...Node) Node`
  - Wraps non-`li` items and adds `list-row` to list entries.
- `Table(headers []string, rows ...[]Node) Node`
- `TableWithVariants(headers []string, rows [][]Node, zebra bool, pinRows bool, pinCols bool, size string) Node`
  - Adds `table-zebra`, `table-pin-rows`, `table-pin-cols`, and
    `table-<size>` according to arguments.
- `Kbd(text string) Node`
- `Code(text string) Node`
- `Indicator(item, target Node) Node`
- `Mask(shapeClass string, child Node) Node`
- `ChatBubble(content Node, end bool) Node`
  - `end` switches from `chat-start` to `chat-end`.
- `Countdown(value int) Node`
- `Dock(items ...Node) Node`
- `DockItem(child Node, active bool) Node`
  - `active` adds `dock-active`.
- `Filter(items ...Node) Node`
- `FilterItem(label string, active bool) Node`
  - Uses radio input markup with `aria-label` set to `label`.
- `Calendar(content Node) Node`
- `CalendarGrid(days ...Node) Node`
- `Diff(before, after Node) Node`
- `TextRotate(words []string, animationClass string) Node`
- `Hover3DCard(content Node) Node`
- `HoverGallery(items ...Node) Node`
