## 3. Context

`Context` is passed to each handler as `func(*Context) Node` and provides request/state access.

### `Param(name string) string`
- Returns path parameter from `Request.PathValue(name)`.
- Returns `""` when `Request` is `nil`.

### `Query(name string) string`
- Returns query parameter from `Request.URL.Query().Get(name)`.
- Returns `""` when `Request` is `nil`.

### `FormValue(name string) string`
- Returns form value from `Request.FormValue(name)`.
- Returns `""` when `Request` is `nil`.

### `Asset(name string) string`
- Builds an asset URL from the parent app.
- Use inside handlers when images, stylesheets, or links should follow the
  registered app asset prefix.

### `Set(key string, value any)`
- Writes shared state.
- If context has parent app, write is synchronized via app mutex.
- If no app is attached, writes directly to `Context.State`.

### `Get(key string) any`
- Reads shared state.
- If context has parent app, read is synchronized via app mutex.
- If no app is attached, reads directly from `Context.State`.

### `GetInt(key string) int`
- `Get` + `int` assertion.
- Returns `0` when value is missing/not `int`.

### Flash APIs

#### `Flashes() []FlashMessage`
- Returns a copied snapshot of currently loaded flashes.
- Returns `nil` when no flash exists.

#### `FlashSuccess(message string)` / `FlashError(message string)` / `FlashInfo(message string)` / `FlashWarn(message string)`
- Convenience wrappers around `AddFlash(level, message)`.
- Level values are implementation constants:
  - `FlashSuccess`, `FlashError`, `FlashInfo`, `FlashWarn`.

#### `AddFlash(level FlashLevel, message string)`
- Trims message; empty after trim means no-op.
- Appends flash into context flash list, serializes to cookie (`marionette_flash`).
- Cookie behavior:
  - `Path=/`
  - `HttpOnly=true`
  - `SameSite=Lax`
  - `Secure` follows `App.SetCookieSecure` (default `false`).
- Serialization failure is ignored (no panic / no status change).

Flash lifecycle on next request:
- Flashes are decoded from cookie into `Context.flashes`.
- Valid entries only: known levels and non-empty messages.
- If flashes were present, cookie is automatically cleared in response.

### Session APIs

#### `SetSession(key, value string)`
- Trims `key`; empty key means no-op.
- Stores/updates a session entry in context memory and writes cookie (`marionette_session`).
- Cookie behavior:
  - `Path=/`
  - `HttpOnly=true`
  - `SameSite=Lax`
  - `Secure` follows `App.SetCookieSecure` (default `false`).

#### `Session(key string) string`
- Reads the session value by key.
- Returns empty string when key is missing.

#### `ClearSession()`
- Replaces the session map with an empty map and writes it to cookie (`marionette_session`).

Session lifecycle on request:
- Session is decoded from cookie into `Context.session` in `newContext`.
- Decode failure falls back to empty session map (no panic / no status change).

Session sample:

```go
app.Page("/session", func(ctx *mb.Context) mf.Node {
    user := ctx.Session("user")
    if user == "" {
        return mf.Form("/session/login", mf.Button("Sign in"))
    }
    return mf.Form("/session/logout", mf.Button("Sign out"))
})
app.Action("session/login", func(ctx *mb.Context) mf.Node {
    ctx.SetSession("user", "Aiko")
    return mf.Paragraph("Signed in")
})
app.Action("session/logout", func(ctx *mb.Context) mf.Node {
    ctx.ClearSession()
    return mf.Paragraph("Signed out")
})
```

Full example: `docs/site-astro/public/examples/go/session.go`.

---
