## 6.2 オーバーレイ・通知・状態表示

- `Modal(props ModalProps) Node`
  - `Body` / `Actions` の描画失敗時はエラーノード。
- `Toast(props ToastProps) Node`
  - `Live`: `polite`（デフォルト）, `assertive`, `off`
- `Alert(props AlertProps) Node`
- `ToastWithPlacement(children []Node, horizontal, vertical, className string) Node`
- `TooltipWithVariants(text string, child Node, placement, color string, open bool) Node`
- `ModalWithPlacement(props ModalProps, placement string) Node`
- `LoadingWithVariants(kind, size string) Node`
- `SwapWithVariants(onNode, offNode Node, active, rotate, flip bool) Node`
- `Skeleton(props SkeletonProps) Node`
  - `Rows <= 0` は `3`。
- `Progress(props ProgressProps) Node`
  - `Max <= 0` は `100`。
  - `Value` は `0..Max` にクランプ。
  - `Indeterminate: true` は不定表示（`value` 属性なし）。
- `EmptyState(props EmptyStateProps) Node`
  - `Rows <= 0` は `3`。
