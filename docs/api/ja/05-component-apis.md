## 6. Component APIs

`templates/components/*` と `frontend/ui_impl.go` で提供される主要コンポーネント API の一覧です。

### 入力・フォーム系
- `Button`, `SubmitButton`, `Link`, `ExternalLink`, `DownloadLink`
- `Input`, `InputWithOptions`, `FileUpload`, `Textarea`, `Select`
- `Form`, `ActionForm`, `HiddenField`, `FormField`

### オーバーレイ・フィードバック
- `Modal`, `Toast`, `Alert`, `Skeleton`, `Progress`, `EmptyState`

### データ表示
- `Table`, `Chart`, `Image`, `Pagination`, `Tabs`, `Breadcrumb`
- `Checkbox`, `RadioGroup`, `Switch`, `Badge`, `TextComponent`
- `DataFrame`, `DataFrameChart`, `DataFrameFromCSV`, `DataFrameFromTSV`

### レイアウト・サーフェス
- `Actions`, `Divider`, `Box`, `AppShell`
- `Stack`, `Grid`, `Split`, `PageHeader`, `Container`, `Region`, `Card`, `Section`

詳細な引数やデフォルト挙動は英語版 `docs/api/05-component-apis.md` を参照してください。
