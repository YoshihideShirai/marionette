package backend

import (
	"html/template"
	"net/http"

	frontend "github.com/YoshihideShirai/marionette/frontend"
)

type shellOptions = frontend.ShellOptions

func shell(content template.HTML) (string, error) {
	return frontend.Shell(content)
}

func shellWithOptions(content template.HTML, options shellOptions) (string, error) {
	return frontend.ShellWithOptions(content, options)
}

func writeHTML(w http.ResponseWriter, body string) {
	frontend.WriteHTML(w, body)
}
