# Framework Overview

## What Marionette optimizes for

- Define UI in Go.
- Compose pages from handlers and reusable components.
- Use htmx for partial updates.
- Keep frontend logic minimal.

## Architecture at a glance

- **Backend (Go):** routing, state updates, event handlers.
- **Frontend (htmx):** transport and incremental swaps.
- **Styling:** Tailwind CSS + daisyUI.
- **Rendering:** `html/template` and Marionette UI DSL.

## Runtime model

- `Page` handlers return full pages.
- `Action` handlers return HTML fragments for partial swaps.
- State is managed server-side through context helpers.

## Read next

- [Quickstart](./getting-started.md)
- [Components gallery](./components/index.md)
- [UI Architecture Policy](../architecture/ui.md)
