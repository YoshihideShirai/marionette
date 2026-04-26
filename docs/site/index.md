# Marionette Documentation

## Framework Overview

Marionette is a Go-first UI framework for building admin interfaces. Application developers compose screens in Go, while htmx handles partial updates in the browser. The runtime is also designed to work in a WebView host for desktop use cases.

- **Goals (summary)**
  - Define UI using Go only.
  - Compose admin screens from pages, actions, forms, and partial updates.
  - Keep browser-side interaction simple with htmx.
  - Support desktop-oriented deployment via WebView hosting.
- **Architecture (summary)**
  - **Backend (Go):** routing, state updates, and event handling.
  - **Frontend (htmx):** transport and incremental HTML swaps.
  - **Styling:** Tailwind CSS + daisyUI (CDN).
  - **UI DSL (Go):** build Nodes declaratively and render via `html/template`.

Related documentation:

- API reference: [`docs_api.md`](../../docs_api.md)
- UI policy: [`docs/architecture/ui.md`](../architecture/ui.md)
- Components gallery: [`docs/site/components/index.md`](./components/index.md)

---

## Core UI Component List

The following is the fixed list of core Marionette UI components.

1. [ComponentButton](#componentbutton)
2. [ComponentInput](#componentinput)
3. [ComponentFormField](#componentformfield)
4. [ComponentSelect](#componentselect)
5. [ComponentModal](#componentmodal)
6. [ComponentEmptyState](#componentemptystate)
7. [ComponentTable](#componenttable)
8. [ComponentPagination](#componentpagination)

---

## Examples

> Each section includes: (1) Go usage code and (2) a rendered sample preview.

### ComponentButton

**Go code**

```go
saveButton := marionette.ComponentButton("Save", marionette.ComponentProps{
    Variant: "secondary",
    Size:    "sm",
})
```

**Rendered sample**

<div style="margin:8px 0 16px;">
  <button type="button" style="padding:6px 12px;border-radius:8px;border:1px solid #7c3aed;background:#8b5cf6;color:white;font-size:14px;">Save</button>
</div>

### ComponentInput

**Go code**

```go
nameInput := marionette.ComponentInputWithOptions("name", "", marionette.InputOptions{
    Type:        "text",
    Placeholder: "Username",
    Required:    true,
    Props:       marionette.ComponentProps{Size: "lg"},
})
```

**Rendered sample**

<div style="margin:8px 0 16px;max-width:360px;">
  <input type="text" placeholder="Username" style="width:100%;padding:10px 12px;border:1px solid #cbd5e1;border-radius:8px;font-size:14px;" />
</div>

### ComponentFormField

**Go code**

```go
field := marionette.ComponentFormField(
    marionette.ComponentInput("email", "", marionette.ComponentProps{}),
    marionette.FormFieldProps{
        Label:    "Email address",
        Required: true,
        Hint:     "Used for notifications",
    },
)
```

**Rendered sample**

<div style="margin:8px 0 16px;max-width:420px;display:grid;gap:6px;">
  <label style="font-weight:600;font-size:14px;">Email address <span style="color:#dc2626;">*</span></label>
  <input type="email" placeholder="name@example.com" style="padding:10px 12px;border:1px solid #cbd5e1;border-radius:8px;font-size:14px;" />
  <small style="color:#64748b;">Used for notifications</small>
</div>

### ComponentSelect

**Go code**

```go
statusSelect := marionette.ComponentSelect("status", []marionette.SelectOption{
    {Label: "Active", Value: "active", Selected: true},
    {Label: "Inactive", Value: "inactive"},
}, marionette.ComponentProps{})
```

**Rendered sample**

<div style="margin:8px 0 16px;max-width:240px;">
  <select style="width:100%;padding:10px 12px;border:1px solid #cbd5e1;border-radius:8px;font-size:14px;">
    <option selected>Active</option>
    <option>Inactive</option>
  </select>
</div>

### ComponentModal

**Go code**

```go
confirmModal := marionette.ComponentModal(marionette.ModalProps{
    Title: "Confirm deletion",
    Body:  marionette.Text("This action cannot be undone."),
    Actions: marionette.ComponentButton("Close", marionette.ComponentProps{
        Variant: "ghost",
    }),
    Open: true,
})
```

**Rendered sample**

<div style="margin:8px 0 16px;max-width:460px;border:1px solid #e2e8f0;border-radius:12px;padding:16px;background:#fff;box-shadow:0 8px 24px rgba(15,23,42,.08);">
  <div style="font-weight:700;font-size:18px;">Confirm deletion</div>
  <p style="margin:8px 0 14px;color:#475569;">This action cannot be undone.</p>
  <div style="display:flex;gap:8px;justify-content:flex-end;">
    <button type="button" style="padding:6px 12px;border:1px solid #cbd5e1;border-radius:8px;background:white;">Close</button>
    <button type="button" style="padding:6px 12px;border:1px solid #b91c1c;border-radius:8px;background:#dc2626;color:white;">Delete</button>
  </div>
</div>

### ComponentEmptyState

**Go code**

```go
empty := marionette.ComponentEmptyState(marionette.EmptyStateProps{
    Title:       "No data available",
    Description: "Try changing your filter criteria.",
    Skeleton:    false,
})
```

**Rendered sample**

<section style="margin:8px 0 16px;padding:24px;border:1px solid #e2e8f0;border-radius:12px;text-align:center;max-width:480px;background:#f8fafc;">
  <h4 style="margin:0 0 8px;font-size:18px;">No data available</h4>
  <p style="margin:0;color:#64748b;">Try changing your filter criteria.</p>
</section>

### ComponentTable

**Go code**

```go
table := marionette.ComponentTable(marionette.TableProps{
    Columns: []marionette.TableColumn{
        {Label: "Name"},
        {Label: "Role"},
    },
    Rows: []marionette.TableComponentRow{
        {Cells: []marionette.Node{marionette.Text("Alice"), marionette.Text("Admin")}},
    },
    EmptyTitle:       "No users",
    EmptyDescription: "Create your first user.",
})
```

**Rendered sample**

<table style="margin:8px 0 16px;border-collapse:collapse;min-width:320px;">
  <thead>
    <tr>
      <th style="text-align:left;padding:8px 10px;border-bottom:1px solid #cbd5e1;">Name</th>
      <th style="text-align:left;padding:8px 10px;border-bottom:1px solid #cbd5e1;">Role</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td style="padding:8px 10px;border-bottom:1px solid #e2e8f0;">Alice</td>
      <td style="padding:8px 10px;border-bottom:1px solid #e2e8f0;">Admin</td>
    </tr>
  </tbody>
</table>

### ComponentPagination

**Go code**

```go
pager := marionette.ComponentPagination(marionette.PaginationProps{
    Page:       2,
    TotalPages: 10,
    PrevHref:   "?page=1",
    NextHref:   "?page=3",
})
```

**Rendered sample**

<nav aria-label="Pagination" style="margin:8px 0 16px;display:flex;gap:8px;align-items:center;">
  <a href="?page=1" style="padding:6px 10px;border:1px solid #cbd5e1;border-radius:8px;text-decoration:none;color:#0f172a;">Previous</a>
  <span style="padding:6px 10px;border:1px solid #2563eb;border-radius:8px;background:#eff6ff;color:#1d4ed8;">2 / 10</span>
  <a href="?page=3" style="padding:6px 10px;border:1px solid #cbd5e1;border-radius:8px;text-decoration:none;color:#0f172a;">Next</a>
</nav>
