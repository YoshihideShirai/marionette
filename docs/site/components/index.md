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
11. [Feedback Demo](#feedback-demo)

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

- Golden sample: [`button.golden.html`](../../../testdata/golden/button.golden.html)
- Template: [`templates/components/button.tmpl`](../../../templates/components/button.tmpl)

## ComponentInput

### Purpose

Single-value input field for text, date, and similar values.

### Key props

- `name`, `value`, `type`, `placeholder`
- `min`, `max`, `required`, `disabled`
- `Props` (`Variant`, `Size`)

### Visual

<iframe src="./examples/input.html" title="ComponentInput example" style="width:100%;min-height:180px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`input.golden.html`](../../../testdata/golden/input.golden.html)
- Template: [`templates/components/input.tmpl`](../../../templates/components/input.tmpl)

## ComponentSelect

### Purpose

Dropdown input for selecting one option from a list.

### Key props

- `name`
- `options` (`label`, `value`, `selected`)
- `Props` (`Variant`, `Size`)

### Visual

<iframe src="./examples/select.html" title="ComponentSelect example" style="width:100%;min-height:200px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`select.golden.html`](../../../testdata/golden/select.golden.html)
- Template: [`templates/components/select.tmpl`](../../../templates/components/select.tmpl)

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

- Golden sample: [`modal_open.golden.html`](../../../testdata/golden/modal_open.golden.html)
- Template: [`templates/components/modal.html`](../../../templates/components/modal.html)

## ComponentEmptyState

### Purpose

Empty-state UI for no-data and initial states.

### Key props

- `Title`
- `Description`
- `Skeleton`, `Rows`

### Visual

<iframe src="./examples/empty_state.html" title="ComponentEmptyState example" style="width:100%;min-height:260px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`empty_state.golden.html`](../../../testdata/golden/empty_state.golden.html)
- Template: [`templates/components/empty_state.html`](../../../templates/components/empty_state.html)

## ComponentTable

### Purpose

Tabular list UI rendered from column and row definitions.

### Key props

- `Columns` (`Label`, `SortKey`, `Sorted`, `Direction`)
- `Rows`
- `EmptyTitle`, `EmptyDescription`

### Visual

<iframe src="./examples/table.html" title="ComponentTable example" style="width:100%;min-height:300px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`table.golden.html`](../../../testdata/golden/table.golden.html)
- Template: [`templates/components/table.html`](../../../templates/components/table.html)

## ComponentPagination

### Purpose

Pagination controls for previous/next navigation and current page display.

### Key props

- `Page`
- `TotalPages`
- `PrevHref`, `NextHref`

### Visual

<iframe src="./examples/pagination.html" title="ComponentPagination example" style="width:100%;min-height:180px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`pagination.golden.html`](../../../testdata/golden/pagination.golden.html)
- Template: [`templates/components/pagination.tmpl`](../../../templates/components/pagination.tmpl)

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

- Golden sample: [`form_field.golden.html`](../../../testdata/golden/form_field.golden.html)
- Template: [`templates/components/form_field.html`](../../../templates/components/form_field.html)

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

- Golden sample: [`tabs.golden.html`](../../../testdata/golden/tabs.golden.html)
- Template: [`templates/components/tabs.html`](../../../templates/components/tabs.html)

## ComponentBreadcrumb

### Purpose

Path-based navigation that shows the current location hierarchy.

### Key props

- `Items` (`Label`, `Href`, `Active`)
- `AriaLabel`
- `Props.Class`

### Visual

<iframe src="./examples/navigation.html" title="ComponentBreadcrumb example" style="width:100%;min-height:420px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Golden sample: [`breadcrumb.golden.html`](../../../testdata/golden/breadcrumb.golden.html)
- Template: [`templates/components/breadcrumb.html`](../../../templates/components/breadcrumb.html)

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

- Demo source: [`ui/feedback/stories.html`](../../../ui/feedback/stories.html)
- Tokens: [`ui/feedback/tokens.css`](../../../ui/feedback/tokens.css)
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

<iframe src="../../../ui/overlay/stories.html" title="Overlay system demo" style="width:100%;min-height:560px;border:1px solid #e5e7eb;border-radius:8px;"></iframe>

- Demo source: [`ui/overlay/stories.html`](../../../ui/overlay/stories.html)
- Core logic: [`ui/overlay/overlay-core.js`](../../../ui/overlay/overlay-core.js)
- Tokens: [`ui/overlay/tokens.css`](../../../ui/overlay/tokens.css)
- Overlay docs: [`ui/overlay/README.md`](../../../ui/overlay/README.md)
