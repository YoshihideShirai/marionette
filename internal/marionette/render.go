package marionette

import "net/http"

func shell(content string) string {
	return `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Marionette</title>
    <script src="https://unpkg.com/htmx.org@1.9.12"></script>
  </head>
  <body>
    <main id="app">` + content + `</main>
  </body>
</html>`
}

func writeHTML(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(body))
}
