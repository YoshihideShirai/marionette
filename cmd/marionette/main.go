package main

import (
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

type user struct {
	ID        int
	Name      string
	Email     string
	Role      string
	StartDate string
}

const (
	startDateMin = "2024-01-01"
	startDateMax = "2026-12-31"
)

type createUserFormState struct {
	Name      string
	Email     string
	Role      string
	StartDate string
	Errors    map[string]string
}

type pagination struct {
	Page       int
	PerPage    int
	TotalPages int
}

func defaultCreateUserFormState() createUserFormState {
	return createUserFormState{Role: "Viewer", Errors: map[string]string{}}
}

func demoUsers() []user {
	return []user{
		{ID: 1, Name: "Aiko Tanaka", Email: "aiko@example.com", Role: "Admin", StartDate: "2024-03-18"},
		{ID: 2, Name: "Ren Sato", Email: "ren@example.com", Role: "Editor", StartDate: "2024-07-01"},
		{ID: 3, Name: "Mina Suzuki", Email: "mina@example.com", Role: "Viewer", StartDate: "2025-01-10"},
	}
}

func main() {
	app := buildApp()
	if err := app.Run("127.0.0.1:8080"); err != nil {
		log.Fatal(err)
	}
}

func buildApp() *mb.App {
	app := mb.New()
	app.AddStyle(customCSSDemoStyles())
	app.Set("nextUserID", 4)
	app.Set("users", demoUsers())
	app.Set("deleteModalOpen", false)
	app.Set("deleteTargetID", 0)
	app.Set("loading", false)

	app.Page("/", func(ctx *mb.Context) mf.Node {
		return renderDashboardPage(ctx)
	})
	app.Page("/users", func(ctx *mb.Context) mf.Node {
		return renderUsersPage(ctx, defaultCreateUserFormState())
	})
	app.Page("/analytics", func(ctx *mb.Context) mf.Node {
		return renderAnalyticsPage(ctx)
	})
	app.Page("/settings", func(ctx *mb.Context) mf.Node {
		return renderSettingsPage(ctx)
	})

	app.Action("users/create", func(ctx *mb.Context) mf.Node {
		ctx.Set("loading", false)
		form := createUserFormState{
			Name:      strings.TrimSpace(ctx.FormValue("name")),
			Email:     strings.TrimSpace(ctx.FormValue("email")),
			Role:      strings.TrimSpace(ctx.FormValue("role")),
			StartDate: strings.TrimSpace(ctx.FormValue("start_date")),
			Errors:    map[string]string{},
		}
		if form.Role == "" {
			form.Role = "Viewer"
		}

		if strings.TrimSpace(form.Name) == "" {
			form.Errors["name"] = "Name is required."
		}
		if strings.TrimSpace(form.Email) == "" {
			form.Errors["email"] = "Email is required."
		}
		if errMsg := validateStartDate(form.StartDate); errMsg != "" {
			form.Errors["start_date"] = errMsg
		}

		if len(form.Errors) == 0 {
			users := getUsers(ctx)
			nextID := ctx.GetInt("nextUserID")
			users = append(users, user{ID: nextID, Name: form.Name, Email: form.Email, Role: form.Role, StartDate: form.StartDate})
			ctx.Set("users", users)
			ctx.Set("nextUserID", nextID+1)
			ctx.FlashSuccess("User was saved successfully.")
			return renderUsersWorkspace(ctx, defaultCreateUserFormState())
		}

		ctx.FlashError("Save failed. Please fix the highlighted fields.")
		return renderUsersWorkspace(ctx, form)
	})

	app.Action("users/delete/prompt", func(ctx *mb.Context) mf.Node {
		ctx.Set("loading", false)
		id, _ := strconv.Atoi(ctx.FormValue("id"))
		ctx.Set("deleteTargetID", id)
		ctx.Set("deleteModalOpen", true)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/delete/cancel", func(ctx *mb.Context) mf.Node {
		ctx.Set("loading", false)
		ctx.Set("deleteModalOpen", false)
		ctx.Set("deleteTargetID", 0)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/delete/confirm", func(ctx *mb.Context) mf.Node {
		ctx.Set("loading", false)
		id, _ := strconv.Atoi(ctx.FormValue("id"))
		users := getUsers(ctx)
		next := users[:0]
		for _, u := range users {
			if u.ID != id {
				next = append(next, u)
			}
		}
		ctx.Set("users", next)
		ctx.Set("deleteModalOpen", false)
		ctx.Set("deleteTargetID", 0)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/loading/start", func(ctx *mb.Context) mf.Node {
		ctx.Set("loading", true)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/loading/stop", func(ctx *mb.Context) mf.Node {
		ctx.Set("loading", false)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/reset", func(ctx *mb.Context) mf.Node {
		ctx.Set("users", demoUsers())
		ctx.Set("nextUserID", 4)
		ctx.Set("deleteModalOpen", false)
		ctx.Set("deleteTargetID", 0)
		ctx.Set("loading", false)
		ctx.FlashInfo("Demo data was reset.")
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	return app
}

func validateStartDate(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return "Start date is required."
	}
	selected, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return "Enter a valid date (YYYY-MM-DD)."
	}
	minDate, _ := time.Parse("2006-01-02", startDateMin)
	maxDate, _ := time.Parse("2006-01-02", startDateMax)
	if selected.Before(minDate) || selected.After(maxDate) {
		return "Start date must be between " + startDateMin + " and " + startDateMax + "."
	}
	return ""
}

func getUsers(ctx *mb.Context) []user {
	users, ok := ctx.Get("users").([]user)
	if !ok {
		return nil
	}
	return users
}

func renderUsersPage(ctx *mb.Context, formState createUserFormState) mf.Node {
	header := mf.PageHeader(mf.PageHeaderProps{
		Title:       "User Administration",
		Description: "A CRUD demo combining forms, tables, and modals.",
		Actions:     themeToggleButton(),
	})
	return renderShell(ctx, "/users", header, renderUsersWorkspace(ctx, formState))
}

func renderSidebar(active string) mf.Node {
	return mf.Sidebar("Marionette", "Admin Console",
		sidebarLink("Dashboard", "/", active == "/"),
		sidebarLink("Users", "/users", active == "/users"),
		sidebarLink("Analytics", "/analytics", active == "/analytics"),
		sidebarLink("Settings", "/settings", active == "/settings"),
	).Note("Demo workspace", "In-memory data for admin UI prototyping.")
}

func sidebarLink(label, href string, active bool) mf.SidebarItem {
	link := mf.SidebarLink(label, href)
	if active {
		return link.Active()
	}
	return link
}

func renderShell(ctx *mb.Context, activePath string, header, content mf.Node) mf.Node {
	return mf.AppShell(mf.AppShellProps{
		Sidebar: renderSidebar(activePath),
		Flashes: mf.FlashAlerts(ctx.Flashes()),
		Header:  header,
		Content: content,
	})
}
func themeToggleButton() mf.Node {
	return mf.ThemeToggleButton(mf.ComponentProps{Variant: "outline", Size: "sm"})
}

func renderDashboardPage(ctx *mb.Context) mf.Node {
	users := getUsers(ctx)
	header := mf.PageHeader(mf.PageHeaderProps{
		Title:       "Operations Dashboard",
		Description: "Top-level overview of operational health and activity.",
		Actions:     themeToggleButton(),
	})
	content := mf.Stack(mf.StackProps{Gap: "lg"},
		renderDashboardOverview(ctx),
		renderCustomCSSDemo(),
		renderDashboardCharts(ctx),
		mf.Grid(mf.GridProps{MinColumnWidth: "lg", Gap: "lg"},
			mf.Card(mf.CardProps{}, mf.EmptyState(mf.EmptyStateProps{
				Title:       "Pending approvals",
				Description: "No approvals are currently waiting.",
				Icon:        "✓",
			})),
			mf.Card(mf.CardProps{}, mf.Table(mf.TableProps{
				Columns: []mf.TableColumn{{Label: "Metric"}, {Label: "Value"}},
				Rows: []mf.TableComponentRow{
					mf.TableRowValues("Total users", len(users)),
					mf.TableRowValues("Latest start", latestStartDate(users)),
				},
			})),
		),
	)
	return renderShell(ctx, "/", header, content)
}

func customCSSDemoStyles() string {
	return `
		html[data-mrn-css-profile="focus"] #marionette-root {
			max-width: 72rem;
		}

		html[data-mrn-css-profile="focus"] .mrn-css-demo-card {
			border-color: color-mix(in oklab, var(--color-primary) 45%, transparent);
			box-shadow: 0 18px 40px color-mix(in oklab, var(--color-primary) 18%, transparent);
		}

		html[data-mrn-css-profile="focus"] .mrn-css-demo-preview {
			background: linear-gradient(135deg, color-mix(in oklab, var(--color-primary) 16%, var(--color-base-100)), var(--color-base-100));
		}

		html[data-mrn-css-profile="compact"] #marionette-root {
			max-width: 82rem;
		}

		html[data-mrn-css-profile="compact"] .mrn-css-demo-card .card-body {
			gap: 0.75rem;
			padding: 1rem;
		}

		html[data-mrn-css-profile="compact"] .mrn-css-demo-preview {
			display: grid;
			grid-template-columns: repeat(3, minmax(0, 1fr));
			gap: 0.5rem;
		}

		html[data-mrn-css-profile="editorial"] .mrn-css-demo-card {
			border-radius: 0.25rem;
			border-color: color-mix(in oklab, var(--color-secondary) 50%, transparent);
		}

		html[data-mrn-css-profile="editorial"] .mrn-css-demo-preview {
			background: var(--color-base-100);
			border-left: 0.35rem solid var(--color-secondary);
		}

		html[data-mrn-css-profile="editorial"] .mrn-css-demo-preview strong {
			font-family: Georgia, serif;
		}
	`
}

func renderCustomCSSDemo() mf.Node {
	return mf.Card(mf.CardProps{
		Title:       "Custom CSS profiles",
		Description: "App-level CSS can reshape Marionette without changing component code.",
		Props:       mf.ComponentProps{Class: "mrn-css-demo-card"},
	},
		mf.Raw(`<script>
			(function() {
				var key = "marionette-css-profile";
				var root = document.documentElement;
				var saved = "";
				try { saved = localStorage.getItem(key) || ""; } catch (e) {}
				if (saved) root.setAttribute("data-mrn-css-profile", saved);
				window.mrnSetCSSProfile = function(profile) {
					if (profile) {
						root.setAttribute("data-mrn-css-profile", profile);
						try { localStorage.setItem(key, profile); } catch (e) {}
						return;
					}
					root.removeAttribute("data-mrn-css-profile");
					try { localStorage.removeItem(key); } catch (e) {}
				};
			})();
		</script>`),
		mf.Box(mf.BoxProps{Border: true, Tone: "base", Padding: "md", Props: mf.ComponentProps{Class: "mrn-css-demo-preview rounded-box"}},
			mf.UIText(mf.TextProps{Text: "Same components, different product feel", Weight: "semibold"}),
			mf.UIText(mf.TextProps{Text: "The CSS is registered once with app.AddStyle and toggled by a data attribute.", Size: "sm", Tone: "muted"}),
			mf.Actions(mf.ActionsProps{Gap: "sm", Wrap: true},
				mf.Badge(mf.BadgeProps{Label: "Tailwind classes stay unchanged", Props: mf.ComponentProps{Variant: "primary"}}),
				mf.Badge(mf.BadgeProps{Label: "Scoped app branding", Props: mf.ComponentProps{Variant: "secondary"}}),
				mf.Badge(mf.BadgeProps{Label: "Runtime switchable", Props: mf.ComponentProps{Variant: "accent"}}),
			),
		),
		mf.Raw(`<div class="flex flex-wrap gap-2">
			<button class="btn btn-sm btn-primary" type="button" onclick="window.mrnSetCSSProfile && window.mrnSetCSSProfile('focus')">Focus</button>
			<button class="btn btn-sm btn-secondary" type="button" onclick="window.mrnSetCSSProfile && window.mrnSetCSSProfile('compact')">Compact</button>
			<button class="btn btn-sm btn-accent" type="button" onclick="window.mrnSetCSSProfile && window.mrnSetCSSProfile('editorial')">Editorial</button>
			<button class="btn btn-sm btn-outline" type="button" onclick="window.mrnSetCSSProfile && window.mrnSetCSSProfile('')">Reset</button>
		</div>`),
	)
}

func renderAnalyticsPage(ctx *mb.Context) mf.Node {
	header := mf.PageHeader(mf.PageHeaderProps{Title: "Analytics", Description: "A chart-first, two-column analytics layout.", Actions: themeToggleButton()})
	content := mf.Split(mf.SplitProps{
		Gap: "lg",
		Main: mf.Stack(mf.StackProps{Gap: "lg"},
			renderDashboardCharts(ctx),
			mf.Alert(mf.AlertProps{Title: "Insight", Description: "Editors are growing steadily month-over-month.", Props: mf.ComponentProps{Variant: "info"}}),
		),
		Aside: mf.Stack(mf.StackProps{Gap: "md"},
			roleMixRow("Admin", roleCounts(getUsers(ctx))["Admin"]),
			roleMixRow("Editor", roleCounts(getUsers(ctx))["Editor"]),
			roleMixRow("Viewer", roleCounts(getUsers(ctx))["Viewer"]),
		),
	})
	return renderShell(ctx, "/analytics", header, content)
}

func renderSettingsPage(ctx *mb.Context) mf.Node {
	header := mf.PageHeader(mf.PageHeaderProps{Title: "Settings", Description: "A form-heavy vertical layout for configuration flows.", Actions: themeToggleButton()})
	content := mf.Stack(mf.StackProps{Gap: "lg"},
		mf.Card(mf.CardProps{},
			mf.ActionForm(mf.ActionFormProps{Action: "users/reset", Target: "#app", Swap: "outerHTML"},
				mf.FormRow(mf.FormRowProps{ID: "org-name", Label: "Organization name", Control: mf.TextField(mf.TextFieldProps{ID: "org-name", Name: "org-name", Value: "Marionette Labs"})}),
				mf.FormRow(mf.FormRowProps{ID: "timezone", Label: "Timezone", Control: mf.Select(mf.SelectFieldProps{ID: "timezone", Name: "timezone", Options: []mf.SelectOption{{Label: "UTC", Value: "UTC", Selected: true}, {Label: "Asia/Tokyo", Value: "Asia/Tokyo"}}})}),
				mf.Switch(mf.SwitchProps{ID: "audit", Name: "audit", Value: "yes", Checked: true, Label: "Enable audit log"}),
				mf.Actions(mf.ActionsProps{Props: mf.ComponentProps{Class: "pt-3"}},
					mf.SubmitButton("Save settings", mf.ComponentProps{Variant: "primary"}),
				),
			),
		),
		mf.Toast(mf.ToastProps{Title: "Tip", Description: "This page demonstrates form-heavy layout variation.", Props: mf.ComponentProps{Variant: "info"}}),
	)
	return renderShell(ctx, "/settings", header, content)
}

func renderUsersWorkspace(ctx *mb.Context, formState createUserFormState) mf.Node {
	return mf.Region(mf.RegionProps{ID: "users-workspace", Props: mf.ComponentProps{Class: "space-y-4"}},
		mf.FlashAlerts(ctx.Flashes()),
		renderDashboardOverview(ctx),
		renderDashboardCharts(ctx),
		mf.Split(mf.SplitProps{
			Main:  renderUsersTable(ctx),
			Aside: renderCreateUserForm(formState),
			Gap:   "lg",
		}),
		renderComponentShowcase(ctx),
		renderDeleteModal(ctx),
	)
}

func renderDashboardCharts(ctx *mb.Context) mf.Node {
	users := getUsers(ctx)
	monthLabels, monthData := monthlyStartCounts(users)
	roleLabels, roleData := roleDistribution(users)

	return mf.Grid(mf.GridProps{MinColumnWidth: "md", Gap: "lg"},
		mf.Chart(mf.ChartProps{
			Type:        mf.ChartTypeLine,
			Title:       "Onboarding trend",
			Description: "Users grouped by start month. This chart re-renders after htmx actions.",
			Labels:      monthLabels,
			Datasets: []mf.ChartDataset{
				{
					Label:           "Starts",
					Data:            monthData,
					BorderColor:     "#2563eb",
					BackgroundColor: "rgba(37, 99, 235, 0.16)",
					Fill:            true,
					Tension:         0.35,
				},
			},
			Options: mf.ChartOptions{
				BeginAtZero: true,
				YAxisLabel:  "Users",
			},
			Height: 260,
		}),
		mf.Chart(mf.ChartProps{
			Type:        mf.ChartTypeBar,
			Title:       "Role distribution",
			Description: "Current permission mix rendered from the same in-memory users.",
			Labels:      roleLabels,
			Datasets: []mf.ChartDataset{
				{
					Label:           "Users",
					Data:            roleData,
					BorderColor:     "#7c3aed",
					BackgroundColor: "rgba(124, 58, 237, 0.24)",
				},
			},
			Options: mf.ChartOptions{
				BeginAtZero: true,
				HideLegend:  true,
				YAxisLabel:  "Users",
			},
			Height: 260,
		}),
	)
}

func renderDashboardOverview(ctx *mb.Context) mf.Node {
	users := getUsers(ctx)
	roles := roleCounts(users)
	latest := latestStartDate(users)

	return mf.Grid(mf.GridProps{Columns: "3"},
		statCard("Users", strconv.Itoa(len(users)), "Active demo records", "primary"),
		statCard("Admins", strconv.Itoa(roles["Admin"]), "High-access seats", "secondary"),
		statCard("Latest start", latest, "Newest onboarding date", "accent"),
	)
}

func statCard(label, value, caption, tone string) mf.Node {
	return mf.Card(mf.CardProps{},
		mf.Stack(mf.StackProps{Gap: "sm"},
			mf.UIText(mf.TextProps{Text: label, Size: "sm", Weight: "medium", Tone: "muted"}),
			mf.Stack(mf.StackProps{Direction: "horizontal", Gap: "md", Align: "end", Justify: "between"},
				mf.UIText(mf.TextProps{Text: value, Size: "3xl", Weight: "bold"}),
				mf.Badge(mf.BadgeProps{Label: "live", Props: mf.ComponentProps{Variant: tone}}),
			),
			mf.UIText(mf.TextProps{Text: caption, Size: "sm", Tone: "muted"}),
		),
	)
}

func renderUsersTable(ctx *mb.Context) mf.Node {
	users := getUsers(ctx)
	loading := isLoading(ctx)
	pg := parsePagination(ctx.Query("page"), ctx.Query("per_page"), len(users))
	tableBody := renderUsersTableBody(users, loading, ctx.Query("sort"), pg)

	return mf.Card(mf.CardProps{
		Title:       "Users",
		Description: "Create and remove users with htmx-backed actions.",
		Actions: mf.Actions(mf.ActionsProps{Gap: "sm"},
			mf.Badge(mf.BadgeProps{Label: strconv.Itoa(len(getUsers(ctx))) + " total", Props: mf.ComponentProps{Variant: "outline"}}),
			mf.ActionForm(mf.ActionFormProps{Action: "users/loading/start", Target: "#users-workspace", Swap: "outerHTML"},
				mf.SubmitButton("Show loading", mf.ComponentProps{Variant: "ghost", Size: "sm", Disabled: loading}),
			),
			mf.ActionForm(mf.ActionFormProps{Action: "users/loading/stop", Target: "#users-workspace", Swap: "outerHTML"},
				mf.SubmitButton("Show data", mf.ComponentProps{Variant: "ghost", Size: "sm", Disabled: !loading}),
			),
			mf.ActionForm(mf.ActionFormProps{Action: "users/reset", Target: "#users-workspace", Swap: "outerHTML"},
				mf.SubmitButton("Reset", mf.ComponentProps{Variant: "secondary", Size: "sm"}),
			),
		),
	},
		mf.Box(mf.BoxProps{Border: true, Padding: "none", Props: mf.ComponentProps{Class: "overflow-hidden rounded-box"}}, tableBody),
		mf.Pagination(mf.PaginationProps{
			Page:       pg.Page,
			TotalPages: pg.TotalPages,
			PrevHref:   pageLink(pg.Page-1, pg.PerPage, ctx.Query("sort"), pg.TotalPages),
			NextHref:   pageLink(pg.Page+1, pg.PerPage, ctx.Query("sort"), pg.TotalPages),
		}),
	)
}

func renderUsersTableBody(users []user, loading bool, sortKey string, pg pagination) mf.Node {
	if loading {
		return mf.EmptyState(mf.EmptyStateProps{
			Skeleton: true,
			Rows:     5,
		})
	}

	sorted := paginateUsers(sortUsers(users, sortKey), pg)
	rows := make([]mf.TableComponentRow, 0, len(sorted))
	for _, u := range sorted {
		rows = append(rows, renderUserRow(u))
	}

	return mf.Table(mf.TableProps{
		Columns:          usersTableColumns(sortKey, pg),
		Rows:             rows,
		EmptyTitle:       "No users yet",
		EmptyDescription: "Create a user from the form to populate this table.",
	})
}

func parsePagination(pageRaw, perPageRaw string, total int) pagination {
	page := 1
	if p, err := strconv.Atoi(strings.TrimSpace(pageRaw)); err == nil && p > 0 {
		page = p
	}
	perPage := 5
	if pp, err := strconv.Atoi(strings.TrimSpace(perPageRaw)); err == nil && pp > 0 {
		perPage = pp
	}
	totalPages := total / perPage
	if total%perPage != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}
	return pagination{Page: page, PerPage: perPage, TotalPages: totalPages}
}

func paginateUsers(users []user, pg pagination) []user {
	start := (pg.Page - 1) * pg.PerPage
	if start >= len(users) {
		return nil
	}
	end := start + pg.PerPage
	if end > len(users) {
		end = len(users)
	}
	return users[start:end]
}

func pageLink(page, perPage int, sortKey string, totalPages int) string {
	if page < 1 || page > totalPages {
		return ""
	}
	query := url.Values{}
	query.Set("page", strconv.Itoa(page))
	query.Set("per_page", strconv.Itoa(perPage))
	if strings.TrimSpace(sortKey) != "" {
		query.Set("sort", sortKey)
	}
	return "/users?" + query.Encode()
}

func isLoading(ctx *mb.Context) bool {
	v, _ := ctx.Get("loading").(bool)
	return v
}

func renderUserRow(u user) mf.TableComponentRow {
	return mf.TableRowValues(
		mf.UIText(mf.TextProps{Text: u.Name, Weight: "medium"}),
		mf.UIText(mf.TextProps{Text: u.Email, Size: "sm", Tone: "subtle"}),
		mf.Badge(mf.BadgeProps{Label: u.Role, Props: mf.ComponentProps{Variant: "ghost"}}),
		mf.UIText(mf.TextProps{Text: u.StartDate, Size: "sm"}),
		mf.ActionForm(mf.ActionFormProps{Action: "users/delete/prompt", Target: "#users-workspace", Swap: "outerHTML"},
			mf.HiddenField("id", strconv.Itoa(u.ID)),
			mf.SubmitButton("Delete", mf.ComponentProps{Variant: "danger", Size: "sm"}),
		),
	)
}

func usersTableColumns(activeSort string, pg pagination) []mf.TableColumn {
	sortedQuery := func(key string) string {
		query := url.Values{}
		query.Set("sort", key)
		query.Set("page", strconv.Itoa(pg.Page))
		query.Set("per_page", strconv.Itoa(pg.PerPage))
		return "/users?" + query.Encode()
	}
	return []mf.TableColumn{
		{Label: "Name", SortKey: "name", SortHref: sortedQuery("name"), SortActive: activeSort == "name"},
		{Label: "Email", SortKey: "email", SortHref: sortedQuery("email"), SortActive: activeSort == "email"},
		{Label: "Role", SortKey: "role", SortHref: sortedQuery("role"), SortActive: activeSort == "role"},
		{Label: "Start date", SortKey: "start_date", SortHref: sortedQuery("start_date"), SortActive: activeSort == "start_date"},
		{Label: ""},
	}
}

func sortUsers(users []user, sortKey string) []user {
	sorted := make([]user, len(users))
	copy(sorted, users)

	less := func(i, j int) bool { return sorted[i].ID < sorted[j].ID }
	switch sortKey {
	case "name":
		less = func(i, j int) bool { return sorted[i].Name < sorted[j].Name }
	case "email":
		less = func(i, j int) bool { return sorted[i].Email < sorted[j].Email }
	case "role":
		less = func(i, j int) bool { return sorted[i].Role < sorted[j].Role }
	case "start_date":
		less = func(i, j int) bool { return sorted[i].StartDate < sorted[j].StartDate }
	}

	sort.SliceStable(sorted, less)
	return sorted
}

func roleCounts(users []user) map[string]int {
	counts := map[string]int{"Admin": 0, "Editor": 0, "Viewer": 0}
	for _, u := range users {
		counts[u.Role]++
	}
	return counts
}

func monthlyStartCounts(users []user) ([]string, []float64) {
	counts := map[string]int{}
	months := make([]string, 0)
	for _, u := range users {
		month := startMonth(u.StartDate)
		if month == "" {
			continue
		}
		if _, ok := counts[month]; !ok {
			months = append(months, month)
		}
		counts[month]++
	}
	sort.Strings(months)
	if len(months) == 0 {
		return []string{"No data"}, []float64{0}
	}

	data := make([]float64, 0, len(months))
	for _, month := range months {
		data = append(data, float64(counts[month]))
	}
	return months, data
}

func startMonth(raw string) string {
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(raw))
	if err != nil {
		return ""
	}
	return parsed.Format("2006-01")
}

func roleDistribution(users []user) ([]string, []float64) {
	roles := roleCounts(users)
	labels := []string{"Admin", "Editor", "Viewer"}
	data := make([]float64, 0, len(labels))
	for _, label := range labels {
		data = append(data, float64(roles[label]))
	}
	return labels, data
}

func latestStartDate(users []user) string {
	latest := "n/a"
	for _, u := range users {
		if u.StartDate > latest || latest == "n/a" {
			latest = u.StartDate
		}
	}
	return latest
}

func roleBadge(role string) mf.Node {
	tone := "ghost"
	switch role {
	case "Admin":
		tone = "primary"
	case "Editor":
		tone = "secondary"
	}
	return mf.Badge(mf.BadgeProps{Label: role, Props: mf.ComponentProps{Variant: tone}})
}

func renderCreateUserForm(form createUserFormState) mf.Node {
	return mf.Card(mf.CardProps{Title: "Create user"},
		mf.ActionForm(mf.ActionFormProps{Action: "users/create", Target: "#users-workspace", Swap: "outerHTML"},
			mf.FormRow(mf.FormRowProps{
				ID:          "name",
				Label:       "Name",
				Description: "Enter the display name.",
				Error:       form.Errors["name"],
				Required:    true,
				Control: mf.TextField(mf.TextFieldProps{
					ID:          "name",
					Name:        "name",
					Value:       form.Name,
					Placeholder: "name",
					Description: "Enter the display name.",
					Error:       form.Errors["name"],
					Required:    true,
				}),
			}),
			mf.FormRow(mf.FormRowProps{
				ID:          "email",
				Label:       "Email",
				Description: "Used for notifications.",
				Error:       form.Errors["email"],
				Required:    true,
				Control: mf.TextField(mf.TextFieldProps{
					ID:          "email",
					Name:        "email",
					Value:       form.Email,
					Placeholder: "email",
					Description: "Used for notifications.",
					Error:       form.Errors["email"],
					Required:    true,
				}),
			}),
			mf.FormRow(mf.FormRowProps{
				ID:          "start_date",
				Label:       "Start date",
				Description: "Select a date in the active fiscal window.",
				Error:       form.Errors["start_date"],
				Required:    true,
				Control: mf.TextField(mf.TextFieldProps{
					ID:          "start_date",
					Name:        "start_date",
					Value:       form.StartDate,
					Type:        "date",
					Description: "Select a date in the active fiscal window.",
					Error:       form.Errors["start_date"],
					Required:    true,
				}),
			}),
			mf.FormRow(mf.FormRowProps{
				ID:          "role",
				Label:       "Role",
				Description: "Choose permission scope for this user.",
				Error:       form.Errors["role"],
				Required:    true,
				Control: mf.Select(mf.SelectFieldProps{
					ID:   "role",
					Name: "role",
					Options: []mf.SelectOption{
						{Label: "Admin", Value: "Admin", Selected: form.Role == "Admin"},
						{Label: "Editor", Value: "Editor", Selected: form.Role == "Editor"},
						{Label: "Viewer", Value: "Viewer", Selected: form.Role == "" || form.Role == "Viewer"},
					},
					Description: "Choose permission scope for this user.",
					Error:       form.Errors["role"],
					Required:    true,
				}),
			}),
			mf.Divider(mf.DividerProps{Spacing: "xs"}),
			mf.FormRow(mf.FormRowProps{
				ID:          "workspace",
				Label:       "Workspace",
				Description: "Demo-only field showing radio group markup.",
				Control: mf.RadioGroup(mf.RadioGroupProps{
					ID:          "workspace",
					Name:        "workspace",
					Value:       "core",
					Description: "Demo-only field showing radio group markup.",
					Options: []mf.RadioOption{
						{Label: "Core", Value: "core"},
						{Label: "Growth", Value: "growth"},
						{Label: "Support", Value: "support"},
					},
				}),
			}),
			mf.FormRow(mf.FormRowProps{
				ID:          "notes",
				Label:       "Notes",
				Description: "Ignored by the demo action, useful for layout coverage.",
				Control: mf.Textarea(mf.TextareaProps{
					ID:          "notes",
					Name:        "notes",
					Placeholder: "Add onboarding context",
					Rows:        3,
					Description: "Ignored by the demo action, useful for layout coverage.",
				}),
			}),
			mf.Stack(mf.StackProps{Gap: "sm"},
				mf.Checkbox(mf.CheckboxProps{
					ID:          "send_invite",
					Name:        "send_invite",
					Value:       "yes",
					Checked:     true,
					Label:       "Send invite email",
					Description: "Invite preference.",
				}),
				mf.Switch(mf.SwitchProps{
					ID:          "provision_access",
					Name:        "provision_access",
					Value:       "yes",
					Checked:     true,
					Label:       "Provision default access",
					Description: "Access preference.",
				}),
			),
			mf.Actions(mf.ActionsProps{Gap: "sm", Wrap: true, Props: mf.ComponentProps{Class: "pt-2"}},
				mf.SubmitButton("Create", mf.ComponentProps{Variant: "primary", Size: "sm"}),
				mf.Button("Preview", mf.ComponentProps{Variant: "ghost", Size: "sm", Disabled: true}),
			),
		),
	)
}

func renderComponentShowcase(ctx *mb.Context) mf.Node {
	users := getUsers(ctx)
	roles := roleCounts(users)
	seatLimit := 8

	return mf.Grid(mf.GridProps{Gap: "lg", Props: mf.ComponentProps{Class: "xl:grid-cols-[minmax(0,1fr)_22rem]"}},
		mf.Card(mf.CardProps{
			Title:       "Component states",
			Description: "Alerts, toasts, skeletons, and empty states rendered from Go.",
		},
			mf.Grid(mf.GridProps{Columns: "2", Gap: "sm"},
				mf.Alert(mf.AlertProps{
					Title:       "Validation feedback",
					Description: "Form errors, success flashes, and informational alerts share the same feedback API.",
					Icon:        "!",
					Props:       mf.ComponentProps{Variant: "info", Size: "sm"},
				}),
				mf.Toast(mf.ToastProps{
					Title:       "Saved",
					Description: "Toast markup is available for transient status messages.",
					Icon:        "OK",
					Props:       mf.ComponentProps{Variant: "success", Size: "sm"},
				}),
			),
			mf.Grid(mf.GridProps{Columns: "2", Gap: "sm"},
				mf.Skeleton(mf.SkeletonProps{Rows: 4, Props: mf.ComponentProps{Size: "sm"}}),
				mf.EmptyState(mf.EmptyStateProps{
					Title:       "No pending reviews",
					Description: "Empty states can replace tables or panels without extra branching in templates.",
					Icon:        "0",
					Props:       mf.ComponentProps{Size: "sm"},
				}),
			),
			mf.Box(mf.BoxProps{Border: true, Tone: "base", Padding: "sm", Props: mf.ComponentProps{Class: "rounded-box"}},
				mf.UIText(mf.TextProps{Text: "Theme toggle component", Size: "sm", Tone: "muted", Props: mf.ComponentProps{Class: "mb-2"}}),
				mf.ThemeToggleButton(mf.ComponentProps{Variant: "outline", Size: "sm"}),
			),
			mf.Box(mf.BoxProps{Border: true, Tone: "base", Padding: "sm", Props: mf.ComponentProps{Class: "rounded-box"}},
				mf.Image(mf.ImageProps{
					Src:         "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=900&q=80",
					Alt:         "Desk with laptop and notebook",
					Caption:     "Image component with responsive framing and a caption.",
					Width:       900,
					Height:      600,
					AspectRatio: "video",
					ObjectFit:   "cover",
				}),
			),
		),
		mf.Card(mf.CardProps{
			Title:       "Role mix",
			Description: "Small repeated views stay plain Go functions.",
		},
			roleMixRow("Admin", roles["Admin"]),
			roleMixRow("Editor", roles["Editor"]),
			roleMixRow("Viewer", roles["Viewer"]),
			mf.Divider(mf.DividerProps{Spacing: "sm"}),
			mf.Progress(mf.ProgressProps{
				Value:     float64(len(users)),
				Max:       float64(seatLimit),
				Label:     "Seats used",
				ShowValue: true,
				Props:     mf.ComponentProps{Variant: "primary", Size: "lg"},
			}),
			mf.Progress(mf.ProgressProps{
				Label:         "Sync in progress",
				Indeterminate: true,
				Props:         mf.ComponentProps{Variant: "info", Size: "sm"},
			}),
		),
	)
}

func roleMixRow(role string, count int) mf.Node {
	return mf.Actions(mf.ActionsProps{Align: "between", Props: mf.ComponentProps{Class: "rounded-box border border-base-300 px-3 py-2"}},
		roleBadge(role),
		mf.UIText(mf.TextProps{Text: strconv.Itoa(count) + " users", Size: "sm", Weight: "medium"}),
	)
}

func renderDeleteModal(ctx *mb.Context) mf.Node {
	targetID, _ := ctx.Get("deleteTargetID").(int)
	targetName := ""
	for _, u := range getUsers(ctx) {
		if u.ID == targetID {
			targetName = u.Name
			break
		}
	}

	return mf.Modal(mf.ModalProps{
		Title: "Delete user",
		Body: mf.Stack(mf.StackProps{Gap: "sm"},
			mf.UIText(mf.TextProps{Text: "Are you sure you want to delete this user?"}),
			mf.UIText(mf.TextProps{Text: targetName, Size: "sm", Tone: "subtle"}),
		),
		Actions: mf.Actions(mf.ActionsProps{Align: "end", Gap: "sm", Props: mf.ComponentProps{Class: "w-full"}},
			mf.ActionForm(mf.ActionFormProps{Action: "users/delete/cancel", Target: "#users-workspace", Swap: "outerHTML"},
				mf.SubmitButton("Cancel", mf.ComponentProps{Variant: "ghost", Size: "sm"}),
			),
			mf.ActionForm(mf.ActionFormProps{Action: "users/delete/confirm", Target: "#users-workspace", Swap: "outerHTML"},
				mf.HiddenField("id", strconv.Itoa(targetID)),
				mf.SubmitButton("Delete", mf.ComponentProps{Variant: "danger", Size: "sm"}),
			),
		),
		Open: ctx.Get("deleteModalOpen") == true,
	})
}
