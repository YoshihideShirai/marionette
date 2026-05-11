# UI Component Selection Guide

English | [日本語](ja/04-ui-components.md)

This document explains how to choose between existing Marionette components, the low-level `frontend/html` HTML API, DaisyUI convenience components, and new components when building Marionette screens. It is intended for admin UIs and internal tools where visual consistency, accessibility, maintainability, and implementation speed all matter.

Check these related documents for the canonical detailed rules.

- [UI Component Guidelines](../ui-component-guidelines.md): Component selection, accessibility, DaisyUI usage, and state-display policies.
- [frontend/ARCHITECTURE.md](../../frontend/ARCHITECTURE.md): Responsibility boundaries under `frontend`, the split between templates and Go implementations, and compatibility rules.
- [Low-level HTML API](../api/03-low-level-html.md): Primitive APIs in `frontend/html`.
- [DaisyUI convenience component API](../api/05-component-apis-daisyui-convenience.md): DaisyUI-based convenience components.

## Basic policy

As a rule, **prefer existing Marionette `frontend` components** and use low-level APIs or application-specific helpers only for the missing parts. Add a new component only when existing APIs cannot preserve the intended meaning, reusability, or accessibility through composition.

Use this order when deciding what to use.

1. Check whether an existing `frontend` component can express the UI.
2. Check whether an existing DaisyUI convenience component naturally fits the UI.
3. If only a small screen-specific gap remains, fill it locally with `frontend/html`.
4. If the same meaningful UI repeats across multiple screens, create an application-level UI helper.
5. Consider a new built-in Marionette component only when the UI is generic enough for Marionette itself and its ownership boundary is clear.

## 1. When to prefer existing components

Use existing components first when any of the following applies.

- **The meaning matches**: The visible UI meaning matches an existing component name such as button, link, table, card, alert, modal, or form row.
- **Accessibility matters**: Attributes and behavior such as `aria-*`, labels, disabled states, focus, and roles are easy to miss when hand-writing markup.
- **You want DaisyUI / Tailwind consistency**: Existing components keep colors, sizes, spacing, and state styles aligned with the visual language.
- **You want tested HTML output**: Existing components may already be covered by compatibility checks or golden tests, so they are safer than direct HTML.
- **The same appearance should be shared across screens**: Lists, details, forms, and feedback patterns are easier to keep consistent by composing existing components.

In particular, prefer existing APIs for screen structure, input, data display, and feedback.

| Purpose | Preferred composition examples |
| --- | --- |
| Page structure | Shell, Navbar, Drawer, Container, Section, Grid |
| Actions | Button, Link, Dropdown, Tabs, Steps |
| Input | Form, FormRow, Input, Select, Checkbox, Toggle, Range |
| Data display | Table, DataFrame, Card, Stats, Badge, Avatar, Progress |
| Feedback | Alert, Toast, Modal, Drawer, Loading, Tooltip |

If an existing component is close but lacks a prop or helper, first check whether a small option addition or application-side composition is enough before creating a new component.

## 2. When to use the low-level `frontend/html` HTML API

`frontend/html` is the low-level API for composing primitives such as `Div`, `Span`, `Button`, `Attr`, `Class`, and `Text` directly. It is flexible, but it also makes the caller responsible for UI semantics, class consistency, and accessibility.

Use it only in cases like these.

- **Filling small gaps around existing components**: Supporting text inside a card, small badge groups inside table cells, or explanatory copy in an empty state.
- **Screen-specific markup with no reuse plan**: Decoration or layout adjustment used by one Page only and not worth turning into a generic component.
- **Building child nodes for existing components**: Content passed into Card, Modal, Alert, Dropdown, and similar components.
- **The semantic HTML primitive is the right abstraction**: `dl` / `dt` / `dd`, `time`, `code`, or a small inline label that does not justify a dedicated component.
- **A local experiment or migration step**: Minimal markup used inside one screen before deciding whether it deserves an API.

Avoid direct low-level API usage when:

- The same markup starts being copied in multiple places.
- The UI needs `aria-*`, keyboard behavior, focus management, or label association.
- The DaisyUI class combination is complex and easy to misuse.
- The UI is no longer screen-specific and has meaning across the app or Marionette itself.

Even when using the low-level API, avoid placing long HTML construction directly in Pages or Actions. Move it into a screen-level render helper or an application-specific helper under `internal/ui`.

## 3. When to use DaisyUI convenience components

DaisyUI convenience components make DaisyUI's class system easier to use from Marionette's Node API. If the UI matches a known DaisyUI pattern, prefer a convenience component over hand-written low-level HTML and class strings.

Use them when these conditions apply.

- **DaisyUI has an official component concept**: Examples include `alert`, `badge`, `card`, `modal`, `drawer`, `dropdown`, `tabs`, `stats`, and `toast`.
- **The variants fit DaisyUI vocabulary**: Examples include `primary`, `secondary`, `success`, `warning`, `error`, `sm`, and `lg`.
- **Visual consistency matters**: Dashboards, admin shells, and settings screens should share the same theme, spacing, and state colors.
- **Existing API props are enough**: The convenience component arguments or variant helpers can express the required appearance and state.
- **You want to hide DaisyUI-specific structure**: Callers should not have to scatter details such as `modal-box`, `toast-end`, or `stats shadow` throughout the app.

However, DaisyUI convenience components are for making DaisyUI components easier to use. Business-specific concepts such as an invoice status card, user-permission badge, or job-run summary should usually be named as application-side helpers instead of being added directly as DaisyUI components.

## 4. Checks before creating a new component

Before creating a new component, check the following.

- **Existing APIs**: Look for close matches in `frontend`, `frontend/daisyui`, the Form APIs, the Data display APIs, and the Overlay / Feedback APIs.
- **Composition first**: Decide whether existing components plus child nodes, or an application-specific helper, are enough.
- **Ownership**: Decide whether this belongs in Marionette itself, in a specific app's `internal/ui`, or as a helper local to one screen.
- **Accessibility**: Decide how labels, roles, aria attributes, focus, keyboard behavior, loading, disabled, and error states are handled.
- **State and API shape**: Decide whether display state, validation errors, loading, empty, selected, and active states should be props or child nodes.
- **DaisyUI dependency**: Decide whether it is a DaisyUI primitive or an app-specific UI that merely uses DaisyUI internally.
- **Compatibility requirements**: Decide whether the component must preserve HTML output compatibility with an existing root API or template-backed component.
- **Testing strategy**: Decide whether to add golden tests, render tests, class-name tests, or checks for accessibility-related attributes.
- **Naming**: Name the component after user-visible or business meaning rather than an HTML tag or CSS class.

When in doubt, start with an application-level helper instead of adding something to Marionette itself. Promote it to a built-in component only after the same meaningful UI has stabilized across multiple screens or apps.

## 5. Recommended compositions by common screen pattern

### List screens

For list screens, separate filters, the table, pagination, row actions, and empty states.

Recommended composition:

- Page: Restore query state and load the list data and total count.
- Layout: Use Container / Section / Card for the list area.
- Filter: Use Form, Input, Select, Button, Tabs, and similar components for search and filtering.
- Data: Use Table or DataFrame, and represent row states with Badge / Status / Progress.
- Action: Prefer Button / Link / Dropdown for per-row detail links, edit buttons, and delete buttons.
- Empty / Loading: Use existing EmptyState, Alert, Loading, or skeleton-style expressions.

Keep low-level HTML limited to small supporting details inside table cells or short help text inside filter forms. If the same filter block appears in multiple lists, move it into an application-specific helper.

### Detail screens

For detail screens, separate the title, primary actions, summary, attribute list, and related data.

Recommended composition:

- Header: Group the page title, breadcrumbs, and primary actions with Section / Button / Link.
- Summary: Use Card, Stats, Badge, Avatar, and Progress for important information near the top.
- Attributes: Use semantic `dl` / `dt` / `dd` markup or an application-specific detail-list helper.
- Related data: Use Table, Timeline, Tabs, or Accordion for history and related resources.
- Feedback: Use Alert / Toast for update results and Modal / confirmation UI for dangerous operations.

Detail screens often repeat business-specific combinations of label, value, and state. If low-level HTML copies start accumulating, create a detail-row or metadata-panel helper under `internal/ui`.

### Input forms

For input forms, make input values, validation errors, submit results, and redraw boundaries explicit.

Recommended composition:

- Page: Load initial values and options, then pass form state to the renderer.
- Form: Use existing input components such as Form / FormRow / Fieldset / Label / Input / Select / Checkbox / Toggle.
- Validation: Render field errors near the relevant FormRow and page-level errors as Alerts.
- Actions: Group submit, cancel, delete, and similar actions in a button group or Card footer.
- Partial update: Even for htmx partial updates, re-render through the same helper from form state.

Forms depend on label-input association, required / disabled / error expressions, and value restoration. Create inputs directly with `frontend/html` only for special controls that existing input components cannot express.

### Dashboards

For dashboards, compose metrics, charts, recent events, and action entry points as cards.

Recommended composition:

- Layout: Use responsive Grid, Section, and Card to control information density.
- Metrics: Use Stats, Card, Badge, Progress, and Chart.
- Trends: Use Chart, Timeline, and Table for time series and history.
- Actions: Put common actions in card headers or footers as Button / Link / Dropdown.
- State: Treat date ranges, selected organizations, and filters as URL query state or shared state.

For DashWind-style admin dashboards, check the DashWind API and existing dashboard-oriented components first. Adding one-off grid or card classes per screen makes theme changes and responsive adjustments harder.

### Modals / toasts / alerts

Choose modals, toasts, and alerts by how interruptive the feedback should be.

Recommended composition:

- Alert: Use for warnings that should remain on the page, validation summaries, save failures, and permission failures.
- Toast: Use for successful saves, lightweight notifications, and feedback that can disappear quickly.
- Modal: Use for destructive-action confirmation, operations that need additional information, or short inputs that should keep the current context.
- Drawer: Use when supplementary information or a detail panel should slide in from the side.
- Loading: Show asynchronous processing near the relevant Button, Card, or Table.

Because these UI elements depend on accessibility, focus, close behavior, and placement consistency, prefer DaisyUI convenience components or the existing Overlay / Feedback APIs. If you use `frontend/html` directly, limit it to local content that existing components cannot express.

## Implementation checklist

- Checked existing `frontend` components and DaisyUI convenience components first.
- Limited `frontend/html` to small, screen-specific supporting markup.
- Moved repeated UI across multiple screens into application-specific helpers.
- Confirmed ownership and compatibility requirements before adding a component to Marionette itself.
- Checked labels, aria attributes, focus, disabled, loading, error, and empty states.
- Followed the recommended patterns for lists, details, forms, dashboards, and feedback.
- Checked [UI Component Guidelines](../ui-component-guidelines.md) and [frontend/ARCHITECTURE.md](../../frontend/ARCHITECTURE.md) when deeper judgment was needed.
