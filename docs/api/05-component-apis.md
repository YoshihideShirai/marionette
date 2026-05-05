## 6. Component APIs

Template-backed component constructors (`templates/components/*`).

### Buttons / inputs / field wrappers
- `Button(label string, props ComponentProps) Node`
- `SubmitButton(label string, props ComponentProps) Node`
- `Link(props LinkProps) Node`
  - renders an anchor for internal, external, and download links.
  - `Icon` renders an `aria-hidden` icon before the label.
  - `External` defaults `target="_blank"` and `rel="noopener noreferrer"`.
  - `Target: "_blank"` also defaults `rel="noopener noreferrer"`.
  - `Download` emits `download`; `Filename` emits `download="<filename>"`.
  - `Props.Disabled` renders an inert `href="#"` link with `aria-disabled`.
- `ExternalLink(label, href string, props ComponentProps) Node`
- `ExternalIconLink(icon, ariaLabel, href string, props ComponentProps) Node`
- `DownloadLink(label, href, filename string, props ComponentProps) Node`
- `Input(name, value string, props ComponentProps) Node`
  - uses `InputWithOptions` with defaults:
    - `Type: "text"`
    - `Placeholder: strings.TrimSpace(name)`.
- `FileUpload(name string, required bool, props ...ComponentProps) Node`
  - renders a file input via `InputWithOptions` with `Type: "file"`.
- `InputWithOptions(name, value string, options InputOptions) Node`
  - blank `options.Type` defaults to `"text"`.
- `Textarea(name, value string, options TextareaOptions) Node`
  - `Rows <= 0` defaults to `3`.
- `Form(action string, children ...Node) *form`
  - renders `<form>` with `ID`, `Class`, `Method`, `Action`, and passthrough `Attrs`.
- `ActionForm(props ActionFormProps, children ...Node) Node`
  - renders a form wired to Marionette/HTMX action updates.
  - blank `Method` defaults to `post`; supported methods are `post` and `get`.
  - renders standard `action`/`method` attributes plus `hx-post` or `hx-get`.
  - optional `Target` and `Swap` render `hx-target` and `hx-swap`.
- `HiddenField(name, value string) Node`
  - renders a hidden form field.
- `FormField(control Node, props FormFieldProps) Node`
  - if `control` rendering fails, returns render error node.
- `Select(name string, options []SelectOption, props ComponentProps) Node`

### Overlay / feedback
- `Modal(props ModalProps) Node`
  - renders `Body` and `Actions` nodes first.
  - if either render fails, returns render error node.
- `Toast(props ToastProps) Node`
  - blank `Live` defaults to `"polite"`.
- `Alert(props AlertProps) Node`
- `Skeleton(props SkeletonProps) Node`
  - `Rows <= 0` defaults to `3`.
- `Progress(props ProgressProps) Node`
  - renders a native `<progress>` element with optional label and percentage text.
  - `Max <= 0` defaults to `100`; `Value` is clamped into `0..Max`.
  - `Indeterminate` omits the `value` attribute so browsers render an indeterminate indicator.
- `EmptyState(props EmptyStateProps) Node`
  - `Rows <= 0` defaults to `3`.

### Data display
- `TableRowValues(values ...any) TableRow`
  - converts `nil` to empty text, `Node` values directly, and other values with `fmt.Sprint`.
- `Table(props TableProps) Node`
  - renders each cell node; any cell render error => render error node.
- `Chart(props ChartProps) Node`
  - renders a Chart.js-backed chart from Go props.
  - blank `Type` defaults to `ChartTypeLine`; blank `Height` defaults to `320`.
  - `ChartDataset.Data` renders scalar values; `ChartDataset.Points` renders `{x,y}` values for scatter-style charts.
  - chart config is JSON-encoded and embedded next to a `<canvas data-mrn-chart>`.
  - includes `role="img"`, an accessible label, canvas fallback text, and a screen-reader fallback table.
- `Image(props ImageProps) Node`
  - renders a responsive `<figure>` with an `<img>` and optional caption.
  - `Src` is required; blank `Src` returns a render error node.
  - blank `Loading` defaults to `"lazy"` and blank `Decoding` defaults to `"async"`.
  - `AspectRatio` supports `square`, `video`, `wide`, and `portrait`; `ObjectFit` supports `cover`/blank, `contain`, `fill`, `none`, and `scale-down`.
- `DataFrame(df *dataframe.DataFrame, props TableProps) Node`
  - renders `github.com/rocketlaunchr/dataframe-go` dataframes through `Table`.
  - `df.Names()` is mapped to `TableColumn.Label` and overrides `props.Columns`.
  - each row is read by `df.Row(row, true, dataframe.SeriesName)`.
  - cell conversion: `nil` => empty text, `Node` => rendered directly, all others => `fmt.Sprint(value)`.
- `DataFrameChart(df *dataframe.DataFrame, props DataFrameChartProps) Node`
  - maps a dataframe label column and numeric series columns into `Chart`.
  - blank `LabelColumn` uses the first dataframe column.
  - blank `Series` renders every column after the label column as a dataset.
- `Pagination(props PaginationProps) Node`
  - `Page < 1` defaults to `1`.
  - `TotalPages < 1` defaults to `1`.
- `Tabs(props TabsProps) Node`
  - blank `AriaLabel` defaults to `"tabs"`.
  - supports active/disabled states and link/button tab items.
- `Breadcrumb(props BreadcrumbProps) Node`
  - blank `AriaLabel` defaults to `"breadcrumb"`.
  - supports active/current breadcrumb items.
- `Checkbox(props CheckboxComponentProps) Node`
- `RadioGroup(props RadioGroupComponentProps) Node`
  - blank `AriaLabel` defaults to `"radio group"`.
- `Switch(props SwitchComponentProps) Node`
- `Badge(props BadgeProps) Node`
  - renders a compact label with `Variant`, `Size`, and custom classes from `ComponentProps`.
- `TextComponent(props TextProps) Node`
  - renders plain text with semantic size, weight, and tone options.
- `DataFrameFromCSV(r io.ReadSeeker, props TableProps, opts ...imports.CSVLoadOptions) (Node, error)`
  - loads CSV via `github.com/rocketlaunchr/dataframe-go/imports.LoadFromCSV`.
- `DataFrameFromTSV(r io.ReadSeeker, props TableProps, opts ...imports.CSVLoadOptions) (Node, error)`
  - same loader with `Comma: '\t'` as default.

### Layout / surfaces
- `Actions(props ActionsProps, children ...Node) Node`
  - renders a horizontal action group.
  - `Align`: `start`/blank, `center`, `end`, `between`.
  - `Gap` uses the same values as `Stack`; `Wrap` adds `flex-wrap`.
- `Divider(props DividerProps) Node`
  - renders a visual divider.
  - `Spacing`: `none`, `xs`, `sm`, `md`/blank, `lg`.
- `Box(props BoxProps, children ...Node) Node`
  - renders a generic surface with optional border, tone, padding, and custom classes.
- `AppShell(props AppShellProps) Node`
  - renders the demo/admin shell with sidebar, flashes, header, and content regions.
- `Stack(props StackProps, children ...Node) Node`
  - flex layout for vertical/horizontal stacks.
  - `Direction`: `vertical`/blank or `horizontal`/`row`.
  - `Gap`: `none`, `xs`, `sm`, `md`/blank, `lg`, `xl`.
  - `Align`: `start`, `center`, `end`, blank=`stretch`.
  - `Justify`: `start`/blank, `center`, `end`, `between`.
  - `Wrap` adds `flex-wrap`; `Props.Class` appends custom classes.
- `Grid(props GridProps, children ...Node) Node`
  - grid layout with `Columns` values `1`, `2`, `3`/blank, `4`.
  - `MinColumnWidth`: `sm`, `md`, `lg` switches to auto-fit minmax columns.
  - `Gap` and `Props.Class` use the same behavior as `Stack`.
- `Split(props SplitProps) Node`
  - responsive main/aside layout.
  - `AsideWidth`: `sm`, `md`/blank, `lg`.
  - `ReverseOnMobile` renders the aside before the main pane visually on mobile.
- `PageHeader(props PageHeaderProps) Node`
  - renders title, description, and optional action node.
- `Container(props ContainerProps, children ...Node) Node`
  - `MaxWidth`: `sm`, `md`, `lg`/blank, `full`.
  - `Padding`: `none`, `sm`, `md`/blank, `lg`.
  - `Centered` adds `mx-auto`.
- `Region(props RegionProps, children ...Node) Node`
  - renders an ID-addressable content region for partial updates.
  - `ID` is required; blank `ID` returns a render error node.
  - `Props.Class` appends custom classes.
- `Card(props CardProps, children ...Node) Node`
  - card surface with optional title, description, and action node.
  - `Gap` controls spacing between card body children.
- `Section(props SectionProps, children ...Node) Node`
  - unframed section wrapper with optional title, description, and action node.

#### Example: Convert CSV/TSV data to `DataFrame`

```go
import (
    "os"

    marionette "github.com/YoshihideShirai/marionette"
)

func tableFromCSV(path string) (mrn.Node, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    return mrn.DataFrameFromCSV(f, mrn.TableProps{
        EmptyTitle:       "No data",
        EmptyDescription: "CSV is empty.",
    })
}

func tableFromTSV(path string) (mrn.Node, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    return mrn.DataFrameFromTSV(f, mrn.TableProps{
        EmptyTitle:       "No data",
        EmptyDescription: "TSV is empty.",
    })
}
```

#### Example: Render a chart

```go
chart := mrn.Chart(mrn.ChartProps{
    Type:        mrn.ChartTypeLine,
    Title:       "Weekly signups",
    Description: "New accounts by weekday.",
    Labels:      []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
    Datasets: []mrn.ChartDataset{
        {
            Label:           "Signups",
            Data:            []float64{12, 19, 14, 22, 18},
            BorderColor:     "#2563eb",
            BackgroundColor: "rgba(37, 99, 235, 0.16)",
            Fill:            true,
            Tension:         0.3,
        },
    },
    Options: mrn.ChartOptions{
        BeginAtZero: true,
        YAxisLabel:  "Users",
    },
})
```

---
