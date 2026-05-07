## 3. Context

`Context` は各ハンドラに `func(*Context) Node` として渡され、リクエスト情報と状態へアクセスするための API を提供します。

### リクエストローカル state と共有 state

- ハンドラ内で `Context.State` を直接読み書きしないでください。後方互換のために残っていますが deprecated です。
- そのリクエスト内だけで使う一時値は `Context.Local` に置きます。
- アプリ全体で共有する値は `Context.Set` / `Context.Get` を使います。親アプリを持つ context では app mutex 経由で同期されます。

```go
app.Page("/", func(ctx *mb.Context) mf.Node {
    ctx.Local["trace_id"] = ctx.Query("trace_id") // リクエスト内だけの一時値
    ctx.Set("last_trace_id", ctx.Local["trace_id"]) // アプリ共有 state
    return mf.Text(ctx.Get("last_trace_id").(string))
})
```

### `Param(name string) string`
- `Request.PathValue(name)` からパスパラメータを返します。
- `Request` が `nil` の場合は `""` を返します。

### `Query(name string) string`
- `Request.URL.Query().Get(name)` からクエリパラメータを返します。
- `Request` が `nil` の場合は `""` を返します。

### `FormValue(name string) string`
- `Request.FormValue(name)` からフォーム値を返します。
- `Request` が `nil` の場合は `""` を返します。

### `Asset(name string) string`
- 親アプリからアセット URL を組み立てます。
- 画像、スタイルシート、リンクなどを登録済みアセットプレフィックスに追従させたい場合にハンドラ内で使います。

### `Local map[string]any`
- アプリが作成する context で初期化される、リクエストローカルな一時値用 map です。
- ここに入れた値は他のリクエストと共有されず、app mutex の保護対象でもありません。

### `State map[string]any`
- Deprecated: `Context.Get` / `Context.Set` または `Context.Local` を使ってください。
- 新しいコードでは `Context.State[...]` を直接触らないでください。

### `Set(key string, value any)`
- アプリ共有 state に書き込みます。
- 親アプリを持つ context では app mutex 経由で同期されます。
- 親アプリがない場合は互換性のため deprecated な `Context.State` map に書き込みます。

### `Get(key string) any`
- アプリ共有 state を読み取ります。
- 親アプリを持つ context では app mutex 経由で同期されます。
- 親アプリがない場合は互換性のため deprecated な `Context.State` map から読み取ります。

### `GetInt(key string) int`
- `Get` + `int` アサーションです。
- 値がない、または `int` でない場合は `0` を返します。

### Flash APIs

#### `Flashes() []FlashMessage`
- 現在読み込まれている flash のコピーを返します。
- flash がない場合は `nil` を返します。

#### `FlashSuccess(message string)` / `FlashError(message string)` / `FlashInfo(message string)` / `FlashWarn(message string)`
- `AddFlash(level, message)` の便利ラッパーです。
- level の値は実装定数です:
  - `FlashSuccess`, `FlashError`, `FlashInfo`, `FlashWarn`。

#### `AddFlash(level FlashLevel, message string)`
- メッセージを trim し、空なら no-op です。
- context の flash リストに追加し、Cookie (`marionette_flash`) にシリアライズします。
- Cookie の挙動:
  - `Path=/`
  - `HttpOnly=true`
  - `SameSite=Lax`
  - `Secure` は `App.SetCookieSecure` に従います（デフォルト `false`）。
- シリアライズ失敗は無視されます（panic / status 変更なし）。

次リクエストでの lifecycle:
- flash は Cookie から `Context.flashes` にデコードされます。
- 既知の level かつ空でないメッセージだけが有効です。
- flash が存在した場合、レスポンスで Cookie は自動的に削除されます。

### Session APIs

#### `SetSession(key, value string)`
- `key` を trim し、空なら no-op です。
- context メモリ上の session エントリを保存/更新し、Cookie (`marionette_session`) に書き込みます。
- Cookie の挙動:
  - `Path=/`
  - `HttpOnly=true`
  - `SameSite=Lax`
  - `Secure` は `App.SetCookieSecure` に従います（デフォルト `false`）。

#### `Session(key string) string`
- key に対応する session 値を読み取ります。
- key がない場合は空文字列を返します。

#### `ClearSession()`
- session map を空の map に置き換え、Cookie (`marionette_session`) に書き込みます。

リクエスト時の lifecycle:
- session は `newContext` で Cookie から `Context.session` にデコードされます。
- デコード失敗時は空の session map にフォールバックします（panic / status 変更なし）。

Session sample:

```go
app.Page("/session", func(ctx *mb.Context) mf.Node {
    user := ctx.Session("user")
    if user == "" {
        return mf.Form("/session/login", mf.Button("Sign in"))
    }
    return mf.Form("/session/logout", mf.Button("Sign out"))
})
app.Action("session/login", func(ctx *mb.Context) mf.Node {
    ctx.SetSession("user", "Aiko")
    return mf.Paragraph("Signed in")
})
app.Action("session/logout", func(ctx *mb.Context) mf.Node {
    ctx.ClearSession()
    return mf.Paragraph("Signed out")
})
```

Full example: `docs/site-astro/public/examples/go/session.go`.
