# Marionette

![Marionette concept art](docs/site/assets/concept.png)

Marionette is a Go-first admin UI framework concept inspired by Streamlit, with an htmx front-end runtime.

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

## Development with hot reload

Use [Air](https://github.com/air-verse/air) to restart the demo app automatically when Go files change.

```bash
go install github.com/air-verse/air@latest
air
```

This repository includes a preconfigured `.air.toml` that builds and runs `cmd/marionette`.

## Example API

## Simple full implementation sample

The runnable minimal sample is maintained outside of the README at:

- `cmd/simple-sample/main.go`

Run it with:

```bash
go run ./cmd/simple-sample
```

Then open `http://127.0.0.1:8081` to view the minimal Tasks demo.

What this sample demonstrates:

- **Entrypoint:** create an app with `mb.New()` and start it with `app.Run`
- **Routing:** `Page` for full-page rendering, `Action` for htmx partial updates
- **State management:** server-side state via `ctx.Get` / `ctx.Set`
- **UI composition:** declarative HTML building with the `frontend` Node DSL

### Import path guidance

Split imports by runtime role (`backend` for app wiring, `frontend` for UI nodes):

```go
import (
    mb "github.com/YoshihideShirai/marionette/backend"
    mf "github.com/YoshihideShirai/marionette/frontend"
)

app := mb.New()
app.Set("users", []User{})

app.Page("/", func(ctx *mb.Context) mf.Node {
    return renderUsersPage(ctx)
})

app.Action("users/create", func(ctx *mb.Context) mf.Node {
    name := ctx.FormValue("name")
    _ = name // update app state
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
