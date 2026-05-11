# Data Tables / Charts Guide

English | [日本語](ja/06-data-tables-charts.md)

This document explains how to design screens that combine list tables, search conditions, paging, sorting, and charts in Marionette applications. It focuses on keeping responsibilities clear when a table and chart share the same `DataQueryState`, such as filtering a list from a chart click.

For general state lifetime and storage decisions, see the [State Management Guide](../state-management.md).

## Basic policy

For data list screens, the standard pattern is to **treat display conditions as explicit query state and separate data loading, data shaping, and rendering**.

- A Page builds display conditions from URL query parameters, session/server state, and defaults.
- Services / repositories receive search conditions, paging, and sorting, then load only the required range from the DB or external APIs.
- State stores the information needed to reproduce the display: search conditions, page number, page size, sort order, selected chart filters, and similar values.
- UI helpers receive loaded rows, columns, pagination metadata, and chart data, then build HTML.
- Actions receive inputs such as search form submissions, page changes, sort changes, and chart clicks, build the next query state, and return to the same rendering path.

Small CSV or in-memory dataframe examples can apply `DataFrameViewProps` in the UI layer. For business data or large result sets, pass conditions into DB queries and limit the amount of data loaded.

## 1. Basic list table structure

List tables are easier to maintain when the following elements are collected into one screen model.

```text
internal/<app>/state/order_list.go      # Query state, pagination, sort, selected filters
internal/<app>/pages/orders.go          # GET display, URL query parsing, initial data loading
internal/<app>/actions/order_list.go    # search / sort / page / chart-click actions
internal/<app>/service/orders.go        # DB search, count, aggregation, permission checks
internal/<app>/ui/order_list.go         # table, filters, pager, chart layout helpers
```

At minimum, the screen model should include:

- Input values to show again in the search form.
- The list of applied filters.
- Current page, page size, total row count, and whether previous / next pages exist.
- Current sort column and direction.
- Table rows and columns.
- Messages for empty states, load failures, permission failures, and similar cases.
- Chart labels / datasets / aggregation results rendered from the same conditions.

Example:

```go
type OrderListState struct {
    Search   string
    Status   string
    Region   string
    Page     int
    PageSize int
    Sort     string
    SortDesc bool

    TotalRows int
    Rows      []OrderRow
    Chart     frontend.ChartProps
    Error     string
}
```

The UI helper should make the search form, table, pager, and chart read from the same `OrderListState`. This lets initial display, redisplay after search, page changes, sort changes, and chart-filtered views use the same function.

## 2. Search conditions, paging, and sorting

### Search conditions

Separate "values entered by the user" from "filters passed to the DB or dataframe".

- Trim input, convert types, and check allowed values.
- Do not create filters for empty or unspecified values.
- Check select / radio values against an allowlist.
- Validate date-range consistency, such as start date not being after end date.
- Always add permission-enforced conditions in services / repositories separately from user input.

```go
func (s OrderListState) Filters() []frontend.DataFrameFilter {
    filters := []frontend.DataFrameFilter{}
    if strings.TrimSpace(s.Status) != "" {
        filters = append(filters, frontend.DataFrameFilter{
            Column: "status",
            Op:     frontend.DataFrameFilterEq,
            Value:  s.Status,
        })
    }
    if strings.TrimSpace(s.Search) != "" {
        filters = append(filters, frontend.DataFrameFilter{
            Column: "customer_name",
            Op:     frontend.DataFrameFilterContains,
            Value:  s.Search,
        })
    }
    return filters
}
```

### Paging

Paging should always have an upper bound.

- Normalize `Page` as a 1-based value.
- Clamp `PageSize` to allowed choices or a maximum value.
- When search conditions change, usually reset `Page` to `1`.
- For DB-backed screens, use `LIMIT` / `OFFSET` or keyset pagination instead of passing all rows to the UI.
- For expensive counts, decide whether the screen needs an exact total count or only needs to know whether a next page exists.

```go
func normalizePage(page, pageSize int) (int, int) {
    if page < 1 {
        page = 1
    }
    switch pageSize {
    case 20, 50, 100:
    default:
        pageSize = 20
    }
    return page, pageSize
}
```

### Sorting

Do not directly couple display labels to DB / dataframe column names. Convert sort input through allowed sort keys.

- Use URL-friendly values such as `sort=created_at` or `sort=-total_amount`.
- Fall back to the default sort when an unallowed column is requested.
- Add a secondary sort when many rows can have the same value, so ordering is stable.
- Decide per screen whether changing sort should reset `Page` to `1` or keep the current page.

```go
var orderSortColumns = map[string]string{
    "created_at":   "orders.created_at",
    "total_amount": "orders.total_amount",
    "status":       "orders.status",
}
```

## 3. Design policy when using `DataQueryState`

Use `DataQueryState` as the minimal representation of data display conditions shared by tables and charts. The current type is centered on a set of filters, so do not put the entire search form, paging, sorting, permission conditions, or DB connection details into it.

Good use cases:

- Applying the category selected in a chart to a table.
- Sharing the same filters across multiple widgets.
- Applying filters to `DataFrameViewProps` so examples or small data views render with the same conditions.
- Encoding shared filters into URL query parameters such as `df.filter`, so the same view can be restored after reload.

Avoid:

- Using `DataQueryState` as the only source of truth for a business search form.
- Relying on UI-derived filters for security boundaries such as user ID, permissions, or tenant ID.
- Loading a large dataset and filtering it only through `DataQueryState` afterward.
- Saving display-only filters without distinguishing them from business conditions that should be durable.

For screens that use `DataQueryState`, separate screen-specific state from shared query state.

```go
type SalesDashboardState struct {
    // Screen-specific input values.
    DateFrom string
    DateTo   string
    Region   string

    // Display filters shared by the table and chart.
    Query frontend.DataQueryState

    Page     int
    PageSize int
    Sort     string

    Rows       []SalesRow
    DataFrame  *rdf.DataFrame
    Chart      frontend.ChartProps
    TotalRows  int
    PageError  string
}
```

When rendering, merge the shared filters into the base view.

```go
baseView := frontend.DataFrameViewProps{
    Page:     state.Page,
    PageSize: state.PageSize,
    Sort: []frontend.DataFrameSort{{
        Column: state.Sort,
        Desc:   strings.HasPrefix(state.Sort, "-"),
    }},
}
view := state.Query.ToView(baseView)
```

For DB-backed screens, pass column names and operators through allowlists when converting `DataQueryState` to repository input. Never concatenate UI-provided column names directly into SQL.

## 4. Linked UI example: filtering a table from a chart click

For screens where a chart and table interact, the flow should be: **add the clicked value to query state, then rerender both the chart and table from the same state**.

Typical flow:

1. The Page restores `DataQueryState` and screen-specific state from the URL query.
2. The Service loads table rows and chart aggregation with the same conditions.
3. The Chart sets `QueryStateName` and `QueryStateLabel` so the selected value's filter meaning is explicit.
4. The user clicks a chart segment / bar.
5. An Action or JavaScript converts the selected value, such as `region=Tokyo`, into `DataQueryState`.
6. The Page / Action returns to the same render helper and updates both the table and chart.

Example:

```go
func BuildSalesView(state SalesDashboardState) frontend.Node {
    queryView := state.Query.ToView(frontend.DataFrameViewProps{
        Page:     state.Page,
        PageSize: state.PageSize,
    })

    chart := state.Chart
    chart.QueryStateName = "sales-filter"
    chart.QueryStateLabel = "region"

    table := frontend.TableProps{
        Columns: []frontend.TableColumn{
            {Label: "Region", SortKey: "region"},
            {Label: "Revenue", SortKey: "revenue"},
        },
        View:           queryView,
        QueryStateName: "sales-filter",
    }

    return frontend.Div(nil,
        frontend.Chart(chart),
        frontend.DataFrame(state.DataFrame, table),
    )
}
```

Example filter added by a click:

```go
func AddChartFilter(state SalesDashboardState, column string, value string) SalesDashboardState {
    next := state
    next.Query.Filters = append(next.Query.Filters, frontend.DataFrameFilter{
        Column: column,
        Op:     frontend.DataFrameFilterEq,
        Value:  value,
    })
    next.Page = 1
    return next
}
```

This design avoids making the chart the only place that knows about the filter. Tables, charts, summary cards, and download links can all build their conditions from the same `DataQueryState`, which keeps the whole screen consistent.

## 5. Deciding whether state belongs in the URL, server state, DB, or cache

For list / chart screens, choose where state lives based on the lifetime you need to restore and which layer is the source of truth.

| Location | Good for | Avoid |
| --- | --- | --- |
| URL query | Search conditions, page, page size, sort, selected display filters | Secrets, large JSON values, permission conditions that would be dangerous if modified |
| Server state / session | Temporary wizard conditions, UI preferences that should not be in the URL, short-lived draft pointers | The source of truth for business data, large search result payloads |
| DB | Saved views, report definitions users return to later, business conditions that require auditability | Click-only temporary filters, aggregation caches that can be recomputed |
| Cache | Expensive aggregation results, external API responses, derived data reused with the same conditions | Data that cannot be restored if lost, authorization facts, user input that should be durable |

### Put it in the URL

URL query should be the first choice when all of the following are true.

- The same list should be restorable by bookmark, share, or reload.
- The value is short and easy to represent as text.
- The value is safe for the user to see.
- Server-side validation and permission checks make tampering safe to reject.

Example: `?q=alice&status=open&page=2&page_size=50&sort=-created_at`.

### Put it in server state / session

Consider server state or session for short-lived UI state that should not appear in the URL.

- Temporary display conditions across multiple steps.
- Dashboard layout preferences that would make the URL too long.
- Flash / selected tab state shown once after post-redirect-get.
- A short key pointing to a search-condition draft, not the full result set.

Avoid storing full search results or large dataframes in session because memory use, stale data, and cross-instance inconsistency become problems.

### Put it in the DB

Use durable storage such as the DB for the following state.

- Reports / saved searches that users name and save.
- Dashboard settings that must survive deploys, restarts, and logout.
- Conditions that affect business workflows, approval, billing, or audit trails.
- Definitions shared by multiple users or multiple app instances.

When running conditions saved in the DB, reevaluate current permissions, tenant scope, and visibility. Conditions that were safe when saved are not automatically safe later.

### Put it in cache

Use cache to speed up derived data.

- Expensive aggregations shown repeatedly with the same search conditions.
- External API responses reused briefly.
- Chart aggregations stored with a TTL.
- Values that can be recomputed from the DB or external API on a cache miss.

Include all values that affect the result in the cache key, such as search conditions, user / tenant / permission scope, aggregation version, and timezone.

## 6. Link to the state management guide

For more general state-management decisions, see the [State Management Guide](../state-management.md). The same criteria for request-local state, session / user state, durable state, cache / derived state, and `App` global state also apply to list tables and chart-linked UIs.
