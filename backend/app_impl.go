package backend

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"

	frontend "github.com/YoshihideShirai/marionette/frontend"
)

// Context gives handlers controlled access to application state and request data.
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	State   map[string]any
	app     *App
	flashes []FlashMessage
	session map[string]string
}

type FlashLevel = frontend.FlashLevel

type FlashMessage = frontend.FlashMessage

const (
	FlashSuccess = frontend.FlashSuccess
	FlashError   = frontend.FlashError
	FlashInfo    = frontend.FlashInfo
	FlashWarn    = frontend.FlashWarn
)

const flashCookieName = "marionette_flash"
const sessionCookieName = "marionette_session"

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

func (c *Context) SetSession(key, value string) {
	key = strings.TrimSpace(key)
	if key == "" {
		return
	}
	if c.session == nil {
		c.session = map[string]string{}
	}
	c.session[key] = value
	c.writeSessionCookie()
}

func (c *Context) Session(key string) string {
	if c.session == nil {
		return ""
	}
	return c.session[key]
}

func (c *Context) ClearSession() {
	c.session = map[string]string{}
	c.writeSessionCookie()
}

func (c *Context) writeSessionCookie() {
	if c.Writer == nil {
		return
	}
	secure := false
	if c.app != nil {
		c.app.mu.RLock()
		secure = c.app.cookieSecure
		c.app.mu.RUnlock()
	}
	encoded, err := encodeSession(c.session)
	if err != nil {
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     sessionCookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	})
}

// Handler transforms state into a UI node in response to a user event.
type Handler func(*Context) frontend.Node

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
	assets       []assetRoute
	cookieSecure bool
	shellAssets  frontend.ShellAssets
}

func New() *App {
	return &App{
		state:        map[string]any{},
		pages:        map[string]pageRoute{},
		actions:      map[string]Handler{},
		assets:       []assetRoute{},
		cookieSecure: false,
		shellAssets:  frontend.ShellAssets{},
	}
}

// UseStyleTemplate replaces framework stylesheet/script imports.
// Call AddStylesheet/AddScript after this if you want extra imports.
func (a *App) UseStyleTemplate(tpl frontend.StyleTemplate) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.shellAssets.UseStyleTemplate(tpl)
}

// UseStyleTemplateByName applies a built-in style template preset.
func (a *App) UseStyleTemplateByName(name string) error {
	tpl, ok := frontend.StyleTemplateByName(strings.TrimSpace(name))
	if !ok {
		return fmt.Errorf("unknown style template: %s", name)
	}
	a.UseStyleTemplate(tpl)
	return nil
}

func (a *App) UseTailAdminTemplate() {
	a.UseStyleTemplate(frontend.TailAdminTemplate)
}

func (a *App) UseTailwindCSSTemplate() {
	a.UseStyleTemplate(frontend.TailwindCSSTemplate)
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
	a.shellAssets.AddStylesheet(href)
}

// AddStyle adds trusted inline CSS to the full-page HTML shell.
func (a *App) AddStyle(css string) {
	css = strings.TrimSpace(css)
	if css == "" {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.shellAssets.AddStyle(template.CSS(css))
}

// AddScript adds an external JavaScript file to the full-page HTML shell.
func (a *App) AddScript(src string) {
	src = strings.TrimSpace(src)
	if src == "" {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.shellAssets.AddScript(src)
}

// AddJavaScript adds trusted inline JavaScript to the full-page HTML shell.
func (a *App) AddJavaScript(js string) {
	js = strings.TrimSpace(js)
	if js == "" {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.shellAssets.AddJavaScript(template.JS(js))
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
	a.registerAssetRoutes(mux)
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
	session := decodeSession(r)
	if len(flashes) > 0 {
		a.mu.RLock()
		secure := a.cookieSecure
		a.mu.RUnlock()
		clearFlashCookie(w, secure)
	}
	return &Context{Writer: w, Request: r, State: a.state, app: a, flashes: flashes, session: session}
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

func (a *App) renderAndWritePage(w http.ResponseWriter, node frontend.Node, pageOptions PageOptions) {
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
		Title:                pageOptions.Title,
		FrameworkStylesheets: append([]string(nil), a.shellAssets.FrameworkStylesheets...),
		FrameworkScripts:     append([]string(nil), a.shellAssets.FrameworkScripts...),
		Stylesheets:          append([]string(nil), a.shellAssets.Stylesheets...),
		Styles:               append([]template.CSS(nil), a.shellAssets.Styles...),
		Scripts:              append([]string(nil), a.shellAssets.Scripts...),
		JavaScripts:          append([]template.JS(nil), a.shellAssets.JavaScripts...),
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

func renderAndWriteFragment(w http.ResponseWriter, node frontend.Node) {
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

func decodeSession(r *http.Request) map[string]string {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return map[string]string{}
	}
	raw, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return map[string]string{}
	}
	var session map[string]string
	if err := json.Unmarshal(raw, &session); err != nil {
		return map[string]string{}
	}
	if session == nil {
		return map[string]string{}
	}
	return session
}

func encodeSession(session map[string]string) (string, error) {
	encodedJSON, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encodedJSON), nil
}
