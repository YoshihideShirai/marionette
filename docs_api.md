# API Documentation

This document is a direct API reference for the current runtime surface in the
module root package, reorganized by runtime layer.

## Import Path

Split imports by runtime role:

```go
import (
    mb "github.com/YoshihideShirai/marionette/backend"
    mf "github.com/YoshihideShirai/marionette/frontend"
)
```

- Recommended aliases: `mb` (marionette backend), `mf` (marionette frontend).
- Use `mb` for app/runtime APIs such as `New`, `App`, `Context`, `Handler`.
- Use `mf` for UI node APIs such as `Node`, `Div`, `Element`, `Table`.

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

### `SetCookieSecure(secure bool)`
- Enables/disables `Secure` on the flash cookie (`marionette_flash`).
- Default is `false`.
- Affects both:
  - flash write (`Context.AddFlash`)
  - flash clear (when flashes are consumed on next request).

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

---

## 4. Core Node

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
- `ComponentButton(label string, props ComponentProps) Node`
- `ComponentSubmitButton(label string, props ComponentProps) Node`
- `ComponentInput(name, value string, props ComponentProps) Node`
  - uses `ComponentInputWithOptions` with defaults:
    - `Type: "text"`
    - `Placeholder: strings.TrimSpace(name)`.
- `ComponentInputWithOptions(name, value string, options InputOptions) Node`
  - blank `options.Type` defaults to `"text"`.
- `ComponentTextarea(name, value string, options TextareaOptions) Node`
  - `Rows <= 0` defaults to `3`.
- `ComponentFormField(control Node, props FormFieldProps) Node`
  - if `control` rendering fails, returns render error node.
- `ComponentSelect(name string, options []SelectOption, props ComponentProps) Node`

### Overlay / feedback
- `ComponentModal(props ModalProps) Node`
  - renders `Body` and `Actions` nodes first.
  - if either render fails, returns render error node.
- `ComponentToast(props ToastProps) Node`
  - blank `Live` defaults to `"polite"`.
- `ComponentAlert(props AlertProps) Node`
- `ComponentSkeleton(props SkeletonProps) Node`
  - `Rows <= 0` defaults to `3`.
- `ComponentEmptyState(props EmptyStateProps) Node`
  - `Rows <= 0` defaults to `3`.

### Data display
- `ComponentTable(props TableProps) Node`
  - renders each cell node; any cell render error => render error node.
- `ComponentDataFrame(df *dataframe.DataFrame, props TableProps) Node`
  - renders `github.com/rocketlaunchr/dataframe-go` dataframes through `ComponentTable`.
  - `df.Names()` is mapped to `TableColumn.Label` and overrides `props.Columns`.
  - each row is read by `df.Row(row, true, dataframe.SeriesName)`.
  - cell conversion: `nil` => empty text, `Node` => rendered directly, all others => `fmt.Sprint(value)`.
- `ComponentPagination(props PaginationProps) Node`
  - `Page < 1` defaults to `1`.
  - `TotalPages < 1` defaults to `1`.
- `ComponentTabs(props TabsProps) Node`
  - blank `AriaLabel` defaults to `"tabs"`.
  - supports active/disabled states and link/button tab items.
- `ComponentBreadcrumb(props BreadcrumbProps) Node`
  - blank `AriaLabel` defaults to `"breadcrumb"`.
  - supports active/current breadcrumb items.
- `ComponentCheckbox(props CheckboxComponentProps) Node`
- `ComponentRadioGroup(props RadioGroupComponentProps) Node`
  - blank `AriaLabel` defaults to `"radio group"`.
- `ComponentSwitch(props SwitchComponentProps) Node`
- `ComponentDataFrameFromCSV(r io.ReadSeeker, props TableProps, opts ...imports.CSVLoadOptions) (Node, error)`
  - loads CSV via `github.com/rocketlaunchr/dataframe-go/imports.LoadFromCSV`.
- `ComponentDataFrameFromTSV(r io.ReadSeeker, props TableProps, opts ...imports.CSVLoadOptions) (Node, error)`
  - same loader with `Comma: '\t'` as default.

#### Example: Convert CSV/TSV data to `ComponentDataFrame`

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

    return mrn.ComponentDataFrameFromCSV(f, mrn.TableProps{
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

    return mrn.ComponentDataFrameFromTSV(f, mrn.TableProps{
        EmptyTitle:       "No data",
        EmptyDescription: "TSV is empty.",
    })
}
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
