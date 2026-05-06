## 6.5 daisyUI 便利コンポーネント

これらのヘルパーは `frontend` から再エクスポートされ、`frontend/daisyui`
からも利用できます。多くは daisyUI の標準的なマークアップへ直接対応し、
`Node` の子要素を受け取ります。

### タイポグラフィ・ボタン
- `TextNode(text string) Node`
  - `span` でテキストを描画。
- `PrimaryButton(label string, props ComponentProps) Node`
- `SecondaryButton(label string, props ComponentProps) Node`
- `GhostButton(label string, props ComponentProps) Node`
  - `Button` に `btn-primary` / `btn-secondary` / `btn-ghost` を付与するショートカット。
- `ButtonWithVariants(label, variant, size, style string, props ComponentProps) Node`
  - 空でない値に対して `btn-<variant>` / `btn-<size>` / `btn-<style>` を追加。

### ナビゲーション・メニュー・レイアウト
- `Navbar(start, center, end Node) Node`
  - `navbar-start` / `navbar-center` / `navbar-end` のスロットを出力。
- `Menu(items ...Node) Node`
  - `li` ではない項目を `li` でラップ。
- `Footer(children ...Node) Node`
- `Drawer(id string, side, content Node) Node`
  - `id` を drawer toggle input と overlay label の紐付けに使用。
- `Dropdown(trigger, menu Node) Node`
- `DropdownWithPlacement(trigger, menu Node, placement string) Node`
  - `placement` が空でない場合 `dropdown-<placement>` を追加。
- `Join(children ...Node) Node`
- `JoinWithDirection(direction string, children ...Node) Node`
  - `direction` が空でない場合 `join-<direction>` を追加。

### Hero・Avatar・Mockup
- `Hero(title, description string, actions ...Node) Node`
- `Avatar(src, alt, class string) Node`
  - 指定した wrapper class の中に `img` を置く daisyUI avatar マークアップ。
- `Stat(title, value, desc string) Node`
- `BrowserMockup(content Node) Node`
- `PhoneMockup(content Node) Node`
- `MockupWindow(title string, content Node) Node`
  - `title` は API の対称性のため受け取りますが、現時点では描画されません。
- `CodeMockup(lines ...string) Node`
  - 各行を `data-prefix="$"` 付きの mock terminal 行として描画。

### 開閉・ステップ・タイムライン・タブ・カルーセル
- `Accordion(title string, content Node, open bool) Node`
  - `name="accordion"` の radio ベース collapse。
- `Collapse(title string, content Node, open bool) Node`
  - checkbox ベース collapse。
- `Steps(items ...Node) Node`
- `Step(label string, active bool) Node`
  - `active` で `step-primary` を追加。
- `StepsWithVariants(items []Node, direction, color, className string) Node`
  - `steps-<direction>` と `className` を追加。`color` 指定時は既存の `step` 子要素へ `step-<color>` を付与。
- `Timeline(items ...Node) Node`
- `TimelineItem(startLabel, endLabel string, content Node) Node`
- `TimelineWithDirection(items []Node, direction string, compact bool, snapIcon bool, className string) Node`
  - `timeline-<direction>` / `timeline-compact` / `timeline-snap-icon` / `className` を条件付きで追加。
- `TabsWithVariants(items []TabsItem, style, placement, size, className string) Node`
  - `tabs-<style>` / `tabs-<placement>` / `tabs-<size>` / `className` を追加。
  - `TabsItem.Active` は `tab-active` と `aria-selected="true"` を設定。
  - `TabsItem.Disabled` は `tab-disabled` と `disabled` を設定。
- `Carousel(items ...Node) Node`
- `CarouselItem(id string, child Node) Node`

### 入力・フォーム補助
- `Range(name string, value int, min int, max int) Node`
- `RangeWithVariants(name string, value int, min int, max int, color, size string) Node`
  - `range-<color>` / `range-<size>` を追加。
- `Rating(name string, max int, checked int) Node`
- `RatingWithVariants(name string, max int, checked int, size string, half bool, allowClear bool) Node`
  - `rating-<size>` と `rating-half` を条件付きで追加。
  - `allowClear` で値 `0` の hidden radio を先頭に追加。
- `Toggle(name string, checked bool) Node`
- `ToggleVariant(name string, checked bool, variant string) Node`
  - `primary`, `secondary`, `accent`, `neutral`, `info`, `success`, `warning`, `danger`/`error` を認識。
- `ToggleWithVariants(name string, checked bool, color string, size string) Node`
  - `toggle-<color>` / `toggle-<size>` を追加。
- `ToggleWithIcons(name string, checked bool, className string) Node`
- `ThemeController(options ...Node) Node`
- `ThemeControllerOption(theme string, checked bool, className string) Node`
  - `theme-controller` 付き、`name="theme-buttons"` の radio input を出力。
- `Fieldset(legend string, fields ...Node) Node`
- `Label(text string) Node`
- `Validator(message string) Node`

### フィードバック・オーバーレイ・視覚効果
- `Tooltip(text string, child Node) Node`
- `TooltipWithVariants(text string, child Node, placement, color string, open bool) Node`
  - `tooltip-<placement>` / `tooltip-<color>` / `tooltip-open` を追加。
- `ToastWithPlacement(children []Node, horizontal, vertical, className string) Node`
  - `toast-<horizontal>` / `toast-<vertical>` / `className` を追加。
- `Loading(sizeClass string) Node`
  - 既定は `loading loading-spinner`。`sizeClass` が空でなければ追加。
- `LoadingWithVariants(kind string, size string) Node`
  - `kind` 空文字は `spinner`。`size` 指定時は `loading-<size>` を追加。
- `RadialProgress(value int, sizeClass string) Node`
- `Status(colorClass string) Node`
- `StatusWithVariants(color string, size string) Node`
  - `status-<color>` / `status-<size>` を追加。
- `Swap(onNode, offNode Node, active bool) Node`
- `SwapWithVariants(onNode, offNode Node, active bool, rotate bool, flip bool) Node`
  - flag に応じて `swap-active` / `swap-rotate` / `swap-flip` を追加。
- `FAB(icon Node, label string) Node`
  - `label` は button の `aria-label` として使用。
- `SpeedDial(trigger Node, items ...Node) Node`

### データ表示・メディア・ユーティリティ
- `List(items ...Node) Node`
  - `li` でない項目をラップし、list entry に `list-row` を追加。
- `Table(headers []string, rows ...[]Node) Node`
- `TableWithVariants(headers []string, rows [][]Node, zebra bool, pinRows bool, pinCols bool, size string) Node`
  - `table-zebra` / `table-pin-rows` / `table-pin-cols` / `table-<size>` を条件付きで追加。
- `Kbd(text string) Node`
- `Code(text string) Node`
- `Indicator(item, target Node) Node`
- `Mask(shapeClass string, child Node) Node`
- `ChatBubble(content Node, end bool) Node`
  - `end` で `chat-start` から `chat-end` に切り替え。
- `Countdown(value int) Node`
- `Dock(items ...Node) Node`
- `DockItem(child Node, active bool) Node`
  - `active` で `dock-active` を追加。
- `Filter(items ...Node) Node`
- `FilterItem(label string, active bool) Node`
  - `aria-label` に `label` を設定した radio input。
- `Calendar(content Node) Node`
- `CalendarGrid(days ...Node) Node`
- `Diff(before, after Node) Node`
- `TextRotate(words []string, animationClass string) Node`
- `Hover3DCard(content Node) Node`
- `HoverGallery(items ...Node) Node`
