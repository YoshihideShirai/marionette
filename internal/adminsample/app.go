package adminsample

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

//go:embed assets/*
var embeddedAssets embed.FS

type order struct {
	ID       string
	Customer string
	Plan     string
	Amount   int
	Risk     string
	Status   string
}

var seedOrders = []order{
	{ID: "ORD-1042", Customer: "Acme Robotics", Plan: "Enterprise", Amount: 128000, Risk: "Low", Status: "Active"},
	{ID: "ORD-1043", Customer: "Northwind Health", Plan: "Growth", Amount: 64000, Risk: "Medium", Status: "Review"},
	{ID: "ORD-1044", Customer: "Riverline Bank", Plan: "Enterprise", Amount: 98000, Risk: "High", Status: "Blocked"},
	{ID: "ORD-1045", Customer: "Sora Logistics", Plan: "Starter", Amount: 22000, Risk: "Low", Status: "Active"},
	{ID: "ORD-1046", Customer: "Bluepeak Energy", Plan: "Growth", Amount: 57000, Risk: "Medium", Status: "Review"},
	{ID: "ORD-1047", Customer: "Kite Retail", Plan: "Starter", Amount: 18000, Risk: "Low", Status: "Active"},
}

func BuildApp() *mb.App {
	app := mb.New()
	app.Set("orders", seedOrders)
	app.Set("selectedStatus", "all")
	app.Set("loggedIn", false)
	app.Set("authError", "")
	app.Set("flash", "")

	assetsFS, err := fs.Sub(embeddedAssets, "assets")
	if err == nil {
		app.Assets("/assets", assetsFS)
		app.AddStylesheet(app.Asset("admin-sample.css"))
	}

	app.Page("/", func(ctx *mb.Context) mf.Node {
		if !ctx.Get("loggedIn").(bool) {
			return loginPage(ctx.Get("authError").(string))
		}
		return dashboardFromState(ctx)
	}, mb.WithTitle("Revenue Ops Console"))

	app.Action("auth/login", func(ctx *mb.Context) mf.Node {
		email := strings.TrimSpace(ctx.FormValue("email"))
		password := strings.TrimSpace(ctx.FormValue("password"))
		if email == "ops@example.com" && password == "marionette" {
			ctx.Set("loggedIn", true)
			ctx.Set("authError", "")
			ctx.Set("flash", "Welcome back. Live controls are ready.")
			return dashboardFromState(ctx)
		}
		ctx.Set("loggedIn", false)
		ctx.Set("authError", "Invalid credentials. Try ops@example.com / marionette")
		return loginPage(ctx.Get("authError").(string))
	})

	app.Action("auth/logout", func(ctx *mb.Context) mf.Node {
		ctx.Set("loggedIn", false)
		ctx.Set("authError", "")
		ctx.Set("flash", "")
		return loginPage("")
	})

	app.Action("orders/filter", func(ctx *mb.Context) mf.Node {
		if !ctx.Get("loggedIn").(bool) {
			return mf.Alert(mf.AlertProps{Title: "Session expired", Description: "Please sign in again.", Props: mf.ComponentProps{Variant: "warning"}})
		}
		status := strings.TrimSpace(ctx.FormValue("status"))
		if status == "" {
			status = "all"
		}
		ctx.Set("selectedStatus", status)
		ctx.Set("flash", fmt.Sprintf("Filter applied: %s", status))
		return dashboardBody(ctx)
	})

	app.Action("orders/toggle-status", func(ctx *mb.Context) mf.Node {
		if !ctx.Get("loggedIn").(bool) {
			return mf.Alert(mf.AlertProps{Title: "Session expired", Description: "Please sign in again.", Props: mf.ComponentProps{Variant: "warning"}})
		}
		id := strings.TrimSpace(ctx.FormValue("id"))
		orders := ctx.Get("orders").([]order)
		for i := range orders {
			if orders[i].ID != id {
				continue
			}
			if orders[i].Status == "Blocked" {
				orders[i].Status = "Active"
				orders[i].Risk = "Medium"
			} else {
				orders[i].Status = "Blocked"
				orders[i].Risk = "High"
			}
			ctx.Set("flash", fmt.Sprintf("%s moved to %s", id, orders[i].Status))
			break
		}
		ctx.Set("orders", orders)
		return dashboardBody(ctx)
	})

	return app
}

func dashboardFromState(ctx *mb.Context) mf.Node {
	return mf.Container(mf.ContainerProps{MaxWidth: "7xl", Centered: true, Props: mf.ComponentProps{Class: "ops-shell py-8"}},
		mf.Region(mf.RegionProps{ID: "app-body"}, dashboardBody(ctx)),
	)
}

func dashboardBody(ctx *mb.Context) mf.Node {
	selectedStatus := ctx.Get("selectedStatus").(string)
	orders := filteredOrders(ctx.Get("orders").([]order), selectedStatus)
	flash := ctx.Get("flash").(string)

	nodes := []mf.Node{
		mf.PageHeader(mf.PageHeaderProps{Title: "Revenue Ops Console", Description: "Prioritize opportunities and monitor risk from a single screen"}),
		mf.ActionForm(mf.ActionFormProps{Action: "/auth/logout", Target: "#app-body", Swap: "outerHTML", Props: mf.ComponentProps{Class: "self-end"}}, mf.SubmitButton("Sign out", mf.ComponentProps{Variant: "ghost"})),
	}
	if flash != "" {
		nodes = append(nodes, mf.Alert(mf.AlertProps{Title: "Live update", Description: flash, Props: mf.ComponentProps{Variant: "success", Class: "ops-live"}}))
	}
	nodes = append(nodes,
		mf.Alert(mf.AlertProps{Title: "Morning briefing", Description: "1 high-risk deal needs legal review before renewal.", Props: mf.ComponentProps{Variant: "warning"}}),
		mf.ActionForm(mf.ActionFormProps{Action: "/orders/filter", Target: "#app-body", Swap: "outerHTML", Props: mf.ComponentProps{Class: "ops-toolbar flex flex-wrap items-end gap-3"}},
			mf.FormRow(mf.FormRowProps{ID: "status", Label: "Deal status", Control: mf.Select(mf.SelectFieldProps{ID: "status", Name: "status", Options: statusOptions(selectedStatus)})}),
			mf.SubmitButton("Apply filter", mf.ComponentProps{}),
		),
		mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, summaryCards(orders), pipelineHealth(orders), pipelineChart(orders), ordersTable(orders)),
	)
	return mf.Stack(mf.StackProps{Direction: "column", Gap: "6"}, nodes...)
}

func loginPage(authError string) mf.Node {
	children := []mf.Node{
		mf.PageHeader(mf.PageHeaderProps{Title: "Revenue Ops Console", Description: "Sign in to access the admin dashboard"}),
		mf.Card(mf.CardProps{}, mf.Stack(mf.StackProps{Direction: "column", Gap: "3"}, mf.Text("Demo credentials"), mf.Text("Email: ops@example.com"), mf.Text("Password: marionette"))),
	}
	if authError != "" {
		children = append(children, mf.Alert(mf.AlertProps{Title: "Login failed", Description: authError, Props: mf.ComponentProps{Variant: "error"}}))
	}
	children = append(children,
		mf.ActionForm(mf.ActionFormProps{Action: "/auth/login", Target: "#app-body", Swap: "outerHTML", Props: mf.ComponentProps{Class: "space-y-3"}},
			mf.FormRow(mf.FormRowProps{ID: "email", Label: "Email", Required: true, Control: mf.TextField(mf.TextFieldProps{ID: "email", Name: "email", Type: "email", Placeholder: "ops@example.com", Required: true})}),
			mf.FormRow(mf.FormRowProps{ID: "password", Label: "Password", Required: true, Control: mf.TextField(mf.TextFieldProps{ID: "password", Name: "password", Type: "password", Placeholder: "••••••••", Required: true})}),
			mf.SubmitButton("Sign in", mf.ComponentProps{}),
		),
	)
	return mf.Container(mf.ContainerProps{MaxWidth: "lg", Centered: true, Props: mf.ComponentProps{Class: "ops-shell py-12"}}, mf.Region(mf.RegionProps{ID: "app-body"}, mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, children...)))
}

func summaryCards(orders []order) mf.Node {
	total := 0
	highRisk := 0
	for _, o := range orders {
		total += o.Amount
		if o.Risk == "High" {
			highRisk++
		}
	}
	return mf.Grid(mf.GridProps{Columns: "md:grid-cols-3", Gap: "4"},
		statCard("Visible deals", fmt.Sprintf("%d", len(orders)), "Deals after filters"),
		statCard("Pipeline value", fmt.Sprintf("$%s", formatNumber(total)), "Total projected revenue"),
		statCard("High-risk deals", fmt.Sprintf("%d", highRisk), "Needs escalation"),
	)
}

func statCard(title, value, desc string) mf.Node {
	return mf.Card(mf.CardProps{Props: mf.ComponentProps{Class: "ops-card"}}, mf.Stack(mf.StackProps{Direction: "column", Gap: "1"}, mf.Text(title), mf.H3(mf.Text(value)), mf.Text(desc)))
}

func pipelineHealth(orders []order) mf.Node {
	active := 0
	for _, o := range orders {
		if o.Status == "Active" {
			active++
		}
	}
	value := 0.0
	if len(orders) > 0 {
		value = float64(active) / float64(len(orders)) * 100
	}
	return mf.Section(mf.SectionProps{Title: "Pipeline health"}, mf.Progress(mf.ProgressProps{Label: "Active ratio", Value: value, Max: 100, ShowValue: true}))
}

func pipelineChart(orders []order) mf.Node {
	statusOrder := []string{"Active", "Review", "Blocked"}
	counts := map[string]float64{"Active": 0, "Review": 0, "Blocked": 0}
	for _, o := range orders {
		counts[o.Status]++
	}
	values := make([]float64, 0, len(statusOrder))
	for _, status := range statusOrder {
		values = append(values, counts[status])
	}
	return mf.Section(mf.SectionProps{Title: "Deal distribution"},
		mf.Chart(mf.ChartProps{
			Type:      mf.ChartTypeBar,
			Title:     "Deals by status",
			Labels:    statusOrder,
			Height:    260,
			AriaLabel: "Bar chart showing deal counts by status",
			Datasets: []mf.ChartDataset{{
				Label:           "Deals",
				Data:            values,
				BackgroundColor: "rgba(56, 189, 248, 0.55)",
				BorderColor:     "rgba(14, 165, 233, 1)",
			}},
			Options: mf.ChartOptions{BeginAtZero: true, HideLegend: true, YAxisLabel: "Count"},
		}),
	)
}

func ordersTable(orders []order) mf.Node {
	rows := make([]mf.TableComponentRow, 0, len(orders))
	for _, o := range orders {
		label := "Escalate"
		if o.Status == "Blocked" {
			label = "Re-open"
		}
		action := mf.ActionForm(mf.ActionFormProps{Action: "/orders/toggle-status", Target: "#app-body", Swap: "outerHTML"},
			mf.TextField(mf.TextFieldProps{ID: "id-" + o.ID, Name: "id", Value: o.ID, Type: "hidden"}),
			mf.SubmitButton(label, mf.ComponentProps{Size: "sm"}),
		)
		rows = append(rows, mf.TableRowValues(o.ID, o.Customer, o.Plan, "$"+formatNumber(o.Amount), o.Risk, o.Status, action))
	}
	return mf.Table(mf.TableProps{Columns: []mf.TableColumn{{Label: "Deal"}, {Label: "Customer"}, {Label: "Plan"}, {Label: "ARR"}, {Label: "Risk"}, {Label: "Status"}, {Label: "Action"}}, Rows: rows, EmptyTitle: "No deals match the selected filters"})
}

func statusOptions(selected string) []mf.SelectOption {
	values := []string{"all", "Active", "Review", "Blocked"}
	opts := make([]mf.SelectOption, 0, len(values))
	for _, v := range values {
		opts = append(opts, mf.SelectOption{Label: v, Value: v, Selected: v == selected})
	}
	return opts
}

func filteredOrders(orders []order, status string) []order {
	if status == "" || status == "all" {
		return orders
	}
	filtered := make([]order, 0, len(orders))
	for _, o := range orders {
		if o.Status == status {
			filtered = append(filtered, o)
		}
	}
	return filtered
}

func formatNumber(v int) string {
	s := fmt.Sprintf("%d", v)
	if len(s) <= 3 {
		return s
	}
	var out []byte
	pre := len(s) % 3
	if pre > 0 {
		out = append(out, s[:pre]...)
		if len(s) > pre {
			out = append(out, ',')
		}
	}
	for i := pre; i < len(s); i += 3 {
		out = append(out, s[i:i+3]...)
		if i+3 < len(s) {
			out = append(out, ',')
		}
	}
	return string(out)
}
