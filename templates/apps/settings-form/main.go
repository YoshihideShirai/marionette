package main

import (
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

type settings struct {
	ProductName string
	NotifyEmail string
	Region      string
	Summary     string
}

func main() {
	app := buildApp()
	if err := app.Run("127.0.0.1:8093"); err != nil {
		panic(err)
	}
}

func buildApp() *mb.App {
	app := mb.New()
	dw.Use(app, dw.Options{Theme: "corporate"})
	app.SetGlobal("settings", settings{
		ProductName: "Marionette Console",
		NotifyEmail: "ops@example.com",
		Region:      "ap-northeast-1",
		Summary:     "Internal operations workspace.",
	})

	app.Page("/", page, mb.WithTitle("Settings Form Template"))
	app.Action("settings/save", saveSettings)

	return app
}

func page(ctx *mb.Context) mf.Node {
	current := ctx.GetGlobal("settings").(settings)
	body := mf.DivProps(mf.ElementProps{Class: "space-y-6"},
		dw.PageHeader(dw.PageHeaderProps{
			Title:       "Settings",
			Description: "A form-heavy layout with validation feedback, preview, and helper panels.",
		}),
		settingsWorkspace(current, ""),
	)
	return appShell("/", body)
}

func saveSettings(ctx *mb.Context) mf.Node {
	next := settings{
		ProductName: strings.TrimSpace(ctx.FormValue("product_name")),
		NotifyEmail: strings.TrimSpace(ctx.FormValue("notify_email")),
		Region:      strings.TrimSpace(ctx.FormValue("region")),
		Summary:     strings.TrimSpace(ctx.FormValue("summary")),
	}
	if next.ProductName == "" {
		current := ctx.GetGlobal("settings").(settings)
		return settingsWorkspace(current, "Product name is required.")
	}
	ctx.SetGlobal("settings", next)
	return settingsWorkspace(next, "")
}

func settingsWorkspace(current settings, productNameError string) mf.Node {
	return mf.Region(mf.RegionProps{ID: "settings-workspace"},
		mf.Split(mf.SplitProps{
			Main: settingsForm(current, productNameError),
			Aside: mf.Stack(mf.StackProps{Direction: "column", Gap: "4"},
				previewCard(current),
				dw.CardPanel(dw.CardPanelProps{Title: "Release checklist", Description: "Useful adjacent content for configuration pages."},
					mf.UlProps(mf.ElementProps{Class: "space-y-3 text-sm"},
						checkItem("Validate required fields", true),
						checkItem("Preview saved values", true),
						checkItem("Add persistence layer", false),
					),
				),
			),
			AsideWidth: "sm",
			Gap:        "6",
		}),
	)
}

func settingsForm(current settings, productNameError string) mf.Node {
	return dw.CardPanel(dw.CardPanelProps{Title: "Workspace settings", Description: "Saving returns a replacement for the settings workspace."},
		mf.ActionForm(mf.ActionFormProps{
			Action: "/settings/save",
			Target: "#settings-workspace",
			Swap:   "outerHTML",
			Props:  mf.ComponentProps{Class: "space-y-4"},
		},
			mf.Grid(mf.GridProps{Columns: "2", Gap: "lg"},
				mf.FormRow(mf.FormRowProps{
					ID:       "product-name",
					Label:    "Product name",
					Error:    productNameError,
					Required: true,
					Control: mf.TextField(mf.TextFieldProps{
						ID:       "product-name",
						Name:     "product_name",
						Value:    current.ProductName,
						Required: true,
						Error:    productNameError,
					}),
				}),
				mf.FormRow(mf.FormRowProps{
					ID:    "notify-email",
					Label: "Notification email",
					Control: mf.TextField(mf.TextFieldProps{
						ID:    "notify-email",
						Name:  "notify_email",
						Type:  "email",
						Value: current.NotifyEmail,
					}),
				}),
			),
			mf.FormRow(mf.FormRowProps{
				ID:    "region",
				Label: "Region",
				Control: mf.Select(mf.SelectFieldProps{
					ID:   "region",
					Name: "region",
					Options: []mf.SelectOption{
						{Label: "Tokyo", Value: "ap-northeast-1", Selected: current.Region == "ap-northeast-1"},
						{Label: "Oregon", Value: "us-west-2", Selected: current.Region == "us-west-2"},
						{Label: "Frankfurt", Value: "eu-central-1", Selected: current.Region == "eu-central-1"},
					},
				}),
			}),
			mf.FormRow(mf.FormRowProps{
				ID:    "summary",
				Label: "Summary",
				Control: mf.Textarea(mf.TextareaProps{
					ID:    "summary",
					Name:  "summary",
					Value: current.Summary,
					Rows:  4,
				}),
			}),
			mf.Actions(mf.ActionsProps{Align: "end", Wrap: true},
				mf.Button("Reset", mf.ComponentProps{Class: "btn-outline"}),
				mf.SubmitButton("Save settings", mf.ComponentProps{Class: "btn-primary"}),
			),
		),
	)
}

func previewCard(current settings) mf.Node {
	return dw.CardPanel(dw.CardPanelProps{Title: "Preview", Description: "Current server-side settings."},
		mf.DescriptionListProps(mf.ElementProps{Class: "space-y-4 text-sm"},
			descriptionItem("Product", current.ProductName),
			descriptionItem("Email", current.NotifyEmail),
			descriptionItem("Region", current.Region),
			descriptionItem("Summary", current.Summary),
		),
	)
}

func checkItem(label string, done bool) mf.Node {
	className := "badge-outline"
	text := "Pending"
	if done {
		className = "badge-success"
		text = "Done"
	}
	return mf.LiProps(mf.ElementProps{Class: "flex items-center justify-between gap-3 rounded-box bg-base-200 px-3 py-2"},
		mf.SpanProps(mf.ElementProps{}, mf.Text(label)),
		mf.Badge(mf.BadgeProps{Label: text, Props: mf.ComponentProps{Class: className}}),
	)
}

func descriptionItem(label, value string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "grid gap-1 border-b border-base-300 pb-3 last:border-b-0 last:pb-0"},
		mf.DescriptionTermProps(mf.ElementProps{Class: "font-medium text-base-content/60"}, mf.Text(label)),
		mf.DescriptionDetailsProps(mf.ElementProps{}, mf.Text(value)),
	)
}

func appShell(currentPath string, body mf.Node) mf.Node {
	return dw.Shell(dw.ShellProps{
		Brand:             dw.Brand{Title: "Config", Subtitle: "Settings template", Mark: "C"},
		CurrentPath:       currentPath,
		SearchPlaceholder: "Search settings",
		Navigation: []dw.NavGroup{{
			Label: "Workspace",
			Items: []dw.NavItem{{Path: "/", Label: "Settings", Icon: "S"}, {Path: "#", Label: "Audit log", Icon: "A", Badge: "Soon"}},
		}},
		User: dw.UserMenu{Name: "Template User", Email: "config@example.com", Initials: "TU"},
	}, body)
}
