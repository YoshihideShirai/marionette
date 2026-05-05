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
- `H1(children ...Node) Node` / `H1Props(props ElementProps, children ...Node) Node`
- `H2(children ...Node) Node` / `H2Props(props ElementProps, children ...Node) Node`
- `H3(children ...Node) Node` / `H3Props(props ElementProps, children ...Node) Node`
- `H4(children ...Node) Node` / `H4Props(props ElementProps, children ...Node) Node`
- `Column(children ...Node) Node`

### Table / layout helpers
- `HTMXTable(headers []string, rows ...TableRowData) Node`
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
- `FileUpload(name string, required bool) Node`
- `HiddenInput(name, value string) Node`
- `Submit(label string) Node`
- `Button(label string) *button`
  - default target selector: `#app`.
- `(*button).Post(action string) *button` (action normalized without leading `/`)
- `(*button).OnClick(action string) *button` (`Post` alias)
- `(*button).Target(selector string) *button`
- `(*button).TargetSelector(selector string) *button`

---
