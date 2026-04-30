package marionette

import (
	"strings"
	"sync"
	"testing"

	mf "github.com/YoshihideShirai/marionette/frontend"
	rdf "github.com/rocketlaunchr/dataframe-go"
)

func TestButtonRenderUsesHTMXMarkup(t *testing.T) {
	html, err := Button("Increment").OnClick("counter/increment").Target("#app").Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	if got == "" {
		t.Fatal("expected non-empty button html")
	}
	if want := `hx-post="/counter/increment"`; !strings.Contains(got, want) {
		t.Fatalf("expected %q in %q", want, got)
	}
}

func TestButtonPostAcceptsLeadingSlash(t *testing.T) {
	html, err := Button("Save").Post("/users/create").Target("#users").Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	if strings.Contains(got, `hx-post="//users/create"`) {
		t.Fatalf("expected normalized post path, got %q", got)
	}
	if want := `hx-target="#users"`; !strings.Contains(got, want) {
		t.Fatalf("expected %q in %q", want, got)
	}
}

func TestFormInputAndSubmitRenderHTMXMarkup(t *testing.T) {
	html, err := Form("users/create", Input("name", `<Aiko>`), Submit("Create")).Target("#users").Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`hx-post="/users/create"`,
		`hx-target="#users"`,
		`name="name"`,
		`value="&lt;Aiko&gt;"`,
		`type="submit"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestSidebarRendersNavigationAndEscapesText(t *testing.T) {
	html, err := Sidebar("Marionette", "Admin <Console>",
		SidebarLink("Users", "/").Active(),
		SidebarLink("Settings", "/settings"),
	).Note("Demo", `<unsafe>`).Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`<aside`,
		`href="/"`,
		`href="/settings"`,
		`btn btn-primary justify-start`,
		`Admin &lt;Console&gt;`,
		`&lt;unsafe&gt;`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
	if strings.Contains(got, "Admin <Console>") || strings.Contains(got, `<unsafe>`) {
		t.Fatalf("expected sidebar text to be escaped, got %q", got)
	}
}

func TestTableRendersHeadersRowsAndEscapesCells(t *testing.T) {
	html, err := Table([]string{"Name", "Role"},
		TableRow(Text(`<Aiko>`), DivClass("badge", Text("Admin"))),
	).Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`<table class="table">`,
		`<th>Name</th>`,
		`<td><span>&lt;Aiko&gt;</span></td>`,
		`class="badge"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
	if strings.Contains(got, `<Aiko>`) {
		t.Fatalf("expected table cell text to be escaped, got %q", got)
	}
}

func TestHiddenInputRenderEscapesValue(t *testing.T) {
	html, err := HiddenInput("id", `"42"`).Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`name="id"`,
		`type="hidden"`,
		`value="&#34;42&#34;"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestElementRenderEscapesText(t *testing.T) {
	html, err := Text(`<script>alert(1)</script>`).Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	if strings.Contains(got, `<script>alert(1)</script>`) {
		t.Fatalf("expected escaped content, got %q", got)
	}
}

func TestDivConstructorsSupportPlainAndAttributedMarkup(t *testing.T) {
	plainHTML, err := Div(Text("Plain")).Render()
	if err != nil {
		t.Fatalf("plain div render failed: %v", err)
	}
	if got, want := string(plainHTML), `<div><span>Plain</span></div>`; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}

	attrsHTML, err := DivAttrs(Attrs{"data-ref": "plain"}, Text("Attrs")).Render()
	if err != nil {
		t.Fatalf("attrs div render failed: %v", err)
	}
	if got := string(attrsHTML); !strings.Contains(got, `data-ref="plain"`) {
		t.Fatalf("expected data-ref attribute in %q", got)
	}

	attrHTML, err := DivProps(ElementProps{
		ID:    "panel",
		Class: "p-4",
		Attrs: Attrs{
			"class":    "rounded-box",
			"data-ref": `<unsafe>`,
		},
	}, Text("Panel")).Render()
	if err != nil {
		t.Fatalf("attributed div render failed: %v", err)
	}
	got := string(attrHTML)
	for _, want := range []string{
		`class="rounded-box p-4"`,
		`data-ref="&lt;unsafe&gt;"`,
		`id="panel"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestTemplatePartialsRenderSharedProps(t *testing.T) {
	buttonHTML, err := ButtonComponent("Send", ComponentProps{Class: "tracking-wide", Variant: "secondary", Size: "sm", Disabled: true}).Render()
	if err != nil {
		t.Fatalf("button render failed: %v", err)
	}
	inputHTML, err := InputComponent("email", "demo@example.com", ComponentProps{Variant: "ghost", Size: "sm", Disabled: true}).Render()
	if err != nil {
		t.Fatalf("input render failed: %v", err)
	}
	selectHTML, err := SelectComponent("role", []SelectOption{{Label: "Viewer", Value: "viewer", Selected: true}}, ComponentProps{Variant: "ghost", Size: "sm", Disabled: true}).Render()
	if err != nil {
		t.Fatalf("select render failed: %v", err)
	}

	for _, tc := range []struct {
		name string
		html string
		want []string
	}{
		{name: "button", html: string(buttonHTML), want: []string{`btn-secondary`, `btn-sm`, `tracking-wide`, `disabled`}},
		{name: "input", html: string(inputHTML), want: []string{`input-ghost`, `input-sm`, `name="email"`, `disabled`}},
		{name: "select", html: string(selectHTML), want: []string{`select-ghost`, `select-sm`, `name="role"`, `selected`, `disabled`}},
	} {
		for _, want := range tc.want {
			if !strings.Contains(tc.html, want) {
				t.Fatalf("%s expected %q in %q", tc.name, want, tc.html)
			}
		}
	}
}

func TestLoadComponentTemplatesCachesParsedTemplates(t *testing.T) {
	cachedTemplates = nil
	cachedTemplatesErr = nil
	componentTemplatesOnce = sync.Once{}

	first, err := loadComponentTemplates()
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}
	second, err := loadComponentTemplates()
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}
	if first == nil {
		t.Fatal("expected first template set to be non-nil")
	}
	if second == nil {
		t.Fatal("expected second template set to be non-nil")
	}
	if first != second {
		t.Fatalf("expected cached template pointer reuse, got %p and %p", first, second)
	}
}

func TestMarkdownRendersHTMLAndFiltersRawHTML(t *testing.T) {
	html, err := Markdown(MarkdownProps{
		Content: "# Title\n\n<script>alert(1)</script>\n\n**bold**",
		Props:   ComponentProps{Class: "max-w-none"},
	}).Render()
	if err != nil {
		t.Fatalf("markdown render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`<article class="prose max-w-none">`,
		`<h1>Title</h1>`,
		`<strong>bold</strong>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
	if strings.Contains(got, "<script>") {
		t.Fatalf("expected raw html to be filtered, got %q", got)
	}
}

func TestComponentInputWithOptionsRendersDateConstraints(t *testing.T) {
	html, err := InputWithOptions("start_date", "2030-01-01", InputOptions{
		Type:     "date",
		Min:      "2024-01-01",
		Max:      "2026-12-31",
		Required: true,
		Props:    ComponentProps{Variant: "default", Size: "sm"},
	}).Render()
	if err != nil {
		t.Fatalf("input render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{`type="date"`, `min="2024-01-01"`, `max="2026-12-31"`, `required`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestComponentFormFieldRendersLabelHintAndError(t *testing.T) {
	html, err := FormFieldComponent(
		InputWithOptions("name", "", InputOptions{Required: true, Props: ComponentProps{Variant: "default", Size: "sm"}}),
		FormFieldProps{
			Label:    "Name",
			Required: true,
			Hint:     "Enter a display name.",
			Error:    "Name is required.",
		},
	).Render()
	if err != nil {
		t.Fatalf("form field render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{`label-text`, `Name`, `*`, `Enter a display name.`, `Name is required.`, `name="name"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestActionFormRendersHTMXActionAttributes(t *testing.T) {
	html, err := ActionForm(ActionFormProps{
		Action: "tasks/create",
		Target: "#task-list",
		Swap:   "innerHTML",
		Props:  ComponentProps{Class: "space-y-3"},
	}, TextField(TextFieldProps{ID: "task-name", Name: "name"})).Render()
	if err != nil {
		t.Fatalf("action form render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`action="/tasks/create"`,
		`method="post"`,
		`hx-post="/tasks/create"`,
		`hx-target="#task-list"`,
		`hx-swap="innerHTML"`,
		`class="space-y-3"`,
		`name="name"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestFormRowAndTextFieldWireA11yAttributes(t *testing.T) {
	html, err := FormRow(FormRowProps{
		ID:          "email",
		Label:       "Email",
		Description: "Used for notifications.",
		Error:       "Email is required.",
		Required:    true,
		Control: TextField(TextFieldProps{
			ID:          "email",
			Name:        "email",
			Value:       "",
			Description: "Used for notifications.",
			Error:       "Email is required.",
			Required:    true,
			Ref:         "register-email",
		}),
	}).Render()
	if err != nil {
		t.Fatalf("form row render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`for="email"`,
		`id="email-description"`,
		`id="email-error"`,
		`aria-describedby="email-description email-error"`,
		`aria-invalid="true"`,
		`name="email"`,
		`data-ref="register-email"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestSelectCheckboxRadioAndSwitchExposeState(t *testing.T) {
	selectHTML, err := Select(SelectFieldProps{
		ID:          "role",
		Name:        "role",
		Options:     []SelectOption{{Label: "Admin", Value: "admin", Selected: true}},
		Description: "Role selection",
		Disabled:    true,
	}).Render()
	if err != nil {
		t.Fatalf("select render failed: %v", err)
	}
	for _, want := range []string{`id="role"`, `name="role"`, `disabled="disabled"`, `selected="selected"`} {
		if !strings.Contains(string(selectHTML), want) {
			t.Fatalf("expected %q in %q", want, selectHTML)
		}
	}

	radioHTML, err := RadioGroup(RadioGroupProps{
		ID:      "access",
		Name:    "access",
		Value:   "write",
		Options: []RadioOption{{Label: "Read", Value: "read"}, {Label: "Write", Value: "write"}},
	}).Render()
	if err != nil {
		t.Fatalf("radio render failed: %v", err)
	}
	if !strings.Contains(string(radioHTML), `checked="checked"`) {
		t.Fatalf("expected checked state in %q", radioHTML)
	}

	switchHTML, err := Switch(SwitchProps{ID: "enabled", Name: "enabled", Label: "Enabled", Checked: true, ReadOnly: true}).Render()
	if err != nil {
		t.Fatalf("switch render failed: %v", err)
	}
	for _, want := range []string{`class="toggle`, `readonly="readonly"`, `checked="checked"`} {
		if !strings.Contains(string(switchHTML), want) {
			t.Fatalf("expected %q in %q", want, switchHTML)
		}
	}

	checkboxHTML, err := Checkbox(CheckboxProps{ID: "tos", Name: "tos", Label: "Accept", Error: "Required"}).Render()
	if err != nil {
		t.Fatalf("checkbox render failed: %v", err)
	}
	if !strings.Contains(string(checkboxHTML), `aria-invalid="true"`) {
		t.Fatalf("expected invalid checkbox in %q", checkboxHTML)
	}
}

func TestComponentModalRendersSSRState(t *testing.T) {
	closedHTML, err := Modal(ModalProps{
		Title:   "Delete user",
		Body:    Text("Confirm deletion"),
		Actions: ButtonComponent("Cancel", ComponentProps{Variant: "ghost", Size: "sm"}),
		Open:    false,
	}).Render()
	if err != nil {
		t.Fatalf("closed modal render failed: %v", err)
	}
	openHTML, err := Modal(ModalProps{
		Title:   "Delete user",
		Body:    Text("Confirm deletion"),
		Actions: ButtonComponent("Delete", ComponentProps{Variant: "danger", Size: "sm"}),
		Open:    true,
	}).Render()
	if err != nil {
		t.Fatalf("open modal render failed: %v", err)
	}

	if strings.Contains(string(closedHTML), "modal-open") {
		t.Fatalf("expected closed modal without modal-open class, got %q", closedHTML)
	}
	for _, want := range []string{`modal-open`, `Delete user`, `Confirm deletion`, `btn-error`} {
		if !strings.Contains(string(openHTML), want) {
			t.Fatalf("expected %q in %q", want, openHTML)
		}
	}
}

func TestComponentEmptyStateRendersSkeletonAndCopy(t *testing.T) {
	skeletonHTML, err := EmptyState(EmptyStateProps{Skeleton: true, Rows: 2}).Render()
	if err != nil {
		t.Fatalf("skeleton render failed: %v", err)
	}
	if !strings.Contains(string(skeletonHTML), `aria-busy="true"`) {
		t.Fatalf("expected skeleton aria-busy state, got %q", skeletonHTML)
	}

	emptyHTML, err := EmptyState(EmptyStateProps{Title: "No users", Description: "Create one first."}).Render()
	if err != nil {
		t.Fatalf("empty render failed: %v", err)
	}
	for _, want := range []string{"No users", "Create one first."} {
		if !strings.Contains(string(emptyHTML), want) {
			t.Fatalf("expected %q in %q", want, emptyHTML)
		}
	}
}

func TestComponentTableRendersSortHeadersAndEmptyState(t *testing.T) {
	emptyHTML, err := TableComponent(TableProps{
		Columns: []TableColumn{
			{Label: "Name", SortKey: "name", SortHref: "/?sort=name", SortActive: true},
			{Label: "Role"},
		},
		EmptyTitle:       "No users",
		EmptyDescription: "Create a user to get started.",
	}).Render()
	if err != nil {
		t.Fatalf("empty table render failed: %v", err)
	}
	got := string(emptyHTML)
	for _, want := range []string{`href="/?sort=name"`, `No users`, `Create a user to get started.`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestComponentPaginationRendersState(t *testing.T) {
	html, err := Pagination(PaginationProps{
		Page:       2,
		TotalPages: 4,
		PrevHref:   "/?page=1&per_page=10",
		NextHref:   "/?page=3&per_page=10",
	}).Render()
	if err != nil {
		t.Fatalf("pagination render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{"Page 2 / 4", `href="/?page=1&amp;per_page=10"`, `href="/?page=3&amp;per_page=10"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestFeedbackComponentsShareVariantSizeAndA11y(t *testing.T) {
	toastHTML, err := Toast(ToastProps{
		Title:       "Saved",
		Description: "All changes were synced.",
		Icon:        "✓",
		Props:       ComponentProps{Variant: "success", Size: "sm"},
	}).Render()
	if err != nil {
		t.Fatalf("toast render failed: %v", err)
	}
	for _, want := range []string{`ui-feedback-success`, `ui-feedback-sm`, `role="status"`, `aria-live="polite"`} {
		if !strings.Contains(string(toastHTML), want) {
			t.Fatalf("expected %q in %q", want, toastHTML)
		}
	}

	alertHTML, err := Alert(AlertProps{Title: "Failed", Description: "Try again.", Icon: "!", Props: ComponentProps{Variant: "error", Size: "md"}}).Render()
	if err != nil {
		t.Fatalf("alert render failed: %v", err)
	}
	for _, want := range []string{`ui-feedback-error`, `ui-feedback-md`, `role="alert"`, `aria-live="assertive"`} {
		if !strings.Contains(string(alertHTML), want) {
			t.Fatalf("expected %q in %q", want, alertHTML)
		}
	}

	skeletonHTML, err := Skeleton(SkeletonProps{Rows: 2, Props: ComponentProps{Variant: "warning", Size: "lg"}}).Render()
	if err != nil {
		t.Fatalf("skeleton render failed: %v", err)
	}
	for _, want := range []string{`ui-feedback-warning`, `ui-feedback-lg`, `aria-busy="true"`} {
		if !strings.Contains(string(skeletonHTML), want) {
			t.Fatalf("expected %q in %q", want, skeletonHTML)
		}
	}
}

func TestLayoutComponentsRenderClassesAndContent(t *testing.T) {
	stackHTML, err := Stack(
		StackProps{Direction: "horizontal", Gap: "sm", Align: "center", Justify: "between", Wrap: true, Props: ComponentProps{Class: "w-full"}},
		Text("Left"),
		Text("Right"),
	).Render()
	if err != nil {
		t.Fatalf("stack render failed: %v", err)
	}
	for _, want := range []string{"flex-row", "gap-2", "items-center", "justify-between", "flex-wrap", "w-full", "Left", "Right"} {
		if !strings.Contains(string(stackHTML), want) {
			t.Fatalf("expected %q in %q", want, stackHTML)
		}
	}

	gridHTML, err := Grid(
		GridProps{Columns: "4", Gap: "lg"},
		Text("A"),
		Text("B"),
	).Render()
	if err != nil {
		t.Fatalf("grid render failed: %v", err)
	}
	for _, want := range []string{"grid", "gap-6", "grid-cols-1 sm:grid-cols-2 xl:grid-cols-4", "A", "B"} {
		if !strings.Contains(string(gridHTML), want) {
			t.Fatalf("expected %q in %q", want, gridHTML)
		}
	}

	splitHTML, err := Split(SplitProps{
		Main:            Text("Main"),
		Aside:           Text("Aside"),
		AsideWidth:      "lg",
		ReverseOnMobile: true,
	}).Render()
	if err != nil {
		t.Fatalf("split render failed: %v", err)
	}
	for _, want := range []string{"lg:grid-cols-[minmax(0,1fr)_28rem]", "order-2 lg:order-1", "order-1 lg:order-2", "Main", "Aside"} {
		if !strings.Contains(string(splitHTML), want) {
			t.Fatalf("expected %q in %q", want, splitHTML)
		}
	}
}

func TestSurfaceLayoutComponentsRenderHeadersActionsAndChildren(t *testing.T) {
	headerHTML, err := PageHeader(PageHeaderProps{
		Title:       "Users",
		Description: "Manage users",
		Actions:     ButtonComponent("Create", ComponentProps{Size: "sm"}),
	}).Render()
	if err != nil {
		t.Fatalf("page header render failed: %v", err)
	}
	for _, want := range []string{"<header", "Users", "Manage users", "btn-sm"} {
		if !strings.Contains(string(headerHTML), want) {
			t.Fatalf("expected %q in %q", want, headerHTML)
		}
	}

	containerHTML, err := ContainerComponent(ContainerProps{MaxWidth: "md", Padding: "sm", Centered: true}, Text("Contained")).Render()
	if err != nil {
		t.Fatalf("container render failed: %v", err)
	}
	for _, want := range []string{"max-w-5xl", "p-3", "mx-auto", "Contained"} {
		if !strings.Contains(string(containerHTML), want) {
			t.Fatalf("expected %q in %q", want, containerHTML)
		}
	}

	regionHTML, err := Region(RegionProps{ID: "task-list", Props: ComponentProps{Class: "space-y-3"}}, Text("Tasks")).Render()
	if err != nil {
		t.Fatalf("region render failed: %v", err)
	}
	for _, want := range []string{`id="task-list"`, `class="space-y-3"`, "Tasks"} {
		if !strings.Contains(string(regionHTML), want) {
			t.Fatalf("expected %q in %q", want, regionHTML)
		}
	}

	badgeHTML, err := Badge(BadgeProps{Label: "Admin", Props: ComponentProps{Variant: "primary", Size: "sm"}}).Render()
	if err != nil {
		t.Fatalf("badge render failed: %v", err)
	}
	for _, want := range []string{"badge", "badge-primary", "badge-sm", "Admin"} {
		if !strings.Contains(string(badgeHTML), want) {
			t.Fatalf("expected %q in %q", want, badgeHTML)
		}
	}

	actionsHTML, err := Actions(ActionsProps{Align: "end", Gap: "sm", Wrap: true, Props: ComponentProps{Class: "w-full"}}, ButtonComponent("Save", ComponentProps{})).Render()
	if err != nil {
		t.Fatalf("actions render failed: %v", err)
	}
	for _, want := range []string{"flex", "items-center", "gap-2", "justify-end", "flex-wrap", "w-full", "Save"} {
		if !strings.Contains(string(actionsHTML), want) {
			t.Fatalf("expected %q in %q", want, actionsHTML)
		}
	}

	dividerHTML, err := Divider(DividerProps{Spacing: "xs"}).Render()
	if err != nil {
		t.Fatalf("divider render failed: %v", err)
	}
	for _, want := range []string{"divider", "my-1"} {
		if !strings.Contains(string(dividerHTML), want) {
			t.Fatalf("expected %q in %q", want, dividerHTML)
		}
	}

	textHTML, err := TextComponent(TextProps{Text: "Muted", Size: "sm", Weight: "medium", Tone: "muted"}).Render()
	if err != nil {
		t.Fatalf("text render failed: %v", err)
	}
	for _, want := range []string{"text-sm", "font-medium", "text-base-content/60", "Muted"} {
		if !strings.Contains(string(textHTML), want) {
			t.Fatalf("expected %q in %q", want, textHTML)
		}
	}

	boxHTML, err := Box(BoxProps{Border: true, Tone: "base", Padding: "sm", Props: ComponentProps{Class: "rounded-box"}}, Text("Panel")).Render()
	if err != nil {
		t.Fatalf("box render failed: %v", err)
	}
	for _, want := range []string{"bg-base-100", "p-3", "border border-base-300", "rounded-box", "Panel"} {
		if !strings.Contains(string(boxHTML), want) {
			t.Fatalf("expected %q in %q", want, boxHTML)
		}
	}

	shellHTML, err := AppShell(AppShellProps{
		Sidebar: Text("Nav"),
		Flashes: Text("Flash"),
		Header:  Text("Header"),
		Content: Text("Main"),
	}).Render()
	if err != nil {
		t.Fatalf("app shell render failed: %v", err)
	}
	for _, want := range []string{`id="app"`, `id="main-content"`, "Nav", "Flash", "Header", "Main"} {
		if !strings.Contains(string(shellHTML), want) {
			t.Fatalf("expected %q in %q", want, shellHTML)
		}
	}

	hiddenHTML, err := HiddenField("id", "42").Render()
	if err != nil {
		t.Fatalf("hidden field render failed: %v", err)
	}
	for _, want := range []string{`type="hidden"`, `name="id"`, `value="42"`} {
		if !strings.Contains(string(hiddenHTML), want) {
			t.Fatalf("expected %q in %q", want, hiddenHTML)
		}
	}

	cardHTML, err := Card(CardProps{
		Title:       "Card",
		Description: "Summary",
		Actions:     ButtonComponent("Edit", ComponentProps{Variant: "ghost", Size: "sm"}),
		Gap:         "sm",
	}, Text("Body")).Render()
	if err != nil {
		t.Fatalf("card render failed: %v", err)
	}
	for _, want := range []string{"card bg-base-100 shadow-sm", "gap-2", "Card", "Summary", "btn-ghost", "Body"} {
		if !strings.Contains(string(cardHTML), want) {
			t.Fatalf("expected %q in %q", want, cardHTML)
		}
	}

	sectionHTML, err := Section(SectionProps{Title: "Section", Description: "Details"}, Text("Content")).Render()
	if err != nil {
		t.Fatalf("section render failed: %v", err)
	}
	for _, want := range []string{"space-y-4", "Section", "Details", "Content"} {
		if !strings.Contains(string(sectionHTML), want) {
			t.Fatalf("expected %q in %q", want, sectionHTML)
		}
	}
}

func TestFrontendLayoutComponentsMatchRootOutput(t *testing.T) {
	rootHTML, err := Grid(GridProps{Columns: "2", Gap: "sm"}, Text("A"), Text("B")).Render()
	if err != nil {
		t.Fatalf("root grid render failed: %v", err)
	}
	frontendHTML, err := mf.Grid(mf.GridProps{Columns: "2", Gap: "sm"}, mf.Text("A"), mf.Text("B")).Render()
	if err != nil {
		t.Fatalf("frontend grid render failed: %v", err)
	}
	if rootHTML != frontendHTML {
		t.Fatalf("expected frontend output to match root\nroot:\n%s\nfrontend:\n%s", rootHTML, frontendHTML)
	}

	rootRegionHTML, err := Region(RegionProps{ID: "tasks"}, Text("Root")).Render()
	if err != nil {
		t.Fatalf("root region render failed: %v", err)
	}
	frontendRegionHTML, err := mf.Region(mf.RegionProps{ID: "tasks"}, mf.Text("Root")).Render()
	if err != nil {
		t.Fatalf("frontend region render failed: %v", err)
	}
	if rootRegionHTML != frontendRegionHTML {
		t.Fatalf("expected frontend region output to match root\nroot:\n%s\nfrontend:\n%s", rootRegionHTML, frontendRegionHTML)
	}
}

func TestTableRowValuesConvertsValuesToCells(t *testing.T) {
	row := TableRowValues(7, "Aiko", nil, Text("Admin"))
	tableHTML, err := TableComponent(TableProps{
		Columns: []TableColumn{{Label: "ID"}, {Label: "Name"}, {Label: "Blank"}, {Label: "Role"}},
		Rows:    []TableComponentRow{row},
	}).Render()
	if err != nil {
		t.Fatalf("table render failed: %v", err)
	}
	for _, want := range []string{"7", "Aiko", "Admin"} {
		if !strings.Contains(string(tableHTML), want) {
			t.Fatalf("expected %q in %q", want, tableHTML)
		}
	}
}

func TestComponentChartRendersConfigAndFallback(t *testing.T) {
	html, err := Chart(ChartProps{
		Type:        ChartTypeBar,
		Title:       "Revenue",
		Description: "Monthly recurring revenue.",
		Labels:      []string{"Jan", "Feb"},
		Datasets: []ChartDataset{
			{Label: "Actual", Data: []float64{12.5, 18}, BackgroundColor: "#2563eb"},
		},
		Options: ChartOptions{
			BeginAtZero: true,
			XAxisLabel:  "Month",
			YAxisLabel:  "USD",
		},
		Height: 240,
	}).Render()
	if err != nil {
		t.Fatalf("chart render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`data-mrn-chart-root`,
		`data-mrn-chart`,
		`data-mrn-chart-config`,
		`"type":"bar"`,
		`"labels":["Jan","Feb"]`,
		`"beginAtZero":true`,
		`Revenue data is available in the fallback table below.`,
		`<th scope="row">Jan</th>`,
		`<td>12.5</td>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestShellIncludesChartRuntime(t *testing.T) {
	html, err := shell(`<div></div>`)
	if err != nil {
		t.Fatalf("shell render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`cdn.jsdelivr.net/npm/chart.js@4`,
		`window.mrnInitCharts`,
		`htmx:afterSwap`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in shell", want)
		}
	}
}

func TestFrontendChartComponentMatchesRootOutput(t *testing.T) {
	rootHTML, err := Chart(ChartProps{
		Type:   ChartTypeBar,
		Title:  "Orders",
		Labels: []string{"Open", "Closed"},
		Datasets: []ChartDataset{
			{Label: "Count", Data: []float64{4, 9}},
		},
	}).Render()
	if err != nil {
		t.Fatalf("root chart render failed: %v", err)
	}
	frontendHTML, err := mf.Chart(mf.ChartProps{
		Type:   mf.ChartTypeBar,
		Title:  "Orders",
		Labels: []string{"Open", "Closed"},
		Datasets: []mf.ChartDataset{
			{Label: "Count", Data: []float64{4, 9}},
		},
	}).Render()
	if err != nil {
		t.Fatalf("frontend chart render failed: %v", err)
	}
	if rootHTML != frontendHTML {
		t.Fatalf("expected frontend output to match root\nroot:\n%s\nfrontend:\n%s", rootHTML, frontendHTML)
	}
}

func TestComponentDataFrameRendersPrimitiveAndNodeValues(t *testing.T) {
	df := rdf.NewDataFrame(
		rdf.NewSeriesString("Name", nil, "Aiko", "Ken"),
		rdf.NewSeriesInt64("Age", nil, int64(42), nil),
		rdf.NewSeriesMixed("Role", nil, DivClass("badge", Text("Admin")), "Viewer"),
	)
	html, err := DataFrameComponent(df, TableProps{
		EmptyTitle:       "No rows",
		EmptyDescription: "Add rows to continue.",
	}).Render()
	if err != nil {
		t.Fatalf("dataframe render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{"Name", "Age", "Role", "Aiko", "42", `class="badge"`, "Ken", "Viewer"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestComponentDataFrameOverridesExplicitColumnsWithDataFrameNames(t *testing.T) {
	df := rdf.NewDataFrame(
		rdf.NewSeriesString("Name", nil, "Aiko"),
		rdf.NewSeriesString("Role", nil, "Admin"),
	)
	html, err := DataFrameComponent(df, TableProps{
		Columns: []TableColumn{{Label: "Display Name"}, {Label: "Team Role"}},
	}).Render()
	if err != nil {
		t.Fatalf("dataframe render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{"Name", "Role", "Aiko", "Admin"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
	for _, notWant := range []string{"Display Name", "Team Role"} {
		if strings.Contains(got, notWant) {
			t.Fatalf("did not expect %q in %q", notWant, got)
		}
	}
}

func TestComponentDataFrameChartBuildsChartDatasets(t *testing.T) {
	df := rdf.NewDataFrame(
		rdf.NewSeriesString("Month", nil, "Jan", "Feb"),
		rdf.NewSeriesInt64("Signups", nil, int64(12), int64(18)),
		rdf.NewSeriesString("Revenue", nil, "120.5", "180.25"),
	)
	html, err := DataFrameChart(df, DataFrameChartProps{
		Chart: ChartProps{
			Type:  ChartTypeBar,
			Title: "Growth",
		},
		LabelColumn: "Month",
		Series: []DataFrameChartSeries{
			{Column: "Signups", Label: "New users"},
			{Column: "Revenue"},
		},
	}).Render()
	if err != nil {
		t.Fatalf("dataframe chart render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{
		`"labels":["Jan","Feb"]`,
		`"label":"New users"`,
		`"data":[12,18]`,
		`"label":"Revenue"`,
		`"data":[120.5,180.25]`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestComponentDataFrameFromCSVUsesDataFrameGoImports(t *testing.T) {
	reader := strings.NewReader("name,role\nAiko,Admin\nKen,Viewer\n")
	node, err := DataFrameFromCSV(reader, TableProps{})
	if err != nil {
		t.Fatalf("csv import failed: %v", err)
	}
	html, err := node.Render()
	if err != nil {
		t.Fatalf("csv table render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{"name", "role", "Aiko", "Admin", "Ken", "Viewer"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestComponentDataFrameFromTSVDefaultsToTabDelimiter(t *testing.T) {
	reader := strings.NewReader("name\trole\nAiko\tAdmin\n")
	node, err := DataFrameFromTSV(reader, TableProps{})
	if err != nil {
		t.Fatalf("tsv import failed: %v", err)
	}
	html, err := node.Render()
	if err != nil {
		t.Fatalf("tsv table render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{"name", "role", "Aiko", "Admin"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}
