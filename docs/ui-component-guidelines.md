# UI Component Guidelines

This document standardizes decision criteria and review checkpoints for adding or changing UI components.

## 1. Criteria for Adding a New Component

Only consider adding a new component when all of the following are true:

- **Expected reuse**: The same or similar UI is likely to be reused in **2 or more features**.
- **Screen count**: The component is currently used in **2 or more screens**, or multi-screen adoption is already planned.
- **Complexity**: The implementation is more than a simple HTML fragment and includes at least one of the following:
  - Multiple states (for example: `loading`, `error`, `disabled`)
  - Variations (for example: `size`, `tone`, `emphasis`)
  - Accessibility requirements (keyboard support, ARIA attributes, focus handling)

> If these conditions are not met, prefer extending an existing component or implementing the UI locally in the screen.

## 2. Naming Conventions

### 2.1 Component and Type Names (Current Standard)

- Align names with existing public API patterns.
- The current accepted naming examples are:
  - Component entry points: `Button`, `Input`, `Modal`
  - Props types: `ComponentProps`, `TextFieldProps`, `FormRowProps`
  - Form primitives: `FormRow`, `TextField`
- `Ui` prefix names are **not required** by the current standard.

### 2.2 New API Compatibility Policy (Single Rule)

- For newly added public components, use the `Component*` prefix as the default.
  - Example: `ComponentDatePicker`, `Tabs`
- Keep existing non-`Component*` APIs (for example `FormRow`, `TextField`) as-is for backward compatibility.
- Do not introduce new prefixes in the same layer (`Ui*`, `Base*`, `Core*` etc.) unless an approved RFC explicitly updates this guideline.

#### Review checklist (naming decision criteria)

When reviewing PRs that add or rename APIs, approve naming only when all checks pass:

1. **Consistency**: New public component names follow `Component*`.
2. **Compatibility**: Existing exported names are not renamed without a staged deprecation plan.
3. **Predictability**: Related props are named `*Props` and stay discoverable next to the component.
4. **No mixed policy**: The same API surface does not mix multiple prefix strategies without explicit design approval.

### 2.3 Variant Names

- Name variants by meaning, not by visual appearance.
- Recommended format:
  - `variant`: `primary`, `secondary`, `danger`, `ghost`
  - `size`: `sm`, `md`, `lg`
  - `state`: `default`, `loading`, `disabled`, `error`
- Avoid:
  - Color-based names (`blue`, `red`)
  - Ambiguous abbreviations (`normal2`, `typeA`)

### 2.4 Parameter Naming

- Use the `is` / `has` prefix for boolean parameters.
  - Examples: `isDisabled`, `isLoading`, `hasIcon`
- Use the `on` prefix for event parameters.
  - Examples: `onClick`, `onClose`, `onChange`
- Use clear nouns for value parameters.
  - Examples: `label`, `helperText`, `errorMessage`, `ariaLabel`

### 2.5 Future Direction (Optional)

- If the project decides to standardize on `Ui*` in the future, define it via RFC and migration plan first.
- Until such approval, treat `Ui*` as a future option, not an enforcement rule.

## 3. Handling Breaking Changes

Breaking changes (API-incompatible changes) must be rolled out in stages.

1. **Deprecation notice**: Mark old APIs as deprecated and provide an alternative.
2. **Deprecation period**: Keep old APIs for at least **2 releases** or **30 days** (whichever is longer).
3. **Migration guidance**: Document what changed, how to replace usage, and include diff examples.
4. **Removal**: Remove old APIs after the deprecation period and record the change in release notes.

### 3.1 Migration Template

- Target: `old parameters / old variant`
- Replacement: `new parameters / new variant`
- Replacement rule: whether it is a one-to-one mechanical replacement or requires manual updates
- Impact: expected affected screens/features
- Verification: re-check accessibility, responsiveness, and theme compatibility

## 4. Required Checks

The following checks are mandatory for every component-related PR.

- **Accessibility**
  - Keyboard operable
  - Appropriate ARIA attributes and labels
  - No critical issues in focus visibility or screen-reader output
- **Responsive behavior**
  - No layout breakage at minimum mobile/tablet/desktop sizes
  - No breakage when content size increases
- **Theme support**
  - Acceptable readability in light/dark (or all defined themes)
  - No obvious contrast problems

## 5. Sample Markup and Anti-Patterns

### 5.1 Recommended Samples

```html
<button class="ui-button ui-button--primary ui-button--md" aria-label="Save">
  Save
</button>
```

```html
<label for="email">Email address</label>
<input id="email" type="email" aria-describedby="email-help" />
<p id="email-help" class="ui-helper-text">Use your work email.</p>
```

### 5.2 Good / Bad Naming Examples

- Good (current rule: `Component*`, `*Props`, existing API names remain intact):

```go
button := Button("Save", ComponentProps{Variant: "primary", Size: "sm"})
row := FormRow(FormRowProps{
    Label: "Email",
    Control: TextField(TextFieldProps{Name: "email", Type: "email"}),
})
```

- Bad (mixed or non-standard prefix policy):

```go
button := UiButton("Save", UiButtonProps{Tone: "primary"})
input := BaseTextField(BaseInputOptions{Kind: "email"})
```

- Bad (unclear parameter names):

```text
typeA=true, normal2="x"
```
