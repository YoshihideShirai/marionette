# UI Architecture Policy

UI implementation in this project must follow the policies below.

## 1) UI component implementation

- UI components must be implemented with **Go templates (`html/template`)**.
- Do not rely on client-framework-first implementations.

## 2) New TypeScript files

- Adding new **`.ts` / `.tsx` files is prohibited**.
- If an exception is necessary, prior agreement with maintainers is required before implementation.

## 3) State transitions and validation

- Implement UI state transitions and input validation primarily in **Go handlers**.
- Keep client-side logic minimal (presentation-only support).

## Acceptance criteria

- Future PRs must not create new `.ts` / `.tsx` files.


## 4) Overlay components (`ui/overlay`)

Overlay UI (`Modal`, `Drawer`, `Popover`, `Tooltip`) must reuse the shared base logic in `ui/overlay/overlay-core.js`.

### Shared behavior requirements

- Close on `Esc` key by default (can be disabled per instance).
- Backdrop-click close must be configurable per instance.
- Focus trap while open, then restore focus to the trigger element on close.
- Render all overlay layers into a single portal root (`#ui-overlay-root`).

### z-index policy

Use the overlay token scale from `ui/overlay/tokens.css`:

- `--ui-overlay-z-base: 1000`
- `--ui-overlay-z-popover: 1100`
- `--ui-overlay-z-tooltip: 1200`
- `--ui-overlay-z-modal: 1300`
- `--ui-overlay-z-drawer: 1400`

Blocking overlays (`Modal`/`Drawer`) must be above hint overlays (`Popover`/`Tooltip`).

### Scroll-lock policy

`Modal` and `Drawer` should enable body scroll lock through `overlay-core`.

- Use reference-counted lock handling for nested overlays.
- Preserve and restore scroll position when the final lock is released.
- Never implement ad-hoc `overflow: hidden` toggles outside of `overlay-core`.
