## 3. Context

`Context` is passed to each handler as `func(*Context) Node` and provides request/state access.

### Request-local vs shared state

- Use `Context.Local` for temporary values that should live only for the current request.
- Use `Context.SetGlobal` / `Context.GetGlobal` for application-shared values. Any API whose name includes `Global` reads or writes state shared by all users of the app. These helpers go through the parent app mutex when the context belongs to an app.

```go
app.Page("/", func(ctx *mb.Context) mf.Node {
    ctx.Local["trace_id"] = ctx.Query("trace_id") // request-local scratch value
    ctx.SetGlobal("last_trace_id", ctx.Local["trace_id"]) // shared app state for all users
    return mf.Text(ctx.GetGlobal("last_trace_id").(string))
})
```

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

### `Local map[string]any`
- Request-local scratch map initialized for contexts created by an app.
- Values stored here are not shared with other requests and are not protected by the app mutex.

### `SetGlobal(key string, value any)`
- Writes application-shared state.
- Because this API is named `Global`, the value is shared by all users and all requests for the app.
- Contexts created by an app synchronize writes via the app mutex.
- If no app is attached, this is a no-op because there is no app-wide state map.

### `GetGlobal(key string) any`
- Reads application-shared state.
- Because this API is named `Global`, the value is shared by all users and all requests for the app.
- Contexts created by an app synchronize reads via the app mutex.
- If no app is attached, this returns `nil` because there is no app-wide state map.
- If the value is a mutable slice, map, or pointer, do not mutate it directly after `GetGlobal` returns; mutate under `UpdateGlobal` or read via `GetGlobalSnapshot` with a clone function.

### `GetGlobalSnapshot(key string, clone func(any) any) any`
- Reads application-shared state and returns `clone(value)` while the app read lock is held.
- Use this for read-only snapshots of slices, maps, or other mutable values before rendering or handing them to code that might mutate them.
- If no app is attached, it calls `clone(nil)` because there is no app-wide state map.

### `UpdateGlobal(key string, fn func(old any) any) any`
- Atomically reads, transforms, and writes application-shared state.
- Contexts created by an app run `old := state[key]`, `next := fn(old)`, and `state[key] = next` while holding the app mutex.
- Use this instead of a separate `GetGlobal` → `SetGlobal` sequence whenever the new value depends on the old value.
- If no app is attached, it calls `fn(nil)` and returns the result without storing it because there is no app-wide state map.

```go
app.Action("counter/increment", func(ctx *mb.Context) mf.Node {
    next := ctx.UpdateGlobal("count", func(old any) any {
        oldInt, _ := old.(int)
        return oldInt + 1
    }).(int)
    return mf.Text(strconv.Itoa(next))
})
```

Avoid this pattern for counters, progress ticks, or append-style updates because another request can change the value between the read and write:

```go
count := ctx.GetGlobalInt("count")
ctx.SetGlobal("count", count+1) // use UpdateGlobal or IncrementGlobalInt instead
```

### `GetGlobalInt(key string) int`
- `GetGlobal` + `int` assertion.
- Returns `0` when value is missing/not `int`.

### `IncrementGlobalInt(key string, delta int) int`
- Convenience wrapper around `UpdateGlobal` for integer counters/progress values.
- Treats missing or non-`int` values as `0`, stores `old + delta`, and returns the new `int`.

### Text stream APIs

Text stream helpers keep chunked text progress in application state so an htmx action can reveal a long response over multiple fragment requests. They are useful for AI-chat-style demos and can be paired with `App.EnableSSE()` plus a hidden element that declares `data-marionette-sse-url`.

#### `StartTextStream(options TextStreamOptions)`
- Starts or replaces a named text stream.
- `TextStreamOptions.Name` identifies the stream. Blank names are ignored.
- `TextStreamOptions.Text` is the full response text.
- `TextStreamOptions.ChunkSize` controls how many word chunks `AdvanceTextStream` reveals per call; values less than 1 use the default.

#### `AdvanceTextStream(name string) TextStreamStep`
- Reveals the next chunk for a named stream and returns the cumulative content.
- `TextStreamStep.Active` is `false` when no stream exists.
- `TextStreamStep.Done` becomes `true` on the final chunk; the stream is cleared after that final step.
- `Cursor` and `Total` report chunk progress.

#### `ResetTextStream(name string)`
- Clears the named stream.

#### `TextStreamActive(name string) bool`
- Reports whether a named stream can still advance.

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
