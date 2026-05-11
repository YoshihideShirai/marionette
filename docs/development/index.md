# Marionette Development Guide

English | [日本語](ja/index.md)

This guide is the entry point for the standard workflow for building admin UIs and internal tools with Marionette. Use it to understand the overall flow first, then jump to the existing detailed documentation that matches the implementation phase you are in.

## How to read the docs

### Read first

- [README.md](../../README.md): Review Marionette's purpose, core ideas, demos, and the documentation entry points.
- [Project Structure Guide](01-project-structure.md): Choose a user-app layout for pages, actions, state, UI helpers, and assets.
- [Routing / Pages / Actions Guide](02-routing-pages-actions.md): Decide how to split URLs, screens, partial updates, and action handlers.
- [UI Component Selection Guide](04-ui-components.md): Choose between existing components, low-level HTML, DaisyUI helpers, and new components.
- [Forms / Validation Guide](05-forms-validation.md): Design form state, server-side validation, validation redisplay, and success responses.
- [Data Tables / Charts Guide](06-data-tables-charts.md): Design list tables, shared query state, paging, sorting, and chart-linked filtering.
- [Errors / Flash / Feedback Guide](08-errors-flash-feedback.md): Separate user-facing and logged errors, and choose inline errors, alerts, toasts, flash, and retry paths.
- [Long-running Jobs Design Guide](09-long-running-jobs.md): Design job models, progress / toast / empty_state usage, timeout, retry, cache TTL, and production storage decisions.
- [Security / Authorization Guide](10-security-authz.md): Handle authenticated users, server-side authorization checks, dangerous-operation confirmations, and audit logs.
- [AI Assisted Development Guide](12-ai-assisted-development.md): Use request templates for AI-assisted screen, Action, form, and table-chart changes.
- [State Management Guide](../state-management.md): Understand the basic policy for keeping pages, actions, and state on the Go side.
- [API documentation](../api/): Start here for the main `backend`, `frontend`, and `html` APIs.

### Read when needed

- [UI Component Guidelines](../ui-component-guidelines.md): Check how to choose components, keep accessibility in mind, and work with DaisyUI.
- [frontend/ARCHITECTURE.md](../../frontend/ARCHITECTURE.md): Review the responsibility boundaries in the `frontend` package, the split between templates and Go implementations, and compatibility rules.
- [Form APIs](../api/04-form-apis.md): Use this when composing form rows, inputs, validation errors, and selection controls.
- [Data display component APIs](../api/05-component-apis-data-display.md): Review tables, stats, avatars, progress, and other display components.
- [Overlay / Feedback APIs](../api/05-component-apis-overlay-feedback.md): Review modals, drawers, toasts, alerts, and other feedback components.
- [Forms / Validation Guide](05-forms-validation.md): Use this when implementing submit flows, server-side validation, or form error redisplay.
- [Data Tables / Charts Guide](06-data-tables-charts.md): Use this when implementing searchable, pageable, sortable tables or linked table/chart dashboards.
- [Errors / Flash / Feedback Guide](08-errors-flash-feedback.md): Use this when designing errors, flash messages, toasts, alerts, long-running work, or retry UX.
- [Long-running Jobs Design Guide](09-long-running-jobs.md): Use this when expanding the Heavy Job Template for production use.
- [Security / Authorization Guide](10-security-authz.md): Use this when designing authentication-derived context, Page / Action authorization checks, dangerous-operation confirmations, or audit logs.
- [AI Assisted Development Guide](12-ai-assisted-development.md): Use this when preparing AI requests for screens, Actions, forms, or table-chart coordination.
- [DashWind API](../api/08-dashwind.md): Use this when building DashWind-style dashboards or admin layouts.

### Check before implementation

- [frontend/ARCHITECTURE.md](../../frontend/ARCHITECTURE.md): Before adding a new UI helper or component, confirm which package owns the responsibility.
- [UI Component Guidelines](../ui-component-guidelines.md): Confirm that new UI follows the existing accessibility, visual, state, and component-selection guidance.
- [State Management Guide](../state-management.md): Decide where URL, session, form input, and temporary UI state should live.
- [API documentation](../api/): Check whether existing APIs already cover the use case so you can avoid unnecessary wrappers or duplicate implementations.
- [Security / Authorization Guide](10-security-authz.md): Confirm that UI visibility controls are backed by server-side authorization and audit logging.
- [AI Assisted Development Guide](12-ai-assisted-development.md): Prepare the target URL, Page, Action, State, partial-update, and error-display contract before asking AI for implementation help.

## Standard development flow

### 1. Decide the project structure

Start by separating the application entry point, pages, state, actions, and display components. Small apps can stay in a single package, but once the number of screens grows, a responsibility-based layout like this keeps the code easier to navigate.

```text
cmd/<app>/main.go        # App startup and route registration
internal/<app>/pages/    # Page functions and screen-level composition
internal/<app>/actions/  # POST/event handling, state updates, redirects
internal/<app>/state/    # Queries, filters, form values, session-derived state
internal/<app>/ui/       # Small app-specific UI composition helpers
```

When changing Marionette's built-in components, the responsibility boundaries under `frontend` are important. Do not mix thin aliases for existing components, DaisyUI-specific implementations, and low-level HTML node generation. Check [frontend/ARCHITECTURE.md](../../frontend/ARCHITECTURE.md) before implementing those changes.

### 2. Design pages and actions

For detailed criteria, see the [Routing / Pages / Actions Guide](02-routing-pages-actions.md). In Marionette, it is usually easiest to treat screen reads as pages and user-triggered changes as actions.

- A page restores state from the request, loads data, and returns HTML composed with `frontend` components.
- An action handles form submissions, button clicks, filter changes, and similar events. It owns validation, state updates, persistence, and partial-update responses.
- Shared input values and filter conditions between pages and actions should be collected in explicit state types.

If you first write down what each page displays and which operations can change it, the htmx partial-update targets become easier to choose. For API basics, see the [API documentation](../api/) and [Context API](../api/02-context.md).

### 3. Decide state management

Place state according to its lifetime and sharing scope.

- URL query: Search terms, page numbers, sort order, and other state that should survive reloads or be shareable.
- Form values: Values being edited or submitted. Re-render them when validation fails.
- Server-side session or durable storage: Authentication, user settings, and business data that must be kept longer.
- Temporary UI state: Short-lived screen state such as modal visibility or toast display.

Use the [State Management Guide](../state-management.md) as the baseline. When multiple widgets, such as a table and a chart, share the same filters, move those conditions into a shared state type so changes stay localized.

### 4. Choose UI components

Build UI with existing `frontend` APIs and DaisyUI-based components first. Before writing low-level HTML directly, check whether existing layout, input, data display, and feedback components can express the screen.

- Screen structure: Shell, Navbar, Drawer, Container, Section, Grid, and similar components.
- Actions: Button, Link, Dropdown, Tabs, Steps, and similar controls.
- Information display: Card, Alert, Badge, Stats, Timeline, and similar components.
- Business UI: Table, DataFrame, Chart, FormRow, and input components.

For practical selection criteria by screen pattern, see the [UI Component Selection Guide](04-ui-components.md). For accessibility and visual consistency details, see the [UI Component Guidelines](../ui-component-guidelines.md). If you are adding a new built-in component, confirm the implementation location in [frontend/ARCHITECTURE.md](../../frontend/ARCHITECTURE.md).

### 5. Build forms and validation

Design forms as one flow that includes input values, validation results, and error messages for re-rendering.

1. Put initial values into state in the page.
2. Use FormRow and input components to align labels, descriptions, required markers, and error areas.
3. Validate submitted values in the action.
4. If validation fails, return the same state with values and errors for re-rendering.
5. On success, save data, show a notification, redirect, or return a partial update.

See [Form APIs](../api/04-form-apis.md) for form component behavior, and use the [API documentation](../api/) entry point to find related sections before adding custom wrappers.

### 6. Design tables and charts

For tables and charts, design filters, aggregation granularity, sorting, and pagination as state before focusing on the rendered data.

- For tables, consider column definitions, row actions, empty states, loading states, and pagination together.
- For charts, clarify aggregation units, time ranges, legends, and click-to-filter behavior.
- When tables and charts respond to the same conditions, render both from shared state.

Review display components in the [Data display component APIs](../api/05-component-apis-data-display.md). If you are using DashWind admin patterns, also see the [DashWind API](../api/08-dashwind.md).

### 7. Decide error display and notifications

Show errors and notifications at a level that helps users understand what to do next.

- Use FormRow / FieldError for field-level input errors.
- Use Alert or EmptyState for whole-screen failures.
- Use Toast / Flash for successful saves and lightweight failures.
- For unrecoverable failures, log the details server-side and show only a safe message in the UI.

For available components, see the [Overlay / Feedback APIs](../api/05-component-apis-overlay-feedback.md) and [Flash APIs](../api/06-flash-apis.md). Also check the [UI Component Guidelines](../ui-component-guidelines.md) for notification styling and accessibility.

### 8. Add security and authorization checks

Security-sensitive screens should validate both display access and operation execution on the server side.

- Restore the authenticated actor from session or token middleware and make it available through `Context`.
- Check route, tenant, organization, project, and target-resource permissions before rendering Pages.
- Re-check actor, action, and target authorization inside Actions even when the UI hides unavailable buttons.
- Add confirmation UI for dangerous operations such as deletion, approval, and re-run.
- Record audit logs for important operation attempts, including denied and failed attempts.

See the [Security / Authorization Guide](10-security-authz.md) for detailed criteria and audit-log fields.

### 9. Prepare tests and debugging

Marionette apps become safer to change when page rendering and action state updates are easy to verify with Go tests.

- Test that page functions produce expected HTML fragments from representative state values.
- Test actions separately for success, validation errors, authorization errors, and persistence errors.
- When adding components, follow the existing rendering-test or golden-test approach.
- For htmx partial updates, verify target element IDs, response boundaries, and error re-rendering.

When changing built-in UI implementation, follow the compatibility policy in [frontend/ARCHITECTURE.md](../../frontend/ARCHITECTURE.md) and add the necessary rendering tests. In application code, keeping state types and actions small makes failures easier to reproduce and debug.

## Pre-implementation checklist

- Page / Action / State responsibilities are separated.
- You classified which state belongs in the URL, the form, the session, or the database.
- You checked whether existing `frontend` components and APIs can cover the use case.
- You decided how to separate form errors, success notifications, and whole-screen errors.
- You decided whether table and chart filters should be represented as shared state.
- Any added UI follows the [UI Component Guidelines](../ui-component-guidelines.md) and [frontend/ARCHITECTURE.md](../../frontend/ARCHITECTURE.md).
- You identified which pages, actions, and component output need tests.
- Sensitive Pages and Actions include authorization checks and audit logs.
