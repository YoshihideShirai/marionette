# Marionette

日本語 | [English](README.md)

![Marionette concept art](docs/assets/concept.png)

Marionette は、**管理画面や社内ツール開発を圧倒的にシンプルにする** Go-first フレームワークです。
画面、状態、アクションを Go で一気通貫に記述し、ブラウザ側の高速な部分更新は htmx が担います。
AI コーディング時代にも、この流れをひとつに保つことでクロススタックなコンテキストを圧縮し、境界や説明量を減らしながらフロントエンドの複雑性を抑えられます。

フロントエンドとバックエンドを別々に保守することに疲れたチームに、Marionette は
**実務に適した、運用しやすい UI アーキテクチャ**を提供します。

## Marionette を使う理由

- Go から離れずに、運用 UI を構築できます。
- **AI-friendly な設計**: 運用 UI の流れを Go に集約することで、言語境界、API スキーマの受け渡し、フロントエンド / バックエンド間の同期、状態同期を減らせます。
- ルーティング、状態更新、イベントハンドラをサーバー側に保てます。
- フル SPA を保守せずに、htmx による部分レンダリングでインタラクティブな UI を実現できます。
- ページ、フォーム、アクション、テーブル、チャート、レイアウトコンポーネントから管理画面を組み立てられます。
- チャートとテーブルで同じ `DataQueryState` を共有できるため、ある領域をクリックすると関連ウィジェット全体をまとめて絞り込めます。
- 同じアプリを Web UI とデスクトップ WebView シェルの両方で動かせます。

## AI-friendly なコンテキスト圧縮

Marionette は、未検証の「トークン削減率」を単独で主張するのではなく、コンテキスト圧縮を重視して設計されています。画面、状態遷移、アクションハンドラを Go に寄せることで、バックエンドとフロントエンドの境界を減らし、スキーマ受け渡しやコンテキストスイッチを抑えられます。その結果、AI ツールに変更内容を伝えるときも、複数スタックの連携方法を説明し直すより、プロダクトの流れに集中しやすくなります。構造的な理由は [AI-friendly architecture ガイド](https://yoshihideshirai.github.io/marionette/ja/ai-friendly/) にまとめています。

- **コアアプリに TypeScript ビルドチェーンは不要**: Marionette はアプリケーションロジックを Go 側に寄せ、ブラウザ側の部分更新には htmx を使うため、コアアプリに TypeScript ツールチェーンは必要ありません。ただし、すべてのクライアント JavaScript をなくすという意味ではありません。overlay などの共有ブラウザヘルパーは、表示上の振る舞いを支えるために残る場合があります。新規 `.ts` / `.tsx` ファイルは [UI architecture policy](docs/architecture/ui.md#2-new-typescript-files) で禁止されています。

## こんなチームにおすすめ

- Go バックエンド中心で進めつつ、フロントエンドの保守負荷を下げたい。
- フル SPA スタックに踏み切らず、管理画面・業務画面を短期間でリリースしたい。
- サーバーサイドの監視、アクセス制御、デバッグしやすさを重視したい。
- まずはブラウザ UI として展開し、将来的にはデスクトップシェルでの配布も視野に入れたい。

## まずは 1 分で体験

代表的な管理画面サンプルを起動します:

```bash
go run ./cmd/admin-sample
```

その後、http://127.0.0.1:8082 を開きます。

![Admin sample dashboard](docs/assets/admin-sample.png)

ソース: [`cmd/admin-sample/main.go`](cmd/admin-sample/main.go)

[`robbins23/daisyui-admin-dashboard-template`](https://github.com/robbins23/daisyui-admin-dashboard-template) を参考にした DashWind 風 DaisyUI ダッシュボードデモを起動します:

```bash
go run ./cmd/dashwind-demo
```

その後、http://127.0.0.1:8083 を開きます。

![DashWind デモダッシュボード](docs/assets/dashwind-demo.png)

ソース: [`cmd/dashwind-demo/main.go`](cmd/dashwind-demo/main.go), [`internal/dashwinddemo/app.go`](internal/dashwinddemo/app.go)

最小構成の DashWind セットアップでは、DaisyUI テンプレート、DashWind CSS、ブラウザヘルパーを 1 回の呼び出しで登録できます:

```go
app := mb.New()
dw.Use(app, dw.Options{})
```

フルデモを起動します:

```bash
go run ./cmd/marionette
```

その後、http://127.0.0.1:8080 を開きます。

最小サンプルを起動します:

```bash
go run ./cmd/simple-sample
```

その後、http://127.0.0.1:8081 を開きます。サンプルは [`cmd/simple-sample/main.go`](cmd/simple-sample/main.go) だけで完結します。

AI チャットサンプルデモを起動します（ストリーミング応答はシミュレーションで、外部 API キーは不要です）:

```bash
go run ./cmd/ai-chat-sample
```

その後、http://127.0.0.1:8084 を開きます。サンプルは [`cmd/ai-chat-sample/main.go`](cmd/ai-chat-sample/main.go) だけで完結します。

デスクトップ WebView サンプルを起動します:

```bash
go run -tags marionette_desktop ./cmd/marionette-desktop
```

デスクトップランタイムは、同じ Marionette アプリモデルを localhost サーバーとネイティブ WebView シェルの背後で使います。Linux で desktop tag をビルドするには、GTK 3 と WebKitGTK の開発パッケージが必要です。

## ドキュメント

README は意図的に小さく保っています。チュートリアル、API の詳細、コンポーネント例はドキュメントサイトを参照してください:

- ドキュメントサイト: https://yoshihideshirai.github.io/marionette/ja/
- チュートリアル: https://yoshihideshirai.github.io/marionette/ja/tutorial/
- API ドキュメント: https://yoshihideshirai.github.io/marionette/ja/api/
- コンポーネントギャラリー: https://yoshihideshirai.github.io/marionette/ja/components/
- AI-friendly architecture: https://yoshihideshirai.github.io/marionette/ja/ai-friendly/
- State Management ガイド: [docs/state-management.ja.md](docs/state-management.ja.md)

英語ドキュメントへは、サイト内の言語切り替えから移動できます。

## 開発

[Air](https://github.com/air-verse/air) を使うと、Go ファイルの変更時にデモアプリを自動で再起動できます:

```bash
go install github.com/air-verse/air@latest
air
```

ドキュメントサイトをローカルで起動します:

```bash
cd docs/site-astro
npm install
npm run dev
```

GitHub Pages workflow は `docs/site-astro/` を GitHub Actions 経由で公開します。

## コンポーネントテンプレートの配置規約

- 正式なコンポーネントテンプレート配置先は `templates/components/` です。
- `frontend/components_template_loader_impl.go` は `internal/componenttmpl` 経由でこのディレクトリからコンポーネントテンプレートを解決します。
- テンプレート名は `components/<basename>`（例: `components/link`, `components/button`）とし、`<basename>` は `.tmpl` / `.html` を除いたファイル名にします。

## Heavy Job Template（データアプリ向け）

`cmd/marionette` には、Analytics ページで「Run aggregation」を実行するサンプルフローが含まれています:

- サーバー側の `Job` モデル: 実行 ID、状態、進捗、結果参照を保持します。
- ライフサイクル UX のために UI テンプレートを組み合わせます:
  - 実行中の `progress`
  - ステータスや通知のための `toast`
  - 初回実行前の `empty_state`
- 入力パラメータのハッシュをキーにしたインメモリキャッシュを TTL（3 分）付きで使い、同一条件での再実行を高速化します。

### リトライとタイムアウトのポリシー

- リトライ: 一時的な失敗には自動リトライを 1 回適用します（合計最大 2 回試行）。
- タイムアウト: ジョブの実行予算は 5 秒です。処理がこの予算を超えた場合は失敗として扱います。
- 失敗時の扱い: エラーを toast に表示し、同じ条件または調整した条件でオペレーターが再実行できるようにします。
