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

// Handler transforms state into a UI node in response to a user event.
type Handler func(*Context) Node

// App is a minimal Go-only UI runtime for htmx driven desktop/web views.
type App struct {
	state    map[string]any
	routes   map[string]Handler
	onRender func(*Context) Node
}

func New() *App {
	return &App{
		state:  map[string]any{},
		routes: map[string]Handler{},
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

// Render defines the main root view for initial load.
func (a *App) Render(fn Handler) {
	a.onRender = fn
}

// Handle registers an htmx endpoint. name should not include leading slash.
func (a *App) Handle(name string, fn Handler) {
	a.routes[strings.TrimPrefix(name, "/")] = fn
}

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.handleIndex)
	for name, fn := range a.routes {
		path := "/" + name
		localFn := fn
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			ctx := &Context{Writer: w, Request: r, State: a.state}
			node := localFn(ctx)
			renderAndWriteNode(w, node)
		})
	}
	return mux
}

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	if a.onRender == nil {
		http.Error(w, "missing app.Render registration", http.StatusInternalServerError)
		return
	}
	ctx := &Context{Writer: w, Request: r, State: a.state}
	root := a.onRender(ctx)
	rootHTML, err := root.Render()
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

func renderAndWriteNode(w http.ResponseWriter, node Node) {
	htmlOut, err := node.Render()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeHTML(w, string(htmlOut))
}
