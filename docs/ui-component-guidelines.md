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

### 2.1 Prefix

- Component names must use the `Ui` prefix.
  - Examples: `UiButton`, `UiInput`, `UiModal`

### 2.2 Variant Names

- Name variants by meaning, not by visual appearance.
- Recommended format:
  - `variant`: `primary`, `secondary`, `danger`, `ghost`
  - `size`: `sm`, `md`, `lg`
  - `state`: `default`, `loading`, `disabled`, `error`
- Avoid:
  - Color-based names (`blue`, `red`)
  - Ambiguous abbreviations (`normal2`, `typeA`)

### 2.3 Props Naming

- Use the `is` / `has` prefix for booleans.
  - Examples: `isDisabled`, `isLoading`, `hasIcon`
- Use the `on` prefix for events.
  - Examples: `onClick`, `onClose`, `onChange`
- Use clear nouns for value props.
  - Examples: `label`, `helperText`, `errorMessage`, `ariaLabel`

## 3. Handling Breaking Changes

Breaking changes (API-incompatible changes) must be rolled out in stages.

1. **Deprecation notice**: Mark old APIs as deprecated and provide an alternative.
2. **Deprecation period**: Keep old APIs for at least **2 releases** or **30 days** (whichever is longer).
3. **Migration guidance**: Document what changed, how to replace usage, and include diff examples.
4. **Removal**: Remove old APIs after the deprecation period and record the change in release notes.

### 3.1 Migration Template

- Target: `old props / old variant`
- Replacement: `new props / new variant`
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

## 5. Sample Code and Anti-Patterns

### 5.1 Recommended Samples

```tsx
<UiButton
  variant="primary"
  size="md"
  isLoading={false}
  onClick={handleSubmit}
  ariaLabel="Save"
>
  Save
</UiButton>
```

```tsx
<UiInput
  label="Email address"
  value={email}
  onChange={setEmail}
  isDisabled={false}
  errorMessage={emailError}
/>
```

### 5.2 Prohibited Examples

- Appearance-based variant naming:

```tsx
<UiButton variant="blue">Submit</UiButton>
```

- Unclear prop names:

```tsx
<UiInput typeA={true} normal2="x" />
```

- Missing accessibility labeling:

```tsx
<UiInput value={email} onChange={setEmail} />
```
