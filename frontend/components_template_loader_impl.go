package frontend

import (
	"html/template"
	"sync"

	"github.com/YoshihideShirai/marionette/internal/componenttmpl"
)

// This file keeps frontend wired to the shared component template renderer.
// Template parsing, execution, and cache ownership live in internal/componenttmpl;
// this package only resolves the repository template root.

type templateNode struct {
	name string
	data any
}

var (
	componentTemplateSource     *componenttmpl.Source
	componentTemplateSourceErr  error
	componentTemplateSourceOnce sync.Once
)

func (n templateNode) Render() (template.HTML, error) {
	source := componentTemplateSourceForPackage()
	if componentTemplateSourceErr != nil {
		return "", componentTemplateSourceErr
	}
	return componenttmpl.Node{
		Source: source,
		Name:   n.name,
		Data:   n.data,
	}.Render()
}

func componentTemplateSourceForPackage() *componenttmpl.Source {
	componentTemplateSourceOnce.Do(func() {
		componentTemplateSource, componentTemplateSourceErr = componenttmpl.NewSourceFromCaller(1, "..", "templates", "components")
	})
	return componentTemplateSource
}
