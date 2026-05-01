package desktop

import (
	"context"
	"errors"
	"time"

	"github.com/YoshihideShirai/marionette/backend"
)

// Run starts app on a private localhost server and opens it in a native WebView.
func Run(app *backend.App, options Options) error {
	if app == nil {
		return errors.New("desktop: app is nil")
	}
	options = normalizeOptions(options)

	server, err := startLocalServer(app.Handler())
	if err != nil {
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	return openWebView(server.URL, options)
}
