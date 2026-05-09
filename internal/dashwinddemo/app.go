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

type statCard struct {
	Title, Value, IconName, Description, Trend string
	TrendTone                                  dw.Tone
}
type lead struct{ Name, Role, Email, CreatedAt, Status, Owner, Avatar string }
type transaction struct {
	Invoice, Customer, Plan, Date, Status string
	Amount                                int
}
type calendarEvent struct {
	Day          int
	Title, Theme string
}
type integrationItem struct {
	Name, Icon, IconURL, Description string
	Active                           bool
}
type bill struct{ InvoiceNo, Amount, Description, Status, GeneratedOn, PaidOn string }
type teamMember struct{ Name, Email, Role, Joined, Avatar string }

var statsData = []statCard{
	{Title: "New Users", Value: "34.7k", IconName: "fa-user-plus", Description: "Acquired this period", Trend: "↗︎ 2300 (22%)", TrendTone: dw.ToneSuccess},
	{Title: "Total Sales", Value: "$34,545", IconName: "fa-dollar-sign", Description: "Current month", Trend: "On track", TrendTone: dw.ToneNeutral},
	{Title: "Pending Leads", Value: "450", IconName: "fa-address-card", Description: "50 in hot leads", Trend: "Needs follow-up", TrendTone: dw.ToneWarning},
	{Title: "Active Users", Value: "5.6k", IconName: "fa-bolt", Description: "Weekly active accounts", Trend: "↙ 300 (18%)", TrendTone: dw.ToneError},
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
var calendarEvents = []calendarEvent{
	{Day: -3, Title: "Product call", Theme: "GREEN"},
	{Day: 1, Title: "Meeting with tech team", Theme: "PINK"},
	{Day: 7, Title: "Meeting with Cristina", Theme: "PURPLE"},
	{Day: 9, Title: "Meeting with Alex", Theme: "BLUE"},
	{Day: 9, Title: "Product Call", Theme: "GREEN"},
	{Day: 9, Title: "Client Meeting", Theme: "PURPLE"},
	{Day: 12, Title: "Client Meeting", Theme: "ORANGE"},
	{Day: 14, Title: "Product meeting", Theme: "PINK"},
	{Day: 17, Title: "Sales Meeting", Theme: "GREEN"},
	{Day: 17, Title: "Product Meeting", Theme: "ORANGE"},
	{Day: 17, Title: "Marketing Meeting", Theme: "PINK"},
	{Day: 17, Title: "Client Meeting", Theme: "GREEN"},
	{Day: 21, Title: "Sales meeting", Theme: "BLUE"},
	{Day: 25, Title: "Client meeting", Theme: "PURPLE"},
}
var integrationList = []integrationItem{
	{Name: "Slack", Icon: "S", IconURL: "https://cdn.simpleicons.org/slack", Description: "Instant messaging and workflow notifications for customer operations.", Active: true},
	{Name: "Facebook", Icon: "f", IconURL: "https://cdn.simpleicons.org/facebook/1877F2", Description: "Meta campaign audience sync and lead-form capture for growth teams.", Active: false},
	{Name: "LinkedIn", Icon: "in", IconURL: "https://cdn.simpleicons.org/linkedin/0A66C2", Description: "Business network enrichment and account-based lead routing.", Active: true},
	{Name: "Google Ads", Icon: "G", IconURL: "https://cdn.simpleicons.org/googleads/4285F4", Description: "Paid-search campaign spend and conversion import for dashboards.", Active: false},
	{Name: "Gmail", Icon: "M", IconURL: "https://cdn.simpleicons.org/gmail/EA4335", Description: "Shared inbox import for support and success handoffs.", Active: false},
	{Name: "Salesforce", Icon: "SF", IconURL: "https://cdn.simpleicons.org/salesforce/00A1E0", Description: "CRM account, opportunity, and forecast synchronization.", Active: false},
	{Name: "HubSpot", Icon: "H", IconURL: "https://cdn.simpleicons.org/hubspot/FF7A59", Description: "Inbound marketing, sales, and customer service contact sync.", Active: false},
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
		{Path: "/", IconNode: navFontIcon("fa-house"), Label: "Dashboard"},
		{Path: "/leads", IconNode: navFontIcon("fa-users"), Label: "Leads"},
		{Path: "/transactions", IconNode: navFontIcon("fa-money-bill-transfer"), Label: "Transactions"},
		{Path: "/analytics", IconNode: navFontIcon("fa-chart-pie"), Label: "Analytics"},
		{Path: "/integration", IconNode: navFontIcon("fa-bolt"), Label: "Integration"},
		{Path: "/calendar", IconNode: navFontIcon("fa-calendar-days"), Label: "Calendar"},
		{IconNode: navFontIcon("fa-file-lines"), Label: "Pages", Children: []dw.NavItem{
			{Path: "/login", IconNode: navFontIcon("fa-right-to-bracket"), Label: "Login"},
			{Path: "/register", IconNode: navFontIcon("fa-user-plus"), Label: "Register"},
			{Path: "/forgot-password", IconNode: navFontIcon("fa-key"), Label: "Forgot Password"},
			{Path: "/blank", IconNode: navFontIcon("fa-square"), Label: "Blank Page"},
			{Path: "/404", IconNode: navFontIcon("fa-triangle-exclamation"), Label: "404"},
		}},
		{IconNode: navFontIcon("fa-gear"), Label: "Settings", Children: []dw.NavItem{
			{Path: "/settings-profile", IconNode: navFontIcon("fa-user-gear"), Label: "Profile"},
			{Path: "/settings-billing", IconNode: navFontIcon("fa-wallet"), Label: "Billing"},
			{Path: "/settings-team", IconNode: navFontIcon("fa-user-group"), Label: "Team Members"},
		}},
		{IconNode: navFontIcon("fa-book-open"), Label: "Documentation", Children: []dw.NavItem{
			{Path: "/getting-started", IconNode: navFontIcon("fa-circle-play"), Label: "Getting Started"},
			{Path: "/features", IconNode: navFontIcon("fa-list-check"), Label: "Features"},
			{Path: "/components", IconNode: navFontIcon("fa-code"), Label: "Components"},
		}},
	}},
}

func navFontIcon(name string) mf.Node {
	return fontIcon(name, "w-6 text-center text-base")
}

func fontIcon(name, className string) mf.Node {
	return mf.FontIcon(mf.FontIconProps{
		Library:    "fa-solid",
		Name:       name,
		Decorative: true,
		Props:      mf.ComponentProps{Class: className},
	})
}

func defaultIntegrationStates() map[string]bool {
	states := map[string]bool{}
	for _, item := range integrationList {
		states[item.Name] = item.Active
	}
	return states
}

func copyIntegrationStates(value any) map[string]bool {
	states, _ := value.(map[string]bool)
	copied := defaultIntegrationStates()
	for name, active := range states {
		copied[name] = active
	}
	return copied
}

const mainTargetID = "dashwind-main"

func BuildApp() *mb.App {
	app := mb.New()
	app.SetGlobal("period", "Last 30 days")
	app.SetGlobal("notice", "")
	app.SetGlobal("integrationStates", defaultIntegrationStates())
	app.SetGlobal("calendarSelectedDay", 8)
	app.SetGlobal("leads", append([]lead(nil), seedLeads...))
	dw.Use(app, dw.Options{})
	app.AddStylesheet("https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.2/css/all.min.css")
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
		ctx.SetGlobal("period", period)
		ctx.SetGlobal("notice", fmt.Sprintf("Period updated to %s", period))
		return mainContent(dashboardPage(ctx))
	})
	app.Action("integration/toggle", func(ctx *mb.Context) mf.Node {
		name := strings.TrimSpace(ctx.FormValue("integration"))
		states := copyIntegrationStates(ctx.GetGlobal("integrationStates"))
		active := !states[name]
		states[name] = active
		ctx.SetGlobal("integrationStates", states)
		status := "disabled"
		if active {
			status = "enabled"
		}
		ctx.SetGlobal("notice", fmt.Sprintf("%s %s", name, status))
		return mainContent(integrationPage(ctx))
	})
	app.Action("calendar/day", func(ctx *mb.Context) mf.Node {
		day, err := strconv.Atoi(strings.TrimSpace(ctx.FormValue("day")))
		if err != nil || day < 1 || day > 31 {
			day = 8
		}
		ctx.SetGlobal("calendarSelectedDay", day)
		ctx.SetGlobal("notice", fmt.Sprintf("May %02d selected", day))
		return mainContent(calendarPage(ctx))
	})
	app.Action("leads/add", func(ctx *mb.Context) mf.Node {
		leads := append([]lead(nil), ctx.GetGlobal("leads").([]lead)...)
		n := len(leads) + 1
		leads = append([]lead{{fmt.Sprintf("Demo Lead %d", n), "New opportunity", fmt.Sprintf("demo%d@example.com", n), "Today", "Open", "Demo", "DL"}}, leads...)
		ctx.SetGlobal("leads", leads)
		ctx.SetGlobal("notice", "Added a demo lead")
		return mainContent(leadsPage(ctx))
	})
	app.Action("leads/delete", func(ctx *mb.Context) mf.Node {
		email := strings.TrimSpace(ctx.FormValue("email"))
		filtered := make([]lead, 0)
		for _, l := range ctx.GetGlobal("leads").([]lead) {
			if l.Email != email {
				filtered = append(filtered, l)
			}
		}
		ctx.SetGlobal("leads", filtered)
		ctx.SetGlobal("notice", "Lead removed")
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
	metrics := make([]dw.Metric, 0, len(statsData))
	for _, s := range statsData {
		metrics = append(metrics, statCardNode(s))
	}
	return div("space-y-6",
		pageTitle("Dashboard", "Marionette rebuild of DashWind DashboardTopBar, Stats, Chart, and UserChannels sections.", periodForm(ctx.GetGlobal("period").(string))),
		notice,
		dw.MetricGrid(dw.MetricGridProps{Items: metrics}),
		div("grid grid-cols-1 gap-6 xl:grid-cols-2",
			chartCard("Revenue", "Monthly recurring revenue", mf.Chart(mf.ChartProps{Type: mf.ChartTypeLine, Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}, Height: 260, Datasets: []mf.ChartDataset{{Label: "MRR", Data: []float64{18, 24, 28, 32, 38, 45}, BorderColor: "#3b82f6", BackgroundColor: "rgba(59,130,246,.18)", Fill: true, Tension: .35}}})),
			chartCard("Pipeline", "Qualified leads by stage", mf.Chart(mf.ChartProps{Type: mf.ChartTypeBar, Labels: []string{"Open", "Progress", "Sold", "Followup"}, Height: 260, Datasets: []mf.ChartDataset{{Label: "Leads", Data: []float64{92, 128, 54, 76}, BackgroundColor: "#6366f1"}}, Options: mf.ChartOptions{BeginAtZero: true, HideLegend: true}})),
		),
		div("grid grid-cols-1 gap-6 xl:grid-cols-2", amountStats(), userChannels()),
	)
}

func statCardNode(s statCard) dw.Metric {
	return dw.Metric{
		Title:       s.Title,
		Value:       s.Value,
		Description: s.Description,
		Trend:       s.Trend,
		TrendTone:   s.TrendTone,
		Icon:        fontIcon(s.IconName, ""),
	}
}

func periodForm(period string) mf.Node {
	period = normalizePeriod(period)
	dateRange := "2026-04-08 ~ 2026-05-07"
	if period == "Last 7 days" {
		dateRange = "2026-05-01 ~ 2026-05-07"
	} else if period == "This quarter" {
		dateRange = "2026-04-01 ~ 2026-06-30"
	}
	return daisy.ActionFormWithOptions(daisy.ActionFormOptions{Action: "/dashboard/period", Target: "#" + mainTargetID, Swap: "outerHTML", Class: "flex flex-col gap-2 sm:flex-row sm:items-center"},
		mf.InputElement(mf.ElementProps{Attrs: mf.Attrs{"class": "input input-bordered input-sm w-full sm:w-72", "name": "dateRange", "readonly": "readonly", "aria-label": "Date range", "value": dateRange}}),
		daisy.ButtonWithAttrs("Refresh Data", mf.ComponentProps{Class: activePeriodClass(period, "Last 30 days")}, map[string]string{"type": "submit", "name": "period", "value": "Last 30 days"}),
		daisy.ButtonWithAttrs("Share", mf.ComponentProps{Class: "btn-sm btn-outline"}, map[string]string{"type": "submit", "name": "period", "value": period}),
		dashboardMoreMenu(),
	)
}

func dashboardMoreMenu() mf.Node {
	return div("dropdown dropdown-end",
		daisy.ButtonContentWithAttrs(mf.ComponentProps{Class: "btn-sm btn-ghost btn-square"}, map[string]string{"type": "button", "tabindex": "0", "aria-label": "Dashboard options"}, fontIcon("fa-ellipsis-vertical", "h-4 w-4")),
		mf.UlProps(mf.ElementProps{Class: "dropdown-content menu bg-base-100 rounded-box z-10 w-44 p-2 shadow", Attrs: mf.Attrs{"tabindex": "0"}},
			mf.Li(daisy.ButtonWithAttrs("Email Digests", mf.ComponentProps{Class: "btn-ghost btn-sm justify-start"}, map[string]string{"type": "submit", "name": "period", "value": "Last 7 days"})),
			mf.Li(daisy.ButtonWithAttrs("Download", mf.ComponentProps{Class: "btn-ghost btn-sm justify-start"}, map[string]string{"type": "submit", "name": "period", "value": "This quarter"})),
		),
	)
}

func activePeriodClass(period, option string) string {
	className := "btn-sm"
	if option == period {
		className += " btn-primary"
	}
	return className
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
	return daisy.StatsWithProps(daisy.StatsProps{Class: "w-full bg-base-100 shadow"},
		daisy.StatItem(daisy.StatProps{Title: "Amount to be Collected", Value: "$25,600", Description: "↗︎ 21% more than last month", ValueClass: "text-primary", Figure: daisy.ButtonWithAttrs("View Users", mf.ComponentProps{Class: "btn-primary btn-sm"}, map[string]string{"type": "button"})}),
		daisy.StatItem(daisy.StatProps{Title: "Cash in hand", Value: "$7,200", Description: "Current operating balance", ValueClass: "text-secondary", Figure: daisy.ButtonWithAttrs("View Members", mf.ComponentProps{Class: "btn-outline btn-sm"}, map[string]string{"type": "button"})}),
	)
}

type signupSource struct {
	Source, Users, Conversion string
}

func userChannels() mf.Node {
	rows := []signupSource{
		{Source: "Facebook Ads", Users: "26,345", Conversion: "10.2%"},
		{Source: "Google Ads", Users: "21,341", Conversion: "11.7%"},
		{Source: "Instagram Ads", Users: "34,379", Conversion: "12.4%"},
		{Source: "Affiliates", Users: "12,359", Conversion: "8.9%"},
		{Source: "Organic", Users: "10,345", Conversion: "7.4%"},
	}
	columns := []dw.Column[signupSource]{
		{Header: "Source", Cell: func(row signupSource) mf.Node { return mf.Text(row.Source) }},
		{Header: "No of Users", HeaderClass: "text-right", Class: "text-right", Cell: func(row signupSource) mf.Node { return mf.Text(row.Users) }},
		{Header: "Conversion", HeaderClass: "text-right", Class: "text-right", Cell: func(row signupSource) mf.Node {
			return daisy.Badge(mf.BadgeProps{Label: row.Conversion, Props: mf.ComponentProps{Class: "badge-success"}})
		}},
	}
	return cardPanel("User Signup Source", "Channel conversion table", dw.DataTable(dw.DataTableProps[signupSource]{Columns: columns, Rows: rows, Compact: true}))
}

func chartCard(title, desc string, chart mf.Node) mf.Node {
	return mf.Card(mf.CardProps{Title: title, Description: desc, Props: mf.ComponentProps{Class: "bg-base-100 shadow"}}, chart)
}

func leadsPage(ctx *mb.Context) mf.Node {
	rows := ctx.GetGlobal("leads").([]lead)
	columns := []dw.Column[lead]{
		{Header: "Name", Cell: leadIdentity, Sortable: true},
		{Header: "Email Id", Cell: func(l lead) mf.Node { return mf.Text(l.Email) }},
		{Header: "Created At", Cell: func(l lead) mf.Node { return mf.Text(l.CreatedAt) }, Sortable: true},
		{Header: "Status", Cell: func(l lead) mf.Node { return statusBadge(l.Status) }},
		{Header: "Assigned To", Cell: func(l lead) mf.Node { return mf.Text(l.Owner) }},
	}
	return div("space-y-6",
		noticeNode(ctx),
		dw.ResourcePage(dw.ResourcePageProps[lead]{
			Title:       "Current Leads",
			Description: "DashWind leads table with htmx-powered Add New and delete actions.",
			Rows:        rows,
			Columns:     columns,
			PrimaryAction: dw.Action{
				Label:  "Add New",
				Action: "/leads/add",
				Target: "#" + mainTargetID,
				Swap:   "outerHTML",
				Class:  "btn-primary btn-sm",
			},
			EmptyState: mf.EmptyStateProps{Title: "No leads found", Description: "Add a demo lead to repopulate this table."},
			RowActions: func(l lead) []dw.Action {
				return []dw.Action{{
					Action:  "/leads/delete",
					Target:  "#" + mainTargetID,
					Swap:    "outerHTML",
					Class:   "btn-square btn-ghost btn-sm",
					Fields:  map[string]string{"email": l.Email},
					Attrs:   mf.Attrs{"aria-label": "Delete " + l.Name},
					Content: fontIcon("fa-trash", "h-4 w-4"),
				}}
			},
		}),
	)
}

func leadIdentity(l lead) mf.Node {
	return div("flex items-center gap-3",
		daisy.AvatarPlaceholder(l.Avatar, "", "mask mask-squircle w-12 bg-neutral text-neutral-content"),
		div("", div("font-bold", mf.Text(l.Name)), div("text-sm opacity-60", mf.Text(l.Role))),
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
	return dw.ResourcePage(dw.ResourcePageProps[transaction]{
		Title:       "Transactions",
		Description: "DashWind-style billing and transactions list.",
		Rows:        seedTransactions,
		Columns:     columns,
		Filters:     []mf.Node{daisy.StatsWithProps(daisy.StatsProps{Class: "shadow"}, daisy.StatItem(daisy.StatProps{Title: "Total", Value: "$" + strconv.Itoa(total), ValueClass: "text-primary"}))},
		EmptyState:  mf.EmptyStateProps{Title: "No transactions found", Description: "Transactions will appear here after invoices are created."},
	})
}

func analyticsPage(ctx *mb.Context) mf.Node {
	return div("space-y-6", pageTitle("Analytics", "DashWind charts page with Chart.js widgets.", nil), div("grid grid-cols-1 gap-6 xl:grid-cols-2",
		chartCard("Doughnut", "Channel mix", mf.Chart(mf.ChartProps{Type: mf.ChartTypeDoughnut, Labels: []string{"Organic", "Social", "Referral", "Ads"}, Height: 280, Datasets: []mf.ChartDataset{{Label: "Users", Data: []float64{46, 24, 18, 12}}}})),
		chartCard("Scatter", "Lead score vs. ARR", mf.Chart(mf.ChartProps{Type: mf.ChartTypeScatter, Height: 280, Datasets: []mf.ChartDataset{{Label: "Accounts", Points: []mf.ChartPoint{{X: 20, Y: 15}, {X: 42, Y: 38}, {X: 60, Y: 72}, {X: 82, Y: 120}}, BackgroundColor: "#14b8a6"}}})),
	))
}

func integrationPage(ctx *mb.Context) mf.Node {
	states := copyIntegrationStates(ctx.GetGlobal("integrationStates"))
	cards := make([]mf.Node, 0, len(integrationList))
	for _, item := range integrationList {
		item.Active = states[item.Name]
		cards = append(cards, integrationCard(item))
	}
	return div("space-y-6",
		pageTitle("Integration", "DashWind-style connected app cards with logo art, htmx toggles, and notification feedback.", nil),
		noticeNode(ctx),
		div("grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-3", cards...),
	)
}

func integrationCard(item integrationItem) mf.Node {
	status := "Disabled"
	buttonClass := "btn-outline"
	if item.Active {
		status = "Enabled"
		buttonClass = "btn-success"
	}
	return cardPanel(item.Name, "",
		div("flex gap-4", integrationLogo(item), paragraph("text-sm text-base-content/70", item.Description)),
		div("mt-6 flex items-center justify-between",
			daisy.Badge(mf.BadgeProps{Label: status, Props: mf.ComponentProps{Class: "badge-ghost"}}),
			daisy.ActionFormWithOptions(daisy.ActionFormOptions{Action: "/integration/toggle", Target: "#" + mainTargetID, Swap: "outerHTML", Class: "inline-flex items-center gap-3"},
				daisy.HiddenField("integration", item.Name),
				daisy.ToggleWithVariants(strings.ToLower(strings.ReplaceAll(item.Name, " ", "-"))+"-enabled", item.Active, "success", "lg"),
				daisy.ButtonWithAttrs("Toggle", mf.ComponentProps{Class: "btn-sm " + buttonClass}, map[string]string{"type": "submit"}),
			),
		),
	)
}

func integrationLogo(item integrationItem) mf.Node {
	if strings.TrimSpace(item.IconURL) == "" {
		return daisy.AvatarPlaceholder(item.Icon, "", "w-12 rounded-box bg-primary/10 text-primary font-bold")
	}
	return div("grid h-12 w-12 shrink-0 place-items-center rounded-box bg-base-200 p-2",
		mf.Element("img", mf.ElementProps{Attrs: mf.Attrs{"src": item.IconURL, "alt": item.Name + " logo", "class": "h-8 w-8 object-contain"}}),
	)
}

func calendarPage(ctx *mb.Context) mf.Node {
	selectedDay, _ := ctx.GetGlobal("calendarSelectedDay").(int)
	if selectedDay < 1 || selectedDay > 31 {
		selectedDay = 9
	}
	return div("space-y-6",
		noticeNode(ctx),
		div("w-full rounded-lg bg-base-100 p-4 shadow",
			calendarToolbar(selectedDay),
			div("divider my-4", mf.Raw("")),
			calendarWeekdayRow(),
			div("mt-1 grid grid-cols-7 place-items-center", calendarDayCells(selectedDay)...),
		),
	)
}

func calendarToolbar(selectedDay int) mf.Node {
	return div("flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between",
		div("flex flex-wrap items-center gap-2 sm:gap-4",
			paragraph("w-48 text-xl font-semibold", "May 2026"),
			span("text-xs", "Beta"),
			daisy.ButtonContentWithAttrs(mf.ComponentProps{Class: "btn-square btn-sm btn-ghost"}, map[string]string{"type": "button", "aria-label": "Previous month"}, fontIcon("fa-chevron-left", "h-5 w-5")),
			daisy.ButtonWithAttrs("Current Month", mf.ComponentProps{Class: "btn-sm btn-ghost normal-case"}, map[string]string{"type": "button"}),
			daisy.ButtonContentWithAttrs(mf.ComponentProps{Class: "btn-square btn-sm btn-ghost"}, map[string]string{"type": "button", "aria-label": "Next month"}, fontIcon("fa-chevron-right", "h-5 w-5")),
		),
		daisy.ActionFormWithOptions(daisy.ActionFormOptions{Action: "/calendar/day", Target: "#" + mainTargetID, Swap: "outerHTML"},
			daisy.HiddenField("day", strconv.Itoa(selectedDay)),
			daisy.ButtonWithAttrs("Add New Event", mf.ComponentProps{Class: "btn-sm btn-ghost btn-outline normal-case"}, map[string]string{"type": "submit"}),
		),
	)
}

func calendarWeekdayRow() mf.Node {
	weekdays := []string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}
	nodes := make([]mf.Node, 0, len(weekdays))
	for _, day := range weekdays {
		nodes = append(nodes, div("text-xs capitalize", mf.Text(day)))
	}
	return div("grid grid-cols-7 gap-6 place-items-center sm:gap-12", nodes...)
}

func calendarDayCells(selectedDay int) []mf.Node {
	days := make([]mf.Node, 0, 42)
	for offset := -4; offset <= 37; offset++ {
		monthDay := offset
		displayDay := offset
		inMonth := offset >= 1 && offset <= 31
		if offset < 1 {
			displayDay = 30 + offset
		}
		if offset > 31 {
			displayDay = offset - 31
		}
		days = append(days, calendarDayCell(monthDay, displayDay, inMonth, selectedDay))
	}
	return days
}

func calendarDayCell(monthDay, displayDay int, inMonth bool, selectedDay int) mf.Node {
	dayClass := "inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-full mx-1 mt-1 text-sm hover:bg-base-300"
	if !inMonth {
		dayClass += " text-slate-400 dark:text-slate-600"
	}
	if monthDay == 9 && inMonth {
		dayClass += " bg-blue-100 dark:bg-blue-400 dark:text-white dark:hover:bg-base-300"
	}
	if monthDay == selectedDay && inMonth {
		dayClass += " ring-2 ring-primary"
	}

	children := []mf.Node{daisy.ActionFormWithOptions(daisy.ActionFormOptions{Action: "/calendar/day", Target: "#" + mainTargetID, Swap: "outerHTML"},
		daisy.HiddenField("day", strconv.Itoa(displayDay)),
		mf.Element("button", mf.ElementProps{Class: dayClass, Attrs: mf.Attrs{"type": "submit", "aria-label": fmt.Sprintf("Select May %d", displayDay)}}, mf.Text(strconv.Itoa(displayDay))),
	)}

	events := eventsForCalendarDay(monthDay)
	visibleEvents := events
	moreCount := 0
	if len(events) > 2 {
		moreCount = len(events) - 2
		visibleEvents = events[:2]
	}
	for _, event := range visibleEvents {
		children = append(children, paragraph("mt-1 truncate px-2 text-xs "+calendarThemeClass(event.Theme), event.Title))
	}
	if moreCount > 0 {
		children = append(children, daisy.ActionFormWithOptions(daisy.ActionFormOptions{Action: "/calendar/day", Target: "#" + mainTargetID, Swap: "outerHTML"},
			daisy.HiddenField("day", strconv.Itoa(displayDay)),
			mf.Element("button", mf.ElementProps{Class: "mt-1 truncate px-2 text-left text-xs font-medium hover:underline", Attrs: mf.Attrs{"type": "submit"}}, mf.Text(fmt.Sprintf("%d more", moreCount))),
		))
	}
	return div("h-28 w-full border border-solid border-base-300 text-left", children...)
}

func eventsForCalendarDay(day int) []calendarEvent {
	events := []calendarEvent{}
	for _, event := range calendarEvents {
		if event.Day == day {
			events = append(events, event)
		}
	}
	return events
}

func calendarThemeClass(theme string) string {
	switch theme {
	case "BLUE":
		return "bg-blue-200 dark:bg-blue-600 dark:text-blue-100"
	case "GREEN":
		return "bg-green-200 dark:bg-green-600 dark:text-green-100"
	case "PURPLE":
		return "bg-purple-200 dark:bg-purple-600 dark:text-purple-100"
	case "ORANGE":
		return "bg-orange-200 dark:bg-orange-600 dark:text-orange-100"
	case "PINK":
		return "bg-pink-200 dark:bg-pink-600 dark:text-pink-100"
	default:
		return ""
	}
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
	authFields := make([]dw.Field, 0, len(fields))
	for _, field := range fields {
		inputType := "text"
		if strings.Contains(strings.ToLower(field), "password") {
			inputType = "password"
		}
		authFields = append(authFields, inputField(dw.Field{Label: field, Type: inputType, Required: true}))
	}
	footerNodes := []mf.Node{}
	if helper != "" {
		footerNodes = append(footerNodes, div("mb-2 text-primary", mf.Text(helper)))
	}
	if footer != "" {
		footerNodes = append(footerNodes, mf.Text(footer))
	}
	return div("space-y-6",
		pageTitle(title, "DashWind user page preview rendered inside the Marionette demo shell.", nil),
		dw.AuthCard(dw.AuthCardProps{
			Title:       title,
			Description: "Preview account access flow",
			Fields:      authFields,
			SubmitLabel: title,
			Footer:      div("", footerNodes...),
		}),
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
			dw.SettingsSection(dw.SettingsSectionProps{
				Title:       "Account",
				Description: "Editable account fields",
				Fields: []dw.Field{
					inputField(dw.Field{Label: "Name", Type: "text", Value: "DashWind Admin", Required: true}),
					inputField(dw.Field{Label: "Email Id", Type: "email", Value: "admin@example.com", Required: true, Help: "Used for notifications and login."}),
					inputField(dw.Field{Label: "Role", Type: "text", Value: "Revenue Ops"}),
				},
			}),
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
	return documentationPage("Components", "DaisyUI primitives used by the sample.", []string{"dashwind.Shell, NavGroup, and NavItem", "dashwind.MetricGrid, MetricCard, CardPanel, PageHeader", "dashwind.DataTable plus DaisyUI ActionFormWithOptions"})
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

func inputField(field dw.Field) dw.Field {
	if strings.TrimSpace(field.Type) == "" {
		field.Type = "text"
	}
	if strings.TrimSpace(field.Name) == "" {
		field.Name = strings.ToLower(strings.ReplaceAll(field.Label, " ", "-"))
	}
	return field
}

func pageTitle(title, desc string, actions mf.Node) mf.Node {
	return dw.PageHeader(dw.PageHeaderProps{Title: title, Description: desc, Actions: actions})
}

func noticeNode(ctx *mb.Context) mf.Node {
	notice, _ := ctx.GetGlobal("notice").(string)
	if strings.TrimSpace(notice) == "" {
		return mf.Raw("")
	}
	ctx.SetGlobal("notice", "")
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
