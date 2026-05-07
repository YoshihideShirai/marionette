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
			{Label: "Dashboard", Href: "/", Icon: "▦"},
			{Label: "Leads", Href: "/leads", Icon: "▣"},
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
			{Label: "Dashboard", Href: "/", Icon: "▦"},
			{Label: "Reports", Href: "/reports", Icon: "▣"},
		}}},
		User: UserMenu{
			Name:     "Ada Lovelace",
			Email:    "ada@example.com",
			Initials: "AL",
			Items:    []NavItem{{Label: "Profile", Href: "/profile"}},
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

func TestStatsGridAndDataTableRenderReusableWidgets(t *testing.T) {
	node := mf.DivProps(mf.ElementProps{},
		StatsGrid(StatsGridProps{Items: []Stat{{Title: "Users", Value: "42", Description: "active", Figure: mf.Text("U")}}}),
		CardPanel(CardPanelProps{Title: "Table"}, DataTable(DataTableProps{Headers: []string{"Name"}, Rows: [][]mf.Node{{mf.Text("Ada")}}})),
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
