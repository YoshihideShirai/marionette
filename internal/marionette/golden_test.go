package marionette

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var updateGolden = flag.Bool("update", false, "update golden files")

func TestTemplateRenderingGolden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		node Node
	}{
		{
			name: "button",
			node: ComponentButton("Save", ComponentProps{Variant: "secondary", Size: "sm", Class: "tracking-wide"}),
		},
		{
			name: "input",
			node: ComponentInputWithOptions("start_date", "2030-01-01", InputOptions{
				Type:        "date",
				Placeholder: "Start date",
				Min:         "2024-01-01",
				Max:         "2035-12-31",
				Required:    true,
				Props:       ComponentProps{Variant: "ghost", Size: "sm", Disabled: true},
			}),
		},
		{
			name: "select",
			node: ComponentSelect("role", []SelectOption{{Label: "Admin", Value: "admin", Selected: true}, {Label: "Viewer", Value: "viewer"}}, ComponentProps{Variant: "ghost", Size: "sm"}),
		},
		{
			name: "form_field",
			node: ComponentFormField(
				ComponentInputWithOptions("name", "", InputOptions{Required: true, Props: ComponentProps{Size: "sm"}}),
				FormFieldProps{Label: "Name", Required: true, Hint: "Enter a display name.", Error: "Name is required."},
			),
		},
		{
			name: "modal_open",
			node: ComponentModal(ModalProps{
				Title:   "Delete user",
				Body:    Text("Confirm deletion"),
				Actions: ComponentButton("Delete", ComponentProps{Variant: "danger", Size: "sm"}),
				Open:    true,
			}),
		},
		{
			name: "empty_state",
			node: ComponentEmptyState(EmptyStateProps{Title: "No users", Description: "Create one first."}),
		},
		{
			name: "table",
			node: ComponentTable(TableProps{
				Columns:          []TableColumn{{Label: "Name", SortKey: "name", SortHref: "/?sort=name", SortActive: true}, {Label: "Role"}},
				Rows:             []TableComponentRow{{Cells: []Node{Text("Aiko"), DivClass("", "badge", Text("Admin"))}}},
				EmptyTitle:       "No users",
				EmptyDescription: "Create a user to get started.",
			}),
		},
		{
			name: "pagination",
			node: ComponentPagination(PaginationProps{Page: 2, TotalPages: 4, PrevHref: "/?page=1&per_page=10", NextHref: "/?page=3&per_page=10"}),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rendered, err := tc.node.Render()
			if err != nil {
				t.Fatalf("render failed: %v", err)
			}

			goldenPath := filepath.Join("testdata", "golden", tc.name+".golden.html")
			got := []byte(strings.TrimSpace(string(rendered)) + "\n")

			if *updateGolden {
				if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
					t.Fatalf("failed to update golden file: %v", err)
				}
			}

			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("failed to read golden file %q: %v", goldenPath, err)
			}

			if string(got) != string(want) {
				t.Fatalf("golden mismatch for %s\nwant:\n%s\ngot:\n%s", tc.name, want, got)
			}
		})
	}
}
