package dashwind

import (
	"strings"
	"testing"

	mf "github.com/YoshihideShirai/marionette/frontend"
)

func TestShellRendersPublicNavigationAndMainRegion(t *testing.T) {
	node := Shell(ShellProps{
		CurrentTitle:      "Dashboard",
		BrandTitle:        "Acme Admin",
		BrandSubtitle:     "Operations",
		SearchPlaceholder: "Search accounts",
		NavGroups: []NavGroup{{Label: "Menu", Items: []NavItem{
			{Label: "Dashboard", Href: "/", Icon: "▦"},
			{Label: "Leads", Href: "/leads", Icon: "▣"},
		}}},
		Content: mf.Text("Hello DashWind"),
	})
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
