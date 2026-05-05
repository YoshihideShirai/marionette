# Marionette

日本語 | [English](README.md)

![Marionette concept art](docs/assets/concept.png)

Marionette は、**管理画面や社内ツール開発を圧倒的にシンプルにする** Go-first フレームワークです。  
画面、状態、アクションを Go で一気通貫に記述し、ブラウザ側の高速な部分更新は htmx が担います。

「フロントとバックを別々に保守して疲れた」──そんなチームに向けた、
**実務最適な UI 開発体験**を提供します。

## Marionette を使う理由

- **Go だけで完結**: ルーティング、状態管理、イベント処理、UI 構築を Go から離れずに実装できます。
- **開発速度が高い**: htmx の部分レンダリングで、重厚な SPA 構成なしにインタラクティブな UI を実現できます。
- **運用しやすい**: サーバー主導の設計で、デバッグ・監視・権限制御を既存の Go 資産と統合しやすくなります。
- **管理画面に必要な部品が揃う**: ページ、フォーム、アクション、テーブル、チャート、レイアウトコンポーネントを標準提供。
- **Web とデスクトップを同時に狙える**: 同じアプリをブラウザでも WebView デスクトップでも動かせます。

## まずは 1 分で体験

フルデモを起動します:

```bash
go run ./cmd/marionette
```

その後、http://127.0.0.1:8080 を開きます。

最小サンプルを起動します:

```bash
go run ./cmd/simple-sample
```

その後、http://127.0.0.1:8081 を開きます。

デスクトップ WebView サンプルを起動します:

```bash
go run -tags marionette_desktop ./cmd/marionette-desktop
```

デスクトップランタイムは、同じ Marionette アプリを localhost サーバーと
ネイティブ WebView シェルで表示します。Linux で desktop tag をビルドするには、
GTK 3 と WebKitGTK の開発パッケージが必要です。

## こんなチームにおすすめ

- Go バックエンド中心で、フロント実装コストを抑えたい
- 管理画面・業務画面を短期間でリリースしたい
- SPA の複雑性を避けつつ、十分にリッチな操作性を実現したい
- Web だけでなく将来的にデスクトップ展開も視野に入れている

## ドキュメント

README は意図的に小さく保っています。チュートリアル、API の詳細、コンポーネント例はドキュメントサイトを参照してください:

- ドキュメントサイト: https://yoshihideshirai.github.io/marionette/ja/
- チュートリアル: https://yoshihideshirai.github.io/marionette/ja/tutorial/
- API ドキュメント: https://yoshihideshirai.github.io/marionette/ja/api/
- コンポーネントギャラリー: https://yoshihideshirai.github.io/marionette/ja/components/

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
- `components.go` と `frontend/components_impl.go` の `loadComponentTemplates` はこのディレクトリのみを読み込みます。
- テンプレート名は `components/<basename>`（例: `components/link`, `components/button`）とし、`<basename>` は `.tmpl` / `.html` を除いたファイル名にします。
