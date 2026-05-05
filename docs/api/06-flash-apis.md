## 7. Flash APIs

### Types / constants
- `type FlashLevel string`
- `type FlashMessage struct { Level FlashLevel; Message string }`
- Levels:
  - `FlashSuccess`
  - `FlashError`
  - `FlashInfo`
  - `FlashWarn`

### UI helper
- `FlashAlerts(flashes []FlashMessage) Node`
  - empty flashes => `DivProps(ElementProps{ID: "flash-alerts", Class: "hidden"})`.
  - non-empty => `DivProps(ElementProps{ID: "flash-alerts", Class: "space-y-2"}, ...)` with level class mapping:
    - `FlashSuccess -> alert-success`
    - `FlashError -> alert-error`
    - `FlashWarn -> alert-warning`
    - default -> `alert-info`.

---
