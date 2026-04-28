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
    <script src="https://unpkg.com/htmx.org@1.9.12"></script>
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
  </head>
  <body class="bg-base-200 min-h-screen">
    <main id="marionette-root" class="container mx-auto p-6">{{.Content}}</main>
  </body>
</html>`))

func shell(content template.HTML) (string, error) {
	view := struct {
		Content template.HTML
	}{Content: content}

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
