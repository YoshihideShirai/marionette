# DashWind API

`frontend/dashwind` は、Marionette と DaisyUI で DashWind 風の管理画面を作るための公開パッケージです。`cmd/dashwind-demo` で使っているシェル、ナビゲーション、ページヘッダー、カード、統計、テーブルの再利用可能な要素を切り出しています。

## インポート

```go
import dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
```

アプリ作成時にパッケージ CSS を 1 回追加します。

```go
app.AddStyle(dw.DefaultCSS)
```

## Shell、NavGroup、NavItem

`Shell` はレスポンシブなドロワーレイアウト、固定トップバー、検索ボックス、サイドバーナビゲーションを描画します。リンクをセクションに分けるには `NavGroup`、各ルートには `NavItem` を使います。

```go
page := dw.Shell(dw.ShellProps{
    Brand:             dw.Brand{Title: "DashWind", Subtitle: "Admin workspace"},
    CurrentPath:       "/",
    SearchPlaceholder: "Search workspace",
    Navigation: []dw.NavGroup{
        {Label: "Menu", Items: []dw.NavItem{
            {Label: "Dashboard", Href: "/", Icon: "▦"},
            {Label: "Leads", Href: "/leads", Icon: "▣"},
        }},
    },
    User: dw.UserMenu{Name: "Ada Lovelace", Email: "ada@example.com", Initials: "AL"},
}, dashboardBody)
```

htmx アクションでメイン領域だけを差し替える場合は、アクションハンドラーから `dw.ShellContent(dw.DefaultMainTargetID, body)` を返します。

## PageHeader

`PageHeader` はタイトル、説明文、任意のアクションノードを並べて描画します。

```go
header := dw.PageHeader(dw.PageHeaderProps{
    Title:       "Current Leads",
    Description: "Add, assign, and remove leads.",
    Actions:     addButton,
})
```

## CardPanel

`CardPanel` はコンテンツを DashWind/DaisyUI のカードサーフェスで包みます。

```go
panel := dw.CardPanel(dw.CardPanelProps{
    Title:       "Recent Transactions",
    Description: "Latest billing events",
}, table)
```

## StatsGrid

`StatsGrid` はメトリクスカードのレスポンシブグリッドを描画します。`Stat` の `Figure` にはアイコンや任意の Marionette ノードを渡せます。

```go
stats := dw.StatsGrid(dw.StatsGridProps{Items: []dw.Stat{
    {
        Title:            "New Users",
        Value:            "34.7k",
        Description:      "↗︎ 2300 (22%)",
        Figure:           mf.Text("U"),
        DescriptionClass: "font-medium text-success",
    },
}})
```

## DataTable

`DataTable` は DaisyUI のテーブルプリミティブを扱いやすくするラッパーです。ヘッダーは文字列、セルは Marionette ノードで指定します。

```go
rows := [][]mf.Node{
    {mf.Text("INV-8842"), mf.Text("Acme Inc."), mf.Text("Paid")},
}

table := dw.DataTable(dw.DataTableProps{
    Headers: []string{"Invoice", "Customer", "Status"},
    Rows:    rows,
    Class:   "table-zebra",
})
```
