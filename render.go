package marionette

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"
)

var shellTmpl = template.Must(template.New("shell").Parse(`<!doctype html>
<html lang="en" data-theme="corporate">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>{{.Title}}</title>
    <link href="https://cdn.jsdelivr.net/npm/daisyui@5" rel="stylesheet" type="text/css" />
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
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
    <script src="https://unpkg.com/htmx.org@1.9.12"></script>
    <script src="https://cdn.jsdelivr.net/npm/chart.js@4"></script>
    <script>
      (function() {
        var root = document.documentElement;
        var key = "marionette-theme";
        var storedTheme = null;
        try {
          storedTheme = localStorage.getItem(key);
        } catch (e) {}
        var prefersDark = window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;
        var theme = storedTheme || (prefersDark ? "dark" : "corporate");
        root.setAttribute("data-theme", theme);
        window.mrnToggleTheme = function() {
          var next = root.getAttribute("data-theme") === "dark" ? "corporate" : "dark";
          root.setAttribute("data-theme", next);
          try {
            localStorage.setItem(key, next);
          } catch (e) {}
        };
      })();
    </script>
    <script>
      (function() {
        var charts = new WeakMap();

        function initCharts(root) {
          if (!window.Chart) return;
          var scope = root || document;
          var canvases = scope.querySelectorAll ? scope.querySelectorAll("[data-mrn-chart]") : [];
          canvases.forEach(function(canvas) {
            var container = canvas.closest("[data-mrn-chart-root]");
            if (!container) return;
            var configEl = container.querySelector("[data-mrn-chart-config]");
            if (!configEl) return;

            var config;
            try {
              config = JSON.parse(configEl.textContent || "{}");
            } catch (e) {
              return;
            }

            var existing = charts.get(canvas);
            if (existing) existing.destroy();
            charts.set(canvas, new window.Chart(canvas, config));
          });
        }

        document.addEventListener("DOMContentLoaded", function() {
          initCharts(document);
        });
        document.addEventListener("htmx:afterSwap", function(event) {
          initCharts(event.detail && event.detail.elt ? event.detail.elt : document);
        });
        window.mrnInitCharts = initCharts;
      })();
    </script>
    {{range .Scripts}}<script src="{{.}}"></script>
    {{end}}{{range .JavaScripts}}<script>{{.}}</script>
    {{end}}
  </head>
  <body class="bg-base-200 min-h-screen">
    <main id="marionette-root" class="container mx-auto p-6">{{.Content}}</main>
  </body>
</html>`))

type shellOptions struct {
	Title       string
	Stylesheets []string
	Styles      []template.CSS
	Scripts     []string
	JavaScripts []template.JS
}

func shell(content template.HTML) (string, error) {
	return shellWithOptions(content, shellOptions{})
}

func shellWithOptions(content template.HTML, options shellOptions) (string, error) {
	title := strings.TrimSpace(options.Title)
	if title == "" {
		title = "Marionette"
	}
	view := struct {
		Title       string
		Content     template.HTML
		Stylesheets []string
		Styles      []template.CSS
		Scripts     []string
		JavaScripts []template.JS
	}{
		Title:       title,
		Content:     content,
		Stylesheets: options.Stylesheets,
		Styles:      options.Styles,
		Scripts:     options.Scripts,
		JavaScripts: options.JavaScripts,
	}

	var out bytes.Buffer
	if err := shellTmpl.Execute(&out, view); err != nil {
		return "", err
	}
	return out.String(), nil
}

func writeHTML(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(body))
}
