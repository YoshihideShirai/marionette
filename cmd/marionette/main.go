package main

import (
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/example/marionette/internal/marionette"
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

func buildApp() *marionette.App {
	app := marionette.New()
	app.Set("nextUserID", 4)
	app.Set("users", demoUsers())
	app.Set("deleteModalOpen", false)
	app.Set("deleteTargetID", 0)
	app.Set("loading", false)

	app.Page("/", func(ctx *marionette.Context) marionette.Node {
		return renderUsersPage(ctx, defaultCreateUserFormState())
	})

	app.Action("users/create", func(ctx *marionette.Context) marionette.Node {
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

	app.Action("users/delete/prompt", func(ctx *marionette.Context) marionette.Node {
		ctx.Set("loading", false)
		id, _ := strconv.Atoi(ctx.FormValue("id"))
		ctx.Set("deleteTargetID", id)
		ctx.Set("deleteModalOpen", true)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/delete/cancel", func(ctx *marionette.Context) marionette.Node {
		ctx.Set("loading", false)
		ctx.Set("deleteModalOpen", false)
		ctx.Set("deleteTargetID", 0)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/delete/confirm", func(ctx *marionette.Context) marionette.Node {
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

	app.Action("users/loading/start", func(ctx *marionette.Context) marionette.Node {
		ctx.Set("loading", true)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/loading/stop", func(ctx *marionette.Context) marionette.Node {
		ctx.Set("loading", false)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/reset", func(ctx *marionette.Context) marionette.Node {
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

func getUsers(ctx *marionette.Context) []user {
	users, ok := ctx.Get("users").([]user)
	if !ok {
		return nil
	}
	return users
}

func renderUsersPage(ctx *marionette.Context, formState createUserFormState) marionette.Node {
	return marionette.DivProps(marionette.ElementProps{ID: "app", Class: "grid gap-6 lg:grid-cols-[16rem_minmax(0,1fr)]"},
		renderSidebar(),
		marionette.DivClass("min-w-0 space-y-6",
			marionette.FlashAlerts(ctx.Flashes()),
			marionette.DivClass("space-y-2",
				marionette.DivClass("text-3xl font-bold tracking-tight", marionette.Text("Marionette Admin UI")),
				marionette.DivClass("text-base-content/70",
					marionette.Text("Go handlers, htmx actions, and daisyUI components for small admin tools."),
				),
			),
			renderUsersWorkspace(ctx, formState),
		),
	)
}

func renderSidebar() marionette.Node {
	return marionette.Sidebar("Marionette", "Admin Console",
		marionette.SidebarLink("Users", "/").Active(),
		marionette.SidebarLink("Teams", "/"),
		marionette.SidebarLink("Settings", "/"),
	).Note("Demo workspace", "In-memory data for admin UI prototyping.")
}

func renderUsersWorkspace(ctx *marionette.Context, formState createUserFormState) marionette.Node {
	return marionette.DivProps(marionette.ElementProps{ID: "users-workspace", Class: "space-y-4"},
		marionette.FlashAlerts(ctx.Flashes()),
		renderDashboardOverview(ctx),
		marionette.DivClass("grid gap-6 lg:grid-cols-[minmax(0,1fr)_22rem]",
			renderUsersTable(ctx),
			renderCreateUserForm(formState),
		),
		renderComponentShowcase(ctx),
		renderDeleteModal(ctx),
	)
}

func renderDashboardOverview(ctx *marionette.Context) marionette.Node {
	users := getUsers(ctx)
	roles := roleCounts(users)
	latest := latestStartDate(users)

	return marionette.DivClass("grid gap-4 md:grid-cols-3",
		statCard("Users", strconv.Itoa(len(users)), "Active demo records", "primary"),
		statCard("Admins", strconv.Itoa(roles["Admin"]), "High-access seats", "secondary"),
		statCard("Latest start", latest, "Newest onboarding date", "accent"),
	)
}

func statCard(label, value, caption, tone string) marionette.Node {
	return marionette.DivClass("card bg-base-100 shadow-sm",
		marionette.DivClass("card-body gap-2",
			marionette.DivClass("text-sm font-medium text-base-content/60", marionette.Text(label)),
			marionette.DivClass("flex items-end justify-between gap-3",
				marionette.DivClass("text-3xl font-bold", marionette.Text(value)),
				marionette.DivClass("badge badge-"+tone, marionette.Text("live")),
			),
			marionette.DivClass("text-sm text-base-content/60", marionette.Text(caption)),
		),
	)
}

func renderUsersTable(ctx *marionette.Context) marionette.Node {
	users := getUsers(ctx)
	loading := isLoading(ctx)
	pg := parsePagination(ctx.Query("page"), ctx.Query("per_page"), len(users))
	tableBody := renderUsersTableBody(users, loading, ctx.Query("sort"), pg)

	return marionette.DivClass("card bg-base-100 shadow-sm",
		marionette.DivClass("card-body gap-4",
			marionette.DivClass("flex items-center justify-between gap-4",
				marionette.DivClass("space-y-1",
					marionette.DivClass("text-xl font-semibold", marionette.Text("Users")),
					marionette.DivClass("text-sm text-base-content/60", marionette.Text("Create and remove users with htmx-backed actions.")),
				),
				marionette.DivClass("flex items-center gap-2",
					marionette.DivClass("badge badge-outline", marionette.Text(strconv.Itoa(len(getUsers(ctx)))+" total")),
					marionette.Form("users/loading/start",
						marionette.ComponentSubmitButton("Show loading", marionette.ComponentProps{Variant: "ghost", Size: "sm", Disabled: loading}),
					).Target("#users-workspace"),
					marionette.Form("users/loading/stop",
						marionette.ComponentSubmitButton("Show data", marionette.ComponentProps{Variant: "ghost", Size: "sm", Disabled: !loading}),
					).Target("#users-workspace"),
					marionette.Form("users/reset",
						marionette.ComponentSubmitButton("Reset", marionette.ComponentProps{Variant: "secondary", Size: "sm"}),
					).Target("#users-workspace"),
				),
			),
			marionette.DivClass("overflow-hidden rounded-box border border-base-300", tableBody),
			marionette.ComponentPagination(marionette.PaginationProps{
				Page:       pg.Page,
				TotalPages: pg.TotalPages,
				PrevHref:   pageLink(pg.Page-1, pg.PerPage, ctx.Query("sort"), pg.TotalPages),
				NextHref:   pageLink(pg.Page+1, pg.PerPage, ctx.Query("sort"), pg.TotalPages),
			}),
		),
	)
}

func renderUsersTableBody(users []user, loading bool, sortKey string, pg pagination) marionette.Node {
	if loading {
		return marionette.ComponentEmptyState(marionette.EmptyStateProps{
			Skeleton: true,
			Rows:     5,
		})
	}

	sorted := paginateUsers(sortUsers(users, sortKey), pg)
	rows := make([]marionette.TableComponentRow, 0, len(sorted))
	for _, u := range sorted {
		rows = append(rows, renderUserRow(u))
	}

	return marionette.ComponentTable(marionette.TableProps{
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

func isLoading(ctx *marionette.Context) bool {
	v, _ := ctx.Get("loading").(bool)
	return v
}

func renderUserRow(u user) marionette.TableComponentRow {
	return marionette.TableComponentRow{
		Cells: []marionette.Node{
			marionette.DivClass("font-medium", marionette.Text(u.Name)),
			marionette.DivClass("text-sm text-base-content/70", marionette.Text(u.Email)),
			marionette.DivClass("badge badge-ghost", marionette.Text(u.Role)),
			marionette.DivClass("text-sm", marionette.Text(u.StartDate)),
			marionette.Form("users/delete/prompt",
				marionette.HiddenInput("id", strconv.Itoa(u.ID)),
				marionette.ComponentSubmitButton("Delete", marionette.ComponentProps{Variant: "danger", Size: "sm"}),
			).Target("#users-workspace"),
		},
	}
}

func usersTableColumns(activeSort string, pg pagination) []marionette.TableColumn {
	sortedQuery := func(key string) string {
		query := url.Values{}
		query.Set("sort", key)
		query.Set("page", strconv.Itoa(pg.Page))
		query.Set("per_page", strconv.Itoa(pg.PerPage))
		return "/?" + query.Encode()
	}
	return []marionette.TableColumn{
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

func roleBadge(role string) marionette.Node {
	tone := "badge-ghost"
	switch role {
	case "Admin":
		tone = "badge-primary"
	case "Editor":
		tone = "badge-secondary"
	}
	return marionette.DivClass("badge "+tone, marionette.Text(role))
}

func renderCreateUserForm(form createUserFormState) marionette.Node {
	return marionette.DivClass("card bg-base-100 shadow-sm",
		marionette.DivClass("card-body",
			marionette.DivClass("text-xl font-semibold", marionette.Text("Create user")),
			marionette.Form("users/create",
				marionette.FormRow(marionette.FormRowProps{
					ID:          "name",
					Label:       "Name",
					Description: "Enter the display name.",
					Error:       form.Errors["name"],
					Required:    true,
					Control: marionette.TextField(marionette.TextFieldProps{
						ID:          "name",
						Name:        "name",
						Value:       form.Name,
						Placeholder: "name",
						Description: "Enter the display name.",
						Error:       form.Errors["name"],
						Required:    true,
					}),
				}),
				marionette.FormRow(marionette.FormRowProps{
					ID:          "email",
					Label:       "Email",
					Description: "Used for notifications.",
					Error:       form.Errors["email"],
					Required:    true,
					Control: marionette.TextField(marionette.TextFieldProps{
						ID:          "email",
						Name:        "email",
						Value:       form.Email,
						Placeholder: "email",
						Description: "Used for notifications.",
						Error:       form.Errors["email"],
						Required:    true,
					}),
				}),
				marionette.FormRow(marionette.FormRowProps{
					ID:          "start_date",
					Label:       "Start date",
					Description: "Select a date in the active fiscal window.",
					Error:       form.Errors["start_date"],
					Required:    true,
					Control: marionette.TextField(marionette.TextFieldProps{
						ID:          "start_date",
						Name:        "start_date",
						Value:       form.StartDate,
						Type:        "date",
						Description: "Select a date in the active fiscal window.",
						Error:       form.Errors["start_date"],
						Required:    true,
					}),
				}),
				marionette.FormRow(marionette.FormRowProps{
					ID:          "role",
					Label:       "Role",
					Description: "Choose permission scope for this user.",
					Error:       form.Errors["role"],
					Required:    true,
					Control: marionette.Select(marionette.SelectFieldProps{
						ID:   "role",
						Name: "role",
						Options: []marionette.SelectOption{
							{Label: "Admin", Value: "Admin", Selected: form.Role == "Admin"},
							{Label: "Editor", Value: "Editor", Selected: form.Role == "Editor"},
							{Label: "Viewer", Value: "Viewer", Selected: form.Role == "" || form.Role == "Viewer"},
						},
						Description: "Choose permission scope for this user.",
						Error:       form.Errors["role"],
						Required:    true,
					}),
				}),
				marionette.DivClass("divider my-1"),
				marionette.FormRow(marionette.FormRowProps{
					ID:          "workspace",
					Label:       "Workspace",
					Description: "Demo-only field showing radio group markup.",
					Control: marionette.RadioGroup(marionette.RadioGroupProps{
						ID:          "workspace",
						Name:        "workspace",
						Value:       "core",
						Description: "Demo-only field showing radio group markup.",
						Options: []marionette.RadioOption{
							{Label: "Core", Value: "core"},
							{Label: "Growth", Value: "growth"},
							{Label: "Support", Value: "support"},
						},
					}),
				}),
				marionette.FormRow(marionette.FormRowProps{
					ID:          "notes",
					Label:       "Notes",
					Description: "Ignored by the demo action, useful for layout coverage.",
					Control: marionette.Textarea(marionette.TextareaProps{
						ID:          "notes",
						Name:        "notes",
						Placeholder: "Add onboarding context",
						Rows:        3,
						Description: "Ignored by the demo action, useful for layout coverage.",
					}),
				}),
				marionette.DivClass("space-y-3",
					marionette.Checkbox(marionette.CheckboxProps{
						ID:          "send_invite",
						Name:        "send_invite",
						Value:       "yes",
						Checked:     true,
						Label:       "Send invite email",
						Description: "Invite preference.",
					}),
					marionette.Switch(marionette.SwitchProps{
						ID:          "provision_access",
						Name:        "provision_access",
						Value:       "yes",
						Checked:     true,
						Label:       "Provision default access",
						Description: "Access preference.",
					}),
				),
				marionette.DivClass("flex flex-wrap gap-2 pt-2",
					marionette.ComponentSubmitButton("Create", marionette.ComponentProps{Variant: "primary", Size: "sm"}),
					marionette.ComponentButton("Preview", marionette.ComponentProps{Variant: "ghost", Size: "sm", Disabled: true}),
				),
			).Target("#users-workspace"),
		),
	)
}

func renderComponentShowcase(ctx *marionette.Context) marionette.Node {
	users := getUsers(ctx)
	roles := roleCounts(users)

	return marionette.DivClass("grid gap-6 xl:grid-cols-[minmax(0,1fr)_22rem]",
		marionette.DivClass("card bg-base-100 shadow-sm",
			marionette.DivClass("card-body gap-4",
				marionette.DivClass("space-y-1",
					marionette.DivClass("text-xl font-semibold", marionette.Text("Component states")),
					marionette.DivClass("text-sm text-base-content/60", marionette.Text("Alerts, toasts, skeletons, and empty states rendered from Go.")),
				),
				marionette.DivClass("grid gap-3 md:grid-cols-2",
					marionette.ComponentAlert(marionette.AlertProps{
						Title:       "Validation feedback",
						Description: "Form errors, success flashes, and informational alerts share the same feedback API.",
						Icon:        "!",
						Props:       marionette.ComponentProps{Variant: "info", Size: "sm"},
					}),
					marionette.ComponentToast(marionette.ToastProps{
						Title:       "Saved",
						Description: "Toast markup is available for transient status messages.",
						Icon:        "OK",
						Props:       marionette.ComponentProps{Variant: "success", Size: "sm"},
					}),
				),
				marionette.DivClass("grid gap-3 md:grid-cols-2",
					marionette.ComponentSkeleton(marionette.SkeletonProps{Rows: 4, Props: marionette.ComponentProps{Size: "sm"}}),
					marionette.ComponentEmptyState(marionette.EmptyStateProps{
						Title:       "No pending reviews",
						Description: "Empty states can replace tables or panels without extra branching in templates.",
						Icon:        "0",
						Props:       marionette.ComponentProps{Size: "sm"},
					}),
				),
			),
		),
		marionette.DivClass("card bg-base-100 shadow-sm",
			marionette.DivClass("card-body gap-4",
				marionette.DivClass("space-y-1",
					marionette.DivClass("text-xl font-semibold", marionette.Text("Role mix")),
					marionette.DivClass("text-sm text-base-content/60", marionette.Text("Small repeated views stay plain Go functions.")),
				),
				roleMixRow("Admin", roles["Admin"]),
				roleMixRow("Editor", roles["Editor"]),
				roleMixRow("Viewer", roles["Viewer"]),
			),
		),
	)
}

func roleMixRow(role string, count int) marionette.Node {
	return marionette.DivClass("flex items-center justify-between rounded-box border border-base-300 px-3 py-2",
		roleBadge(role),
		marionette.DivClass("text-sm font-medium", marionette.Text(strconv.Itoa(count)+" users")),
	)
}

func renderDeleteModal(ctx *marionette.Context) marionette.Node {
	targetID, _ := ctx.Get("deleteTargetID").(int)
	targetName := ""
	for _, u := range getUsers(ctx) {
		if u.ID == targetID {
			targetName = u.Name
			break
		}
	}

	return marionette.ComponentModal(marionette.ModalProps{
		Title: "Delete user",
		Body: marionette.DivClass("space-y-2",
			marionette.Text("Are you sure you want to delete this user?"),
			marionette.DivClass("text-sm text-base-content/70", marionette.Text(targetName)),
		),
		Actions: marionette.DivClass("flex w-full justify-end gap-2",
			marionette.Form("users/delete/cancel",
				marionette.ComponentSubmitButton("Cancel", marionette.ComponentProps{Variant: "ghost", Size: "sm"}),
			).Target("#users-workspace"),
			marionette.Form("users/delete/confirm",
				marionette.HiddenInput("id", strconv.Itoa(targetID)),
				marionette.ComponentSubmitButton("Delete", marionette.ComponentProps{Variant: "danger", Size: "sm"}),
			).Target("#users-workspace"),
		),
		Open: ctx.Get("deleteModalOpen") == true,
	})
}
