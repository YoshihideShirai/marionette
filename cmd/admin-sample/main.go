package main

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

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

func main() {
	app := buildApp()
	if err := app.Run("127.0.0.1:8082"); err != nil {
		panic(err)
	}
}

func buildApp() *mb.App {
	app := mb.New()
	app.SetGlobal("orders", seedOrders)
	app.SetGlobal("selectedStatus", "all")
	app.SetGlobal("loggedIn", false)
	app.SetGlobal("authError", "")
	app.SetGlobal("flash", "")
	app.SetGlobal("currentPage", "overview")
	app.SetGlobal("currentOrderID", "")

	app.AddStyle(`
		#marionette-root { width: 100%; max-width: none; padding: 0; }
		#marionette-root > * { animation: none; }
	`)

	app.Page("/", func(ctx *mb.Context) mf.Node {
		if !ctx.GetGlobal("loggedIn").(bool) {
			return loginPage(ctx.GetGlobal("authError").(string))
		}
		ctx.SetGlobal("currentPage", "overview")
		ctx.SetGlobal("currentOrderID", "")
		return dashboardFromState(ctx, "overview")
	}, mb.WithTitle("Admin Sample"))

	app.Page("/pipeline", func(ctx *mb.Context) mf.Node {
		if !ctx.GetGlobal("loggedIn").(bool) {
			return loginPage(ctx.GetGlobal("authError").(string))
		}
		ctx.SetGlobal("currentPage", "pipeline")
		ctx.SetGlobal("currentOrderID", "")
		return dashboardFromState(ctx, "pipeline")
	}, mb.WithTitle("Pipeline - Admin Sample"))

	app.Page("/playbooks", func(ctx *mb.Context) mf.Node {
		if !ctx.GetGlobal("loggedIn").(bool) {
			return loginPage(ctx.GetGlobal("authError").(string))
		}
		ctx.SetGlobal("currentPage", "playbooks")
		ctx.SetGlobal("currentOrderID", "")
		return dashboardFromState(ctx, "playbooks")
	}, mb.WithTitle("Playbooks - Admin Sample"))

	app.Page("/orders/detail", func(ctx *mb.Context) mf.Node {
		if !ctx.GetGlobal("loggedIn").(bool) {
			return loginPage(ctx.GetGlobal("authError").(string))
		}
		ctx.SetGlobal("currentPage", "order-detail")
		ctx.SetGlobal("currentOrderID", strings.TrimSpace(ctx.Query("id")))
		return dashboardFromState(ctx, "order-detail")
	}, mb.WithTitle("Deal detail - Admin Sample"))

	app.Action("auth/login", func(ctx *mb.Context) mf.Node {
		provider := strings.TrimSpace(ctx.FormValue("provider"))
		if provider == "demo-sso" {
			ctx.SetGlobal("loggedIn", true)
			ctx.SetGlobal("authError", "")
			ctx.SetGlobal("flash", "Signed in with Demo SSO")
			ctx.SetGlobal("currentPage", "overview")
			return dashboardFromState(ctx, "overview")
		}
		ctx.SetGlobal("loggedIn", false)
		ctx.SetGlobal("authError", "External authentication failed. Please try again.")
		return loginPage(ctx.GetGlobal("authError").(string))
	})

	app.Action("auth/logout", func(ctx *mb.Context) mf.Node {
		ctx.SetGlobal("loggedIn", false)
		ctx.SetGlobal("authError", "")
		ctx.SetGlobal("flash", "")
		return loginPage("")
	})

	app.Action("orders/filter", func(ctx *mb.Context) mf.Node {
		if !ctx.GetGlobal("loggedIn").(bool) {
			return dashboardMainContent(sessionExpiredAlert())
		}
		status := strings.TrimSpace(ctx.FormValue("status"))
		if status == "" {
			status = "all"
		}
		ctx.SetGlobal("selectedStatus", status)
		ctx.SetGlobal("flash", fmt.Sprintf("Filter applied: %s", status))
		return dashboardMainContent(dashboardBody(ctx, ctx.GetGlobal("currentPage").(string)))
	})

	app.Action("orders/toggle-status", func(ctx *mb.Context) mf.Node {
		if !ctx.GetGlobal("loggedIn").(bool) {
			return dashboardMainContent(sessionExpiredAlert())
		}
		id := strings.TrimSpace(ctx.FormValue("id"))
		orders := ctx.GetGlobal("orders").([]order)
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
			ctx.SetGlobal("flash", fmt.Sprintf("%s -> %s", id, orders[i].Status))
			break
		}
		ctx.SetGlobal("orders", orders)
		return dashboardMainContent(dashboardBody(ctx, ctx.GetGlobal("currentPage").(string)))
	})

	return app
}

func dashboardFromState(ctx *mb.Context, currentPage string) mf.Node {
	content := mf.Container(mf.ContainerProps{MaxWidth: "full", Props: mf.ComponentProps{Class: "py-6 px-4 lg:px-6"}},
		dashboardMainContent(dashboardBody(ctx, currentPage)),
	)
	return mf.Region(mf.RegionProps{ID: "app-body"}, drawerLayout(topbar(currentPage), content, drawerMenu(currentPage)))
}

func drawerLayout(navbar, content, menu mf.Node) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "drawer min-h-screen bg-base-200"},
		mf.InputElement(mf.ElementProps{ID: drawerID, Class: "drawer-toggle", Attrs: mf.Attrs{"type": "checkbox"}}),
		mf.DivProps(mf.ElementProps{Class: "drawer-content flex min-h-screen flex-col"}, navbar, content),
		mf.DivProps(mf.ElementProps{Class: "drawer-side z-40"},
			mf.LabelElementProps(mf.ElementProps{Class: "drawer-overlay", Attrs: mf.Attrs{"for": drawerID, "aria-label": "Close navigation drawer"}}),
			menu,
		),
	)
}

func topbar(currentPage string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "navbar sticky top-0 z-30 border-b border-base-300 bg-base-100/95 shadow-sm backdrop-blur"},
		mf.DivProps(mf.ElementProps{Class: "navbar-start gap-2"},
			mf.LabelElementProps(mf.ElementProps{Class: "btn btn-square btn-ghost", Attrs: mf.Attrs{"for": drawerID, "aria-label": "Open navigation drawer"}}, menuIcon()),
			mf.DivProps(mf.ElementProps{Class: "flex flex-col leading-tight"},
				mf.SpanProps(mf.ElementProps{Class: "text-xs font-semibold uppercase tracking-wide text-base-content/50"}, mf.Text("Admin Sample")),
				mf.SpanProps(mf.ElementProps{Class: "font-bold"}, mf.Text(pageTitle(currentPage))),
			),
		),
		mf.DivProps(mf.ElementProps{Class: "navbar-center hidden lg:flex"},
			mf.Tabs(mf.TabsProps{Items: navigationTabs(currentPage), Props: mf.ComponentProps{Class: "tabs-boxed bg-base-200"}}),
		),
		mf.DivProps(mf.ElementProps{Class: "navbar-end gap-2"},
			mf.Badge(mf.BadgeProps{Label: "Live", Props: mf.ComponentProps{Class: "badge-success hidden sm:inline-flex"}}),
			mf.ThemeToggleButton(mf.ComponentProps{}),
		),
	)
}

func drawerMenu(currentPage string) mf.Node {
	return mf.AsideProps(mf.ElementProps{Class: "min-h-full w-80 bg-base-100 p-4 text-base-content"},
		mf.Stack(mf.StackProps{Direction: "column", Gap: "4"},
			mf.DivProps(mf.ElementProps{Class: "rounded-box bg-base-200 p-4"},
				mf.TextComponent(mf.TextProps{Text: "Workspace", Size: "text-xs", Weight: "font-semibold", Props: mf.ComponentProps{Class: "uppercase tracking-wide text-base-content/50"}}),
				mf.H2Props(mf.ElementProps{Class: "text-xl font-bold"}, mf.Text("Revenue Ops")),
				mf.TextComponent(mf.TextProps{Text: "Navigate without a fixed sidebar.", Props: mf.ComponentProps{Class: "text-sm text-base-content/60"}}),
			),
			mf.UlProps(mf.ElementProps{Class: "menu gap-1 rounded-box bg-base-100 p-0"}, drawerItem("Overview", "/", currentPage == "overview"), drawerItem("Pipeline", "/pipeline", currentPage == "pipeline"), drawerItem("Playbooks", "/playbooks", currentPage == "playbooks")),
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
	return mf.Li(mf.AnchorProps(mf.ElementProps{Class: className, Attrs: mf.Attrs{"href": href}}, mf.Text(label)))
}

func navigationTabs(currentPage string) []mf.TabsItem {
	return []mf.TabsItem{
		{Label: "Overview", Href: "/", Active: currentPage == "overview"},
		{Label: "Pipeline", Href: "/pipeline", Active: currentPage == "pipeline"},
		{Label: "Playbooks", Href: "/playbooks", Active: currentPage == "playbooks"},
	}
}

func menuIcon() mf.Node {
	return mf.MenuIcon(mf.ComponentProps{})
}

func dashboardMainContent(content mf.Node) mf.Node {
	return mf.Region(mf.RegionProps{ID: "main-content"}, content)
}

func sessionExpiredAlert() mf.Node {
	return mf.Alert(mf.AlertProps{Title: "Session expired", Description: "Please sign in again.", Props: mf.ComponentProps{Variant: "warning"}})
}

func dashboardBody(ctx *mb.Context, currentPage string) mf.Node {
	selectedStatus := ctx.GetGlobal("selectedStatus").(string)
	allOrders := ctx.GetGlobal("orders").([]order)
	orders := filteredOrders(allOrders, selectedStatus)
	flash := ctx.GetGlobal("flash").(string)
	currentOrderID := stateString(ctx, "currentOrderID")
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
	children = append(children, pageContent(currentPage, selectedStatus, orders, allOrders, currentOrderID))
	return mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, children...)
}

func pageTitle(currentPage string) string {
	switch currentPage {
	case "pipeline":
		return "Pipeline"
	case "playbooks":
		return "Playbooks"
	case "order-detail":
		return "Deal detail"
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
	case "order-detail":
		return "Inspect one deal, review account context, and update its workflow status."
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
		mf.ActionForm(mf.ActionFormProps{Action: "/orders/filter", Target: "#main-content", Swap: "outerHTML", Props: mf.ComponentProps{Class: "mx-auto grid w-full max-w-sm gap-3 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-end"}},
			mf.FormRow(mf.FormRowProps{ID: "status", Label: "Status", Control: mf.Select(mf.SelectFieldProps{ID: "status", Name: "status", Options: statusOptions(selectedStatus)})}),
			mf.IconButton(mf.IconButtonProps{Type: "submit", Label: "Apply", IconSVG: template.HTML(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-[1.2em]"><path d="M5 12h14"/><path d="M12 5l7 7-7 7"/></svg>`), Props: mf.ComponentProps{Variant: "primary", Class: "w-fit"}}),
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
	return mf.Region(mf.RegionProps{ID: "app-body"},
		mf.Container(mf.ContainerProps{MaxWidth: "lg", Centered: true, Props: mf.ComponentProps{Class: "py-12 px-4"}},
			mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, children...),
		),
	)
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

func pageContent(currentPage, selectedStatus string, orders, allOrders []order, currentOrderID string) mf.Node {
	if currentPage == "order-detail" {
		return orderDetailPage(currentOrderID, allOrders)
	}
	sideRail := dashboardSideRail(selectedStatus, orders)
	switch currentPage {
	case "pipeline":
		return wideDashboardGrid(
			mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, pipelineChart(orders), ordersTable(orders)),
			sideRail,
		)
	case "playbooks":
		return wideDashboardGrid(
			mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, playbookSteps(), playbookNotes(), ordersTable(orders)),
			sideRail,
		)
	default:
		return wideDashboardGrid(
			mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, summaryCards(orders), pipelineChart(orders), ordersTable(orders)),
			sideRail,
		)
	}
}

func wideDashboardGrid(main, aside mf.Node) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "grid gap-4 xl:grid-cols-[minmax(0,1fr)_24rem] 2xl:grid-cols-[minmax(0,1fr)_28rem] xl:items-start"}, main, aside)
}

func dashboardSideRail(selectedStatus string, orders []order) mf.Node {
	return mf.AsideProps(mf.ElementProps{Class: "grid gap-4 xl:sticky xl:top-24"},
		filterPanel(selectedStatus),
		pipelineHealth(orders),
		riskChart(orders),
	)
}

func orderDetailPage(id string, orders []order) mf.Node {
	o, ok := findOrder(orders, id)
	if !ok {
		return mf.Stack(mf.StackProps{Direction: "column", Gap: "4"},
			mf.Alert(mf.AlertProps{Title: "Deal not found", Description: "The selected order is no longer available.", Props: mf.ComponentProps{Variant: "warning"}}),
			mf.Link(mf.LinkProps{Label: "Back to dashboard", Href: "/", Props: mf.ComponentProps{Variant: "outline", Size: "sm"}}),
		)
	}

	return mf.Stack(mf.StackProps{Direction: "column", Gap: "4"},
		mf.Link(mf.LinkProps{Label: "Back to dashboard", Href: "/", Props: mf.ComponentProps{Variant: "ghost", Size: "sm"}}),
		mf.Grid(mf.GridProps{Columns: "cards", Gap: "4"},
			mf.Stat("ARR", "$"+formatNumber(o.Amount), o.Plan+" plan"),
			mf.Stat("Risk", o.Risk, "Current account signal"),
			mf.Stat("Status", o.Status, "Workflow state"),
		),
		mf.DivProps(mf.ElementProps{Class: "grid gap-4 xl:grid-cols-[minmax(0,1fr)_22rem] xl:items-start"},
			orderOverviewCard(o),
			orderActionPanel(o),
		),
	)
}

func orderOverviewCard(o order) mf.Node {
	return mf.Section(mf.SectionProps{Title: o.ID + " - " + o.Customer, Description: "Account and commercial details for the selected deal."},
		mf.DescriptionListProps(mf.ElementProps{Class: "grid gap-4 md:grid-cols-2 xl:grid-cols-3"},
			detailField("Customer", o.Customer),
			detailField("Plan", o.Plan),
			detailField("Projected ARR", "$"+formatNumber(o.Amount)),
			detailField("Risk", o.Risk),
			detailField("Status", o.Status),
			detailField("Next step", nextStep(o)),
		),
	)
}

func orderActionPanel(o order) mf.Node {
	label := "Block deal"
	variant := "secondary"
	description := "Move this deal into the blocked workflow for review."
	if o.Status == "Blocked" {
		label = "Re-open deal"
		variant = "success"
		description = "Return this deal to the active workflow."
	}
	return mf.Card(mf.CardProps{
		Title:       "Workflow action",
		Description: description,
		Props:       mf.ComponentProps{Class: "border border-base-300 shadow-none"},
	},
		mf.ActionForm(mf.ActionFormProps{Action: "/orders/toggle-status", Target: "#main-content", Swap: "outerHTML", Props: mf.ComponentProps{Class: "grid gap-3"}},
			mf.TextField(mf.TextFieldProps{ID: "detail-id-" + o.ID, Name: "id", Value: o.ID, Type: "hidden"}),
			mf.SubmitButton(label, mf.ComponentProps{Variant: variant, Class: "w-full"}),
		),
		mf.PProps(mf.ElementProps{Class: "text-sm text-base-content/60"}, mf.Text("Inline updates keep you on this detail view.")),
	)
}

func detailField(label, value string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "rounded-box border border-base-300 bg-base-100 p-4"},
		mf.DescriptionTermProps(mf.ElementProps{Class: "text-xs font-semibold uppercase tracking-wide text-base-content/50"}, mf.Text(label)),
		mf.DescriptionDetailsProps(mf.ElementProps{Class: "mt-1 font-semibold text-base-content"}, mf.Text(value)),
	)
}

func nextStep(o order) string {
	switch o.Status {
	case "Blocked":
		return "Escalate owner and confirm resolution path."
	case "Review":
		return "Schedule follow-up and validate risk signals."
	default:
		return "Keep momentum and confirm close plan."
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
		dealLink := mf.Link(mf.LinkProps{Label: o.ID, Href: orderDetailHref(o.ID), Props: mf.ComponentProps{Class: "font-semibold"}})
		action := mf.ActionForm(mf.ActionFormProps{Action: "/orders/toggle-status", Target: "#main-content", Swap: "outerHTML"}, mf.TextField(mf.TextFieldProps{ID: "id-" + o.ID, Name: "id", Value: o.ID, Type: "hidden"}), mf.SubmitButton(label, mf.ComponentProps{Variant: variant, Size: "sm"}))
		risk := mf.Badge(mf.BadgeProps{Label: o.Risk, Props: mf.ComponentProps{Class: badgeClass(o.Risk)}})
		rows = append(rows, mf.TableRowValues(dealLink, o.Customer, o.Plan, "$"+formatNumber(o.Amount), risk, o.Status, action))
	}
	return mf.Section(mf.SectionProps{Title: "Deals", Description: "Inline actions update the visible fragment."},
		mf.Table(mf.TableProps{Columns: []mf.TableColumn{{Label: "Deal"}, {Label: "Customer"}, {Label: "Plan"}, {Label: "ARR"}, {Label: "Risk"}, {Label: "Status"}, {Label: "Action"}}, Rows: rows, EmptyTitle: "No deals", EmptyDescription: "Try another status filter."}),
	)
}

func findOrder(orders []order, id string) (order, bool) {
	for _, o := range orders {
		if o.ID == id {
			return o, true
		}
	}
	return order{}, false
}

func orderDetailHref(id string) string {
	return "/orders/detail?id=" + url.QueryEscape(id)
}

func stateString(ctx *mb.Context, key string) string {
	v, ok := ctx.GetGlobal(key).(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(v)
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
