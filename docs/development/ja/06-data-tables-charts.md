# データテーブル / チャート設計ガイド

[English](../06-data-tables-charts.md) | 日本語

このドキュメントは、Marionette アプリケーションで一覧テーブル、検索条件、ページング、ソート、チャートを組み合わせる画面の設計方針を説明します。特に、同じ `DataQueryState` をテーブルとチャートで共有し、チャートクリックから一覧を絞り込むような連動 UI を作るときの責務分離を明確にします。

状態全般の lifetime と保存先の考え方は、[状態管理ガイド](../../state-management.ja.md) も参照してください。

## 基本方針

データ一覧画面では、**表示条件を明示的な query state として扱い、データ取得、整形、表示を分ける**ことを標準にします。

- Page は URL query、session/server state、初期値から表示条件を作る。
- Service / repository は検索条件、ページング、ソートを受け取り、DB や外部 API から必要な範囲だけ取得する。
- State は検索条件、ページ番号、ページサイズ、ソート、選択中のチャート filter など、表示を再現するための情報を持つ。
- UI helper は、取得済み rows、columns、pagination metadata、chart data を受け取って HTML を組み立てる。
- Action は検索フォーム、ページ移動、ソート変更、チャートクリックなどの入力を受け取り、次の query state を作って同じ描画 path に戻す。

小さな CSV や in-memory dataframe のサンプルでは UI 側で `DataFrameViewProps` を適用できますが、業務データや件数が多い画面では DB query に条件を渡し、取得量を制限してください。

## 1. 一覧テーブルの基本構成

一覧テーブルは、次の要素を 1 つの画面モデルにまとめると扱いやすくなります。

```text
internal/<app>/state/order_list.go      # Query state, pagination, sort, selected filters
internal/<app>/pages/orders.go          # GET display, URL query parsing, initial data loading
internal/<app>/actions/order_list.go    # search / sort / page / chart-click actions
internal/<app>/service/orders.go        # DB search, count, aggregation, permission checks
internal/<app>/ui/order_list.go         # table, filters, pager, chart layout helpers
```

画面モデルには最低限、次の情報を含めます。

- 検索フォームに再表示する入力値。
- 適用済み filter の一覧。
- 現在のページ、ページサイズ、総件数、次 / 前ページの有無。
- 現在のソート列と方向。
- テーブルの rows と columns。
- 空状態、読み込み失敗、権限不足などを表示するための message。
- 同じ条件で描画するチャート用の labels / datasets / aggregation result。

例:

```go
type OrderListState struct {
    Search   string
    Status   string
    Region   string
    Page     int
    PageSize int
    Sort     string
    SortDesc bool

    TotalRows int
    Rows      []OrderRow
    Chart     frontend.ChartProps
    Error     string
}
```

UI helper では、検索フォーム、テーブル、ページャー、チャートが同じ `OrderListState` を参照するようにします。これにより、初期表示、検索後の再表示、ページ移動、ソート変更、チャート filter 適用後の表示を同じ関数で扱えます。

## 2. 検索条件、ページング、ソートの扱い

### 検索条件

検索条件は「ユーザーが入力した値」と「DB / dataframe に渡す filter」を分けます。

- 入力値は trim、型変換、許可値チェックを行う。
- 空文字や未指定値は filter に変換しない。
- select / radio の値は許可リストと照合する。
- 日付範囲は開始日と終了日の整合性を確認する。
- 権限で固定される条件は、ユーザー入力とは別に service / repository で必ず付与する。

```go
func (s OrderListState) Filters() []frontend.DataFrameFilter {
    filters := []frontend.DataFrameFilter{}
    if strings.TrimSpace(s.Status) != "" {
        filters = append(filters, frontend.DataFrameFilter{
            Column: "status",
            Op:     frontend.DataFrameFilterEq,
            Value:  s.Status,
        })
    }
    if strings.TrimSpace(s.Search) != "" {
        filters = append(filters, frontend.DataFrameFilter{
            Column: "customer_name",
            Op:     frontend.DataFrameFilterContains,
            Value:  s.Search,
        })
    }
    return filters
}
```

### ページング

ページングは必ず上限を決めます。

- `Page` は 1 始まりに正規化する。
- `PageSize` は許可された候補または最大値で clamp する。
- 検索条件が変わったら原則として `Page = 1` に戻す。
- DB では `LIMIT` / `OFFSET` または keyset pagination を使い、全件を UI に渡さない。
- 総件数が高コストな画面では、正確な total count が必要か、次ページ有無だけでよいかを決める。

```go
func normalizePage(page, pageSize int) (int, int) {
    if page < 1 {
        page = 1
    }
    switch pageSize {
    case 20, 50, 100:
    default:
        pageSize = 20
    }
    return page, pageSize
}
```

### ソート

ソートは、画面に表示する label と DB / dataframe の列名を直接結び付けず、許可された sort key に変換します。

- `sort=created_at` や `sort=-total_amount` のように URL に載せやすい表現にする。
- 許可していない列名は default sort に戻す。
- 同値が多い列では安定した順序になるように secondary sort を追加する。
- ソート変更時は `Page = 1` に戻すか、現在ページを維持するかを画面ごとに決める。

```go
var orderSortColumns = map[string]string{
    "created_at":   "orders.created_at",
    "total_amount": "orders.total_amount",
    "status":       "orders.status",
}
```

## 3. `DataQueryState` を使う場合の設計方針

`DataQueryState` は、テーブルとチャートが共有する「データ表示条件」の最小表現として使います。現状の型は filter の集合を中心にしているため、検索フォーム全体、ページング、ソート、権限条件、DB 接続情報まで詰め込まないでください。

使いどころ:

- チャートで選択したカテゴリをテーブルにも適用する。
- 複数 widget が同じ filter を共有する。
- `DataFrameViewProps` に filter を適用して、サンプルや小規模 data を同じ条件で描画する。
- URL query に `df.filter` のような共有 filter を encode して、再読み込み後も同じ view を復元する。

避けること:

- `DataQueryState` を業務検索フォームの唯一の source of truth にする。
- user ID、permission、tenant ID などの security boundary を UI 由来の filter だけに任せる。
- 大量データを全件取得してから `DataQueryState` で絞り込む。
- 表示専用 filter と永続化すべき業務条件を区別せず保存する。

`DataQueryState` を使う画面では、次のように「画面固有 state」と「共有 query state」を分けます。

```go
type SalesDashboardState struct {
    // 画面固有の入力値
    DateFrom string
    DateTo   string
    Region   string

    // テーブル / チャートで共有する表示 filter
    Query frontend.DataQueryState

    Page     int
    PageSize int
    Sort     string

    Rows       []SalesRow
    DataFrame  *rdf.DataFrame
    Chart      frontend.ChartProps
    TotalRows  int
    PageError  string
}
```

描画時は、共有 filter を基本 view に合成します。

```go
baseView := frontend.DataFrameViewProps{
    Page:     state.Page,
    PageSize: state.PageSize,
    Sort: []frontend.DataFrameSort{{
        Column: state.Sort,
        Desc:   strings.HasPrefix(state.Sort, "-"),
    }},
}
view := state.Query.ToView(baseView)
```

DB-backed な画面では、`DataQueryState` から repository input に変換するときに、列名と operator の allowlist を通してください。UI から届いた column 名をそのまま SQL に連結してはいけません。

## 4. チャートクリックでテーブルを絞り込む連動 UI の例

チャートとテーブルを連動させる画面では、**クリックで選ばれた値を query state に追加し、同じ state でチャートとテーブルを再描画する**流れにします。

典型的な flow:

1. Page が URL query から `DataQueryState` と画面固有 state を復元する。
2. Service が同じ条件で table rows と chart aggregation を取得する。
3. Chart に `QueryStateName` と `QueryStateLabel` を設定し、どの filter として扱うかを明示する。
4. ユーザーが chart segment / bar をクリックする。
5. Action または JavaScript が `region=Tokyo` のような選択値を `DataQueryState` に変換する。
6. Page / Action が同じ render helper に戻り、table と chart の両方を更新する。

例:

```go
func BuildSalesView(state SalesDashboardState) frontend.Node {
    queryView := state.Query.ToView(frontend.DataFrameViewProps{
        Page:     state.Page,
        PageSize: state.PageSize,
    })

    chart := state.Chart
    chart.QueryStateName = "sales-filter"
    chart.QueryStateLabel = "region"

    table := frontend.TableProps{
        Columns: []frontend.TableColumn{
            {Label: "Region", SortKey: "region"},
            {Label: "Revenue", SortKey: "revenue"},
        },
        View:           queryView,
        QueryStateName: "sales-filter",
    }

    return frontend.Div(nil,
        frontend.Chart(chart),
        frontend.DataFrame(state.DataFrame, table),
    )
}
```

クリックで追加する filter の例:

```go
func AddChartFilter(state SalesDashboardState, column string, value string) SalesDashboardState {
    next := state
    next.Query.Filters = append(next.Query.Filters, frontend.DataFrameFilter{
        Column: column,
        Op:     frontend.DataFrameFilterEq,
        Value:  value,
    })
    next.Page = 1
    return next
}
```

この設計では、チャートだけが filter を知っている状態を避けられます。テーブル、チャート、summary card、download link が同じ `DataQueryState` から条件を作れるため、画面全体の整合性を保ちやすくなります。

## 5. 状態を URL、server state、DB、cache のどこに置くかの判断

一覧 / チャート画面では、状態の保存先を「復元したい lifetime」と「誰が source of truth か」で決めます。

| 保存先 | 向いている状態 | 避ける状態 |
| --- | --- | --- |
| URL query | 検索条件、ページ、ページサイズ、ソート、選択中の表示 filter | 秘密情報、大きな JSON、改ざんされると危険な権限条件 |
| server state / session | wizard 中の一時条件、URL に出したくない UI preference、短命の draft pointer | 業務 data の source of truth、大きな検索結果そのもの |
| DB | 保存済み view、ユーザーが後で戻る report 定義、audit 対象の業務条件 | クリック直後だけ必要な一時 filter、再計算できる集計 cache |
| cache | 重い集計結果、外部 API response、同じ条件で再利用する derived data | 失うと復元できない data、authorization の根拠、永続化すべき user input |

### URL に置く

次の条件を満たすものは URL query に置くのが第一候補です。

- bookmark / share / reload で同じ一覧を復元したい。
- 値が短く、文字列として表現しやすい。
- ユーザーが見ても問題ない。
- 改ざんされても server-side validation と permission check で安全に拒否できる。

例: `?q=alice&status=open&page=2&page_size=50&sort=-created_at`。

### server state / session に置く

URL に出したくないが短命な UI 状態は server state または session を検討します。

- 複数 step にまたがる一時的な表示条件。
- URL が長くなりすぎる dashboard layout preference。
- post-redirect-get の後に一度だけ表示する flash / selected tab。
- 検索結果そのものではなく、検索条件 draft への短い key。

ただし、session に検索結果全体や大きな dataframe を保存すると、メモリ使用量、stale data、複数 instance 間の不整合が問題になります。

### DB に置く

次の状態は DB などの durable storage に置きます。

- ユーザーが名前を付けて保存する report / saved search。
- deploy、restart、logout 後も残す必要がある dashboard 設定。
- 業務 workflow、承認、請求、監査に影響する条件。
- 複数ユーザーや複数 app instance で共有する定義。

DB に保存した条件を実行するときも、現在の権限、tenant、公開範囲を再評価してください。保存時に安全だった条件が、後で安全とは限りません。

### cache に置く

cache は derived data の高速化に使います。

- 同じ検索条件で重い集計を何度も表示する。
- 外部 API response を短時間だけ再利用する。
- chart aggregation を TTL 付きで保存する。
- cache miss の場合に DB や外部 API から再計算できる。

cache key には、検索条件だけでなく user / tenant / permission scope、集計 version、timezone など、結果に影響する値を含めます。

## 6. 状態管理ガイドへのリンク

より一般的な状態管理の判断基準は、[状態管理ガイド](../../state-management.ja.md) を参照してください。特に、request-local state、session / user state、durable state、cache / derived state、`App` global state の使い分けは、一覧テーブルやチャート連動 UI でも同じ基準で判断します。
