## 5. Form APIs

### `FormRow(props FormRowProps) Node`
- Required:
  - `props.ID` must be non-empty; otherwise render error node.
  - `props.Control` must be non-nil; otherwise render error node.
- Behavior:
  - optional label, required marker (`*`), description, error row.
  - when `props.Error` exists, internally renders `FieldError`.

### `FieldError(props FieldErrorProps) Node`
- Empty `Message` => returns empty `Raw("")`.
- Non-empty message requires non-empty `ID`; empty ID => render error node.

### `TextField(props TextFieldProps) Node`
- Defaults:
  - `Type` defaults to `"text"` when blank.
- Behavior:
  - applies `aria-describedby` from description/error presence.
  - sets `aria-invalid=true` when error exists.
  - supports `required`, `disabled`, `readonly`, `data-ref`.

### `Textarea(props TextareaProps) Node`
- `Rows > 0` adds `rows` attribute; otherwise omitted.
- Same accessibility/error behavior as `TextField`.

### `Select(props SelectFieldProps) Node`
- Renders `<option value="...">` list from `Options`.
- `SelectOption.Selected` sets `selected="selected"`.
- Same accessibility/error behavior as text controls.

### `Checkbox(props CheckboxProps) Node`
- Renders label + checkbox input.
- `Checked=true` sets `checked="checked"`.
- Supports `disabled`, `readonly`, `data-ref`, error aria attrs.

### `RadioGroup(props RadioGroupProps) Node`
- Required: non-empty `props.ID`; empty ID => render error node.
- Generates child IDs as `<groupID>-<index>`.
- Marks checked option where `props.Value == option.Value`.

### `Switch(props SwitchProps) Node`
- Implemented as checkbox input with switch class.
- `Checked=true` sets `checked="checked"`.
- Supports same checkable attrs as checkbox/radio.

---
