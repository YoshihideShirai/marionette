# Marionette

日本語 | [English](README.md)

![Marionette concept art](docs/assets/concept.png)

Marionette は、管理画面や社内ツールを Go で構築するための Go-first なフレームワークです。
画面、状態、アクションを Go で記述し、ブラウザ側の部分更新は htmx が担います。

## Marionette を使う理由

- Go から離れずに業務向け UI を構築できます。
- ルーティング、状態更新、イベントハンドラをサーバー側に置けます。
- フル SPA を保守せず、htmx による部分レンダリングを使えます。
- ページ、フォーム、アクション、テーブル、チャート、レイアウトコンポーネントから管理画面を組み立てられます。
- 同じアプリを Web UI としても、デスクトップ WebView シェル内でも動かせます。

## 試す

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
