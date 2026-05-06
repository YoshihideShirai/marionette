package componenttmpl

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// Source owns component template path resolution, parsing, and caching.
// Callers keep one Source per package-level template root and render templates
// through Node so the template loader remains centralized in this package.
type Source struct {
	dir string

	once sync.Once
	tmpl *template.Template
	err  error
}

// NewSource returns a cached component template source rooted at dir.
func NewSource(dir string) *Source {
	return &Source{dir: dir}
}

// NewSourceFromCaller returns a cached template source by resolving pathParts
// relative to the file at runtime.Caller(skip).
func NewSourceFromCaller(skip int, pathParts ...string) (*Source, error) {
	_, currentFile, _, ok := runtime.Caller(skip)
	if !ok {
		return nil, fmt.Errorf("failed to resolve component template path for %s", filepath.Join(pathParts...))
	}
	parts := append([]string{filepath.Dir(currentFile)}, pathParts...)
	return NewSource(filepath.Join(parts...)), nil
}

// Load parses component templates from the source directory once and returns
// the cached parsed template set.
func (s *Source) Load() (*template.Template, error) {
	if s == nil {
		return nil, fmt.Errorf("component template source is nil")
	}
	s.once.Do(func() {
		s.tmpl, s.err = Load(s.dir)
	})
	return s.tmpl, s.err
}

// Node executes a named template from Source.
type Node struct {
	Source *Source
	Name   string
	Data   any
}

// Render executes the node template and returns rendered HTML.
func (n Node) Render() (template.HTML, error) {
	tmpl, err := n.Source.Load()
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	if err := tmpl.ExecuteTemplate(&out, n.Name, n.Data); err != nil {
		return "", err
	}
	return template.HTML(out.String()), nil
}

// Load parses component templates from dir and returns parsed templates.
// Template names are expected to be referenced as "components/<basename>".
func Load(dir string) (*template.Template, error) {
	patterns := []string{"*.tmpl", "*.html"}
	files := make([]string, 0, 16)
	checked := make([]string, 0, len(patterns))

	for _, pattern := range patterns {
		glob := filepath.Join(dir, pattern)
		checked = append(checked, glob)
		matched, err := filepath.Glob(glob)
		if err != nil {
			return nil, fmt.Errorf("failed to glob component templates (%s): %w", glob, err)
		}
		files = append(files, matched...)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no component templates found; checked paths: %s", strings.Join(checked, ", "))
	}

	sort.Strings(files)
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse component templates from %s (template names must follow components/<basename>): %w", dir, err)
	}
	return tmpl, nil
}
