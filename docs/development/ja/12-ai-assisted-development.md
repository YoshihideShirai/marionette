# AI Assisted Development ガイド

[English](../12-ai-assisted-development.md) | 日本語

このドキュメントは、Marionette アプリケーションの開発を AI に依頼するときに、依頼内容を実装可能な粒度へそろえるためのテンプレート集です。AI には「何を作るか」だけでなく、Page、Action、State、partial update、エラー表示、確認すべきリスクを明示して渡します。

AI に依頼する前に、関連する開発ガイドも確認してください。

- [Routing / Pages / Actions ガイド](02-routing-pages-actions.md): URL、Page、Action、partial update の責務分担。
- [Forms / Validation ガイド](05-forms-validation.md): form state、バリデーション、エラー再表示。
- [データテーブル / チャート設計ガイド](06-data-tables-charts.md): テーブル、チャート、共有 filter state。
- [Security / Authorization ガイド](10-security-authz.md): 認可、危険操作、監査ログ。

## 基本方針

AI への依頼は、次の情報が 1 つの文章に混ざらないように分けて書きます。

- 目的: ユーザーが何を達成する画面・操作か。
- 対象: URL、Page 名、Action 名、関連ファイル。
- 状態: URL query、form state、session、DB のどこに何を保存するか。
- 更新: full page reload、redirect、partial update、toast / flash のどれを期待するか。
- 失敗: validation error、権限 error、保存 error、外部 API error をどう表示するか。
- 制約: 認可、監査ログ、データ破壊操作、本番設定など、人間が最終確認する項目。

依頼文では「いい感じに」「既存に合わせて」だけにせず、既存ファイル名、関数名、表示したい状態、テスト観点を具体的に書きます。

## 1. 新しい画面を追加するときの依頼テンプレート

```text
Marionette アプリに新しい画面を追加してください。

目的:
- <誰が何を確認・操作する画面か>

対象 URL:
- GET <例: /admin/projects/:projectID/members>

Page 名:
- <例: ProjectMembersPage>

表示内容:
- <例: プロジェクト名、メンバー一覧、ロール、招待状態、最終ログイン日時>
- <空状態、ローディング状態、権限がない場合の表示も書く>

State の保存先:
- URL query: <例: q, role, page, sort>
- Page state struct: <例: ProjectMembersPageState>
- session / DB: <画面表示に必要な永続化データ>

Action 名:
- <この画面から呼ぶ Action があれば列挙。例: InviteMember, RemoveMember>

期待する partial update:
- <例: filter 変更時は #members-table だけ更新する>
- <例: 招待成功時は #members-table と toast を更新する>

エラー時の表示:
- <例: 一覧取得失敗は Alert、権限なしは 403 相当のメッセージ、招待失敗はフォーム内 field error>

実装方針:
- 既存の frontend コンポーネントを優先してください。
- Page / Action / State / UI helper の責務を分けてください。
- 必要なテストを追加してください。

人間が確認する項目:
- 認可、監査ログ、データ破壊操作、本番設定は実装後にレビューします。
```

### 補足

新しい画面では、URL と Page state の責務を先に固定すると実装が安定します。テーブルやフォームが含まれる場合も、まず Page が読み込むデータと state の形を決め、その後に Action と partial update を追加します。

## 2. 既存 Action を変更するときの依頼テンプレート

```text
既存 Action を変更してください。

対象 URL:
- <例: POST /admin/projects/:projectID/members/:memberID/role>

Page 名:
- <例: ProjectMembersPage>

Action 名:
- <例: UpdateMemberRole>

現在の挙動:
- <例: ロール変更後に一覧ページへ redirect している>

変更したい挙動:
- <例: ロール変更後に #members-table の対象行だけ partial update し、成功 toast を表示する>

入力:
- path parameter: <例: projectID, memberID>
- form value / JSON / query: <例: role>

State の保存先:
- URL query: <維持する条件。例: q, page, sort>
- Action state / form state: <例: MemberRoleForm>
- DB: <更新する永続化データ>

期待する partial update:
- 成功時: <例: #member-row-<memberID> を差し替える>
- validation error 時: <例: 同じ row 内に field error を表示する>
- service error 時: <例: #members-feedback に Alert を表示する>

エラー時の表示:
- validation error: <ユーザーが修正できる field error>
- 認可 error: <安全なメッセージ。詳細はログへ>
- 保存 error: <再試行できるメッセージ>

維持してほしいこと:
- <既存の公開 API、既存テスト、既存 UI の互換性など>

テスト観点:
- 正常系、validation error、認可 error、保存 error、partial update の target を確認してください。
```

### 補足

既存 Action の変更では、現在の redirect / partial update / flash の流れを明示します。AI が Action の戻り方を勝手に変えると、画面遷移や htmx target が壊れやすいためです。

## 3. フォームとバリデーションを追加するときの依頼テンプレート

```text
フォームとサーバー側バリデーションを追加してください。

対象 URL:
- GET <例: /admin/projects/new>
- POST <例: /admin/projects>

Page 名:
- <例: NewProjectPage>

Action 名:
- <例: CreateProject>

フォーム項目:
- <例: name: 必須、1-80 文字>
- <例: slug: 必須、英数字とハイフンのみ、重複不可>
- <例: visibility: public/private のいずれか>

State の保存先:
- form state struct: <例: ProjectForm>
- field errors: <例: ProjectForm.Errors>
- page-level error: <例: ProjectForm.PageError>
- DB: <保存先の model / repository / service>

バリデーション:
- 入力形式: <form state または validator で確認する内容>
- 業務ルール: <service / domain で確認する内容>
- 重複確認: <必要なら repository / service で確認する内容>

期待する partial update:
- 初回表示: <full page または form fragment>
- validation error 時: <例: #project-form を入力値とエラー付きで再描画する>
- 成功時: <例: 詳細ページへ redirect、または一覧に row を追加して toast>

エラー時の表示:
- field error: <各 FormRow に表示する文言>
- page-level error: <フォーム上部 Alert に表示する文言>
- unexpected error: <安全な汎用メッセージ。詳細はログへ>

UI 方針:
- 既存の FormRow、Input、Select、Textarea、Button、Alert を優先してください。
- 入力値は validation error 後も保持してください。

テスト観点:
- 初期表示、正常 submit、各 field error、業務エラー、成功時の戻り方を確認してください。
```

### 補足

フォーム依頼では、field error と page-level error を分けて指定します。入力形式のエラーはユーザーがその場で修正できるようにし、DB や外部 API 由来の失敗は page-level error や toast に寄せます。

## 4. テーブル・チャート連動を追加するときの依頼テンプレート

```text
テーブルとチャートが同じ filter state で連動する機能を追加してください。

目的:
- <例: 売上推移チャートをクリックすると、下の注文テーブルを同じ期間・カテゴリで絞り込む>

対象 URL:
- GET <例: /admin/reports/sales>

Page 名:
- <例: SalesReportPage>

Action 名:
- <例: UpdateSalesReportFilters, SelectSalesChartBucket>

共有 State の保存先:
- URL query: <例: from, to, granularity, category, status, page, sort>
- Page state struct: <例: SalesReportState>
- table state: <例: SalesTableState>
- chart state: <例: SalesChartState>
- DB / cache: <集計結果や一覧取得の保存・取得方針>

連動ルール:
- filter 変更時: <例: chart と table を同じ条件で再取得する>
- chart click 時: <例: clicked bucket を URL query に反映し、table を 1 ページ目へ戻す>
- table page / sort 変更時: <例: chart 条件は維持し、table だけ更新する>

期待する partial update:
- filter 変更: <例: #sales-chart と #sales-table を更新する>
- chart click: <例: #sales-table、#active-filters、URL query を更新する>
- table pagination: <例: #sales-table のみ更新する>

エラー時の表示:
- chart 取得失敗: <chart 領域に Alert または EmptyState>
- table 取得失敗: <table 領域に Alert または EmptyState>
- filter 不正: <filter form に field error、または安全な default に戻す>

UI 方針:
- 共有 state から chart と table を描画してください。
- filter、active filter 表示、empty state、loading state を用意してください。

テスト観点:
- 共有 state の parse / serialize、chart click 後の table 条件、pagination 時の chart 条件維持を確認してください。
```

### 補足

テーブル・チャート連動では、表示データより先に filter state を設計します。URL query に残す条件と、一時的に画面内だけで使う状態を分けると、reload、共有リンク、partial update の挙動が説明しやすくなります。

## 5. AI に必ず渡すべき情報

AI に実装を依頼するときは、少なくとも次を渡します。

| 情報 | 書く内容 | 例 |
| --- | --- | --- |
| 対象 URL | HTTP method、path、path parameter、query parameter。 | `GET /admin/projects/:projectID/members?q=&page=` |
| Page 名 | 追加・変更する Page 関数や画面 state の名前。 | `ProjectMembersPage`, `ProjectMembersPageState` |
| Action 名 | submit、クリック、filter 変更などを処理する Action 名。 | `InviteMember`, `UpdateMemberRole` |
| State の保存先 | URL query、form state、Page state、session、DB、cache のどこに置くか。 | `q` と `page` は URL query、入力値は `InviteMemberForm` |
| 期待する partial update | 成功・失敗・filter 変更時に更新する DOM target と範囲。 | 成功時は `#members-table` と toast、失敗時は `#invite-member-form` |
| エラー時の表示 | field error、page-level Alert、toast、EmptyState、403 / 404 の出し分け。 | validation は FormRow、認可失敗は安全な Alert、詳細はログ |

追加で渡すと精度が上がる情報は次のとおりです。

- 参考にしてほしい既存ファイル、関数、テスト。
- 使ってほしい frontend コンポーネント。
- 変更してはいけない public API、既存 URL、既存 CSS class、既存テスト。
- 正常系だけでなく、空状態、権限なし、保存失敗、外部 API 失敗の期待表示。
- 人間がレビューする前提のリスク項目。

## 6. AI に任せず人間が確認すべき項目

AI は実装のたたき台を作れますが、次の項目は必ず人間が確認します。

### 認可

- Page 表示前の認可チェックが、対象 tenant、organization、project、resource に対して正しいか。
- UI でボタンを隠すだけでなく、Action 内でも同じ操作権限を確認しているか。
- 他ユーザーや他 tenant の ID を直接指定された場合に拒否できるか。
- 403 / 404 の使い分けで、存在してはいけない情報を漏らしていないか。

### 監査ログ

- 作成、更新、削除、承認、権限変更、再実行などの重要操作が記録されるか。
- 成功だけでなく、拒否、失敗、validation error、外部連携失敗を記録すべきか判断しているか。
- actor、action、target、before / after、request ID、IP、user agent など必要な項目が含まれるか。
- ログに secret、token、個人情報、不要な本文が残らないか。

### データ破壊操作

- 削除、取り消し、上書き、bulk update、再同期などの破壊的操作に確認 UI があるか。
- soft delete / hard delete、復元可否、cascade、関連データへの影響が明確か。
- 二重送信、リトライ、部分失敗、途中中断時の扱いが決まっているか。
- transaction、idempotency key、楽観ロックなどが必要か判断しているか。

### 本番設定

- 本番向けの環境変数、secret、外部 API endpoint、callback URL、CORS、cookie 設定が正しいか。
- debug mode、サンプル認証、開発用 seed、モック API が本番で有効にならないか。
- rate limit、timeout、retry、queue、cache TTL、監視・アラートの値が運用に合っているか。
- migration、feature flag、rollback 手順、データバックフィル手順が用意されているか。

## 依頼前チェックリスト

- 対象 URL、Page 名、Action 名を書いた。
- State の保存先を書いた。
- 期待する partial update の DOM target と更新範囲を書いた。
- エラー時の表示を validation、認可、保存失敗、予期しない失敗に分けて書いた。
- 既存ファイル、参考実装、変更してはいけない互換性を書いた。
- 認可、監査ログ、データ破壊操作、本番設定は人間が確認すると明記した。
