## 2. App

### `New() *App`
- Returns a new `*App`.
- Initializes internal state:
  - `state: map[string]any{}`
  - `pages: map[string]Handler{}`
  - `actions: map[string]Handler{}`
  - `cookieSecure: false` (default).

### `Run(addr string) error`
- Starts an HTTP server with `http.ListenAndServe(addr, a.Handler())`.
- Logs `marionette listening at http://<addr>` to stdout before serving.
- Returns any `ListenAndServe` error as-is.

### Desktop runtime: `desktop.Run(app *backend.App, options desktop.Options) error`
- Starts the app on a private `127.0.0.1:0` server using `app.Handler()`.
- Opens that local URL in a native WebView shell.
- Shuts the local server down when the WebView exits.
- Requires building with `-tags marionette_desktop` to enable the native
  WebView adapter.

### Desktop options
- `Title string`: native window title; defaults to `"Marionette"`.
- `Width int`: native window width; defaults to `1200`.
- `Height int`: native window height; defaults to `800`.
- `Debug bool`: passes WebView debug mode through to the adapter.

### Page options
- `type PageOptions struct { Title string }`
- `type PageOption func(*PageOptions)`
- `WithTitle(title string) PageOption`
  - Trims whitespace and sets the HTML `<title>` for that page route.

### `SetCookieSecure(secure bool)`
- Enables/disables `Secure` on the flash cookie (`marionette_flash`).
- Default is `false`.
- Affects both:
  - flash write (`Context.AddFlash`)
  - flash clear (when flashes are consumed on next request).

### `Assets(prefix string, fsys fs.FS, options ...AssetOption)`
- Serves static files from `fsys` under `prefix`.
- Use `os.DirFS("assets")` for local files or `embed.FS` plus `fs.Sub` for
  single-binary assets.
- Directory index responses are disabled by default; pass
  `WithAssetIndex(true)` only when directory browsing is intentional.
- Asset routes are included in `Handler()` and `Run()`.

### `Downloads(prefix string, fsys fs.FS, options ...AssetOption)`
- Serves static files from `fsys` under `prefix` with
  `Content-Disposition: attachment`.
- Equivalent to `Assets(prefix, fsys, WithAssetDownload(), options...)`.
- Uses the requested file basename as the attachment filename.

### `Asset(name string) string`
- Builds an asset URL from the first registered asset prefix.
- Example: after `app.Assets("/assets", ...)`, `app.Asset("hero.jpg")`
  returns `"/assets/hero.jpg"`.
- Absolute `http://`, `https://`, and `data:` URLs pass through unchanged.

### Asset options
- `WithAssetCache(maxAge time.Duration)` emits
  `Cache-Control: public, max-age=<seconds>`.
- `WithAssetImmutable()` adds `immutable` when asset caching is enabled.
- `WithAssetIndex(enabled bool)` allows directory index responses from the
  underlying file server.
- `WithAssetDownload()` serves matching files as attachment downloads.
- `WithAssetContentTypes(types map[string]string)` sets `Content-Type` by
  extension before serving.

### `AddStylesheet(href string)`
- Adds a custom stylesheet link to the full-page HTML shell.
- Empty/whitespace-only values are ignored.
- Stylesheets are emitted after the built-in Tailwind/daisyUI assets.

### `AddStyle(css string)`
- Adds trusted inline CSS to the full-page HTML shell.
- Empty/whitespace-only values are ignored.
- Use for small app-level overrides or CSS variables.

### `AddScript(src string)`
- Adds an external JavaScript file to the full-page HTML shell.
- Empty/whitespace-only values are ignored.
- Scripts are emitted after Marionette's built-in JavaScript.

### `AddJavaScript(js string)`
- Adds trusted inline JavaScript to the full-page HTML shell.
- Empty/whitespace-only values are ignored.
- Inline JavaScript is emitted after custom external scripts, so it can use
  libraries registered with `AddScript`.

### `Handler() http.Handler`
- Builds and returns `*http.ServeMux` with all registered routes.
- `Page` routes:
  - path match is strict (`r.URL.Path` must equal registered path, otherwise `404`).
  - method is `GET` only, otherwise `405 Method Not Allowed`.
  - renders full HTML shell.
- `Action` routes:
  - method is `POST` only, otherwise `405 Method Not Allowed`.
  - executes `r.ParseForm()`, parse failure is `400 Bad Request`.
  - renders HTML fragment only.
- If `/` is not registered by `Page`/`Render`:
  - `GET /` returns `500 Internal Server Error` with configuration message.
  - non-root unmatched paths are `404`.

### `Page(path string, fn Handler, options ...PageOption)`
- Registers full-page handler for `GET`.
- Path normalization:
  - `""` -> `"/"`
  - no leading slash -> leading slash is added.
- Render mode: handler `Node` is wrapped in Marionette shell HTML.

### `Action(name string, fn Handler)`
- Registers fragment handler for `POST`.
- Name normalization:
  - always stored as `"/" + strings.TrimPrefix(name, "/")`.
- Parse failure in request form body -> `400`.
- Render mode: handler `Node` is returned as fragment HTML.

### `Render(fn Handler, options ...PageOption)`
- Compatibility alias for root page registration.
- Equivalent to `Page("/", fn)`.

### `Handle(name string, fn Handler)`
- Compatibility alias for action registration.
- Equivalent to `Action(name, fn)`.

### App state helpers

#### `Set(key string, value any)`
- Writes into app shared state map with lock.

#### `Get(key string) any`
- Reads from app shared state map with lock.

#### `GetInt(key string) int`
- Reads app shared state and type-asserts to `int`.
- Returns `0` when value is missing or not `int`.

---
