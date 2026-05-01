package backend

import (
	"bytes"
	"html/template"
	"net/http"
)

var shellTmpl = template.Must(template.New("shell").Parse(`<!doctype html>
<html lang="en" data-theme="corporate">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Marionette</title>
    <link href="https://cdn.jsdelivr.net/npm/daisyui@5" rel="stylesheet" type="text/css" />
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
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
  </head>
  <body class="bg-base-200 min-h-screen">
    <main id="marionette-root" class="container mx-auto p-6">{{.Content}}</main>
  </body>
</html>`))

type shellOptions struct {
	Stylesheets []string
	Styles      []template.CSS
}

func shell(content template.HTML) (string, error) {
	return shellWithOptions(content, shellOptions{})
}

func shellWithOptions(content template.HTML, options shellOptions) (string, error) {
	view := struct {
		Content     template.HTML
		Stylesheets []string
		Styles      []template.CSS
	}{
		Content:     content,
		Stylesheets: options.Stylesheets,
		Styles:      options.Styles,
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
