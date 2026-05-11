# Project Structure Guide

English | [日本語](ja/01-project-structure.md)

This document explains the recommended structure for **user applications** built with Marionette. It is about application layout for admin UIs and internal tools that use Marionette, not about developing Marionette itself.

## Basic policy

In a Marionette app, separate Pages that compose screens, Actions that handle user operations, State restored from requests or sessions, and app-specific UI helpers. This keeps responsibilities easy to follow as the number of screens grows.

- **Page**: Screen-level read processing, data loading, and HTML composition.
- **Action**: Change handling for form submits, button clicks, filter changes, and similar user operations.
- **State**: Query parameters, form values, session-derived values, and conditions shared across screens.
- **UI helper**: Layouts, shared panels, business-specific small display parts, and other composition helpers used only inside the app.

Keep the file count low for small apps. Split by responsibility once screens and operations start to grow.

## 1. Minimal structure for small apps

For a small internal tool with only a few screens, it is fine to start with a single `internal/app` package.

```text
my-admin-app/
  cmd/server/main.go
  internal/app/
    app.go
    routes.go
    pages.go
    actions.go
    state.go
    ui.go
```

The files have these responsibilities.

| File | Responsibility |
| --- | --- |
| `cmd/server/main.go` | Executable entry point. Loads configuration, creates the Marionette app, and starts the HTTP server. |
| `internal/app/app.go` | App initialization, middleware, and dependency assembly. |
| `internal/app/routes.go` | Connects URL paths to Page and Action handlers. |
| `internal/app/pages.go` | Screen-level render functions such as dashboards and settings pages. |
| `internal/app/actions.go` | POST or htmx action handlers for form submits and state-changing operations. |
| `internal/app/state.go` | Types and parsers for values restored from query parameters, filters, form input, and sessions. |
| `internal/app/ui.go` | Small app-specific UI helpers such as app shells, sections, and cards. |

This structure works well when:

- There are only a few screens and only a few Page / Action handlers.
- The app handles one domain or one workflow.
- Initial implementation speed and readability matter more than preparing for future package splits.

When `pages.go` or `actions.go` becomes large enough that handlers are hard to find, move to the split structure for medium and larger apps.

## 2. Split structure for medium and larger admin apps

When the number of screens grows and the app has multiple feature areas such as dashboards, user management, job management, and settings, split packages by responsibility.

```text
my-admin-app/
  cmd/server/main.go
  internal/app/
    app.go
    routes.go
  internal/pages/
    dashboard.go
    users.go
    settings.go
  internal/actions/
    user_actions.go
    job_actions.go
  internal/state/
    query_state.go
    session_state.go
  internal/ui/
    layout.go
    components.go
  internal/assets/
    assets.go
```

The directories have these responsibilities.

| Directory / file | Responsibility |
| --- | --- |
| `cmd/server/main.go` | Binary entry point. Handles production / development configuration, server startup, graceful shutdown, and similar concerns. |
| `internal/app/app.go` | Creates the Marionette app, injects dependencies, installs common middleware, and performs app-wide initialization. |
| `internal/app/routes.go` | Central route registration. Connects `pages` and `actions` to URLs. |
| `internal/pages/` | Handlers and render functions for displaying screens. Split files by screen or feature area. |
| `internal/actions/` | Handles operations that change state. Owns validation, persistence, redirects, and partial-update responses. |
| `internal/state/` | Types and restoration logic for URL queries, sessions, form input, filters, pagination, and similar values. |
| `internal/ui/` | App-specific layouts and small composition helpers. These are not generic components added to Marionette itself. |
| `internal/assets/` | Registration and serving helpers for app-specific embedded CSS, images, JavaScript, favicons, and other assets. |

### When to split

Move from the minimal structure to the split structure when you see signs like these:

- Multiple independent screens are mixed into one `pages.go` file.
- Action handlers have grown enough that it is hard to tell which Page calls which Action.
- Query state, form state, and session state types are increasing.
- Shared layout or card helpers are used by multiple screens.
- Route registration has become long and hard to read as a screen list.

## 3. Where to put Page / Action / State / UI helper code

### Put Pages in `internal/pages/`

A Page usually restores display state from the HTTP request, loads the data it needs, and composes HTML with Marionette `frontend` components.

Example:

```text
internal/pages/
  dashboard.go  # dashboard page
  users.go      # user list / user detail page
  settings.go   # settings page
```

Put these in Pages:

- GET request handlers.
- Screen-level render functions.
- Whole-screen composition for page titles, breadcrumbs, tables, forms, and similar structures.
- Small helpers that are specific to one Page and are not reused by other screens.

Avoid putting these in Pages:

- Change processing such as database updates, writes to external APIs, or email sending.
- Query or session parsers shared by multiple screens.
- Layout helpers used across the whole app.

### Put Actions in `internal/actions/`

Actions handle operations that change state in response to user input. POST requests, htmx requests, button clicks, form submits, and filter changes are typical examples.

Example:

```text
internal/actions/
  user_actions.go  # create, update, delete, role change
  job_actions.go   # enqueue, retry, cancel
```

Put these in Actions:

- Form-submit validation.
- Updates to databases or durable state.
- Session updates.
- Redirects and flash message setup.
- Response composition for htmx partial updates.

If Actions call large Page render functions directly too often, dependencies become hard to follow. Move the small display helpers needed for partial updates into `internal/ui/`, or extract small functions near the target Page.

### Put State in `internal/state/`

State is an explicit type for values restored from requests and for inputs shared by Pages and Actions.

Example:

```text
internal/state/
  query_state.go    # search, sort, pagination, filter
  session_state.go  # login user, selected workspace, flash-like values
```

Put these in State:

- Filters, sort options, and pagination restored from URL queries.
- Structs that receive form values.
- Form state that includes validation errors.
- User context or workspace context restored from sessions.
- Request parsing logic shared by Pages and Actions.

State is not only about where a value is stored; it also describes the value's lifetime. Keep shareable URL state, request-local form state, server-side session state, and business data stored in durable storage separate from each other.

### Put UI helpers in `internal/ui/`

UI helpers combine Marionette `frontend` components into app-specific appearance and repeated structures.

Example:

```text
internal/ui/
  layout.go      # app shell, navigation, page frame
  components.go  # status badge, metric card, empty state helper
```

Put these in UI helpers:

- Layouts such as app shells, sidebars, navbars, and footers.
- Business-specific badges, cards, empty states, and toolbars.
- Sections, panels, and form wrappers used by multiple Pages.
- App-only composition built by combining Marionette components.

Do not put these in UI helpers:

- Generic components that should be added to Marionette itself.
- Page-specific business logic.
- Action validation or persistence logic.
- State restoration from sessions or databases.

## 4. Do not confuse this with Marionette's own structure

This repository contains framework implementation directories such as `backend/`, `frontend/`, and `templates/components/`. These are **framework-side implementation locations**. A normal user app does not need to create the same directories.

| Marionette directory | Role | User-app counterpart |
| --- | --- | --- |
| `backend/` | Implementation of Marionette backend APIs such as App, Context, routing, and runtime. | Put app initialization and route registration in `internal/app/`. |
| `frontend/` | Implementation of the Go UI component APIs provided by Marionette. | Use the `frontend` package and put app-specific combinations in `internal/ui/`. |
| `templates/components/` | Framework implementation related to Marionette component templates or generation sources. | Put user-app screen helpers in `internal/pages/` or `internal/ui/`. |

Important notes:

- You do not need to create `backend/` or `frontend/` directories in a user app.
- Do not add app-specific layouts or helpers to Marionette's `frontend/` directory.
- Treat adding a generic component to Marionette itself as a different task from composing screens in a user app.
- `templates/components/` is related to Marionette's own component implementation. Do not use it as the location for user-app Pages or partial templates.
- In user apps, use existing `frontend` components first, then add only the missing app-specific combinations as helpers in `internal/ui/`.

## Recommended structure summary

Start small, then split as the number of screens and responsibilities grows.

```text
my-admin-app/
  cmd/server/main.go
  internal/app/
    app.go
    routes.go
  internal/pages/
    dashboard.go
    users.go
    settings.go
  internal/actions/
    user_actions.go
    job_actions.go
  internal/state/
    query_state.go
    session_state.go
  internal/ui/
    layout.go
    components.go
  internal/assets/
    assets.go
```

With this structure, `internal/app` handles app assembly and routing, `internal/pages` handles reads and screen composition, `internal/actions` handles changes, `internal/state` handles request / session / form state, and `internal/ui` handles app-specific display helpers. Marionette's own `backend/`, `frontend/`, and `templates/components/` directories are framework implementation directories with different responsibilities from user-app structure.
