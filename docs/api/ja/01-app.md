## 2. App

### `New() *App`
- 新しい `*App` を返します。
- 内部状態を初期化します:
  - `state: map[string]any{}`
  - `pages: map[string]Handler{}`
  - `actions: map[string]Handler{}`
  - `cookieSecure: false` (デフォルト)。

### `Run(addr string) error`
- `http.ListenAndServe(addr, a.Handler())` で HTTP サーバーを起動します。
- サーブ開始前に `marionette listening at http://<addr>` を stdout にログ出力します。
- `ListenAndServe` のエラーをそのまま返します。

### デスクトップランタイム: `desktop.Run(app *backend.App, options desktop.Options) error`
- `app.Handler()` を使って private な `127.0.0.1:0` サーバーで app を起動します。
- その local URL を native WebView shell で開きます。
- WebView の終了時に local server をシャットダウンします。
- native WebView adapter を有効にするには `-tags marionette_desktop` 付きでビルドする必要があります。

### デスクトップオプション
- `Title string`: native window のタイトルです。デフォルトは `"Marionette"` です。
- `Width int`: native window の幅です。デフォルトは `1200` です。
- `Height int`: native window の高さです。デフォルトは `800` です。
- `Debug bool`: WebView の debug mode を adapter に渡します。

### ページオプション
- `type PageOptions struct { Title string }`
- `type PageOption func(*PageOptions)`
- `WithTitle(title string) PageOption`
  - 前後空白を除去し、その page route の HTML `<title>` を設定します。

### `SetCookieSecure(secure bool)`
- flash cookie (`marionette_flash`) の `Secure` を有効/無効にします。
- デフォルトは `false` です。
- 次の両方に影響します:
  - flash の書き込み (`Context.AddFlash`)
  - flash の削除 (次の request で flash が消費されるとき)。

### `Assets(prefix string, fsys fs.FS, options ...AssetOption)`
- `prefix` 配下で `fsys` から静的ファイルを配信します。
- local files には `os.DirFS("assets")`、single-binary assets には `embed.FS` と `fs.Sub` を使います。
- directory index response はデフォルトで無効です。directory browsing を意図している場合のみ `WithAssetIndex(true)` を渡してください。
- asset routes は `Handler()` と `Run()` に含まれます。


### `UseAssetProvider(provider assets.AssetProvider)` / `UseAssets(provider assets.AssetProvider)`
- 生成される shell 内で Marionette 組み込み framework/library の CSS と JavaScript URL を解決する provider を置き換えます。
- DaisyUI、Tailwind browser、HTMX、Chart.js の URL を完全に制御したい場合は custom provider を使ってください。
- `UseAssets` は既存の shorthand です。`UseAssetProvider` は同じ挙動の説明的な alias です。

### `UseOfflineAssets(basePath string)`
- CDN URL の代わりに、`basePath` 配下にある Marionette 標準の local vendor file names を使います。デフォルトの logical file names は次のとおりです:
  - `daisyui.css`
  - `tailwindcss-browser.js`
  - `htmx.min.js`
  - `chart.umd.js`
- local directory または embedded filesystem からファイルを配信するため、`Assets` と組み合わせて使います:

```go
package main

import (
    "embed"
    "io/fs"
    "log"
    "time"

    "github.com/YoshihideShirai/marionette/backend"
)

//go:embed frontend/assets/vendor/*
var embeddedAssets embed.FS

func main() {
    app := backend.New()

    vendorFS, err := fs.Sub(embeddedAssets, "frontend/assets/vendor")
    if err != nil {
        log.Fatal(err)
    }

    app.Assets("/vendor", vendorFS, backend.WithAssetCache(24*time.Hour), backend.WithAssetImmutable())
    app.UseOfflineAssets("/vendor")

    // register pages/actions...
    log.Fatal(app.Run("127.0.0.1:8080"))
}
```

### `Downloads(prefix string, fsys fs.FS, options ...AssetOption)`
- `prefix` 配下で `fsys` から `Content-Disposition: attachment` 付きの静的ファイルを配信します。
- `Assets(prefix, fsys, WithAssetDownload(), options...)` と同等です。
- request された file basename を attachment filename として使います。

### `Asset(name string) string`
- 最初に登録された asset prefix から asset URL を組み立てます。
- 例: `app.Assets("/assets", ...)` の後、`app.Asset("hero.jpg")` は `"/assets/hero.jpg"` を返します。
- absolute な `http://`、`https://`、`data:` URL は変更されずそのまま通ります。

### Asset options
- `WithAssetCache(maxAge time.Duration)` は `Cache-Control: public, max-age=<seconds>` を出力します。
- `WithAssetImmutable()` は asset caching が有効なときに `immutable` を追加します。
- `WithAssetIndex(enabled bool)` は underlying file server からの directory index response を許可します。
- `WithAssetDownload()` は一致する files を attachment download として配信します。
- `WithAssetContentTypes(types map[string]string)` は serving 前に extension に基づいて `Content-Type` を設定します。

### `AddStylesheet(href string)`
- full-page HTML shell に custom stylesheet link を追加します。
- 空または空白のみの値は無視されます。
- stylesheets は組み込み Tailwind/daisyUI assets の後に出力されます。

### `AddStyle(css string)`
- full-page HTML shell に trusted inline CSS を追加します。
- 空または空白のみの値は無視されます。
- 小さな app-level override や CSS variables に使います。

### `AddScript(src string)`
- full-page HTML shell に external JavaScript file を追加します。
- 空または空白のみの値は無視されます。
- scripts は Marionette 組み込み JavaScript の後に出力されます。

### `AddJavaScript(js string)`
- full-page HTML shell に trusted inline JavaScript を追加します。
- 空または空白のみの値は無視されます。
- inline JavaScript は custom external scripts の後に出力されるため、`AddScript` で登録した libraries を利用できます。

### `Handler() http.Handler`
- 登録済み routes をすべて含む `*http.ServeMux` を構築して返します。
- `Page` routes:
  - path match は厳密です (`r.URL.Path` が registered path と等しくなければ `404`)。
  - method は `GET` のみです。それ以外は `405 Method Not Allowed` です。
  - full HTML shell を render します。
- `Action` routes:
  - method は `POST` のみです。それ以外は `405 Method Not Allowed` です。
  - `r.ParseForm()` を実行し、parse failure は `400 Bad Request` です。
  - HTML fragment のみを render します。
- `/` が `Page`/`Render` で登録されていない場合:
  - `GET /` は configuration message 付きの `500 Internal Server Error` を返します。
  - root 以外の unmatched paths は `404` です。

### `Page(path string, fn Handler, options ...PageOption)`
- `GET` 用の full-page handler を登録します。
- Path normalization:
  - `""` -> `"/"`
  - 先頭 slash がない場合 -> leading slash が追加されます。
- Render mode: handler の `Node` は Marionette shell HTML で wrap されます。

### `Action(name string, fn Handler)`
- `POST` 用の fragment handler を登録します。
- Name normalization:
  - 常に `"/" + strings.TrimPrefix(name, "/")` として保存されます。
- request form body の parse failure -> `400`。
- Render mode: handler の `Node` は fragment HTML として返されます。

### `Render(fn Handler, options ...PageOption)`
- root page registration の compatibility alias です。
- `Page("/", fn)` と同等です。

### `Handle(name string, fn Handler)`
- action registration の compatibility alias です。
- `Action(name, fn)` と同等です。

### App state helpers

名前に `Global` が含まれる App state helpers は、すべての users とすべての requests で共有される app-wide state を読み書きします。

#### `SetGlobal(key string, value any)`
- lock で保護しながら app shared state map に書き込みます。
- この API は `Global` という名前のため、値は app のすべての users で共有されます。

#### `GetGlobal(key string) any`
- lock で保護しながら app shared state map から読み取ります。
- この API は `Global` という名前のため、値は app のすべての users で共有されます。
- 値が mutable な slice、map、pointer の場合、`GetGlobal` が返した後に直接変更しないでください。`UpdateGlobal` の中で変更するか、clone function 付きの `GetGlobalSnapshot` で読み取ってください。

#### `GetGlobalSnapshot(key string, clone func(any) any) any`
- app read lock を保持したまま app shared state から読み取り、`clone(value)` を返します。
- rendering 前や、値を変更する可能性のある code に渡す前に、slices、maps、その他の mutable values の read-only snapshots を作るために使います。

#### `UpdateGlobal(key string, fn func(old any) any) any`
- app mutex を保持したまま app shared state を atomically に読み取り、変換し、書き込みます。
- counters、progress ticks、append-style updates など、次の値が古い値に依存する場合に使います。

#### `GetGlobalInt(key string) int`
- app shared state を読み取り、`int` に type-assert します。
- 値が missing または `int` でない場合は `0` を返します。

#### `IncrementGlobalInt(key string, delta int) int`
- integer counters/progress values 向けの `UpdateGlobal` convenience wrapper です。
- missing または non-`int` values を `0` として扱い、`old + delta` を保存し、新しい `int` を返します。

#### `Set(key string, value any)` / `Get(key string) any` / `GetInt(key string) int`
- Deprecated: app-wide state にアクセスするときは `SetGlobal` / `GetGlobal` / `GetGlobalInt` / `UpdateGlobal` を使ってください。
- これらの compatibility aliases は、現在もすべての users で共有される同じ app-wide state にアクセスします。
