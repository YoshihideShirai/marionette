package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type componentIndex struct {
	Components []componentEntry `json:"components"`
}

type componentEntry struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Group       string `json:"group"`
	Description string `json:"description"`
	Golden      string `json:"golden"`
	Example     string `json:"example"`
	Template    string `json:"template"`
}

type componentGroup struct {
	Name       string
	Components []componentEntry
}

var pageTemplate = template.Must(template.New("docs").Parse(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>{{.Title}}</title>
    <style>
      :root { color-scheme: light dark; }
      body {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
        margin: 2rem auto;
        max-width: 880px;
        line-height: 1.65;
        padding: 0 1rem;
      }
      a { color: #2563eb; }
      pre {
        background: color-mix(in oklab, canvas 92%, black 8%);
        border-radius: 10px;
        overflow-x: auto;
        padding: 0.75rem;
      }
      code { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
      iframe { background: white; }
    </style>
  </head>
  <body>
    {{.Body}}
  </body>
</html>
`))

var componentTemplate = template.Must(template.New("component-docs").Funcs(template.FuncMap{
	"componentHref": componentHref,
}).Parse(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>{{.Title}}</title>
    <style>
      :root {
        color-scheme: light dark;
        --border: color-mix(in oklab, canvasText 16%, transparent);
        --muted: color-mix(in oklab, canvasText 68%, transparent);
        --surface: color-mix(in oklab, canvas 94%, canvasText 6%);
        --active: color-mix(in oklab, #2563eb 14%, canvas);
      }
      * { box-sizing: border-box; }
      body {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
        margin: 0;
        line-height: 1.65;
      }
      a { color: #2563eb; }
      .docs-shell {
        display: grid;
        grid-template-columns: 18rem minmax(0, 1fr);
        min-height: 100vh;
      }
      .docs-sidebar {
        border-right: 1px solid var(--border);
        padding: 1.25rem;
        position: sticky;
        top: 0;
        height: 100vh;
        overflow: auto;
      }
      .docs-brand {
        display: inline-flex;
        font-size: 0.95rem;
        font-weight: 700;
        margin-bottom: 1rem;
        text-decoration: none;
      }
      .docs-group { margin-top: 1.25rem; }
      .docs-group-title {
        color: var(--muted);
        font-size: 0.75rem;
        font-weight: 700;
        letter-spacing: 0;
        margin: 0 0 0.35rem;
        text-transform: uppercase;
      }
      .docs-nav {
        display: grid;
        gap: 0.1rem;
      }
      .docs-nav a {
        border-radius: 8px;
        color: inherit;
        display: block;
        font-size: 0.92rem;
        line-height: 1.3;
        padding: 0.45rem 0.55rem;
        text-decoration: none;
      }
      .docs-nav a:hover,
      .docs-nav a[aria-current="page"] {
        background: var(--active);
        color: #1d4ed8;
      }
      .docs-main {
        min-width: 0;
        padding: 2rem;
      }
      .docs-content {
        max-width: 960px;
      }
      .docs-index-grid {
        display: grid;
        gap: 1rem;
        grid-template-columns: repeat(auto-fit, minmax(16rem, 1fr));
        margin-top: 1.25rem;
      }
      .docs-link-card {
        border: 1px solid var(--border);
        border-radius: 8px;
        display: block;
        padding: 1rem;
        text-decoration: none;
      }
      .docs-link-card strong { color: canvasText; display: block; }
      .docs-link-card span { color: var(--muted); display: block; font-size: 0.9rem; margin-top: 0.25rem; }
      pre {
        background: var(--surface);
        border-radius: 8px;
        overflow-x: auto;
        padding: 0.75rem;
      }
      code { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
      iframe {
        background: white;
        max-width: 100%;
      }
      @media (max-width: 760px) {
        .docs-shell { display: block; }
        .docs-sidebar {
          border-bottom: 1px solid var(--border);
          border-right: 0;
          height: auto;
          max-height: 55vh;
          position: static;
        }
        .docs-main { padding: 1.25rem; }
      }
    </style>
  </head>
  <body>
    <div class="docs-shell">
      <aside class="docs-sidebar">
        <a class="docs-brand" href="./index.html">Components Gallery</a>
        {{range .Groups}}
          <section class="docs-group">
            <p class="docs-group-title">{{.Name}}</p>
            <nav class="docs-nav" aria-label="{{.Name}}">
              {{range .Components}}
                <a href="{{componentHref .}}"{{if eq $.ActiveID .ID}} aria-current="page"{{end}}>{{.Name}}</a>
              {{end}}
            </nav>
          </section>
        {{end}}
      </aside>
      <main class="docs-main">
        <div class="docs-content">
          {{.Body}}
        </div>
      </main>
    </div>
  </body>
</html>
`))

var mdHrefRe = regexp.MustCompile(`href="([^"]+?)\.md"`)
var iframeRe = regexp.MustCompile(`<iframe[^>]*></iframe>`)

func main() {
	if err := generate("docs/site/index.md"); err != nil {
		fmt.Fprintf(os.Stderr, "generate docs/site/index.md: %v\n", err)
		os.Exit(1)
	}
	if err := generateComponents("docs/site/components/index.md", "docs/site/components/_index.json"); err != nil {
		fmt.Fprintf(os.Stderr, "generate components: %v\n", err)
		os.Exit(1)
	}
}

func generate(path string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	renderedBody, err := renderMarkdown(src)
	if err != nil {
		return err
	}
	title := docTitle(src, filepath.Base(path))

	var out bytes.Buffer
	if err := pageTemplate.Execute(&out, map[string]any{
		"Title": title,
		"Body":  renderedBody,
	}); err != nil {
		return err
	}

	htmlPath := path[:len(path)-len(filepath.Ext(path))] + ".html"
	return os.WriteFile(htmlPath, out.Bytes(), 0o644)
}

func generateComponents(indexPath, catalogPath string) error {
	src, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	catalog, err := readComponentIndex(catalogPath)
	if err != nil {
		return err
	}
	groups := groupedComponents(catalog.Components)
	sections := componentSections(src)

	body, err := componentIndexBody(groups)
	if err != nil {
		return err
	}
	if err := writeComponentPage(indexPath[:len(indexPath)-len(filepath.Ext(indexPath))]+".html", "Components Gallery", "", groups, body); err != nil {
		return err
	}

	for _, entry := range catalog.Components {
		section, ok := sections[entry.Name]
		if !ok {
			section = fallbackComponentSection(entry)
		}
		section = ensureExamplePreview(section, entry)
		section = ensureUsageExample(section, entry)

		rendered, err := renderMarkdown([]byte(section))
		if err != nil {
			return fmt.Errorf("render %s: %w", entry.Name, err)
		}

		path := filepath.Join(filepath.Dir(indexPath), entry.ID+".html")
		if err := writeComponentPage(path, entry.Name, entry.ID, groups, rendered); err != nil {
			return err
		}
	}

	return nil
}

func readComponentIndex(path string) (componentIndex, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return componentIndex{}, err
	}

	var catalog componentIndex
	if err := json.Unmarshal(src, &catalog); err != nil {
		return componentIndex{}, err
	}
	for i := range catalog.Components {
		if catalog.Components[i].Group == "" {
			catalog.Components[i].Group = "Components"
		}
	}
	return catalog, nil
}

func groupedComponents(entries []componentEntry) []componentGroup {
	order := []string{"Components", "Layout"}
	seen := map[string]bool{}
	for _, entry := range entries {
		if !seen[entry.Group] && !slices.Contains(order, entry.Group) {
			order = append(order, entry.Group)
		}
		seen[entry.Group] = true
	}

	var groups []componentGroup
	for _, name := range order {
		var group componentGroup
		group.Name = name
		for _, entry := range entries {
			if entry.Group == name {
				group.Components = append(group.Components, entry)
			}
		}
		if len(group.Components) > 0 {
			groups = append(groups, group)
		}
	}
	return groups
}

func componentSections(src []byte) map[string]string {
	lines := strings.Split(string(src), "\n")
	sections := map[string]string{}

	for i := 0; i < len(lines); i++ {
		if !strings.HasPrefix(lines[i], "## ") {
			continue
		}

		title := strings.TrimSpace(strings.TrimPrefix(lines[i], "## "))
		if title == "Contents" {
			continue
		}

		start := i
		end := len(lines)
		for j := i + 1; j < len(lines); j++ {
			if strings.HasPrefix(lines[j], "## ") {
				end = j
				break
			}
		}

		sectionLines := append([]string{}, lines[start:end]...)
		sectionLines[0] = "# " + title
		sections[title] = strings.TrimSpace(strings.Join(sectionLines, "\n")) + "\n"
		i = end - 1
	}

	return sections
}

func componentIndexBody(groups []componentGroup) (template.HTML, error) {
	var body bytes.Buffer
	body.WriteString(`<h1>Components Gallery</h1>
<p>Marionette UI docs are split into one page per component. Use the sidebar to jump between pages, grouped by general components and layout primitives.</p>
`)

	for _, group := range groups {
		fmt.Fprintf(&body, "<h2>%s</h2>\n<div class=\"docs-index-grid\">\n", template.HTMLEscapeString(group.Name))
		for _, entry := range group.Components {
			detail := entry.Description
			if detail == "" {
				detail = entry.Template
			}
			fmt.Fprintf(
				&body,
				"<a class=\"docs-link-card\" href=\"%s\"><strong>%s</strong><span>%s</span></a>\n",
				template.HTMLEscapeString(componentHref(entry)),
				template.HTMLEscapeString(entry.Name),
				template.HTMLEscapeString(detail),
			)
		}
		body.WriteString("</div>\n")
	}

	return template.HTML(body.String()), nil
}

func fallbackComponentSection(entry componentEntry) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", entry.Name)
	if entry.Example != "" && strings.HasPrefix(entry.Example, "docs/site/components/") {
		fmt.Fprintf(&b, "### Visual\n\n%s\n\n", iframeFor(entry))
	}
	if entry.Golden != "" {
		fmt.Fprintf(&b, "- Golden sample: [`%s`](https://github.com/YoshihideShirai/marionette/blob/main/%s)\n", filepath.Base(entry.Golden), entry.Golden)
	}
	if entry.Template != "" {
		fmt.Fprintf(&b, "- Template: [`%s`](https://github.com/YoshihideShirai/marionette/blob/main/%s)\n", entry.Template, entry.Template)
	}
	return b.String()
}

func ensureExamplePreview(section string, entry componentEntry) string {
	if entry.Example == "" || !strings.HasPrefix(entry.Example, "docs/site/components/") {
		return section
	}

	iframe := iframeFor(entry)
	src := strings.TrimPrefix(entry.Example, "docs/site/components/")
	if strings.Contains(section, "<iframe") {
		if strings.Contains(section, `src="./`+src+`"`) {
			return section
		}
		return iframeRe.ReplaceAllString(section, iframe)
	}

	insert := "\n### Visual\n\n" + iframe + "\n"
	marker := "- Golden sample:"
	if idx := strings.Index(section, marker); idx >= 0 {
		return strings.TrimSpace(section[:idx]) + insert + "\n" + strings.TrimSpace(section[idx:]) + "\n"
	}
	return strings.TrimSpace(section) + insert
}

func ensureUsageExample(section string, entry componentEntry) string {
	snippet := usageSnippet(entry.ID)
	if snippet == "" || strings.Contains(section, "### Go usage") {
		return section
	}

	block := "### Go usage\n\n```go\n" + strings.TrimSpace(snippet) + "\n```\n"
	iframe := iframeRe.FindString(section)
	if iframe == "" {
		marker := "- Golden sample:"
		if idx := strings.Index(section, marker); idx >= 0 {
			return strings.TrimSpace(section[:idx]) + "\n\n" + block + "\n" + strings.TrimSpace(section[idx:]) + "\n"
		}
		return strings.TrimSpace(section) + "\n\n" + block
	}

	return strings.Replace(section, iframe, iframe+"\n\n"+block, 1)
}

func iframeFor(entry componentEntry) string {
	height := "320px"
	switch entry.ID {
	case "button":
		height = "160px"
	case "input", "pagination":
		height = "180px"
	case "select":
		height = "200px"
	case "empty-state", "form-field":
		height = "260px"
	case "checkbox", "radio-group", "switch", "tabs", "breadcrumb":
		height = "220px"
	case "table", "modal":
		height = "320px"
	case "chart":
		height = "420px"
	case "textarea":
		height = "240px"
	case "feedback":
		height = "420px"
	case "stack", "page-header", "container", "card", "section":
		height = "320px"
	case "grid", "split":
		height = "380px"
	}

	src := strings.TrimPrefix(entry.Example, "docs/site/components/")
	return fmt.Sprintf(`<iframe src="./%s" title="%s example" style="width:100%%;min-height:%s;border:1px solid #e5e7eb;border-radius:8px;"></iframe>`, src, entry.Name, height)
}

func usageSnippet(id string) string {
	snippets := map[string]string{
		"button": `saveButton := mf.ComponentButton("Save", mf.ComponentProps{
    Variant: "secondary",
    Size:    "sm",
})`,
		"theme-toggle-button": `themeToggle := mf.ComponentThemeToggleButton("/toggle-theme", "Toggle theme", mf.ComponentProps{
    Variant: "ghost",
    Size:    "sm",
})`,
		"input": `startDate := mf.ComponentInputWithOptions("start_date", "2030-01-01", mf.InputOptions{
    Type:        "date",
    Placeholder: "Start date",
    Min:         "2024-01-01",
    Max:         "2035-12-31",
    Required:    true,
    Props: mf.ComponentProps{
        Variant:  "ghost",
        Size:     "sm",
        Disabled: true,
    },
})`,
		"select": `roleSelect := mf.ComponentSelect("role", []mf.SelectOption{
    {Label: "Admin", Value: "admin", Selected: true},
    {Label: "Viewer", Value: "viewer"},
}, mf.ComponentProps{
    Variant: "ghost",
    Size:    "sm",
})`,
		"modal": `confirmModal := mf.ComponentModal(mf.ModalProps{
    Title: "Delete user",
    Body:  mf.Text("Confirm deletion"),
    Actions: mf.ComponentButton("Delete", mf.ComponentProps{
        Variant: "error",
        Size:    "sm",
    }),
    Open: true,
})`,
		"empty-state": `emptyUsers := mf.ComponentEmptyState(mf.EmptyStateProps{
    Title:       "No users",
    Description: "Create one first.",
})`,
		"table": `usersTable := mf.ComponentTable(mf.TableProps{
    Columns: []mf.TableColumn{
        {Label: "Name", SortKey: "name", SortHref: "/?sort=name", SortActive: true},
        {Label: "Role"},
    },
    Rows: []mf.TableComponentRow{
        {
            Cells: []mf.Node{
                mf.Text("Aiko"),
                mf.DivClass("badge", mf.Text("Admin")),
            },
        },
    },
})`,
		"chart": `signupsChart := mf.ComponentChart(mf.ChartProps{
    Type:        mf.ChartTypeLine,
    Title:       "Weekly signups",
    Description: "New accounts by weekday.",
    Labels:      []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
    Datasets: []mf.ChartDataset{
        {
            Label:           "Signups",
            Data:            []float64{12, 19, 14, 22, 18},
            BorderColor:     "#2563eb",
            BackgroundColor: "rgba(37, 99, 235, 0.16)",
            Fill:            true,
            Tension:         0.3,
        },
    },
    Options: mf.ChartOptions{
        BeginAtZero: true,
        YAxisLabel:  "Users",
    },
})`,
		"pagination": `pager := mf.ComponentPagination(mf.PaginationProps{
    Page:       2,
    TotalPages: 4,
    PrevHref:   "/?page=1&per_page=10",
    NextHref:   "/?page=3&per_page=10",
})`,
		"form": `profileForm := mf.ComponentForm("/users", "post", mf.ComponentButton("Save", mf.ComponentProps{
    Variant: "primary",
    Size:    "sm",
}))`,
		"form-field": `nameField := mf.ComponentFormField(
    mf.ComponentInput("name", "", mf.ComponentProps{Size: "sm"}),
    mf.FormFieldProps{
        Label:    "Name",
        Required: true,
        Hint:     "Enter a display name.",
        Error:    "Name is required.",
    },
)`,
		"tabs": `tabs := mf.ComponentTabs(mf.TabsProps{
    AriaLabel: "user sections",
    Items: []mf.TabsItem{
        {Label: "Profile", Href: "/users/1/profile", Active: true},
        {Label: "Permissions", Href: "/users/1/permissions"},
        {Label: "Audit", Disabled: true},
    },
})`,
		"breadcrumb": `breadcrumb := mf.ComponentBreadcrumb(mf.BreadcrumbProps{
    Items: []mf.BreadcrumbItem{
        {Label: "Home", Href: "/"},
        {Label: "Users", Href: "/users"},
        {Label: "Aiko", Active: true},
    },
})`,
		"textarea": `notes := mf.ComponentTextarea("notes", "hello", mf.TextareaOptions{
    Placeholder: "Memo",
    Rows:        4,
    Required:    true,
    Props: mf.ComponentProps{
        Variant: "ghost",
        Size:    "sm",
    },
})`,
		"checkbox": `activeUser := mf.ComponentCheckbox(mf.CheckboxComponentProps{
    Name:    "active",
    Value:   "1",
    Label:   "Active user",
    Checked: true,
    Props:   mf.ComponentProps{Size: "sm"},
})`,
		"radio-group": `roleGroup := mf.ComponentRadioGroup(mf.RadioGroupComponentProps{
    Name:      "role",
    AriaLabel: "role",
    Items: []mf.RadioItem{
        {Label: "Admin", Value: "admin", Checked: true},
        {Label: "Editor", Value: "editor"},
        {Label: "Viewer", Value: "viewer", Disabled: true},
    },
    Props: mf.ComponentProps{Size: "sm"},
})`,
		"switch": `notifications := mf.ComponentSwitch(mf.SwitchComponentProps{
    Name:    "notify",
    Value:   "1",
    Label:   "Enable notifications",
    Checked: true,
    Props:   mf.ComponentProps{Size: "sm"},
})`,
		"toast": `toast := mf.ComponentToast(mf.ToastProps{
    Title:       "Saved",
    Description: "User settings were updated.",
    Icon:        "✓",
    Props:       mf.ComponentProps{Variant: "success", Size: "md"},
})`,
		"alert": `alert := mf.ComponentAlert(mf.AlertProps{
    Title:       "Permission denied",
    Description: "You do not have access to this workspace.",
    Icon:        "!",
    Props:       mf.ComponentProps{Variant: "error", Size: "md"},
})`,
		"skeleton": `loadingCard := mf.ComponentCard(mf.CardProps{Title: "Loading"}, mf.ComponentSkeleton(mf.SkeletonProps{
    Lines: 3,
    Props: mf.ComponentProps{Size: "md"},
}))`,
		"stack": `stack := mf.ComponentStack(
    mf.StackProps{
        Direction: "horizontal",
        Gap:       "sm",
        Align:     "center",
        Justify:   "between",
        Wrap:      true,
    },
    mf.Text("Aiko Tanaka"),
    mf.DivClass("badge badge-primary", mf.Text("Admin")),
    mf.ComponentButton("Open", mf.ComponentProps{Variant: "secondary", Size: "sm"}),
)`,
		"grid": `summaryGrid := mf.ComponentGrid(
    mf.GridProps{Columns: "3", Gap: "lg"},
    mf.ComponentCard(mf.CardProps{}, mf.Text("Users: 24")),
    mf.ComponentCard(mf.CardProps{}, mf.Text("Admins: 4")),
    mf.ComponentCard(mf.CardProps{}, mf.Text("Pending: 7")),
)`,
		"split": `workspace := mf.ComponentSplit(mf.SplitProps{
    Main: mf.ComponentCard(
        mf.CardProps{Title: "Main workspace"},
        mf.Text("Aiko / Admin / Active"),
    ),
    Aside: mf.ComponentSection(
        mf.SectionProps{Title: "Aside panel"},
        mf.ComponentButton("Apply", mf.ComponentProps{Variant: "primary", Size: "sm"}),
    ),
    AsideWidth: "md",
    Gap:        "lg",
})`,
		"page-header": `header := mf.ComponentPageHeader(mf.PageHeaderProps{
    Title:       "Users",
    Description: "Manage account records, roles, and invitations.",
    Actions: mf.ComponentStack(
        mf.StackProps{Direction: "horizontal", Gap: "sm"},
        mf.ComponentButton("Export", mf.ComponentProps{Variant: "ghost", Size: "sm"}),
        mf.ComponentButton("Create", mf.ComponentProps{Variant: "primary", Size: "sm"}),
    ),
})`,
		"container": `page := mf.ComponentContainer(
    mf.ContainerProps{
        MaxWidth: "md",
        Padding:  "md",
        Centered: true,
    },
    mf.Text("Centered page container"),
)`,
		"card": `card := mf.ComponentCard(
    mf.CardProps{
        Title:       "Workspace summary",
        Description: "Header, description, actions, then body content.",
        Actions:     mf.ComponentButton("Edit", mf.ComponentProps{Variant: "ghost", Size: "sm"}),
    },
    mf.Text("Active: 24"),
)`,
		"section": `section := mf.ComponentSection(
    mf.SectionProps{
        Title:       "Recent activity",
        Description: "An unframed content section with consistent header spacing.",
        Actions:     mf.ComponentButton("View all", mf.ComponentProps{Variant: "secondary", Size: "sm"}),
    },
    mf.Text("Aiko updated a role"),
)`,
		"feedback": `toast := mf.ComponentToast(mf.ToastProps{
    Title:       "Toast / Default",
    Description: "Saved successfully.",
    Icon:        "✓",
    Props:       mf.ComponentProps{Variant: "success", Size: "md"},
})

alert := mf.ComponentAlert(mf.AlertProps{
    Title:       "Alert / Long content",
    Description: "An error occurred. Please wait a moment and try again.",
    Icon:        "!",
    Props:       mf.ComponentProps{Variant: "error", Size: "lg"},
})`,
		"chart-line": `trend := mf.ComponentChart(mf.ChartProps{
    Type:   mf.ChartTypeLine,
    Labels: []string{"Mon", "Tue", "Wed"},
    Datasets: []mf.ChartDataset{{
        Label: "Signups", Data: []float64{12, 19, 14},
    }},
})`,
		"chart-bar": `breakdown := mf.ComponentChart(mf.ChartProps{
    Type:   mf.ChartTypeBar,
    Labels: []string{"Admin", "Editor", "Viewer"},
    Datasets: []mf.ChartDataset{{
        Label: "Users", Data: []float64{4, 7, 13},
    }},
})`,
		"chart-pie": `share := mf.ComponentChart(mf.ChartProps{
    Type:   mf.ChartTypePie,
    Labels: []string{"Desktop", "Mobile", "Tablet"},
    Datasets: []mf.ChartDataset{{
        Data: []float64{62, 29, 9},
    }},
})`,
		"chart-doughnut": `conversion := mf.ComponentChart(mf.ChartProps{
    Type:   mf.ChartTypeDoughnut,
    Labels: []string{"Converted", "Dropped"},
    Datasets: []mf.ChartDataset{{
        Data: []float64{72, 28},
    }},
})`,
		"chart-scatter": `correlation := mf.ComponentChart(mf.ChartProps{
    Type: mf.ChartTypeScatter,
    Datasets: []mf.ChartDataset{{
        Label: "Sessions",
        Points: []mf.ChartPoint{
            {X: 4, Y: 12},
            {X: 8, Y: 26},
            {X: 12, Y: 33},
        },
    }},
})`,
		"overlay-system": `modal := mf.ComponentModal(mf.ModalProps{
    Title: "Overlay demo",
    Body:  mf.Text("Use modal, drawer, and popover together for rich flows."),
    Open:  true,
})`,
		"dataframe": `df := dataframe.LoadRecords(
    []struct{ Name string; Role string }{
        {Name: "Aiko", Role: "Admin"},
        {Name: "Ken", Role: "Viewer"},
    },
)

table := mf.ComponentDataFrame(df, mf.DataFrameProps{})`,
		"dataframe-chart": `chart := mf.ComponentDataFrameChart(df, mf.DataFrameChartProps{
    Type:      mf.ChartTypeBar,
    LabelCol:  "month",
    ValueCols: []string{"signups"},
})`,
	}
	return snippets[id]
}

func writeComponentPage(path, title, activeID string, groups []componentGroup, body template.HTML) error {
	var out bytes.Buffer
	if err := componentTemplate.Execute(&out, map[string]any{
		"Title":    title,
		"ActiveID": activeID,
		"Groups":   groups,
		"Body":     body,
	}); err != nil {
		return err
	}
	return os.WriteFile(path, out.Bytes(), 0o644)
}

func componentHref(entry componentEntry) string {
	return "./" + entry.ID + ".html"
}

func renderMarkdown(src []byte) (template.HTML, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	var body bytes.Buffer
	if err := md.Convert(src, &body); err != nil {
		return "", err
	}
	return template.HTML(rewriteMarkdownLinks(body.String())), nil
}

func rewriteMarkdownLinks(s string) string {
	return mdHrefRe.ReplaceAllStringFunc(s, func(m string) string {
		parts := mdHrefRe.FindStringSubmatch(m)
		if len(parts) < 2 {
			return m
		}

		href := parts[1]
		if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
			return m
		}

		return `href="` + href + `.html"`
	})
}

func docTitle(src []byte, fallback string) string {
	for _, line := range bytes.Split(src, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("# ")) {
			return string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("# "))))
		}
	}
	return fallback
}
