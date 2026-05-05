## 4. 低レベルHTML (`frontend/html`)

### `type Node interface { Render() (template.HTML, error) }`
- すべての UI ノードは安全な HTML を返します。

### 基本コンストラクタ
- `Text(v string) Node`
- `Raw(html string)`（`type Raw string`）
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

### テーブル / レイアウト
- `HTMXTable(headers []string, rows ...TableRowData) Node`
- `TableRow(cells ...Node) TableRowData`
- `Sidebar(brand, title string, items ...SidebarItem) *sidebar`
- `SidebarLink(label, href string) SidebarItem`
