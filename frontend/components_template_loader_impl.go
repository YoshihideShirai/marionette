package frontend

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/YoshihideShirai/marionette/internal/componenttmpl"
)

// このファイルはコンポーネントテンプレートの読み込みとキャッシュを管理する。
// テンプレート実行ノードの共通処理もここに置く。

type templateNode struct {
	name string
	data any
}

var (
	cachedTemplates        *template.Template
	cachedTemplatesErr     error
	componentTemplatesOnce sync.Once
)

func (n templateNode) Render() (template.HTML, error) {
	tmpl, err := loadComponentTemplates()
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	if err := tmpl.ExecuteTemplate(&out, n.name, n.data); err != nil {
		return "", err
	}
	return template.HTML(out.String()), nil
}

func loadComponentTemplates() (*template.Template, error) {
	componentTemplatesOnce.Do(func() {
		_, currentFile, _, ok := runtime.Caller(0)
		if !ok {
			cachedTemplatesErr = fmt.Errorf("failed to resolve component template path for %s", "templates/components")
			return
		}
		componentsDir := filepath.Join(filepath.Dir(currentFile), "..", "templates", "components")
		cachedTemplates, cachedTemplatesErr = componenttmpl.Load(componentsDir)
	})
	return cachedTemplates, cachedTemplatesErr
}
