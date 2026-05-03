package frontend

import "html/template"

type ShellAssets struct {
	FrameworkStylesheets []string
	FrameworkScripts     []string
	Stylesheets          []string
	Styles               []template.CSS
	Scripts              []string
	JavaScripts          []template.JS
}

func (a *ShellAssets) UseStyleTemplate(tpl StyleTemplate) {
	a.FrameworkStylesheets = append([]string(nil), tpl.FrameworkStylesheets...)
	a.FrameworkScripts = append([]string(nil), tpl.FrameworkScripts...)
}

func (a *ShellAssets) AddStylesheet(href string) { a.Stylesheets = append(a.Stylesheets, href) }
func (a *ShellAssets) AddStyle(css template.CSS) { a.Styles = append(a.Styles, css) }
func (a *ShellAssets) AddScript(src string)      { a.Scripts = append(a.Scripts, src) }
func (a *ShellAssets) AddJavaScript(js template.JS) {
	a.JavaScripts = append(a.JavaScripts, js)
}
