# Security / Authorization Guide

English | [日本語](ja/10-security-authz.md)

This document summarizes the basic policy for handling authenticated users and authorization checks in Marionette applications. In admin screens and internal tools, do not rely only on screen-level display controls. Always perform server-side checks in both Pages and Actions.

For related responsibility boundaries, see the [Routing / Pages / Actions Guide](02-routing-pages-actions.md). For where to keep sessions and state, see the [State Management Guide](../state-management.md).

## Basic policy

Design authentication and authorization by separating the following boundaries.

- **Authentication**: Determines who made the request.
- **Session / Context**: Keeps the authenticated user available to handlers.
- **Authorization**: Determines whether that actor may perform an operation on the target resource.
- **UI display control**: Hides or disables buttons and links that the actor cannot use.
- **Audit log**: Keeps a traceable record of actor, action, target, timestamp, and result for important operations.

UI controls are only a user-experience aid; they are not a security boundary. The final allow / deny decision must always happen on the server side in the Page, Action, or service layer.

## 1. Handling authenticated users with Context / session

Restore authenticated user information once per request in middleware, then expose it through `Context`. Avoid implementations where each Page or Action directly re-reads cookies, headers, or the session store.

A recommended flow is:

1. Authentication middleware verifies the session ID or token.
2. The middleware restores actor information from the session store or user repository.
3. The middleware stores actor ID, tenant ID, role, permissions, session ID, and similar data in `Context`.
4. Pages, Actions, and services get the actor from `Context` and pass it to authorization checks.

Keep only the minimum information needed for authorization and display decisions in `Context`. Values such as email address, display name, and role are often useful for screens. Secrets such as password hashes, access tokens, refresh tokens, and external API secrets must not be passed around directly.

Example:

```go
type CurrentUser struct {
    ID        string
    TenantID  string
    Email     string
    Role      string
    SessionID string
}

func CurrentUserFromContext(ctx *backend.Context) (CurrentUser, bool) {
    user, ok := ctx.Local["current_user"].(CurrentUser)
    return user, ok
}
```

For routes that require authentication, stop unauthenticated requests in middleware before they reach a Page or Action. For unauthenticated requests, redirect to login, or return a partial response that indicates login is required when the request is a partial-update request.

## 2. Authorization checks when displaying Pages

A Page should first check whether the screen may be displayed. Avoid loading target data first and only hiding parts of the UI afterward for users who do not have permission.

Check the following when displaying a Page:

- Whether the actor can access the route.
- Whether the actor belongs to the target tenant / organization / project.
- Whether the actor can view the target resource.
- For list pages, whether the query includes only data that the actor is allowed to see.
- For detail pages, whether the ID in the path parameter is inside the actor's authorization scope.

Example:

```go
func (p *UserPages) Detail(ctx *backend.Context) error {
    actor, ok := CurrentUserFromContext(ctx)
    if !ok {
        return ctx.Redirect("/login")
    }

    userID := ctx.Param("id")
    if allowed := p.Authorizer.CanViewUser(ctx.Request.Context(), actor, userID); !allowed {
        ctx.Writer.WriteHeader(http.StatusForbidden)
        return ui.ForbiddenPage()
    }

    user, err := p.Users.GetVisibleUser(ctx.Request.Context(), actor, userID)
    if err != nil {
        return err
    }
    return ctx.HTML(ui.UserDetailPage(user))
}
```

For lists and searches, include authorization constraints in the query from the beginning instead of filtering afterward. This reduces the risk that pagination totals, aggregate values, search suggestions, or empty states reveal the existence of data outside the actor's permissions.

## 3. Authorization checks when executing Actions

An Action must always run authorization checks immediately before execution. Even if the Page showed the button, the user's permissions or the target resource's state may change before submission. HTTP requests can also be sent directly without going through the browser UI.

Check the following when executing an Action:

- Whether the actor is authenticated.
- Whether the actor has the role / permission required for the action.
- Whether the target is inside the actor's tenant / organization / project scope.
- Whether the target's current state accepts the operation.
- Whether the target specified by form, query, or path parameters has been tampered with.
- Whether operation-specific protections such as CSRF protection or idempotency keys are enabled.

Example:

```go
func (a *UserActions) Delete(ctx *backend.Context) error {
    actor, ok := CurrentUserFromContext(ctx)
    if !ok {
        return ctx.Redirect("/login")
    }

    userID := ctx.Param("id")
    if allowed := a.Authorizer.CanDeleteUser(ctx.Request.Context(), actor, userID); !allowed {
        a.Audit.Record(ctx.Request.Context(), audit.Event{
            Actor:  actor.ID,
            Action: "user.delete",
            Target: userID,
            Result: "denied",
        })
        ctx.Writer.WriteHeader(http.StatusForbidden)
        return ui.ForbiddenAlert()
    }

    if err := a.Users.Delete(ctx.Request.Context(), actor, userID); err != nil {
        a.Audit.Record(ctx.Request.Context(), audit.Event{
            Actor:  actor.ID,
            Action: "user.delete",
            Target: userID,
            Result: "failed",
        })
        return err
    }

    a.Audit.Record(ctx.Request.Context(), audit.Event{
        Actor:  actor.ID,
        Action: "user.delete",
        Target: userID,
        Result: "succeeded",
    })
    ctx.FlashSuccess("Deleted the user")
    return ctx.Redirect("/admin/users")
}
```

Treat authorization errors separately from validation errors. Users cannot fix them by editing fields, so use a page-level alert, a 403 response, or a safe redirect instead of field errors.

## 4. Hiding buttons in the UI is not enough

Hiding delete or approval buttons in the UI is useful for reducing mistakes and keeping screens understandable. However, it is not a substitute for authorization checks.

Hiding buttons is insufficient because:

- Users can send requests directly to Action URLs with developer tools, curl, or scripts.
- Roles, permissions, or target state can change after the Page is displayed.
- htmx partial endpoints or JSON endpoints may be called through separate paths.
- Disabled attributes and hidden elements can be changed on the client side.

Use a two-layer approach:

1. **Page / UI**: Adjust visible buttons, disabled states, and explanatory text according to the actor's permissions.
2. **Action / service**: Validate the actor, action, and target on every request, and reject operations that are not allowed.

On the UI side, use disabled buttons, tooltips, alerts, or empty states when needed so users understand why an operation is unavailable. On the Action side, do not trust the UI state; validate every value sent by the client.

## 5. Confirmation UI for dangerous operations

Provide confirmation UI before dangerous operations such as deletion, approval, re-run, deactivation, permission changes, payments, or external sends. Confirmation UI prevents accidental operations, and it is required separately from authorization checks and audit logs.

Show the following in confirmation UI for dangerous operations:

- The action name to execute.
- The target name, ID, and count.
- Whether the operation can be undone.
- The impact range and side effects.
- Whether notifications, external sends, or jobs will be triggered after execution.

Examples by operation type:

- **Deletion**: Show the target name and, if needed, require re-entering the target name.
- **Approval**: Clearly state if approval will publish data, start billing, send notifications, or cause other effects.
- **Re-run**: Explain whether existing results will be overwritten, whether duplicate execution can occur, and whether external APIs will be called again.
- **Permission changes**: Confirm the roles / permissions being granted or removed and the affected users.

Confirmation UI can be implemented with modals, drawers, confirmation pages, and similar patterns. A modal is enough for lightweight operations, but irreversible or high-impact operations should use a dedicated confirmation page that clearly shows target information and warnings.

Even when an Action is submitted from confirmation UI, run server-side validation and authorization again. Do not trust a hidden input as proof that the confirmation screen was passed.

## 6. Information to keep in audit logs

Audit logs let you later answer who performed which action, on which target, when, and with what result. Keep them as searchable structured data separate from ordinary application logs.

At minimum, record the following information.

| Field | Content | Example |
| --- | --- | --- |
| actor | The subject that performed the operation, such as a user ID, service account ID, or session ID. | `user_123` |
| action | The attempted operation. Use a stable name that includes a namespace. | `invoice.approve` |
| target | The operation target. Keeping resource type and ID separately makes searching easier. | `invoice:inv_456` |
| timestamp | The server-side timestamp when the event occurred. | `2026-05-11T10:15:30Z` |
| result | The result, such as success, failure, denial, or cancellation. | `succeeded`, `failed`, `denied` |

Add the following when needed:

- Request ID / trace ID.
- Actor tenant ID / organization ID.
- Client IP and user agent.
- A summary of before / after changes. Do not include secrets or excessive personal information.
- A failure reason category. Use safe values such as `permission_denied`, `validation_failed`, or `conflict` instead of detailed stack traces.
- Idempotency key, job ID, or external request ID.

Record denied and failed operations as well as successful operations. Permission failures, access attempts to resources in another tenant, repeated failures, and cancellations of dangerous operations are especially useful for later investigations.

## Implementation checklist

- Authenticated users are restored in middleware and consistently available through `Context`.
- Secrets are not passed around directly in `Context`.
- Page display checks route, tenant, and target view permissions.
- List and search queries include authorization constraints.
- Action execution re-validates the actor, action, and target combination.
- UI button display control is not treated as a substitute for server-side authorization.
- Dangerous operations such as deletion, approval, and re-run have confirmation UI.
- Audit logs record actor, action, target, timestamp, and result as structured data.
- Tests cover authorization errors, success, failure, and denial paths.
