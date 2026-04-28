package main

import (
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	backend "github.com/YoshihideShirai/marionette/backend"
	frontend "github.com/YoshihideShirai/marionette/frontend"
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

func buildApp() *backend.App {
	app := backend.New()
	app.Set("nextUserID", 4)
	app.Set("users", demoUsers())
	app.Set("deleteModalOpen", false)
	app.Set("deleteTargetID", 0)
	app.Set("loading", false)

	app.Page("/", func(ctx *backend.Context) frontend.Node {
		return renderUsersPage(ctx, defaultCreateUserFormState())
	})

	app.Action("users/create", func(ctx *backend.Context) frontend.Node {
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

	app.Action("users/delete/prompt", func(ctx *backend.Context) frontend.Node {
		ctx.Set("loading", false)
		id, _ := strconv.Atoi(ctx.FormValue("id"))
		ctx.Set("deleteTargetID", id)
		ctx.Set("deleteModalOpen", true)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/delete/cancel", func(ctx *backend.Context) frontend.Node {
		ctx.Set("loading", false)
		ctx.Set("deleteModalOpen", false)
		ctx.Set("deleteTargetID", 0)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/delete/confirm", func(ctx *backend.Context) frontend.Node {
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

	app.Action("users/loading/start", func(ctx *backend.Context) frontend.Node {
		ctx.Set("loading", true)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/loading/stop", func(ctx *backend.Context) frontend.Node {
		ctx.Set("loading", false)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/reset", func(ctx *backend.Context) frontend.Node {
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

func getUsers(ctx *backend.Context) []user {
	users, ok := ctx.Get("users").([]user)
	if !ok {
		return nil
	}
	return users
}

func renderUsersPage(ctx *backend.Context, formState createUserFormState) frontend.Node {
	return frontend.DivProps(frontend.ElementProps{ID: "app", Class: "grid gap-6 lg:grid-cols-[16rem_minmax(0,1fr)]"},
		renderSidebar(),
		frontend.DivClass("min-w-0 space-y-6",
			frontend.FlashAlerts(ctx.Flashes()),
			frontend.DivClass("flex items-start justify-between gap-4",
				frontend.DivClass("space-y-2",
					frontend.DivClass("text-3xl font-bold tracking-tight", frontend.Text("Marionette Admin UI")),
					frontend.DivClass("text-base-content/70",
						frontend.Text("Go handlers, htmx actions, and daisyUI components for small admin tools."),
					),
				),
				frontend.Element("button", frontend.ElementProps{
					Class: "btn btn-outline btn-sm",
					Attrs: frontend.Attrs{
						"type":    "button",
						"onclick": "window.mrnToggleTheme && window.mrnToggleTheme()",
					},
				}, frontend.Text("🌓 Theme")),
			),
			renderUsersWorkspace(ctx, formState),
		),
	)
}

func renderSidebar() frontend.Node {
	return frontend.Sidebar("Marionette", "Admin Console",
		frontend.SidebarLink("Users", "/").Active(),
		frontend.SidebarLink("Teams", "/"),
		frontend.SidebarLink("Settings", "/"),
	).Note("Demo workspace", "In-memory data for admin UI prototyping.")
}

func renderUsersWorkspace(ctx *backend.Context, formState createUserFormState) frontend.Node {
	return frontend.DivProps(frontend.ElementProps{ID: "users-workspace", Class: "space-y-4"},
		frontend.FlashAlerts(ctx.Flashes()),
		renderDashboardOverview(ctx),
		frontend.DivClass("grid gap-6 lg:grid-cols-[minmax(0,1fr)_22rem]",
			renderUsersTable(ctx),
			renderCreateUserForm(formState),
		),
		renderComponentShowcase(ctx),
		renderDeleteModal(ctx),
	)
}

func renderDashboardOverview(ctx *backend.Context) frontend.Node {
	users := getUsers(ctx)
	roles := roleCounts(users)
	latest := latestStartDate(users)

	return frontend.DivClass("grid gap-4 md:grid-cols-3",
		statCard("Users", strconv.Itoa(len(users)), "Active demo records", "primary"),
		statCard("Admins", strconv.Itoa(roles["Admin"]), "High-access seats", "secondary"),
		statCard("Latest start", latest, "Newest onboarding date", "accent"),
	)
}

func statCard(label, value, caption, tone string) frontend.Node {
	return frontend.DivClass("card bg-base-100 shadow-sm",
		frontend.DivClass("card-body gap-2",
			frontend.DivClass("text-sm font-medium text-base-content/60", frontend.Text(label)),
			frontend.DivClass("flex items-end justify-between gap-3",
				frontend.DivClass("text-3xl font-bold", frontend.Text(value)),
				frontend.DivClass("badge badge-"+tone, frontend.Text("live")),
			),
			frontend.DivClass("text-sm text-base-content/60", frontend.Text(caption)),
		),
	)
}

func renderUsersTable(ctx *backend.Context) frontend.Node {
	users := getUsers(ctx)
	loading := isLoading(ctx)
	pg := parsePagination(ctx.Query("page"), ctx.Query("per_page"), len(users))
	tableBody := renderUsersTableBody(users, loading, ctx.Query("sort"), pg)

	return frontend.DivClass("card bg-base-100 shadow-sm",
		frontend.DivClass("card-body gap-4",
			frontend.DivClass("flex items-center justify-between gap-4",
				frontend.DivClass("space-y-1",
					frontend.DivClass("text-xl font-semibold", frontend.Text("Users")),
					frontend.DivClass("text-sm text-base-content/60", frontend.Text("Create and remove users with htmx-backed actions.")),
				),
				frontend.DivClass("flex items-center gap-2",
					frontend.DivClass("badge badge-outline", frontend.Text(strconv.Itoa(len(getUsers(ctx)))+" total")),
					frontend.Form("users/loading/start",
						frontend.ComponentSubmitButton("Show loading", frontend.ComponentProps{Variant: "ghost", Size: "sm", Disabled: loading}),
					).Target("#users-workspace"),
					frontend.Form("users/loading/stop",
						frontend.ComponentSubmitButton("Show data", frontend.ComponentProps{Variant: "ghost", Size: "sm", Disabled: !loading}),
					).Target("#users-workspace"),
					frontend.Form("users/reset",
						frontend.ComponentSubmitButton("Reset", frontend.ComponentProps{Variant: "secondary", Size: "sm"}),
					).Target("#users-workspace"),
				),
			),
			frontend.DivClass("overflow-hidden rounded-box border border-base-300", tableBody),
			frontend.ComponentPagination(frontend.PaginationProps{
				Page:       pg.Page,
				TotalPages: pg.TotalPages,
				PrevHref:   pageLink(pg.Page-1, pg.PerPage, ctx.Query("sort"), pg.TotalPages),
				NextHref:   pageLink(pg.Page+1, pg.PerPage, ctx.Query("sort"), pg.TotalPages),
			}),
		),
	)
}

func renderUsersTableBody(users []user, loading bool, sortKey string, pg pagination) frontend.Node {
	if loading {
		return frontend.ComponentEmptyState(frontend.EmptyStateProps{
			Skeleton: true,
			Rows:     5,
		})
	}

	sorted := paginateUsers(sortUsers(users, sortKey), pg)
	rows := make([]frontend.TableComponentRow, 0, len(sorted))
	for _, u := range sorted {
		rows = append(rows, renderUserRow(u))
	}

	return frontend.ComponentTable(frontend.TableProps{
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
	return "/?" + query.Encode()
}

func isLoading(ctx *backend.Context) bool {
	v, _ := ctx.Get("loading").(bool)
	return v
}

func renderUserRow(u user) frontend.TableComponentRow {
	return frontend.TableComponentRow{
		Cells: []frontend.Node{
			frontend.DivClass("font-medium", frontend.Text(u.Name)),
			frontend.DivClass("text-sm text-base-content/70", frontend.Text(u.Email)),
			frontend.DivClass("badge badge-ghost", frontend.Text(u.Role)),
			frontend.DivClass("text-sm", frontend.Text(u.StartDate)),
			frontend.Form("users/delete/prompt",
				frontend.HiddenInput("id", strconv.Itoa(u.ID)),
				frontend.ComponentSubmitButton("Delete", frontend.ComponentProps{Variant: "danger", Size: "sm"}),
			).Target("#users-workspace"),
		},
	}
}

func usersTableColumns(activeSort string, pg pagination) []frontend.TableColumn {
	sortedQuery := func(key string) string {
		query := url.Values{}
		query.Set("sort", key)
		query.Set("page", strconv.Itoa(pg.Page))
		query.Set("per_page", strconv.Itoa(pg.PerPage))
		return "/?" + query.Encode()
	}
	return []frontend.TableColumn{
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

func latestStartDate(users []user) string {
	latest := "n/a"
	for _, u := range users {
		if u.StartDate > latest || latest == "n/a" {
			latest = u.StartDate
		}
	}
	return latest
}

func roleBadge(role string) frontend.Node {
	tone := "badge-ghost"
	switch role {
	case "Admin":
		tone = "badge-primary"
	case "Editor":
		tone = "badge-secondary"
	}
	return frontend.DivClass("badge "+tone, frontend.Text(role))
}

func renderCreateUserForm(form createUserFormState) frontend.Node {
	return frontend.DivClass("card bg-base-100 shadow-sm",
		frontend.DivClass("card-body",
			frontend.DivClass("text-xl font-semibold", frontend.Text("Create user")),
			frontend.Form("users/create",
				frontend.FormRow(frontend.FormRowProps{
					ID:          "name",
					Label:       "Name",
					Description: "Enter the display name.",
					Error:       form.Errors["name"],
					Required:    true,
					Control: frontend.TextField(frontend.TextFieldProps{
						ID:          "name",
						Name:        "name",
						Value:       form.Name,
						Placeholder: "name",
						Description: "Enter the display name.",
						Error:       form.Errors["name"],
						Required:    true,
					}),
				}),
				frontend.FormRow(frontend.FormRowProps{
					ID:          "email",
					Label:       "Email",
					Description: "Used for notifications.",
					Error:       form.Errors["email"],
					Required:    true,
					Control: frontend.TextField(frontend.TextFieldProps{
						ID:          "email",
						Name:        "email",
						Value:       form.Email,
						Placeholder: "email",
						Description: "Used for notifications.",
						Error:       form.Errors["email"],
						Required:    true,
					}),
				}),
				frontend.FormRow(frontend.FormRowProps{
					ID:          "start_date",
					Label:       "Start date",
					Description: "Select a date in the active fiscal window.",
					Error:       form.Errors["start_date"],
					Required:    true,
					Control: frontend.TextField(frontend.TextFieldProps{
						ID:          "start_date",
						Name:        "start_date",
						Value:       form.StartDate,
						Type:        "date",
						Description: "Select a date in the active fiscal window.",
						Error:       form.Errors["start_date"],
						Required:    true,
					}),
				}),
				frontend.FormRow(frontend.FormRowProps{
					ID:          "role",
					Label:       "Role",
					Description: "Choose permission scope for this user.",
					Error:       form.Errors["role"],
					Required:    true,
					Control: frontend.Select(frontend.SelectFieldProps{
						ID:   "role",
						Name: "role",
						Options: []frontend.SelectOption{
							{Label: "Admin", Value: "Admin", Selected: form.Role == "Admin"},
							{Label: "Editor", Value: "Editor", Selected: form.Role == "Editor"},
							{Label: "Viewer", Value: "Viewer", Selected: form.Role == "" || form.Role == "Viewer"},
						},
						Description: "Choose permission scope for this user.",
						Error:       form.Errors["role"],
						Required:    true,
					}),
				}),
				frontend.DivClass("divider my-1"),
				frontend.FormRow(frontend.FormRowProps{
					ID:          "workspace",
					Label:       "Workspace",
					Description: "Demo-only field showing radio group markup.",
					Control: frontend.RadioGroup(frontend.RadioGroupProps{
						ID:          "workspace",
						Name:        "workspace",
						Value:       "core",
						Description: "Demo-only field showing radio group markup.",
						Options: []frontend.RadioOption{
							{Label: "Core", Value: "core"},
							{Label: "Growth", Value: "growth"},
							{Label: "Support", Value: "support"},
						},
					}),
				}),
				frontend.FormRow(frontend.FormRowProps{
					ID:          "notes",
					Label:       "Notes",
					Description: "Ignored by the demo action, useful for layout coverage.",
					Control: frontend.Textarea(frontend.TextareaProps{
						ID:          "notes",
						Name:        "notes",
						Placeholder: "Add onboarding context",
						Rows:        3,
						Description: "Ignored by the demo action, useful for layout coverage.",
					}),
				}),
				frontend.DivClass("space-y-3",
					frontend.Checkbox(frontend.CheckboxProps{
						ID:          "send_invite",
						Name:        "send_invite",
						Value:       "yes",
						Checked:     true,
						Label:       "Send invite email",
						Description: "Invite preference.",
					}),
					frontend.Switch(frontend.SwitchProps{
						ID:          "provision_access",
						Name:        "provision_access",
						Value:       "yes",
						Checked:     true,
						Label:       "Provision default access",
						Description: "Access preference.",
					}),
				),
				frontend.DivClass("flex flex-wrap gap-2 pt-2",
					frontend.ComponentSubmitButton("Create", frontend.ComponentProps{Variant: "primary", Size: "sm"}),
					frontend.ComponentButton("Preview", frontend.ComponentProps{Variant: "ghost", Size: "sm", Disabled: true}),
				),
			).Target("#users-workspace"),
		),
	)
}

func renderComponentShowcase(ctx *backend.Context) frontend.Node {
	users := getUsers(ctx)
	roles := roleCounts(users)

	return frontend.DivClass("grid gap-6 xl:grid-cols-[minmax(0,1fr)_22rem]",
		frontend.DivClass("card bg-base-100 shadow-sm",
			frontend.DivClass("card-body gap-4",
				frontend.DivClass("space-y-1",
					frontend.DivClass("text-xl font-semibold", frontend.Text("Component states")),
					frontend.DivClass("text-sm text-base-content/60", frontend.Text("Alerts, toasts, skeletons, and empty states rendered from Go.")),
				),
				frontend.DivClass("grid gap-3 md:grid-cols-2",
					frontend.ComponentAlert(frontend.AlertProps{
						Title:       "Validation feedback",
						Description: "Form errors, success flashes, and informational alerts share the same feedback API.",
						Icon:        "!",
						Props:       frontend.ComponentProps{Variant: "info", Size: "sm"},
					}),
					frontend.ComponentToast(frontend.ToastProps{
						Title:       "Saved",
						Description: "Toast markup is available for transient status messages.",
						Icon:        "OK",
						Props:       frontend.ComponentProps{Variant: "success", Size: "sm"},
					}),
				),
				frontend.DivClass("grid gap-3 md:grid-cols-2",
					frontend.ComponentSkeleton(frontend.SkeletonProps{Rows: 4, Props: frontend.ComponentProps{Size: "sm"}}),
					frontend.ComponentEmptyState(frontend.EmptyStateProps{
						Title:       "No pending reviews",
						Description: "Empty states can replace tables or panels without extra branching in templates.",
						Icon:        "0",
						Props:       frontend.ComponentProps{Size: "sm"},
					}),
				),
			),
		),
		frontend.DivClass("card bg-base-100 shadow-sm",
			frontend.DivClass("card-body gap-4",
				frontend.DivClass("space-y-1",
					frontend.DivClass("text-xl font-semibold", frontend.Text("Role mix")),
					frontend.DivClass("text-sm text-base-content/60", frontend.Text("Small repeated views stay plain Go functions.")),
				),
				roleMixRow("Admin", roles["Admin"]),
				roleMixRow("Editor", roles["Editor"]),
				roleMixRow("Viewer", roles["Viewer"]),
			),
		),
	)
}

func roleMixRow(role string, count int) frontend.Node {
	return frontend.DivClass("flex items-center justify-between rounded-box border border-base-300 px-3 py-2",
		roleBadge(role),
		frontend.DivClass("text-sm font-medium", frontend.Text(strconv.Itoa(count)+" users")),
	)
}

func renderDeleteModal(ctx *backend.Context) frontend.Node {
	targetID, _ := ctx.Get("deleteTargetID").(int)
	targetName := ""
	for _, u := range getUsers(ctx) {
		if u.ID == targetID {
			targetName = u.Name
			break
		}
	}

	return frontend.ComponentModal(frontend.ModalProps{
		Title: "Delete user",
		Body: frontend.DivClass("space-y-2",
			frontend.Text("Are you sure you want to delete this user?"),
			frontend.DivClass("text-sm text-base-content/70", frontend.Text(targetName)),
		),
		Actions: frontend.DivClass("flex w-full justify-end gap-2",
			frontend.Form("users/delete/cancel",
				frontend.ComponentSubmitButton("Cancel", frontend.ComponentProps{Variant: "ghost", Size: "sm"}),
			).Target("#users-workspace"),
			frontend.Form("users/delete/confirm",
				frontend.HiddenInput("id", strconv.Itoa(targetID)),
				frontend.ComponentSubmitButton("Delete", frontend.ComponentProps{Variant: "danger", Size: "sm"}),
			).Target("#users-workspace"),
		),
		Open: ctx.Get("deleteModalOpen") == true,
	})
}
