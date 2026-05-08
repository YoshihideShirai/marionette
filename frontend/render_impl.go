package frontend

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"github.com/YoshihideShirai/marionette/frontend/assets"
)

var shellTmpl = template.Must(template.New("shell").Parse(`<!doctype html>
<html lang="en" data-theme="corporate">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>{{.Title}}</title>
    {{range .FrameworkStylesheets}}<link href="{{.}}" rel="stylesheet" type="text/css" />
    {{end}}{{range .FrameworkScripts}}<script src="{{.}}"></script>
    {{end}}
    <style>
      :root {
        --mrn-page-max-width: 80rem;
        --mrn-page-padding: clamp(1rem, 2vw + 0.5rem, 2rem);
        --mrn-focus-ring: 0 0 0 3px color-mix(in oklab, var(--color-primary) 28%, transparent);
      }

      html {
        min-height: 100%;
        background: var(--color-base-200);
      }

      body {
        min-height: 100vh;
        margin: 0;
        background:
          radial-gradient(circle at top left, color-mix(in oklab, var(--color-primary) 16%, transparent), transparent 28rem),
          radial-gradient(circle at bottom right, color-mix(in oklab, var(--color-secondary) 10%, transparent), transparent 24rem),
          linear-gradient(180deg, var(--color-base-100) 0%, var(--color-base-200) 42%, var(--color-base-200) 100%);
        color: var(--color-base-content);
        font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
        font-feature-settings: "cv02", "cv03", "cv04", "cv11";
        text-rendering: optimizeLegibility;
        -webkit-font-smoothing: antialiased;
      }

      #marionette-root {
        width: min(100%, var(--mrn-page-max-width));
        padding: var(--mrn-page-padding);
      }

      #marionette-root > * {
        animation: mrn-page-enter 160ms ease-out both;
      }

      :where(a, button, input, select, textarea, [tabindex]):focus-visible {
        outline: none;
        box-shadow: var(--mrn-focus-ring);
      }

      ::selection {
        background: color-mix(in oklab, var(--color-primary) 24%, transparent);
      }

      @keyframes mrn-page-enter {
        from {
          opacity: 0;
          transform: translateY(0.25rem);
        }
        to {
          opacity: 1;
          transform: translateY(0);
        }
      }

      @media (prefers-reduced-motion: reduce) {
        #marionette-root > * {
          animation: none;
        }
      }
    </style>
    {{range .Stylesheets}}<link href="{{.}}" rel="stylesheet" type="text/css" />
    {{end}}{{range .Styles}}<style>{{.}}</style>
    {{end}}
    {{range .Scripts}}<script src="{{.}}"></script>
    {{end}}{{range .JavaScripts}}<script>{{.}}</script>
    {{end}}
  </head>
  <body class="bg-base-200 min-h-screen">
    <main id="marionette-root" class="container mx-auto p-6">{{.Content}}</main>
  </body>
</html>`))

type shellOptions struct {
	Title                string
	StyleTemplate        StyleTemplate
	FrameworkStylesheets []string
	FrameworkScripts     []string
	Stylesheets          []string
	Styles               []template.CSS
	AssetProvider        assets.AssetProvider
	Scripts              []string
	JavaScripts          []template.JS
	DisableHTMX          bool
	// DisableCharts keeps Chart.js out of pages that do not need charts.
	// TODO: replace this opt-out with render-context feature tracking when Chart nodes can mark Chart.js as required.
	DisableCharts bool
}

// ShellOptions configures the HTML document shell rendered by ShellWithOptions.
type ShellOptions = shellOptions

// Shell renders content inside the default Marionette HTML document shell.
func Shell(content template.HTML) (string, error) {
	return shell(content)
}

func shell(content template.HTML) (string, error) {
	return shellWithOptions(content, shellOptions{})
}

// ShellWithOptions renders content inside the Marionette HTML document shell.
func ShellWithOptions(content template.HTML, options ShellOptions) (string, error) {
	return shellWithOptions(content, shellOptions(options))
}

func shellWithOptions(content template.HTML, options shellOptions) (string, error) {
	title := strings.TrimSpace(options.Title)
	if title == "" {
		title = "Marionette"
	}
	provider := assetProvider(options.AssetProvider)
	frameworkStylesheets, frameworkScripts := resolveFrameworkAssets(options, provider)
	scripts := append(resolveFeatureScriptAssets(options, provider), options.Scripts...)
	javaScripts := []template.JS{template.JS(assets.ThemeBootstrapJS)}
	if !options.DisableCharts {
		javaScripts = append(javaScripts, template.JS(assets.ChartBootstrapJS))
	}
	javaScripts = append(javaScripts, options.JavaScripts...)

	view := struct {
		Title                string
		Content              template.HTML
		FrameworkStylesheets []string
		FrameworkScripts     []string
		Stylesheets          []string
		Styles               []template.CSS
		Scripts              []string
		JavaScripts          []template.JS
	}{
		Title:                title,
		Content:              content,
		FrameworkStylesheets: frameworkStylesheets,
		FrameworkScripts:     frameworkScripts,
		Stylesheets:          options.Stylesheets,
		Styles:               options.Styles,
		Scripts:              scripts,
		JavaScripts:          javaScripts,
	}

	var out bytes.Buffer
	if err := shellTmpl.Execute(&out, view); err != nil {
		return "", err
	}
	return out.String(), nil
}

func resolveFeatureScriptAssets(options shellOptions, provider assets.AssetProvider) []string {
	names := make([]assets.AssetName, 0, 2)
	if !options.DisableHTMX {
		names = append(names, assets.HTMX)
	}
	if !options.DisableCharts {
		names = append(names, assets.ChartJS)
	}
	return resolveScriptAssets(provider, names)
}

func resolveFrameworkAssets(options shellOptions, provider assets.AssetProvider) ([]string, []string) {
	if len(options.FrameworkStylesheets) > 0 || len(options.FrameworkScripts) > 0 {
		return append([]string(nil), options.FrameworkStylesheets...), append([]string(nil), options.FrameworkScripts...)
	}
	styleTemplate := options.StyleTemplate
	if styleTemplate.Name == "" && len(styleTemplate.FrameworkStylesheets) == 0 && len(styleTemplate.FrameworkScripts) == 0 && len(styleTemplate.FrameworkStylesheetAssets) == 0 && len(styleTemplate.FrameworkScriptAssets) == 0 {
		styleTemplate = DefaultStyleTemplate()
	}
	stylesheets := resolveStylesheetAssets(provider, styleTemplate.FrameworkStylesheetAssets)
	if len(styleTemplate.FrameworkStylesheetAssets) == 0 {
		stylesheets = append(stylesheets, styleTemplate.FrameworkStylesheets...)
	}
	scripts := resolveScriptAssets(provider, styleTemplate.FrameworkScriptAssets)
	if len(styleTemplate.FrameworkScriptAssets) == 0 {
		scripts = append(scripts, styleTemplate.FrameworkScripts...)
	}
	return stylesheets, scripts
}

func assetProvider(provider assets.AssetProvider) assets.AssetProvider {
	if provider != nil {
		return provider
	}
	return assets.DefaultProvider
}

func resolveStylesheetAssets(provider assets.AssetProvider, names []assets.AssetName) []string {
	out := make([]string, 0, len(names))
	for _, name := range names {
		if url, ok := provider.StylesheetURL(name); ok && url != "" {
			out = append(out, url)
		}
	}
	return out
}

func resolveScriptAssets(provider assets.AssetProvider, names []assets.AssetName) []string {
	out := make([]string, 0, len(names))
	for _, name := range names {
		if url, ok := provider.ScriptURL(name); ok && url != "" {
			out = append(out, url)
		}
	}
	return out
}

// WriteHTML writes an HTML response with Marionette's standard content type.
func WriteHTML(w http.ResponseWriter, body string) {
	writeHTML(w, body)
}

func writeHTML(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(body))
}
