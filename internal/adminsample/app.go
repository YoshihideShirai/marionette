package adminsample

import (
	"embed"
	"fmt"
	"html/template"
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

const drawerID = "admin-nav-drawer"

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
	app.AddStyle(`
		#marionette-root { width: 100%; max-width: none; padding: 0; }
		#marionette-root > * { animation: none; }
	`)

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
			ctx.Set("currentPage", "overview")
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
	content := mf.Container(mf.ContainerProps{MaxWidth: "7xl", Centered: true, Props: mf.ComponentProps{Class: "py-6 px-4 lg:px-6"}},
		dashboardMainContent(dashboardBody(ctx, currentPage)),
	)
	return mf.Region(mf.RegionProps{ID: "app-body"}, drawerLayout(topbar(currentPage), content, drawerMenu(currentPage)))
}

func drawerLayout(navbar, content, menu mf.Node) mf.Node {
	return mf.Element("div", mf.ElementProps{Class: "drawer min-h-screen bg-base-200"},
		mf.Element("input", mf.ElementProps{ID: drawerID, Class: "drawer-toggle", Attrs: mf.Attrs{"type": "checkbox"}}),
		mf.Element("div", mf.ElementProps{Class: "drawer-content flex min-h-screen flex-col"}, navbar, content),
		mf.Element("div", mf.ElementProps{Class: "drawer-side z-40"},
			mf.Element("label", mf.ElementProps{Class: "drawer-overlay", Attrs: mf.Attrs{"for": drawerID, "aria-label": "Close navigation drawer"}}),
			menu,
		),
	)
}

func topbar(currentPage string) mf.Node {
	return mf.Element("div", mf.ElementProps{Class: "navbar sticky top-0 z-30 border-b border-base-300 bg-base-100/95 shadow-sm backdrop-blur"},
		mf.Element("div", mf.ElementProps{Class: "navbar-start gap-2"},
			mf.Element("label", mf.ElementProps{Class: "btn btn-square btn-ghost", Attrs: mf.Attrs{"for": drawerID, "aria-label": "Open navigation drawer"}}, menuIcon()),
			mf.Element("div", mf.ElementProps{Class: "flex flex-col leading-tight"},
				mf.Element("span", mf.ElementProps{Class: "text-xs font-semibold uppercase tracking-wide text-base-content/50"}, mf.Text("Admin Sample")),
				mf.Element("span", mf.ElementProps{Class: "font-bold"}, mf.Text(pageTitle(currentPage))),
			),
		),
		mf.Element("div", mf.ElementProps{Class: "navbar-center hidden lg:flex"},
			mf.Tabs(mf.TabsProps{Items: navigationTabs(currentPage), Props: mf.ComponentProps{Class: "tabs-boxed bg-base-200"}}),
		),
		mf.Element("div", mf.ElementProps{Class: "navbar-end gap-2"},
			mf.Badge(mf.BadgeProps{Label: "Live", Props: mf.ComponentProps{Class: "badge-success hidden sm:inline-flex"}}),
			mf.ThemeToggleButton(mf.ComponentProps{}),
		),
	)
}

func drawerMenu(currentPage string) mf.Node {
	return mf.Element("aside", mf.ElementProps{Class: "min-h-full w-80 bg-base-100 p-4 text-base-content"},
		mf.Stack(mf.StackProps{Direction: "column", Gap: "4"},
			mf.Element("div", mf.ElementProps{Class: "rounded-box bg-base-200 p-4"},
				mf.TextComponent(mf.TextProps{Text: "Workspace", Size: "text-xs", Weight: "font-semibold", Props: mf.ComponentProps{Class: "uppercase tracking-wide text-base-content/50"}}),
				mf.H2Props(mf.ElementProps{Class: "text-xl font-bold"}, mf.Text("Revenue Ops")),
				mf.TextComponent(mf.TextProps{Text: "Navigate without a fixed sidebar.", Props: mf.ComponentProps{Class: "text-sm text-base-content/60"}}),
			),
			mf.Element("ul", mf.ElementProps{Class: "menu gap-1 rounded-box bg-base-100 p-0"}, drawerItem("Overview", "/", currentPage == "overview"), drawerItem("Pipeline", "/pipeline", currentPage == "pipeline"), drawerItem("Playbooks", "/playbooks", currentPage == "playbooks")),
			mf.Card(mf.CardProps{Title: "Today", Description: "Demo SSO is active for this session.", Props: mf.ComponentProps{Class: "border border-base-300 shadow-none"}},
				mf.Progress(mf.ProgressProps{Value: 72, Max: 100, Label: "Readiness", ShowValue: true, Props: mf.ComponentProps{Variant: "success"}}),
			),
		),
	)
}

func drawerItem(label, href string, active bool) mf.Node {
	className := ""
	if active {
		className = "active"
	}
	return mf.Element("li", mf.ElementProps{}, mf.Element("a", mf.ElementProps{Class: className, Attrs: mf.Attrs{"href": href}}, mf.Text(label)))
}

func navigationTabs(currentPage string) []mf.TabsItem {
	return []mf.TabsItem{
		{Label: "Overview", Href: "/", Active: currentPage == "overview"},
		{Label: "Pipeline", Href: "/pipeline", Active: currentPage == "pipeline"},
		{Label: "Playbooks", Href: "/playbooks", Active: currentPage == "playbooks"},
	}
}

func menuIcon() mf.Node {
	return mf.Raw(`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block h-6 w-6 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/></svg>`)
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
			Title:       pageTitle(currentPage),
			Description: pageDescription(currentPage),
			Actions:     mf.ActionForm(mf.ActionFormProps{Action: "/auth/logout", Target: "#app-body", Swap: "outerHTML"}, mf.IconButton(mf.IconButtonProps{Type: "submit", Label: "Sign out", IconSVG: template.HTML(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-[1.2em]"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>`), Props: mf.ComponentProps{Variant: "outline", Size: "sm"}})),
		}),
		mf.Breadcrumb(mf.BreadcrumbProps{Items: []mf.BreadcrumbItem{{Label: "Home", Href: "/"}, {Label: pageTitle(currentPage), Active: true}}}),
	}
	if flash != "" {
		children = append(children, mf.Toast(mf.ToastProps{Title: "Updated", Description: flash, Props: mf.ComponentProps{Variant: "info"}}))
	}
	children = append(children, filterPanel(selectedStatus), pageContent(currentPage, orders))
	return mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, children...)
}

func pageTitle(currentPage string) string {
	switch currentPage {
	case "pipeline":
		return "Pipeline"
	case "playbooks":
		return "Playbooks"
	default:
		return "Dashboard"
	}
}

func pageDescription(currentPage string) string {
	switch currentPage {
	case "pipeline":
		return "Review deal movement, risk, and the accounts that need attention."
	case "playbooks":
		return "Turn account signals into repeatable next actions for the sales team."
	default:
		return "A compact admin workspace using daisyUI components and an overlay drawer."
	}
}

func filterPanel(selectedStatus string) mf.Node {
	return mf.Card(mf.CardProps{
		Title:       "Filters",
		Description: "Narrow the deal list without leaving the current view.",
		Props:       mf.ComponentProps{Class: "border border-base-300 shadow-none"},
	},
		mf.ActionForm(mf.ActionFormProps{Action: "/orders/filter", Target: "#main-content", Swap: "outerHTML", Props: mf.ComponentProps{Class: "grid gap-3 md:grid-cols-[1fr_auto] md:items-end"}},
			mf.FormRow(mf.FormRowProps{ID: "status", Label: "Status", Control: mf.Select(mf.SelectFieldProps{ID: "status", Name: "status", Options: statusOptions(selectedStatus)})}),
			mf.IconButton(mf.IconButtonProps{Type: "submit", Label: "Apply", IconSVG: template.HTML(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-[1.2em]"><path d="M5 12h14"/><path d="M12 5l7 7-7 7"/></svg>`), Props: mf.ComponentProps{Variant: "primary", Class: "w-full md:w-fit"}}),
		),
	)
}

func loginPage(authError string) mf.Node {
	children := []mf.Node{
		mf.PageHeader(mf.PageHeaderProps{Title: "Admin Login", Description: "Sign in via an external identity provider."}),
	}
	if authError != "" {
		children = append(children, mf.Alert(mf.AlertProps{Title: "Login failed", Description: authError, Props: mf.ComponentProps{Variant: "error"}}))
	}
	children = append(children,
		mf.Card(mf.CardProps{Title: "Secure workspace", Description: "Use Demo SSO to enter the admin console.", Props: mf.ComponentProps{Class: "border border-base-300"}}, mf.Stack(mf.StackProps{Direction: "column", Gap: "4"},
			mf.Hero("Revenue Ops Console", "A focused sample dashboard with drawer navigation, filters, charts, and action forms.", mf.Badge(mf.BadgeProps{Label: "Demo", Props: mf.ComponentProps{Class: "badge-primary"}})),
			mf.ActionForm(mf.ActionFormProps{Action: "/auth/login", Target: "#app-body", Swap: "outerHTML", Props: mf.ComponentProps{Class: "space-y-3"}},
				mf.HiddenField("provider", "demo-sso"),
				mf.IconButton(mf.IconButtonProps{Type: "submit", Label: "Continue with Demo SSO", IconSVG: template.HTML(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-[1.2em]"><path d="M12 2l7 7-7 7-7-7 7-7z"/><path d="M12 22V10"/></svg>`), Props: mf.ComponentProps{Variant: "primary", Class: "w-full"}}),
			),
		)),
	)
	return mf.Container(mf.ContainerProps{MaxWidth: "lg", Centered: true, Props: mf.ComponentProps{Class: "py-12 px-4"}}, mf.Region(mf.RegionProps{ID: "app-body"}, mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, children...)))
}

func summaryCards(orders []order) mf.Node {
	return mf.Grid(mf.GridProps{Columns: "md:grid-cols-3", Gap: "4"},
		mf.Stat("Deals", fmt.Sprintf("%d", len(orders)), "Visible in current filter"),
		mf.Stat("Value", "$"+formatNumber(totalAmount(orders)), "Projected ARR"),
		mf.Stat("High risk", fmt.Sprintf("%d", highRiskCount(orders)), "Needs review"),
	)
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
		return mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, pipelineHealth(orders), pipelineChart(orders), riskChart(orders), ordersTable(orders))
	case "playbooks":
		return mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, playbookSteps(), playbookNotes(), ordersTable(orders))
	default:
		return mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, summaryCards(orders), pipelineHealth(orders), mf.Split(mf.SplitProps{Main: pipelineChart(orders), Aside: riskChart(orders), AsideWidth: "md", Gap: "4"}), ordersTable(orders))
	}
}

func playbookSteps() mf.Node {
	return mf.Section(mf.SectionProps{Title: "Deal response", Description: "A lightweight operating rhythm for the current book."},
		mf.Steps(mf.Step("Qualify", true), mf.Step("Review", true), mf.Step("Escalate", false), mf.Step("Close", false)),
	)
}

func playbookNotes() mf.Node {
	return mf.Grid(mf.GridProps{Columns: "md:grid-cols-3", Gap: "4"},
		mf.Collapse("Blocked enterprise deals", mf.Text("Confirm executive sponsor, legal owner, and security review date."), true),
		mf.Collapse("Review deals", mf.Text("Schedule follow-up and capture the next measurable commitment."), false),
		mf.Collapse("Procurement", mf.Text("Validate contact details and renewal timing before forecast submission."), false),
	)
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
	return mf.Section(mf.SectionProps{Title: "Pipeline health", Description: "Active deals compared with the currently visible pipeline."}, mf.Progress(mf.ProgressProps{Value: value, Max: 100, Label: "Active ratio", ShowValue: true, Props: mf.ComponentProps{Variant: "success"}}))
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
	return mf.Section(mf.SectionProps{Title: "Deals by status", Description: "Volume by current workflow stage."}, mf.Chart(mf.ChartProps{Type: mf.ChartTypeBar, Labels: statusOrder, Height: 260, Datasets: []mf.ChartDataset{{Label: "Deals", Data: values}}, Options: mf.ChartOptions{BeginAtZero: true, HideLegend: true}}))
}

func riskChart(orders []order) mf.Node {
	labels := []string{"Low", "Medium", "High"}
	counts := map[string]float64{"Low": 0, "Medium": 0, "High": 0}
	for _, o := range orders {
		counts[o.Risk]++
	}
	values := []float64{counts["Low"], counts["Medium"], counts["High"]}
	return mf.Section(mf.SectionProps{Title: "Risk distribution", Description: "Open risk signals across visible deals."}, mf.Chart(mf.ChartProps{Type: mf.ChartTypeDoughnut, Labels: labels, Height: 260, Datasets: []mf.ChartDataset{{Label: "Deals", Data: values, BackgroundColor: "rgba(59,130,246,0.35)", BorderColor: "#2563eb"}}, Options: mf.ChartOptions{HideLegend: false}}))
}

func ordersTable(orders []order) mf.Node {
	rows := make([]mf.TableComponentRow, 0, len(orders))
	for _, o := range orders {
		label := "Block"
		variant := "secondary"
		if o.Status == "Blocked" {
			label = "Re-open"
			variant = "success"
		}
		action := mf.ActionForm(mf.ActionFormProps{Action: "/orders/toggle-status", Target: "#main-content", Swap: "outerHTML"}, mf.TextField(mf.TextFieldProps{ID: "id-" + o.ID, Name: "id", Value: o.ID, Type: "hidden"}), mf.SubmitButton(label, mf.ComponentProps{Variant: variant, Size: "sm"}))
		risk := mf.Badge(mf.BadgeProps{Label: o.Risk, Props: mf.ComponentProps{Class: badgeClass(o.Risk)}})
		rows = append(rows, mf.TableRowValues(o.ID, o.Customer, o.Plan, "$"+formatNumber(o.Amount), risk, o.Status, action))
	}
	return mf.Section(mf.SectionProps{Title: "Deals", Description: "Inline actions update the visible fragment."},
		mf.Table(mf.TableProps{Columns: []mf.TableColumn{{Label: "Deal"}, {Label: "Customer"}, {Label: "Plan"}, {Label: "ARR"}, {Label: "Risk"}, {Label: "Status"}, {Label: "Action"}}, Rows: rows, EmptyTitle: "No deals", EmptyDescription: "Try another status filter."}),
	)
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
