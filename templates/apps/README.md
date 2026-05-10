# Go app layout templates

These templates are small DashWind-based Marionette applications you can copy
when starting a new Go-first UI. Each directory is intentionally self-contained
and runnable.

```bash
go run ./templates/apps/minimal
go run ./templates/apps/crud-list
go run ./templates/apps/dashboard
go run ./templates/apps/settings-form
go run ./templates/apps/master-detail
```

The default ports are different so you can compare layouts side by side:

- `minimal`: `127.0.0.1:8090`
- `crud-list`: `127.0.0.1:8091`
- `dashboard`: `127.0.0.1:8092`
- `settings-form`: `127.0.0.1:8093`
- `master-detail`: `127.0.0.1:8094`

## Choosing a template

- Start with `minimal` when you only need a page shell and a few content blocks.
- Start with `crud-list` for forms, action handlers, and table fragments.
- Start with `dashboard` for KPI cards, grids, and operational summaries.
- Start with `settings-form` for configuration screens and validation feedback.
- Start with `master-detail` for a list with a focused detail panel.

Copy the directory into `cmd/<your-app>/`, change the title, port, data model,
navigation, and action names, then run it with `go run ./cmd/<your-app>`.
