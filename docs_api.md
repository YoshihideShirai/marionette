# API Documentation

This document is a direct API reference for the current runtime surface in the
module root package, reorganized by runtime layer.

## Import Path

Split imports by runtime role:

```go
import (
    mb "github.com/YoshihideShirai/marionette/backend"
    md "github.com/YoshihideShirai/marionette/desktop"
    mf "github.com/YoshihideShirai/marionette/frontend"
    mh "github.com/YoshihideShirai/marionette/frontend/html"
)
```

- Recommended aliases: `mb` (marionette backend), `mf` (marionette frontend).
- Use `mb` for app/runtime APIs such as `New`, `App`, `Context`, `Handler`.
- Use `md` for desktop runtime APIs such as `desktop.Run`.
- Use `mf` for component APIs such as `Button`, `Card`, `Table`, `FormRow`.
- Use `mh` for advanced low-level node APIs such as `Node`, `Div`, `Element`, `Raw`.
- Low-level constructors are intentionally not exposed from `frontend`; import
  `frontend/html` when custom markup is needed.

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

### `Page(path string, fn Handler)`
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

### `Render(fn Handler)`
- Compatibility alias for root page registration.
- Equivalent to `Page("/", fn)`.

### `Handle(name string, fn Handler)`
- Compatibility alias for action registration.
- Equivalent to `Action(name, fn)`.

### App state helpers

#### `Set(key string, value any)`
- Writes into app shared state map with lock.

#### `GetInt(key string) int`
- Reads app shared state and type-asserts to `int`.
- Returns `0` when value is missing or not `int`.

---

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

---

## 4. Low-level HTML (`frontend/html`)

Low-level HTML constructors live in `github.com/YoshihideShirai/marionette/frontend/html`.
They are intended for advanced users and component internals. The `frontend`
package exposes component APIs; use the `mh` import shown above for custom markup.

### `type Node interface { Render() (template.HTML, error) }`
- Every UI node renders itself to safe HTML.
- Rendering failure eventually becomes `500 Internal Server Error` in HTTP responses.

### Basic node constructors
- `Text(v string) Node`
- `Raw(html string)` (`type Raw string`, trusted HTML passthrough)
- `type Attrs map[string]string`
- `type ElementProps struct { ID string; Class string; Attrs Attrs }`
- `Element(tag string, props ElementProps, children ...Node) Node`
- `Div(children ...Node) Node`
- `DivID(id string, children ...Node) Node`
- `DivClass(className string, children ...Node) Node`
- `DivAttrs(attrs Attrs, children ...Node) Node`
- `DivProps(props ElementProps, children ...Node) Node`
- `Column(children ...Node) Node`

### Table / layout helpers
- `Table(headers []string, rows ...TableRowData) Node`
- `TableRow(cells ...Node) TableRowData`
- `Sidebar(brand, title string, items ...SidebarItem) *sidebar`
- `SidebarLink(label, href string) SidebarItem`
- `(SidebarItem).Active() SidebarItem`
- `(*sidebar).Note(title, text string) *sidebar`

### Legacy form/button nodes
- `Form(action string, children ...Node) *form`
  - default target selector: `#app`.
  - rendered attrs include `hx-post`, `hx-target`, `hx-swap="outerHTML"`.
- `(*form).Target(selector string) *form`
- `Input(name, value string) Node`
- `HiddenInput(name, value string) Node`
- `Submit(label string) Node`
- `Button(label string) *button`
  - default target selector: `#app`.
- `(*button).Post(action string) *button` (action normalized without leading `/`)
- `(*button).OnClick(action string) *button` (`Post` alias)
- `(*button).Target(selector string) *button`
- `(*button).TargetSelector(selector string) *button`

---

## 5. Form APIs

### `FormRow(props FormRowProps) Node`
- Required:
  - `props.ID` must be non-empty; otherwise render error node.
  - `props.Control` must be non-nil; otherwise render error node.
- Behavior:
  - optional label, required marker (`*`), description, error row.
  - when `props.Error` exists, internally renders `FieldError`.

### `FieldError(props FieldErrorProps) Node`
- Empty `Message` => returns empty `Raw("")`.
- Non-empty message requires non-empty `ID`; empty ID => render error node.

### `TextField(props TextFieldProps) Node`
- Defaults:
  - `Type` defaults to `"text"` when blank.
- Behavior:
  - applies `aria-describedby` from description/error presence.
  - sets `aria-invalid=true` when error exists.
  - supports `required`, `disabled`, `readonly`, `data-ref`.

### `Textarea(props TextareaProps) Node`
- `Rows > 0` adds `rows` attribute; otherwise omitted.
- Same accessibility/error behavior as `TextField`.

### `Select(props SelectFieldProps) Node`
- Renders `<option value="...">` list from `Options`.
- `SelectOption.Selected` sets `selected="selected"`.
- Same accessibility/error behavior as text controls.

### `Checkbox(props CheckboxProps) Node`
- Renders label + checkbox input.
- `Checked=true` sets `checked="checked"`.
- Supports `disabled`, `readonly`, `data-ref`, error aria attrs.

### `RadioGroup(props RadioGroupProps) Node`
- Required: non-empty `props.ID`; empty ID => render error node.
- Generates child IDs as `<groupID>-<index>`.
- Marks checked option where `props.Value == option.Value`.

### `Switch(props SwitchProps) Node`
- Implemented as checkbox input with switch class.
- `Checked=true` sets `checked="checked"`.
- Supports same checkable attrs as checkbox/radio.

---

## 6. Component APIs

Template-backed component constructors (`templates/components/*`).

### Buttons / inputs / field wrappers
- `Button(label string, props ComponentProps) Node`
- `SubmitButton(label string, props ComponentProps) Node`
- `Link(props LinkProps) Node`
  - renders an anchor for internal, external, and download links.
  - `Icon` renders an `aria-hidden` icon before the label.
  - `External` defaults `target="_blank"` and `rel="noopener noreferrer"`.
  - `Target: "_blank"` also defaults `rel="noopener noreferrer"`.
  - `Download` emits `download`; `Filename` emits `download="<filename>"`.
  - `Props.Disabled` renders an inert `href="#"` link with `aria-disabled`.
- `ExternalLink(label, href string, props ComponentProps) Node`
- `ExternalIconLink(icon, ariaLabel, href string, props ComponentProps) Node`
- `DownloadLink(label, href, filename string, props ComponentProps) Node`
- `Input(name, value string, props ComponentProps) Node`
  - uses `InputWithOptions` with defaults:
    - `Type: "text"`
    - `Placeholder: strings.TrimSpace(name)`.
- `InputWithOptions(name, value string, options InputOptions) Node`
  - blank `options.Type` defaults to `"text"`.
- `Textarea(name, value string, options TextareaOptions) Node`
  - `Rows <= 0` defaults to `3`.
- `Form(action string, children ...Node) *form`
  - renders `<form>` with `ID`, `Class`, `Method`, `Action`, and passthrough `Attrs`.
- `ActionForm(props ActionFormProps, children ...Node) Node`
  - renders a form wired to Marionette/HTMX action updates.
  - blank `Method` defaults to `post`; supported methods are `post` and `get`.
  - renders standard `action`/`method` attributes plus `hx-post` or `hx-get`.
  - optional `Target` and `Swap` render `hx-target` and `hx-swap`.
- `HiddenField(name, value string) Node`
  - renders a hidden form field.
- `FormField(control Node, props FormFieldProps) Node`
  - if `control` rendering fails, returns render error node.
- `Select(name string, options []SelectOption, props ComponentProps) Node`

### Overlay / feedback
- `Modal(props ModalProps) Node`
  - renders `Body` and `Actions` nodes first.
  - if either render fails, returns render error node.
- `Toast(props ToastProps) Node`
  - blank `Live` defaults to `"polite"`.
- `Alert(props AlertProps) Node`
- `Skeleton(props SkeletonProps) Node`
  - `Rows <= 0` defaults to `3`.
- `Progress(props ProgressProps) Node`
  - renders a native `<progress>` element with optional label and percentage text.
  - `Max <= 0` defaults to `100`; `Value` is clamped into `0..Max`.
  - `Indeterminate` omits the `value` attribute so browsers render an indeterminate indicator.
- `EmptyState(props EmptyStateProps) Node`
  - `Rows <= 0` defaults to `3`.

### Data display
- `TableRowValues(values ...any) TableRow`
  - converts `nil` to empty text, `Node` values directly, and other values with `fmt.Sprint`.
- `Table(props TableProps) Node`
  - renders each cell node; any cell render error => render error node.
- `Chart(props ChartProps) Node`
  - renders a Chart.js-backed chart from Go props.
  - blank `Type` defaults to `ChartTypeLine`; blank `Height` defaults to `320`.
  - `ChartDataset.Data` renders scalar values; `ChartDataset.Points` renders `{x,y}` values for scatter-style charts.
  - chart config is JSON-encoded and embedded next to a `<canvas data-mrn-chart>`.
  - includes `role="img"`, an accessible label, canvas fallback text, and a screen-reader fallback table.
- `Image(props ImageProps) Node`
  - renders a responsive `<figure>` with an `<img>` and optional caption.
  - `Src` is required; blank `Src` returns a render error node.
  - blank `Loading` defaults to `"lazy"` and blank `Decoding` defaults to `"async"`.
  - `AspectRatio` supports `square`, `video`, `wide`, and `portrait`; `ObjectFit` supports `cover`/blank, `contain`, `fill`, `none`, and `scale-down`.
- `DataFrame(df *dataframe.DataFrame, props TableProps) Node`
  - renders `github.com/rocketlaunchr/dataframe-go` dataframes through `Table`.
  - `df.Names()` is mapped to `TableColumn.Label` and overrides `props.Columns`.
  - each row is read by `df.Row(row, true, dataframe.SeriesName)`.
  - cell conversion: `nil` => empty text, `Node` => rendered directly, all others => `fmt.Sprint(value)`.
- `DataFrameChart(df *dataframe.DataFrame, props DataFrameChartProps) Node`
  - maps a dataframe label column and numeric series columns into `Chart`.
  - blank `LabelColumn` uses the first dataframe column.
  - blank `Series` renders every column after the label column as a dataset.
- `Pagination(props PaginationProps) Node`
  - `Page < 1` defaults to `1`.
  - `TotalPages < 1` defaults to `1`.
- `Tabs(props TabsProps) Node`
  - blank `AriaLabel` defaults to `"tabs"`.
  - supports active/disabled states and link/button tab items.
- `Breadcrumb(props BreadcrumbProps) Node`
  - blank `AriaLabel` defaults to `"breadcrumb"`.
  - supports active/current breadcrumb items.
- `Checkbox(props CheckboxComponentProps) Node`
- `RadioGroup(props RadioGroupComponentProps) Node`
  - blank `AriaLabel` defaults to `"radio group"`.
- `Switch(props SwitchComponentProps) Node`
- `Badge(props BadgeProps) Node`
  - renders a compact label with `Variant`, `Size`, and custom classes from `ComponentProps`.
- `UIText(props TextProps) Node`
  - renders plain text with semantic size, weight, and tone options.
- `DataFrameFromCSV(r io.ReadSeeker, props TableProps, opts ...imports.CSVLoadOptions) (Node, error)`
  - loads CSV via `github.com/rocketlaunchr/dataframe-go/imports.LoadFromCSV`.
- `DataFrameFromTSV(r io.ReadSeeker, props TableProps, opts ...imports.CSVLoadOptions) (Node, error)`
  - same loader with `Comma: '\t'` as default.

### Layout / surfaces
- `Actions(props ActionsProps, children ...Node) Node`
  - renders a horizontal action group.
  - `Align`: `start`/blank, `center`, `end`, `between`.
  - `Gap` uses the same values as `Stack`; `Wrap` adds `flex-wrap`.
- `Divider(props DividerProps) Node`
  - renders a visual divider.
  - `Spacing`: `none`, `xs`, `sm`, `md`/blank, `lg`.
- `Box(props BoxProps, children ...Node) Node`
  - renders a generic surface with optional border, tone, padding, and custom classes.
- `AppShell(props AppShellProps) Node`
  - renders the demo/admin shell with sidebar, flashes, header, and content regions.
- `Stack(props StackProps, children ...Node) Node`
  - flex layout for vertical/horizontal stacks.
  - `Direction`: `vertical`/blank or `horizontal`/`row`.
  - `Gap`: `none`, `xs`, `sm`, `md`/blank, `lg`, `xl`.
  - `Align`: `start`, `center`, `end`, blank=`stretch`.
  - `Justify`: `start`/blank, `center`, `end`, `between`.
  - `Wrap` adds `flex-wrap`; `Props.Class` appends custom classes.
- `Grid(props GridProps, children ...Node) Node`
  - grid layout with `Columns` values `1`, `2`, `3`/blank, `4`.
  - `MinColumnWidth`: `sm`, `md`, `lg` switches to auto-fit minmax columns.
  - `Gap` and `Props.Class` use the same behavior as `Stack`.
- `Split(props SplitProps) Node`
  - responsive main/aside layout.
  - `AsideWidth`: `sm`, `md`/blank, `lg`.
  - `ReverseOnMobile` renders the aside before the main pane visually on mobile.
- `PageHeader(props PageHeaderProps) Node`
  - renders title, description, and optional action node.
- `Container(props ContainerProps, children ...Node) Node`
  - `MaxWidth`: `sm`, `md`, `lg`/blank, `full`.
  - `Padding`: `none`, `sm`, `md`/blank, `lg`.
  - `Centered` adds `mx-auto`.
- `Region(props RegionProps, children ...Node) Node`
  - renders an ID-addressable content region for partial updates.
  - `ID` is required; blank `ID` returns a render error node.
  - `Props.Class` appends custom classes.
- `Card(props CardProps, children ...Node) Node`
  - card surface with optional title, description, and action node.
  - `Gap` controls spacing between card body children.
- `Section(props SectionProps, children ...Node) Node`
  - unframed section wrapper with optional title, description, and action node.

#### Example: Convert CSV/TSV data to `DataFrame`

```go
import (
    "os"

    marionette "github.com/YoshihideShirai/marionette"
)

func tableFromCSV(path string) (mrn.Node, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    return mrn.DataFrameFromCSV(f, mrn.TableProps{
        EmptyTitle:       "No data",
        EmptyDescription: "CSV is empty.",
    })
}

func tableFromTSV(path string) (mrn.Node, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    return mrn.DataFrameFromTSV(f, mrn.TableProps{
        EmptyTitle:       "No data",
        EmptyDescription: "TSV is empty.",
    })
}
```

#### Example: Render a chart

```go
chart := mrn.Chart(mrn.ChartProps{
    Type:        mrn.ChartTypeLine,
    Title:       "Weekly signups",
    Description: "New accounts by weekday.",
    Labels:      []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
    Datasets: []mrn.ChartDataset{
        {
            Label:           "Signups",
            Data:            []float64{12, 19, 14, 22, 18},
            BorderColor:     "#2563eb",
            BackgroundColor: "rgba(37, 99, 235, 0.16)",
            Fill:            true,
            Tension:         0.3,
        },
    },
    Options: mrn.ChartOptions{
        BeginAtZero: true,
        YAxisLabel:  "Users",
    },
})
```

---

## 7. Flash APIs

### Types / constants
- `type FlashLevel string`
- `type FlashMessage struct { Level FlashLevel; Message string }`
- Levels:
  - `FlashSuccess`
  - `FlashError`
  - `FlashInfo`
  - `FlashWarn`

### UI helper
- `FlashAlerts(flashes []FlashMessage) Node`
  - empty flashes => `DivProps(ElementProps{ID: "flash-alerts", Class: "hidden"})`.
  - non-empty => `DivProps(ElementProps{ID: "flash-alerts", Class: "space-y-2"}, ...)` with level class mapping:
    - `FlashSuccess -> alert-success`
    - `FlashError -> alert-error`
    - `FlashWarn -> alert-warning`
    - default -> `alert-info`.

---

## 8. Runtime

### Handler type
```go
type Handler func(*Context) Node
```
- Shared by `Page`, `Action`, `Render`, `Handle`.

### HTTP / rendering failure behavior (consolidated)
- `Page` endpoint + non-`GET` method: `405 Method Not Allowed`.
- `Action` endpoint + non-`POST` method: `405 Method Not Allowed`.
- `Action` request `ParseForm` failure: `400 Bad Request`.
- Node render failure (page/fragment): `500 Internal Server Error`.
- Shell render failure (page): `500 Internal Server Error`.
- Missing root page (`/` not registered): `GET /` returns `500 Internal Server Error`.

### Response content type
- Both full page and fragment writes set:
  - `Content-Type: text/html; charset=utf-8`.
