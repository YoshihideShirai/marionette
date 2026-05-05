# APIドキュメント（日本語版）

このディレクトリは API リファレンスの日本語版です。
現在は主要セクションから順次翻訳しています。

## インポートパス

```go
import (
    mb "github.com/YoshihideShirai/marionette/backend"
    md "github.com/YoshihideShirai/marionette/desktop"
    mf "github.com/YoshihideShirai/marionette/frontend"
    mh "github.com/YoshihideShirai/marionette/frontend/html"
)
```

- 推奨エイリアス: `mb`（backend）, `mf`（frontend）。
- `mb`: `New`, `App`, `Context`, `Handler` など実行時 API。
- `md`: `desktop.Run` などデスクトップ実行 API。
- `mf`: `Button`, `Card`, `Table`, `FormRow` などコンポーネント API。
- `mh`: `Node`, `Div`, `Element`, `Raw` など低レベル HTML API。
