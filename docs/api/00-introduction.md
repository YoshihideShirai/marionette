# API Documentation

This document is a direct API reference for the current runtime surface in the
module root package, reorganized by runtime layer.

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
