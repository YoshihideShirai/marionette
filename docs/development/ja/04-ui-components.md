# UI コンポーネント選定ガイド

[English](../04-ui-components.md) | 日本語

このドキュメントは、Marionette アプリケーションで画面を組み立てるときに、既存コンポーネント、`frontend/html` の低レベル HTML API、DaisyUI convenience component、新規コンポーネントのどれを選ぶかを判断するためのガイドです。admin UI や社内ツールで、見た目・アクセシビリティ・保守性を揃えながら実装速度を落とさないことを目的にします。

関連する詳細ルールは、先に次のドキュメントを確認してください。

- [UI Component Guidelines](../../ui-component-guidelines.md): コンポーネント選定、アクセシビリティ、DaisyUI 利用、状態表現の基本方針。
- [frontend/ARCHITECTURE.md](../../../frontend/ARCHITECTURE.md): `frontend` 配下の責務分離、template と Go 実装の境界、互換性ルール。
- [Low-level HTML API](../../api/ja/03-low-level-html.md): `frontend/html` の primitive API。
- [DaisyUI convenience component API](../../api/ja/05-component-apis-daisyui-convenience.md): DaisyUI ベースの convenience component。

## 基本方針

原則として、**既存の Marionette `frontend` コンポーネントを優先**し、足りない部分だけを低レベル API やアプリ固有 helper で補います。新規コンポーネントを追加するのは、既存 API の組み合わせでは意味・再利用性・アクセシビリティを保ちにくい場合に限定します。

選定の順序は次の通りです。

1. 既存の `frontend` コンポーネントで表現できるか確認する。
2. DaisyUI の既存 convenience component で自然に表現できるか確認する。
3. 画面固有の小さな差分だけなら、`frontend/html` で局所的に補う。
4. 複数画面で同じ意味の UI が繰り返されるなら、アプリ内 UI helper を作る。
5. Marionette 本体に入れるべき汎用性があり、責務の置き場所が明確な場合だけ新規 built-in component を検討する。

## 1. 既存コンポーネントを優先する基準

次の条件に当てはまる場合は、まず既存コンポーネントを使います。

- **意味が一致している**: button、link、table、card、alert、modal、form row など、ユーザーに見える UI の意味が既存 component 名と一致する。
- **アクセシビリティが必要**: `aria-*`、label、disabled、focus、role などの属性を毎回手書きすると漏れやすい。
- **DaisyUI / Tailwind の class 体系に乗せたい**: 色、サイズ、余白、状態表現を既存の visual language に揃えられる。
- **テスト済みの HTML 出力に乗りたい**: 既存 component は互換性や golden test の対象になっている場合があり、直接 HTML を書くより安全。
- **複数画面で同じ見た目にしたい**: 一覧、詳細、フォーム、フィードバックなどのパターンは既存 component を組み合わせる方が差分を抑えやすい。

特に画面構造、入力、データ表示、フィードバックは既存 API を優先します。

| 目的 | 優先する構成例 |
| --- | --- |
| ページ骨格 | Shell、Navbar、Drawer、Container、Section、Grid |
| 操作 | Button、Link、Dropdown、Tabs、Steps |
| 入力 | Form、FormRow、Input、Select、Checkbox、Toggle、Range |
| データ表示 | Table、DataFrame、Card、Stats、Badge、Avatar、Progress |
| フィードバック | Alert、Toast、Modal、Drawer、Loading、Tooltip |

既存 component を使っていて props や helper が不足する場合は、いきなり新規 component を作る前に「小さな option の追加で済むか」「アプリ側 helper で合成できるか」を確認します。

## 2. `frontend/html` の低レベル HTML API を使う基準

`frontend/html` は、`Div`、`Span`、`Button`、`Attr`、`Class`、`Text` のような primitive を直接組み合わせるための低レベル API です。自由度が高い一方で、UI の意味付け、class の統一、アクセシビリティを実装者が管理する必要があります。

次のような場合に限定して使います。

- **既存 component の隙間を小さく埋める**: card 内の補助テキスト、table cell 内の小さな badge 群、空状態の説明文など。
- **画面固有で再利用予定がない**: 1 つの Page だけで使う装飾やレイアウト調整で、汎用 component にするほどではない。
- **既存 component に渡す child node を作る**: Card、Modal、Alert、Dropdown などの中身を組み立てる。
- **意味的に HTML primitive そのものが適切**: `dl` / `dt` / `dd`、`time`、`code`、小さな inline label など、専用 component を作るほどではない semantic HTML。
- **実験・移行中の局所実装**: API 化する前に画面内で最小限の markup を試し、必要性を確認する。

反対に、次の場合は低レベル API の直接利用を避けます。

- 同じ markup が複数箇所にコピーされ始めた。
- `aria-*`、keyboard 操作、focus 管理、label 関連付けが必要。
- DaisyUI の class 組み合わせが複雑で、誤用しやすい。
- 画面固有ではなく、アプリ全体や Marionette 本体に意味のある UI になっている。

低レベル API を使う場合でも、Page や Action に長い HTML 組み立てを直接置きすぎず、画面内 render helper や `internal/ui` の app-specific helper に切り出してください。

## 3. DaisyUI convenience component を使う基準

DaisyUI convenience component は、DaisyUI の class 体系を Marionette の Node API から扱いやすくするための helper です。DaisyUI の既知パターンに合う UI は、低レベル HTML で class を手書きするより convenience component を優先します。

次の条件に当てはまる場合に使います。

- **DaisyUI に公式の component 概念がある**: `alert`、`badge`、`card`、`modal`、`drawer`、`dropdown`、`tabs`、`stats`、`toast` など。
- **variant が DaisyUI の語彙で表せる**: `primary`、`secondary`、`success`、`warning`、`error`、`sm`、`lg` など。
- **見た目の統一が重要**: dashboard、admin shell、設定画面など、同じ theme / spacing / state color に揃えたい。
- **既存 API の props で十分**: convenience component の引数や variant helper で必要な見た目・状態を表せる。
- **DaisyUI 固有の構造を隠したい**: caller 側に `modal-box`、`toast-end`、`stats shadow` のような class 詳細を散らしたくない。

ただし、DaisyUI convenience component は「DaisyUI の component を使いやすくするもの」です。業務固有の概念、たとえば「請求ステータスカード」「ユーザー権限バッジ」「ジョブ実行サマリー」は、DaisyUI component を直接増やすのではなく、アプリ側 helper として名前を付ける方が自然です。

## 4. 新規コンポーネントを作る前の確認事項

新規 component を作る前に、次の項目を確認します。

- **既存 API の確認**: `frontend`、`frontend/daisyui`、Form API、Data display API、Overlay / Feedback API に近い component がないか。
- **合成で足りるか**: 既存 component と child node の組み合わせ、または app-specific helper で十分ではないか。
- **責務の置き場所**: Marionette 本体の汎用 component か、特定アプリの `internal/ui` helper か、画面内 helper か。
- **アクセシビリティ**: label、role、aria、focus、keyboard 操作、loading / disabled / error 状態をどう扱うか。
- **状態と API の形**: 表示 state、validation error、loading、empty、selected、active などを props と child node のどちらで受けるか。
- **DaisyUI 依存の有無**: DaisyUI primitive なのか、DaisyUI を内部実装として使うだけのアプリ固有 UI なのか。
- **互換性要件**: 既存 root API や template-backed component と HTML 出力互換性が必要か。
- **テスト方針**: golden test、render test、class name test、アクセシビリティに関係する属性の確認をどこまで行うか。
- **命名**: HTML タグ名や class 名ではなく、ユーザーや業務から見た意味で名前を付けられるか。

判断に迷う場合は、最初から Marionette 本体へ追加せず、アプリ内 helper として始めます。複数アプリや複数画面で同じ意味の UI が安定してから built-in component 化を検討してください。

## 5. よく使う画面パターン別の推奨構成

### 一覧画面

一覧画面は、検索条件、テーブル、ページング、行アクション、空状態を分けて考えます。

推奨構成:

- Page: query state を復元し、一覧データと件数を読み込む。
- Layout: Container / Section / Card で一覧領域を作る。
- Filter: Form、Input、Select、Button、Tabs などで検索・絞り込みを表現する。
- Data: Table または DataFrame を使い、行内の状態は Badge / Status / Progress で表す。
- Action: 行ごとの詳細リンク、編集ボタン、削除ボタンは Button / Link / Dropdown に寄せる。
- Empty / Loading: EmptyState、Alert、Loading、Skeleton 相当の既存表現を使う。

低レベル HTML は、table cell 内の小さな補助情報や、検索フォーム内の短い説明などに限定します。複数一覧で同じ filter block が出る場合は、app-specific helper に切り出します。

### 詳細画面

詳細画面は、タイトル、主要アクション、概要、属性リスト、関連データを分けて配置します。

推奨構成:

- Header: Page title、breadcrumb、primary action を Section / Button / Link でまとめる。
- Summary: Card、Stats、Badge、Avatar、Progress で重要情報を上部に出す。
- Attributes: semantic HTML の `dl` / `dt` / `dd` または app-specific detail list helper を使う。
- Related data: Table、Timeline、Tabs、Accordion で履歴や関連リソースを整理する。
- Feedback: 更新結果は Alert / Toast、危険操作は Modal / Confirm 系 UI に寄せる。

詳細画面では、業務固有の「表示名 + 値 + 状態」の組み合わせが繰り返されやすいため、低レベル HTML のコピーが増えたら `internal/ui` に detail row / metadata panel helper を作ります。

### 入力フォーム

入力フォームは、入力値、validation error、submit 結果、再描画の境界を明確にします。

推奨構成:

- Page: 初期値と選択肢を読み込んで form state を渡す。
- Form: Form / FormRow / Fieldset / Label / Input / Select / Checkbox / Toggle など既存入力 component を使う。
- Validation: field error は該当 FormRow の近くに、page-level error は Alert として表示する。
- Actions: submit、cancel、delete などは Button group や Card footer にまとめる。
- Partial update: htmx で一部だけ更新する場合も、form state から同じ render helper を呼び直す。

フォームでは label と input の関連付け、required / disabled / error 表現、入力値の再表示が重要です。`frontend/html` で input を直接作るのは、既存入力 component では表せない特殊入力に限定します。

### ダッシュボード

ダッシュボードは、指標、チャート、直近イベント、操作導線をカード単位で構成します。

推奨構成:

- Layout: responsive Grid、Section、Card で情報密度を調整する。
- Metrics: Stats、Card、Badge、Progress、Chart を使う。
- Trends: Chart、Timeline、Table で時系列や履歴を表示する。
- Actions: よく使う操作は Button / Link / Dropdown としてカードの header / footer に置く。
- State: 期間、対象組織、フィルターは URL query または共有 state として扱う。

DashWind 風の admin dashboard では、まず DashWind API や既存 dashboard 向け component を確認してください。見た目のためだけに画面ごとに独自 grid / card class を増やすと、テーマ変更や responsive 調整が難しくなります。

### モーダル / トースト / アラート

モーダル、トースト、アラートは、ユーザーへの割り込み度合いで使い分けます。

推奨構成:

- Alert: ページ内に残すべき注意、validation summary、保存失敗、権限不足など。
- Toast: 保存成功、軽い通知、短時間で消えてよい feedback。
- Modal: 破壊的操作の確認、追加情報が必要な操作、現在の文脈を保ったまま行う短い入力。
- Drawer: 画面横から補助情報や詳細パネルを出したい場合。
- Loading: 非同期処理中の状態を Button / Card / Table の近くに出す。

これらはアクセシビリティや focus、閉じる操作、表示位置の一貫性が重要なため、DaisyUI convenience component または既存 Overlay / Feedback API を優先します。`frontend/html` で直接実装する場合は、既存 component で表せない局所的な中身だけに留めます。

## 実装チェックリスト

- 既存 `frontend` component と DaisyUI convenience component を先に確認した。
- `frontend/html` は画面固有の小さな補助 markup に限定した。
- 複数画面で繰り返す UI は app-specific helper に切り出した。
- Marionette 本体へ追加する component は、責務の置き場所と互換性要件を確認した。
- label、aria、focus、disabled、loading、error、empty state を確認した。
- 一覧、詳細、フォーム、ダッシュボード、フィードバックのパターンに沿って構成した。
- 詳細判断が必要な場合は [UI Component Guidelines](../../ui-component-guidelines.md) と [frontend/ARCHITECTURE.md](../../../frontend/ARCHITECTURE.md) を確認した。
