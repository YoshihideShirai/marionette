# Frontend Policy Audit

This audit checks the current implementation against the architecture policy:

- `frontend/html`: primitive HTML/HTMX tags and basic tag utilities (internal layer)
- `frontend/daisyui`: definitions for daisyUI official components only
- Users should not need to directly depend on `frontend/html`

## Findings

## 1) `frontend` package re-exports low-level `frontend/html` primitives

### Evidence
- `frontend/ui_impl.go` aliases low-level types from `frontend/html`:
  - `type Node = lowhtml.Node`
  - `type Attrs = lowhtml.Attrs`
  - `type ElementProps = lowhtml.ElementProps`
  - `type Raw = lowhtml.Raw`

### Why this is misaligned
Even if users do not import `frontend/html` directly, this still exposes the primitive layer through the public `frontend` API surface and weakens the intended internal boundary.

### Refactor direction
- Introduce frontend-owned interfaces/types and keep `frontend/html` details internal.
- Limit or deprecate direct exposure of primitive-specific types in public APIs.

---

## 2) `frontend/daisyui` includes many pass-through wrappers that are not daisyUI-specific definitions

### Evidence
Large sections of `frontend/daisyui/components.go` are direct delegations to generic `frontend` APIs (e.g. `Button`, `Card`, `Input`, `Form`, `Container`, `Text`, `H1`-`H4`, etc.).

### Why this is misaligned
The policy states `frontend/daisyui` should contain daisyUI component definitions. Pure pass-through wrappers for generic frontend primitives/components blur the boundary and make `frontend/daisyui` act as a broad alias layer.

### Refactor direction
- Keep only daisyUI-specific component constructors in `frontend/daisyui`.
- Move generic wrappers to `frontend` (or remove redundant aliases).
- For compatibility, mark wrappers as deprecated before removal.

---

## 3) `frontend/daisyui` mixes two styles: delegating wrappers and low-level node assembly

### Evidence
`components.go` contains both:
- delegation-style wrappers (`return frontend.UI...`), and
- direct `lowhtml.ElementNode` composition (e.g. Navbar/Hero and many others).

### Why this is misaligned
Mixed construction styles make the intended layer responsibility harder to maintain and review. It is unclear whether daisyUI components should be composed from frontend-level APIs or from low-level primitives.

### Refactor direction
- Pick one composition strategy for `frontend/daisyui` (recommended: daisyUI-specific composition helpers).
- Isolate repetitive low-level assembly with local helpers and conventions.

---

## Prioritized Refactor Backlog

1. **Boundary hardening (high impact)**
   - Stop leaking primitive types through `frontend` public aliases.
2. **API surface cleanup (high impact)**
   - Remove/deprecate non-daisyUI pass-through wrappers from `frontend/daisyui`.
3. **Consistency pass (medium impact)**
   - Standardize component construction style within `frontend/daisyui`.
