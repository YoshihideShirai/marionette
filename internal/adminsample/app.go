package adminsample

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	du "github.com/YoshihideShirai/marionette/frontend/daisyui"
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
	app.Set("currentPage", "overview")

	assetsFS, err := fs.Sub(embeddedAssets, "assets")
	if err == nil {
		app.Assets("/assets", assetsFS)
	}

	app.Page("/", func(ctx *mb.Context) mf.Node {
		if !ctx.Get("loggedIn").(bool) {
			return loginPage(ctx.Get("authError").(string))
		}
		ctx.Set("currentPage", "overview")
		return dashboardFromState(ctx, "overview")
	}, mb.WithTitle("Admin Sample"))

	app.Page("/pipeline", func(ctx *mb.Context) mf.Node {
		if !ctx.Get("loggedIn").(bool) {
			return loginPage(ctx.Get("authError").(string))
		}
		ctx.Set("currentPage", "pipeline")
		return dashboardFromState(ctx, "pipeline")
	}, mb.WithTitle("Pipeline - Admin Sample"))

	app.Page("/playbooks", func(ctx *mb.Context) mf.Node {
		if !ctx.Get("loggedIn").(bool) {
			return loginPage(ctx.Get("authError").(string))
		}
		ctx.Set("currentPage", "playbooks")
		return dashboardFromState(ctx, "playbooks")
	}, mb.WithTitle("Playbooks - Admin Sample"))

	app.Action("auth/login", func(ctx *mb.Context) mf.Node {
		email := strings.TrimSpace(ctx.FormValue("email"))
		password := strings.TrimSpace(ctx.FormValue("password"))
		if email == "ops@example.com" && password == "marionette" {
			ctx.Set("loggedIn", true)
			ctx.Set("authError", "")
			ctx.Set("flash", "Signed in successfully")
			return dashboardFromState(ctx, "overview")
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
			return dashboardMainContent(sessionExpiredAlert())
		}
		status := strings.TrimSpace(ctx.FormValue("status"))
		if status == "" {
			status = "all"
		}
		ctx.Set("selectedStatus", status)
		ctx.Set("flash", fmt.Sprintf("Filter applied: %s", status))
		return dashboardMainContent(dashboardBody(ctx, ctx.Get("currentPage").(string)))
	})

	app.Action("orders/toggle-status", func(ctx *mb.Context) mf.Node {
		if !ctx.Get("loggedIn").(bool) {
			return dashboardMainContent(sessionExpiredAlert())
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
			ctx.Set("flash", fmt.Sprintf("%s -> %s", id, orders[i].Status))
			break
		}
		ctx.Set("orders", orders)
		return dashboardMainContent(dashboardBody(ctx, ctx.Get("currentPage").(string)))
	})

	return app
}

func dashboardFromState(ctx *mb.Context, currentPage string) mf.Node {
	return du.Region(mf.RegionProps{ID: "app-body"},
		du.DrawerLayout("my-drawer-2", topbar(currentPage), du.Container(mf.ContainerProps{MaxWidth: "7xl", Centered: true, Props: mf.ComponentProps{Class: "py-6"}}, dashboardMainContent(dashboardBody(ctx, currentPage))), sidebarItems(currentPage)),
	)
}

func topbar(currentPage string) mf.Node {
	return du.DrawerNavbar("my-drawer-2", "admin-sample", []mf.Node{
		navLink("Overview", "/", currentPage == "overview"),
		navLink("Pipeline", "/pipeline", currentPage == "pipeline"),
	})
}

func dashboardMainContent(content mf.Node) mf.Node {
	return du.Region(mf.RegionProps{ID: "main-content"}, content)
}

func sessionExpiredAlert() mf.Node {
	return du.Alert("Session expired", "Please sign in again.", mf.ComponentProps{Variant: "warning"})
}

func dashboardBody(ctx *mb.Context, currentPage string) mf.Node {
	selectedStatus := ctx.Get("selectedStatus").(string)
	orders := filteredOrders(ctx.Get("orders").([]order), selectedStatus)
	flash := ctx.Get("flash").(string)
	children := []mf.Node{
		du.PageHeader(mf.PageHeaderProps{Title: "Dashboard", Description: "Default daisyUI style admin layout"}),
		du.ActionForm(mf.ActionFormProps{Action: "/auth/logout", Target: "#app-body", Swap: "outerHTML"}, du.Button("Sign out", mf.ComponentProps{Variant: "outline"})),
	}
	if flash != "" {
		children = append(children, du.Alert("Updated", flash, mf.ComponentProps{Variant: "info"}))
	}
	children = append(children,
		du.ActionForm(mf.ActionFormProps{Action: "/orders/filter", Target: "#main-content", Swap: "outerHTML", Props: mf.ComponentProps{Class: "card bg-base-100 border border-base-300 p-4"}},
			mf.FormRow(mf.FormRowProps{ID: "status", Label: "Status", Control: mf.Select(mf.SelectFieldProps{ID: "status", Name: "status", Options: statusOptions(selectedStatus)})}),
			du.Button("Apply", mf.ComponentProps{Variant: "primary"}),
		),
		pageContent(currentPage, orders),
	)
	return du.Stack(mf.StackProps{Direction: "column", Gap: "4"}, children...)
}

func loginPage(authError string) mf.Node {
	children := []mf.Node{
		du.PageHeader(mf.PageHeaderProps{Title: "Admin Login", Description: "Use demo account to continue"}),
		du.Alert("Demo credentials", "ops@example.com / marionette", mf.ComponentProps{Variant: "info"}),
	}
	if authError != "" {
		children = append(children, du.Alert("Login failed", authError, mf.ComponentProps{Variant: "error"}))
	}
	children = append(children,
		du.Card("", "", nil, []mf.Node{du.ActionForm(mf.ActionFormProps{Action: "/auth/login", Target: "#app-body", Swap: "outerHTML", Props: mf.ComponentProps{Class: "space-y-3"}},
			mf.FormRow(mf.FormRowProps{ID: "email", Label: "Email", Required: true, Control: mf.TextField(mf.TextFieldProps{ID: "email", Name: "email", Type: "email", Placeholder: "ops@example.com", Required: true})}),
			mf.FormRow(mf.FormRowProps{ID: "password", Label: "Password", Required: true, Control: mf.TextField(mf.TextFieldProps{ID: "password", Name: "password", Type: "password", Placeholder: "••••••••", Required: true})}),
			du.Button("Sign in", mf.ComponentProps{Variant: "primary"}),
		)}, mf.ComponentProps{}),
	)
	return du.Container(mf.ContainerProps{MaxWidth: "lg", Centered: true, Props: mf.ComponentProps{Class: "py-12"}}, du.Region(mf.RegionProps{ID: "app-body"}, du.Stack(mf.StackProps{Direction: "column", Gap: "4"}, children...)))
}

func summaryCards(orders []order) mf.Node {
	return du.Grid(mf.GridProps{Columns: "md:grid-cols-3", Gap: "4"}, statCard("Deals", fmt.Sprintf("%d", len(orders)), "Visible"), statCard("Value", "$"+formatNumber(totalAmount(orders)), "Projected ARR"), statCard("High risk", fmt.Sprintf("%d", highRiskCount(orders)), "Needs review"))
}

func totalAmount(orders []order) int {
	t := 0
	for _, o := range orders {
		t += o.Amount
	}
	return t
}
func highRiskCount(orders []order) int {
	c := 0
	for _, o := range orders {
		if o.Risk == "High" {
			c++
		}
	}
	return c
}

func pageContent(currentPage string, orders []order) mf.Node {
	switch currentPage {
	case "pipeline":
		return du.Stack(mf.StackProps{Direction: "column", Gap: "4"}, pipelineChart(orders), ordersTable(orders))
	case "playbooks":
		return du.Stack(mf.StackProps{Direction: "column", Gap: "4"}, du.Card("", "", nil, []mf.Node{du.Stack(mf.StackProps{Direction: "column", Gap: "2"}, mf.Text("Escalate blocked enterprise deals"), mf.Text("Follow up on review deals"), mf.Text("Validate procurement contacts"))}, mf.ComponentProps{}), ordersTable(orders))
	default:
		return du.Stack(mf.StackProps{Direction: "column", Gap: "4"}, summaryCards(orders), pipelineHealth(orders), pipelineChart(orders), ordersTable(orders))
	}
}

func sidebarItems(currentPage string) []mf.Node {
	return []mf.Node{
		navLink("Overview", "/", currentPage == "overview"),
		navLink("Pipeline", "/pipeline", currentPage == "pipeline"),
		navLink("Playbooks", "/playbooks", currentPage == "playbooks"),
	}
}
func navLink(label, href string, active bool) mf.Node {
	return mf.Link(mf.LinkProps{Label: label, Href: href, Props: mf.ComponentProps{Class: navClass(active)}})
}

func navClass(active bool) string {
	if active {
		return "active"
	}
	return ""
}
func statCard(title, value, desc string) mf.Node {
	return du.Card("", "", nil, []mf.Node{du.Stack(mf.StackProps{Direction: "column", Gap: "1"}, mf.Text(title), mf.H3(mf.Text(value)), mf.Text(desc))}, mf.ComponentProps{})
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
	return du.Section(mf.SectionProps{Title: "Pipeline health"}, du.Progress(value, 100, "Active ratio", mf.ComponentProps{}))
}

func pipelineChart(orders []order) mf.Node { /* unchanged behavior */
	statusOrder := []string{"Active", "Review", "Blocked"}
	counts := map[string]float64{"Active": 0, "Review": 0, "Blocked": 0}
	for _, o := range orders {
		counts[o.Status]++
	}
	values := make([]float64, 0, len(statusOrder))
	for _, status := range statusOrder {
		values = append(values, counts[status])
	}
	return du.Section(mf.SectionProps{Title: "Deals by status"}, du.Chart(mf.ChartProps{Type: mf.ChartTypeBar, Labels: statusOrder, Height: 260, Datasets: []mf.ChartDataset{{Label: "Deals", Data: values}}, Options: mf.ChartOptions{BeginAtZero: true, HideLegend: true}}))
}

func ordersTable(orders []order) mf.Node {
	rows := make([]mf.TableComponentRow, 0, len(orders))
	for _, o := range orders {
		label := "Block"
		if o.Status == "Blocked" {
			label = "Re-open"
		}
		action := mf.ActionForm(mf.ActionFormProps{Action: "/orders/toggle-status", Target: "#main-content", Swap: "outerHTML"}, mf.TextField(mf.TextFieldProps{ID: "id-" + o.ID, Name: "id", Value: o.ID, Type: "hidden"}), mf.Button(label, mf.ComponentProps{Variant: "secondary"}))
		rows = append(rows, mf.TableRowValues(o.ID, o.Customer, o.Plan, "$"+formatNumber(o.Amount), o.Risk, o.Status, action))
	}
	return mf.Table(mf.TableProps{Columns: []mf.TableColumn{{Label: "Deal"}, {Label: "Customer"}, {Label: "Plan"}, {Label: "ARR"}, {Label: "Risk"}, {Label: "Status"}, {Label: "Action"}}, Rows: rows, EmptyTitle: "No deals"})
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
