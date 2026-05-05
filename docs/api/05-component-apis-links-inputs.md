## 6.1 Links, buttons, and inputs

### Button / SubmitButton
- `Button(label string, props ComponentProps) Node`
- `ButtonWithVariants(label, variant, size, style string, props ComponentProps) Node`
- `SubmitButton(label string, props ComponentProps) Node`
- Extend with `props.Class`, `props.ID`, and `props.Attrs`.

### Link family
- `Link(props LinkProps) Node`
- `ExternalLink(label, href string, props ComponentProps) Node`
- `ExternalIconLink(icon, ariaLabel, href string, props ComponentProps) Node`
- `DownloadLink(label, href, filename string, props ComponentProps) Node`

`Link` behaviors:
- `External: true` defaults to `target="_blank"` + `rel="noopener noreferrer"`.
- `Target: "_blank"` also defaults `rel="noopener noreferrer"`.
- `Download: true` emits `download`.
- `Filename` emits `download="<filename>"`.
- `Props.Disabled: true` renders inert `href="#"` + `aria-disabled`.

### Input / Textarea / Select
- `Input(name, value string, props ComponentProps) Node`
  - Defaults: `Type: "text"`, `Placeholder: strings.TrimSpace(name)`.
- `FileUpload(name string, required bool, props ...ComponentProps) Node`
  - Emits `type="file"`.
- `InputWithOptions(name, value string, options InputOptions) Node`
  - Blank `options.Type` defaults to `"text"`.
  - Common values: `text`, `email`, `password`, `number`, `search`, `url`, `tel`, `date`, `datetime-local`, `time`, `month`, `week`, `file`, `hidden`.
- `Textarea(name, value string, options TextareaOptions) Node`
  - `Rows <= 0` defaults to `3`.
- `Select(name string, options []SelectOption, props ComponentProps) Node`
- `SelectWithVariants(name string, options []SelectOption, props SelectVariantProps) Node`
  - Uses `SelectOption.Value`, `Label`, `Selected`, `Disabled`.

### Form family
- `Form(action string, children ...Node) *form`
- `ActionForm(props ActionFormProps, children ...Node) Node`
  - `Method`: `post` (default), `get`
  - `Target` -> `hx-target`, `Swap` -> `hx-swap`
- `HiddenField(name, value string) Node`
- `FormField(control Node, props FormFieldProps) Node`
