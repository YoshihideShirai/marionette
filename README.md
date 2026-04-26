# Marionette

Marionette is a Go-first GUI framework concept inspired by Electron/Tauri, with an htmx front-end runtime.

## Goals

- Framework users build the UI with Go only.
- htmx handles browser-side interaction and partial swaps.
- The runtime can be hosted in a WebView for desktop apps.

## Architecture

- **Backend (Go):** routing, state updates, event handlers.
- **Frontend (htmx):** transport and incremental HTML swaps.
- **Declarative Go UI DSL:** `Text`, `Div`, `Column`, `Button(...).OnClick(...)`.

## Run example

```bash
go run ./cmd/marionette
```

Open http://127.0.0.1:8080 and click **Increment**.

## Example API

```go
app := marionette.New()
app.Set("count", 0)

app.Render(func(ctx *marionette.Context) marionette.Node {
    return renderCounter(app.GetInt("count"))
})

app.Handle("counter/increment", func(ctx *marionette.Context) marionette.Node {
    count := app.GetInt("count") + 1
    app.Set("count", count)
    return renderCounter(count)
})
```
