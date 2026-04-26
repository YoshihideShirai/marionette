package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/example/marionette/internal/marionette"
)

type user struct {
	ID    int
	Name  string
	Email string
	Role  string
}

func main() {
	app := marionette.New()
	app.Set("nextUserID", 4)
	app.Set("users", []user{
		{ID: 1, Name: "Aiko Tanaka", Email: "aiko@example.com", Role: "Admin"},
		{ID: 2, Name: "Ren Sato", Email: "ren@example.com", Role: "Editor"},
		{ID: 3, Name: "Mina Suzuki", Email: "mina@example.com", Role: "Viewer"},
	})
	app.Set("deleteModalOpen", false)
	app.Set("deleteTargetID", 0)

	app.Page("/", func(ctx *marionette.Context) marionette.Node {
		return renderUsersPage(ctx)
	})

	app.Action("users/create", func(ctx *marionette.Context) marionette.Node {
		name := strings.TrimSpace(ctx.FormValue("name"))
		email := strings.TrimSpace(ctx.FormValue("email"))
		role := strings.TrimSpace(ctx.FormValue("role"))
		if name != "" && email != "" {
			if role == "" {
				role = "Viewer"
			}
			users := getUsers(ctx)
			nextID := ctx.GetInt("nextUserID")
			users = append(users, user{ID: nextID, Name: name, Email: email, Role: role})
			ctx.Set("users", users)
			ctx.Set("nextUserID", nextID+1)
			ctx.FlashSuccess("User was saved successfully.")
		} else {
			ctx.FlashError("Save failed. Name and email are required.")
		}
		return renderUsersWorkspace(ctx)
	})

	app.Action("users/delete/prompt", func(ctx *marionette.Context) marionette.Node {
		id, _ := strconv.Atoi(ctx.FormValue("id"))
		ctx.Set("deleteTargetID", id)
		ctx.Set("deleteModalOpen", true)
		return renderUsersWorkspace(ctx)
	})

	app.Action("users/delete/cancel", func(ctx *marionette.Context) marionette.Node {
		ctx.Set("deleteModalOpen", false)
		ctx.Set("deleteTargetID", 0)
		return renderUsersWorkspace(ctx)
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
		return renderUsersWorkspace(ctx)
	})

	if err := app.Run("127.0.0.1:8080"); err != nil {
		log.Fatal(err)
	}
}

func getUsers(ctx *marionette.Context) []user {
	users, ok := ctx.Get("users").([]user)
	if !ok {
		return nil
	}
	return users
}

func renderUsersPage(ctx *marionette.Context) marionette.Node {
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
			renderUsersWorkspace(ctx),
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

func renderUsersWorkspace(ctx *marionette.Context) marionette.Node {
	return marionette.DivClass("users-workspace", "space-y-4",
		marionette.FlashAlerts(ctx.Flashes()),
		marionette.DivClass("", "grid gap-6 lg:grid-cols-[minmax(0,1fr)_22rem]",
			renderUsersTable(ctx),
			renderCreateUserForm(),
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
	return marionette.Table([]string{"Name", "Email", "Role", ""}, rows...)
}

func renderUserRow(u user) marionette.TableRowData {
	return marionette.TableRow(
		marionette.DivClass("", "font-medium", marionette.Text(u.Name)),
		marionette.DivClass("", "text-sm text-base-content/70", marionette.Text(u.Email)),
		marionette.DivClass("", "badge badge-ghost", marionette.Text(u.Role)),
		marionette.Form("users/delete/prompt",
			marionette.HiddenInput("id", strconv.Itoa(u.ID)),
			marionette.ComponentSubmitButton("Delete", marionette.ComponentProps{Variant: "danger", Size: "sm"}),
		).Target("#users-workspace"),
	)
}

func renderCreateUserForm() marionette.Node {
	return marionette.DivClass("", "card bg-base-100 shadow-sm",
		marionette.DivClass("", "card-body",
			marionette.DivClass("", "text-xl font-semibold", marionette.Text("Create user")),
			marionette.Form("users/create",
				marionette.ComponentInput("name", "", marionette.ComponentProps{Variant: "default", Size: "sm"}),
				marionette.ComponentInput("email", "", marionette.ComponentProps{Variant: "default", Size: "sm"}),
				marionette.ComponentSelect("role", []marionette.SelectOption{
					{Label: "Admin", Value: "Admin"},
					{Label: "Editor", Value: "Editor"},
					{Label: "Viewer", Value: "Viewer", Selected: true},
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
