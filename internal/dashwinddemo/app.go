package dashwinddemo

import (
	"fmt"
	"strconv"
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
	daisy "github.com/YoshihideShirai/marionette/frontend/daisyui"
	dw "github.com/YoshihideShirai/marionette/frontend/dashwind"
)

type statCard struct{ Title, Value, Icon, Description, TrendClass string }
type lead struct{ Name, Role, Email, CreatedAt, Status, Owner, Avatar string }
type transaction struct {
	Invoice, Customer, Plan, Date, Status string
	Amount                                int
}
type integrationItem struct {
	Name, Icon, Description string
	Active                  bool
}
type bill struct{ InvoiceNo, Amount, Description, Status, GeneratedOn, PaidOn string }
type teamMember struct{ Name, Email, Role, Joined, Avatar string }

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
var integrationList = []integrationItem{
	{"Slack", "S", "Instant messaging and workflow notifications for customer operations.", true},
	{"Facebook", "f", "Meta campaign audience sync and lead-form capture for growth teams.", false},
	{"LinkedIn", "in", "Business network enrichment and account-based lead routing.", true},
	{"Google Ads", "G", "Paid-search campaign spend and conversion import for dashboards.", false},
	{"Gmail", "M", "Shared inbox import for support and success handoffs.", false},
	{"Salesforce", "SF", "CRM account, opportunity, and forecast synchronization.", false},
	{"HubSpot", "H", "Inbound marketing, sales, and customer service contact sync.", false},
}
var bills = []bill{
	{"#4567", "23,989", "Product usages", "Pending", "06 Apr 2026", "-"},
	{"#4523", "34,989", "Product usages", "Pending", "07 Mar 2026", "-"},
	{"#4453", "39,989", "Product usages", "Paid", "05 Feb 2026", "12 Apr 2026"},
	{"#4359", "28,927", "Product usages", "Paid", "06 Jan 2026", "13 Mar 2026"},
	{"#3359", "28,927", "Product usages", "Paid", "07 Dec 2025", "14 Feb 2026"},
	{"#3367", "28,927", "Product usages", "Paid", "07 Nov 2025", "15 Jan 2026"},
}
var teamMembers = []teamMember{
	{"Olivia Martin", "olivia@example.com", "Admin", "12 Jan 2026", "OM"},
	{"Noah Williams", "noah@example.com", "Billing", "18 Jan 2026", "NW"},
	{"Emma Brown", "emma@example.com", "Support", "02 Feb 2026", "EB"},
	{"Liam Johnson", "liam@example.com", "Developer", "20 Feb 2026", "LJ"},
}
var routes = dw.Navigation{
	{Label: "Menu", Items: []dw.NavItem{
		{Path: "/", Icon: "▦", Label: "Dashboard"},
		{Path: "/leads", Icon: "▣", Label: "Leads"},
		{Path: "/transactions", Icon: "$", Label: "Transactions"},
		{Path: "/analytics", Icon: "◒", Label: "Analytics"},
		{Path: "/integration", Icon: "⚡", Label: "Integration"},
		{Path: "/calendar", Icon: "◷", Label: "Calendar"},
	}},
	{Label: "Pages", Items: []dw.NavItem{
		{Path: "/login", Icon: "↪", Label: "Login"},
		{Path: "/register", Icon: "U", Label: "Register"},
		{Path: "/forgot-password", Icon: "K", Label: "Forgot Password"},
		{Path: "/blank", Icon: "□", Label: "Blank Page"},
		{Path: "/404", Icon: "!", Label: "404"},
	}},
	{Label: "Settings", Items: []dw.NavItem{
		{Path: "/settings-profile", Icon: "⚙", Label: "Profile"},
		{Path: "/settings-billing", Icon: "W", Label: "Billing"},
		{Path: "/settings-team", Icon: "◎", Label: "Team Members"},
	}},
	{Label: "Documentation", Items: []dw.NavItem{
		{Path: "/getting-started", Icon: "D", Label: "Getting Started"},
		{Path: "/features", Icon: "▤", Label: "Features"},
		{Path: "/components", Icon: "<> ", Label: "Components"},
	}},
}

const mainTargetID = "dashwind-main"

func BuildApp() *mb.App {
	app := mb.New()
	app.Set("period", "Last 30 days")
	app.Set("notice", "")
	app.Set("leads", append([]lead(nil), seedLeads...))
	app.AddStyle(dw.DefaultCSS)
	pageBodies := map[string]func(*mb.Context) mf.Node{
		"/": dashboardPage, "/leads": leadsPage, "/transactions": transactionsPage, "/analytics": analyticsPage, "/integration": integrationPage, "/calendar": calendarPage,
		"/login": loginPage, "/register": registerPreviewPage, "/forgot-password": forgotPasswordPage, "/blank": blankPage, "/404": notFoundPage,
		"/settings-profile": profilePage, "/settings-billing": billingPage, "/settings-team": teamPage,
		"/getting-started": gettingStartedPage, "/features": featuresPage, "/components": componentsPage,
	}
	for _, item := range routeItems(routes) {
		item := item
		body := pageBodies[item.Path]
		if body == nil {
			continue
		}
		title := item.Label + " - DashWind Demo"
		if item.Path == "/" {
			title = "DashWind Demo"
		}
		app.Page(item.Path, func(ctx *mb.Context) mf.Node { return shell(item.Path, body(ctx)) }, mb.WithTitle(title))
	}

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

func routeItems(nav dw.Navigation) []dw.NavItem {
	items := []dw.NavItem{}
	for _, group := range nav {
		items = appendRouteItems(items, group.Items)
	}
	return items
}

func appendRouteItems(dst []dw.NavItem, items []dw.NavItem) []dw.NavItem {
	for _, item := range items {
		dst = append(dst, item)
		dst = appendRouteItems(dst, item.Children)
	}
	return dst
}

func shell(currentPath string, body mf.Node) mf.Node {
	return dw.Shell(dw.ShellProps{
		CurrentPath:       currentPath,
		Brand:             dw.Brand{Title: "DashWind", Subtitle: "DaisyUI admin template demo"},
		SearchPlaceholder: "Search DashWind demo",
		Navigation:        routes,
		SidebarFooter:     div("mt-6 rounded-box bg-primary/10 p-4 text-sm", paragraph("font-semibold", "Marionette port"), paragraph("mt-1 opacity-70", "React/Redux template patterns rebuilt as Go handlers and htmx fragments.")),
	}, body)
}

func mainContent(body mf.Node) mf.Node {
	return dw.ShellContent(mainTargetID, body)
}

func dashboardPage(ctx *mb.Context) mf.Node {
	notice := noticeNode(ctx)
	stats := make([]dw.Stat, 0, len(statsData))
	for _, s := range statsData {
		stats = append(stats, dashwindStat(s))
	}
	return div("space-y-6",
		pageTitle("Dashboard", "Marionette rebuild of DashWind DashboardTopBar, Stats, Chart, and UserChannels sections.", periodForm(ctx.Get("period").(string))),
		notice,
		dw.StatsGrid(dw.StatsGridProps{Items: stats}),
		div("grid grid-cols-1 gap-6 xl:grid-cols-2",
			chartCard("Revenue", "Monthly recurring revenue", mf.Chart(mf.ChartProps{Type: mf.ChartTypeLine, Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}, Height: 260, Datasets: []mf.ChartDataset{{Label: "MRR", Data: []float64{18, 24, 28, 32, 38, 45}, BorderColor: "#3b82f6", BackgroundColor: "rgba(59,130,246,.18)", Fill: true, Tension: .35}}})),
			chartCard("Pipeline", "Qualified leads by stage", mf.Chart(mf.ChartProps{Type: mf.ChartTypeBar, Labels: []string{"Open", "Progress", "Sold", "Followup"}, Height: 260, Datasets: []mf.ChartDataset{{Label: "Leads", Data: []float64{92, 128, 54, 76}, BackgroundColor: "#6366f1"}}, Options: mf.ChartOptions{BeginAtZero: true, HideLegend: true}})),
		),
		div("grid grid-cols-1 gap-6 xl:grid-cols-2", amountStats(), userChannels()),
	)
}

func dashwindStat(s statCard) dw.Stat {
	return dw.Stat{
		Title:            s.Title,
		Value:            s.Value,
		Description:      s.Description,
		Figure:           mf.Text(s.Icon),
		DescriptionClass: "font-medium " + s.TrendClass,
	}
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
	type channelRow struct {
		Source  string
		Users   string
		Percent int
	}
	rows := []channelRow{}
	for _, row := range []string{"Organic search|12,432|46%", "Twitter|8,120|24%", "Newsletter|5,420|18%", "Partners|2,804|12%"} {
		parts := strings.Split(row, "|")
		percent, _ := strconv.Atoi(strings.TrimSuffix(parts[2], "%"))
		rows = append(rows, channelRow{Source: parts[0], Users: parts[1], Percent: percent})
	}
	columns := []dw.Column[channelRow]{
		{Header: "Source", Cell: func(row channelRow) mf.Node { return mf.Text(row.Source) }},
		{Header: "Users", Cell: func(row channelRow) mf.Node { return mf.Text(row.Users) }},
		{Header: "Progress", Cell: func(row channelRow) mf.Node {
			return daisy.ProgressWithClass(float64(row.Percent), 100, "progress-primary w-32")
		}},
		{Header: "Share", Cell: func(row channelRow) mf.Node { return mf.Text(strconv.Itoa(row.Percent) + "%") }},
	}
	return cardPanel("User Channels", "Traffic source breakdown", dw.DataTable(dw.DataTableProps[channelRow]{Columns: columns, Rows: rows, Compact: true}))
}

func chartCard(title, desc string, chart mf.Node) mf.Node {
	return mf.Card(mf.CardProps{Title: title, Description: desc, Props: mf.ComponentProps{Class: "bg-base-100 shadow"}}, chart)
}

func leadsPage(ctx *mb.Context) mf.Node {
	rows := ctx.Get("leads").([]lead)
	columns := []dw.Column[lead]{
		{Header: "Name", Cell: leadIdentity, Sortable: true},
		{Header: "Email Id", Cell: func(l lead) mf.Node { return mf.Text(l.Email) }},
		{Header: "Created At", Cell: func(l lead) mf.Node { return mf.Text(l.CreatedAt) }, Sortable: true},
		{Header: "Status", Cell: func(l lead) mf.Node { return statusBadge(l.Status) }},
		{Header: "Assigned To", Cell: func(l lead) mf.Node { return mf.Text(l.Owner) }},
		{HeaderClass: "w-12", Class: "text-right", Cell: func(l lead) mf.Node { return deleteLeadForm(l.Email) }},
	}
	return div("space-y-6",
		pageTitle("Current Leads", "DashWind leads table with htmx-powered Add New and delete actions.", daisy.ButtonWithAttrs("Add New", mf.ComponentProps{Class: "btn-primary btn-sm"}, map[string]string{"type": "button", "hx-post": "/leads/add", "hx-target": "#" + mainTargetID, "hx-swap": "outerHTML"})),
		noticeNode(ctx),
		cardPanel("Leads List", "Rendered from Marionette server state instead of a Redux slice/API call.", dw.DataTable(dw.DataTableProps[lead]{Columns: columns, Rows: rows, Empty: mf.Text("No leads found")})),
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
	total := 0
	for _, t := range seedTransactions {
		total += t.Amount
	}
	columns := []dw.Column[transaction]{
		{Header: "Invoice", Cell: func(t transaction) mf.Node { return mf.Text(t.Invoice) }, Sortable: true},
		{Header: "Customer", Cell: func(t transaction) mf.Node { return mf.Text(t.Customer) }},
		{Header: "Plan", Cell: func(t transaction) mf.Node { return mf.Text(t.Plan) }},
		{Header: "Date", Cell: func(t transaction) mf.Node { return mf.Text(t.Date) }, Sortable: true},
		{Header: "Status", Cell: func(t transaction) mf.Node { return statusBadge(t.Status) }},
		{Header: "Amount", HeaderClass: "text-right", Class: "text-right", Cell: func(t transaction) mf.Node { return div("font-semibold", mf.Text("$"+strconv.Itoa(t.Amount))) }},
	}
	return div("space-y-6",
		pageTitle("Transactions", "DashWind-style billing and transactions list.", daisy.StatsWithProps(daisy.StatsProps{Class: "shadow"}, daisy.StatItem(daisy.StatProps{Title: "Total", Value: "$" + strconv.Itoa(total), ValueClass: "text-primary"}))),
		cardPanel("Recent Transactions", "", dw.DataTable(dw.DataTableProps[transaction]{Columns: columns, Rows: seedTransactions, Zebra: true})),
	)
}

func analyticsPage(ctx *mb.Context) mf.Node {
	return div("space-y-6", pageTitle("Analytics", "DashWind charts page with Chart.js widgets.", nil), div("grid grid-cols-1 gap-6 xl:grid-cols-2",
		chartCard("Doughnut", "Channel mix", mf.Chart(mf.ChartProps{Type: mf.ChartTypeDoughnut, Labels: []string{"Organic", "Social", "Referral", "Ads"}, Height: 280, Datasets: []mf.ChartDataset{{Label: "Users", Data: []float64{46, 24, 18, 12}}}})),
		chartCard("Scatter", "Lead score vs. ARR", mf.Chart(mf.ChartProps{Type: mf.ChartTypeScatter, Height: 280, Datasets: []mf.ChartDataset{{Label: "Accounts", Points: []mf.ChartPoint{{X: 20, Y: 15}, {X: 42, Y: 38}, {X: 60, Y: 72}, {X: 82, Y: 120}}, BackgroundColor: "#14b8a6"}}})),
	))
}

func integrationPage(ctx *mb.Context) mf.Node {
	cards := make([]mf.Node, 0, len(integrationList))
	for _, item := range integrationList {
		cards = append(cards, integrationCard(item))
	}
	return div("space-y-6",
		pageTitle("Integration", "DashWind-style connected app cards with DaisyUI toggles.", nil),
		div("grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-3", cards...),
	)
}

func integrationCard(item integrationItem) mf.Node {
	return cardPanel(item.Name, "",
		div("flex gap-4", daisy.AvatarPlaceholder(item.Icon, "", "w-12 rounded-box bg-primary/10 text-primary font-bold"), paragraph("text-sm text-base-content/70", item.Description)),
		div("mt-6 text-right", daisy.ToggleWithVariants(strings.ToLower(item.Name)+"-enabled", item.Active, "success", "lg")),
	)
}

func calendarPage(ctx *mb.Context) mf.Node {
	days := []mf.Node{}
	for day := 1; day <= 35; day++ {
		className := "min-h-20 rounded-box border border-base-300 bg-base-100 p-2 text-sm"
		label := strconv.Itoa(day)
		if day > 31 {
			className += " opacity-30"
			label = strconv.Itoa(day - 31)
		}
		if day == 8 || day == 13 || day == 20 {
			className += " ring-2 ring-primary/30"
		}
		days = append(days, div(className, div("font-semibold", mf.Text(label))))
	}
	return div("space-y-6",
		pageTitle("Calendar", "Calendar view and customer-success events matching the DashWind sample area.", nil),
		div("grid grid-cols-1 gap-6 xl:grid-cols-3",
			cardPanel("May 2026", "Monthly schedule", div("grid grid-cols-7 gap-2", days...)),
			cardPanel("Upcoming events", "Right drawer event feed", bulletList("May 08 - Enterprise QBR", "May 13 - Renewal review", "May 20 - Product webinar", "May 28 - Campaign retrospective")),
		),
	)
}

func loginPage(ctx *mb.Context) mf.Node {
	return authPreviewPage("Login", []string{"Email Id", "Password"}, "Forgot Password?", "Don't have an account yet? Register")
}

func registerPreviewPage(ctx *mb.Context) mf.Node {
	return authPreviewPage("Register", []string{"Name", "Email Id", "Password"}, "", "Already have an account? Login")
}

func forgotPasswordPage(ctx *mb.Context) mf.Node {
	return authPreviewPage("Forgot Password", []string{"Email Id"}, "", "Remembered your password? Login")
}

func authPreviewPage(title string, fields []string, helper string, footer string) mf.Node {
	controls := []mf.Node{}
	for _, field := range fields {
		inputType := "text"
		if strings.Contains(strings.ToLower(field), "password") {
			inputType = "password"
		}
		controls = append(controls, inputField(field, inputType, ""))
	}
	if helper != "" {
		controls = append(controls, div("text-right text-primary text-sm", mf.Text(helper)))
	}
	controls = append(controls, daisy.ButtonWithAttrs(title, mf.ComponentProps{Class: "mt-2 w-full btn-primary"}, map[string]string{"type": "button"}), div("text-center mt-4", mf.Text(footer)))
	return div("space-y-6",
		pageTitle(title, "DashWind user page preview rendered inside the Marionette demo shell.", nil),
		div("card mx-auto w-full max-w-5xl shadow-xl bg-base-100",
			div("grid grid-cols-1 rounded-xl md:grid-cols-2",
				div("rounded-l-xl bg-primary p-10 text-primary-content", mf.H2Props(mf.ElementProps{Class: "text-3xl font-bold"}, mf.Text("DashWind")), paragraph("mt-4 opacity-80", "Build admin dashboards with DaisyUI, htmx fragments, and Go state.")),
				div("py-16 px-10", mf.H2Props(mf.ElementProps{Class: "mb-4 text-center text-2xl font-semibold"}, mf.Text(title)), formBlock(controls...)),
			),
		),
	)
}

func blankPage(ctx *mb.Context) mf.Node {
	return div("space-y-6", pageTitle("Blank Page", "A clean DashWind starter surface for new features.", nil), cardPanel("Blank", "Start composing your next Marionette screen here.", div("min-h-64 rounded-box border border-dashed border-base-300")))
}

func notFoundPage(ctx *mb.Context) mf.Node {
	return div("space-y-6", pageTitle("404", "DashWind not-found screen pattern.", nil), div("hero min-h-96 rounded-box bg-base-100 shadow", div("hero-content text-center", div("max-w-md", mf.H1Props(mf.ElementProps{Class: "text-8xl font-bold text-primary"}, mf.Text("404")), paragraph("py-6 text-base-content/70", "The page you are looking for does not exist in this demo."), daisy.ButtonWithAttrs("Back to dashboard", mf.ComponentProps{Class: "btn-primary"}, map[string]string{"type": "button"})))))
}

func profilePage(ctx *mb.Context) mf.Node {
	return div("space-y-6",
		pageTitle("Profile", "Profile settings page with DashWind-style input cards.", daisy.ButtonWithAttrs("Save", mf.ComponentProps{Class: "btn-primary btn-sm"}, map[string]string{"type": "button"})),
		div("grid grid-cols-1 gap-6 xl:grid-cols-3",
			cardPanel("Profile photo", "Avatar and account role", div("flex items-center gap-4", daisy.AvatarPlaceholder("DW", "", "w-20 rounded-full bg-primary text-primary-content text-xl"), div("", mf.H3Props(mf.ElementProps{Class: "font-semibold"}, mf.Text("DashWind Admin")), paragraph("text-sm opacity-70", "Revenue Ops")))),
			cardPanel("Account", "Editable account fields", inputField("Name", "text", "DashWind Admin"), inputField("Email Id", "email", "admin@example.com"), inputField("Role", "text", "Revenue Ops")),
		),
	)
}

func billingPage(ctx *mb.Context) mf.Node {
	columns := []dw.Column[bill]{
		{Header: "Invoice No", Cell: func(b bill) mf.Node { return mf.Text(b.InvoiceNo) }},
		{Header: "Amount", Cell: func(b bill) mf.Node { return mf.Text("$" + b.Amount) }, HeaderClass: "text-right", Class: "text-right"},
		{Header: "Description", Cell: func(b bill) mf.Node { return mf.Text(b.Description) }},
		{Header: "Status", Cell: func(b bill) mf.Node { return statusBadge(b.Status) }},
		{Header: "Generated On", Cell: func(b bill) mf.Node { return mf.Text(b.GeneratedOn) }, Sortable: true},
		{Header: "Paid On", Cell: func(b bill) mf.Node { return mf.Text(b.PaidOn) }},
	}
	return div("space-y-6",
		pageTitle("Billing", "Billing table following the DashWind settings page.", daisy.ButtonWithAttrs("Download All", mf.ComponentProps{Class: "btn-primary btn-sm"}, map[string]string{"type": "button"})),
		cardPanel("Billing History", "Product usage invoices", dw.DataTable(dw.DataTableProps[bill]{Columns: columns, Rows: bills, Zebra: true})),
	)
}

func teamPage(ctx *mb.Context) mf.Node {
	columns := []dw.Column[teamMember]{
		{Header: "Name", Cell: func(member teamMember) mf.Node {
			return div("flex items-center gap-3", daisy.AvatarPlaceholder(member.Avatar, "", "w-10 rounded-full bg-neutral text-neutral-content"), div("", div("font-bold", mf.Text(member.Name)), div("text-sm opacity-60", mf.Text(member.Email))))
		}, Sortable: true},
		{Header: "Role", Cell: func(member teamMember) mf.Node { return mf.Text(member.Role) }},
		{Header: "Joined", Cell: func(member teamMember) mf.Node { return mf.Text(member.Joined) }, Sortable: true},
		{Header: "Status", Cell: func(member teamMember) mf.Node {
			return daisy.Badge(mf.BadgeProps{Label: "Active", Props: mf.ComponentProps{Class: "badge-success"}})
		}},
	}
	return div("space-y-6",
		pageTitle("Team Members", "Team settings list with member roles and status badges.", daisy.ButtonWithAttrs("Add New", mf.ComponentProps{Class: "btn-primary btn-sm"}, map[string]string{"type": "button"})),
		cardPanel("Current Team", "DashWind team table", dw.DataTable(dw.DataTableProps[teamMember]{Columns: columns, Rows: teamMembers})),
	)
}

func gettingStartedPage(ctx *mb.Context) mf.Node {
	return documentationPage("Getting Started", "Quick start steps from the DashWind documentation section.", []string{"Install dependencies", "Run the Marionette DashWind demo", "Customize routes, stats, and tables"})
}

func featuresPage(ctx *mb.Context) mf.Node {
	return documentationPage("Features", "Template feature checklist mirrored in the demo.", []string{"Sidebar and submenu navigation", "Dashboard stats and charts", "Leads, integrations, calendar, billing, and auth pages"})
}

func componentsPage(ctx *mb.Context) mf.Node {
	return documentationPage("Components", "DaisyUI primitives used by the sample.", []string{"dashwind.Shell, NavGroup, and NavItem", "dashwind.StatsGrid, CardPanel, PageHeader", "dashwind.DataTable plus DaisyUI ActionFormWithOptions"})
}

func documentationPage(title, desc string, items []string) mf.Node {
	return div("space-y-6", pageTitle(title, desc, nil), div("grid grid-cols-1 gap-6 lg:grid-cols-3", cardPanel("Navigation", "Documentation submenu", bulletList("Getting Started", "Features", "Components")), cardPanel(title, "Content", bulletList(items...)), cardPanel("Code pointers", "Marionette implementation", bulletList("cmd/dashwind-demo", "internal/dashwinddemo", "frontend/daisyui"))))
}

func bulletList(items ...string) mf.Node {
	children := make([]mf.Node, 0, len(items))
	for _, item := range items {
		children = append(children, mf.Li(mf.Text(item)))
	}
	return mf.UlProps(mf.ElementProps{Class: "list-disc space-y-2 pl-5"}, children...)
}

func formBlock(children ...mf.Node) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "space-y-4"}, children...)
}

func inputField(label string, inputType string, value string) mf.Node {
	if strings.TrimSpace(inputType) == "" {
		inputType = "text"
	}
	name := strings.ToLower(strings.ReplaceAll(label, " ", "-"))
	return mf.LabelElementProps(mf.ElementProps{Class: "form-control w-full"},
		span("label-text mb-1", label),
		mf.InputElement(mf.ElementProps{Class: "input input-bordered w-full", Attrs: mf.Attrs{"type": inputType, "name": name, "value": value}}),
	)
}

func pageTitle(title, desc string, actions mf.Node) mf.Node {
	return dw.PageHeader(dw.PageHeaderProps{Title: title, Description: desc, Actions: actions})
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
	return dw.CardPanel(dw.CardPanelProps{Title: title, Description: desc}, children...)
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
