package dashwind

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	mf "github.com/YoshihideShirai/marionette/frontend"
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
