package backend

import (
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type assetRoute struct {
	prefix  string
	fsys    fs.FS
	options AssetOptions
}

type AssetOptions struct {
	MaxAge       time.Duration
	Immutable    bool
	Index        bool
	Download     bool
	ContentTypes map[string]string
}

type AssetOption func(*AssetOptions)

// WithAssetCache adds a public Cache-Control max-age to served assets.
func WithAssetCache(maxAge time.Duration) AssetOption {
	return func(options *AssetOptions) {
		options.MaxAge = maxAge
	}
}

// WithAssetImmutable adds immutable to Cache-Control when asset caching is enabled.
func WithAssetImmutable() AssetOption {
	return func(options *AssetOptions) {
		options.Immutable = true
	}
}

// WithAssetIndex allows directory index responses from the underlying file server.
func WithAssetIndex(enabled bool) AssetOption {
	return func(options *AssetOptions) {
		options.Index = enabled
	}
}

// WithAssetDownload serves assets with Content-Disposition: attachment.
func WithAssetDownload() AssetOption {
	return func(options *AssetOptions) {
		options.Download = true
	}
}

// WithAssetContentTypes sets Content-Type by file extension before serving assets.
func WithAssetContentTypes(types map[string]string) AssetOption {
	return func(options *AssetOptions) {
		if len(types) == 0 {
			return
		}
		options.ContentTypes = make(map[string]string, len(types))
		for ext, contentType := range types {
			ext = normalizeAssetExtension(ext)
			contentType = strings.TrimSpace(contentType)
			if ext == "" || contentType == "" {
				continue
			}
			options.ContentTypes[ext] = contentType
		}
	}
}

// Assets serves files from fsys under prefix, for example /assets/app.css.
func (a *App) Assets(prefix string, fsys fs.FS, options ...AssetOption) {
	if fsys == nil {
		return
	}
	normalized := normalizeAssetPrefix(prefix)
	if normalized == "" {
		return
	}
	route := assetRoute{
		prefix:  normalized,
		fsys:    fsys,
		options: applyAssetOptions(options),
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.assets = append(a.assets, route)
}

// Downloads serves files from fsys under prefix as attachment downloads.
func (a *App) Downloads(prefix string, fsys fs.FS, options ...AssetOption) {
	downloadOptions := make([]AssetOption, 0, len(options)+1)
	downloadOptions = append(downloadOptions, WithAssetDownload())
	downloadOptions = append(downloadOptions, options...)
	a.Assets(prefix, fsys, downloadOptions...)
}

// Asset returns an application asset URL using the first registered asset prefix.
func (a *App) Asset(name string) string {
	assetName := normalizeAssetName(name)
	if assetName == "" {
		return ""
	}
	if isAbsoluteAssetURL(assetName) {
		return assetName
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.assets) == 0 {
		return "/" + escapeAssetPath(assetName)
	}
	return a.assets[0].prefix + "/" + escapeAssetPath(assetName)
}

// Asset returns an application asset URL using the parent app.
func (c *Context) Asset(name string) string {
	if c.app == nil {
		return "/" + escapeAssetPath(normalizeAssetName(name))
	}
	return c.app.Asset(name)
}

func (a *App) registerAssetRoutes(mux *http.ServeMux) {
	a.mu.RLock()
	routes := append([]assetRoute(nil), a.assets...)
	a.mu.RUnlock()
	for _, route := range routes {
		localRoute := route
		mux.Handle(localRoute.prefix+"/", localRoute.handler())
	}
}

func (r assetRoute) handler() http.Handler {
	fileServer := http.StripPrefix(r.prefix+"/", http.FileServer(http.FS(r.fsys)))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet && req.Method != http.MethodHead {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name, err := url.PathUnescape(strings.TrimPrefix(req.URL.Path, r.prefix+"/"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		name = strings.TrimPrefix(path.Clean("/"+name), "/")
		if name == "." || name == "" {
			http.NotFound(w, req)
			return
		}
		info, err := fs.Stat(r.fsys, name)
		if err != nil {
			if fsErrNotExist(err) {
				http.NotFound(w, req)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if info.IsDir() && !r.options.Index {
			http.NotFound(w, req)
			return
		}
		r.setHeaders(w, name)
		fileServer.ServeHTTP(w, req)
	})
}

func (r assetRoute) setHeaders(w http.ResponseWriter, name string) {
	if contentType := r.contentType(name); contentType != "" && w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", contentType)
	}
	if r.options.Download {
		w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{
			"filename": path.Base(name),
		}))
	}
	if r.options.MaxAge > 0 {
		value := fmt.Sprintf("public, max-age=%d", int(r.options.MaxAge.Seconds()))
		if r.options.Immutable {
			value += ", immutable"
		}
		w.Header().Set("Cache-Control", value)
	}
}

func (r assetRoute) contentType(name string) string {
	ext := normalizeAssetExtension(path.Ext(name))
	if ext == "" {
		return ""
	}
	if contentType := r.options.ContentTypes[ext]; contentType != "" {
		return contentType
	}
	return mime.TypeByExtension(ext)
}

func applyAssetOptions(options []AssetOption) AssetOptions {
	var out AssetOptions
	for _, option := range options {
		if option != nil {
			option(&out)
		}
	}
	return out
}

func normalizeAssetPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" || prefix == "/" {
		return ""
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	return strings.TrimSuffix(path.Clean(prefix), "/")
}

func normalizeAssetName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" || isAbsoluteAssetURL(name) {
		return name
	}
	return strings.TrimPrefix(path.Clean("/"+name), "/")
}

func normalizeAssetExtension(ext string) string {
	ext = strings.TrimSpace(strings.ToLower(ext))
	if ext == "" {
		return ""
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return ext
}

func escapeAssetPath(name string) string {
	segments := strings.Split(name, "/")
	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}
	return strings.Join(segments, "/")
}

func isAbsoluteAssetURL(name string) bool {
	return strings.HasPrefix(name, "http://") ||
		strings.HasPrefix(name, "https://") ||
		strings.HasPrefix(name, "data:")
}

func fsErrNotExist(err error) bool {
	return errors.Is(err, fs.ErrNotExist)
}
