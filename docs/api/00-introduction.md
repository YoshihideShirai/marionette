# API Documentation

This document is a direct API reference for the current runtime surface,
organized by runtime layer. The module root no longer exposes a Go package;
import the layer-specific packages below.

## Import Path

Split imports by runtime role:

```go
import (
    mb "github.com/YoshihideShirai/marionette/backend"
    md "github.com/YoshihideShirai/marionette/desktop"
    mf "github.com/YoshihideShirai/marionette/frontend"
    mh "github.com/YoshihideShirai/marionette/frontend/html"
)
```

- Recommended aliases: `mb` (marionette backend), `mf` (marionette frontend).
- Use `mb` for app/runtime APIs such as `New`, `App`, `Context`, `Handler`.
- Use `md` for desktop runtime APIs such as `desktop.Run`.
- Use `mf` for component APIs such as `Button`, `Card`, `Table`, `FormRow`.
- Use `mh` for advanced low-level node APIs such as `Node`, `Div`, `Element`, `Raw`.
- Low-level constructors are intentionally not exposed from `frontend`; import
  `frontend/html` when custom markup is needed.
