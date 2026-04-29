# docs/site-astro

Astro-based static documentation prototype for GitHub Pages.

## Run

```bash
cd docs/site-astro
npm install
npm run dev
```

## Build

```bash
npm run build
```

Output is generated in `dist/` and can be deployed via GitHub Actions to Pages.

## `components-index.json` path conventions

In `src/data/components-index.json`, use the following conventions for each `components[]` entry:

- `example`: Public URL path (for example: `/marionette/examples/button.html`).
  - Do not apply path conversion in UI code.
  - If a component sample is not ready, set `example: null` explicitly.
- `template`: Repository-root-relative path (for example: `templates/components/button.tmpl`).
- `golden`: Repository-root-relative path (for example: `testdata/golden/button.golden.html`).

These conventions keep path semantics in data and reduce UI-side maintenance overhead.
