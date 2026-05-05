# daisyUI coverage audit (2026-05-05)

## Scope
- Official reference: https://daisyui.com/components/ (v5.5.19 page, 65 components)
- Local implementation: `frontend/daisyui/components_core.go`, `frontend/daisyui/components_extra.go`

## Result summary
- **Component-level coverage**: 65 / 65 components have at least one corresponding helper in this repository.
- **Variation-level coverage**: many official variations (size/color/style/modifier/placement/direction) are not exposed as explicit typed helpers and are only reachable via `props.Class` (or not wired at all for some components).

## Components with explicit variation handling found
- Toggle: `ToggleVariant` supports color variants (`primary`, `secondary`, `accent`, `neutral`, `info`, `success`, `warning`, `error`).
- Tabs: `Tabs` supports active state per tab (`tab-active`).
- Modal: `Modal` supports open state (`modal-open`).
- Steps: `Step` supports active-like coloring (`step-primary`).

## Notable gaps vs official variation docs
- **Button**: no dedicated helpers for size/style/color variants (`btn-xs/sm/lg/xl`, `btn-outline`, etc.).
- **Input / Select / Textarea**: no dedicated typed size/color/style helpers; largely class-string driven.
- **Checkbox / Radio / Range / Rating / Toggle**: only Toggle has partial typed variants; other input families mostly missing typed variant APIs.
- **Alert / Badge / Progress / Tooltip / Status / Stat / Table / Timeline / Toast**: official placement/style/color/size/direction variants are not comprehensively modeled with typed helpers.
- **Drawer / Dropdown / Menu / Navbar / Pagination / Tabs / Steps**: orientation/placement/state variants are partially or not explicitly represented.
- **Mockups (browser/code/phone/window)**: base helpers exist; official sub-parts/modifiers are not fully modeled.

## Conclusion
- If the requirement is "official component + official documented variations" parity, this codebase is **not yet complete**.
- To remove gaps, add typed option structs/enums per component for official variants instead of relying on ad-hoc `Class` strings.
