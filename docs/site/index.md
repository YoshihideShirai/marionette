# Marionette Documentation

## Framework Overview

Marionette is a Go-first UI framework for building admin interfaces. Application developers compose screens in Go, while htmx handles partial updates in the browser. The runtime is also designed to work in a WebView host for desktop use cases.

- **Goals (summary)**
  - Define UI using Go only.
  - Compose admin screens from pages, actions, forms, and partial updates.
  - Keep browser-side interaction simple with htmx.
  - Support desktop-oriented deployment via WebView hosting.
- **Architecture (summary)**
  - **Backend (Go):** routing, state updates, and event handling.
  - **Frontend (htmx):** transport and incremental HTML swaps.
  - **Styling:** Tailwind CSS + daisyUI (CDN).
  - **UI DSL (Go):** build Nodes declaratively and render via `html/template`.

Related documentation:

- API reference: [`docs_api.md`](https://github.com/YoshihideShirai/marionette/blob/main/docs_api.md)
- UI policy: [`docs/architecture/ui.md`](https://github.com/YoshihideShirai/marionette/blob/main/docs/architecture/ui.md)
- Components gallery: [`docs/site/components/index.md`](./components/index.md)
