package marionette

import (
	"fmt"
	"net/http"
	"strings"
)

// Context gives handlers controlled access to application state and request data.
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	State   map[string]any
}

func (c *Context) Param(name string) string {
	if c.Request == nil {
		return ""
	}
	return c.Request.PathValue(name)
}

func (c *Context) FormValue(name string) string {
	if c.Request == nil {
		return ""
	}
	return c.Request.FormValue(name)
}

func (c *Context) Query(name string) string {
	if c.Request == nil {
		return ""
	}
	return c.Request.URL.Query().Get(name)
}

func (c *Context) Set(key string, value any) {
	c.State[key] = value
}

func (c *Context) Get(key string) any {
	return c.State[key]
}

func (c *Context) GetInt(key string) int {
	v, ok := c.State[key].(int)
	if !ok {
		return 0
	}
	return v
}

// Handler transforms state into a UI node in response to a user event.
type Handler func(*Context) Node

// App is a minimal Go-only UI runtime for htmx driven desktop/web views.
type App struct {
	state   map[string]any
	pages   map[string]Handler
	actions map[string]Handler
}

func New() *App {
	return &App{
		state:   map[string]any{},
		pages:   map[string]Handler{},
		actions: map[string]Handler{},
	}
}

func (a *App) Set(key string, value any) { a.state[key] = value }

func (a *App) GetInt(key string) int {
	v, ok := a.state[key].(int)
	if !ok {
		return 0
	}
	return v
}

// Page registers a full-page GET view.
func (a *App) Page(path string, fn Handler) {
	a.pages[normalizePagePath(path)] = fn
}

// Action registers a POST-only htmx endpoint. name should not include leading slash.
func (a *App) Action(name string, fn Handler) {
	a.actions[normalizeActionPath(name)] = fn
}

// Render defines the main root view for initial load.
func (a *App) Render(fn Handler) {
	a.Page("/", fn)
}

// Handle registers an htmx endpoint. name should not include leading slash.
func (a *App) Handle(name string, fn Handler) {
	a.Action(name, fn)
}

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	for path, fn := range a.pages {
		localPath := path
		localFn := fn
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != localPath {
				http.NotFound(w, r)
				return
			}
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			ctx := a.newContext(w, r)
			renderAndWritePage(w, localFn(ctx))
		})
	}
	for path, fn := range a.actions {
		localFn := fn
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ctx := a.newContext(w, r)
			renderAndWriteFragment(w, localFn(ctx))
		})
	}
	if _, ok := a.pages["/"]; !ok {
		mux.HandleFunc("/", a.handleMissingIndex)
	}
	return mux
}

func (a *App) newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{Writer: w, Request: r, State: a.state}
}

func (a *App) handleMissingIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Error(w, "missing app.Page or app.Render registration for /", http.StatusInternalServerError)
}

func renderAndWritePage(w http.ResponseWriter, node Node) {
	rootHTML, err := node.Render()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page, err := shell(rootHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeHTML(w, page)
}

func (a *App) Run(addr string) error {
	fmt.Printf("marionette listening at http://%s\n", addr)
	return http.ListenAndServe(addr, a.Handler())
}

func renderAndWriteFragment(w http.ResponseWriter, node Node) {
	htmlOut, err := node.Render()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeHTML(w, string(htmlOut))
}

func normalizePagePath(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func normalizeActionPath(name string) string {
	return "/" + strings.TrimPrefix(name, "/")
}
