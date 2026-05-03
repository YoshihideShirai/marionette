package componenttmpl

import (
	"fmt"
	"html/template"
	"path/filepath"
	"sort"
	"strings"
)

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
