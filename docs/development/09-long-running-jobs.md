# Long-running Jobs Design Guide

English | [日本語](ja/09-long-running-jobs.md)

This guide explains how to handle work that should not be completed within a single HTTP request in Marionette applications, such as data aggregation, CSV imports, external API syncs, and report generation. It expands the **Heavy Job Template (data apps)** from README.md into practical guidance for job models, UI patterns, timeout / retry / cache TTL policies, and production storage decisions.

Long-running jobs should be designed as separate lifecycle steps: **start request, progress inspection, result display, and retry after failure**. The user should always know the current state, whether to wait, fix input, or rerun, and where to inspect the result. The server should keep an execution ID and state transitions that can be used for support and operations.

## Basic policy

- Keep job state on the Go side, just like ordinary Pages and Actions. UI code should render the current job state rather than owning it.
- The start Action should create or enqueue a job, then immediately return to a progress page or a job panel on the originating screen.
- Do not keep the request handler blocked while the heavy work runs. Move the work to a worker, goroutine, queue, or external job runner.
- Use htmx polling or manual refresh to update only the progress region instead of reloading the full page.
- Treat success, failure, timeout, cancel, and retry as explicit state transitions.
- Do not expose internal errors in user-facing messages. Log execution ID, request ID, actor, input hash, and external request IDs for investigation.

## 1. Standard fields for the Job model

A Job is used both for UI display and operational investigation, so keep a standard set of fields. Application-specific input and result types can live in separate structs, but fields shared by lists, detail pages, retry flows, and audit logs should remain on the Job model.

```go
type JobState string

const (
    JobStateQueued    JobState = "queued"
    JobStateRunning   JobState = "running"
    JobStateSucceeded JobState = "succeeded"
    JobStateFailed    JobState = "failed"
    JobStateTimedOut  JobState = "timed_out"
    JobStateCanceled  JobState = "canceled"
)

type Job struct {
    ID             string
    ExecutionID    string
    State          JobState
    Progress       JobProgress
    ResultRef      *JobResultRef
    Error          *JobError
    InputHash      string
    Attempt        int
    MaxAttempts    int
    CreatedAt      time.Time
    StartedAt      *time.Time
    FinishedAt     *time.Time
    ExpiresAt      *time.Time
    RequestedBy    string
}

type JobProgress struct {
    Current int
    Total   int
    Percent int
    Message string
}

type JobResultRef struct {
    Kind string // "table", "csv", "report", "object", "url", etc.
    URI  string
}

type JobError struct {
    Code       string
    Message    string
    Retryable  bool
    DetailRef  string
}
```

### execution ID

The `execution ID` is the correlation ID that ties together user support requests, logs, external API requests, and worker traces. It can be the same as the user-visible job ID, but retries are easier to investigate when you separate the business job from each execution attempt.

- `Job.ID`: the job record for the user operation. It can stay the same across retries if retries are considered the same business operation.
- `ExecutionID`: one concrete execution attempt. Generate a new value for each retry.
- `ParentJobID` or `RetryOf`: add one of these when you need to trace retry ancestry.

Show a short `execution ID` or copyable `job ID` in the UI so users can include it in support requests.

### state

`state` drives display, allowed actions, and worker recovery. Do not let arbitrary strings accumulate; define the allowed states and transitions.

| state | Meaning | Main UI |
| --- | --- | --- |
| `queued` | Accepted and waiting to run | toast / waiting progress |
| `running` | Currently executing | progress, running message, optional cancel |
| `succeeded` | Completed successfully | result link, summary, success toast |
| `failed` | Failed | alert, error summary, retry button |
| `timed_out` | Execution budget exceeded | timeout alert, retry / input adjustment |
| `canceled` | Canceled by a user or administrator | neutral alert, rerun path |

Move from `running` to `succeeded` / `failed` / `timed_out` / `canceled`, and do not silently move a terminal job back to `running`. For retries, prefer creating a new execution instead of rewinding the existing job state; this is easier to audit.

### progress

`progress` helps users decide whether to keep waiting. For measurable work, keep `Current`, `Total`, and `Percent`. For work with no known total, `Message` plus the state may be enough.

- Clamp `Percent` to 0–100.
- If `Total` is unknown, render indeterminate progress.
- Use `Message` for the current step, such as “Aggregating 4,200 of 12,000 rows.”
- Avoid overly frequent updates. Batch them every few hundred milliseconds to a few seconds, or by logical stage, to reduce DB / cache writes.

### result reference

A `result reference` points to where the result can be read; it should not always contain the result itself. Large aggregation outputs and files make job records expensive to list and retry if they are embedded directly.

- Small summaries: may be saved as a result summary on the job record.
- Table output: reference a query state, snapshot ID, or materialized table ID.
- Files: reference an object storage path, an object key used to create a signed URL, or a download route.
- Reports: reference a report ID, dashboard URL, or external BI URL.

If the referenced result can expire, store `ExpiresAt` or a result TTL and show a regeneration path when the result is no longer available.

### error

`error` separates user-facing display from operational investigation.

- `Code`: a category such as `external_timeout`, `invalid_input`, or `rate_limited`.
- `Message`: a short explanation that is safe to show to users.
- `Retryable`: whether the same conditions may be rerun.
- `DetailRef`: an internal reference such as a log trace, error report, or external request ID.

Do not display `err.Error()` or stack traces directly in the UI. Prefer messages that explain the state and next action, such as “The aggregation could not finish because the external API was slow. Try again later.”

## 2. UI patterns for start, progress, success, failure, and retry

Design long-running job screens as a lifecycle, not as a single submit response.

### Starting execution

The start Action validates input, creates the job, prevents duplicate submission, and directs the user to progress.

- Put the run button in loading / disabled state while submitting.
- Build an `InputHash` from the input conditions and use it for cache hits or duplicate detection.
- After enqueue succeeds, show “Job started” with a toast or flash message.
- If there is a progress page, redirect to `/admin/jobs/{id}`.
- If progress is shown inline, return a partial that includes the job ID.

```go
func (h *AnalyticsActions) RunAggregation(ctx *backend.Context) error {
    input, formErr := state.AggregationInputFromRequest(ctx)
    if formErr != nil {
        return ctx.HTML(ui.AggregationForm(input, formErr))
    }

    job, err := h.Jobs.Enqueue(ctx.Request.Context(), input)
    if err != nil {
        h.Logger.Error("enqueue aggregation failed", "err", err, "input_hash", input.Hash())
        ctx.FlashError("Aggregation could not be started. Check the conditions and try again.")
        return ctx.Redirect("/admin/analytics")
    }

    ctx.FlashInfo("Aggregation started")
    return ctx.Redirect("/admin/analytics/jobs/" + job.ID)
}
```

### Showing progress

Keep progress display in a Page or partial endpoint that reads job state. If you use htmx polling, avoid overly short intervals and stop polling after the job reaches a terminal state.

- `queued` / `running`: show progress, current step, execution ID, and start time.
- `succeeded`: stop polling and switch to result links or a summary.
- `failed` / `timed_out`: stop polling and show an alert with retry guidance.
- Provide a job detail URL so users can leave the screen and come back later.

```go
func AggregationJobPanel(job service.Job) frontend.Node {
    switch job.State {
    case service.JobStateQueued, service.JobStateRunning:
        return frontend.Section(frontend.SectionProps{},
            frontend.Progress(frontend.ProgressProps{Value: job.Progress.Percent, Max: 100}),
            frontend.P(frontend.Text(job.Progress.Message)),
            frontend.P(frontend.Text("Execution ID: "+job.ExecutionID)),
        )
    case service.JobStateSucceeded:
        return ui.AggregationResultCard(job.ResultRef)
    default:
        return ui.AggregationFailureCard(job)
    }
}
```

### Success

On success, do not stop at a short success notification. Leave a durable place to inspect the result.

- Show a result summary with cards, stats, or a table.
- Add links to result files or detail pages.
- Show input conditions, duration, generated time, and whether the result came from cache.
- If users may need to regenerate the same result, separate “rerun” from “rerun ignoring cache.”

A success toast is useful when the user is not on the progress page. On the progress page itself, prefer durable in-page UI such as a result card or alert.

### Failure

On failure, leave the cause category and next action on the page.

- Invalid input: return to the form and identify the fields to fix.
- External API timeout / rate limit: explain that retry may work later.
- Authorization failure: do not offer retry; show the required permission or support path.
- Partial success: show succeeded count, failed count, and optionally a failed-row download.
- Support needed: make job ID / execution ID / request ID easy to copy.

If the failure explanation is longer than a sentence, use an alert, details section, or log-summary link instead of a toast.

### Retry

For retry, prefer creating a new execution derived from the original input and context instead of overwriting the same job record.

- Do not show a retry button when `Retryable` is false.
- Show the current input conditions before retry and allow edits when appropriate.
- Use an idempotency key to prevent double-click duplicate executions.
- After retry starts, show the new execution ID and keep a link to the original execution.
- Record automatic retries separately from manual retries.

```go
func (h *AnalyticsActions) RetryAggregation(ctx *backend.Context) error {
    original, err := h.Jobs.Find(ctx.Request.Context(), ctx.Param("id"))
    if err != nil || !original.CanRetry() {
        ctx.FlashError("This aggregation cannot be retried. Check the conditions and start a new run.")
        return ctx.Redirect("/admin/analytics/jobs/" + ctx.Param("id"))
    }

    retried, err := h.Jobs.Retry(ctx.Request.Context(), original.ID)
    if err != nil {
        ctx.FlashError("Retry could not be started. Try again later.")
        return ctx.Redirect("/admin/analytics/jobs/" + original.ID)
    }

    ctx.FlashInfo("Retry started")
    return ctx.Redirect("/admin/analytics/jobs/" + retried.ID)
}
```

## 3. Choosing progress / toast / empty_state

The Heavy Job Template combines `progress`, `toast`, and `empty_state`, but each one has a different role.

| UI | Use for | Avoid for |
| --- | --- | --- |
| `progress` | Running or queued work whose state continues to change | Before the first run, after-result explanation, failure details users must read |
| `toast` | Temporary notifications such as start, completion, and lightweight failures | Errors that need rereading, reasons actions are blocked, important audit information |
| `empty_state` | Before first run, no matching result, expired result | Running jobs, detailed errors, screens that already have a result |

### progress

Use `progress` to explain why the user is waiting. For imports or aggregations with known totals, use determinate progress. When the total is unknown, combine indeterminate loading with a step message.

### toast

Use `toast` for lifecycle milestones.

- Start: “Aggregation started”
- Cache hit: “Showing a result from the same conditions”
- Success: “Aggregation completed”
- Lightweight failure: “Retry could not be started”

However, failure reasons, retry eligibility, and support IDs should also remain visible in the page.

### empty_state

Use `empty_state` to explain why no result is available and what to do next.

- Before first run: “Choose conditions and click Run aggregation.”
- No matching data: “Widen the date range or change filters.”
- Expired result: “The saved result expired. Run the job again.”

Do not mix empty states and failure states. Use alert plus retry guidance for failures; use empty_state for not-yet-run or no-result situations.

## 4. Design criteria for timeout / retry / cache TTL

The README sample uses a 5-second timeout, 1 retry, and a 3-minute cache TTL. Those values are intentionally short for a demo. In production, choose values based on processing cost, external dependencies, user expectations, and operational cost.

### timeout

A timeout is the execution budget for one run, not the amount of time the user is willing to stare at a page. Set separate budgets for HTTP requests, worker processing, external API calls, and DB queries.

- UI request timeout: keep progress partial requests short.
- Enqueue timeout: creating the job should finish within a few seconds.
- Worker execution timeout: define a budget for each processing unit.
- External API timeout: align with the other API's SLA and retry policy.
- DB timeout: set a query budget to prevent accidental full-table scans.

Examples:

| Work | Timeout approach |
| --- | --- |
| Small aggregation / preview | A few seconds. If it exceeds that, move it to a background job. |
| CSV import | Worker timeout based on file size or row count. UI should show progress. |
| External API sync | Use per-API timeouts and rate-limit behavior, plus a separate overall timeout. |
| Report generation | Allow a larger generation budget and store the output as a result reference. |

### retry

Retry only transient failures. Invalid input, missing permissions, and business-rule violations will not become successful just by retrying.

- Limit automatic retries to one or a small number of attempts with backoff.
- Treat external API `429` / `503` / network timeout as retry candidates.
- Do not retry validation errors, authorization errors, or not found errors.
- Record attempt number, execution ID, and error code for each retry.
- Require idempotency keys and duplicate detection for non-idempotent work.

Manual retry should let users confirm the conditions first. If automatic retries are exhausted and you still show a manual retry button, explain why, for example: “This may be a temporary outage. Try again later.”

### cache TTL

Cache TTL is a business decision about how long a result produced from the same inputs may be reused. If it is too short, you recalculate too often. If it is too long, users see stale results.

- Canonicalize input conditions and hash them into the cache key.
- Show `GeneratedAt` and `ExpiresAt` with the result.
- For freshness-sensitive screens, keep TTL short or provide an explicit refresh button.
- For expensive aggregations, use a longer TTL and show cache hits in the UI.
- Include tenant ID / actor scope in the key so cached results do not cross authorization boundaries.

Examples:

| Data type | TTL guideline |
| --- | --- |
| Development demo / sample | A few minutes. Prioritize easy verification. |
| Operational dashboard | Tens of seconds to a few minutes. Balance freshness and load. |
| Daily report | Several hours to 1 day. Show generated time clearly. |
| Billing / audit data | Prefer accuracy and saved snapshots / versions over a simple cache. |

When a user opens an expired result, do not show a blank page. Explain the expiration with empty_state and offer a rerun button.

## 5. When production should use DB or external cache instead of in-memory

The Heavy Job Template uses an in-memory cache as a sample. In production, consider a DB, Redis / Memcached, object storage, and a queue depending on process restarts, multiple instances, worker separation, audit needs, and result size.

### When in-memory is enough

- Development environments, demos, and short-lived previews in a single process.
- Caches that can be lost without business impact.
- Lightweight work that users can quickly rerun.
- Screens that do not require job history or audit logs.

### When to consider a DB

- You need to keep job history, state transitions, and the actor who started the job.
- Support needs to investigate past executions after failures.
- Retry, cancel, timeout, and partial success must be managed reliably.
- Multiple web instances or workers need to read the same job state.
- Progress and result references must survive restarts.
- Tenant / organization permissions must be enforced on job records.

Store job records, attempts, state transitions, input hashes, result references, and error codes in the DB. Do not store large result payloads directly in the job row; reference object storage or aggregated tables instead.

### When to consider an external cache

- The same aggregation result should be shared across multiple instances.
- Cache TTL must survive server restarts.
- Progress updates are frequent and you want to avoid excessive DB writes.
- You need rate limiting or distributed locks.
- Web processes and queue workers need lightweight shared state.

Even with Redis or another external cache, keep the final state of auditable jobs in the DB. Treat cache as the fast current value and DB as the durable history.

### When to consider a queue / object storage

- Job execution should be separated from the web process.
- CPU- or memory-heavy work should run in a worker pool.
- The job generates large CSV, image, PDF, or report files.
- You need retry, dead-letter queue, rate limiting, or priority controls.

Use queues for execution delivery, DB for state and history, object storage for large artifacts, and cache for short-lived progress / result summaries.

## Implementation checklist

- The Job model has execution ID, state, progress, result reference, and error fields.
- State transitions are explicit, and terminal jobs are not moved directly back to running.
- The start Action validates input, prevents duplicate submission, enqueues work, and provides a progress path.
- Progress display uses progress and step messages, and polling stops after completion.
- Success states show result reference, generated time, input conditions, and cache-hit status.
- Failure states keep an alert plus retry / fix / support guidance visible in the page.
- progress, toast, and empty_state are not used interchangeably.
- timeout, retry, and cache TTL are chosen for the business requirements, not copied blindly from demo values.
- Retry is limited to retryable errors and uses idempotency keys to prevent duplicate execution.
- tenant / actor / input hash are included in cache keys or DB records so results do not cross authorization boundaries.
- If production uses in-memory storage, you have confirmed that losing it on restart or across multiple instances is acceptable.
- If audit, multiple instances, worker separation, or large results are required, you have considered DB / external cache / queue / object storage.
