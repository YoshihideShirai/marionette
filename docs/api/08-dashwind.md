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
    CurrentTitle:      "Dashboard",
    BrandTitle:        "DashWind",
    BrandSubtitle:     "Admin workspace",
    SearchPlaceholder: "Search workspace",
    NavGroups: []dw.NavGroup{
        {Label: "Menu", Items: []dw.NavItem{
            {Label: "Dashboard", Href: "/", Icon: "▦"},
            {Label: "Leads", Href: "/leads", Icon: "▣"},
        }},
    },
    Content: dashboardBody,
})
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

## StatsGrid

`StatsGrid` renders a responsive grid of metric cards. Each `Stat` can provide an icon or any Marionette node as `Figure`.

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
