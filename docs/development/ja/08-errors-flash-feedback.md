# Errors / Flash / Feedback ガイド

[English](../08-errors-flash-feedback.md) | 日本語

このドキュメントは、Marionette アプリケーションでエラー、flash message、alert、toast、inline error、再実行導線をどう使い分けるかの基準をまとめます。目的は、ユーザーには次に取るべき行動を示し、開発者には調査に必要な情報をログとして残すことです。

フォームの validation error については [Forms / Validation ガイド](05-forms-validation.md) も参照してください。Action の責務分割は [Routing / Pages / Actions ガイド](02-routing-pages-actions.md) を基準にします。

## 基本方針

エラーとフィードバックは、**誰が読む情報か**と**どのくらい画面に残す必要があるか**で分けます。

- ユーザーに見せるメッセージは、原因の概要、現在の状態、次にできる操作を短く伝える。
- ログに残すエラーは、調査に必要な内部エラー、ID、入力の要約、外部 API の状態を含める。
- validation error は inline error として入力欄の近くに出す。
- 操作結果の成功通知は flash または toast を使う。
- 失敗後にユーザーが直せる場合は、再入力・再実行の導線を残す。
- ユーザーが直せない内部エラーは、詳細を隠し、ログと trace ID で追えるようにする。

## 1. ユーザーに見せるエラーとログに残すエラーの分離

ユーザー向けエラーとログ向けエラーは、同じ `err` から発生しても内容を分けます。

### ユーザーに見せるメッセージ

ユーザー向けメッセージには、次を含めます。

- 何が完了しなかったか。
- ユーザーが修正できる場合は、どの入力や条件を確認すべきか。
- 再試行できる場合は、再試行できること。
- 問い合わせが必要な場合は、問い合わせ時に伝える ID。

避けるべき内容:

- SQL、stack trace、外部 API の raw response。
- 秘密情報、token、メールアドレス以外の個人情報、内部ホスト名。
- 開発者にしか分からない型名や package 名。

例:

```go
form.PageError = "ユーザーを保存できませんでした。入力内容を確認して、もう一度お試しください。"
```

### ログに残す情報

ログには、調査に必要な情報を構造化して残します。

- 内部エラーの詳細。
- request ID、job ID、user ID、resource ID。
- 呼び出した service / external API の名前。
- retry 回数、timeout、HTTP status などの技術情報。
- 入力値そのものではなく、必要に応じた要約や安全な識別子。

```go
h.Logger.Error("create user failed",
    "err", err,
    "request_id", requestID,
    "actor_id", actorID,
)
form.PageError = "ユーザーを作成できませんでした。時間をおいて再度お試しください。"
```

Action は service error を受け取ったら、まずログに残す情報を決め、次にユーザーへ返す表現へ変換します。`err.Error()` をそのまま画面に出す実装は避けます。

## 2. alert / toast / inline error の使い分け

フィードバック UI は、メッセージの重要度、表示範囲、残す必要性で選びます。

### inline error

inline error は、特定の入力欄や行に紐づく修正可能な問題に使います。

- 必須項目が空。
- メールアドレス形式が不正。
- 選択肢が無効。
- 行内編集でその行だけ保存できなかった。

原則として、field error は該当 field の近くに表示します。フォーム全体に影響する validation summary は alert と併用できます。

### alert

alert は、画面内に残してユーザーが確認すべき問題や状態に使います。

- 保存失敗、権限不足、外部連携失敗。
- ページ全体に関係する validation summary。
- 削除や無効化など、注意が必要な状態。
- 長時間処理の失敗や一部成功の説明。

alert は画面内の文脈に置き、再読できるようにします。入力フォームでは上部の page-level error として使うと分かりやすくなります。

### toast

Toast は、短時間で消えても問題ない軽い結果通知に使います。

- 保存しました。
- コピーしました。
- キューに投入しました。
- バックグラウンド処理を開始しました。

Toast は作業を中断しない一方で、あとから読み返しにくい UI です。失敗理由を読んで修正が必要なエラーや、法務・請求・権限など重要な通知は alert として残します。

## 3. flash API がある場合の標準利用例

Marionette の `backend.Context` には flash message 用の helper があります。Action で `ctx.FlashSuccess`、`ctx.FlashError`、`ctx.FlashInfo`、`ctx.FlashWarn` を呼び、redirect 後の Page で `ctx.Flashes()` を読み出して表示する形を標準にします。

### Action 側

```go
func (h *UserActions) Update(ctx *backend.Context) error {
    form := state.UserFormFromRequest(ctx)
    if !form.Validate() {
        return ctx.HTML(ui.UserForm(form))
    }

    if err := h.Users.Update(ctx.Request.Context(), form.ToInput()); err != nil {
        h.Logger.Error("update user failed", "err", err, "user_id", ctx.Param("id"))
        form.PageError = "ユーザーを更新できませんでした。時間をおいて再度お試しください。"
        return ctx.HTML(ui.UserForm(form))
    }

    ctx.FlashSuccess("ユーザーを更新しました")
    return ctx.Redirect("/admin/users/" + ctx.Param("id"))
}
```

### Page / layout 側

```go
func UserDetailPage(ctx *backend.Context) frontend.Node {
    flashes := ctx.Flashes()
    return ui.AdminLayout(ui.AdminLayoutProps{
        Flashes: flashes,
        Content: frontend.Div(
            frontend.FlashAlerts(flashes),
            userDetailContent(ctx),
        ),
    })
}
```

flash は redirect をまたいで一度だけ表示したいメッセージに向いています。validation error のように同じ POST response でフォームを再描画する場合は、flash ではなく form state の `Errors` や `PageError` を使います。

## 4. 長時間処理や失敗時の再実行導線

長時間処理は、submit 直後にすべての結果が分からないことを前提に設計します。

### 長時間処理の開始時

- ボタンを loading / disabled にして二重送信を防ぐ。
- 「処理を開始しました」と toast または flash で伝える。
- job ID や進捗ページへのリンクがある場合は表示する。
- 画面を離れても結果を確認できるなら、履歴・通知・詳細ページへの導線を置く。

```go
job, err := h.Jobs.Enqueue(ctx.Request.Context(), input)
if err != nil {
    h.Logger.Error("enqueue import failed", "err", err)
    ctx.FlashError("インポートを開始できませんでした。もう一度お試しください。")
    return ctx.Redirect("/admin/imports/new")
}

ctx.FlashInfo("インポートを開始しました")
return ctx.Redirect("/admin/imports/" + job.ID)
```

### 失敗時の再実行導線

失敗時は、ユーザーが次にできる操作を画面内に残します。

- 入力を直して再送信できるフォームを表示する。
- 同じ条件で再実行するボタンを置く。
- 失敗した job の詳細、ログ要約、再実行可否を表示する。
- 一時的な外部 API 障害なら、時間をおいて再試行できることを伝える。
- 問い合わせが必要なら、request ID / job ID を表示する。

```go
func ImportJobDetail(job service.ImportJob) frontend.Node {
    children := []frontend.Node{frontend.H2(frontend.Text("インポート結果"))}
    if job.Failed {
        children = append(children,
            frontend.Alert(frontend.AlertProps{
                Title:       "インポートに失敗しました",
                Description: "内容を確認して再実行できます。",
                Props:       frontend.ComponentProps{Variant: "error"},
            }),
            frontend.Link(frontend.LinkProps{
                Label: "再実行する",
                Href:  "/admin/imports/" + job.ID + "/retry",
            }),
        )
    }
    return frontend.Section(frontend.SectionProps{}, children...)
}
```

再実行できない失敗の場合も、理由と次の導線を明示します。例: 「この請求書はすでに確定済みのため再実行できません。新しい請求書を作成してください。」のように、できない理由と代替操作をセットで表示します。

## 実装チェックリスト

- `err.Error()` や stack trace をそのままユーザーに表示していない。
- ユーザー向けメッセージは、状態と次の行動を短く示している。
- ログには request ID、resource ID、内部エラーなど調査情報を構造化して残している。
- field error は inline、画面や操作全体の問題は alert、軽い成功通知は toast に分けている。
- redirect 後に一度だけ表示する通知は flash を使っている。
- validation 失敗は flash ではなく form state に戻している。
- 長時間処理では開始、進捗、結果、再実行の導線を用意している。
- 失敗時にユーザーが修正・再試行・問い合わせのどれをすべきか判断できる。
