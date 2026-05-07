# DashWind API

`frontend/dashwind` is a public package for building DashWind-style admin pages with Marionette and DaisyUI. It extracts the reusable shell, navigation, page header, card, stats, and table primitives that power `cmd/dashwind-demo`.

## Import

```go
import dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
```

Add the package CSS once when creating the app:

```go
app.AddStyle(dw.DefaultCSS)
```

## Shell, NavGroup, and NavItem

`Shell` renders the responsive drawer layout, sticky top bar, search box, and sidebar navigation. Use `NavGroup` to split links into labeled sections and `NavItem` for each route.

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

For htmx actions that replace only the main content area, return `dw.ShellContent(dw.DefaultMainTargetID, body)` from the action handler.

## PageHeader

`PageHeader` renders a title, secondary description, and optional action node.

```go
header := dw.PageHeader(dw.PageHeaderProps{
    Title:       "Current Leads",
    Description: "Add, assign, and remove leads.",
    Actions:     addButton,
})
```

## CardPanel

`CardPanel` wraps content in a DashWind/DaisyUI card surface.

```go
panel := dw.CardPanel(dw.CardPanelProps{
    Title:       "Recent Transactions",
    Description: "Latest billing events",
}, table)
```

## MetricGrid

`MetricGrid` renders a responsive grid of KPI cards. Use `MetricCard` when you need a single card, and set `TrendTone` with `ToneSuccess`, `ToneWarning`, `ToneError`, or `ToneNeutral`.

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

Short KPI dashboard sample:

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

`StatsGrid` and `Stat` remain available for older code, but new dashboards should use `MetricGrid` and `Metric`.

## DataTable

`DataTable` is a convenience wrapper around the DaisyUI table primitive, using string headers and Marionette nodes for cells.

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
