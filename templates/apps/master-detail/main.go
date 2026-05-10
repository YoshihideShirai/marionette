package main

import (
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

type customer struct {
	ID      string
	Name    string
	Company string
	Plan    string
	Health  string
	Notes   string
}

var customers = []customer{
	{ID: "c-1001", Name: "Aki Tanaka", Company: "Northstar Labs", Plan: "Enterprise", Health: "Healthy", Notes: "Preparing the Q2 rollout."},
	{ID: "c-1002", Name: "Mina Sato", Company: "Blue Harbor", Plan: "Growth", Health: "Watch", Notes: "Wants tighter billing exports."},
	{ID: "c-1003", Name: "Ren Ito", Company: "Fieldworks", Plan: "Starter", Health: "Onboarding", Notes: "Evaluating admin permissions."},
}

func main() {
	app := mb.New()
	dw.Use(app, dw.Options{Theme: "corporate"})
	app.SetGlobal("selected_customer", customers[0].ID)

	app.Page("/", page, mb.WithTitle("Master Detail Template"))
	app.Action("customers/select", selectCustomer)

	if err := app.Run("127.0.0.1:8094"); err != nil {
		panic(err)
	}
}

func page(ctx *mb.Context) mf.Node {
	selected := selectedCustomer(ctx)
	body := mf.DivProps(mf.ElementProps{Class: "space-y-6"},
		dw.PageHeader(dw.PageHeaderProps{
			Title:       "Customers",
			Description: "A master-detail layout with list actions and a focused detail region.",
			Actions:     mf.Button("New customer", mf.ComponentProps{Class: "btn-primary btn-sm"}),
		}),
		dw.MetricGrid(dw.MetricGridProps{Items: []dw.Metric{
			{Title: "Accounts", Value: "3", Description: "Demo records", Trend: "In memory", TrendTone: dw.ToneNeutral, Icon: metricIcon("A")},
			{Title: "Enterprise", Value: "1", Description: "High touch", Trend: "Expansion ready", TrendTone: dw.ToneSuccess, Icon: metricIcon("E")},
			{Title: "Watch", Value: "1", Description: "Needs follow-up", Trend: "Visible in detail", TrendTone: dw.ToneWarning, Icon: metricIcon("W")},
			{Title: "Region", Value: "JP", Description: "Template data", Trend: "Customize freely", TrendTone: dw.ToneNeutral, Icon: metricIcon("R")},
		}}),
		mf.Split(mf.SplitProps{
			Main:            mf.Region(mf.RegionProps{ID: "customer-detail"}, detailPanel(selected)),
			Aside:           customerList(selected.ID),
			AsideWidth:      "sm",
			Gap:             "6",
			ReverseOnMobile: true,
		}),
	)
	return appShell("/", body)
}

func selectCustomer(ctx *mb.Context) mf.Node {
	id := strings.TrimSpace(ctx.FormValue("id"))
	if findCustomer(id).ID != "" {
		ctx.SetGlobal("selected_customer", id)
	}
	return detailPanel(selectedCustomer(ctx))
}

func customerList(selectedID string) mf.Node {
	items := make([]mf.Node, 0, len(customers))
	for _, c := range customers {
		className := "btn btn-ghost h-auto w-full justify-start px-3 py-3 text-left"
		if c.ID == selectedID {
			className = "btn btn-primary h-auto w-full justify-start px-3 py-3 text-left"
		}
		items = append(items,
			mf.ActionForm(mf.ActionFormProps{Action: "/customers/select", Target: "#customer-detail", Swap: "innerHTML"},
				mf.HiddenField("id", c.ID),
				mf.SubmitButton(c.Name+" - "+c.Company, mf.ComponentProps{Class: className}),
			),
		)
	}
	return dw.CardPanel(dw.CardPanelProps{Title: "Customer list", Description: "Select a row to replace the detail panel."},
		mf.Stack(mf.StackProps{Direction: "column", Gap: "2"}, items...),
	)
}

func detailPanel(c customer) mf.Node {
	if c.ID == "" {
		return dw.CardPanel(dw.CardPanelProps{Title: "Customer detail"},
			mf.EmptyState(mf.EmptyStateProps{Title: "Select a customer"}),
		)
	}
	return dw.CardPanel(dw.CardPanelProps{
		Title:       c.Name,
		Description: c.Company,
		Actions: mf.Actions(mf.ActionsProps{Align: "end", Wrap: true},
			mf.Badge(mf.BadgeProps{Label: c.Plan, Props: mf.ComponentProps{Class: "badge-primary"}}),
			healthBadge(c.Health),
		),
	},
		mf.DescriptionListProps(mf.ElementProps{Class: "grid gap-4 md:grid-cols-2"},
			descriptionItem("Customer ID", c.ID),
			descriptionItem("Plan", c.Plan),
			descriptionItem("Company", c.Company),
			descriptionItem("Health", c.Health),
			descriptionItem("Notes", c.Notes),
		),
		mf.DivProps(mf.ElementProps{Class: "mt-6 grid gap-4 md:grid-cols-3"},
			summaryTile("ARR", "$84K"),
			summaryTile("Seats", "42"),
			summaryTile("Renewal", "Aug 12"),
		),
	)
}

func selectedCustomer(ctx *mb.Context) customer {
	id, _ := ctx.GetGlobal("selected_customer").(string)
	return findCustomer(id)
}

func findCustomer(id string) customer {
	for _, c := range customers {
		if c.ID == id {
			return c
		}
	}
	return customer{}
}

func healthBadge(health string) mf.Node {
	className := "badge-info"
	if health == "Healthy" {
		className = "badge-success"
	}
	if health == "Watch" {
		className = "badge-warning"
	}
	return mf.Badge(mf.BadgeProps{Label: health, Props: mf.ComponentProps{Class: className}})
}

func summaryTile(label, value string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "rounded-box bg-base-200 p-4"},
		mf.PProps(mf.ElementProps{Class: "text-sm text-base-content/60"}, mf.Text(label)),
		mf.DivProps(mf.ElementProps{Class: "mt-1 text-2xl font-bold"}, mf.Text(value)),
	)
}

func descriptionItem(label, value string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "space-y-1 rounded-box border border-base-300 p-4"},
		mf.DescriptionTermProps(mf.ElementProps{Class: "text-sm font-medium text-base-content/60"}, mf.Text(label)),
		mf.DescriptionDetailsProps(mf.ElementProps{Class: "text-base"}, mf.Text(value)),
	)
}

func appShell(currentPath string, body mf.Node) mf.Node {
	return dw.Shell(dw.ShellProps{
		Brand:             dw.Brand{Title: "Accounts", Subtitle: "Master-detail template", Mark: "A"},
		CurrentPath:       currentPath,
		SearchPlaceholder: "Search customers",
		Navigation: []dw.NavGroup{{
			Label: "Workspace",
			Items: []dw.NavItem{{Path: "/", Label: "Customers", Icon: "C"}, {Path: "#", Label: "Segments", Icon: "S", Badge: "Soon"}},
		}},
		User: dw.UserMenu{Name: "Template User", Email: "accounts@example.com", Initials: "TU"},
	}, body)
}

func metricIcon(label string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "grid h-10 w-10 place-items-center rounded-box bg-primary/10 font-bold text-primary"}, mf.Text(label))
}
