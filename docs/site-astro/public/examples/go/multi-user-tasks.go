package goexamples

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

// Task belongs to one authenticated user.
type Task struct {
	ID     int
	UserID string
	Title  string
	Done   bool
}

// TaskRepository is the boundary between Marionette handlers and business data.
// A production implementation can use SQL, a service client, or another durable store.
type TaskRepository interface {
	ListByUser(userID string) ([]Task, error)
	Create(userID, title string) (Task, error)
}

// MemoryTaskRepository is only for the documentation sample. It is injected into
// handlers and is not stored in App global state.
type MemoryTaskRepository struct {
	mu     sync.Mutex
	nextID int
	byUser map[string][]Task
}

func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		nextID: 3,
		byUser: map[string][]Task{
			"user_aiko": {
				{ID: 1, UserID: "user_aiko", Title: "Review team handoff", Done: false},
			},
			"user_ren": {
				{ID: 2, UserID: "user_ren", Title: "Confirm vendor invoice", Done: true},
			},
		},
	}
}

func (r *MemoryTaskRepository) ListByUser(userID string) ([]Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks := append([]Task(nil), r.byUser[userID]...)
	sort.SliceStable(tasks, func(i, j int) bool { return tasks[i].ID > tasks[j].ID })
	return tasks, nil
}

func (r *MemoryTaskRepository) Create(userID, title string) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task := Task{ID: r.nextID, UserID: userID, Title: title}
	r.nextID++
	r.byUser[userID] = append(r.byUser[userID], task)
	return task, nil
}

// RegisterMultiUserTasksExample wires a small multi-user task app.
func RegisterMultiUserTasksExample(app *mb.App, repo TaskRepository) {
	if repo == nil {
		repo = NewMemoryTaskRepository()
	}

	app.Page("/business-tasks", func(ctx *mb.Context) mf.Node {
		userID, ok := currentAuthenticatedUserID(ctx)
		if !ok {
			return signInCard()
		}

		tasks, err := repo.ListByUser(userID)
		if err != nil {
			return errorCard(err)
		}
		return tasksPage(userID, tasks, taskFormState{})
	})

	app.Action("business-tasks/demo-login", func(ctx *mb.Context) mf.Node {
		// The session value is an opaque handle, not a trusted user ID. In a real
		// app this should be a signed/encrypted session cookie or a random session
		// ID looked up on the server.
		ctx.SetSession("session_handle", "signed-or-opaque-demo-session-a")
		ctx.FlashSuccess("Signed in with an opaque session handle.")
		return mf.Toast(mf.ToastProps{Title: "Signed in", Description: "Reload /business-tasks to view your tasks.", Props: mf.ComponentProps{Variant: "success"}})
	})

	app.Action("business-tasks/create", func(ctx *mb.Context) mf.Node {
		userID, ok := currentAuthenticatedUserID(ctx)
		if !ok {
			ctx.FlashError("Please sign in before creating tasks.")
			return mf.Toast(mf.ToastProps{Title: "Not signed in", Description: "Authentication is required.", Props: mf.ComponentProps{Variant: "error"}})
		}

		form := validateTaskForm(ctx.FormValue("title"))
		if form.Error != "" {
			tasks, _ := repo.ListByUser(userID)
			ctx.FlashError(form.Error)
			return tasksWorkspace(userID, tasks, form,
				mf.Toast(mf.ToastProps{Title: "Validation failed", Description: form.Error, Props: mf.ComponentProps{Variant: "error"}}),
			)
		}

		if _, err := repo.Create(userID, form.Title); err != nil {
			ctx.FlashError("Could not save the task.")
			return tasksWorkspace(userID, nil, form,
				mf.Toast(mf.ToastProps{Title: "Save failed", Description: err.Error(), Props: mf.ComponentProps{Variant: "error"}}),
			)
		}

		tasks, err := repo.ListByUser(userID)
		if err != nil {
			ctx.FlashError("Task was saved, but the list could not be refreshed.")
			return tasksWorkspace(userID, nil, taskFormState{},
				mf.Toast(mf.ToastProps{Title: "Refresh failed", Description: err.Error(), Props: mf.ComponentProps{Variant: "warning"}}),
			)
		}

		ctx.FlashSuccess("Task created.")
		return tasksWorkspace(userID, tasks, taskFormState{},
			mf.Toast(mf.ToastProps{Title: "Task created", Description: form.Title, Props: mf.ComponentProps{Variant: "success"}}),
		)
	})
}

// currentAuthenticatedUserID demonstrates the recommended shape: the session
// contains only a signed/opaque session handle, then trusted user identity is
// resolved by a server-side lookup. Do not store a raw authenticated user ID in
// the browser session and trust it directly from handlers.
func currentAuthenticatedUserID(ctx *mb.Context) (string, bool) {
	sessionHandle := ctx.Session("session_handle")
	return lookupUserIDBySessionHandle(sessionHandle)
}

func lookupUserIDBySessionHandle(sessionHandle string) (string, bool) {
	// Pseudocode for production:
	//   claims, ok := signedSessions.Verify(sessionHandle)
	//   return sessionStore.LookupUserID(claims.SessionID)
	// or:
	//   return serverSideSessionTable.LookupUserID(sessionHandle)
	if sessionHandle == "signed-or-opaque-demo-session-a" {
		return "user_aiko", true
	}
	if sessionHandle == "signed-or-opaque-demo-session-b" {
		return "user_ren", true
	}
	return "", false
}

type taskFormState struct {
	Title string
	Error string
}

func validateTaskForm(rawTitle string) taskFormState {
	title := strings.TrimSpace(rawTitle)
	if title == "" {
		return taskFormState{Error: "Task title is required."}
	}
	if len([]rune(title)) > 120 {
		return taskFormState{Title: title, Error: "Task title must be 120 characters or fewer."}
	}
	return taskFormState{Title: title}
}

func signInCard() mf.Node {
	return mf.Container(mf.ContainerProps{MaxWidth: "xl", Centered: true},
		mf.Card(mf.CardProps{Title: "Business tasks", Description: "Sign in stores an opaque session handle, then the server resolves the user."},
			mf.ActionForm(mf.ActionFormProps{Action: "/business-tasks/demo-login", Target: "#auth-feedback", Swap: "innerHTML"},
				mf.SubmitButton("Sign in as Aiko", mf.ComponentProps{Variant: "primary"}),
			),
			mf.Region(mf.RegionProps{ID: "auth-feedback"}),
		),
	)
}

func errorCard(err error) mf.Node {
	return mf.Card(mf.CardProps{Title: "Tasks unavailable", Description: err.Error(), Props: mf.ComponentProps{Variant: "error"}})
}

func tasksPage(userID string, tasks []Task, form taskFormState) mf.Node {
	return mf.Container(mf.ContainerProps{MaxWidth: "4xl", Centered: true},
		mf.PageHeader(mf.PageHeaderProps{Title: "Business tasks", Description: "Repository-backed tasks scoped to the authenticated user."}),
		tasksWorkspace(userID, tasks, form),
	)
}

func tasksWorkspace(userID string, tasks []Task, form taskFormState, notices ...mf.Node) mf.Node {
	children := []mf.Node{
		mf.P(mf.Text("Current user: " + userID)),
		mf.ActionForm(mf.ActionFormProps{Action: "/business-tasks/create", Target: "#business-tasks-workspace", Swap: "outerHTML", Props: mf.ComponentProps{Class: "space-y-3"}},
			mf.FormRow(mf.FormRowProps{
				ID:       "task-title",
				Label:    "Task title",
				Error:    form.Error,
				Required: true,
				Control:  mf.TextField(mf.TextFieldProps{ID: "task-title", Name: "title", Value: form.Title, Placeholder: "Follow up with customer", Required: true, Error: form.Error}),
			}),
			mf.SubmitButton("Add task", mf.ComponentProps{Variant: "primary"}),
		),
		taskTable(tasks),
	}
	children = append(children, notices...)

	return mf.Region(mf.RegionProps{ID: "business-tasks-workspace", Props: mf.ComponentProps{Class: "space-y-4"}},
		mf.Stack(mf.StackProps{Direction: "column", Gap: "4"}, children...),
	)
}

func taskTable(tasks []Task) mf.Node {
	if len(tasks) == 0 {
		return mf.EmptyState(mf.EmptyStateProps{Title: "No tasks yet", Description: "Create a task to populate this user-scoped list."})
	}

	rows := make([]mf.TableComponentRow, 0, len(tasks))
	for _, task := range tasks {
		status := "Open"
		if task.Done {
			status = "Done"
		}
		rows = append(rows, mf.TableRowValues(fmt.Sprintf("#%d", task.ID), task.Title, status))
	}
	return mf.Table(mf.TableProps{
		Columns: []mf.TableColumn{{Label: "ID"}, {Label: "Task"}, {Label: "Status"}},
		Rows:    rows,
	})
}
