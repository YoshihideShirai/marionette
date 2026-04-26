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
