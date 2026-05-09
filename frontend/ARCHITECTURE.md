# Frontend Implementation Architecture

## Purpose
Clarify responsibilities within `frontend` and preserve clear boundaries between implementation layers and user-facing APIs.

## Directory Responsibilities

### `frontend/html`
- Defines primitive HTML and HTMX tags.
- Provides foundational utilities for building HTML tags.
- This directory is primarily for internal implementation and is not intended for direct use by Marionette users.

### `frontend/daisyui`
- Contains definitions only for components listed in official daisyUI.
- Centralizes definitions required to represent and compose daisyUI components.
- Primitive HTML construction logic should remain in `frontend/html` to avoid responsibility overlap.

### `frontend` (root)
- Provides user-facing aliases only.
- Should not contain daisyUI component rendering details directly.
- Public helpers in this layer should delegate to the daisyUI/UI implementation layer.

## Separation Rules
- **Primitive layer**: `frontend/html`
- **Design-system layer (daisyUI)**: `frontend/daisyui`
- User-facing APIs should hide `frontend/html` internal details whenever possible.
- daisyUI-specific implementations should remain encapsulated in `frontend/daisyui` and not leak into other layers.
- Keep root-level alias files (e.g. `component_aliases_impl.go`) free of rendering logic.
- Exception: compatibility-critical aliases may keep thin markup adapters when test-verified legacy HTML output (ARIA attributes, pagination labels, or theme toggle hooks) must remain stable during migration.

## Compatibility Exceptions in Current Migration
- `ThemeToggleButton`: the daisyUI implementation keeps the legacy `onclick="globalThis.mrnToggleTheme()"` hook while exposing the active `system`/`light`/`dark` mode through an icon-only control so existing shell JavaScript and tests keep working while the root API remains a thin alias.
- `EmptyState`: the daisyUI implementation keeps the legacy skeleton branch (`aria-busy`, `aria-live`, and skeleton row markup) because callers use the same public `EmptyStateProps` for both empty copy and loading placeholders.

## Component Template Loader Boundary

Decision: keep `templates/components/*.tmpl` and `templates/components/*.html` at the repository-level `templates/components` path for this migration. They are shared compatibility templates for the root `marionette` package and the `frontend` package, so moving them under `frontend/daisyui` would incorrectly make compatibility markup look daisyUI-package-private.

Responsibility boundary:
- `internal/componenttmpl` is the only owner of template discovery, parsing, execution, and cache ownership.
- `components_template_loader.go` and `frontend/components_template_loader_impl.go` may only resolve their package-relative template root and delegate rendering/loading to `internal/componenttmpl`.
- Component renderers may return a local `templateNode` adapter when they need a legacy shared template, but they must not parse files or manage template caches directly.
- `frontend/daisyui` should prefer Go `ElementNode` composition for daisyUI primitives and should not load repository templates.

## New Component Implementation Rules

Use a repository template (`templates/components/*.tmpl` or `*.html`) when all of the following are true:
- The markup is compatibility-sensitive for the root `marionette` package or shared between root and `frontend` APIs.
- The component has a mostly static HTML skeleton with conditional classes/attributes that are clearer in HTML template form.
- Golden output stability is more important than programmatic node composition.
- The component is not a daisyUI-only primitive owned exclusively by `frontend/daisyui`.

Use Go `ElementNode` composition when any of the following are true:
- The component is a daisyUI primitive or convenience helper that belongs to `frontend/daisyui`.
- The component mainly composes child `Node` values, dynamically assembles attributes, or branches over nested structures.
- The implementation needs type-safe reuse of low-level HTML helpers or is easier to test as node composition.
- The component has no root-package compatibility template requirement.

When both options are plausible, choose the smallest owner that matches the public API boundary: root/shared compatibility markup stays in templates through `internal/componenttmpl`; daisyUI-only behavior stays in `frontend/daisyui` as `ElementNode` composition.

## Mixed Implementation Inventory

Representative components that currently have both template-backed compatibility implementations and Go `ElementNode` daisyUI implementations:

| Component | Template-backed compatibility path | Go `ElementNode` path | Migration note |
| --- | --- | --- | --- |
| `Button` | `componentButton` returns `templateNode` for `templates/components/button.tmpl` in root and `frontend` render-core files. | `frontend/daisyui.Button`, plus variant helpers in `frontend/daisyui/components_extra.go`, build button nodes directly. | Keep the template for legacy API/golden compatibility; use daisyUI Go helpers for daisyUI-only aliases and new daisyUI helpers. |
| `Input` | Root `inputWithOptionsComponent` returns `templateNode` for `templates/components/input.tmpl`; the root `Input` path delegates there. | `frontend.Input` and `frontend/daisyui.Input` build `<input>` via `ElementNode`. | Do not add new template loaders; decide per API boundary whether legacy root output or daisyUI behavior owns the change. |
| `Card` | Root `Card` and `CardWithVariants` return `templateNode` for `templates/components/card.html`. | `frontend.Card` aliases to `frontend/daisyui.Card`, which builds card markup with `ElementNode`. | Treat root card template output as compatibility markup; evolve frontend daisyUI cards in Go. |
| `Grid` | Root `Grid` returns `templateNode` for `templates/components/grid.html`. | `frontend/components.Grid` and `frontend/daisyui.Grid` build grid containers with `ElementNode`. | Keep layout helpers close to their package owner; avoid moving the shared root template into daisyUI. |
| `Progress` | Root `Progress`/`ProgressWithVariants` return `templateNode` for `templates/components/progress.html`. | `frontend.Progress` aliases to `frontend/daisyui.Progress`, and daisyUI variant helpers build progress nodes in Go. | Preserve root template compatibility; prefer Go for daisyUI progress variants and dynamic attributes. |
