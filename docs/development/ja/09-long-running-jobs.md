# 長時間ジョブ設計ガイド

[English](../09-long-running-jobs.md) | 日本語

このドキュメントは、データ集計、CSV インポート、外部 API 同期、レポート生成のように、HTTP request の寿命だけでは完了しない処理を Marionette アプリケーションで扱うための設計ガイドです。README.ja.md の **Heavy Job Template（データアプリ向け）** を、実務で使える job model、UI pattern、timeout / retry / cache TTL、永続化判断まで広げて説明します。

長時間ジョブは、**開始 request、進捗確認、結果表示、失敗時の再実行を分ける**ことを標準にします。ユーザーには「今どの状態か」「待つべきか、直すべきか、再実行すべきか」「結果をどこで確認できるか」を明示し、サーバー側には調査できる execution ID と状態遷移を残します。

## 基本方針

- Job は通常の Page / Action と同じく Go 側で状態を管理し、UI は job state を描画するだけにします。
- 実行開始 Action は job を作成または enqueue し、すぐに進捗ページまたは元画面の job panel へ戻します。
- 重い処理そのものは request handler の中で同期的に待ち続けず、worker、goroutine、queue、外部 job runner などに分離します。
- htmx polling や手動 refresh で進捗領域だけを更新し、ページ全体の再読み込みを避けます。
- 成功、失敗、timeout、cancel、retry は state transition として明示的に扱います。
- ユーザー向けメッセージには内部エラーをそのまま出さず、ログには execution ID、request ID、actor、入力 hash、外部 request ID を残します。

## 1. Job モデルの標準フィールド

Job は UI 表示と運用調査の両方に使うため、最低限の標準フィールドを揃えます。アプリ固有の入力値や結果の型は別 struct に分けても構いませんが、一覧、詳細、retry、監査ログで共通して参照する項目は Job model に残します。

```go
type JobState string

const (
    JobStateQueued    JobState = "queued"
    JobStateRunning   JobState = "running"
    JobStateSucceeded JobState = "succeeded"
    JobStateFailed    JobState = "failed"
    JobStateTimedOut  JobState = "timed_out"
    JobStateCanceled  JobState = "canceled"
)

type Job struct {
    ID             string
    ExecutionID    string
    State          JobState
    Progress       JobProgress
    ResultRef      *JobResultRef
    Error          *JobError
    InputHash      string
    Attempt        int
    MaxAttempts    int
    CreatedAt      time.Time
    StartedAt      *time.Time
    FinishedAt     *time.Time
    ExpiresAt      *time.Time
    RequestedBy    string
}

type JobProgress struct {
    Current int
    Total   int
    Percent int
    Message string
}

type JobResultRef struct {
    Kind string // "table", "csv", "report", "object", "url" など
    URI  string
}

type JobError struct {
    Code       string
    Message    string
    Retryable  bool
    DetailRef  string
}
```

### execution ID

`execution ID` は、ユーザー問い合わせ、ログ検索、外部 API request、worker trace をつなぐ相関 ID です。UI に表示する job ID と同じでもよいですが、retry で新しい実行を作る場合は次のように分けると追跡しやすくなります。

- `Job.ID`: ユーザー操作としての job record。retry 後も同じ業務操作として扱う場合に使います。
- `ExecutionID`: 実際の 1 回の試行。retry のたびに新しい値にします。
- `ParentJobID` または `RetryOf`: 再実行元をたどりたい場合に追加します。

UI には短い `execution ID` または copy しやすい `job ID` を表示し、support へ問い合わせるときに使えるようにします。

### state

`state` は画面表示、操作可否、worker の再開判断の基準です。文字列を自由に増やすのではなく、許可された状態遷移を決めます。

| state | 意味 | 主な UI |
| --- | --- | --- |
| `queued` | 受け付け済みで実行待ち | toast / progress の待機表示 |
| `running` | 実行中 | progress、実行中メッセージ、必要なら cancel |
| `succeeded` | 成功 | result link、summary、success toast |
| `failed` | 失敗 | alert、error summary、retry button |
| `timed_out` | 実行予算超過 | timeout 用 alert、retry / 条件調整 |
| `canceled` | ユーザーまたは管理者が中止 | neutral alert、再実行導線 |

`running` から `succeeded` / `failed` / `timed_out` / `canceled` に進むようにし、完了済み state から勝手に `running` に戻さないでください。retry は既存 job の state を巻き戻すのではなく、新しい execution を作るほうが監査しやすくなります。

### progress

`progress` は、ユーザーが待つ判断をするための情報です。数値化できる処理では `Current` / `Total` / `Percent` を持ち、数値化できない処理では `Message` と state だけでも構いません。

- `Percent` は 0〜100 に clamp します。
- `Total` が不明な場合は indeterminate progress として表示します。
- `Message` は「12,000 件中 4,200 件を集計中」のように現在の段階を表します。
- 細かすぎる更新は DB / cache 書き込みを増やすため、数百 ms〜数秒単位または段階ごとにまとめます。

### result reference

`result reference` は結果そのものではなく、結果を取りに行くための参照です。大きな集計結果やファイルを Job record に直接詰め込むと、一覧表示や retry 管理が重くなります。

- 小さな summary: Job record の result summary に保存してもよい。
- テーブル表示: query state、snapshot ID、materialized table ID を参照する。
- ファイル: object storage path、signed URL の元になる object key、download route を参照する。
- レポート: report ID、dashboard URL、外部 BI URL を参照する。

参照先が消える可能性がある場合は `ExpiresAt` や result TTL を明示し、期限切れ時の再生成導線を UI に出します。

### error

`error` はユーザー向け表示と運用調査の境界を分けます。

- `Code`: `external_timeout`、`invalid_input`、`rate_limited` のような分類。
- `Message`: ユーザーに見せてもよい短い説明。
- `Retryable`: 同じ条件で再実行してよいか。
- `DetailRef`: log trace、error report、外部 request ID など内部調査用の参照。

`err.Error()` や stack trace を UI にそのまま表示しないでください。画面には「外部 API の応答が遅いため集計できませんでした。時間をおいて再実行してください。」のように、状態と次の行動を出します。

## 2. 実行開始、進捗表示、成功、失敗、再実行の UI パターン

長時間ジョブの画面は、1 つの submit 結果ではなく lifecycle として設計します。

### 実行開始

実行開始 Action では、入力を検証し、job を作成し、二重送信を防ぎ、進捗確認先へ誘導します。

- 実行ボタンは submit 中に loading / disabled にします。
- 入力条件から `InputHash` を作り、cache hit または重複実行の判定に使います。
- enqueue に成功したら toast または flash で「処理を開始しました」と伝えます。
- 進捗ページがある場合は `/admin/jobs/{id}` へ redirect します。
- 同じ画面内で進捗 panel を出す場合は、job ID を含む partial を返します。

```go
func (h *AnalyticsActions) RunAggregation(ctx *backend.Context) error {
    input, formErr := state.AggregationInputFromRequest(ctx)
    if formErr != nil {
        return ctx.HTML(ui.AggregationForm(input, formErr))
    }

    job, err := h.Jobs.Enqueue(ctx.Request.Context(), input)
    if err != nil {
        h.Logger.Error("enqueue aggregation failed", "err", err, "input_hash", input.Hash())
        ctx.FlashError("集計を開始できませんでした。条件を確認してもう一度お試しください。")
        return ctx.Redirect("/admin/analytics")
    }

    ctx.FlashInfo("集計を開始しました")
    return ctx.Redirect("/admin/analytics/jobs/" + job.ID)
}
```

### 進捗表示

進捗表示は、job state を読む Page または partial endpoint に集約します。htmx polling を使う場合も、polling interval は短くしすぎず、完了後は停止します。

- `queued` / `running`: progress、現在の step、execution ID、開始時刻を表示します。
- `succeeded`: polling を止め、結果リンクや summary に切り替えます。
- `failed` / `timed_out`: polling を止め、alert と retry 導線を表示します。
- 画面を離れても戻れるように job detail URL を用意します。

```go
func AggregationJobPanel(job service.Job) frontend.Node {
    switch job.State {
    case service.JobStateQueued, service.JobStateRunning:
        return frontend.Section(frontend.SectionProps{},
            frontend.Progress(frontend.ProgressProps{Value: job.Progress.Percent, Max: 100}),
            frontend.P(frontend.Text(job.Progress.Message)),
            frontend.P(frontend.Text("Execution ID: "+job.ExecutionID)),
        )
    case service.JobStateSucceeded:
        return ui.AggregationResultCard(job.ResultRef)
    default:
        return ui.AggregationFailureCard(job)
    }
}
```

### 成功

成功時は、短い成功通知だけで終わらせず、結果を確認する場所を残します。

- result summary を card / stats / table で表示する。
- 結果ファイルや詳細ページへの link を置く。
- 入力条件、実行時間、生成日時、cache hit かどうかを表示する。
- 同じ条件で再生成する必要がある画面では「再実行」または「キャッシュを無視して再実行」を分けます。

成功 toast は、進捗ページ以外にいるユーザーへ「完了しました」と知らせる用途に向いています。進捗ページ上では、result card や alert など画面内に残る表示を主にします。

### 失敗

失敗時は、原因分類と次の行動を画面内に残します。

- 入力が不正: フォームに戻して修正箇所を示す。
- 外部 API timeout / rate limit: 時間をおいて retry できることを示す。
- 権限不足: retry ではなく必要な権限や問い合わせ先を示す。
- 部分成功: 成功件数、失敗件数、失敗 row download などを表示する。
- 問い合わせが必要: job ID / execution ID / request ID を copy できるようにする。

失敗理由が長い場合は、toast ではなく alert、details、log summary への link を使います。

### 再実行

retry は「同じ job record を上書きする」よりも、「元 job から入力と context を引き継いだ新しい execution を作る」設計を推奨します。

- `Retryable` が false の error では retry button を出さない。
- retry 前に現在の入力条件を表示し、必要なら編集できるようにする。
- retry は idempotency key を使って二重クリックを防ぐ。
- retry 後は新しい execution ID を表示し、元 execution への link を残す。
- 自動 retry と手動 retry は別々に記録する。

```go
func (h *AnalyticsActions) RetryAggregation(ctx *backend.Context) error {
    original, err := h.Jobs.Find(ctx.Request.Context(), ctx.Param("id"))
    if err != nil || !original.CanRetry() {
        ctx.FlashError("この集計は再実行できません。条件を確認して新しく実行してください。")
        return ctx.Redirect("/admin/analytics/jobs/" + ctx.Param("id"))
    }

    retried, err := h.Jobs.Retry(ctx.Request.Context(), original.ID)
    if err != nil {
        ctx.FlashError("再実行を開始できませんでした。時間をおいてお試しください。")
        return ctx.Redirect("/admin/analytics/jobs/" + original.ID)
    }

    ctx.FlashInfo("再実行を開始しました")
    return ctx.Redirect("/admin/analytics/jobs/" + retried.ID)
}
```

## 3. progress / toast / empty_state の使い分け

Heavy Job Template では `progress`、`toast`、`empty_state` を組み合わせますが、それぞれの役割は分けます。

| UI | 使う場面 | 避ける場面 |
| --- | --- | --- |
| `progress` | 実行中または待機中で、状態が変わり続けることを示す | 初回実行前、完了後の結果説明、読ませたい失敗理由 |
| `toast` | 開始、完了、軽い失敗など、一時的な通知 | 再読が必要な error、操作不能の理由、重要な監査情報 |
| `empty_state` | 初回実行前、条件に一致する結果がない、期限切れで結果がない | 実行中、詳細な error、既に結果が存在する画面 |

### progress

`progress` は「待っている理由」を表します。件数が分かる import や aggregation では determinate progress を使い、外部 API の応答待ちのように全体量が不明な場合は indeterminate loading と step message を組み合わせます。

### toast

`toast` は lifecycle の節目を軽く伝えるために使います。

- 開始: 「集計を開始しました」
- cache hit: 「同じ条件の集計結果を表示しました」
- 成功: 「集計が完了しました」
- 軽い失敗: 「再実行を開始できませんでした」

ただし、失敗理由、retry 可否、問い合わせ用 ID は画面内にも残してください。

### empty_state

`empty_state` は「まだ結果がない理由」と「次の行動」を示します。

- 初回実行前: 「条件を選んで Run aggregation を押してください」
- 条件に一致するデータがない: 「期間を広げるか filter を変更してください」
- 結果の TTL 切れ: 「結果の保存期限が切れました。再実行してください」

空状態と失敗状態を混ぜないことが重要です。失敗は alert と retry 導線、未実行や結果なしは empty_state を使います。

## 4. timeout / retry / cache TTL の設計基準

README.ja.md のサンプルでは、timeout 5 秒、retry 1 回、cache TTL 3 分という短い値を使っています。これはデモとして分かりやすい設定であり、本番では処理の重さ、外部依存、ユーザー期待、コストに応じて決めます。

### timeout

Timeout は「ユーザーが待つ時間」ではなく、「1 execution に許す実行予算」です。HTTP request、worker 処理、外部 API 呼び出し、DB query で別々に設定します。

- UI request timeout: 進捗 partial の取得などは短く保つ。
- Enqueue timeout: job 作成だけなので数秒以内にする。
- Worker execution timeout: 処理単位ごとの予算を決める。
- External API timeout: 相手 API の SLA と retry 方針に合わせる。
- DB timeout: 全表 scan を防ぐため、query budget を明示する。

基準例:

| 処理 | timeout の考え方 |
| --- | --- |
| 小さな集計 / preview | 数秒。超えたら background job に切り替える。 |
| CSV import | ファイルサイズや行数に応じた worker timeout。UI は進捗表示にする。 |
| 外部 API 同期 | API ごとの timeout と rate limit を尊重し、全体 timeout を別に持つ。 |
| レポート生成 | 生成予算を長めに取り、結果は result reference として保存する。 |

### retry

Retry は一時的な失敗にだけ使います。入力不備、権限不足、業務ルール違反は retry しても成功しないため、ユーザーに修正を促します。

- 自動 retry は短い backoff 付きで 1〜数回に制限する。
- 外部 API の `429` / `503` / network timeout は retry 候補にする。
- validation error、authorization error、not found は retry しない。
- retry のたびに attempt、execution ID、error code を記録する。
- 非冪等な処理では idempotency key と重複検知を必須にする。

手動 retry は、ユーザーが条件を確認してから実行できるようにします。自動 retry が尽きたあとに手動 retry を出す場合は、「一時的な障害の可能性があります。時間をおいて再実行してください。」のように理由を説明します。

### cache TTL

Cache TTL は、同じ入力条件の結果をどれくらい再利用してよいかという業務判断です。短すぎると再計算が増え、長すぎると古い結果を見せます。

- 入力条件を canonicalize して hash 化し、cache key にする。
- 結果に `GeneratedAt` と `ExpiresAt` を表示する。
- データの鮮度が重要な画面では TTL を短くするか、明示的な refresh button を置く。
- 高コストな集計では TTL を長めにし、cache hit を UI に表示する。
- 権限や tenant が違うユーザー間で cache が混ざらないよう、tenant ID / actor scope を key に含める。

基準例:

| データの性質 | TTL の目安 |
| --- | --- |
| 開発デモ / サンプル | 数分。動作確認しやすさを優先する。 |
| 運用 dashboard | 数十秒〜数分。鮮度と負荷の balance を取る。 |
| 日次レポート | 数時間〜1 日。生成日時を明示する。 |
| 請求・監査データ | cache より正確性を優先し、snapshot / version を保存する。 |

TTL 切れの結果にアクセスされた場合は、いきなり空白にせず、empty_state で期限切れを説明し、再実行 button を出します。

## 5. 本番では in-memory ではなく DB や外部 cache を検討する条件

README.ja.md の Heavy Job Template は、サンプルとして in-memory cache を使います。本番では、プロセス再起動、複数 instance、worker 分離、監査、結果サイズを考慮し、DB、Redis / Memcached、object storage、queue を検討してください。

### in-memory で足りる条件

- 開発環境、デモ、単一プロセスの短命な preview。
- 失われても業務影響がない cache。
- ユーザーが再実行すればすぐ復元できる軽い処理。
- job 履歴や監査ログを残す必要がない画面。

### DB を検討する条件

- job 履歴、状態遷移、誰が実行したかを残す必要がある。
- 失敗時に support が過去 execution を調査する。
- retry、cancel、timeout、部分成功を確実に管理したい。
- 複数 web instance / worker から同じ job state を参照する。
- 再起動後も進捗や結果参照を失ってはいけない。
- tenant / organization ごとの権限チェックを job record に適用したい。

DB には Job record、attempt、state transition、input hash、result reference、error code を保存します。大きな結果 payload は DB に直接入れず、object storage や集計 table を参照する形にします。

### 外部 cache を検討する条件

- 同じ条件の集計結果を複数 instance で共有したい。
- cache TTL をサーバー再起動後も維持したい。
- progress を高頻度に更新するが、DB への書き込みを増やしたくない。
- rate limit や distributed lock が必要。
- queue worker と web process の間で軽量な状態共有をしたい。

Redis などの外部 cache を使う場合でも、監査が必要な job の最終 state は DB に残す設計を推奨します。cache は「速く読むための現在値」、DB は「正として残す履歴」と分けると、障害時の復旧がしやすくなります。

### queue / object storage を検討する条件

- job 実行を web process から分離したい。
- CPU / memory を使う処理を worker pool に流したい。
- 大きな CSV、画像、PDF、レポートを生成する。
- 再試行、dead-letter queue、rate limit、優先度制御が必要。

Queue は execution の配送、DB は状態と履歴、object storage は大きな成果物、cache は短命な progress / result summary という役割分担にします。

## 実装チェックリスト

- Job model に execution ID、state、progress、result reference、error がある。
- state transition が明示され、完了済み job を直接 running に戻していない。
- 実行開始 Action は入力検証、二重送信防止、enqueue、進捗導線を持つ。
- 進捗表示は progress と step message を使い、完了後は polling を止める。
- 成功時は result reference、生成日時、入力条件、cache hit の有無を表示する。
- 失敗時は alert と retry / 修正 / 問い合わせ導線を画面内に残す。
- progress、toast、empty_state の役割を混ぜていない。
- timeout、retry、cache TTL はデモ値のままではなく、業務要件に合わせて決めている。
- retry は retryable error のみに限定し、idempotency key で二重実行を防いでいる。
- tenant / actor / input hash を cache key または DB record に含め、権限境界をまたいで結果を共有していない。
- 本番で in-memory を使う場合、再起動や複数 instance で失われても問題ないことを確認している。
- 監査、複数 instance、worker 分離、大きな結果が必要な場合は DB / 外部 cache / queue / object storage を検討している。
