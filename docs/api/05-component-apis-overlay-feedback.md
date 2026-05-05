## 6.2 Overlay, feedback, and state components

- `Modal(props ModalProps) Node`
  - Render failures in `Body`/`Actions` return a render error node.
- `Toast(props ToastProps) Node`
  - `Live`: `polite` (default), `assertive`, `off`
- `Alert(props AlertProps) Node`
- `Skeleton(props SkeletonProps) Node`
  - `Rows <= 0` defaults to `3`.
- `Progress(props ProgressProps) Node`
  - `Max <= 0` defaults to `100`.
  - `Value` is clamped to `0..Max`.
  - `Indeterminate: true` omits `value`.
- `EmptyState(props EmptyStateProps) Node`
  - `Rows <= 0` defaults to `3`.
