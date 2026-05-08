package marionette

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	backend "github.com/YoshihideShirai/marionette/backend"
	frontend "github.com/YoshihideShirai/marionette/frontend"
	"github.com/YoshihideShirai/marionette/frontend/assets"
)

// Context gives handlers controlled access to application state and request data.
type Context = backend.Context

type FlashLevel = backend.FlashLevel

type FlashMessage = backend.FlashMessage

const (
	FlashSuccess = backend.FlashSuccess
	FlashError   = backend.FlashError
	FlashInfo    = backend.FlashInfo
	FlashWarn    = backend.FlashWarn
)

const flashCookieName = "marionette_flash"
const sessionCookieName = "marionette_session"

// Handler transforms state into a UI node in response to a user event.
type Handler = backend.Handler

// PageOptions configures the full-page HTML shell for a page route.
type PageOptions = backend.PageOptions

// PageOption updates page route options.
type PageOption = backend.PageOption

// WithTitle sets the HTML document title for a page route.
func WithTitle(title string) PageOption { return backend.WithTitle(title) }

// AssetPolicy controls which asset URLs Marionette may emit into generated shells.
type AssetPolicy = assets.AssetPolicy

// AssetMode describes whether generated shells may depend on network-hosted assets.
type AssetMode = backend.AssetMode

const (
	// AssetModeOnline allows Marionette's default CDN-backed asset resolution.
	AssetModeOnline = backend.AssetModeOnline
	// AssetModeOffline rejects http:// and https:// asset URLs during shell rendering.
	AssetModeOffline = backend.AssetModeOffline
)

// StyleTemplate configures framework-level shell stylesheet and script imports.
type StyleTemplate = frontend.StyleTemplate

var (
	// DaisyUITemplate is the default shell style preset used by the backend runtime.
	DaisyUITemplate = frontend.DaisyUITemplate
	// TailwindCSSTemplate is a Tailwind CSS-only shell style preset.
	TailwindCSSTemplate = frontend.TailwindCSSTemplate
)

// StyleTemplateByName returns a built-in style template preset by name.
func StyleTemplateByName(name string) (StyleTemplate, bool) {
	return frontend.StyleTemplateByName(name)
}

// App is the root-package compatibility alias for backend.App.
//
// New code should prefer importing github.com/YoshihideShirai/marionette/backend
// directly. The root package delegates runtime behavior to backend.App so the
// two import paths do not drift apart.
type App = backend.App

// New creates a Marionette application backed by backend.App.
func New() *App { return backend.New() }

// decodeSession is retained for root-package tests and compatibility with
// unexported historical helpers. Runtime session handling lives in backend.App.
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

// encodeSession is retained for root-package tests and compatibility with
// unexported historical helpers. Runtime session handling lives in backend.App.
func encodeSession(session map[string]string) (string, error) {
	encodedJSON, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encodedJSON), nil
}
