package marionette

import (
	"time"

	backend "github.com/YoshihideShirai/marionette/backend"
)

type AssetOptions = backend.AssetOptions

type AssetOption = backend.AssetOption

// WithAssetCache adds a public Cache-Control max-age to served assets.
func WithAssetCache(maxAge time.Duration) AssetOption { return backend.WithAssetCache(maxAge) }

// WithAssetImmutable adds immutable to Cache-Control when asset caching is enabled.
func WithAssetImmutable() AssetOption { return backend.WithAssetImmutable() }

// WithAssetIndex allows directory index responses from the underlying file server.
func WithAssetIndex(enabled bool) AssetOption { return backend.WithAssetIndex(enabled) }

// WithAssetDownload serves assets with Content-Disposition: attachment.
func WithAssetDownload() AssetOption { return backend.WithAssetDownload() }

// WithAssetContentTypes sets Content-Type by file extension before serving assets.
func WithAssetContentTypes(types map[string]string) AssetOption {
	return backend.WithAssetContentTypes(types)
}
