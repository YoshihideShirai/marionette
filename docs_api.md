# API Documentation

This document provides an organized reference for the public API in the `marionette` package.

## 1. Application setup

```go
app := marionette.New()
```

`New` initializes:

- a shared state store (`map[string]any`)
- page (GET) handlers
- action (POST) handlers

## 2. Handlers and context

### 2.1 Handler type

```go
type Handler func(*Context) Node
```

- Every handler returns a `Node`.
- In `Page`, the returned `Node` is rendered as a full HTML page (wrapped in the Marionette shell).
- In `Action`, the returned `Node` is rendered as an HTML fragment.

### 2.2 Context

`Context` provides access to request data and shared state.

```go
ctx.Param("id")       // path parameter
ctx.FormValue("name") // form value
ctx.Query("q")        // query parameter

ctx.Set("count", 3)
count := ctx.GetInt("count")
```

> `GetInt` returns `0` when the underlying value is not an `int`.

## 3. Routing APIs

### 3.1 Page

```go
app.Page("/users", func(ctx *marionette.Context) marionette.Node {
    return marionette.Div("app", marionette.Text("Users"))
})
```

- Accepts GET only.
- If `path` does not start with `/`, it is normalized automatically (for example, `users` becomes `/users`).

### 3.2 Action

```go
app.Action("users/create", func(ctx *marionette.Context) marionette.Node {
    name := ctx.FormValue("name")
    _ = name
    return marionette.Div("app", marionette.Text("created"))
})
```

- Accepts POST only.
- `name` is normalized internally whether or not it includes a leading `/`.
- If `r.ParseForm()` fails, Marionette returns HTTP 400.

### 3.3 Compatibility APIs

```go
app.Render(fn)       // same as app.Page("/", fn)
app.Handle(name, fn) // same as app.Action(name, fn)
```

These remain available for compatibility with earlier examples.

### 3.4 Exporting the HTTP handler

```go
mux := app.Handler()
http.ListenAndServe(":8080", mux)
```

- If `/` is not registered, `GET /` returns HTTP 500 with a clear configuration error.

## 4. State management API

```go
app.Set("users", users)

ctx.Set("users", nextUsers) // handlers update the same shared store
```

- `App` and `Context` reference the same backing map.
- If you need concurrent access safety, add synchronization at the application level.

## 5. UI node APIs

All UI components implement `Node` and render HTML via `Render() (template.HTML, error)`.

### 5.1 Basic nodes

- `Text(v string) Node`
- `Div(id string, children ...Node) Node`
- `DivClass(id, className string, children ...Node) Node`
- `Column(children ...Node) Node`
- `Raw(html string)` for trusted HTML snippets

### 5.2 Tables

```go
rows := []marionette.TableRowData{
    marionette.TableRow(marionette.Text("alice"), marionette.Text("admin")),
}

table := marionette.Table([]string{"name", "role"}, rows...)
```

- `Table(headers []string, rows ...TableRowData) Node`
- `TableRow(cells ...Node) TableRowData`

### 5.3 Sidebar

```go
sidebar := marionette.
    Sidebar("Marionette", "Admin",
        marionette.SidebarLink("Users", "/").Active(),
        marionette.SidebarLink("Reports", "/reports"),
    ).
    Note("Tip", "Use htmx actions for partial updates")
```

- `Sidebar(brand, title string, items ...SidebarItem) *sidebar`
- `SidebarLink(label, href string) SidebarItem`
- `SidebarItem.Active() SidebarItem`
- `(*sidebar).Note(title, text string) *sidebar`

### 5.4 Forms

```go
form := marionette.
    Form("users/create",
        marionette.Input("name", ""),
        marionette.Submit("Create"),
    ).
    Target("#workspace")
```

- `Form(action string, children ...Node) *form`
- `(*form).Target(selector string) *form`
- `Input(name, value string) Node`
- `HiddenInput(name, value string) Node`
- `Submit(label string) Node`

`Form` automatically sets `hx-post`, `hx-target`, and `hx-swap="outerHTML"`.

### 5.5 Buttons

```go
btn := marionette.
    Button("Refresh").
    Post("users/refresh").
    Target("#workspace")
```

- `Button(label string) *button`
- `(*button).Post(action string) *button`
- `(*button).OnClick(action string) *button` (`Post` alias)
- `(*button).Target(selector string) *button`
- `(*button).TargetSelector(selector string) *button`

## 6. Minimal example

```go
package main

import (
    "fmt"
    "net/http"

    "marionette/internal/marionette"
)

func main() {
    app := marionette.New()
    app.Set("count", 0)

    app.Page("/", func(ctx *marionette.Context) marionette.Node {
        count := ctx.GetInt("count")
        return marionette.Div("app",
            marionette.Text("count"),
            marionette.Text(": "),
            marionette.Text(fmt.Sprintf("%d", count)),
            marionette.Button("+1").Post("counter/inc"),
        )
    })

    app.Action("counter/inc", func(ctx *marionette.Context) marionette.Node {
        next := ctx.GetInt("count") + 1
        ctx.Set("count", next)
        return marionette.Div("app",
            marionette.Text(fmt.Sprintf("count: %d", next)),
            marionette.Button("+1").Post("counter/inc"),
        )
    })

    _ = http.ListenAndServe(":8080", app.Handler())
}
```

## 7. HTTP behavior summary

- `Page`: GET only (`405 Method Not Allowed` otherwise)
- `Action`: POST only (`405 Method Not Allowed` otherwise)
- Action form parsing failure: `400 Bad Request`
- Node/template rendering failure: `500 Internal Server Error`
- `GET /` without page registration: `500 Internal Server Error`
