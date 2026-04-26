# ui/overlay

`ui/overlay` provides a shared implementation layer for `Modal`, `Drawer`, `Popover`, and `Tooltip`.

## overlay-core

`overlay-core.js` provides base logic shared by all overlay variants.

- Close on Esc key
- Configurable backdrop-click close behavior
- Focus trapping and focus restoration on close
- Unified portal mount target (`#ui-overlay-root`)
- Reference-counted scroll locking

### API

```js
import { createOverlayCore } from './overlay-core.js';

const core = createOverlayCore({ portalId: 'ui-overlay-root' });
const layer = core.open({
  type: 'modal',
  panel,
  closeOnEsc: true,
  closeOnBackdrop: true,
  lockBodyScroll: true,
  zIndex: 1300,
  restoreFocusTo: triggerElement,
  onRequestClose: (reason) => console.log(reason),
});

layer.close('programmatic');
```

## z-index rules

Reuse the token scale in `tokens.css` across other components.

- `--ui-overlay-z-base: 1000`
- `--ui-overlay-z-popover: 1100`
- `--ui-overlay-z-tooltip: 1200`
- `--ui-overlay-z-modal: 1300`
- `--ui-overlay-z-drawer: 1400`

> Rule: place blocking overlays (`Modal`/`Drawer`) above lightweight hint overlays (`Popover`/`Tooltip`).

## Scroll-lock behavior

Overlays with `lockBodyScroll: true` apply the following behavior.

1. Set `data-overlay-scroll-lock="true"` on `body`.
2. Lock the viewport with `position: fixed` + `top: -scrollY`.
3. Manage nested overlays through a reference count.
4. When the final overlay closes, release styles and restore scroll position.

This prevents scroll-position corruption when `Modal` and `Drawer` overlap.
