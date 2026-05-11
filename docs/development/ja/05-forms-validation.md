# Forms / Validation ガイド

[English](../05-forms-validation.md) | 日本語

このドキュメントは、Marionette アプリケーションで入力フォーム、サーバー側バリデーション、Action の戻り方をどう設計するかの基準をまとめます。フォームは「入力値」「エラー」「保存結果」「再描画の範囲」が混ざりやすいため、Page、Action、state、service、UI helper の責務を明確に分けます。

関連する画面分割や Action の責務は [Routing / Pages / Actions ガイド](02-routing-pages-actions.md) を、入力コンポーネントの詳細は [Form APIs](../../api/04-form-apis.md) を参照してください。

## 基本方針

Marionette のフォームでは、**入力 state を明示的な型に集め、Action で検証し、同じ render helper に state を戻す**ことを標準にします。

- Page は初期値、選択肢、既存データを読み込み、空または既存値入りの form state を作る。
- Action は request から入力を取得し、サーバー側バリデーションを呼び出し、成功・失敗に応じた表示を選ぶ。
- state は入力値、field error、page-level error、submit 後に再表示したい値を保持する。
- service / domain は保存、更新、外部 API 連携、業務ルールの検証を担当する。
- UI helper は form state だけを受け取り、入力値とエラーを同じ構造で描画する。

クライアント側の `required` や入力 type は補助として使えますが、信頼できる検証は必ずサーバー側で行います。

## 1. 入力フォームの基本構成

入力フォームは、Page が渡す state と、state を描画する helper を中心に組み立てます。

推奨する構成は次のとおりです。

```text
internal/<app>/state/user_form.go   # 入力値、エラー、Validate
internal/<app>/pages/users.go       # GET 表示、初期値・選択肢の読み込み
internal/<app>/actions/users.go     # POST submit、検証、保存、戻り方の選択
internal/<app>/ui/user_form.go      # form state から HTML を組み立てる helper
internal/<app>/service/users.go     # DB 更新、重複確認、業務ルール
```

form state には、少なくとも次を持たせます。

- ユーザーが入力した値。
- field ごとのエラーメッセージ。
- フォーム全体に関係する page-level error。
- select / radio などの選択肢を再描画するために必要な補助データ。
- submit 後に部分更新する場合の target ID や表示モード。

例:

```go
type UserForm struct {
    Name  string
    Email string
    Role  string

    RoleOptions []frontend.SelectOption
    Errors      map[string]string
    PageError   string
}

func NewUserForm(roleOptions []frontend.SelectOption) UserForm {
    return UserForm{
        RoleOptions: roleOptions,
        Errors:      map[string]string{},
    }
}
```

描画側は、`Name` や `Email` の値と `Errors["email"]` を同じ state から読むようにします。これにより、初回表示、validation 失敗後の再表示、編集画面の初期値表示を同じ helper で扱えます。

## 2. サーバー側バリデーションの置き場所

サーバー側バリデーションは、**入力形式の検証**と**業務ルールの検証**を分けて置きます。

### form state / validator に置くもの

次のような、request 入力だけで判断できる検証は form state または validator に置きます。

- 必須入力。
- 文字数、数値範囲、日付範囲。
- メールアドレスや slug などの基本形式。
- select / radio の値が許可リストに含まれるか。
- 2 つの入力値の単純な整合性。例: 開始日が終了日より後ではない。

```go
func (f *UserForm) Validate() bool {
    f.Errors = map[string]string{}

    if strings.TrimSpace(f.Name) == "" {
        f.Errors["name"] = "名前を入力してください"
    }
    if !strings.Contains(f.Email, "@") {
        f.Errors["email"] = "有効なメールアドレスを入力してください"
    }
    if f.Role == "" {
        f.Errors["role"] = "ロールを選択してください"
    }

    return len(f.Errors) == 0
}
```

### service / domain に置くもの

次のような、永続化データ、外部 API、権限、複数オブジェクトの状態が必要な検証は service / domain に置きます。

- メールアドレスやコードの重複確認。
- 現在のユーザーが対象リソースを変更できるか。
- 契約状態、在庫、上限数、ワークフロー状態などの業務ルール。
- DB transaction 内でしか確定できない整合性。
- 外部サービスの応答に依存する検証。

Action は `form.Validate()` で入力形式を確認したあと、service に渡して業務エラーを受け取り、ユーザーに見せるメッセージへ変換します。

## 3. Action での入力取得と検証の流れ

Action の標準フローは次の順序にします。

1. `ctx.FormValue`、path parameter、query parameter から入力を取得する。
2. 取得した値で form state を作る。
3. `Validate` を呼び、field error があれば同じ画面を再描画する。
4. service / domain に入力を渡す。
5. service error を user-facing error、log-only error、再試行可能な error に分類する。
6. 成功時は redirect、partial update、toast のいずれかを選ぶ。

例:

```go
func UserFormFromRequest(ctx *backend.Context, roles []frontend.SelectOption) UserForm {
    form := NewUserForm(roles)
    form.Name = ctx.FormValue("name")
    form.Email = ctx.FormValue("email")
    form.Role = ctx.FormValue("role")
    return form
}

func (h *UserActions) Create(ctx *backend.Context) error {
    roles := h.Users.RoleOptions(ctx.Request.Context())
    form := state.UserFormFromRequest(ctx, roles)

    if !form.Validate() {
        return ctx.HTML(ui.UserForm(form))
    }

    user, err := h.Users.Create(ctx.Request.Context(), service.CreateUserInput{
        Name:  form.Name,
        Email: form.Email,
        Role:  form.Role,
    })
    if err != nil {
        form.PageError = "ユーザーを作成できませんでした。時間をおいて再度お試しください。"
        h.Logger.Error("create user failed", "err", err)
        return ctx.HTML(ui.UserForm(form))
    }

    ctx.FlashSuccess("ユーザーを作成しました")
    return ctx.Redirect("/admin/users/" + user.ID)
}
```

Action 内で SQL、複雑な HTML、細かい業務判定を増やさないことが重要です。Action は入力を集め、検証と保存を呼び、表示結果を選ぶ orchestration 層に留めます。

## 4. 入力エラーを同じ画面に返す例

入力エラーは、ユーザーが直前に入力した値を保ったまま同じ form helper に戻します。field error は該当項目の近くに、フォーム全体のエラーは Alert として上部に出します。

```go
func UserForm(form state.UserForm) frontend.Node {
    children := []frontend.Node{}
    if form.PageError != "" {
        children = append(children, frontend.Alert(frontend.AlertProps{
            Title:       "保存できませんでした",
            Description: form.PageError,
            Props:       frontend.ComponentProps{Variant: "error"},
        }))
    }

    children = append(children,
        frontend.FormRow(frontend.FormRowProps{
            ID:    "user-name",
            Label: "名前",
            Error: form.Errors["name"],
            Control: frontend.TextField(frontend.TextFieldProps{
                ID:    "user-name",
                Name:  "name",
                Value: form.Name,
            }),
        }),
        frontend.FormRow(frontend.FormRowProps{
            ID:    "user-email",
            Label: "メールアドレス",
            Error: form.Errors["email"],
            Control: frontend.TextField(frontend.TextFieldProps{
                ID:    "user-email",
                Name:  "email",
                Type:  "email",
                Value: form.Email,
            }),
        }),
        frontend.Submit("保存"),
    )
    return frontend.Form("/admin/users", children...)
}
```

同じ helper を Page と Action の両方から呼ぶと、初回表示とエラー表示の markup 差分が小さくなります。`FormRow` や入力 component の error 表現を使い、`aria-invalid` や説明文との関連付けを component に任せます。

partial update の場合も考え方は同じです。Action は form 全体、form section、または target になっている field group を再描画して返し、入力値とエラーを state から復元します。

## 5. 成功時に redirect / partial update / toast を使い分ける基準

成功時の戻り方は、ユーザーが次に見るべき状態で選びます。

### redirect を使う場合

redirect は、操作後に URL や画面の主文脈が変わるときに使います。

- 作成後に詳細画面へ移動する。
- 編集後に一覧へ戻る。
- 削除後に対象詳細ページから離れる。
- reload や bookmark で同じ完了状態を再現したい。
- POST の再送信を避けたい。

redirect する場合は、次の画面で表示する一時メッセージを `ctx.FlashSuccess` などの flash に入れます。

### partial update を使う場合

partial update は、URL を変えずに画面の一部だけを最新化したいときに使います。

- 設定画面の一部の card だけ保存する。
- 一覧の行内ステータスや詳細 panel だけ更新する。
- modal 内のフォームを送信し、modal 本体または背景の table だけ差し替える。
- htmx target が明確で、更新範囲が小さい。

partial update では、更新後の HTML に成功状態を含めるか、別 target の toast 領域も同時に更新する設計にします。

### toast を使う場合

Toast は、操作は成功したが、ユーザーの主作業を中断したくないときに使います。

- 自動保存、軽い設定変更、コピー完了。
- inline edit の保存成功。
- 画面遷移なしの短い成功通知。
- 失敗しても同じ場所からすぐ再試行できる操作。

Toast だけで完了状態を伝える場合は、画面上の実データも更新済みであることが前提です。データが変わっていないのに toast だけを出すと、ユーザーが成功状態を確認できません。

## 実装チェックリスト

- 入力値、field error、page-level error を form state に集めた。
- Page と Action が同じ form render helper を呼ぶ構成にした。
- request 入力だけで判断できる検証は form state / validator に置いた。
- DB、権限、外部 API、業務ルールに依存する検証は service / domain に置いた。
- validation 失敗時は、入力値を保持して同じ画面または同じ partial target を再描画した。
- 成功時の戻り方を redirect、partial update、toast の基準で選んだ。
- ユーザー向けメッセージとログ向けエラーを混ぜて表示していない。
