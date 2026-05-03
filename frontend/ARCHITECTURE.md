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

## Separation Rules
- **Primitive layer**: `frontend/html`
- **Design-system layer (daisyUI)**: `frontend/daisyui`
- User-facing APIs should hide `frontend/html` internal details whenever possible.
- daisyUI-specific implementations should remain encapsulated in `frontend/daisyui` and not leak into other layers.
