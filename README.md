# Marionette

Marionette is a Go-first admin UI framework concept inspired by Electron/Tauri, with an htmx front-end runtime.

## Goals

- Framework users build the UI with Go only.
- Admin screens can be composed from pages, actions, forms, and partial updates.
- htmx handles browser-side interaction and partial swaps.
- The runtime can be hosted in a WebView for desktop apps.

## Architecture

- **Backend (Go):** routing, state updates, event handlers.
- **Frontend (htmx):** transport and incremental HTML swaps.
- **Styling:** Tailwind CSS + daisyUI (CDN).
- **Declarative Go UI DSL:** `Text`, `Div`, `Column`, `Form`, `Input`, `Submit`, `Button(...).Post(...)` (rendered via `html/template`).

## Run example

```bash
go run ./cmd/marionette
```

Open http://127.0.0.1:8080 and try the users admin demo.

## Example API

```go
app := marionette.New()
app.Set("users", []User{})

app.Page("/", func(ctx *marionette.Context) marionette.Node {
    return renderUsersPage(ctx)
})

app.Action("users/create", func(ctx *marionette.Context) marionette.Node {
    name := ctx.FormValue("name")
    // update app state
    return renderUsersWorkspace(ctx)
})
```

`Page` handlers return full pages wrapped in the Marionette shell. `Action`
handlers are POST-only htmx endpoints and return HTML fragments for partial
swaps. The older `Render` and `Handle` APIs still work for small examples.

## API Documentation

Detailed API reference is available here:

- [`docs_api.md`](./docs_api.md)


## GitHub Pages deployment

- The Pages deploy workflow publishes `docs/site/` via GitHub Actions.
- Before enabling deploys, set repository **Settings → Pages → Build and deployment → Source** to **GitHub Actions**.
