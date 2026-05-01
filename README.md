# Marionette

![Marionette concept art](docs/assets/concept.png)

Marionette is a Go-first framework for building admin UIs and internal tools.
It lets you describe screens, state, and actions in Go while htmx handles
partial updates in the browser.

## Why Marionette

- Build operational UI without leaving Go.
- Keep routing, state updates, and event handlers on the server.
- Use htmx-powered partial rendering instead of maintaining a full SPA.
- Compose admin screens from pages, forms, actions, tables, charts, and layout components.
- Run the same app as a web UI or inside a desktop WebView shell.

## Try it

Run the full demo:

```bash
go run ./cmd/marionette
```

Then open http://127.0.0.1:8080.

Run the minimal sample:

```bash
go run ./cmd/simple-sample
```

Then open http://127.0.0.1:8081.

## Documentation

The README is intentionally small. Use the documentation site for tutorials,
API details, and component examples:

- Docs site: https://yoshihideshirai.github.io/marionette/
- Tutorial: https://yoshihideshirai.github.io/marionette/en/tutorial/
- API docs: https://yoshihideshirai.github.io/marionette/en/api/
- Components gallery: https://yoshihideshirai.github.io/marionette/en/components/

Japanese docs are available from the language switcher on the site.

## Development

Use [Air](https://github.com/air-verse/air) to restart the demo app when Go
files change:

```bash
go install github.com/air-verse/air@latest
air
```

Run the documentation site locally:

```bash
cd docs/site-astro
npm install
npm run dev
```

The GitHub Pages workflow publishes `docs/site-astro/` via GitHub Actions.
