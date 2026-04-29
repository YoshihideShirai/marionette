# ADR 0001: Documentation Ecosystem Structure

- Status: Accepted
- Date: 2026-04-29

## Context

Marionette needs both:

1. A user-friendly docs site for discovery and onboarding.
2. Stable engineering documentation for contributor consistency.

Mixing both concerns in a single surface makes navigation and maintenance harder.

## Decision

Adopt a 3-layer documentation ecosystem:

1. **Public docs (`docs/site-astro/`)**
   - Quickstart, framework overview, components, cookbook, changelog.
2. **Engineering docs (`docs/`)**
   - Architecture policy, component guidelines, future ADRs.
3. **Repo front door (`README.md`, `CONTRIBUTING.md`, `.github/`)**
   - Entry points, contribution flow, ownership, PR checklist.

## Consequences

### Positive

- Faster onboarding for new users.
- Clear ownership and review boundaries.
- Better long-term maintainability with explicit conventions.

### Trade-offs

- More files and structure to keep synchronized.
- Requires PR discipline to keep docs updated.

## Follow-ups

- Add CI checks for markdown links and formatting.
- Expand cookbook with concrete recipes linked to runnable examples.
