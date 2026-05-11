# AI Assisted Development Guide

English | [日本語](ja/12-ai-assisted-development.md)

This document provides request templates for asking AI to help develop Marionette applications. Give the AI not only what to build, but also the expected Page, Action, State, partial-update behavior, error display, and risks that require human review.

Before asking AI to implement a change, also review the related development guides.

- [Routing / Pages / Actions Guide](02-routing-pages-actions.md): Responsibility boundaries for URLs, Pages, Actions, and partial updates.
- [Forms / Validation Guide](05-forms-validation.md): Form state, validation, and error re-rendering.
- [Data Tables / Charts Guide](06-data-tables-charts.md): Tables, charts, and shared filter state.
- [Security / Authorization Guide](10-security-authz.md): Authorization, dangerous operations, and audit logs.

## Basic policy

Write AI requests by separating the following information instead of mixing it into one vague paragraph.

- Purpose: What the user should be able to view or operate.
- Target: URL, Page name, Action name, and related files.
- State: Which data belongs in the URL query, form state, session, or database.
- Update behavior: Whether to use a full page reload, redirect, partial update, toast, or flash message.
- Failure behavior: How to display validation errors, authorization errors, persistence errors, and external API errors.
- Constraints: Items that humans must review, such as authorization, audit logs, destructive operations, and production settings.

Do not rely only on phrases such as “make it nice” or “follow the existing style.” Include concrete existing file names, function names, display states, and testing expectations.

## 1. Request template for adding a new screen

```text
Add a new screen to a Marionette application.

Purpose:
- <Who uses this screen, and what should they view or operate?>

Target URL:
- GET <Example: /admin/projects/:projectID/members>

Page name:
- <Example: ProjectMembersPage>

Displayed content:
- <Example: project name, member list, roles, invitation status, last login time>
- <Also describe the empty state, loading state, and display for unauthorized users.>

State location:
- URL query: <Example: q, role, page, sort>
- Page state struct: <Example: ProjectMembersPageState>
- session / DB: <Persistent data required for rendering the screen>

Action names:
- <List Actions called from this screen, if any. Example: InviteMember, RemoveMember>

Expected partial update:
- <Example: update only #members-table when filters change>
- <Example: update #members-table and a toast after a successful invitation>

Error display:
- <Example: show an Alert when list loading fails, a 403-style message when unauthorized, and field errors inside the invitation form>

Implementation policy:
- Prefer existing frontend components.
- Keep Page / Action / State / UI helper responsibilities separate.
- Add the necessary tests.

Human review items:
- Authorization, audit logs, destructive operations, and production settings will be reviewed after implementation.
```

### Notes

For a new screen, first fix the responsibilities of the URL and Page state. Even when the screen includes tables or forms, decide what data the Page loads and what the state shape looks like before adding Actions and partial updates.

## 2. Request template for changing an existing Action

```text
Change an existing Action.

Target URL:
- <Example: POST /admin/projects/:projectID/members/:memberID/role>

Page name:
- <Example: ProjectMembersPage>

Action name:
- <Example: UpdateMemberRole>

Current behavior:
- <Example: redirects back to the list page after changing a role>

Desired behavior:
- <Example: after changing a role, partially update only the target row in #members-table and show a success toast>

Input:
- path parameters: <Example: projectID, memberID>
- form values / JSON / query: <Example: role>

State location:
- URL query: <Conditions to preserve. Example: q, page, sort>
- Action state / form state: <Example: MemberRoleForm>
- DB: <Persistent data to update>

Expected partial update:
- Success: <Example: replace #member-row-<memberID>>
- Validation error: <Example: show field errors inside the same row>
- Service error: <Example: show an Alert in #members-feedback>

Error display:
- Validation error: <Field error the user can fix>
- Authorization error: <Safe message for the user; details go to logs>
- Persistence error: <Retryable message>

Behavior to preserve:
- <Existing public APIs, tests, UI compatibility, and similar constraints>

Testing expectations:
- Cover the success path, validation error, authorization error, persistence error, and partial-update target.
```

### Notes

When changing an existing Action, explicitly describe the current redirect / partial-update / flash flow. If AI changes the return behavior without being told to do so, screen transitions and htmx targets can break easily.

## 3. Request template for adding forms and validation

```text
Add a form and server-side validation.

Target URL:
- GET <Example: /admin/projects/new>
- POST <Example: /admin/projects>

Page name:
- <Example: NewProjectPage>

Action name:
- <Example: CreateProject>

Form fields:
- <Example: name: required, 1-80 characters>
- <Example: slug: required, alphanumeric characters and hyphens only, must be unique>
- <Example: visibility: one of public/private>

State location:
- form state struct: <Example: ProjectForm>
- field errors: <Example: ProjectForm.Errors>
- page-level error: <Example: ProjectForm.PageError>
- DB: <Destination model / repository / service>

Validation:
- Input format: <Checks performed by form state or a validator>
- Business rules: <Checks performed by the service / domain layer>
- Uniqueness: <Checks performed by a repository / service, if needed>

Expected partial update:
- Initial display: <full page or form fragment>
- Validation error: <Example: re-render #project-form with the submitted values and errors>
- Success: <Example: redirect to the detail page, or add a row to the list and show a toast>

Error display:
- Field error: <Message shown in each FormRow>
- Page-level error: <Message shown in an Alert above the form>
- Unexpected error: <Safe generic message; details go to logs>

UI policy:
- Prefer existing FormRow, Input, Select, Textarea, Button, and Alert components.
- Preserve entered values after validation errors.

Testing expectations:
- Cover initial display, successful submit, each field error, business errors, and success return behavior.
```

### Notes

For form requests, specify field errors and page-level errors separately. Input-format errors should be shown where the user can fix them, while failures from the database or an external API usually belong in a page-level error or toast.

## 4. Request template for adding table-chart coordination

```text
Add a feature where a table and chart coordinate through the same filter state.

Purpose:
- <Example: clicking a sales trend chart narrows the orders table below to the same period and category>

Target URL:
- GET <Example: /admin/reports/sales>

Page name:
- <Example: SalesReportPage>

Action names:
- <Example: UpdateSalesReportFilters, SelectSalesChartBucket>

Shared state location:
- URL query: <Example: from, to, granularity, category, status, page, sort>
- Page state struct: <Example: SalesReportState>
- table state: <Example: SalesTableState>
- chart state: <Example: SalesChartState>
- DB / cache: <How aggregated results and table rows are stored or fetched>

Coordination rules:
- Filter change: <Example: reload both chart and table with the same conditions>
- Chart click: <Example: reflect the clicked bucket in the URL query and reset the table to page 1>
- Table page / sort change: <Example: keep chart conditions and update only the table>

Expected partial update:
- Filter change: <Example: update #sales-chart and #sales-table>
- Chart click: <Example: update #sales-table, #active-filters, and the URL query>
- Table pagination: <Example: update only #sales-table>

Error display:
- Chart load failure: <Alert or EmptyState in the chart area>
- Table load failure: <Alert or EmptyState in the table area>
- Invalid filter: <Field error in the filter form, or fall back to a safe default>

UI policy:
- Render the chart and table from shared state.
- Include filters, active-filter display, empty states, and loading states.

Testing expectations:
- Cover shared-state parsing / serialization, table conditions after chart clicks, and preserved chart conditions during pagination.
```

### Notes

For table-chart coordination, design the filter state before focusing on the rendered data. Separating conditions that belong in the URL query from temporary screen-only state makes reloads, shareable links, and partial updates easier to explain.

## 5. Information you must always provide to AI

When asking AI to implement a change, provide at least the following information.

| Information | What to write | Example |
| --- | --- | --- |
| Target URL | HTTP method, path, path parameters, and query parameters. | `GET /admin/projects/:projectID/members?q=&page=` |
| Page name | The Page function and screen state names to add or change. | `ProjectMembersPage`, `ProjectMembersPageState` |
| Action name | Action names for submit, click, filter changes, and similar events. | `InviteMember`, `UpdateMemberRole` |
| State location | Whether state belongs in URL query, form state, Page state, session, DB, or cache. | `q` and `page` are URL query; submitted values are in `InviteMemberForm` |
| Expected partial update | DOM targets and response boundaries for success, failure, and filter changes. | Success updates `#members-table` and a toast; failure updates `#invite-member-form` |
| Error display | How to split field errors, page-level Alerts, toasts, EmptyStates, 403, and 404. | Validation uses FormRow; authorization failure uses a safe Alert; details go to logs |

The following information improves accuracy further.

- Existing files, functions, and tests to use as references.
- frontend components you want to use.
- Public APIs, existing URLs, CSS classes, and tests that must not change.
- Expected displays for normal results, empty states, unauthorized users, persistence failures, and external API failures.
- Risk items that humans will review.

## 6. Items humans should verify instead of delegating to AI

AI can produce an implementation draft, but humans must verify the following items.

### Authorization

- Whether Page-level authorization checks are correct for the target tenant, organization, project, and resource.
- Whether Actions check the same operation permissions, instead of relying only on hidden UI buttons.
- Whether direct IDs for other users or other tenants are rejected.
- Whether the 403 / 404 behavior avoids leaking information that should not be visible.

### Audit logs

- Whether important operations such as creation, update, deletion, approval, permission changes, and reruns are recorded.
- Whether denied attempts, failures, validation errors, and external integration failures should be recorded in addition to successful operations.
- Whether required fields such as actor, action, target, before / after, request ID, IP, and user agent are included.
- Whether logs avoid secrets, tokens, personal information, and unnecessary request bodies.

### Destructive operations

- Whether destructive operations such as delete, cancel, overwrite, bulk update, and resync have confirmation UI.
- Whether soft delete / hard delete behavior, recoverability, cascades, and impact on related data are clear.
- Whether double submits, retries, partial failures, and interrupted execution have defined behavior.
- Whether transactions, idempotency keys, or optimistic locking are needed.

### Production settings

- Whether production environment variables, secrets, external API endpoints, callback URLs, CORS, and cookie settings are correct.
- Whether debug mode, sample authentication, development seeds, or mock APIs cannot be enabled in production.
- Whether rate limits, timeouts, retries, queues, cache TTL, monitoring, and alert settings match production operations.
- Whether migration, feature flag, rollback, and data backfill procedures are prepared.

## Pre-request checklist

- You wrote the target URL, Page name, and Action name.
- You wrote where state is stored.
- You wrote the DOM target and response boundary for the expected partial update.
- You split error display into validation, authorization, persistence failure, and unexpected failure.
- You wrote existing files, reference implementations, and compatibility constraints that must not change.
- You explicitly stated that humans will review authorization, audit logs, destructive operations, and production settings.
