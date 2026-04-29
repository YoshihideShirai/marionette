# Components Gallery

This page lists Marionette UI components. The target component set is managed in one place via `docs/site/components/_index.json`.

- Source of truth: [`_index.json`](./_index.json)
- Golden source files: `testdata/golden/*.golden.html`
- Generated sample output: `docs/site/components/examples/*.html`

## Contents

1. [ComponentButton](#componentbutton)
2. [ComponentInput](#componentinput)
3. [ComponentSelect](#componentselect)
4. [ComponentModal](#componentmodal)
5. [ComponentEmptyState](#componentemptystate)
6. [ComponentTable](#componenttable)
7. [ComponentPagination](#componentpagination)
8. [ComponentFormField](#componentformfield)
9. [ComponentTabs](#componenttabs)
10. [ComponentBreadcrumb](#componentbreadcrumb)
11. [ComponentTextarea](#componenttextarea)
12. [ComponentCheckbox](#componentcheckbox)
13. [ComponentRadioGroup](#componentradiogroup)
14. [ComponentSwitch](#componentswitch)
15. [ComponentStack](#componentstack)
16. [ComponentGrid](#componentgrid)
17. [ComponentSplit](#componentsplit)
18. [ComponentPageHeader](#componentpageheader)
19. [ComponentContainer](#componentcontainer)
20. [ComponentCard](#componentcard)
21. [ComponentSection](#componentsection)
22. [Feedback Demo](#feedback-demo)

---

## ComponentButton

### Purpose

Clickable button for primary and secondary actions.

### Key props

- `Variant` (`primary`, `secondary`, `ghost`, `error`, etc.)
- `Size` (`sm`, `md`, `lg`)
- `Disabled`

### Visual

<iframe src="./examples/button.html" title="ComponentButton example" style="width:100%;min-height:160px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`button.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/button.golden.html)
- Template: [`templates/components/button.tmpl`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/button.tmpl)

## ComponentInput

### Purpose

Single-value input field for text, date, and similar values.

### Key props

- `name`, `value`, `type`, `placeholder`
- `min`, `max`, `required`, `disabled`
- `Props` (`Variant`, `Size`)

### Visual

<iframe src="./examples/input.html" title="ComponentInput example" style="width:100%;min-height:180px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`input.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/input.golden.html)
- Template: [`templates/components/input.tmpl`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/input.tmpl)

## ComponentSelect

### Purpose

Dropdown input for selecting one option from a list.

### Key props

- `name`
- `options` (`label`, `value`, `selected`)
- `Props` (`Variant`, `Size`)

### Visual

<iframe src="./examples/select.html" title="ComponentSelect example" style="width:100%;min-height:200px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`select.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/select.golden.html)
- Template: [`templates/components/select.tmpl`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/select.tmpl)

## ComponentModal

### Purpose

Overlay dialog for confirmations and detail views.

### Key props

- `Title`
- `Body`
- `Actions`
- `Open`

### Visual

<iframe src="./examples/modal.html" title="ComponentModal example" style="width:100%;min-height:320px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`modal_open.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/modal_open.golden.html)
- Template: [`templates/components/modal.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/modal.html)

## ComponentEmptyState

### Purpose

Empty-state UI for no-data and initial states.

### Key props

- `Title`
- `Description`
- `Skeleton`, `Rows`

### Visual

<iframe src="./examples/empty_state.html" title="ComponentEmptyState example" style="width:100%;min-height:260px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`empty_state.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/empty_state.golden.html)
- Template: [`templates/components/empty_state.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/empty_state.html)

## ComponentTable

### Purpose

Tabular list UI rendered from column and row definitions.

### Key props

- `Columns` (`Label`, `SortKey`, `Sorted`, `Direction`)
- `Rows`
- `EmptyTitle`, `EmptyDescription`

### Visual

<iframe src="./examples/table.html" title="ComponentTable example" style="width:100%;min-height:300px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`table.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/table.golden.html)
- Template: [`templates/components/table.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/table.html)

## ComponentPagination

### Purpose

Pagination controls for previous/next navigation and current page display.

### Key props

- `Page`
- `TotalPages`
- `PrevHref`, `NextHref`

### Visual

<iframe src="./examples/pagination.html" title="ComponentPagination example" style="width:100%;min-height:180px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`pagination.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/pagination.golden.html)
- Template: [`templates/components/pagination.tmpl`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/pagination.tmpl)

## ComponentFormField

### Purpose

Form wrapper that combines label, input, hint text, and error text.

### Key props

- `Label`
- `Required`
- `Hint`
- `Error`
- `Child` (nested input/select, etc.)

### Visual

<iframe src="./examples/form_field.html" title="ComponentFormField example" style="width:100%;min-height:260px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`form_field.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/form_field.golden.html)
- Template: [`templates/components/form_field.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/form_field.html)

---

## ComponentTabs

### Purpose

Tabbed navigation for related views and sections.

### Key props

- `Items` (`Label`, `Href`, `Active`, `Disabled`)
- `AriaLabel`
- `Props.Class`

### Visual

<iframe src="./examples/navigation.html" title="ComponentTabs example" style="width:100%;min-height:420px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`tabs.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/tabs.golden.html)
- Template: [`templates/components/tabs.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/tabs.html)

## ComponentBreadcrumb

### Purpose

Path-based navigation that shows the current location hierarchy.

### Key props

- `Items` (`Label`, `Href`, `Active`)
- `AriaLabel`
- `Props.Class`

### Visual

<iframe src="./examples/navigation.html" title="ComponentBreadcrumb example" style="width:100%;min-height:420px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`breadcrumb.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/breadcrumb.golden.html)
- Template: [`templates/components/breadcrumb.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/breadcrumb.html)

---

## ComponentTextarea

### Purpose

Multiline text input with shared component styling options.

### Key props

- `name`, `value`
- `Placeholder`, `Rows`, `Required`
- `Props` (`Variant`, `Size`, `Disabled`, `Class`)

- Golden sample: [`textarea.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/textarea.golden.html)
- Template: [`templates/components/textarea.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/textarea.html)

## ComponentCheckbox

### Purpose

Single boolean selection input with label.

### Key props

- `Name`, `Value`, `Label`
- `Checked`
- `Props` (`Size`, `Disabled`, `Class`)

- Golden sample: [`checkbox.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/checkbox.golden.html)
- Template: [`templates/components/checkbox.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/checkbox.html)

## ComponentRadioGroup

### Purpose

Exclusive selection group rendered from item definitions.

### Key props

- `Name`, `AriaLabel`
- `Items` (`Label`, `Value`, `Checked`, `Disabled`)
- `Props` (`Size`, `Disabled`, `Class`)

- Golden sample: [`radio_group.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/radio_group.golden.html)
- Template: [`templates/components/radio_group.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/radio_group.html)

## ComponentSwitch

### Purpose

Switch-style boolean control for settings and toggles.

### Key props

- `Name`, `Value`, `Label`
- `Checked`
- `Props` (`Size`, `Disabled`, `Class`)

- Golden sample: [`switch.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/switch.golden.html)
- Template: [`templates/components/switch.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/switch.html)

---

## ComponentStack

### Purpose

Flex stack for vertical and horizontal layout rhythm.

### Key props

- `Direction` (`vertical`, `horizontal`)
- `Gap` (`none`, `xs`, `sm`, `md`, `lg`, `xl`)
- `Align`, `Justify`, `Wrap`
- `Props.Class`

### Visual

<iframe src="./examples/layout.html" title="Layout components example" style="width:100%;min-height:520px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`stack.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/stack.golden.html)
- Template: [`templates/components/stack.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/stack.html)

## ComponentGrid

### Purpose

Responsive grid for cards, summaries, and repeated panels.

### Key props

- `Columns` (`1`, `2`, `3`, `4`)
- `MinColumnWidth` (`sm`, `md`, `lg`)
- `Gap`
- `Props.Class`

- Golden sample: [`grid.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/grid.golden.html)
- Template: [`templates/components/grid.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/grid.html)

## ComponentSplit

### Purpose

Responsive main/aside layout for admin workspaces.

### Key props

- `Main`, `Aside`
- `AsideWidth` (`sm`, `md`, `lg`)
- `ReverseOnMobile`
- `Gap`, `Props.Class`

- Golden sample: [`split.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/split.golden.html)
- Template: [`templates/components/split.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/split.html)

## ComponentPageHeader

### Purpose

Page heading layout with title, description, and actions.

### Key props

- `Title`, `Description`
- `Actions`
- `Props.Class`

- Golden sample: [`page_header.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/page_header.golden.html)
- Template: [`templates/components/page_header.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/page_header.html)

## ComponentContainer

### Purpose

Width and padding wrapper for page-level content.

### Key props

- `MaxWidth` (`sm`, `md`, `lg`, `full`)
- `Padding` (`none`, `sm`, `md`, `lg`)
- `Centered`
- `Props.Class`

- Golden sample: [`container.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/container.golden.html)
- Template: [`templates/components/container.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/container.html)

## ComponentCard

### Purpose

Card surface with optional header actions.

### Key props

- `Title`, `Description`
- `Actions`
- `Props.Class`

- Golden sample: [`card.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/card.golden.html)
- Template: [`templates/components/card.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/card.html)

## ComponentSection

### Purpose

Unframed section wrapper with consistent header spacing.

### Key props

- `Title`, `Description`
- `Actions`
- `Props.Class`

- Golden sample: [`section.golden.html`](https://github.com/YoshihideShirai/marionette/blob/main/testdata/golden/section.golden.html)
- Template: [`templates/components/section.html`](https://github.com/YoshihideShirai/marionette/blob/main/templates/components/section.html)

---


## Feedback Demo

### Purpose

Unified demo for `Toast`, `Alert`, `EmptyState`, and `Skeleton` added under `ui/feedback/`.

- Variant: `success` / `info` / `warning` / `error`
- Size: `sm` / `md` / `lg`
- Accessibility: `role`, `aria-live`, `aria-busy`
- Tokens: `ui/feedback/tokens.css`

### Visual

<iframe src="./examples/feedback.html" title="Feedback components demo" style="width:100%;min-height:420px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Demo source: [`ui/feedback/stories.html`](https://github.com/YoshihideShirai/marionette/blob/main/ui/feedback/stories.html)
- Tokens: [`ui/feedback/tokens.css`](https://github.com/YoshihideShirai/marionette/blob/main/ui/feedback/tokens.css)
- Templates: `templates/components/{toast,alert,empty_state,skeleton}.html`

## Optional demo screenshot

Place a demo screenshot from `cmd/marionette/main.go` under `docs/site/assets/` when browser-capture tooling is available.


## Overlay System Demo

### Purpose

A demo that validates `Modal`, `Drawer`, `Popover`, and `Tooltip` in `ui/overlay/`, powered by shared `overlay-core` logic.

- Esc-key close
- Configurable backdrop-click close
- Focus trap and focus restoration
- Unified portal target (`#ui-overlay-root`)
- Reusable z-index / scroll-lock rules

### Visual

> Live iframe preview is omitted in GitHub Pages because `ui/overlay/` is outside the published `docs/site/` artifact. Use the links below to view source assets.

- Demo source: [`ui/overlay/stories.html`](https://github.com/YoshihideShirai/marionette/blob/main/ui/overlay/stories.html)
- Core logic: [`ui/overlay/overlay-core.js`](https://github.com/YoshihideShirai/marionette/blob/main/ui/overlay/overlay-core.js)
- Tokens: [`ui/overlay/tokens.css`](https://github.com/YoshihideShirai/marionette/blob/main/ui/overlay/tokens.css)
- Overlay docs: [`ui/overlay/README.md`](https://github.com/YoshihideShirai/marionette/blob/main/ui/overlay/README.md)
