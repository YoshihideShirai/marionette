# プロジェクト構成ガイド

[English](../01-project-structure.md) | 日本語

このドキュメントは、Marionette を使う **ユーザーアプリ側** の推奨構成を説明します。ここで扱うのは、Marionette 本体の開発ではなく、Marionette を利用して管理画面や社内ツールを作るアプリケーションの配置方針です。

## 基本方針

Marionette アプリでは、画面を組み立てる Page、ユーザー操作を処理する Action、リクエストやセッションから復元する State、アプリ固有の UI helper を分けておくと、画面数が増えても責務を追いやすくなります。

- **Page**: 画面単位の読み取り処理、データ取得、HTML の組み立て。
- **Action**: form submit、button click、filter change など、ユーザー操作による変更処理。
- **State**: query parameter、form value、session 由来の値、画面間で共有する条件。
- **UI helper**: layout、共通 panel、業務固有の小さな表示部品など、アプリ内だけで使う composition helper。

小規模なうちはファイル数を増やしすぎず、画面や操作が増えてきたら責務ごとに package を分けます。

## 1. 小規模アプリ向けの最小構成

画面が数個で、管理者や開発者向けの小さな internal tool であれば、最初は 1 つの `internal/app` package にまとめても構いません。

```text
my-admin-app/
  cmd/server/main.go
  internal/app/
    app.go
    routes.go
    pages.go
    actions.go
    state.go
    ui.go
```

各ファイルの役割は次のとおりです。

| ファイル | 役割 |
| --- | --- |
| `cmd/server/main.go` | executable の entry point。設定の読み込み、Marionette app の生成、HTTP server の起動を行います。 |
| `internal/app/app.go` | app の初期化、middleware、依存 object の組み立てなどを置きます。 |
| `internal/app/routes.go` | URL path と Page / Action handler の対応を登録します。 |
| `internal/app/pages.go` | dashboard や設定画面など、画面単位の render 関数を置きます。 |
| `internal/app/actions.go` | form submit や状態変更を伴う POST / htmx action handler を置きます。 |
| `internal/app/state.go` | query parameter、filter、form input、session から復元する値の型と parser を置きます。 |
| `internal/app/ui.go` | app shell、section、card など、アプリ固有の小さな UI helper を置きます。 |

この構成は、以下のような場合に向いています。

- 画面数が少なく、Page / Action の数も少ない。
- 1 つの domain、または 1 つの workflow だけを扱う。
- 将来の分割よりも、最初の実装速度と見通しを優先したい。

ただし、`pages.go` や `actions.go` が大きくなり、目的の handler を探しにくくなったら、中規模以上の構成へ移行します。

## 2. 中規模以上の管理画面向けの分割構成

画面数が増え、dashboard、user 管理、job 管理、settings など複数の機能領域を持つ場合は、責務ごとに package を分ける構成を推奨します。

```text
my-admin-app/
  cmd/server/main.go
  internal/app/
    app.go
    routes.go
  internal/pages/
    dashboard.go
    users.go
    settings.go
  internal/actions/
    user_actions.go
    job_actions.go
  internal/state/
    query_state.go
    session_state.go
  internal/ui/
    layout.go
    components.go
  internal/assets/
    assets.go
```

各 directory の役割は次のとおりです。

| Directory / file | 役割 |
| --- | --- |
| `cmd/server/main.go` | binary の entry point。production / development の設定、server 起動、graceful shutdown などを扱います。 |
| `internal/app/app.go` | Marionette app の生成、依存関係の注入、共通 middleware、app-wide な初期化を担当します。 |
| `internal/app/routes.go` | route registration の中心です。`pages` と `actions` を URL に接続します。 |
| `internal/pages/` | 画面を表示する handler や render 関数を置きます。画面単位、または機能領域単位でファイルを分けます。 |
| `internal/actions/` | 変更を伴う操作を処理します。validation、保存、redirect、partial update response を扱います。 |
| `internal/state/` | URL query、session、form input、filter、pagination などの型と復元処理を置きます。 |
| `internal/ui/` | アプリ固有の layout と小さな composition helper を置きます。Marionette 本体に追加する汎用 component ではありません。 |
| `internal/assets/` | embed した CSS、画像、JavaScript、favicon など、アプリ固有 asset の登録や配信 helper を置きます。 |

### 分割の目安

次の兆候が出たら、最小構成から分割構成へ移行するタイミングです。

- 1 つの `pages.go` に複数の独立した画面が混在している。
- Action handler が増え、どの Page から呼ばれるか追いにくい。
- query state、form state、session state の型が増えてきた。
- 共通 layout や card helper が複数画面から使われ始めた。
- route registration が長くなり、画面の一覧として読みにくい。

## 3. Page / Action / State / UI helper の置き場所

### Page は `internal/pages/`

Page は、基本的に HTTP request から表示に必要な state を復元し、必要なデータを読み、Marionette の `frontend` component で HTML を組み立てます。

例:

```text
internal/pages/
  dashboard.go  # dashboard page
  users.go      # user list / user detail page
  settings.go   # settings page
```

Page に置くもの:

- GET request の handler。
- 画面単位の render 関数。
- Page title、breadcrumb、table、form など、画面全体の composition。
- Page 固有で、他の画面から再利用しない小さな helper。

Page に置きすぎないもの:

- DB 更新、外部 API への書き込み、メール送信などの変更処理。
- 複数画面で共有する query / session parser。
- アプリ全体で使う layout helper。

### Action は `internal/actions/`

Action は、ユーザー操作によって状態を変更する処理を担当します。POST request、htmx request、button click、form submit、filter change などが該当します。

例:

```text
internal/actions/
  user_actions.go  # create, update, delete, role change
  job_actions.go   # enqueue, retry, cancel
```

Action に置くもの:

- form submit の validation。
- database や durable state の更新。
- session の更新。
- redirect や flash message の設定。
- htmx partial update 用 response の組み立て。

Action から Page の大きな render 関数を直接呼びすぎると依存が複雑になります。partial update に必要な小さな表示 helper は `internal/ui/` に寄せるか、対象 Page と同じ package 内の小さな関数として切り出します。

### State は `internal/state/`

State は、request から復元できる値や、Page と Action の両方で共有する入力値を明示的な型にしたものです。

例:

```text
internal/state/
  query_state.go    # search, sort, pagination, filter
  session_state.go  # login user, selected workspace, flash-like values
```

State に置くもの:

- URL query から復元する filter / sort / pagination。
- form value を受け取る struct。
- validation error を持つ form state。
- session から復元する user context や workspace context。
- Page と Action で共有する request parsing logic。

State は「どこに保存するか」だけではなく「どの lifetime の値か」を表す場所でもあります。share 可能な URL state、request-local な form state、server-side session、durable storage に保存する業務 data を混同しないようにします。

### UI helper は `internal/ui/`

UI helper は、Marionette の `frontend` component を組み合わせて、アプリ固有の見た目や繰り返し構造を作るための関数です。

例:

```text
internal/ui/
  layout.go      # app shell, navigation, page frame
  components.go  # status badge, metric card, empty state helper
```

UI helper に置くもの:

- app shell、sidebar、navbar、footer などの layout。
- 業務固有の badge、card、empty state、toolbar。
- 複数 Page で使う section、panel、form wrapper。
- Marionette の component を組み合わせただけの、アプリ専用 composition。

UI helper に置くべきでないもの:

- Marionette 本体へ追加すべき汎用 component。
- Page 固有の business logic。
- Action の validation や保存処理。
- session や database からの state 復元。

## 4. Marionette 本体の構成と混同しない

このリポジトリには Marionette 本体の実装として、`backend/`、`frontend/`、`templates/components/` などの directory があります。これらは **フレームワーク側の実装場所** であり、通常のユーザーアプリが真似して同じ directory を作る必要はありません。

| Marionette 本体 | 役割 | ユーザーアプリ側での対応 |
| --- | --- | --- |
| `backend/` | Marionette の App、Context、routing、runtime などの backend API 実装。 | ユーザーアプリでは `internal/app/` に app 初期化や route registration を置きます。 |
| `frontend/` | Marionette が提供する Go 側 UI component API の実装。 | ユーザーアプリでは `frontend` package を利用し、アプリ固有の組み合わせを `internal/ui/` に置きます。 |
| `templates/components/` | Marionette 本体の component template や生成元に関する実装。 | ユーザーアプリの画面 helper は `internal/pages/` または `internal/ui/` に置きます。 |

注意点:

- ユーザーアプリに `backend/` や `frontend/` という directory を作る必要はありません。
- アプリ固有の layout や helper を Marionette 本体の `frontend/` に追加しないでください。
- Marionette 本体へ汎用 component を追加する作業と、ユーザーアプリで画面を組み立てる作業は分けて考えます。
- `templates/components/` は Marionette 本体の component 実装に関係する場所です。ユーザーアプリの Page や partial template の置き場所として扱わないでください。
- ユーザーアプリでは、まず `frontend` の既存 component を使い、足りない組み合わせだけを `internal/ui/` に app-specific helper として追加します。

## 推奨構成のまとめ

最初は小さく始め、画面数や責務が増えたら分割します。

```text
my-admin-app/
  cmd/server/main.go
  internal/app/
    app.go
    routes.go
  internal/pages/
    dashboard.go
    users.go
    settings.go
  internal/actions/
    user_actions.go
    job_actions.go
  internal/state/
    query_state.go
    session_state.go
  internal/ui/
    layout.go
    components.go
  internal/assets/
    assets.go
```

この構成では、`internal/app` が app の組み立てと routing を担当し、`internal/pages` が読み取りと画面構成、`internal/actions` が変更処理、`internal/state` が request / session / form state、`internal/ui` がアプリ固有の表示 helper を担当します。Marionette 本体の `backend/`、`frontend/`、`templates/components/` は framework 実装の directory であり、ユーザーアプリ側の構成とは責務が異なります。
