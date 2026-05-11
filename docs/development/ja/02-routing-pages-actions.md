# Routing / Pages / Actions ガイド

[English](../02-routing-pages-actions.md) | 日本語

このドキュメントは、Marionette アプリケーションで URL、Page、Action をどう分けるかの判断基準を説明します。特に admin UI や社内ツールで、画面が増えても route と handler の責務を追いやすくするための設計メモです。

State の置き場所、URL query・form・session・global state の使い分けは [State Management ガイド](../../state-management.ja.md) を参照してください。このドキュメントでは、State の詳細には踏み込みすぎず、Page と Action の境界に集中します。

## 基本方針

Marionette アプリでは、**Page は表示の入口**、**Action はユーザー操作の入口**として分けます。

- **Page**: request から表示 state を復元し、必要なデータを読み、画面または部分画面の HTML を返す。
- **Action**: form submit、button click、filter change などの入力を受け取り、検証と業務処理を呼び出し、結果の表示方法を決める。
- **State**: Page と Action の両方で共有する query、form 値、validation error、session 由来の値を型として表す。
- **Service / domain**: DB 更新、外部 API 呼び出し、業務ルール、トランザクションなどを担当する。
- **UI helper / component**: HTML の組み立て、共通 layout、panel、table、form section などを担当する。

handler を薄く保つほど、route 一覧から「どの URL がどの画面・操作につながるか」を追いやすくなります。

## 1. Page の分割基準

Page は「ユーザーから見た表示単位」を基準に分けます。実装上は、URL 単位、業務画面単位、partial update 単位の 3 つを組み合わせて判断します。

### URL 単位

まずは route に対応する URL ごとに Page を置きます。URL が変わるということは、ブラウザ履歴、bookmark、reload、権限、title、breadcrumb が変わる可能性があるためです。

例:

```text
GET /admin/users          -> pages.UsersIndex
GET /admin/users/{id}     -> pages.UserDetail
GET /admin/settings       -> pages.Settings
```

URL 単位で Page を分けると、route 定義から画面構成を読み取りやすくなります。ただし、同じ URL 内にある小さな card や table body まで無条件に Page と呼ぶ必要はありません。

### 業務画面単位

URL が同じでも、業務上の意味が大きく異なる表示は Page 関数や render 関数を分けます。

例:

- ユーザー一覧とユーザー詳細。
- 請求書一覧と請求書作成。
- ジョブ実行履歴とジョブ設定。
- 管理者向け設定と一般ユーザー向け設定。

業務画面単位で分けると、必要な state、権限、読み込むデータ、表示 component が自然にまとまります。1 つの Page が複数の業務文脈を持ち始めたら、ファイルまたは関数を分ける合図です。

### partial update 単位

htmx などで画面の一部だけを更新する場合は、**partial update の対象ごとに再描画関数を切り出す**と扱いやすくなります。

例:

- 検索条件を変えたときに table body だけ更新する。
- 保存後に detail panel だけ更新する。
- filter と chart が同じ state から再描画される。
- modal の中身だけを lazy load する。

ただし、partial update 用の関数は「独立した Page」ではなく、「Page の一部を再描画する helper」として扱うと責務が膨らみにくくなります。URL、権限、主要な state 復元は親 Page 側に置き、partial は必要な state と data を受け取って HTML を返す形に寄せます。

## 2. Action の責務

Action は、ユーザー操作を受け取り、処理結果を画面に返すための薄い orchestration 層です。典型的な流れは次の 4 つです。

1. **入力取得**: `ctx.FormValue`、path parameter、query parameter、session 由来の値などを読み取る。
2. **バリデーション呼び出し**: form state や validator に入力値を渡し、field error や page-level error を得る。
3. **ドメイン処理呼び出し**: service / domain 層のメソッドを呼び出し、保存、削除、外部連携などを実行する。
4. **結果表示**: redirect、flash、toast、partial HTML、form 再描画など、ユーザーに返す表示を選ぶ。

Action 自体は「何を呼ぶか」と「結果をどう返すか」を決める場所です。DB の更新手順、外部 API の request/response 詳細、複雑な HTML 構築を Action に直接書き始めると、テストしづらくなり、route の見通しも悪くなります。

## 3. Action に書きすぎないための基準

Action が長くなったら、次の基準で責務を逃がします。

### DB 更新や外部 API 呼び出しは service 層へ逃がす

次の処理は Action に直接書かず、service / domain 層へ寄せます。

- DB transaction、insert / update / delete、複数 table の整合性維持。
- 外部 API 呼び出し、retry、timeout、response mapping。
- 権限を含む業務ルールの判定。
- 監査ログ、通知、queue 投入など、画面操作の副作用。

Action は service に必要な input を渡し、成功・失敗を受け取って表示へ変換します。これにより、業務処理を HTTP handler なしでテストできます。

### 表示 HTML の組み立ては UI helper / component に寄せる

次の処理は UI helper / component へ寄せます。

- table、card、form section、empty state、alert の組み立て。
- validation error 付き form の再描画。
- htmx target に返す partial HTML。
- 複数 Page / Action で使う共通 layout や panel。

Action 内に HTML ノードを大量に並べると、入力処理と表示構造が混ざります。Action では `renderUserForm(state)` や `ui.UserTable(rows, filters)` のような helper を呼び出すだけにすると、差分が読みやすくなります。

## 4. 良い例 / 避けたい例

### 良い例: Action は入力、検証、service 呼び出し、結果表示に集中する

```go
func (h *UserActions) Create(ctx *backend.Context) error {
    form := state.UserFormFromRequest(ctx)
    if result := form.Validate(); !result.OK() {
        form.Errors = result.Errors
        return ctx.HTML(ui.UserForm(form))
    }

    user, err := h.Users.Create(ctx.Request().Context(), service.CreateUserInput{
        Name:  form.Name,
        Email: form.Email,
    })
    if err != nil {
        form.PageError = "ユーザーを作成できませんでした"
        return ctx.HTML(ui.UserForm(form))
    }

    ctx.Flash("success", "ユーザーを作成しました")
    return ctx.Redirect("/admin/users/" + user.ID)
}
```

この例では、Action が form の取得、validation、service 呼び出し、結果表示だけを担当しています。保存処理は `h.Users.Create`、HTML の組み立ては `ui.UserForm` に寄せています。

### 避けたい例: Action に DB 更新と HTML 構築を詰め込む

```go
func CreateUser(ctx *backend.Context, db *sql.DB) error {
    name := ctx.FormValue("name")
    email := ctx.FormValue("email")
    if name == "" || !strings.Contains(email, "@") {
        return ctx.HTML(html.Div(
            html.P("入力内容を確認してください"),
            html.Input(html.Name("name"), html.Value(name)),
            html.Input(html.Name("email"), html.Value(email)),
        ))
    }

    tx, err := db.BeginTx(ctx.Request().Context(), nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    row := tx.QueryRowContext(ctx.Request().Context(),
        "insert into users (name, email) values ($1, $2) returning id",
        name, email,
    )

    var id string
    if err := row.Scan(&id); err != nil {
        return err
    }
    if err := tx.Commit(); err != nil {
        return err
    }

    return ctx.HTML(html.Div(
        html.H2("作成しました"),
        html.A(html.Href("/admin/users/"+id), html.Text("詳細へ")),
    ))
}
```

この例は、validation、DB transaction、SQL、HTML 構築が 1 つの Action に混ざっています。画面変更でも業務処理変更でも同じ関数を触る必要があり、テスト対象も大きくなります。

## 実装前チェックリスト

- URL が変わる表示は Page として分かれている。
- 業務画面ごとに state、data load、render 関数が追いやすい粒度になっている。
- partial update は親 Page の state を使って再描画できる helper に切り出している。
- Action は入力取得、validation 呼び出し、service 呼び出し、結果表示に集中している。
- DB 更新、外部 API 呼び出し、業務ルールは service / domain 層に寄せている。
- HTML の組み立ては UI helper / component に寄せている。
- State の寿命や置き場所に迷うものは [State Management ガイド](../../state-management.ja.md) で確認している。
