package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateWritesHTMLFromMarkdown(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "index.md")
	if err := os.WriteFile(path, []byte("# Docs Home\n\nSee [Components](components/index.md).\n\n## Overview\n\nWelcome to Marionette."), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := generate(path); err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	body := readFile(t, filepath.Join(dir, "index.html"))
	for _, want := range []string{
		"<title>Docs Home</title>",
		"<h1 id=\"docs-home\">Docs Home</h1>",
		`href="components/index.html"`,
		"Welcome to Marionette.",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected generated docs page to contain %q, got %q", want, body)
		}
	}
}

func TestGenerateComponentsWritesIndexAndComponentPages(t *testing.T) {
	dir := t.TempDir()
	componentsDir := filepath.Join(dir, "components")
	if err := os.MkdirAll(componentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	indexPath := filepath.Join(componentsDir, "index.md")
	catalogPath := filepath.Join(componentsDir, "_index.json")

	indexMarkdown := `# Components

## Contents

- [Button](button.md)

## Button

Button overview from docs.

- Golden sample: [button.golden.html](../../testdata/golden/button.golden.html)
`
	if err := os.WriteFile(indexPath, []byte(indexMarkdown), 0o644); err != nil {
		t.Fatal(err)
	}
	catalogJSON := `{"components":[{"id":"button","name":"Button","group":"Components","description":"Trigger an action.","golden":"testdata/golden/button.golden.html","example":"docs/site/components/button.html","template":"templates/components/button.tmpl"}]}`
	if err := os.WriteFile(catalogPath, []byte(catalogJSON), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := generateComponents(indexPath, catalogPath); err != nil {
		t.Fatalf("generateComponents failed: %v", err)
	}

	indexHTML := readFile(t, filepath.Join(componentsDir, "index.html"))
	for _, want := range []string{"Components Gallery", `href="./button.html"`, "Trigger an action."} {
		if !strings.Contains(indexHTML, want) {
			t.Fatalf("expected component index to contain %q, got %q", want, indexHTML)
		}
	}

	buttonHTML := readFile(t, filepath.Join(componentsDir, "button.html"))
	for _, want := range []string{
		"<title>Button</title>",
		`aria-current="page"`,
		"Button overview from docs.",
		"Go usage",
		"saveButton := mf.Button",
		`<iframe src="./button.html" title="Button example"`,
	} {
		if !strings.Contains(buttonHTML, want) {
			t.Fatalf("expected button page to contain %q, got %q", want, buttonHTML)
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(body)
}
