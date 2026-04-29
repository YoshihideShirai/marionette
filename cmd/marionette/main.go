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
	header := mf.PageHeaderComponent(mf.PageHeaderProps{
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
	return mf.DivProps(mf.ElementProps{ID: "app", Class: "grid gap-6 lg:grid-cols-[16rem_minmax(0,1fr)]"},
		renderSidebar(activePath),
		mf.DivClass("min-w-0 space-y-6",
			mf.FlashAlerts(ctx.Flashes()),
			header,
			mf.DivProps(mf.ElementProps{ID: "main-content", Class: "space-y-6"}, content),
		),
	)
}
func themeToggleButton() mf.Node {
	return mf.ThemeToggleButtonComponent(mf.ComponentProps{Variant: "outline", Size: "sm"})
}

func renderDashboardPage(ctx *mb.Context) mf.Node {
	users := getUsers(ctx)
	header := mf.PageHeaderComponent(mf.PageHeaderProps{
		Title:       "Operations Dashboard",
		Description: "Top-level overview of operational health and activity.",
		Actions:     themeToggleButton(),
	})
	content := mf.StackComponent(mf.StackProps{Gap: "lg"},
		renderDashboardOverview(ctx),
		renderDashboardCharts(ctx),
		mf.GridComponent(mf.GridProps{MinColumnWidth: "lg", Gap: "lg"},
			mf.CardComponent(mf.CardProps{}, mf.EmptyStateComponent(mf.EmptyStateProps{
				Title:       "Pending approvals",
				Description: "No approvals are currently waiting.",
				Icon:        "✓",
			})),
			mf.CardComponent(mf.CardProps{}, mf.TableComponent(mf.TableProps{
				Columns: []mf.TableColumn{{Label: "Metric"}, {Label: "Value"}},
				Rows: []mf.TableComponentRow{
					{Cells: []mf.Node{mf.Text("Total users"), mf.Text(strconv.Itoa(len(users)))}},
					{Cells: []mf.Node{mf.Text("Latest start"), mf.Text(latestStartDate(users))}},
				},
			})),
		),
	)
	return renderShell(ctx, "/", header, content)
}

func renderAnalyticsPage(ctx *mb.Context) mf.Node {
	header := mf.PageHeaderComponent(mf.PageHeaderProps{Title: "Analytics", Description: "A chart-first, two-column analytics layout.", Actions: themeToggleButton()})
	content := mf.SplitComponent(mf.SplitProps{
		Gap: "lg",
		Main: mf.StackComponent(mf.StackProps{Gap: "lg"},
			renderDashboardCharts(ctx),
			mf.AlertComponent(mf.AlertProps{Title: "Insight", Description: "Editors are growing steadily month-over-month.", Props: mf.ComponentProps{Variant: "info"}}),
		),
		Aside: mf.StackComponent(mf.StackProps{Gap: "md"},
			roleMixRow("Admin", roleCounts(getUsers(ctx))["Admin"]),
			roleMixRow("Editor", roleCounts(getUsers(ctx))["Editor"]),
			roleMixRow("Viewer", roleCounts(getUsers(ctx))["Viewer"]),
		),
	})
	return renderShell(ctx, "/analytics", header, content)
}

func renderSettingsPage(ctx *mb.Context) mf.Node {
	header := mf.PageHeaderComponent(mf.PageHeaderProps{Title: "Settings", Description: "A form-heavy vertical layout for configuration flows.", Actions: themeToggleButton()})
	content := mf.StackComponent(mf.StackProps{Gap: "lg"},
		mf.CardComponent(mf.CardProps{},
			mf.Form("users/reset",
				mf.FormRow(mf.FormRowProps{ID: "org-name", Label: "Organization name", Control: mf.TextField(mf.TextFieldProps{ID: "org-name", Name: "org-name", Value: "Marionette Labs"})}),
				mf.FormRow(mf.FormRowProps{ID: "timezone", Label: "Timezone", Control: mf.Select(mf.SelectFieldProps{ID: "timezone", Name: "timezone", Options: []mf.SelectOption{{Label: "UTC", Value: "UTC", Selected: true}, {Label: "Asia/Tokyo", Value: "Asia/Tokyo"}}})}),
				mf.Switch(mf.SwitchProps{ID: "audit", Name: "audit", Value: "yes", Checked: true, Label: "Enable audit log"}),
				mf.DivClass("pt-3", mf.SubmitButtonComponent("Save settings", mf.ComponentProps{Variant: "primary"})),
			),
		),
		mf.ToastComponent(mf.ToastProps{Title: "Tip", Description: "This page demonstrates form-heavy layout variation.", Props: mf.ComponentProps{Variant: "info"}}),
	)
	return renderShell(ctx, "/settings", header, content)
}

func renderUsersWorkspace(ctx *mb.Context, formState createUserFormState) mf.Node {
	return mf.DivProps(mf.ElementProps{ID: "users-workspace", Class: "space-y-4"},
		mf.FlashAlerts(ctx.Flashes()),
		renderDashboardOverview(ctx),
		renderDashboardCharts(ctx),
		mf.SplitComponent(mf.SplitProps{
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

	return mf.GridComponent(mf.GridProps{MinColumnWidth: "md", Gap: "lg"},
		mf.ChartComponent(mf.ChartProps{
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
		mf.ChartComponent(mf.ChartProps{
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

	return mf.GridComponent(mf.GridProps{Columns: "3"},
		statCard("Users", strconv.Itoa(len(users)), "Active demo records", "primary"),
		statCard("Admins", strconv.Itoa(roles["Admin"]), "High-access seats", "secondary"),
		statCard("Latest start", latest, "Newest onboarding date", "accent"),
	)
}

func statCard(label, value, caption, tone string) mf.Node {
	return mf.CardComponent(mf.CardProps{},
		mf.StackComponent(mf.StackProps{Gap: "sm"},
			mf.DivClass("text-sm font-medium text-base-content/60", mf.Text(label)),
			mf.StackComponent(mf.StackProps{Direction: "horizontal", Gap: "md", Align: "end", Justify: "between"},
				mf.DivClass("text-3xl font-bold", mf.Text(value)),
				mf.DivClass("badge badge-"+tone, mf.Text("live")),
			),
			mf.DivClass("text-sm text-base-content/60", mf.Text(caption)),
		),
	)
}

func renderUsersTable(ctx *mb.Context) mf.Node {
	users := getUsers(ctx)
	loading := isLoading(ctx)
	pg := parsePagination(ctx.Query("page"), ctx.Query("per_page"), len(users))
	tableBody := renderUsersTableBody(users, loading, ctx.Query("sort"), pg)

	return mf.DivClass("card bg-base-100 shadow-sm",
		mf.DivClass("card-body gap-4",
			mf.DivClass("flex items-center justify-between gap-4",
				mf.DivClass("space-y-1",
					mf.DivClass("text-xl font-semibold", mf.Text("Users")),
					mf.DivClass("text-sm text-base-content/60", mf.Text("Create and remove users with htmx-backed actions.")),
				),
				mf.DivClass("flex items-center gap-2",
					mf.DivClass("badge badge-outline", mf.Text(strconv.Itoa(len(getUsers(ctx)))+" total")),
					mf.Form("users/loading/start",
						mf.SubmitButtonComponent("Show loading", mf.ComponentProps{Variant: "ghost", Size: "sm", Disabled: loading}),
					).Target("#users-workspace"),
					mf.Form("users/loading/stop",
						mf.SubmitButtonComponent("Show data", mf.ComponentProps{Variant: "ghost", Size: "sm", Disabled: !loading}),
					).Target("#users-workspace"),
					mf.Form("users/reset",
						mf.SubmitButtonComponent("Reset", mf.ComponentProps{Variant: "secondary", Size: "sm"}),
					).Target("#users-workspace"),
				),
			),
			mf.DivClass("overflow-hidden rounded-box border border-base-300", tableBody),
			mf.PaginationComponent(mf.PaginationProps{
				Page:       pg.Page,
				TotalPages: pg.TotalPages,
				PrevHref:   pageLink(pg.Page-1, pg.PerPage, ctx.Query("sort"), pg.TotalPages),
				NextHref:   pageLink(pg.Page+1, pg.PerPage, ctx.Query("sort"), pg.TotalPages),
			}),
		),
	)
}

func renderUsersTableBody(users []user, loading bool, sortKey string, pg pagination) mf.Node {
	if loading {
		return mf.EmptyStateComponent(mf.EmptyStateProps{
			Skeleton: true,
			Rows:     5,
		})
	}

	sorted := paginateUsers(sortUsers(users, sortKey), pg)
	rows := make([]mf.TableComponentRow, 0, len(sorted))
	for _, u := range sorted {
		rows = append(rows, renderUserRow(u))
	}

	return mf.TableComponent(mf.TableProps{
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
	return mf.TableComponentRow{
		Cells: []mf.Node{
			mf.DivClass("font-medium", mf.Text(u.Name)),
			mf.DivClass("text-sm text-base-content/70", mf.Text(u.Email)),
			mf.DivClass("badge badge-ghost", mf.Text(u.Role)),
			mf.DivClass("text-sm", mf.Text(u.StartDate)),
			mf.Form("users/delete/prompt",
				mf.HiddenInput("id", strconv.Itoa(u.ID)),
				mf.SubmitButtonComponent("Delete", mf.ComponentProps{Variant: "danger", Size: "sm"}),
			).Target("#users-workspace"),
		},
	}
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
	tone := "badge-ghost"
	switch role {
	case "Admin":
		tone = "badge-primary"
	case "Editor":
		tone = "badge-secondary"
	}
	return mf.DivClass("badge "+tone, mf.Text(role))
}

func renderCreateUserForm(form createUserFormState) mf.Node {
	return mf.DivClass("card bg-base-100 shadow-sm",
		mf.DivClass("card-body",
			mf.DivClass("text-xl font-semibold", mf.Text("Create user")),
			mf.Form("users/create",
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
				mf.DivClass("divider my-1"),
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
				mf.DivClass("space-y-3",
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
				mf.DivClass("flex flex-wrap gap-2 pt-2",
					mf.SubmitButtonComponent("Create", mf.ComponentProps{Variant: "primary", Size: "sm"}),
					mf.ButtonComponent("Preview", mf.ComponentProps{Variant: "ghost", Size: "sm", Disabled: true}),
				),
			).Target("#users-workspace"),
		),
	)
}

func renderComponentShowcase(ctx *mb.Context) mf.Node {
	users := getUsers(ctx)
	roles := roleCounts(users)

	return mf.DivClass("grid gap-6 xl:grid-cols-[minmax(0,1fr)_22rem]",
		mf.DivClass("card bg-base-100 shadow-sm",
			mf.DivClass("card-body gap-4",
				mf.DivClass("space-y-1",
					mf.DivClass("text-xl font-semibold", mf.Text("Component states")),
					mf.DivClass("text-sm text-base-content/60", mf.Text("Alerts, toasts, skeletons, and empty states rendered from Go.")),
				),
				mf.DivClass("grid gap-3 md:grid-cols-2",
					mf.AlertComponent(mf.AlertProps{
						Title:       "Validation feedback",
						Description: "Form errors, success flashes, and informational alerts share the same feedback API.",
						Icon:        "!",
						Props:       mf.ComponentProps{Variant: "info", Size: "sm"},
					}),
					mf.ToastComponent(mf.ToastProps{
						Title:       "Saved",
						Description: "Toast markup is available for transient status messages.",
						Icon:        "OK",
						Props:       mf.ComponentProps{Variant: "success", Size: "sm"},
					}),
				),
				mf.DivClass("grid gap-3 md:grid-cols-2",
					mf.SkeletonComponent(mf.SkeletonProps{Rows: 4, Props: mf.ComponentProps{Size: "sm"}}),
					mf.EmptyStateComponent(mf.EmptyStateProps{
						Title:       "No pending reviews",
						Description: "Empty states can replace tables or panels without extra branching in templates.",
						Icon:        "0",
						Props:       mf.ComponentProps{Size: "sm"},
					}),
				),
				mf.DivClass("rounded-box border border-base-300 bg-base-100 p-3",
					mf.DivClass("mb-2 text-sm text-base-content/60", mf.Text("Theme toggle component")),
					mf.ThemeToggleButtonComponent(mf.ComponentProps{Variant: "outline", Size: "sm"}),
				),
			),
		),
		mf.DivClass("card bg-base-100 shadow-sm",
			mf.DivClass("card-body gap-4",
				mf.DivClass("space-y-1",
					mf.DivClass("text-xl font-semibold", mf.Text("Role mix")),
					mf.DivClass("text-sm text-base-content/60", mf.Text("Small repeated views stay plain Go functions.")),
				),
				roleMixRow("Admin", roles["Admin"]),
				roleMixRow("Editor", roles["Editor"]),
				roleMixRow("Viewer", roles["Viewer"]),
			),
		),
	)
}

func roleMixRow(role string, count int) mf.Node {
	return mf.DivClass("flex items-center justify-between rounded-box border border-base-300 px-3 py-2",
		roleBadge(role),
		mf.DivClass("text-sm font-medium", mf.Text(strconv.Itoa(count)+" users")),
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

	return mf.ModalComponent(mf.ModalProps{
		Title: "Delete user",
		Body: mf.DivClass("space-y-2",
			mf.Text("Are you sure you want to delete this user?"),
			mf.DivClass("text-sm text-base-content/70", mf.Text(targetName)),
		),
		Actions: mf.DivClass("flex w-full justify-end gap-2",
			mf.Form("users/delete/cancel",
				mf.SubmitButtonComponent("Cancel", mf.ComponentProps{Variant: "ghost", Size: "sm"}),
			).Target("#users-workspace"),
			mf.Form("users/delete/confirm",
				mf.HiddenInput("id", strconv.Itoa(targetID)),
				mf.SubmitButtonComponent("Delete", mf.ComponentProps{Variant: "danger", Size: "sm"}),
			).Target("#users-workspace"),
		),
		Open: ctx.Get("deleteModalOpen") == true,
	})
}
