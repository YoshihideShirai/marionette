# Security / Authorization ガイド

[English](../10-security-authz.md) | 日本語

このドキュメントは、Marionette アプリケーションで認証済みユーザー情報と認可チェックを扱うときの基本方針をまとめます。管理画面や社内ツールでは、画面上の表示制御だけでなく、Page と Action の両方でサーバー側の検証を徹底してください。

関連する責務分割は [Routing / Pages / Actions Guide](02-routing-pages-actions.md) を、session や state の置き場所は [State Management Guide](../../state-management.ja.md) を参照します。

## 基本方針

認証と認可は、次の境界を分けて設計します。

- **Authentication**: リクエストが誰によるものかを確定する。
- **Session / Context**: 認証済みユーザーを handler から参照できる形で保持する。
- **Authorization**: その actor が対象 resource に対して操作してよいかを判定する。
- **UI 表示制御**: 操作できないボタンやリンクを表示しない、または disabled にする。
- **Audit log**: 重要操作の actor、action、target、timestamp、result を後から追えるようにする。

UI はユーザー体験を整えるための補助であり、セキュリティ境界ではありません。最終的な可否判定は、必ずサーバー側の Page / Action / service 層で行います。

## 1. 認証済みユーザー情報を Context / session で扱う方針

認証済みユーザー情報は、リクエストごとに middleware で復元し、`Context` から参照できる形にします。Page や Action が cookie、header、session store を個別に読み直す実装は避けます。

推奨する流れは次の通りです。

1. 認証 middleware が session ID や token を検証する。
2. session store や user repository から actor 情報を復元する。
3. `Context` に actor ID、tenant ID、role、permission、session ID などを格納する。
4. Page / Action / service は `Context` から actor を取得して認可判定に渡す。

`Context` に入れる情報は、認可判定と表示分岐に必要な最小限にします。メールアドレス、表示名、role のような画面表示に必要な値は持ってよい一方で、password hash、access token、refresh token、外部 API secret などの秘密情報は直接持ち回らないでください。

例:

```go
type CurrentUser struct {
    ID        string
    TenantID  string
    Email     string
    Role      string
    SessionID string
}

func CurrentUserFromContext(ctx *backend.Context) (CurrentUser, bool) {
    user, ok := ctx.Local["current_user"].(CurrentUser)
    return user, ok
}
```

認証が必須の route では、Page / Action に入る前に middleware で未認証リクエストを止めます。未認証の場合は login へ redirect するか、部分更新リクエストであれば login が必要であることを示す partial response を返します。

## 2. Page 表示時の認可チェック

Page は、画面を表示してよいかを最初に確認します。権限のないユーザーに対して、対象データを読み込んでから UI だけ隠す実装は避けます。

Page 表示時に確認すること:

- actor がその route にアクセスできるか。
- actor が対象 tenant / organization / project に属しているか。
- actor が対象 resource を閲覧できるか。
- 一覧画面では、actor が見てよいデータだけを query 条件に含めているか。
- 詳細画面では、path parameter の ID が actor の権限範囲内か。

例:

```go
func (p *UserPages) Detail(ctx *backend.Context) error {
    actor, ok := CurrentUserFromContext(ctx)
    if !ok {
        return ctx.Redirect("/login")
    }

    userID := ctx.Param("id")
    if allowed := p.Authorizer.CanViewUser(ctx.Request.Context(), actor, userID); !allowed {
        ctx.Writer.WriteHeader(http.StatusForbidden)
        return ui.ForbiddenPage()
    }

    user, err := p.Users.GetVisibleUser(ctx.Request.Context(), actor, userID)
    if err != nil {
        return err
    }
    return ctx.HTML(ui.UserDetailPage(user))
}
```

一覧や検索では、後段でフィルタするのではなく、最初から権限条件を query に含めます。これにより、ページング件数、集計値、検索候補、empty state から権限外データの存在が推測されるリスクを下げられます。

## 3. Action 実行時の認可チェック

Action は、実行直前に必ず認可チェックを行います。Page 表示時にボタンを出していたとしても、ユーザーの権限や対象 resource の状態は送信までに変わる可能性があります。また、HTTP リクエストはブラウザの UI を経由せずに直接送信できます。

Action 実行時に確認すること:

- actor が認証済みか。
- actor がその action を実行できる role / permission を持つか。
- target が actor の tenant / organization / project の範囲内か。
- target が現在その操作を受け付ける状態か。
- form / query / path parameter で指定された target が改ざんされていないか。
- CSRF 対策や idempotency key など、操作種別に応じた保護が有効か。

例:

```go
func (a *UserActions) Delete(ctx *backend.Context) error {
    actor, ok := CurrentUserFromContext(ctx)
    if !ok {
        return ctx.Redirect("/login")
    }

    userID := ctx.Param("id")
    if allowed := a.Authorizer.CanDeleteUser(ctx.Request.Context(), actor, userID); !allowed {
        a.Audit.Record(ctx.Request.Context(), audit.Event{
            Actor:  actor.ID,
            Action: "user.delete",
            Target: userID,
            Result: "denied",
        })
        ctx.Writer.WriteHeader(http.StatusForbidden)
        return ui.ForbiddenAlert()
    }

    if err := a.Users.Delete(ctx.Request.Context(), actor, userID); err != nil {
        a.Audit.Record(ctx.Request.Context(), audit.Event{
            Actor:  actor.ID,
            Action: "user.delete",
            Target: userID,
            Result: "failed",
        })
        return err
    }

    a.Audit.Record(ctx.Request.Context(), audit.Event{
        Actor:  actor.ID,
        Action: "user.delete",
        Target: userID,
        Result: "succeeded",
    })
    ctx.FlashSuccess("ユーザーを削除しました")
    return ctx.Redirect("/admin/users")
}
```

認可エラーは validation error とは分けて扱います。入力の修正で解決できないため、field error ではなく page-level alert、403 response、または安全な redirect を使います。

## 4. UI でボタンを隠すだけでは不十分

UI で削除ボタンや承認ボタンを隠すことは、誤操作を減らし、画面を分かりやすくするために有効です。ただし、これは認可チェックの代替にはなりません。

ボタンを隠すだけでは不十分な理由:

- ユーザーは開発者ツール、curl、スクリプトから Action URL に直接リクエストできる。
- Page 表示後に role、permission、target の状態が変わることがある。
- htmx の partial endpoint や JSON endpoint が別経路で呼ばれることがある。
- disabled 属性や hidden 要素はクライアント側で書き換えられる。

したがって、次の 2 段構えにします。

1. **Page / UI**: actor の権限に応じて、操作ボタンの表示、disabled、説明文を調整する。
2. **Action / service**: リクエストごとに actor、action、target の組み合わせを検証し、許可されない操作を拒否する。

UI 側では「なぜ操作できないか」がユーザーに分かるよう、必要に応じて disabled button、tooltip、alert、empty state を使います。一方で、Action 側では UI 表示状態を信用せず、送られてきた値をすべて検証します。

## 5. 危険操作の確認 UI

削除、承認、再実行、無効化、権限変更、支払い、外部送信などの危険操作は、実行前に確認 UI を用意します。確認 UI は誤操作を防ぐためのものであり、認可チェックや監査ログとは別に必要です。

危険操作の確認 UI で示すこと:

- 実行する action 名。
- 対象 target の名前、ID、件数。
- 操作後に戻せるかどうか。
- 影響範囲と副作用。
- 実行後に通知、外部送信、ジョブ起動が発生するか。

操作種別ごとの例:

- **削除**: 対象名を表示し、必要に応じて対象名の再入力を求める。
- **承認**: 承認後に公開、請求、通知などが発生する場合は明記する。
- **再実行**: 既存結果が上書きされるか、重複実行になるか、外部 API が再度呼ばれるかを示す。
- **権限変更**: 付与・剥奪される role / permission と影響ユーザーを確認する。

確認 UI の実装では、modal、drawer、confirmation page などを使えます。軽微な操作は modal で十分ですが、取り消し不能な操作や影響範囲が大きい操作は専用の確認 page にして、対象情報と注意事項を読みやすく表示します。

確認 UI から送信された Action でも、再度 server-side validation と authorization を実行します。確認画面を通ったことを hidden input だけで信用してはいけません。

## 6. 監査ログに残すべき情報

監査ログは、重要操作について「誰が、何を、どの対象に、いつ、どういう結果で実行したか」を後から確認するための記録です。通常の application log とは別に、検索しやすい構造化データとして残します。

最低限、次の情報を記録します。

| 項目 | 内容 | 例 |
| --- | --- | --- |
| actor | 操作した主体。user ID、service account ID、session ID など。 | `user_123` |
| action | 実行しようとした操作。namespace を含む安定した名前にする。 | `invoice.approve` |
| target | 操作対象。resource type と ID を分けて持つと検索しやすい。 | `invoice:inv_456` |
| timestamp | サーバー側で記録した発生時刻。 | `2026-05-11T10:15:30Z` |
| result | 成功、失敗、拒否、キャンセルなどの結果。 | `succeeded`, `failed`, `denied` |

必要に応じて、次の情報も追加します。

- request ID / trace ID。
- actor の tenant ID / organization ID。
- client IP、user agent。
- 変更前後の要約。ただし秘密情報や過剰な個人情報は含めない。
- 失敗理由の分類。詳細な stack trace ではなく、`permission_denied`、`validation_failed`、`conflict` などの安全な値にする。
- idempotency key、job ID、external request ID。

監査ログは、成功時だけでなく拒否や失敗も記録します。特に権限不足、対象外 tenant へのアクセス、連続した失敗、危険操作のキャンセルは、後から調査できるようにしておくと有用です。

## 実装前チェックリスト

- 認証済みユーザーを middleware で復元し、`Context` から一貫して参照できる。
- `Context` に秘密情報を直接持ち回っていない。
- Page 表示時に route、tenant、target の閲覧権限を確認している。
- 一覧や検索の query に権限条件を含めている。
- Action 実行時に actor、action、target の組み合わせを再検証している。
- UI のボタン表示制御を、server-side authorization の代替として扱っていない。
- 削除、承認、再実行などの危険操作に確認 UI がある。
- 監査ログに actor、action、target、timestamp、result を構造化して残している。
- 認可エラー、成功、失敗、拒否の各パスをテストしている。
