package adminsample

import (
	"embed"
	"fmt"
	"html/template"
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
		provider := strings.TrimSpace(ctx.FormValue("provider"))
		if provider == "demo-sso" {
			ctx.Set("loggedIn", true)
			ctx.Set("authError", "")
			ctx.Set("flash", "Signed in with Demo SSO")
			return dashboardFromState(ctx, "overview")
		}
		ctx.Set("loggedIn", false)
		ctx.Set("authError", "External authentication failed. Please try again.")
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
	return mf.Region(mf.RegionProps{ID: "app-body"},
		du.DrawerLayout("my-drawer-2", topbar(currentPage), mf.Container(mf.ContainerProps{MaxWidth: "7xl", Centered: true, Props: mf.ComponentProps{Class: "py-6"}}, dashboardMainContent(dashboardBody(ctx, currentPage))), sidebarItems(currentPage)),
	)
}

func topbar(currentPage string) mf.Node {
	return du.DrawerNavbar("my-drawer-2", "admin-sample", []mf.Node{
		navLink("Overview", "/", currentPage == "overview"),
		navLink("Pipeline", "/pipeline", currentPage == "pipeline"),
	})
}

func dashboardMainContent(content mf.Node) mf.Node {
	return mf.Region(mf.RegionProps{ID: "main-content"}, content)
}

func sessionExpiredAlert() mf.Node {
	return mf.Alert(mf.AlertProps{Title: "Session expired", Description: "Please sign in again.", Props: mf.ComponentProps{Variant: "warning"}})
}

func dashboardBody(ctx *mb.Context, currentPage string) mf.Node {
	selectedStatus := ctx.Get("selectedStatus").(string)
	orders := filteredOrders(ctx.Get("orders").([]order), selectedStatus)
	flash := ctx.Get("flash").(string)
	children := []mf.Node{
		mf.PageHeader(mf.PageHeaderProps{
			Title:       "Dashboard",
			Description: "Default daisyUI style admin layout",
			Actions: mf.ActionForm(mf.ActionFormProps{Action: "/auth/logout", Target: "#app-body", Swap: "outerHTML"}, mf.IconButton(mf.IconButtonProps{Type: "submit", Label: "Sign out", IconSVG: template.HTML(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-[1.2em]"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>`), Props: mf.ComponentProps{Variant: "outline"}})),
		}),
		mf.Breadcrumb(mf.BreadcrumbProps{Items: []mf.BreadcrumbItem{{Label: "Home", Href: "/"}, {Label: "Dashboard", Active: true}}}),
		mf.Tabs(mf.TabsProps{Items: []mf.TabsItem{
			{Label: "Overview", Href: "/", Active: currentPage == "overview"},
			{Label: "Pipeline", Href: "/pipeline", Active: currentPage == "pipeline"},
			{Label: "Playbooks", Href: "/playbooks", Active: currentPage == "playbooks"},
		}}),
		mf.Box(mf.BoxProps{Padding: "4", Border: true, Props: mf.ComponentProps{Class: "bg-base-100"}},
			mf.Text("Key metrics for your business are shown below."),
			mf.Actions(mf.ActionsProps{Align: "end", Gap: "2"},
				mf.Badge(mf.BadgeProps{Label: "Live", Props: mf.ComponentProps{Class: "badge-success"}}),
				mf.Badge(mf.BadgeProps{Label: "Updated", Props: mf.ComponentProps{Class: "badge-outline"}}),
			),
		),
	}
	if flash != "" {
		children = append(children, mf.Alert(mf.AlertProps{Title: "Updated", Description: flash, Props: mf.ComponentProps{Variant: "info"}}))
	}
	children = append(children,
		mf.ActionForm(mf.ActionFormProps{Action: "/orders/filter", Target: "#main-content", Swap: "outerHTML", Props: mf.ComponentProps{Class: "card bg-base-100 border border-base-300 p-4"}},
			mf.FormRow(mf.FormRowProps{ID: "status", Label: "Status", Control: mf.Select(mf.SelectFieldProps{ID: "status", Name: "status", Options: statusOptions(selectedStatus)})}),
			mf.Switch(mf.SwitchComponentProps{Name: "high_risk_only", Value: "1", Label: "High risk only", Checked: false, Props: mf.ComponentProps{Class: "mt-2"}}),
			mf.Actions(mf.ActionsProps{Align: "end", Gap: "2"}, mf.IconButton(mf.IconButtonProps{Type: "submit", Label: "Apply", IconSVG: template.HTML(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-[1.2em]"><path d="M5 12h14"/><path d="M12 5l7 7-7 7"/></svg>`), Props: mf.ComponentProps{Variant: "primary"}})),
		),
		pageContent(currentPage, orders),
	)
	return mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, children...)
}

func loginPage(authError string) mf.Node {
	children := []mf.Node{
		mf.PageHeader(mf.PageHeaderProps{Title: "Admin Login", Description: "Sign in via an external identity provider."}),
	}
	if authError != "" {
		children = append(children, mf.Alert(mf.AlertProps{Title: "Login failed", Description: authError, Props: mf.ComponentProps{Variant: "error"}}))
	}
	children = append(children,
		mf.Card(mf.CardProps{Title: "", Description: "", Actions: nil, Gap: "", Props: mf.ComponentProps{}}, mf.Stack(mf.StackProps{Direction: "column", Gap: "4"},
			mf.Image(mf.ImageProps{Src: "https://placehold.co/720x280/0d9488/ffffff?text=Secure+SSO", Alt: "Secure SSO illustration", Props: mf.ComponentProps{Class: "rounded-xl shadow-lg"}}),
			mf.Text("Sign in with an external identity provider to continue."),
			mf.ActionForm(mf.ActionFormProps{Action: "/auth/login", Target: "#app-body", Swap: "outerHTML", Props: mf.ComponentProps{Class: "space-y-3"}},
				mf.HiddenField("provider", "demo-sso"),
				mf.IconButton(mf.IconButtonProps{Type: "submit", Label: "Continue with Demo SSO", IconSVG: template.HTML(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-[1.2em]"><path d="M12 2l7 7-7 7-7-7 7-7z"/><path d="M12 22V10"/></svg>`), Props: mf.ComponentProps{Variant: "primary", Class: "w-full"}}),
			),
		)),
	)
	return mf.Container(mf.ContainerProps{MaxWidth: "lg", Centered: true, Props: mf.ComponentProps{Class: "py-12"}}, mf.Region(mf.RegionProps{ID: "app-body"}, mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, children...)))
}

func summaryCards(orders []order) mf.Node {
	return mf.Grid(mf.GridProps{Columns: "md:grid-cols-3", Gap: "4"}, statCard("Deals", fmt.Sprintf("%d", len(orders)), "Visible"), statCard("Value", "$"+formatNumber(totalAmount(orders)), "Projected ARR"), statCard("High risk", fmt.Sprintf("%d", highRiskCount(orders)), "Needs review"))
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
		return mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, pipelineChart(orders), riskChart(orders), ordersTable(orders))
	case "playbooks":
		return mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, mf.Card(mf.CardProps{Title: "", Description: "", Actions: nil, Gap: "", Props: mf.ComponentProps{}}, mf.Stack(mf.StackProps{Direction: "column", Gap: "2"}, mf.Text("Escalate blocked enterprise deals"), mf.Text("Follow up on review deals"), mf.Text("Validate procurement contacts"))), ordersTable(orders))
	default:
		return mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, summaryCards(orders), pipelineHealth(orders), pipelineChart(orders), riskChart(orders), ordersTable(orders))
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
	iconText := "?"
	switch title {
	case "Deals":
		iconText = "D"
	case "Value":
		iconText = "$"
	case "High risk":
		iconText = "!"
	}
	return mf.Card(mf.CardProps{Title: "", Description: "", Actions: nil, Gap: "", Props: mf.ComponentProps{Class: "text-center"}}, mf.Image(mf.ImageProps{Src: "https://placehold.co/80x80/ffffff/0d9488?text=" + iconText, Alt: title + " icon", Width: 80, Height: 80, Props: mf.ComponentProps{Class: "mx-auto rounded-full bg-base-200"}}), mf.Stack(mf.StackProps{Direction: "column", Gap: "1", Align: "center"}, mf.Text(title), mf.H3(mf.Text(value)), mf.Text(desc)))
}

func badgeClass(risk string) string {
	switch risk {
	case "High":
		return "badge-error"
	case "Medium":
		return "badge-warning"
	default:
		return "badge-success"
	}
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
	return mf.Section(mf.SectionProps{Title: "Pipeline health"}, mf.Progress(mf.ProgressProps{Value: value, Max: 100, Label: "Active ratio", Props: mf.ComponentProps{}}))
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
	return mf.Section(mf.SectionProps{Title: "Deals by status"}, mf.Chart(mf.ChartProps{Type: mf.ChartTypeBar, Labels: statusOrder, Height: 260, Datasets: []mf.ChartDataset{{Label: "Deals", Data: values}}, Options: mf.ChartOptions{BeginAtZero: true, HideLegend: true}}))
}

func riskChart(orders []order) mf.Node {
	labels := []string{"Low", "Medium", "High"}
	counts := map[string]float64{"Low": 0, "Medium": 0, "High": 0}
	for _, o := range orders {
		counts[o.Risk]++
	}
	values := []float64{counts["Low"], counts["Medium"], counts["High"]}
	return mf.Section(mf.SectionProps{Title: "Risk distribution"}, mf.Chart(mf.ChartProps{Type: mf.ChartTypeDoughnut, Labels: labels, Height: 260, Datasets: []mf.ChartDataset{{Label: "Deals", Data: values, BackgroundColor: "rgba(59,130,246,0.35)", BorderColor: "#2563eb"}}, Options: mf.ChartOptions{HideLegend: false}}))
}

func ordersTable(orders []order) mf.Node {
	rows := make([]mf.TableComponentRow, 0, len(orders))
	for _, o := range orders {
		label := "Block"
		if o.Status == "Blocked" {
			label = "Re-open"
		}
		action := mf.ActionForm(mf.ActionFormProps{Action: "/orders/toggle-status", Target: "#main-content", Swap: "outerHTML"}, mf.TextField(mf.TextFieldProps{ID: "id-" + o.ID, Name: "id", Value: o.ID, Type: "hidden"}), mf.SubmitButton(label, mf.ComponentProps{Variant: "secondary"}))
		risk := mf.Badge(mf.BadgeProps{Label: o.Risk, Props: mf.ComponentProps{Class: badgeClass(o.Risk)}})
		rows = append(rows, mf.TableRowValues(o.ID, o.Customer, o.Plan, "$"+formatNumber(o.Amount), risk, o.Status, action))
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
