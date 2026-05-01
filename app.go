package marionette

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
)

// Context gives handlers controlled access to application state and request data.
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	State   map[string]any
	app     *App
	flashes []FlashMessage
}

type FlashLevel string

const (
	FlashSuccess FlashLevel = "success"
	FlashError   FlashLevel = "error"
	FlashInfo    FlashLevel = "info"
	FlashWarn    FlashLevel = "warn"
)

type FlashMessage struct {
	Level   FlashLevel `json:"level"`
	Message string     `json:"message"`
}

const flashCookieName = "marionette_flash"

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
	if c.app == nil {
		c.State[key] = value
		return
	}
	c.app.mu.Lock()
	defer c.app.mu.Unlock()
	c.app.state[key] = value
}

func (c *Context) Get(key string) any {
	if c.app == nil {
		return c.State[key]
	}
	c.app.mu.RLock()
	defer c.app.mu.RUnlock()
	return c.app.state[key]
}

func (c *Context) GetInt(key string) int {
	v, ok := c.Get(key).(int)
	if !ok {
		return 0
	}
	return v
}

func (c *Context) AddFlash(level FlashLevel, message string) {
	secure := false
	if c.app != nil {
		c.app.mu.RLock()
		secure = c.app.cookieSecure
		c.app.mu.RUnlock()
	}

	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		return
	}
	c.flashes = append(c.flashes, FlashMessage{Level: level, Message: trimmed})
	encoded, err := encodeFlashes(c.flashes)
	if err != nil {
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     flashCookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	})
}

func (c *Context) FlashSuccess(message string) { c.AddFlash(FlashSuccess, message) }
func (c *Context) FlashError(message string)   { c.AddFlash(FlashError, message) }
func (c *Context) FlashInfo(message string)    { c.AddFlash(FlashInfo, message) }
func (c *Context) FlashWarn(message string)    { c.AddFlash(FlashWarn, message) }

func (c *Context) Flashes() []FlashMessage {
	if len(c.flashes) == 0 {
		return nil
	}
	out := make([]FlashMessage, len(c.flashes))
	copy(out, c.flashes)
	return out
}

// Handler transforms state into a UI node in response to a user event.
type Handler func(*Context) Node

// PageOptions configures the full-page HTML shell for a page route.
type PageOptions struct {
	Title string
}

// PageOption updates page route options.
type PageOption func(*PageOptions)

// WithTitle sets the HTML document title for a page route.
func WithTitle(title string) PageOption {
	return func(options *PageOptions) {
		options.Title = strings.TrimSpace(title)
	}
}

type pageRoute struct {
	handler Handler
	options PageOptions
}

// App is a minimal Go-only UI runtime for htmx driven desktop/web views.
type App struct {
	mu           sync.RWMutex
	state        map[string]any
	pages        map[string]pageRoute
	actions      map[string]Handler
	cookieSecure bool
	stylesheets  []string
	styles       []template.CSS
	scripts      []string
	javascripts  []template.JS
}

func New() *App {
	return &App{
		state:        map[string]any{},
		pages:        map[string]pageRoute{},
		actions:      map[string]Handler{},
		cookieSecure: false,
		stylesheets:  []string{},
		styles:       []template.CSS{},
		scripts:      []string{},
		javascripts:  []template.JS{},
	}
}

func (a *App) SetCookieSecure(secure bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.cookieSecure = secure
}

// AddStylesheet adds a stylesheet link to the full-page HTML shell.
func (a *App) AddStylesheet(href string) {
	href = strings.TrimSpace(href)
	if href == "" {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.stylesheets = append(a.stylesheets, href)
}

// AddStyle adds trusted inline CSS to the full-page HTML shell.
func (a *App) AddStyle(css string) {
	css = strings.TrimSpace(css)
	if css == "" {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.styles = append(a.styles, template.CSS(css))
}

// AddScript adds an external JavaScript file to the full-page HTML shell.
func (a *App) AddScript(src string) {
	src = strings.TrimSpace(src)
	if src == "" {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.scripts = append(a.scripts, src)
}

// AddJavaScript adds trusted inline JavaScript to the full-page HTML shell.
func (a *App) AddJavaScript(js string) {
	js = strings.TrimSpace(js)
	if js == "" {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.javascripts = append(a.javascripts, template.JS(js))
}

func (a *App) Set(key string, value any) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.state[key] = value
}

func (a *App) GetInt(key string) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	v, ok := a.state[key].(int)
	if !ok {
		return 0
	}
	return v
}

// Page registers a full-page GET view.
func (a *App) Page(path string, fn Handler, options ...PageOption) {
	a.pages[normalizePagePath(path)] = pageRoute{handler: fn, options: applyPageOptions(options)}
}

// Action registers a POST-only htmx endpoint. name should not include leading slash.
func (a *App) Action(name string, fn Handler) {
	a.actions[normalizeActionPath(name)] = fn
}

// Render defines the main root view for initial load.
func (a *App) Render(fn Handler, options ...PageOption) {
	a.Page("/", fn, options...)
}

// Handle registers an htmx endpoint. name should not include leading slash.
func (a *App) Handle(name string, fn Handler) {
	a.Action(name, fn)
}

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	for path, route := range a.pages {
		localPath := path
		localRoute := route
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
			a.renderAndWritePage(w, localRoute.handler(ctx), localRoute.options)
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
	flashes := decodeFlashes(r)
	if len(flashes) > 0 {
		a.mu.RLock()
		secure := a.cookieSecure
		a.mu.RUnlock()
		clearFlashCookie(w, secure)
	}
	return &Context{Writer: w, Request: r, State: a.state, app: a, flashes: flashes}
}

func clearFlashCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     flashCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	})
}

func decodeFlashes(r *http.Request) []FlashMessage {
	cookie, err := r.Cookie(flashCookieName)
	if err != nil || cookie.Value == "" {
		return nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil
	}
	var flashes []FlashMessage
	if err := json.Unmarshal(raw, &flashes); err != nil {
		return nil
	}
	out := make([]FlashMessage, 0, len(flashes))
	for _, f := range flashes {
		if strings.TrimSpace(f.Message) == "" {
			continue
		}
		switch f.Level {
		case FlashSuccess, FlashError, FlashInfo, FlashWarn:
			out = append(out, f)
		}
	}
	return out
}

func encodeFlashes(flashes []FlashMessage) (string, error) {
	encodedJSON, err := json.Marshal(flashes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encodedJSON), nil
}

func (a *App) handleMissingIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Error(w, "missing app.Page or app.Render registration for /", http.StatusInternalServerError)
}

func (a *App) renderAndWritePage(w http.ResponseWriter, node Node, pageOptions PageOptions) {
	rootHTML, err := node.Render()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page, err := shellWithOptions(rootHTML, a.shellOptions(pageOptions))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeHTML(w, page)
}

func (a *App) shellOptions(pageOptions PageOptions) shellOptions {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return shellOptions{
		Title:       pageOptions.Title,
		Stylesheets: append([]string(nil), a.stylesheets...),
		Styles:      append([]template.CSS(nil), a.styles...),
		Scripts:     append([]string(nil), a.scripts...),
		JavaScripts: append([]template.JS(nil), a.javascripts...),
	}
}

func applyPageOptions(options []PageOption) PageOptions {
	var pageOptions PageOptions
	for _, option := range options {
		if option != nil {
			option(&pageOptions)
		}
	}
	return pageOptions
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
