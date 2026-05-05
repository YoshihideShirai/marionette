## 6.1 リンク・ボタン・入力

### Button / SubmitButton
- `Button(label string, props ComponentProps) Node`
- `SubmitButton(label string, props ComponentProps) Node`
- `props.Class`, `props.ID`, `props.Attrs` で拡張できます。

### Link 系
- `Link(props LinkProps) Node`
- `ExternalLink(label, href string, props ComponentProps) Node`
- `ExternalIconLink(icon, ariaLabel, href string, props ComponentProps) Node`
- `DownloadLink(label, href, filename string, props ComponentProps) Node`

`Link` の主要挙動:
- `External: true` で `target="_blank"` + `rel="noopener noreferrer"` を補完。
- `Target: "_blank"` 指定時も `rel="noopener noreferrer"` を補完。
- `Download: true` で `download` 属性を出力。
- `Filename` 指定時は `download="<filename>"`。
- `Props.Disabled: true` では `href="#"` + `aria-disabled` の無効表示。

### Input / Textarea / Select
- `Input(name, value string, props ComponentProps) Node`
  - デフォルト: `Type: "text"`, `Placeholder: strings.TrimSpace(name)`
- `FileUpload(name string, required bool, props ...ComponentProps) Node`
  - `type="file"` を出力。
- `InputWithOptions(name, value string, options InputOptions) Node`
  - `options.Type` 空文字は `"text"`。
  - 代表値: `text`, `email`, `password`, `number`, `search`, `url`, `tel`, `date`, `datetime-local`, `time`, `month`, `week`, `file`, `hidden`。
- `Textarea(name, value string, options TextareaOptions) Node`
  - `Rows <= 0` は `3`。
- `Select(name string, options []SelectOption, props ComponentProps) Node`
- `SelectWithVariants(name string, options []SelectOption, color, size, style string, props ComponentProps) Node`
  - `SelectOption` の `Value`, `Label`, `Selected`, `Disabled` を反映。

### Form 系
- `Form(action string, children ...Node) *form`
- `ActionForm(props ActionFormProps, children ...Node) Node`
  - `Method`: `post`（デフォルト）, `get`
  - `Target` → `hx-target`, `Swap` → `hx-swap`
- `HiddenField(name, value string) Node`
- `FormField(control Node, props FormFieldProps) Node`
