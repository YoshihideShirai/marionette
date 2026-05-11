# Errors / Flash / Feedback Guide

English | [日本語](ja/08-errors-flash-feedback.md)

This document explains how to use errors, flash messages, alerts, toasts, inline errors, and retry paths in Marionette applications. The goal is to show users the next action they can take while leaving enough information in logs for developers to investigate failures.

For validation errors in forms, also see the [Forms / Validation Guide](05-forms-validation.md). Use the [Routing / Pages / Actions Guide](02-routing-pages-actions.md) as the baseline for Action responsibility boundaries.

## Basic policy

Separate errors and feedback by **who reads the information** and **how long it needs to remain visible**.

- User-facing messages should briefly explain what happened, the current state, and the next available action.
- Logged errors should include internal errors, IDs, input summaries, and external API status needed for investigation.
- Validation errors belong near the relevant inputs as inline errors.
- Successful operation results should use flash or toast.
- If users can fix a failure, leave a path to edit input or retry.
- If users cannot fix an internal error, hide the details and make it traceable through logs and trace IDs.

## 1. Separate user-facing errors from logged errors

User-facing and logged errors should use different content, even when they originate from the same `err`.

### User-facing messages

User-facing messages should include:

- What did not complete.
- Which input or condition to check when the user can fix it.
- Whether retrying is possible.
- Which ID to provide when contacting support, if needed.

Avoid showing:

- SQL, stack traces, or raw external API responses.
- Secrets, tokens, personal information beyond what is necessary, or internal host names.
- Type names or package names that only developers understand.

Example:

```go
form.PageError = "Could not save the user. Check the input and try again."
```

### Logged information

Keep investigation details in structured logs.

- Internal error details.
- Request ID, job ID, user ID, and resource ID.
- The service or external API that was called.
- Retry count, timeout, HTTP status, and other technical details.
- Safe identifiers or summaries instead of raw input values when possible.

```go
h.Logger.Error("create user failed",
    "err", err,
    "request_id", requestID,
    "actor_id", actorID,
)
form.PageError = "Could not create the user. Please try again later."
```

When an Action receives a service error, decide what to log first, then convert it into a user-facing representation. Avoid rendering `err.Error()` directly on the screen.

## 2. Choosing alert / toast / inline error

Choose feedback UI based on message importance, scope, and whether the message needs to remain visible.

### inline error

Use inline errors for fixable problems tied to a specific input or row.

- A required field is empty.
- An email address format is invalid.
- A selected value is no longer valid.
- Only one inline-edited row failed to save.

As a rule, show field errors near the relevant field. You can combine them with an alert when the whole form needs a validation summary.

### alert

Use alerts for problems or states that should remain visible in the screen context.

- Save failure, permission denial, or external integration failure.
- Page-level validation summary.
- States that need attention, such as deletion or disabling.
- Long-running task failure or partial success explanation.

Place alerts within the relevant screen context so users can reread them. In forms, alerts work well as page-level errors near the top.

### toast

Use Toast for lightweight result notifications that can disappear after a short time.

- Saved.
- Copied.
- Queued.
- Background task started.

Toast does not interrupt the user's work, but it is hard to reread later. Use an alert for errors that require users to read the reason and fix something, or for important legal, billing, or permission messages.

## 3. Standard usage when a flash API exists

`backend.Context` provides helpers for flash messages. The standard pattern is to call `ctx.FlashSuccess`, `ctx.FlashError`, `ctx.FlashInfo`, or `ctx.FlashWarn` in the Action, then read `ctx.Flashes()` on the Page after redirect.

### Action side

```go
func (h *UserActions) Update(ctx *backend.Context) error {
    form := state.UserFormFromRequest(ctx)
    if !form.Validate() {
        return ctx.HTML(ui.UserForm(form))
    }

    if err := h.Users.Update(ctx.Request.Context(), form.ToInput()); err != nil {
        h.Logger.Error("update user failed", "err", err, "user_id", ctx.Param("id"))
        form.PageError = "Could not update the user. Please try again later."
        return ctx.HTML(ui.UserForm(form))
    }

    ctx.FlashSuccess("User updated")
    return ctx.Redirect("/admin/users/" + ctx.Param("id"))
}
```

### Page / layout side

```go
func UserDetailPage(ctx *backend.Context) frontend.Node {
    flashes := ctx.Flashes()
    return ui.AdminLayout(ui.AdminLayoutProps{
        Flashes: flashes,
        Content: frontend.Div(
            frontend.FlashAlerts(flashes),
            userDetailContent(ctx),
        ),
    })
}
```

Flash is a good fit for a one-time message that should appear after redirect. For validation errors that re-render the form in the same POST response, use form state fields such as `Errors` and `PageError` instead of flash.

## 4. Long-running work and retry paths after failure

Design long-running work with the assumption that the final result is not known immediately after submit.

### When starting long-running work

- Put the button in a loading / disabled state to prevent double submit.
- Tell the user that processing started with toast or flash.
- Show a job ID or link to a progress page when available.
- If the user can leave the screen and check later, provide a path to history, notifications, or a detail page.

```go
job, err := h.Jobs.Enqueue(ctx.Request.Context(), input)
if err != nil {
    h.Logger.Error("enqueue import failed", "err", err)
    ctx.FlashError("Could not start the import. Please try again.")
    return ctx.Redirect("/admin/imports/new")
}

ctx.FlashInfo("Import started")
return ctx.Redirect("/admin/imports/" + job.ID)
```

### Retry paths after failure

After a failure, leave the next action visible on the screen.

- Show a form where the user can fix input and submit again.
- Provide a button to retry with the same conditions.
- Show failed job details, log summaries, and whether retry is allowed.
- For temporary external API failures, explain that the user can retry later.
- If support is needed, show a request ID or job ID.

```go
func ImportJobDetail(job service.ImportJob) frontend.Node {
    children := []frontend.Node{frontend.H2(frontend.Text("Import result"))}
    if job.Failed {
        children = append(children,
            frontend.Alert(frontend.AlertProps{
                Title:       "Import failed",
                Description: "Review the details and retry the import.",
                Props:       frontend.ComponentProps{Variant: "error"},
            }),
            frontend.Link(frontend.LinkProps{
                Label: "Retry",
                Href:  "/admin/imports/" + job.ID + "/retry",
            }),
        )
    }
    return frontend.Section(frontend.SectionProps{}, children...)
}
```

Even when a failure cannot be retried, explain why and provide an alternative action. For example: "This invoice is already finalized, so it cannot be retried. Create a new invoice instead."

## Implementation checklist

- `err.Error()` and stack traces are not shown directly to users.
- User-facing messages briefly explain the state and next action.
- Logs contain structured investigation details such as request IDs, resource IDs, and internal errors.
- Field errors use inline errors, whole-screen or operation-level problems use alerts, and lightweight success messages use toasts.
- Notifications that should appear once after redirect use flash.
- Validation failures return through form state instead of flash.
- Long-running work has start, progress, result, and retry paths.
- Users can tell whether they should fix input, retry, or contact support after a failure.
