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
			name: "toast",
			node: ComponentToast(ToastProps{Title: "Saved", Description: "All changes were synced.", Icon: "✓", Props: ComponentProps{Variant: "success", Size: "sm"}}),
		},
		{
			name: "alert",
			node: ComponentAlert(AlertProps{Title: "Request failed", Description: "Try again later.", Icon: "!", Props: ComponentProps{Variant: "error", Size: "md"}}),
		},
		{
			name: "skeleton",
			node: ComponentSkeleton(SkeletonProps{Rows: 2, Props: ComponentProps{Variant: "warning", Size: "lg"}}),
		},
		{
			name: "empty_state",
			node: ComponentEmptyState(EmptyStateProps{Title: "No users", Description: "Create one first."}),
		},
		{
			name: "table",
			node: ComponentTable(TableProps{
				Columns:          []TableColumn{{Label: "Name", SortKey: "name", SortHref: "/?sort=name", SortActive: true}, {Label: "Role"}},
				Rows:             []TableComponentRow{{Cells: []Node{Text("Aiko"), DivClass("badge", Text("Admin"))}}},
				EmptyTitle:       "No users",
				EmptyDescription: "Create a user to get started.",
			}),
		},
		{
			name: "chart",
			node: ComponentChart(ChartProps{
				Type:        ChartTypeLine,
				Title:       "Weekly signups",
				Description: "New accounts by weekday.",
				Labels:      []string{"Mon", "Tue", "Wed"},
				Datasets: []ChartDataset{
					{
						Label:           "Signups",
						Data:            []float64{12, 19, 14},
						BorderColor:     "#2563eb",
						BackgroundColor: "rgba(37, 99, 235, 0.16)",
						Fill:            true,
						Tension:         0.3,
					},
				},
				Options: ChartOptions{
					BeginAtZero: true,
					YAxisLabel:  "Users",
				},
				Height: 260,
			}),
		},
		{
			name: "pagination",
			node: ComponentPagination(PaginationProps{Page: 2, TotalPages: 4, PrevHref: "/?page=1&per_page=10", NextHref: "/?page=3&per_page=10"}),
		},
		{
			name: "tabs",
			node: ComponentTabs(TabsProps{
				AriaLabel: "user sections",
				Items: []TabsItem{
					{Label: "Profile", Href: "/users/1/profile", Active: true},
					{Label: "Permissions", Href: "/users/1/permissions"},
					{Label: "Audit", Disabled: true},
				},
			}),
		},
		{
			name: "breadcrumb",
			node: ComponentBreadcrumb(BreadcrumbProps{
				Items: []BreadcrumbItem{
					{Label: "Home", Href: "/"},
					{Label: "Users", Href: "/users"},
					{Label: "Aiko", Active: true},
				},
			}),
		},
		{
			name: "textarea",
			node: ComponentTextarea("notes", "hello", TextareaOptions{
				Placeholder: "Memo",
				Rows:        4,
				Required:    true,
				Props:       ComponentProps{Variant: "ghost", Size: "sm"},
			}),
		},
		{
			name: "checkbox",
			node: ComponentCheckbox(CheckboxComponentProps{
				Name:    "active",
				Value:   "1",
				Label:   "Active user",
				Checked: true,
				Props:   ComponentProps{Size: "sm"},
			}),
		},
		{
			name: "radio_group",
			node: ComponentRadioGroup(RadioGroupComponentProps{
				Name:      "role",
				AriaLabel: "role",
				Items: []RadioItem{
					{Label: "Admin", Value: "admin", Checked: true},
					{Label: "Viewer", Value: "viewer"},
				},
				Props: ComponentProps{Size: "sm"},
			}),
		},
		{
			name: "switch",
			node: ComponentSwitch(SwitchComponentProps{
				Name:    "notify",
				Value:   "1",
				Label:   "Enable notifications",
				Checked: true,
				Props:   ComponentProps{Size: "sm"},
			}),
		},
		{
			name: "stack",
			node: ComponentStack(
				StackProps{Direction: "horizontal", Gap: "sm", Align: "center", Justify: "between", Wrap: true, Props: ComponentProps{Class: "w-full"}},
				Text("Left"),
				ComponentButton("Right", ComponentProps{Variant: "secondary", Size: "sm"}),
			),
		},
		{
			name: "grid",
			node: ComponentGrid(
				GridProps{Columns: "3", Gap: "lg"},
				DivClass("card bg-base-100 p-4", Text("One")),
				DivClass("card bg-base-100 p-4", Text("Two")),
				DivClass("card bg-base-100 p-4", Text("Three")),
			),
		},
		{
			name: "split",
			node: ComponentSplit(SplitProps{
				Main:            DivClass("card bg-base-100 p-4", Text("Main")),
				Aside:           DivClass("card bg-base-100 p-4", Text("Aside")),
				AsideWidth:      "md",
				ReverseOnMobile: true,
				Gap:             "lg",
			}),
		},
		{
			name: "page_header",
			node: ComponentPageHeader(PageHeaderProps{
				Title:       "Users",
				Description: "Manage account records.",
				Actions:     ComponentButton("Create", ComponentProps{Size: "sm"}),
			}),
		},
		{
			name: "container",
			node: ComponentContainer(
				ContainerProps{MaxWidth: "md", Padding: "sm", Centered: true},
				Text("Contained"),
			),
		},
		{
			name: "card",
			node: ComponentCard(
				CardProps{
					Title:       "Summary",
					Description: "Current workspace status.",
					Actions:     ComponentButton("Edit", ComponentProps{Variant: "ghost", Size: "sm"}),
				},
				Text("Ready"),
			),
		},
		{
			name: "section",
			node: ComponentSection(
				SectionProps{Title: "Details", Description: "Supporting information."},
				Text("Section body"),
			),
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
