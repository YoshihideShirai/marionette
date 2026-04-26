package main

import (
	"log"
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

func defaultCreateUserFormState() createUserFormState {
	return createUserFormState{Role: "Viewer", Errors: map[string]string{}}
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
	app.Set("users", []user{
		{ID: 1, Name: "Aiko Tanaka", Email: "aiko@example.com", Role: "Admin", StartDate: "2024-03-18"},
		{ID: 2, Name: "Ren Sato", Email: "ren@example.com", Role: "Editor", StartDate: "2024-07-01"},
		{ID: 3, Name: "Mina Suzuki", Email: "mina@example.com", Role: "Viewer", StartDate: "2025-01-10"},
	})
	app.Set("deleteModalOpen", false)
	app.Set("deleteTargetID", 0)

	app.Page("/", func(ctx *marionette.Context) marionette.Node {
		return renderUsersPage(ctx, defaultCreateUserFormState())
	})

	app.Action("users/create", func(ctx *marionette.Context) marionette.Node {
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
		id, _ := strconv.Atoi(ctx.FormValue("id"))
		ctx.Set("deleteTargetID", id)
		ctx.Set("deleteModalOpen", true)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/delete/cancel", func(ctx *marionette.Context) marionette.Node {
		ctx.Set("deleteModalOpen", false)
		ctx.Set("deleteTargetID", 0)
		return renderUsersWorkspace(ctx, defaultCreateUserFormState())
	})

	app.Action("users/delete/confirm", func(ctx *marionette.Context) marionette.Node {
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
	return marionette.DivClass("app", "grid gap-6 lg:grid-cols-[16rem_minmax(0,1fr)]",
		renderSidebar(),
		marionette.DivClass("", "min-w-0 space-y-6",
			marionette.FlashAlerts(ctx.Flashes()),
			marionette.DivClass("", "space-y-2",
				marionette.DivClass("", "text-3xl font-bold tracking-tight", marionette.Text("Marionette Admin UI")),
				marionette.DivClass("", "text-base-content/70",
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
	return marionette.DivClass("users-workspace", "space-y-4",
		marionette.FlashAlerts(ctx.Flashes()),
		marionette.DivClass("", "grid gap-6 lg:grid-cols-[minmax(0,1fr)_22rem]",
			renderUsersTable(ctx),
			renderCreateUserForm(formState),
		),
		renderDeleteModal(ctx),
	)
}

func renderUsersTable(ctx *marionette.Context) marionette.Node {
	users := getUsers(ctx)
	tableBody := renderUsersTableBody(users)

	return marionette.DivClass("", "card bg-base-100 shadow-sm",
		marionette.DivClass("", "card-body gap-4",
			marionette.DivClass("", "flex items-center justify-between gap-4",
				marionette.DivClass("", "space-y-1",
					marionette.DivClass("", "text-xl font-semibold", marionette.Text("Users")),
					marionette.DivClass("", "text-sm text-base-content/60", marionette.Text("Create and remove users with htmx-backed actions.")),
				),
				marionette.DivClass("", "badge badge-outline", marionette.Text(strconv.Itoa(len(getUsers(ctx)))+" total")),
			),
			marionette.DivClass("", "overflow-hidden rounded-box border border-base-300", tableBody),
		),
	)
}

func renderUsersTableBody(users []user) marionette.Node {
	if len(users) == 0 {
		return marionette.DivClass("", "px-4 py-8 text-center text-base-content/60", marionette.Text("No users yet."))
	}

	rows := make([]marionette.TableRowData, 0, len(users))
	for _, u := range users {
		rows = append(rows, renderUserRow(u))
	}
	return marionette.Table([]string{"Name", "Email", "Role", "Start date", ""}, rows...)
}

func renderUserRow(u user) marionette.TableRowData {
	return marionette.TableRow(
		marionette.DivClass("", "font-medium", marionette.Text(u.Name)),
		marionette.DivClass("", "text-sm text-base-content/70", marionette.Text(u.Email)),
		marionette.DivClass("", "badge badge-ghost", marionette.Text(u.Role)),
		marionette.DivClass("", "text-sm", marionette.Text(u.StartDate)),
		marionette.Form("users/delete/prompt",
			marionette.HiddenInput("id", strconv.Itoa(u.ID)),
			marionette.ComponentSubmitButton("Delete", marionette.ComponentProps{Variant: "danger", Size: "sm"}),
		).Target("#users-workspace"),
	)
}

func renderCreateUserForm(form createUserFormState) marionette.Node {
	return marionette.DivClass("", "card bg-base-100 shadow-sm",
		marionette.DivClass("", "card-body",
			marionette.DivClass("", "text-xl font-semibold", marionette.Text("Create user")),
			marionette.Form("users/create",
				marionette.ComponentInputWithOptions("name", form.Name, marionette.InputOptions{
					Type:        "text",
					Placeholder: "name",
					Required:    true,
					Error:       form.Errors["name"],
					Props:       marionette.ComponentProps{Variant: "default", Size: "sm"},
				}),
				marionette.ComponentInputWithOptions("email", form.Email, marionette.InputOptions{
					Type:        "text",
					Placeholder: "email",
					Required:    true,
					Error:       form.Errors["email"],
					Props:       marionette.ComponentProps{Variant: "default", Size: "sm"},
				}),
				marionette.ComponentInputWithOptions("start_date", form.StartDate, marionette.InputOptions{
					Type:     "date",
					Min:      startDateMin,
					Max:      startDateMax,
					Required: true,
					Error:    form.Errors["start_date"],
					Props:    marionette.ComponentProps{Variant: "default", Size: "sm"},
				}),
				marionette.ComponentSelect("role", []marionette.SelectOption{
					{Label: "Admin", Value: "Admin", Selected: form.Role == "Admin"},
					{Label: "Editor", Value: "Editor", Selected: form.Role == "Editor"},
					{Label: "Viewer", Value: "Viewer", Selected: form.Role == "" || form.Role == "Viewer"},
				}, marionette.ComponentProps{Variant: "default", Size: "sm"}),
				marionette.ComponentSubmitButton("Create", marionette.ComponentProps{Variant: "primary", Size: "sm"}),
			).Target("#users-workspace"),
			marionette.DivClass("", "pt-2", marionette.ComponentButton("Preview (disabled)", marionette.ComponentProps{Variant: "ghost", Size: "sm", Disabled: true})),
		),
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
		Body: marionette.DivClass("", "space-y-2",
			marionette.Text("Are you sure you want to delete this user?"),
			marionette.DivClass("", "text-sm text-base-content/70", marionette.Text(targetName)),
		),
		Actions: marionette.DivClass("", "flex w-full justify-end gap-2",
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
