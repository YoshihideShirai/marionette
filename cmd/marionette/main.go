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
		}
		return renderUsersWorkspace(ctx)
	})

	app.Action("users/delete", func(ctx *marionette.Context) marionette.Node {
		id, _ := strconv.Atoi(ctx.FormValue("id"))
		users := getUsers(ctx)
		next := users[:0]
		for _, u := range users {
			if u.ID != id {
				next = append(next, u)
			}
		}
		ctx.Set("users", next)
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
	return marionette.DivClass("app", "space-y-6",
		marionette.DivClass("", "space-y-2",
			marionette.DivClass("", "text-3xl font-bold tracking-tight", marionette.Text("Marionette Admin UI")),
			marionette.DivClass("", "text-base-content/70",
				marionette.Text("Go handlers, htmx actions, and daisyUI components for small admin tools."),
			),
		),
		renderUsersWorkspace(ctx),
	)
}

func renderUsersWorkspace(ctx *marionette.Context) marionette.Node {
	return marionette.DivClass("users-workspace", "grid gap-6 lg:grid-cols-[minmax(0,1fr)_22rem]",
		renderUsersTable(ctx),
		renderCreateUserForm(),
	)
}

func renderUsersTable(ctx *marionette.Context) marionette.Node {
	rows := []marionette.Node{
		marionette.Raw(`<div class="grid grid-cols-[1fr_1.4fr_.8fr_auto] gap-3 border-b border-base-300 px-4 py-3 text-sm font-semibold text-base-content/70">
<span>Name</span><span>Email</span><span>Role</span><span></span>
</div>`),
	}
	for _, u := range getUsers(ctx) {
		rows = append(rows, renderUserRow(u))
	}
	if len(rows) == 1 {
		rows = append(rows, marionette.DivClass("", "px-4 py-8 text-center text-base-content/60", marionette.Text("No users yet.")))
	}

	return marionette.DivClass("", "card bg-base-100 shadow-sm",
		marionette.DivClass("", "card-body gap-4",
			marionette.DivClass("", "flex items-center justify-between gap-4",
				marionette.DivClass("", "space-y-1",
					marionette.DivClass("", "text-xl font-semibold", marionette.Text("Users")),
					marionette.DivClass("", "text-sm text-base-content/60", marionette.Text("Create and remove users with htmx-backed actions.")),
				),
				marionette.DivClass("", "badge badge-outline", marionette.Text(strconv.Itoa(len(getUsers(ctx)))+" total")),
			),
			marionette.DivClass("", "overflow-hidden rounded-box border border-base-300", rows...),
		),
	)
}

func renderUserRow(u user) marionette.Node {
	return marionette.Raw(`<div class="grid grid-cols-[1fr_1.4fr_.8fr_auto] items-center gap-3 border-b border-base-200 px-4 py-3 last:border-b-0">
<div class="font-medium">` + escape(u.Name) + `</div>
<div class="text-sm text-base-content/70">` + escape(u.Email) + `</div>
<div><span class="badge badge-ghost">` + escape(u.Role) + `</span></div>
<form hx-post="/users/delete" hx-target="#users-workspace" hx-swap="outerHTML">
  <input type="hidden" name="id" value="` + strconv.Itoa(u.ID) + `" />
  <button class="btn btn-ghost btn-sm text-error" type="submit">Delete</button>
</form>
</div>`)
}

func renderCreateUserForm() marionette.Node {
	return marionette.DivClass("", "card bg-base-100 shadow-sm",
		marionette.DivClass("", "card-body",
			marionette.DivClass("", "text-xl font-semibold", marionette.Text("Create user")),
			marionette.Form("users/create",
				marionette.Input("name", ""),
				marionette.Input("email", ""),
				marionette.Input("role", "Viewer"),
				marionette.Submit("Create"),
			).Target("#users-workspace"),
		),
	)
}

func escape(v string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&#34;",
		"'", "&#39;",
	)
	return replacer.Replace(v)
}
