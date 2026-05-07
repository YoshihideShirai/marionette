# DashWind API

`frontend/dashwind` は、Marionette と DaisyUI で DashWind 風の管理画面を作るための公開パッケージです。`cmd/dashwind-demo` で使っているシェル、ナビゲーション、ページヘッダー、カード、統計、テーブルの再利用可能な要素を切り出しています。

## インポート

```go
import dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
```

アプリ作成時に DaisyUI style template と DashWind アセットを 1 回登録します。

```go
app := mb.New()
dw.Use(app, dw.Options{})
```

`Options` では `Theme`、`CustomCSS`、`DisableDefaultCSS`、`AssetsBasePath` を指定でき、テーマ初期値、追加の trusted CSS、標準 CSS の無効化、DaisyUI/Tailwind import のセルフホスト化を制御できます。

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

## MetricGrid

`MetricGrid` は KPI カードのレスポンシブグリッドを描画します。単体カードには `MetricCard` を使い、`TrendTone` には `ToneSuccess`、`ToneWarning`、`ToneError`、`ToneNeutral` を指定できます。

```go
metrics := dw.MetricGrid(dw.MetricGridProps{Items: []dw.Metric{
    {
        Title:       "New Users",
        Value:       "34.7k",
        Description: "Acquired this period",
        Trend:       "↗︎ 2300 (22%)",
        TrendTone:   dw.ToneSuccess,
        Icon:        mf.Text("U"),
        Href:        "/analytics",
    },
}})
```

短い KPI dashboard サンプル:

```go
body := mf.DivProps(mf.ElementProps{Class: "space-y-6"},
    dw.PageHeader(dw.PageHeaderProps{Title: "KPI Dashboard", Description: "Today at a glance"}),
    dw.MetricGrid(dw.MetricGridProps{Items: []dw.Metric{
        {Title: "Revenue", Value: "$128K", Trend: "12% up", TrendTone: dw.ToneSuccess, Icon: mf.Text("$"), Href: "/revenue"},
        {Title: "Open Alerts", Value: "7", Trend: "3 urgent", TrendTone: dw.ToneWarning, Icon: mf.Text("!"), Href: "/alerts"},
        {Title: "Failed Jobs", Value: "2", Trend: "Needs review", TrendTone: dw.ToneError, Icon: mf.Text("×")},
        {Title: "NPS", Value: "61", Trend: "Stable", TrendTone: dw.ToneNeutral, Icon: mf.Text("★")},
    }}),
)
```

既存コード向けに `StatsGrid` と `Stat` も残していますが、新しいダッシュボードでは `MetricGrid` と `Metric` を使ってください。

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
