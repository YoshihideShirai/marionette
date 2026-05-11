# Forms / Validation Guide

English | [日本語](ja/05-forms-validation.md)

This document explains how to design input forms, server-side validation, and Action responses in Marionette applications. Forms often mix input values, errors, save results, and re-render boundaries, so keep the responsibilities of Pages, Actions, state, services, and UI helpers explicit.

For screen and Action boundaries, see the [Routing / Pages / Actions Guide](02-routing-pages-actions.md). For input component details, see [Form APIs](../api/04-form-apis.md).

## Basic policy

For Marionette forms, the standard pattern is to **collect input state in an explicit type, validate it in the Action, and pass the state back to the same render helper**.

- A Page loads initial values, options, and existing records, then creates an empty or prefilled form state.
- An Action reads input from the request, calls server-side validation, and chooses the display response for success or failure.
- State stores input values, field errors, page-level errors, and values that should be shown again after submit.
- Services / domain code handle saving, updating, external API integration, and business-rule validation.
- UI helpers receive only form state and render input values and errors from the same structure.

Client-side `required` attributes and input types are useful aids, but trusted validation must always happen on the server.

## 1. Basic input form structure

Build input forms around the state passed by the Page and the helper that renders that state.

A recommended structure is:

```text
internal/<app>/state/user_form.go   # Input values, errors, Validate
internal/<app>/pages/users.go       # GET display, initial values, option loading
internal/<app>/actions/users.go     # POST submit, validation, save, response choice
internal/<app>/ui/user_form.go      # Helper that builds HTML from form state
internal/<app>/service/users.go     # DB updates, uniqueness checks, business rules
```

At minimum, form state should contain:

- Values entered by the user.
- Error messages for each field.
- Page-level errors that affect the whole form.
- Supporting data needed to re-render controls such as selects or radio groups.
- Target IDs or display mode when a submit should partially update the screen.

Example:

```go
type UserForm struct {
    Name  string
    Email string
    Role  string

    RoleOptions []frontend.SelectOption
    Errors      map[string]string
    PageError   string
}

func NewUserForm(roleOptions []frontend.SelectOption) UserForm {
    return UserForm{
        RoleOptions: roleOptions,
        Errors:      map[string]string{},
    }
}
```

The rendering code should read values such as `Name` and `Email` and errors such as `Errors["email"]` from the same state. This lets initial display, redisplay after validation failure, and edit screens with existing values use the same helper.

## 2. Where to place server-side validation

Separate server-side validation into **input-shape validation** and **business-rule validation**.

### Put these in form state / validators

Put checks that can be decided from request input alone in form state or validators.

- Required fields.
- String length, numeric ranges, and date ranges.
- Basic formats such as email addresses or slugs.
- Whether a select / radio value is in the allowed list.
- Simple consistency checks between two inputs. Example: start date is not after end date.

```go
func (f *UserForm) Validate() bool {
    f.Errors = map[string]string{}

    if strings.TrimSpace(f.Name) == "" {
        f.Errors["name"] = "Enter a name"
    }
    if !strings.Contains(f.Email, "@") {
        f.Errors["email"] = "Enter a valid email address"
    }
    if f.Role == "" {
        f.Errors["role"] = "Select a role"
    }

    return len(f.Errors) == 0
}
```

### Put these in services / domain code

Put checks that require persistent data, external APIs, permissions, or multiple object states in services / domain code.

- Uniqueness checks for email addresses or codes.
- Whether the current user can change the target resource.
- Business rules for contracts, inventory, quotas, workflow states, and similar concepts.
- Consistency that can only be confirmed inside a DB transaction.
- Validation that depends on an external service response.

After `form.Validate()` confirms the input shape, the Action passes the input to a service and converts business errors into user-facing messages.

## 3. Action input and validation flow

Use this standard order in Actions.

1. Read input from `ctx.FormValue`, path parameters, and query parameters.
2. Build form state from the acquired values.
3. Call `Validate`; if there are field errors, re-render the same screen.
4. Pass the input to services / domain code.
5. Classify service errors as user-facing, log-only, or retryable errors.
6. On success, choose redirect, partial update, or toast.

Example:

```go
func UserFormFromRequest(ctx *backend.Context, roles []frontend.SelectOption) UserForm {
    form := NewUserForm(roles)
    form.Name = ctx.FormValue("name")
    form.Email = ctx.FormValue("email")
    form.Role = ctx.FormValue("role")
    return form
}

func (h *UserActions) Create(ctx *backend.Context) error {
    roles := h.Users.RoleOptions(ctx.Request.Context())
    form := state.UserFormFromRequest(ctx, roles)

    if !form.Validate() {
        return ctx.HTML(ui.UserForm(form))
    }

    user, err := h.Users.Create(ctx.Request.Context(), service.CreateUserInput{
        Name:  form.Name,
        Email: form.Email,
        Role:  form.Role,
    })
    if err != nil {
        form.PageError = "Could not create the user. Please try again later."
        h.Logger.Error("create user failed", "err", err)
        return ctx.HTML(ui.UserForm(form))
    }

    ctx.FlashSuccess("User created")
    return ctx.Redirect("/admin/users/" + user.ID)
}
```

Avoid adding SQL, complex HTML, or detailed business decisions directly inside Actions. Keep Actions as orchestration layers that gather input, call validation and persistence, and choose the display result.

## 4. Returning input errors to the same screen

For input errors, preserve the values the user just entered and pass them back to the same form helper. Display field errors near the relevant fields and page-level errors as an Alert near the top.

```go
func UserForm(form state.UserForm) frontend.Node {
    children := []frontend.Node{}
    if form.PageError != "" {
        children = append(children, frontend.Alert(frontend.AlertProps{
            Title:       "Could not save",
            Description: form.PageError,
            Props:       frontend.ComponentProps{Variant: "error"},
        }))
    }

    children = append(children,
        frontend.FormRow(frontend.FormRowProps{
            ID:    "user-name",
            Label: "Name",
            Error: form.Errors["name"],
            Control: frontend.TextField(frontend.TextFieldProps{
                ID:    "user-name",
                Name:  "name",
                Value: form.Name,
            }),
        }),
        frontend.FormRow(frontend.FormRowProps{
            ID:    "user-email",
            Label: "Email address",
            Error: form.Errors["email"],
            Control: frontend.TextField(frontend.TextFieldProps{
                ID:    "user-email",
                Name:  "email",
                Type:  "email",
                Value: form.Email,
            }),
        }),
        frontend.Submit("Save"),
    )
    return frontend.Form("/admin/users", children...)
}
```

When both Pages and Actions call the same helper, initial and error markup stay aligned. Use the error handling built into `FormRow` and input components so components can manage `aria-invalid` and description associations.

The same idea applies to partial updates. The Action re-renders the form, form section, or field group targeted by the update, restoring input values and errors from state.

## 5. Choosing redirect / partial update / toast on success

Choose the success response based on what the user should see next.

### Use redirect when

Use redirect when the operation changes the URL or primary screen context.

- After creating a record, move to its detail screen.
- After editing, return to the list.
- After deleting, leave the deleted detail page.
- The completed state should be reproducible on reload or bookmark.
- You want to avoid resubmitting POST on refresh.

When redirecting, put the one-time message for the next screen in flash with helpers such as `ctx.FlashSuccess`.

### Use partial update when

Use partial update when the URL should stay the same and only part of the screen should refresh.

- Save one card in a settings screen.
- Update only a list row status or detail panel.
- Submit a modal form and replace only the modal body or the background table.
- The htmx target is clear and the update scope is small.

For partial updates, include the success state in the returned HTML or update a separate toast target at the same time.

### Use toast when

Use Toast when the operation succeeded but should not interrupt the user's main task.

- Autosave, small setting changes, or copy completion.
- Successful inline edit save.
- Short success notification without navigation.
- Operations that can be retried from the same place if they fail.

If toast is the only completion signal, the actual data on screen should already be updated. Showing only a toast without updating changed data makes it hard for users to verify the success state.

## Implementation checklist

- Form state contains input values, field errors, and page-level errors.
- Pages and Actions call the same form render helper.
- Validation that only depends on request input lives in form state / validators.
- Validation that depends on DB state, permissions, external APIs, or business rules lives in services / domain code.
- Validation failures preserve input values and re-render the same screen or partial target.
- The success response is chosen using redirect, partial update, and toast criteria.
- User-facing messages are separated from log-only errors.
