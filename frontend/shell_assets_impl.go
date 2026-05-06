package frontend

import "html/template"

type ShellAssets struct {
	StyleTemplate        StyleTemplate
	FrameworkStylesheets []string
	FrameworkScripts     []string
	Stylesheets          []string
	Styles               []template.CSS
	Scripts              []string
	JavaScripts          []template.JS
}

func (a *ShellAssets) UseStyleTemplate(tpl StyleTemplate) {
	a.StyleTemplate = StyleTemplate{
		Name:                 tpl.Name,
		FrameworkStylesheets: append([]string(nil), tpl.FrameworkStylesheets...),
		FrameworkScripts:     append([]string(nil), tpl.FrameworkScripts...),
	}
	a.FrameworkStylesheets = nil
	a.FrameworkScripts = nil
}

func (a *ShellAssets) AddStylesheet(href string) { a.Stylesheets = append(a.Stylesheets, href) }
func (a *ShellAssets) AddStyle(css template.CSS) { a.Styles = append(a.Styles, css) }
func (a *ShellAssets) AddScript(src string)      { a.Scripts = append(a.Scripts, src) }
func (a *ShellAssets) AddJavaScript(js template.JS) {
	a.JavaScripts = append(a.JavaScripts, js)
}
