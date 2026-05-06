package dashwinddemo

import (
	"fmt"
	"strconv"
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	daisy "github.com/YoshihideShirai/marionette/frontend/daisyui"
)

type statCard struct{ Title, Value, Icon, Description, TrendClass string }
type lead struct{ Name, Role, Email, CreatedAt, Status, Owner, Avatar string }
type transaction struct {
	Invoice, Customer, Plan, Date, Status string
	Amount                                int
}
type routeItem struct{ Path, Icon, Name, Group string }

var statsData = []statCard{
	{"New Users", "34.7k", "U", "↗︎ 2300 (22%)", "text-success"},
	{"Total Sales", "$34,545", "$", "Current month", "text-base-content/60"},
	{"Pending Leads", "450", "L", "50 in hot leads", "text-warning"},
	{"Active Users", "5.6k", "↯", "↙ 300 (18%)", "text-error"},
}
var seedLeads = []lead{
	{"Alex Morgan", "Product buyer", "alex@example.com", "02 May 26", "In Progress", "Olivia", "AM"},
	{"Priya Shah", "Growth lead", "priya@example.com", "28 Apr 26", "Sold", "Noah", "PS"},
	{"Ken Tanaka", "Platform owner", "ken@example.com", "22 Apr 26", "Need Followup", "Emma", "KT"},
	{"Maria Garcia", "Operations", "maria@example.com", "18 Apr 26", "Open", "Liam", "MG"},
	{"Sam Wilson", "Finance", "sam@example.com", "15 Apr 26", "Not Interested", "Ava", "SW"},
}
var seedTransactions = []transaction{
	{"INV-8842", "Acme Inc.", "Enterprise", "May 05, 2026", "Paid", 12500},
	{"INV-8843", "Northwind", "Growth", "May 04, 2026", "Pending", 6200},
	{"INV-8844", "Blue Peak", "Starter", "May 02, 2026", "Failed", 1200},
	{"INV-8845", "Sora Labs", "Enterprise", "Apr 30, 2026", "Paid", 15200},
}
var routes = []routeItem{{"/", "▦", "Dashboard", "main"}, {"/leads", "▣", "Leads", "main"}, {"/transactions", "$", "Transactions", "main"}, {"/analytics", "◒", "Analytics", "main"}, {"/integration", "⚡", "Integration", "main"}, {"/calendar", "◷", "Calendar", "main"}, {"/settings-profile", "⚙", "Profile", "settings"}, {"/settings-team", "◎", "Team Members", "settings"}}

const mainTargetID = "dashwind-main"

func BuildApp() *mb.App {
	app := mb.New()
	app.Set("period", "Last 30 days")
	app.Set("notice", "")
	app.Set("leads", append([]lead(nil), seedLeads...))
	app.AddStyle(dashwindCSS)
	registerPage := func(path, current, title string, body func(*mb.Context) mf.Node) {
		app.Page(path, func(ctx *mb.Context) mf.Node { return shell(current, body(ctx)) }, mb.WithTitle(title))
	}
	registerPage("/", "Dashboard", "DashWind Demo", dashboardPage)
	registerPage("/leads", "Leads", "Leads - DashWind Demo", leadsPage)
	registerPage("/transactions", "Transactions", "Transactions - DashWind Demo", transactionsPage)
	registerPage("/analytics", "Analytics", "Analytics - DashWind Demo", analyticsPage)
	registerPage("/integration", "Integration", "Integration - DashWind Demo", integrationPage)
	registerPage("/calendar", "Calendar", "Calendar - DashWind Demo", calendarPage)
	registerPage("/settings-profile", "Profile", "Profile - DashWind Demo", profilePage)
	registerPage("/settings-team", "Team Members", "Team - DashWind Demo", teamPage)

	app.Action("dashboard/period", func(ctx *mb.Context) mf.Node {
		period := normalizePeriod(ctx.FormValue("period"))
		ctx.Set("period", period)
		ctx.Set("notice", fmt.Sprintf("Period updated to %s", period))
		return mainContent(dashboardPage(ctx))
	})
	app.Action("leads/add", func(ctx *mb.Context) mf.Node {
		leads := append([]lead(nil), ctx.Get("leads").([]lead)...)
		n := len(leads) + 1
		leads = append([]lead{{fmt.Sprintf("Demo Lead %d", n), "New opportunity", fmt.Sprintf("demo%d@example.com", n), "Today", "Open", "Demo", "DL"}}, leads...)
		ctx.Set("leads", leads)
		ctx.Set("notice", "Added a demo lead")
		return mainContent(leadsPage(ctx))
	})
	app.Action("leads/delete", func(ctx *mb.Context) mf.Node {
		email := strings.TrimSpace(ctx.FormValue("email"))
		filtered := make([]lead, 0)
		for _, l := range ctx.Get("leads").([]lead) {
			if l.Email != email {
				filtered = append(filtered, l)
			}
		}
		ctx.Set("leads", filtered)
		ctx.Set("notice", "Lead removed")
		return mainContent(leadsPage(ctx))
	})
	return app
}

func shell(current string, body mf.Node) mf.Node {
	content := div("flex min-h-screen flex-col bg-base-200", topbar(current), mainContent(body))
	side := sidebar(current)
	return mf.Region(mf.RegionProps{ID: "dashwind-app"}, daisy.DrawerWithProps(daisy.DrawerProps{
		ID:           "dashwind-drawer",
		Class:        "lg:drawer-open dashwind-shell",
		ContentClass: "flex min-h-screen flex-col bg-base-200",
		SideClass:    "z-40",
		Content:      content,
		Side:         side,
	}))
}

func mainContent(body mf.Node) mf.Node {
	return mf.Region(mf.RegionProps{ID: mainTargetID, Props: mf.ComponentProps{Class: "flex-1 p-4 md:p-6 lg:p-8 space-y-6"}}, body)
}

func topbar(current string) mf.Node {
	return daisy.NavbarWithProps(daisy.NavbarProps{Class: "sticky top-0 z-30 border-b border-base-300 bg-base-100/90 backdrop-blur"},
		div("flex-none lg:hidden", mf.LabelElementProps(mf.ElementProps{Class: "btn btn-square btn-ghost", Attrs: mf.Attrs{"for": "dashwind-drawer", "aria-label": "open sidebar"}}, mf.Text("☰"))),
		div("flex-1", mf.H1Props(mf.ElementProps{Class: "text-xl font-semibold"}, mf.Text(current))),
		div("hidden max-w-md flex-1 md:block", searchInput()),
		div("flex-none gap-2",
			daisy.ButtonWithAttrs("◐", mf.ComponentProps{Class: "btn-ghost btn-circle"}, map[string]string{"type": "button", "onclick": "mrnToggleTheme()", "aria-label": "toggle theme"}),
			daisy.ButtonContentWithAttrs(mf.ComponentProps{Class: "btn-ghost btn-circle indicator"}, map[string]string{"type": "button", "aria-label": "notifications"}, span("indicator-item badge badge-primary badge-xs", ""), mf.Text("🔔")),
			daisy.AvatarPlaceholder("DW", "", "bg-primary text-primary-content w-10 rounded-full"),
		),
	)
}

func searchInput() mf.Node {
	return mf.LabelElementProps(mf.ElementProps{Class: "input input-bordered flex items-center gap-2"},
		span("opacity-60", "⌕"),
		mf.InputElement(mf.ElementProps{Class: "grow", Attrs: mf.Attrs{"type": "search", "placeholder": "Search DashWind demo"}}),
	)
}

func sidebar(current string) mf.Node {
	return div("min-h-full w-80 bg-base-100 text-base-content shadow-xl",
		div("p-5",
			div("mb-6 flex items-center gap-3",
				div("grid h-11 w-11 place-items-center rounded-2xl bg-primary text-xl font-black text-primary-content", mf.Text("D")),
				div("", mf.H2Props(mf.ElementProps{Class: "text-lg font-bold"}, mf.Text("DashWind")), paragraph("text-xs text-base-content/60", "DaisyUI admin template demo")),
			),
			menuGroup("Menu", "main", current),
			menuGroup("Settings", "settings", current),
			div("mt-6 rounded-box bg-primary/10 p-4 text-sm", paragraph("font-semibold", "Marionette port"), paragraph("mt-1 opacity-70", "React/Redux template patterns rebuilt as Go handlers and htmx fragments.")),
		),
	)
}

func menuGroup(label, group, current string) mf.Node {
	items := []mf.Node{daisy.MenuTitle(label)}
	for _, r := range routes {
		if r.Group != group {
			continue
		}
		items = append(items, daisy.MenuLink(daisy.MenuLinkProps{Label: r.Name, Href: r.Path, Icon: r.Icon, Active: r.Name == current}))
	}
	return daisy.MenuWithProps(daisy.MenuProps{Class: "rounded-box gap-1 p-0"}, items...)
}

func dashboardPage(ctx *mb.Context) mf.Node {
	notice := noticeNode(ctx)
	cards := make([]mf.Node, 0, len(statsData))
	for _, s := range statsData {
		cards = append(cards, statCardNode(s))
	}
	return div("space-y-6",
		pageTitle("Dashboard", "Marionette rebuild of DashWind DashboardTopBar, Stats, Chart, and UserChannels sections.", periodForm(ctx.Get("period").(string))),
		notice,
		div("grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-4", cards...),
		div("grid grid-cols-1 gap-6 xl:grid-cols-2",
			chartCard("Revenue", "Monthly recurring revenue", mf.Chart(mf.ChartProps{Type: mf.ChartTypeLine, Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}, Height: 260, Datasets: []mf.ChartDataset{{Label: "MRR", Data: []float64{18, 24, 28, 32, 38, 45}, BorderColor: "#3b82f6", BackgroundColor: "rgba(59,130,246,.18)", Fill: true, Tension: .35}}})),
			chartCard("Pipeline", "Qualified leads by stage", mf.Chart(mf.ChartProps{Type: mf.ChartTypeBar, Labels: []string{"Open", "Progress", "Sold", "Followup"}, Height: 260, Datasets: []mf.ChartDataset{{Label: "Leads", Data: []float64{92, 128, 54, 76}, BackgroundColor: "#6366f1"}}, Options: mf.ChartOptions{BeginAtZero: true, HideLegend: true}})),
		),
		div("grid grid-cols-1 gap-6 xl:grid-cols-2", amountStats(), userChannels()),
	)
}

func statCardNode(s statCard) mf.Node {
	return daisy.StatsWithProps(daisy.StatsProps{Class: "shadow bg-base-100"}, daisy.StatItem(daisy.StatProps{
		Title:            s.Title,
		Value:            s.Value,
		Description:      s.Description,
		Figure:           mf.Text(s.Icon),
		FigureClass:      "text-primary text-3xl",
		ValueClass:       "text-primary",
		DescriptionClass: "font-medium " + s.TrendClass,
	}))
}

func periodForm(period string) mf.Node {
	period = normalizePeriod(period)
	buttons := []mf.Node{}
	for _, option := range []string{"Last 7 days", "Last 30 days", "This quarter"} {
		className := "btn-sm"
		if option == period {
			className += " btn-primary"
		}
		buttons = append(buttons, daisy.ButtonWithAttrs(option, mf.ComponentProps{Class: className}, map[string]string{"type": "submit", "name": "period", "value": option}))
	}
	return daisy.ActionFormWithOptions(daisy.ActionFormOptions{Action: "/dashboard/period", Target: "#" + mainTargetID, Swap: "outerHTML", Class: "join"}, buttons...)
}

func normalizePeriod(period string) string {
	period = strings.TrimSpace(period)
	for _, allowed := range []string{"Last 7 days", "Last 30 days", "This quarter"} {
		if period == allowed {
			return period
		}
	}
	return "Last 30 days"
}

func amountStats() mf.Node {
	return daisy.StatsWithProps(daisy.StatsProps{Class: "stats-vertical lg:stats-horizontal bg-base-100 shadow"},
		daisy.StatItem(daisy.StatProps{Title: "Total Likes", Value: "25.6K", Description: "21% more than last month"}),
		daisy.StatItem(daisy.StatProps{Title: "Page Views", Value: "2.6M", Description: "14% more than last month"}),
	)
}

func userChannels() mf.Node {
	rows := [][]mf.Node{}
	for _, row := range []string{"Organic search|12,432|46%", "Twitter|8,120|24%", "Newsletter|5,420|18%", "Partners|2,804|12%"} {
		parts := strings.Split(row, "|")
		percent, _ := strconv.Atoi(strings.TrimSuffix(parts[2], "%"))
		rows = append(rows, []mf.Node{
			mf.Text(parts[0]),
			mf.Text(parts[1]),
			daisy.ProgressWithClass(float64(percent), 100, "progress-primary w-32"),
			mf.Text(parts[2]),
		})
	}
	return cardPanel("User Channels", "Traffic source breakdown", daisy.TableWithProps(daisy.TableProps{Rows: rows}))
}

func chartCard(title, desc string, chart mf.Node) mf.Node {
	return mf.Card(mf.CardProps{Title: title, Description: desc, Props: mf.ComponentProps{Class: "bg-base-100 shadow"}}, chart)
}

func leadsPage(ctx *mb.Context) mf.Node {
	rows := [][]mf.Node{}
	for _, l := range ctx.Get("leads").([]lead) {
		rows = append(rows, []mf.Node{leadIdentity(l), mf.Text(l.Email), mf.Text(l.CreatedAt), statusBadge(l.Status), mf.Text(l.Owner), deleteLeadForm(l.Email)})
	}
	return div("space-y-6",
		pageTitle("Current Leads", "DashWind leads table with htmx-powered Add New and delete actions.", daisy.ButtonWithAttrs("Add New", mf.ComponentProps{Class: "btn-primary btn-sm"}, map[string]string{"type": "button", "hx-post": "/leads/add", "hx-target": "#" + mainTargetID, "hx-swap": "outerHTML"})),
		noticeNode(ctx),
		cardPanel("Leads List", "Rendered from Marionette server state instead of a Redux slice/API call.", daisy.TableWithProps(daisy.TableProps{Headers: []string{"Name", "Email Id", "Created At", "Status", "Assigned To", ""}, Rows: rows})),
	)
}

func leadIdentity(l lead) mf.Node {
	return div("flex items-center gap-3",
		daisy.AvatarPlaceholder(l.Avatar, "", "mask mask-squircle w-12 bg-neutral text-neutral-content"),
		div("", div("font-bold", mf.Text(l.Name)), div("text-sm opacity-60", mf.Text(l.Role))),
	)
}

func deleteLeadForm(email string) mf.Node {
	return daisy.ActionFormWithOptions(daisy.ActionFormOptions{Action: "/leads/delete", Target: "#" + mainTargetID, Swap: "outerHTML"},
		daisy.HiddenField("email", email),
		daisy.ButtonWithAttrs("✕", mf.ComponentProps{Class: "btn-square btn-ghost btn-sm"}, map[string]string{"type": "submit"}),
	)
}

func transactionsPage(ctx *mb.Context) mf.Node {
	rows := [][]mf.Node{}
	total := 0
	for _, t := range seedTransactions {
		total += t.Amount
		rows = append(rows, []mf.Node{mf.Text(t.Invoice), mf.Text(t.Customer), mf.Text(t.Plan), mf.Text(t.Date), statusBadge(t.Status), div("font-semibold", mf.Text("$"+strconv.Itoa(t.Amount)))})
	}
	return div("space-y-6",
		pageTitle("Transactions", "DashWind-style billing and transactions list.", daisy.StatsWithProps(daisy.StatsProps{Class: "shadow"}, daisy.StatItem(daisy.StatProps{Title: "Total", Value: "$" + strconv.Itoa(total), ValueClass: "text-primary"}))),
		cardPanel("Recent Transactions", "", daisy.TableWithProps(daisy.TableProps{Headers: []string{"Invoice", "Customer", "Plan", "Date", "Status", "Amount"}, Rows: rows, Class: "table-zebra"})),
	)
}

func analyticsPage(ctx *mb.Context) mf.Node {
	return div("space-y-6", pageTitle("Analytics", "DashWind charts page with Chart.js widgets.", nil), div("grid grid-cols-1 gap-6 xl:grid-cols-2",
		chartCard("Doughnut", "Channel mix", mf.Chart(mf.ChartProps{Type: mf.ChartTypeDoughnut, Labels: []string{"Organic", "Social", "Referral", "Ads"}, Height: 280, Datasets: []mf.ChartDataset{{Label: "Users", Data: []float64{46, 24, 18, 12}}}})),
		chartCard("Scatter", "Lead score vs. ARR", mf.Chart(mf.ChartProps{Type: mf.ChartTypeScatter, Height: 280, Datasets: []mf.ChartDataset{{Label: "Accounts", Points: []mf.ChartPoint{{X: 20, Y: 15}, {X: 42, Y: 38}, {X: 60, Y: 72}, {X: 82, Y: 120}}, BackgroundColor: "#14b8a6"}}})),
	))
}

func integrationPage(ctx *mb.Context) mf.Node {
	return placeholderPage("Integration", "Connected apps", []string{"Stripe billing webhook", "Slack notifications", "HubSpot CRM sync"})
}
func calendarPage(ctx *mb.Context) mf.Node {
	return placeholderPage("Calendar", "Upcoming customer-success events", []string{"May 08 - Enterprise QBR", "May 13 - Renewal review", "May 20 - Product webinar"})
}
func profilePage(ctx *mb.Context) mf.Node {
	return placeholderPage("Profile", "Settings page example", []string{"Name: DashWind Admin", "Role: Revenue Ops", "Theme: DaisyUI corporate/dark"})
}
func teamPage(ctx *mb.Context) mf.Node {
	return placeholderPage("Team Members", "Team settings submenu example", []string{"Olivia - Admin", "Noah - Billing", "Emma - Support"})
}
func placeholderPage(title, desc string, items []string) mf.Node {
	listItems := make([]mf.Node, 0, len(items))
	for _, item := range items {
		listItems = append(listItems, mf.Li(mf.Text(item)))
	}
	return div("space-y-6", pageTitle(title, desc, nil), cardPanel(title, "", mf.UlProps(mf.ElementProps{Class: "list-disc space-y-2 pl-5"}, listItems...)))
}

func pageTitle(title, desc string, actions mf.Node) mf.Node {
	children := []mf.Node{div("", mf.H1Props(mf.ElementProps{Class: "text-3xl font-bold"}, mf.Text(title)), paragraph("mt-1 text-base-content/60", desc))}
	if actions != nil {
		children = append(children, actions)
	}
	return div("flex flex-col gap-4 md:flex-row md:items-center md:justify-between", children...)
}

func noticeNode(ctx *mb.Context) mf.Node {
	notice, _ := ctx.Get("notice").(string)
	if strings.TrimSpace(notice) == "" {
		return mf.Raw("")
	}
	ctx.Set("notice", "")
	return daisy.Alert(notice, "", mf.ComponentProps{Class: "alert-success shadow"})
}

func cardPanel(title, desc string, children ...mf.Node) mf.Node {
	return daisy.CardPanel(daisy.CardPanelProps{Title: title, Description: desc, Class: "bg-base-100 shadow"}, children...)
}

func statusBadge(status string) mf.Node {
	class := "badge-ghost"
	switch status {
	case "In Progress", "Pending":
		class = "badge-primary"
	case "Sold", "Paid":
		class = "badge-success"
	case "Need Followup":
		class = "badge-accent"
	case "Failed", "Not Interested":
		class = "badge-error"
	}
	return daisy.Badge(mf.BadgeProps{Label: status, Props: mf.ComponentProps{Class: class}})
}

func div(className string, children ...mf.Node) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: className}, children...)
}

func paragraph(className string, text string) mf.Node {
	return mf.PProps(mf.ElementProps{Class: className}, mf.Text(text))
}

func span(className string, text string) mf.Node {
	return mf.SpanProps(mf.ElementProps{Class: className}, mf.Text(text))
}

const dashwindCSS = `
#marionette-root { width: 100%; max-width: none; padding: 0; }
#marionette-root > * { animation: none; }
.dashwind-shell .drawer-side .menu a.active { background: var(--color-primary); color: var(--color-primary-content); font-weight: 700; }
.dashwind-shell .card, .dashwind-shell .stats { border: 1px solid color-mix(in oklab, var(--color-base-content) 10%, transparent); }
.dashwind-shell .stat-value { font-size: clamp(1.75rem, 2vw, 2.25rem); }
`
