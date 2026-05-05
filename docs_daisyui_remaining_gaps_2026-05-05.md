# daisyUI remaining variation gaps (after helper expansion)

As of 2026-05-05, these official component variation areas are still not modeled as typed helpers.

## High-priority remaining gaps
- Alert: color variants (`alert-info/success/warning/error`), outline/soft styles.
- Avatar: placeholder/status/ring sizing patterns.
- Carousel: `carousel-center`, `carousel-end`, vertical mode, snap control variants.
- Collapse/Accordion: `collapse-plus`, `collapse-open`, `collapse-close` APIs.
- Drawer: `drawer-end`, overlay/side variants.
- Dropdown: placement modifiers (`dropdown-top/end/left/right`), open state classes.
- Footer: direction/layout variants.
- Indicator: placement variants (`indicator-top/bottom/start/end`).
- Join/Pagination: explicit `join-vertical/horizontal` wrappers.
- Link: style/color variants (`link-primary` etc.).
- Loading: type variants (`loading-dots/ring/ball/bars/infinity`) + size variants.
- Menu: orientation + compact/dense variants.
- Mockups: variant-specific wrappers for browser/code/phone/window parts and states.
- Navbar: wrappers for responsive layout variants.
- Progress: size variants not yet explicit.
- Radial progress: size/thickness API options as typed params.
- Select/Input/Textarea: still missing some official style variants beyond current wrappers.
- Skeleton: `skeleton-text` helper missing.
- Stat: `stats-vertical` / richer part helpers.
- Status: size + color typed wrappers.
- Swap: rotate/flip modifiers + indeterminate part helper.
- Toggle: size variants (`toggle-xs..xl`) helper missing.
- Validator: `validator` class application helper for input/select/textarea.

## Components present but mostly class-string driven
- Breadcrumbs, Calendar, Chat, Countdown, Diff, Dock, FAB, Fieldset, Filter, Hover 3D,
  Hover Gallery, Kbd, Label, List, Mask, Theme Controller, Text Rotate.

## Recommendation
- Move from many ad-hoc helper signatures to per-component option structs with enum-like constants,
  so invalid combinations can be constrained and docs parity can be audited mechanically.
