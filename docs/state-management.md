# State Management Guide

Marionette keeps the request/response model explicit: page and action handlers run
on the Go server, htmx swaps HTML fragments in the browser, and `App` global
state is shared by every request handled by that `App` instance. Use the smallest
state scope that satisfies the feature, then move data outward only when it needs
to survive longer or be shared more broadly.

## State categories

| Category | Lifetime and scope | Put here | Avoid putting here |
| --- | --- | --- | --- |
| request-local | One HTTP request/action invocation. It disappears when the handler returns. | Parsed form values, validation errors for the current render, temporary variables, request-scoped service results. | Anything needed after a redirect, across tabs, or by a later request. |
| session/user state | One browser/user session, usually stored as a signed/encrypted cookie or a server-side session keyed by a cookie. | Minimal identity/session pointers, CSRF tokens, short flash messages, locale/theme preference when it is user-specific. | Complete user profiles, permissions snapshots that must stay fresh, large form drafts, business records. |
| app-global state | One Go process and one `App` instance. Shared by all users and all requests in that process. | App configuration, metrics, short-lived in-process caches, demo counters. | Login users, permissions, per-user form input, business data that must be persisted, data that must be shared across multiple app instances. |
| durable state | A database, durable queue, object store, or another system of record. Survives restarts and deploys. | Users, permissions, orders, invoices, audit logs, saved form drafts, workflow state, any data with compliance or recovery requirements. | Ephemeral rendering details, derived values that can be recomputed cheaply, data that must expire quickly without database cleanup. |
| cache/derived state | Recomputable data with an explicit expiration or invalidation rule. Often external when shared across instances. | Expensive query results, API responses with TTLs, rendered summaries, rate-limit counters, cross-instance coordination hints. | The only copy of business data, values whose staleness would break authorization or correctness. |

## Using `App` global state safely

`SetGlobal`, `GetGlobal`, `UpdateGlobal`, `GetGlobalInt`, and
`IncrementGlobalInt` operate on state shared by all users of the current `App`
instance. That makes them useful for process-wide concerns, but unsafe for
per-user or durable data.

Good examples for `App` global state:

- **App configuration** loaded at startup and read by handlers.
- **Metrics** such as counters, gauges, and development diagnostics.
- **Short-lived in-process caches** where losing the value on restart is fine and
  each instance may hold its own copy.
- **Demo counters** or tutorial values where cross-user sharing is intentional.

Do not put these in `App` global state:

- **Logged-in user** or current account identity.
- **Permissions** or authorization decisions.
- **Per-user form input** or drafts.
- **Business data that must be persisted**, such as orders, invoices, workflow
  state, or user-generated records.
- **Data that must be shared across multiple app instances**, because each
  process has its own `App` memory.

When the next value depends on the previous value, prefer `UpdateGlobal` or a
specific helper such as `IncrementGlobalInt` so the read/modify/write sequence is
atomic within the process.

### Mutable values: slices, maps, and pointers

The state mutex protects the `GetGlobal` or `SetGlobal` operation itself, but it
does not make a returned slice, map, or pointer safe to mutate after `GetGlobal`
returns. Treat values returned by `GetGlobal` as read-only references unless the
value has its own synchronization. Follow these rules:

- Do **not** directly modify maps or slices returned by `GetGlobal`.
- Perform every mutation in an `UpdateGlobal` closure so the read/modify/write
  sequence stays under the app lock.
- When rendering or returning a mutable collection to code that might modify it,
  return a clone. Use `GetGlobalSnapshot(key, clone)` to clone while the state
  lock is still held.
- If global state stores a pointer to a mutable object, that object must provide
  its own locking or immutable/snapshot methods. Otherwise store immutable values
  and replace them with `UpdateGlobal`.

Safe `App.UpdateGlobal` examples:

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
// or, for custom counter logic:
count = app.UpdateGlobal("count", func(old any) any {
    current, _ := old.(int)
    return current + 1
}).(int)
```

Safe `Context.UpdateGlobal` examples inside handlers:

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

Use `GetGlobalSnapshot` for read paths that need a collection snapshot:

```go
func cloneStrings(old any) any {
    values, _ := old.([]string)
    return append([]string(nil), values...)
}

messages := ctx.GetGlobalSnapshot("messages", cloneStrings).([]string)
```

## Decision tree

Use this checklist when deciding where a value belongs.

1. **Is the value needed only while rendering the current request or action
   response?**
   - Yes: keep it request-local.
   - No: continue.
2. **Must the value survive process restarts, deploys, or crashes?**
   - Yes: store it in durable state such as a database, durable queue, or object
     store.
   - No: continue.
3. **Is it business data, audit data, permissions, user profile data, or a saved
   draft users expect to recover later?**
   - Yes: store it in the database or another durable system of record.
   - No: continue.
4. **Is it required for authentication/session continuity but not valuable on its
   own?**
   - Yes: put only the minimum pointer in the cookie session, such as a session
     ID, user ID reference, CSRF token, or flash key. Load authoritative user and
     permission data from durable storage on demand.
   - No: continue.
5. **Is the value user-specific but temporary, such as a wizard step or UI
   preference?**
   - Yes: use the smallest session/user state representation. Prefer a compact
     session key or server-side session record for anything large or sensitive.
   - No: continue.
6. **Is the value derived from durable data or an external API and safe to
   recompute?**
   - Yes: use cache/derived state. Use an external cache when the value should be
     shared across app instances or survive a single process restart; otherwise an
     in-process cache can be acceptable.
   - No: continue.
7. **Is cross-instance consistency required?**
   - Yes: use the database, an external cache, or a coordination service instead
     of `App` global state.
   - No: continue.
8. **Is the value process-wide, non-user-specific, and safe to lose on restart?**
   - Yes: `App` global state is appropriate.
   - No: choose durable state or session/user state based on the closest lifetime
     and ownership requirement.

## Storage guidance

### Put state in the database when

- It is the source of truth for a business workflow.
- Users expect it to be available after logout, restart, or deployment.
- It affects authorization, billing, auditability, reporting, or recovery.
- Multiple app instances must read and update it consistently.

### Put only minimal state in the cookie session when

- The browser needs to resume a session on later requests.
- The value is a pointer to authoritative server-side data, not the data itself.
- The value is compact and safe to send with every request.
- Staleness can be resolved by reloading from durable storage.

Typical cookie session contents are a session ID, user ID reference, CSRF token,
post-redirect flash message key, or small preference. Avoid storing full user
objects, permission sets, large drafts, or sensitive business records in cookies.

### Put state in an external cache when

- The value is expensive to compute but can be recomputed from durable data.
- Multiple app instances should share the cached value.
- You need TTL-based expiration, rate limiting, or cross-instance coordination.
- Losing the cache causes slower responses, not data loss or authorization bugs.
