package marionette

import (
	"bytes"
	"html/template"
	"net/http"
)

var shellTmpl = template.Must(template.New("shell").Parse(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Marionette</title>
    <script src="https://unpkg.com/htmx.org@1.9.12"></script>
  </head>
  <body>
    <main id="app">{{.Content}}</main>
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
