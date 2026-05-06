package dashwinddemo

import (
	"fmt"
	"strconv"
	"strings"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
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
	return mf.Region(mf.RegionProps{ID: "dashwind-app"}, mf.DivProps(mf.ElementProps{Class: "drawer lg:drawer-open dashwind-shell"},
		mf.Raw(`<input id="dashwind-drawer" type="checkbox" class="drawer-toggle" />`),
		mf.DivProps(mf.ElementProps{Class: "drawer-content flex min-h-screen flex-col bg-base-200"}, topbar(current), mainContent(body)),
		mf.DivProps(mf.ElementProps{Class: "drawer-side z-40"}, mf.Raw(`<label for="dashwind-drawer" aria-label="close sidebar" class="drawer-overlay"></label>`), sidebar(current)),
	))
}
func mainContent(body mf.Node) mf.Node {
	return mf.Region(mf.RegionProps{ID: mainTargetID, Props: mf.ComponentProps{Class: "flex-1 p-4 md:p-6 lg:p-8 space-y-6"}}, body)
}
func topbar(current string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "navbar sticky top-0 z-30 border-b border-base-300 bg-base-100/90 backdrop-blur"}, mf.DivProps(mf.ElementProps{Class: "flex-none lg:hidden"}, mf.Raw(`<label for="dashwind-drawer" class="btn btn-square btn-ghost" aria-label="open sidebar">☰</label>`)), mf.DivProps(mf.ElementProps{Class: "flex-1"}, mf.H1Props(mf.ElementProps{Class: "text-xl font-semibold"}, mf.Text(current))), mf.DivProps(mf.ElementProps{Class: "hidden max-w-md flex-1 md:block"}, mf.Raw(`<label class="input input-bordered flex items-center gap-2"><span class="opacity-60">⌕</span><input type="search" class="grow" placeholder="Search DashWind demo" /></label>`)), mf.DivProps(mf.ElementProps{Class: "flex-none gap-2"}, mf.Raw(`<button class="btn btn-ghost btn-circle" onclick="mrnToggleTheme()" aria-label="toggle theme">◐</button><button class="btn btn-ghost btn-circle indicator" aria-label="notifications"><span class="indicator-item badge badge-primary badge-xs"></span>🔔</button><div class="avatar placeholder"><div class="bg-primary text-primary-content w-10 rounded-full"><span>DW</span></div></div>`)))
}
func sidebar(current string) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "min-h-full w-80 bg-base-100 text-base-content shadow-xl"}, mf.DivProps(mf.ElementProps{Class: "p-5"}, mf.DivProps(mf.ElementProps{Class: "mb-6 flex items-center gap-3"}, mf.DivProps(mf.ElementProps{Class: "grid h-11 w-11 place-items-center rounded-2xl bg-primary text-xl font-black text-primary-content"}, mf.Text("D")), mf.Div(mf.H2Props(mf.ElementProps{Class: "text-lg font-bold"}, mf.Text("DashWind")), mf.Raw(`<p class="text-xs text-base-content/60">DaisyUI admin template demo</p>`))), menuGroup("Menu", "main", current), menuGroup("Settings", "settings", current), mf.Raw(`<div class="mt-6 rounded-box bg-primary/10 p-4 text-sm"><p class="font-semibold">Marionette port</p><p class="mt-1 opacity-70">React/Redux template patterns rebuilt as Go handlers and htmx fragments.</p></div>`)))
}
func menuGroup(label, group, current string) mf.Node {
	items := []mf.Node{mf.Raw(`<li class="menu-title"><span>` + label + `</span></li>`)}
	for _, r := range routes {
		if r.Group != group {
			continue
		}
		active := ""
		if r.Name == current {
			active = " active"
		}
		items = append(items, mf.Raw(`<li><a class="`+active+`" href="`+r.Path+`"><span class="w-6 text-center">`+r.Icon+`</span>`+r.Name+`</a></li>`))
	}
	return mf.Element("ul", mf.ElementProps{Class: "menu rounded-box gap-1 p-0"}, items...)
}

func dashboardPage(ctx *mb.Context) mf.Node {
	notice := noticeNode(ctx)
	cards := make([]mf.Node, 0, len(statsData))
	for _, s := range statsData {
		cards = append(cards, statCardNode(s))
	}
	return mf.DivProps(mf.ElementProps{Class: "space-y-6"}, pageTitle("Dashboard", "Marionette rebuild of DashWind DashboardTopBar, Stats, Chart, and UserChannels sections.", periodForm(ctx.Get("period").(string))), notice, mf.DivProps(mf.ElementProps{Class: "grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-4"}, cards...), mf.DivProps(mf.ElementProps{Class: "grid grid-cols-1 gap-6 xl:grid-cols-2"}, chartCard("Revenue", "Monthly recurring revenue", mf.Chart(mf.ChartProps{Type: mf.ChartTypeLine, Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}, Height: 260, Datasets: []mf.ChartDataset{{Label: "MRR", Data: []float64{18, 24, 28, 32, 38, 45}, BorderColor: "#3b82f6", BackgroundColor: "rgba(59,130,246,.18)", Fill: true, Tension: .35}}})), chartCard("Pipeline", "Qualified leads by stage", mf.Chart(mf.ChartProps{Type: mf.ChartTypeBar, Labels: []string{"Open", "Progress", "Sold", "Followup"}, Height: 260, Datasets: []mf.ChartDataset{{Label: "Leads", Data: []float64{92, 128, 54, 76}, BackgroundColor: "#6366f1"}}, Options: mf.ChartOptions{BeginAtZero: true, HideLegend: true}}))), mf.DivProps(mf.ElementProps{Class: "grid grid-cols-1 gap-6 xl:grid-cols-2"}, amountStats(), userChannels()))
}
func statCardNode(s statCard) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "stats shadow bg-base-100"}, mf.DivProps(mf.ElementProps{Class: "stat"}, mf.DivProps(mf.ElementProps{Class: "stat-figure text-primary text-3xl"}, mf.Text(s.Icon)), mf.DivProps(mf.ElementProps{Class: "stat-title"}, mf.Text(s.Title)), mf.DivProps(mf.ElementProps{Class: "stat-value text-primary"}, mf.Text(s.Value)), mf.DivProps(mf.ElementProps{Class: "stat-desc font-medium " + s.TrendClass}, mf.Text(s.Description))))
}
func periodForm(period string) mf.Node {
	period = normalizePeriod(period)
	opts := []string{"Last 7 days", "Last 30 days", "This quarter"}
	children := []mf.Node{}
	for _, o := range opts {
		cls := "btn btn-sm"
		if o == period {
			cls += " btn-primary"
		}
		children = append(children, mf.Raw(`<button class="`+cls+`" name="period" value="`+o+`" hx-post="/dashboard/period" hx-target="#`+mainTargetID+`" hx-swap="outerHTML">`+o+`</button>`))
	}
	return mf.Element("form", mf.ElementProps{Class: "join", Attrs: mf.Attrs{"method": "post"}}, children...)
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
	return mf.DivProps(mf.ElementProps{Class: "stats stats-vertical lg:stats-horizontal bg-base-100 shadow"}, mf.Raw(`<div class="stat"><div class="stat-title">Total Likes</div><div class="stat-value">25.6K</div><div class="stat-desc">21% more than last month</div></div><div class="stat"><div class="stat-title">Page Views</div><div class="stat-value">2.6M</div><div class="stat-desc">14% more than last month</div></div>`))
}
func userChannels() mf.Node {
	rows := []string{"Organic search|12,432|46%", "Twitter|8,120|24%", "Newsletter|5,420|18%", "Partners|2,804|12%"}
	trs := ""
	for _, row := range rows {
		p := strings.Split(row, "|")
		trs += `<tr><td>` + p[0] + `</td><td>` + p[1] + `</td><td><progress class="progress progress-primary w-32" value="` + strings.TrimSuffix(p[2], "%") + `" max="100"></progress></td><td>` + p[2] + `</td></tr>`
	}
	return cardRaw("User Channels", "Traffic source breakdown", `<div class="overflow-x-auto"><table class="table"><tbody>`+trs+`</tbody></table></div>`)
}
func chartCard(title, desc string, chart mf.Node) mf.Node {
	return mf.Card(mf.CardProps{Title: title, Description: desc, Props: mf.ComponentProps{Class: "bg-base-100 shadow"}}, chart)
}

func leadsPage(ctx *mb.Context) mf.Node {
	leads := ctx.Get("leads").([]lead)
	rows := ""
	for _, l := range leads {
		rows += `<tr><td><div class="flex items-center gap-3"><div class="avatar placeholder"><div class="mask mask-squircle w-12 bg-neutral text-neutral-content"><span>` + l.Avatar + `</span></div></div><div><div class="font-bold">` + l.Name + `</div><div class="text-sm opacity-60">` + l.Role + `</div></div></div></td><td>` + l.Email + `</td><td>` + l.CreatedAt + `</td><td>` + statusBadge(l.Status) + `</td><td>` + l.Owner + `</td><td><form method="post"><input type="hidden" name="email" value="` + l.Email + `"><button class="btn btn-square btn-ghost btn-sm" hx-post="/leads/delete" hx-target="#` + mainTargetID + `" hx-swap="outerHTML">✕</button></form></td></tr>`
	}
	return mf.DivProps(mf.ElementProps{Class: "space-y-6"}, pageTitle("Current Leads", "DashWind leads table with htmx-powered Add New and delete actions.", mf.Raw(`<button class="btn btn-primary btn-sm" hx-post="/leads/add" hx-target="#`+mainTargetID+`" hx-swap="outerHTML">Add New</button>`)), noticeNode(ctx), cardRaw("Leads List", "Rendered from Marionette server state instead of a Redux slice/API call.", `<div class="overflow-x-auto"><table class="table"><thead><tr><th>Name</th><th>Email Id</th><th>Created At</th><th>Status</th><th>Assigned To</th><th></th></tr></thead><tbody>`+rows+`</tbody></table></div>`))
}
func transactionsPage(ctx *mb.Context) mf.Node {
	rows := ""
	total := 0
	for _, t := range seedTransactions {
		total += t.Amount
		rows += `<tr><td>` + t.Invoice + `</td><td>` + t.Customer + `</td><td>` + t.Plan + `</td><td>` + t.Date + `</td><td>` + statusBadge(t.Status) + `</td><td class="font-semibold">$` + strconv.Itoa(t.Amount) + `</td></tr>`
	}
	return mf.DivProps(mf.ElementProps{Class: "space-y-6"}, pageTitle("Transactions", "DashWind-style billing and transactions list.", mf.Raw(`<div class="stats shadow"><div class="stat"><div class="stat-title">Total</div><div class="stat-value text-primary">$`+strconv.Itoa(total)+`</div></div></div>`)), cardRaw("Recent Transactions", "", `<div class="overflow-x-auto"><table class="table table-zebra"><thead><tr><th>Invoice</th><th>Customer</th><th>Plan</th><th>Date</th><th>Status</th><th>Amount</th></tr></thead><tbody>`+rows+`</tbody></table></div>`))
}
func analyticsPage(ctx *mb.Context) mf.Node {
	return mf.DivProps(mf.ElementProps{Class: "space-y-6"}, pageTitle("Analytics", "DashWind charts page with Chart.js widgets.", nil), mf.DivProps(mf.ElementProps{Class: "grid grid-cols-1 gap-6 xl:grid-cols-2"}, chartCard("Doughnut", "Channel mix", mf.Chart(mf.ChartProps{Type: mf.ChartTypeDoughnut, Labels: []string{"Organic", "Social", "Referral", "Ads"}, Height: 280, Datasets: []mf.ChartDataset{{Label: "Users", Data: []float64{46, 24, 18, 12}}}})), chartCard("Scatter", "Lead score vs. ARR", mf.Chart(mf.ChartProps{Type: mf.ChartTypeScatter, Height: 280, Datasets: []mf.ChartDataset{{Label: "Accounts", Points: []mf.ChartPoint{{X: 20, Y: 15}, {X: 42, Y: 38}, {X: 60, Y: 72}, {X: 82, Y: 120}}, BackgroundColor: "#14b8a6"}}}))))
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
	lis := ""
	for _, it := range items {
		lis += `<li>` + it + `</li>`
	}
	return mf.DivProps(mf.ElementProps{Class: "space-y-6"}, pageTitle(title, desc, nil), cardRaw(title, "", `<ul class="list-disc space-y-2 pl-5">`+lis+`</ul>`))
}
func pageTitle(title, desc string, actions mf.Node) mf.Node {
	children := []mf.Node{mf.Div(mf.H1Props(mf.ElementProps{Class: "text-3xl font-bold"}, mf.Text(title)), mf.Raw(`<p class="mt-1 text-base-content/60">`+desc+`</p>`))}
	if actions != nil {
		children = append(children, actions)
	}
	return mf.DivProps(mf.ElementProps{Class: "flex flex-col gap-4 md:flex-row md:items-center md:justify-between"}, children...)
}
func noticeNode(ctx *mb.Context) mf.Node {
	notice, _ := ctx.Get("notice").(string)
	if strings.TrimSpace(notice) == "" {
		return mf.Raw("")
	}
	ctx.Set("notice", "")
	return mf.Raw(`<div class="alert alert-success shadow"><span>` + notice + `</span></div>`)
}
func cardRaw(title, desc, body string) mf.Node {
	sub := ""
	if desc != "" {
		sub = `<p class="text-sm text-base-content/60">` + desc + `</p>`
	}
	return mf.Raw(`<div class="card bg-base-100 shadow"><div class="card-body"><h2 class="card-title">` + title + `</h2>` + sub + body + `</div></div>`)
}
func statusBadge(status string) string {
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
	return `<span class="badge ` + class + `">` + status + `</span>`
}

const dashwindCSS = `
#marionette-root { width: 100%; max-width: none; padding: 0; }
#marionette-root > * { animation: none; }
.dashwind-shell .drawer-side .menu a.active { background: var(--color-primary); color: var(--color-primary-content); font-weight: 700; }
.dashwind-shell .card, .dashwind-shell .stats { border: 1px solid color-mix(in oklab, var(--color-base-content) 10%, transparent); }
.dashwind-shell .stat-value { font-size: clamp(1.75rem, 2vw, 2.25rem); }
`
