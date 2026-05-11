# Routing / Pages / Actions Guide

English | [日本語](ja/02-routing-pages-actions.md)

This document explains how to split URLs, Pages, and Actions in Marionette applications. It is intended as a design note for admin UIs and internal tools where routes and handler responsibilities need to stay easy to follow as the number of screens grows.

For details about where state should live, including URL query values, form values, sessions, and global state, see the [State Management Guide](../state-management.md). This document focuses on the boundary between Pages and Actions rather than covering state management in depth.

## Basic policy

In a Marionette app, treat **Pages as display entry points** and **Actions as user-operation entry points**.

- **Page**: Restores display state from the request, loads the required data, and returns HTML for a full screen or partial screen.
- **Action**: Receives input from form submissions, button clicks, filter changes, and similar operations, then calls validation and business processing and decides how to display the result.
- **State**: Represents query values, form values, validation errors, and session-derived values shared by Pages and Actions.
- **Service / domain**: Owns database updates, external API calls, business rules, transactions, and similar concerns.
- **UI helper / component**: Owns HTML composition, common layouts, panels, tables, form sections, and similar display pieces.

The thinner the handlers are, the easier it is to understand which URL leads to which screen or operation from the route list.

## 1. How to split Pages

Split Pages around the display unit seen by the user. In practice, combine three criteria: URL unit, business-screen unit, and partial-update unit.

### URL unit

Start by creating Pages per URL route. When the URL changes, browser history, bookmarks, reload behavior, permissions, titles, and breadcrumbs may also change.

Example:

```text
GET /admin/users          -> pages.UsersIndex
GET /admin/users/{id}     -> pages.UserDetail
GET /admin/settings       -> pages.Settings
```

Splitting Pages by URL makes the screen structure easy to read from route definitions. However, small cards or table bodies inside the same URL do not need to be treated as Pages automatically.

### Business-screen unit

Even under the same URL area, split Page functions or render functions when the business meaning of the display is substantially different.

Examples:

- User list and user detail.
- Invoice list and invoice creation.
- Job execution history and job settings.
- Administrator settings and regular-user settings.

Splitting by business screen keeps the required state, permissions, loaded data, and display components naturally grouped. If one Page starts to contain multiple business contexts, it is a sign to split the file or function.

### Partial-update unit

When htmx or similar behavior updates only part of the screen, it is often useful to extract a render function for each **partial-update target**.

Examples:

- Update only the table body when search conditions change.
- Update only the detail panel after saving.
- Re-render filters and charts from the same state.
- Lazy-load only the contents of a modal.

Treat partial-update functions as helpers that redraw part of a Page, not necessarily as independent Pages. Keep URL handling, permissions, and primary state restoration in the parent Page, and let the partial receive the required state and data before returning HTML.

## 2. Action responsibilities

An Action is a thin orchestration layer that receives a user operation and returns the processing result to the screen. A typical Action has four steps.

1. **Acquire input**: Read `ctx.FormValue`, path parameters, query parameters, session-derived values, and similar input.
2. **Call validation**: Pass input values to form state or validators and receive field errors or page-level errors.
3. **Call domain processing**: Call service / domain methods to save, delete, integrate with external systems, or perform other business operations.
4. **Display the result**: Choose a redirect, flash message, toast, partial HTML response, form re-render, or other user-facing result.

The Action itself should decide what to call and how to return the result. If database update steps, external API request/response details, or complex HTML construction are written directly inside Actions, they become harder to test and the route structure becomes harder to understand.

## 3. Criteria for keeping Actions small

When an Action grows, move responsibilities out using the following criteria.

### Move database updates and external API calls to the service layer

Do not write the following directly in Actions. Move them to the service / domain layer instead.

- Database transactions, inserts / updates / deletes, and consistency across multiple tables.
- External API calls, retries, timeouts, and response mapping.
- Business-rule decisions, including authorization-related rules.
- Side effects of screen operations, such as audit logs, notifications, and queue enqueueing.

Actions should pass the required input to services, receive success or failure, and convert that result into a display response. This makes business processing testable without HTTP handlers.

### Move display HTML composition to UI helpers / components

Move the following to UI helpers / components.

- Building tables, cards, form sections, empty states, and alerts.
- Re-rendering forms with validation errors.
- Returning partial HTML for htmx targets.
- Common layouts and panels used by multiple Pages or Actions.

When many HTML nodes are listed directly inside an Action, input handling and display structure become mixed together. Keep Actions easy to read by calling helpers such as `renderUserForm(state)` or `ui.UserTable(rows, filters)` instead.

## 4. Good example / example to avoid

### Good example: the Action focuses on input, validation, service calls, and result display

```go
func (h *UserActions) Create(ctx *backend.Context) error {
    form := state.UserFormFromRequest(ctx)
    if result := form.Validate(); !result.OK() {
        form.Errors = result.Errors
        return ctx.HTML(ui.UserForm(form))
    }

    user, err := h.Users.Create(ctx.Request().Context(), service.CreateUserInput{
        Name:  form.Name,
        Email: form.Email,
    })
    if err != nil {
        form.PageError = "Could not create the user"
        return ctx.HTML(ui.UserForm(form))
    }

    ctx.Flash("success", "User created")
    return ctx.Redirect("/admin/users/" + user.ID)
}
```

In this example, the Action only acquires the form, validates it, calls the service, and displays the result. Persistence is delegated to `h.Users.Create`, and HTML composition is delegated to `ui.UserForm`.

### Example to avoid: database updates and HTML construction are packed into the Action

```go
func CreateUser(ctx *backend.Context, db *sql.DB) error {
    name := ctx.FormValue("name")
    email := ctx.FormValue("email")
    if name == "" || !strings.Contains(email, "@") {
        return ctx.HTML(html.Div(
            html.P("Check the input values"),
            html.Input(html.Name("name"), html.Value(name)),
            html.Input(html.Name("email"), html.Value(email)),
        ))
    }

    tx, err := db.BeginTx(ctx.Request().Context(), nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    row := tx.QueryRowContext(ctx.Request().Context(),
        "insert into users (name, email) values ($1, $2) returning id",
        name, email,
    )

    var id string
    if err := row.Scan(&id); err != nil {
        return err
    }
    if err := tx.Commit(); err != nil {
        return err
    }

    return ctx.HTML(html.Div(
        html.H2("Created"),
        html.A(html.Href("/admin/users/"+id), html.Text("View details")),
    ))
}
```

This example mixes validation, database transactions, SQL, and HTML composition in one Action. Both screen changes and business-processing changes require touching the same function, and the test target becomes larger.

## Implementation checklist

- Displays that change URLs are split as Pages.
- State, data loading, and render functions are easy to follow for each business screen.
- Partial updates are extracted into helpers that can redraw from the parent Page state.
- Actions focus on input acquisition, validation calls, service calls, and result display.
- Database updates, external API calls, and business rules live in the service / domain layer.
- HTML composition lives in UI helpers / components.
- Questions about state lifetime and placement are checked against the [State Management Guide](../state-management.md).
