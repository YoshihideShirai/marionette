# State Management ガイド

Marionette では request/response モデルを明示的に扱います。ページハンドラと action ハンドラは Go サーバー上で動き、htmx がブラウザ側で HTML fragment を差し替えます。また、`App` の global state は、その `App` instance が処理するすべての request で共有されます。まずは機能を満たす最小の state scope を選び、より長く保持する必要や、より広く共有する必要が出たときだけ外側の保存先へ移してください。

## State の分類

| 分類 | 寿命と scope | 置くもの | 置かないもの |
| --- | --- | --- | --- |
| request-local | 1 回の HTTP request/action 呼び出し内。ハンドラが返ると消えます。 | parse 済み form 値、現在の render 用 validation error、一時変数、request scoped な service 結果。 | redirect 後、別 tab、後続 request で必要なもの。 |
| session/user state | 1 browser/user session。多くの場合、署名/暗号化 cookie または cookie key に紐づく server-side session に置きます。 | 最小限の identity/session pointer、CSRF token、短い flash message、user-specific な locale/theme 設定。 | 完全な user profile、常に最新であるべき permission snapshot、大きな form draft、業務 record。 |
| app-global state | 1 Go process かつ 1 `App` instance。process 内の全 user・全 request で共有されます。 | app 設定、metrics、短命な process 内 cache、demo 用 counter。 | ログイン user、権限、user 別 form 入力、永続化が必要な業務 data、複数 instance 間で共有すべき data。 |
| durable state | database、durable queue、object store など system of record。restart や deploy をまたいで残ります。 | user、permission、order、invoice、audit log、保存済み form draft、workflow state、compliance/recovery が必要な data。 | 一時的な render 詳細、安価に再計算できる derived value、database cleanup なしで短時間 expire させたい data。 |
| cache/derived state | 明示的な expire/invalidation rule を持つ再計算可能 data。instance 間で共有する場合は外部 cache が多いです。 | 重い query 結果、TTL 付き API response、render 済み summary、rate-limit counter、cross-instance coordination hint。 | 業務 data の唯一の copy、stale になると authorization や正しさが壊れる値。 |

## `App` global state を安全に使う

`SetGlobal`、`GetGlobal`、`UpdateGlobal`、`GetGlobalInt`、`IncrementGlobalInt` は、現在の `App` instance の全 user で共有される state を操作します。そのため process-wide な関心事には便利ですが、user 別 data や durable data には向きません。

`App` global state に置いてよい例:

- 起動時に読み込み、handler から参照する **app 設定**。
- counter、gauge、開発用 diagnostic などの **metrics**。
- restart で失われてもよく、instance ごとに copy を持ってよい **短命な process 内 cache**。
- cross-user で共有されることが意図された **demo 用 counter** や tutorial 用の値。

`App` global state に置くべきでない例:

- **ログイン user** や現在の account identity。
- **権限** や authorization decision。
- **user 別 form 入力** や draft。
- order、invoice、workflow state、user-generated record など、**永続化が必要な業務 data**。
- 各 process が別々の `App` memory を持つため、**複数 instance 間で共有すべき data**。

新しい値が古い値に依存する場合は、process 内で read/modify/write を atomic に行うため、`UpdateGlobal` または `IncrementGlobalInt` のような専用 helper を優先してください。

### mutable な値: slice、map、pointer

state の mutex が保護するのは `GetGlobal` や `SetGlobal` の操作そのものです。
`GetGlobal` が返した slice、map、pointer をその後に変更しても、その変更は lock
の外で行われます。返ってきた値は、その値自身が同期機構を持っている場合を除き、
読み取り専用の reference として扱ってください。ルールは次のとおりです。

- `GetGlobal` で取り出した map/slice を直接変更しない。
- 変更は必ず `UpdateGlobal` の closure 内で行い、read/modify/write を app lock の内側に収める。
- render や後続処理に mutable collection を渡す必要がある場合は、clone してから返す。state lock を保持したまま clone するには `GetGlobalSnapshot(key, clone)` を使う。
- global state に mutable object への pointer を置く場合、その object 自身が lock または immutable/snapshot method を持つ必要がある。そうでない場合は immutable な値を保存し、`UpdateGlobal` で値ごと差し替える。

安全な `App.UpdateGlobal` の例:

```go
app.UpdateGlobal("messages", func(old any) any {
    messages, _ := old.([]string)
    next := append([]string(nil), messages...)
    return append(next, "new message")
})

app.UpdateGlobal("labels", func(old any) any {
    labels, _ := old.(map[string]string)
    next := make(map[string]string, len(labels)+1)
    for key, value := range labels {
        next[key] = value
    }
    next["status"] = "ready"
    return next
})

count := app.IncrementGlobalInt("count", 1)
// 独自の counter logic が必要な場合:
count = app.UpdateGlobal("count", func(old any) any {
    current, _ := old.(int)
    return current + 1
}).(int)
```

handler 内での安全な `Context.UpdateGlobal` の例:

```go
ctx.UpdateGlobal("messages", func(old any) any {
    messages, _ := old.([]string)
    next := append([]string(nil), messages...)
    return append(next, ctx.FormValue("message"))
})

ctx.UpdateGlobal("labels", func(old any) any {
    labels, _ := old.(map[string]string)
    next := make(map[string]string, len(labels)+1)
    for key, value := range labels {
        next[key] = value
    }
    next["last_user"] = ctx.FormValue("name")
    return next
})

count := ctx.IncrementGlobalInt("count", 1)
```

collection の snapshot を読む path では `GetGlobalSnapshot` を使います。

```go
func cloneStrings(old any) any {
    values, _ := old.([]string)
    return append([]string(nil), values...)
}

messages := ctx.GetGlobalSnapshot("messages", cloneStrings).([]string)
```

## Decision tree

値をどこに置くべきか迷ったら、次の順に判断します。

1. **現在の request または action response を render する間だけ必要か?**
   - はい: request-local に置きます。
   - いいえ: 次へ。
2. **process restart、deploy、crash をまたいで残す必要があるか?**
   - はい: database、durable queue、object store などの durable state に置きます。
   - いいえ: 次へ。
3. **業務 data、audit data、permission、user profile、後で復元できるべき saved draft か?**
   - はい: database または system of record になる durable storage に置きます。
   - いいえ: 次へ。
4. **authentication/session 継続に必要だが、それ自体に業務価値はないか?**
   - はい: cookie session には session ID、user ID reference、CSRF token、flash key など最小限の pointer だけを置きます。authoritative な user/permission data は必要に応じて durable storage から読みます。
   - いいえ: 次へ。
5. **wizard step や UI preference など、user-specific だが一時的な値か?**
   - はい: 最小の session/user state 表現を使います。大きい値や sensitive な値は、cookie に直接入れず server-side session record を優先します。
   - いいえ: 次へ。
6. **durable data や外部 API から導出でき、再計算しても安全か?**
   - はい: cache/derived state を使います。app instance 間で共有したい、または単一 process restart をまたぎたい場合は外部 cache を使い、それ以外は in-process cache も選択肢になります。
   - いいえ: 次へ。
7. **cross-instance consistency が必要か?**
   - はい: `App` global state ではなく、database、外部 cache、coordination service を使います。
   - いいえ: 次へ。
8. **process-wide、non-user-specific、restart で失ってよい値か?**
   - はい: `App` global state が適しています。
   - いいえ: 最も近い lifetime と ownership に基づいて durable state または session/user state を選びます。

## 保存先別の指針

### DB に置くべき状態

- 業務 workflow の source of truth である。
- user が logout、restart、deploy 後も利用できることを期待する。
- authorization、billing、auditability、reporting、recovery に影響する。
- 複数 app instance から一貫して読み書きする必要がある。

### cookie session に置くべき最小状態

- browser が後続 request で session を再開するために必要である。
- 値そのものではなく、authoritative な server-side data への pointer である。
- 小さく、すべての request に付いて送信されても問題ない。
- stale になっても durable storage から読み直して解決できる。

典型例は session ID、user ID reference、CSRF token、post-redirect flash message key、小さな preference です。完全な user object、permission set、大きな draft、sensitive な業務 record は cookie に保存しないでください。

### 外部 cache に置くべき状態

- 計算は重いが、durable data から再計算できる。
- 複数 app instance で同じ cache value を共有したい。
- TTL expire、rate limiting、cross-instance coordination が必要である。
- cache を失っても response が遅くなるだけで、data loss や authorization bug にならない。
