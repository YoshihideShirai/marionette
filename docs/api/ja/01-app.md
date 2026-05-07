## 2. App

### `New() *App`
- 新しい `*App` を返します。
- 初期状態:
  - `state: map[string]any{}`
  - `pages: map[string]Handler{}`
  - `actions: map[string]Handler{}`
  - `cookieSecure: false`

### ページオプション
- `type PageOptions struct { Title string }`
- `type PageOption func(*PageOptions)`
- `WithTitle(title string) PageOption`
  - 前後空白を除去し、ページの HTML `<title>` を設定します。

### `Page(path string, fn Handler, options ...PageOption)`
- `GET` 用のフルページハンドラを登録します。
- `path` は `""` の場合 `"/"`、先頭 `/` がなければ補完されます。

### `Render(fn Handler, options ...PageOption)`
- ルートページ登録の互換エイリアスです。
- `Page("/", fn, options...)` と同等です。

### `Action(name string, fn Handler)`
- `POST` 用のフラグメントハンドラを登録します。

### `Handle(name string, fn Handler)`
- `Action` の互換エイリアスです。

### `SetCookieSecure(secure bool)`
- `marionette_flash` / `marionette_session` Cookie の `Secure` を切り替えます。

### `Assets(prefix string, fsys fs.FS, options ...AssetOption)`
- 静的ファイル配信を登録します。

### `Downloads(prefix string, fsys fs.FS, options ...AssetOption)`
- 添付ダウンロードとして静的ファイル配信を登録します。

### `Asset(name string) string`
- 登録済みアセットプレフィックスに基づく URL を返します。

### App state helpers
- `Global` が名前に付く App state helper は、アプリの全ユーザー・全リクエストで共有される state を読み書きします。
- `SetGlobal(key string, value any)`: mutex で保護しながらアプリ共有 state に書き込みます。
- `GetGlobal(key string) any`: mutex で保護しながらアプリ共有 state を読み取ります。
- `GetGlobalInt(key string) int`: `GetGlobal` の結果を `int` に type assertion し、値がない/`int` でない場合は `0` を返します。
- Deprecated: `Set(key string, value any)` / `Get(key string) any` / `GetInt(key string) int` は互換エイリアスです。アプリ全体の state へアクセスするときは `SetGlobal` / `GetGlobal` / `GetGlobalInt` を使ってください。
