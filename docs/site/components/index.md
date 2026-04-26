# Components Gallery

This page lists Marionette UI components. The target component set is managed in one place via `docs/site/components/_index.json`.

- Source of truth: [`_index.json`](./_index.json)
- Golden source files: `internal/marionette/testdata/golden/*.golden.html`
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

- Golden sample: [`button.golden.html`](../../../internal/marionette/testdata/golden/button.golden.html)
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

- Golden sample: [`input.golden.html`](../../../internal/marionette/testdata/golden/input.golden.html)
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

- Golden sample: [`select.golden.html`](../../../internal/marionette/testdata/golden/select.golden.html)
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

- Golden sample: [`modal_open.golden.html`](../../../internal/marionette/testdata/golden/modal_open.golden.html)
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

- Golden sample: [`empty_state.golden.html`](../../../internal/marionette/testdata/golden/empty_state.golden.html)
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

- Golden sample: [`table.golden.html`](../../../internal/marionette/testdata/golden/table.golden.html)
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

- Golden sample: [`pagination.golden.html`](../../../internal/marionette/testdata/golden/pagination.golden.html)
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

- Golden sample: [`form_field.golden.html`](../../../internal/marionette/testdata/golden/form_field.golden.html)
- Template: [`templates/components/form_field.html`](../../../templates/components/form_field.html)

---

## Optional demo screenshot

Place a demo screenshot from `cmd/marionette/main.go` under `docs/site/assets/` when browser-capture tooling is available.
