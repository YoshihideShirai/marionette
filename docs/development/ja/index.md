# Marionette 開発ガイド

[English](../index.md) | 日本語

このガイドは、Marionette で管理画面・社内ツールを作るときの標準的な進め方と、参照すべき既存ドキュメントへの導線をまとめた入口ページです。最初に全体像をつかみ、実装フェーズごとに必要な詳細ドキュメントへ移動できるようにしています。

## ドキュメントの読み方

### まず読む

- [README.ja.md](../../../README.ja.md): Marionette の目的、基本思想、デモ、ドキュメント全体への入口を確認します。
- [プロジェクト構成ガイド](01-project-structure.md): Page、Action、State、UI helper、asset のユーザーアプリ側配置を決めます。
- [Routing / Pages / Actions ガイド](02-routing-pages-actions.md): URL、業務画面、partial update、Action handler の分割基準を決めます。
- [UI コンポーネント選定ガイド](04-ui-components.md): 既存コンポーネント、低レベル HTML、DaisyUI helper、新規コンポーネントの使い分けを決めます。
- [State Management ガイド](../../state-management.ja.md): ページ、アクション、状態を Go 側に集約する基本方針を確認します。
- [API ドキュメント（日本語版）](../../api/ja/): `backend` / `frontend` / `html` など主要 API の入口です。

### 必要になったら読む

- [UI Component Guidelines](../../ui-component-guidelines.md): コンポーネントの選び方、アクセシビリティ、DaisyUI との付き合い方を確認します。
- [frontend/ARCHITECTURE.md](../../../frontend/ARCHITECTURE.md): `frontend` パッケージの責務境界、テンプレートと Go 実装の使い分け、互換性方針を確認します。
- [フォーム API](../../api/04-form-apis.md): フォーム行、入力、エラー表示、選択 UI を組み立てるときに参照します。
- [データ表示コンポーネント API](../../api/ja/05-component-apis-data-display.md): テーブル、統計、アバター、プログレスなどの表示部品を確認します。
- [Overlay / Feedback API](../../api/ja/05-component-apis-overlay-feedback.md): モーダル、ドロワー、トースト、アラートなどの通知・フィードバック部品を確認します。
- [DashWind API](../../api/ja/08-dashwind.md): DashWind ベースのダッシュボードや管理画面レイアウトを使うときに参照します。

### 実装前に確認する

- [frontend/ARCHITECTURE.md](../../../frontend/ARCHITECTURE.md): 新しい UI ヘルパーやコンポーネントを追加する前に、どのパッケージが責務を持つか確認します。
- [UI Component Guidelines](../../ui-component-guidelines.md): 新規 UI のアクセシビリティ、見た目、状態表現、コンポーネント選定が既存方針とずれていないか確認します。
- [State Management ガイド](../../state-management.ja.md): URL、セッション、フォーム入力、一時的な UI 状態をどこに置くか確認します。
- [API ドキュメント（日本語版）](../../api/ja/): 既存 API で実現できるかを確認し、不要なラッパーや重複実装を避けます。

## 標準的な開発フロー

### 1. プロジェクト構成を決める

最初に、アプリの入口、ページ、状態、アクション、表示コンポーネントの置き場所を分けます。小さなアプリでは 1 パッケージにまとめても構いませんが、画面数が増える場合は次のように責務で分けると見通しがよくなります。

```text
cmd/<app>/main.go        # アプリ起動とルーティング登録
internal/<app>/pages/    # Page 関数、画面単位の構成
internal/<app>/actions/  # POST/イベント処理、状態更新、リダイレクト
internal/<app>/state/    # クエリ、フィルタ、フォーム値、セッション由来の状態
internal/<app>/ui/       # アプリ固有の小さな UI 合成ヘルパー
```

Marionette 本体のコンポーネントを変更する場合は、`frontend` 配下の責務境界が重要です。既存コンポーネントの薄いエイリアス、DaisyUI 固有の実装、低レベル HTML ノード生成を混ぜないように、実装前に [frontend/ARCHITECTURE.md](../../../frontend/ARCHITECTURE.md) を確認してください。

### 2. Page / Action を設計する

詳しい分割基準は [Routing / Pages / Actions ガイド](02-routing-pages-actions.md) を参照してください。Marionette では、画面の読み取りは Page、ユーザー操作による変更は Action として考えると整理しやすくなります。

- Page はリクエストから状態を復元し、データ取得を行い、`frontend` コンポーネントで HTML を返します。
- Action はフォーム送信、ボタンクリック、フィルタ変更などを受け取り、検証、状態更新、永続化、部分更新レスポンスを担当します。
- Page と Action の間で共有する入力値やフィルタ条件は、明示的な state 型にまとめます。

ページ単位で「何を表示するか」と「どの操作で何が変わるか」を先に書き出すと、htmx による部分更新の対象も決めやすくなります。API の基本は [API ドキュメント（日本語版）](../../api/ja/) と [Context API](../../api/ja/02-context.md) を参照してください。

### 3. State 管理を決める

状態は、寿命と共有範囲に応じて置き場所を分けます。

- URL クエリ: 検索条件、ページ番号、ソート順など、共有・再読み込みしたい状態。
- フォーム値: 入力中または送信された値。バリデーションエラー時に再表示します。
- サーバー側セッションまたは永続化層: 認証、ユーザー設定、長く保持する業務データ。
- 一時的な UI 状態: モーダル開閉やトースト表示など、画面上の短命な状態。

詳細な考え方は [State Management ガイド](../../state-management.ja.md) を基準にします。テーブルとチャートのように複数ウィジェットが同じ絞り込み条件を共有する場合は、共通の state 型に寄せると変更範囲を抑えられます。

### 4. UI コンポーネントを選定する

UI は、まず既存の `frontend` API と DaisyUI ベースのコンポーネントで組み立てます。低レベル HTML を直接書く前に、既存のレイアウト、入力、データ表示、フィードバック部品で表現できるかを確認してください。

- 画面骨格: Shell、Navbar、Drawer、Container、Section、Grid など。
- 操作: Button、Link、Dropdown、Tabs、Steps など。
- 情報表示: Card、Alert、Badge、Stats、Timeline など。
- 業務 UI: Table、DataFrame、Chart、FormRow、Input 系コンポーネントなど。

画面パターン別の実践的な選定基準は [UI コンポーネント選定ガイド](04-ui-components.md) を参照します。アクセシビリティ、見た目の一貫性は [UI Component Guidelines](../../ui-component-guidelines.md) を参照します。新しいコンポーネントを本体に追加する場合は、実装先を [frontend/ARCHITECTURE.md](../../../frontend/ARCHITECTURE.md) で確認します。

### 5. フォーム・バリデーションを組み立てる

フォームは、入力値、検証結果、再表示するエラーメッセージを 1 つの流れとして設計します。

1. Page で初期値を state に詰める。
2. FormRow と入力コンポーネントでラベル、説明、必須表示、エラー領域を揃える。
3. Action で送信値を検証する。
4. エラーがあれば同じ state に値とエラーを戻して再描画する。
5. 成功時は保存、通知、リダイレクト、または部分更新を行う。

フォーム部品の仕様は [フォーム API](../../api/04-form-apis.md) を参照してください。日本語版 API に未翻訳の項目がある場合も、[API ドキュメント（日本語版）](../../api/ja/) を入口に関連セクションを確認します。

### 6. テーブル・チャートを設計する

テーブルやチャートは、表示データそのものより先に「絞り込み、集計粒度、ソート、ページング」を state として設計します。

- テーブルは、列定義、行アクション、空状態、ローディング表示、ページングをまとめて考えます。
- チャートは、集計単位、期間、凡例、クリック時の絞り込み動作を明確にします。
- テーブルとチャートが同じ条件で連動する場合は、共通 state から両方を描画します。

データ表示部品は [データ表示コンポーネント API](../../api/ja/05-component-apis-data-display.md) を確認してください。DashWind の管理画面パターンを使う場合は [DashWind API](../../api/ja/08-dashwind.md) も参照します。

### 7. エラー表示・通知を決める

エラーと通知は、ユーザーが次に何をすればよいか分かる粒度で出します。

- フィールド単位の入力エラーは FormRow / FieldError に寄せます。
- 画面全体の失敗は Alert や EmptyState で説明します。
- 保存成功や軽い失敗は Toast / Flash などの短い通知にします。
- 復旧不能な失敗は、ログに詳細を残し、画面には安全なメッセージだけを表示します。

表示部品は [Overlay / Feedback API](../../api/ja/05-component-apis-overlay-feedback.md) と [Flash API](../../api/06-flash-apis.md) を参照してください。通知の見た目やアクセシビリティは [UI Component Guidelines](../../ui-component-guidelines.md) も確認します。

### 8. テスト・デバッグを用意する

Marionette アプリでは、画面生成と Action の状態更新を Go のテストで確認しやすくしておくと、UI の変更が安全になります。

- Page 関数は、代表的な state から期待する HTML 断片が出ることをテストします。
- Action は、正常系、バリデーションエラー、権限エラー、永続化エラーを分けてテストします。
- コンポーネント追加時は、既存のレンダリングテストや golden テストの方針に合わせます。
- htmx の部分更新は、対象要素 ID、レスポンス範囲、エラー時の再描画を確認します。

本体の UI 実装を変更する場合は [frontend/ARCHITECTURE.md](../../../frontend/ARCHITECTURE.md) の互換性方針を守り、必要なレンダリングテストを追加します。アプリ側では、状態型と Action を小さく保つほど、失敗時の再現とデバッグが簡単になります。

## 実装前チェックリスト

- Page / Action / State の責務が分かれている。
- URL に残す状態、フォームに残す状態、セッションや DB に置く状態を分類した。
- 既存の `frontend` コンポーネントと API で実現できるか確認した。
- フォームのエラー表示、成功通知、画面全体エラーの出し分けを決めた。
- テーブル・チャートのフィルタ条件を共有 state として扱うか判断した。
- 追加する UI が [UI Component Guidelines](../../ui-component-guidelines.md) と [frontend/ARCHITECTURE.md](../../../frontend/ARCHITECTURE.md) に沿っている。
- テスト対象の Page、Action、コンポーネント出力を決めた。
