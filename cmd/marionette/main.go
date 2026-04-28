package main

import (
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	mrn "github.com/YoshihideShirai/marionette"
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

func buildApp() *mrn.App {
	app := mrn.New()
	app.Set("nextUserID", 4)
	app.Set("users", demoUsers())
	app.Set("deleteModalOpen", false)
	app.Set("deleteTargetID", 0)
	app.Set("loading", false)

	app.Page("/", func(ctx *mrn.Context) mrn.Node {
		return renderUsersPage(ctx, defaultCreateUserFormState())
	})

	app.Action("users/create", func(ctx *mrn.Context) mrn.Node {
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

	app.Action("users/delete/prompt", func(ctx *mrn.Context) mrn.Node {
		ctx.Set("loading", false)
		id, _ := strconv.Atoi(ctx.FormValue("id"))
		ctx.Set("deleteTargetID", id)
		ctx.Set("deleteModalOpen", true)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/delete/cancel", func(ctx *mrn.Context) mrn.Node {
		ctx.Set("loading", false)
		ctx.Set("deleteModalOpen", false)
		ctx.Set("deleteTargetID", 0)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/delete/confirm", func(ctx *mrn.Context) mrn.Node {
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

	app.Action("users/loading/start", func(ctx *mrn.Context) mrn.Node {
		ctx.Set("loading", true)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/loading/stop", func(ctx *mrn.Context) mrn.Node {
		ctx.Set("loading", false)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/reset", func(ctx *mrn.Context) mrn.Node {
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

func getUsers(ctx *mrn.Context) []user {
	users, ok := ctx.Get("users").([]user)
	if !ok {
		return nil
	}
	return users
}

func renderUsersPage(ctx *mrn.Context, formState createUserFormState) mrn.Node {
	return mrn.DivProps(mrn.ElementProps{ID: "app", Class: "grid gap-6 lg:grid-cols-[16rem_minmax(0,1fr)]"},
		renderSidebar(),
		mrn.DivClass("min-w-0 space-y-6",
			mrn.FlashAlerts(ctx.Flashes()),
			mrn.DivClass("space-y-2",
				mrn.DivClass("text-3xl font-bold tracking-tight", mrn.Text("Marionette Admin UI")),
				mrn.DivClass("text-base-content/70",
					mrn.Text("Go handlers, htmx actions, and daisyUI components for small admin tools."),
				),
			),
			renderUsersWorkspace(ctx, formState),
		),
	)
}

func renderSidebar() mrn.Node {
	return mrn.Sidebar("Marionette", "Admin Console",
		mrn.SidebarLink("Users", "/").Active(),
		mrn.SidebarLink("Teams", "/"),
		mrn.SidebarLink("Settings", "/"),
	).Note("Demo workspace", "In-memory data for admin UI prototyping.")
}

func renderUsersWorkspace(ctx *mrn.Context, formState createUserFormState) mrn.Node {
	return mrn.DivProps(mrn.ElementProps{ID: "users-workspace", Class: "space-y-4"},
		mrn.FlashAlerts(ctx.Flashes()),
		renderDashboardOverview(ctx),
		mrn.DivClass("grid gap-6 lg:grid-cols-[minmax(0,1fr)_22rem]",
			renderUsersTable(ctx),
			renderCreateUserForm(formState),
		),
		renderComponentShowcase(ctx),
		renderDeleteModal(ctx),
	)
}

func renderDashboardOverview(ctx *mrn.Context) mrn.Node {
	users := getUsers(ctx)
	roles := roleCounts(users)
	latest := latestStartDate(users)

	return mrn.DivClass("grid gap-4 md:grid-cols-3",
		statCard("Users", strconv.Itoa(len(users)), "Active demo records", "primary"),
		statCard("Admins", strconv.Itoa(roles["Admin"]), "High-access seats", "secondary"),
		statCard("Latest start", latest, "Newest onboarding date", "accent"),
	)
}

func statCard(label, value, caption, tone string) mrn.Node {
	return mrn.DivClass("card bg-base-100 shadow-sm",
		mrn.DivClass("card-body gap-2",
			mrn.DivClass("text-sm font-medium text-base-content/60", mrn.Text(label)),
			mrn.DivClass("flex items-end justify-between gap-3",
				mrn.DivClass("text-3xl font-bold", mrn.Text(value)),
				mrn.DivClass("badge badge-"+tone, mrn.Text("live")),
			),
			mrn.DivClass("text-sm text-base-content/60", mrn.Text(caption)),
		),
	)
}

func renderUsersTable(ctx *mrn.Context) mrn.Node {
	users := getUsers(ctx)
	loading := isLoading(ctx)
	pg := parsePagination(ctx.Query("page"), ctx.Query("per_page"), len(users))
	tableBody := renderUsersTableBody(users, loading, ctx.Query("sort"), pg)

	return mrn.DivClass("card bg-base-100 shadow-sm",
		mrn.DivClass("card-body gap-4",
			mrn.DivClass("flex items-center justify-between gap-4",
				mrn.DivClass("space-y-1",
					mrn.DivClass("text-xl font-semibold", mrn.Text("Users")),
					mrn.DivClass("text-sm text-base-content/60", mrn.Text("Create and remove users with htmx-backed actions.")),
				),
				mrn.DivClass("flex items-center gap-2",
					mrn.DivClass("badge badge-outline", mrn.Text(strconv.Itoa(len(getUsers(ctx)))+" total")),
					mrn.Form("users/loading/start",
						mrn.ComponentSubmitButton("Show loading", mrn.ComponentProps{Variant: "ghost", Size: "sm", Disabled: loading}),
					).Target("#users-workspace"),
					mrn.Form("users/loading/stop",
						mrn.ComponentSubmitButton("Show data", mrn.ComponentProps{Variant: "ghost", Size: "sm", Disabled: !loading}),
					).Target("#users-workspace"),
					mrn.Form("users/reset",
						mrn.ComponentSubmitButton("Reset", mrn.ComponentProps{Variant: "secondary", Size: "sm"}),
					).Target("#users-workspace"),
				),
			),
			mrn.DivClass("overflow-hidden rounded-box border border-base-300", tableBody),
			mrn.ComponentPagination(mrn.PaginationProps{
				Page:       pg.Page,
				TotalPages: pg.TotalPages,
				PrevHref:   pageLink(pg.Page-1, pg.PerPage, ctx.Query("sort"), pg.TotalPages),
				NextHref:   pageLink(pg.Page+1, pg.PerPage, ctx.Query("sort"), pg.TotalPages),
			}),
		),
	)
}

func renderUsersTableBody(users []user, loading bool, sortKey string, pg pagination) mrn.Node {
	if loading {
		return mrn.ComponentEmptyState(mrn.EmptyStateProps{
			Skeleton: true,
			Rows:     5,
		})
	}

	sorted := paginateUsers(sortUsers(users, sortKey), pg)
	rows := make([]mrn.TableComponentRow, 0, len(sorted))
	for _, u := range sorted {
		rows = append(rows, renderUserRow(u))
	}

	return mrn.ComponentTable(mrn.TableProps{
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

func isLoading(ctx *mrn.Context) bool {
	v, _ := ctx.Get("loading").(bool)
	return v
}

func renderUserRow(u user) mrn.TableComponentRow {
	return mrn.TableComponentRow{
		Cells: []mrn.Node{
			mrn.DivClass("font-medium", mrn.Text(u.Name)),
			mrn.DivClass("text-sm text-base-content/70", mrn.Text(u.Email)),
			mrn.DivClass("badge badge-ghost", mrn.Text(u.Role)),
			mrn.DivClass("text-sm", mrn.Text(u.StartDate)),
			mrn.Form("users/delete/prompt",
				mrn.HiddenInput("id", strconv.Itoa(u.ID)),
				mrn.ComponentSubmitButton("Delete", mrn.ComponentProps{Variant: "danger", Size: "sm"}),
			).Target("#users-workspace"),
		},
	}
}

func usersTableColumns(activeSort string, pg pagination) []mrn.TableColumn {
	sortedQuery := func(key string) string {
		query := url.Values{}
		query.Set("sort", key)
		query.Set("page", strconv.Itoa(pg.Page))
		query.Set("per_page", strconv.Itoa(pg.PerPage))
		return "/?" + query.Encode()
	}
	return []mrn.TableColumn{
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

func roleBadge(role string) mrn.Node {
	tone := "badge-ghost"
	switch role {
	case "Admin":
		tone = "badge-primary"
	case "Editor":
		tone = "badge-secondary"
	}
	return mrn.DivClass("badge "+tone, mrn.Text(role))
}

func renderCreateUserForm(form createUserFormState) mrn.Node {
	return mrn.DivClass("card bg-base-100 shadow-sm",
		mrn.DivClass("card-body",
			mrn.DivClass("text-xl font-semibold", mrn.Text("Create user")),
			mrn.Form("users/create",
				mrn.FormRow(mrn.FormRowProps{
					ID:          "name",
					Label:       "Name",
					Description: "Enter the display name.",
					Error:       form.Errors["name"],
					Required:    true,
					Control: mrn.TextField(mrn.TextFieldProps{
						ID:          "name",
						Name:        "name",
						Value:       form.Name,
						Placeholder: "name",
						Description: "Enter the display name.",
						Error:       form.Errors["name"],
						Required:    true,
					}),
				}),
				mrn.FormRow(mrn.FormRowProps{
					ID:          "email",
					Label:       "Email",
					Description: "Used for notifications.",
					Error:       form.Errors["email"],
					Required:    true,
					Control: mrn.TextField(mrn.TextFieldProps{
						ID:          "email",
						Name:        "email",
						Value:       form.Email,
						Placeholder: "email",
						Description: "Used for notifications.",
						Error:       form.Errors["email"],
						Required:    true,
					}),
				}),
				mrn.FormRow(mrn.FormRowProps{
					ID:          "start_date",
					Label:       "Start date",
					Description: "Select a date in the active fiscal window.",
					Error:       form.Errors["start_date"],
					Required:    true,
					Control: mrn.TextField(mrn.TextFieldProps{
						ID:          "start_date",
						Name:        "start_date",
						Value:       form.StartDate,
						Type:        "date",
						Description: "Select a date in the active fiscal window.",
						Error:       form.Errors["start_date"],
						Required:    true,
					}),
				}),
				mrn.FormRow(mrn.FormRowProps{
					ID:          "role",
					Label:       "Role",
					Description: "Choose permission scope for this user.",
					Error:       form.Errors["role"],
					Required:    true,
					Control: mrn.Select(mrn.SelectFieldProps{
						ID:   "role",
						Name: "role",
						Options: []mrn.SelectOption{
							{Label: "Admin", Value: "Admin", Selected: form.Role == "Admin"},
							{Label: "Editor", Value: "Editor", Selected: form.Role == "Editor"},
							{Label: "Viewer", Value: "Viewer", Selected: form.Role == "" || form.Role == "Viewer"},
						},
						Description: "Choose permission scope for this user.",
						Error:       form.Errors["role"],
						Required:    true,
					}),
				}),
				mrn.DivClass("divider my-1"),
				mrn.FormRow(mrn.FormRowProps{
					ID:          "workspace",
					Label:       "Workspace",
					Description: "Demo-only field showing radio group markup.",
					Control: mrn.RadioGroup(mrn.RadioGroupProps{
						ID:          "workspace",
						Name:        "workspace",
						Value:       "core",
						Description: "Demo-only field showing radio group markup.",
						Options: []mrn.RadioOption{
							{Label: "Core", Value: "core"},
							{Label: "Growth", Value: "growth"},
							{Label: "Support", Value: "support"},
						},
					}),
				}),
				mrn.FormRow(mrn.FormRowProps{
					ID:          "notes",
					Label:       "Notes",
					Description: "Ignored by the demo action, useful for layout coverage.",
					Control: mrn.Textarea(mrn.TextareaProps{
						ID:          "notes",
						Name:        "notes",
						Placeholder: "Add onboarding context",
						Rows:        3,
						Description: "Ignored by the demo action, useful for layout coverage.",
					}),
				}),
				mrn.DivClass("space-y-3",
					mrn.Checkbox(mrn.CheckboxProps{
						ID:          "send_invite",
						Name:        "send_invite",
						Value:       "yes",
						Checked:     true,
						Label:       "Send invite email",
						Description: "Invite preference.",
					}),
					mrn.Switch(mrn.SwitchProps{
						ID:          "provision_access",
						Name:        "provision_access",
						Value:       "yes",
						Checked:     true,
						Label:       "Provision default access",
						Description: "Access preference.",
					}),
				),
				mrn.DivClass("flex flex-wrap gap-2 pt-2",
					mrn.ComponentSubmitButton("Create", mrn.ComponentProps{Variant: "primary", Size: "sm"}),
					mrn.ComponentButton("Preview", mrn.ComponentProps{Variant: "ghost", Size: "sm", Disabled: true}),
				),
			).Target("#users-workspace"),
		),
	)
}

func renderComponentShowcase(ctx *mrn.Context) mrn.Node {
	users := getUsers(ctx)
	roles := roleCounts(users)

	return mrn.DivClass("grid gap-6 xl:grid-cols-[minmax(0,1fr)_22rem]",
		mrn.DivClass("card bg-base-100 shadow-sm",
			mrn.DivClass("card-body gap-4",
				mrn.DivClass("space-y-1",
					mrn.DivClass("text-xl font-semibold", mrn.Text("Component states")),
					mrn.DivClass("text-sm text-base-content/60", mrn.Text("Alerts, toasts, skeletons, and empty states rendered from Go.")),
				),
				mrn.DivClass("grid gap-3 md:grid-cols-2",
					mrn.ComponentAlert(mrn.AlertProps{
						Title:       "Validation feedback",
						Description: "Form errors, success flashes, and informational alerts share the same feedback API.",
						Icon:        "!",
						Props:       mrn.ComponentProps{Variant: "info", Size: "sm"},
					}),
					mrn.ComponentToast(mrn.ToastProps{
						Title:       "Saved",
						Description: "Toast markup is available for transient status messages.",
						Icon:        "OK",
						Props:       mrn.ComponentProps{Variant: "success", Size: "sm"},
					}),
				),
				mrn.DivClass("grid gap-3 md:grid-cols-2",
					mrn.ComponentSkeleton(mrn.SkeletonProps{Rows: 4, Props: mrn.ComponentProps{Size: "sm"}}),
					mrn.ComponentEmptyState(mrn.EmptyStateProps{
						Title:       "No pending reviews",
						Description: "Empty states can replace tables or panels without extra branching in templates.",
						Icon:        "0",
						Props:       mrn.ComponentProps{Size: "sm"},
					}),
				),
			),
		),
		mrn.DivClass("card bg-base-100 shadow-sm",
			mrn.DivClass("card-body gap-4",
				mrn.DivClass("space-y-1",
					mrn.DivClass("text-xl font-semibold", mrn.Text("Role mix")),
					mrn.DivClass("text-sm text-base-content/60", mrn.Text("Small repeated views stay plain Go functions.")),
				),
				roleMixRow("Admin", roles["Admin"]),
				roleMixRow("Editor", roles["Editor"]),
				roleMixRow("Viewer", roles["Viewer"]),
			),
		),
	)
}

func roleMixRow(role string, count int) mrn.Node {
	return mrn.DivClass("flex items-center justify-between rounded-box border border-base-300 px-3 py-2",
		roleBadge(role),
		mrn.DivClass("text-sm font-medium", mrn.Text(strconv.Itoa(count)+" users")),
	)
}

func renderDeleteModal(ctx *mrn.Context) mrn.Node {
	targetID, _ := ctx.Get("deleteTargetID").(int)
	targetName := ""
	for _, u := range getUsers(ctx) {
		if u.ID == targetID {
			targetName = u.Name
			break
		}
	}

	return mrn.ComponentModal(mrn.ModalProps{
		Title: "Delete user",
		Body: mrn.DivClass("space-y-2",
			mrn.Text("Are you sure you want to delete this user?"),
			mrn.DivClass("text-sm text-base-content/70", mrn.Text(targetName)),
		),
		Actions: mrn.DivClass("flex w-full justify-end gap-2",
			mrn.Form("users/delete/cancel",
				mrn.ComponentSubmitButton("Cancel", mrn.ComponentProps{Variant: "ghost", Size: "sm"}),
			).Target("#users-workspace"),
			mrn.Form("users/delete/confirm",
				mrn.HiddenInput("id", strconv.Itoa(targetID)),
				mrn.ComponentSubmitButton("Delete", mrn.ComponentProps{Variant: "danger", Size: "sm"}),
			).Target("#users-workspace"),
		),
		Open: ctx.Get("deleteModalOpen") == true,
	})
}
