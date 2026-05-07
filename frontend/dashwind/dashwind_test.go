package dashwind

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	"github.com/YoshihideShirai/marionette/frontend/assets"
)

func TestShellRendersPublicNavigationAndMainRegion(t *testing.T) {
	node := Shell(ShellProps{
		Brand:             Brand{Title: "Acme Admin", Subtitle: "Operations"},
		CurrentPath:       "/",
		SearchPlaceholder: "Search accounts",
		Navigation: []NavGroup{{Label: "Menu", Items: []NavItem{
			{Label: "Dashboard", Path: "/", Icon: "▦"},
			{Label: "Leads", Path: "/leads", Icon: "▣"},
		}}},
	}, mf.Text("Hello DashWind"))
	html, err := node.Render()
	if err != nil {
		t.Fatalf("render shell: %v", err)
	}
	body := string(html)
	for _, want := range []string{"dashwind-shell", `id="dashwind-main"`, "Acme Admin", "Search accounts", `class="active" href="/"`, "Hello DashWind"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected shell to contain %q, got %q", want, body)
		}
	}
}

func TestShellGolden(t *testing.T) {
	node := Shell(ShellProps{
		DrawerID:          "account-drawer",
		MainTargetID:      "account-main",
		Brand:             Brand{Title: "Acme Admin", Subtitle: "Operations", Mark: "A", Href: "/"},
		CurrentPath:       "/reports",
		SearchPlaceholder: "Search reports",
		Navigation: []NavGroup{{Label: "Menu", Items: []NavItem{
			{Label: "Dashboard", Path: "/", Icon: "▦"},
			{Label: "Reports", Path: "/reports", Icon: "▣"},
		}}},
		User: UserMenu{
			Name:     "Ada Lovelace",
			Email:    "ada@example.com",
			Initials: "AL",
			Items:    []NavItem{{Label: "Profile", Path: "/profile"}},
		},
	}, mf.Text("Report body"))
	rendered, err := node.Render()
	if err != nil {
		t.Fatalf("render shell: %v", err)
	}
	got := strings.TrimSpace(string(rendered)) + "\n"
	goldenPath := filepath.Join("testdata", "golden", "shell.golden.html")
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch\nwant:\n%s\ngot:\n%s", want, got)
	}
	for _, want := range []string{"drawer", "navbar", "menu", `class="active" href="/reports"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected golden shell to contain %q, got %q", want, got)
		}
	}
}

func TestRenderNavigationMarksCurrentPathActive(t *testing.T) {
	node := RenderNavigation(Navigation{{Label: "Menu", Items: []NavItem{
		{Label: "Dashboard", Path: "/", Icon: "▦"},
		{Label: "Leads", Path: "/leads", Icon: "▣", Badge: "3"},
	}}}, "/leads")
	html, err := node.Render()
	if err != nil {
		t.Fatalf("render navigation: %v", err)
	}
	body := string(html)
	if !strings.Contains(body, `class="active" href="/leads"`) {
		t.Fatalf("expected current path item to have active class, got %q", body)
	}
	if strings.Contains(body, `class="active" href="/"`) {
		t.Fatalf("expected non-current item to be inactive, got %q", body)
	}
}

func TestMetricGridRendersToneAndLinks(t *testing.T) {
	node := MetricGrid(MetricGridProps{Items: []Metric{{
		Title:       "Revenue",
		Value:       "$42K",
		Description: "MRR",
		Trend:       "12% up",
		TrendTone:   ToneSuccess,
		Icon:        mf.Text("$"),
		Href:        "/analytics",
	}}})
	html, err := node.Render()
	if err != nil {
		t.Fatalf("render metric grid: %v", err)
	}
	body := string(html)
	for _, want := range []string{`href="/analytics"`, "Revenue", "$42K", "MRR", "12% up", "text-success"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected metric grid to contain %q, got %q", want, body)
		}
	}
}

func TestStatsGridAndDataTableRenderReusableWidgets(t *testing.T) {
	type person struct {
		Name string
	}
	node := mf.DivProps(mf.ElementProps{},
		StatsGrid(StatsGridProps{Items: []Stat{{Title: "Users", Value: "42", Description: "active", Figure: mf.Text("U")}}}),
		CardPanel(CardPanelProps{Title: "Table"}, DataTable(DataTableProps[person]{
			Columns: []Column[person]{{Header: "Name", Cell: func(row person) mf.Node { return mf.Text(row.Name) }}},
			Rows:    []person{{Name: "Ada"}},
		})),
	)
	html, err := node.Render()
	if err != nil {
		t.Fatalf("render widgets: %v", err)
	}
	body := string(html)
	for _, want := range []string{"stat-title", "Users", "card-title", "Table", "<th>Name</th>", "<td><span>Ada</span></td>"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected widgets to contain %q, got %q", want, body)
		}
	}
}

func TestAuthCardAndSettingsSectionRenderPresetFields(t *testing.T) {
	node := mf.DivProps(mf.ElementProps{},
		AuthCard(AuthCardProps{
			Title:       "Login",
			Description: "Access your workspace",
			Action:      "/login",
			Fields: []Field{
				{Name: "email", Label: "Email Id", Type: "email", Required: true, Help: "Use your work email."},
				{Name: "password", Label: "Password", Type: "password", Error: "Required"},
			},
			SubmitLabel: "Sign in",
			Footer:      mf.Text("Forgot Password?"),
		}),
		SettingsSection(SettingsSectionProps{
			Title:       "Account",
			Description: "Editable account fields",
			Fields:      []Field{{Label: "Role", Value: "Revenue Ops"}},
			SubmitLabel: "Save",
		}),
	)
	html, err := node.Render()
	if err != nil {
		t.Fatalf("render form presets: %v", err)
	}
	body := string(html)
	for _, want := range []string{"Login", "Access your workspace", `action="/login"`, `name="email"`, `type="email"`, `required="required"`, "Use your work email.", "input-error", "Forgot Password?", "Account", "Editable account fields", `name="role"`, `value="Revenue Ops"`, "Save"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected form presets to contain %q, got %q", want, body)
		}
	}
}

func TestResourcePageRendersActionsAndRowActions(t *testing.T) {
	type person struct {
		Email string
		Name  string
	}
	node := ResourcePage(ResourcePageProps[person]{
		Title:       "People",
		Description: "Manage people",
		Rows:        []person{{Email: "ada@example.com", Name: "Ada"}},
		Columns: []Column[person]{
			{Header: "Name", Cell: func(row person) mf.Node { return mf.Text(row.Name) }},
		},
		PrimaryAction: Action{Label: "Add", Action: "/people/add", Target: "#dashwind-main", Swap: "outerHTML", Class: "btn-primary btn-sm"},
		RowActions: func(row person) []Action {
			return []Action{{Label: "Delete", Action: "/people/delete", Target: "#dashwind-main", Swap: "outerHTML", Class: "btn-error btn-sm", Fields: map[string]string{"email": row.Email}}}
		},
	})
	html, err := node.Render()
	if err != nil {
		t.Fatalf("render resource page: %v", err)
	}
	body := string(html)
	for _, want := range []string{"People", "Manage people", `hx-post="/people/add"`, `hx-post="/people/delete"`, `name="email"`, `value="ada@example.com"`, "Delete"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected resource page to contain %q, got %q", want, body)
		}
	}
}

func TestResourcePageRendersEmptyState(t *testing.T) {
	type person struct {
		Name string
	}
	node := ResourcePage(ResourcePageProps[person]{
		Title:       "People",
		Description: "Manage people",
		Columns:     []Column[person]{{Header: "Name", Cell: func(row person) mf.Node { return mf.Text(row.Name) }}},
		EmptyState:  mf.EmptyStateProps{Title: "No people", Description: "Invite someone first."},
	})
	html, err := node.Render()
	if err != nil {
		t.Fatalf("render empty resource page: %v", err)
	}
	body := string(html)
	for _, want := range []string{"No people", "Invite someone first.", "hero bg-base-200"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected empty resource page to contain %q, got %q", want, body)
		}
	}
	if strings.Contains(body, "<table") {
		t.Fatalf("expected empty resource page to render empty state instead of table, got %q", body)
	}
}

func TestUseRegistersDashWindAssets(t *testing.T) {
	app := backend.New()
	Use(app, Options{Theme: "night", CustomCSS: ".dashwind-custom { color: red; }"})
	app.Page("/", func(ctx *backend.Context) mf.Node { return mf.Text("Dashboard") })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	for _, want := range []string{
		`href="` + assets.DaisyUICSSURL + `"`,
		"dashwind-shell .drawer-side .menu a.active",
		".dashwind-custom { color: red; }",
		`document.documentElement.setAttribute("data-theme", "night");`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected Use output to contain %q, got %q", want, body)
		}
	}
}

func TestUseCanDisableDefaultCSSAndUseAssetBasePath(t *testing.T) {
	app := backend.New()
	Use(app, Options{DisableDefaultCSS: true, AssetsBasePath: "/assets/dashwind"})
	app.Page("/", func(ctx *backend.Context) mf.Node { return mf.Text("Dashboard") })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	for _, want := range []string{
		`href="/assets/dashwind/daisyui.css"`,
		`src="/assets/dashwind/tailwindcss-browser.js"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected Use output to contain %q, got %q", want, body)
		}
	}
	if strings.Contains(body, "dashwind-shell .drawer-side .menu a.active") {
		t.Fatalf("did not expect default DashWind CSS when disabled, got %q", body)
	}
	if strings.Contains(body, assets.DaisyUICSSURL) {
		t.Fatalf("did not expect CDN DaisyUI URL when AssetsBasePath is set, got %q", body)
	}
}
