package marionette

import (
	"html/template"
	"net/http"

	mf "github.com/YoshihideShirai/marionette/frontend"
)

type shellOptions = mf.ShellOptions

func shell(content template.HTML) (string, error) {
	return mf.Shell(content)
}

func shellWithOptions(content template.HTML, options shellOptions) (string, error) {
	return mf.ShellWithOptions(content, options)
}

func writeHTML(w http.ResponseWriter, body string) {
	mf.WriteHTML(w, body)
}
