package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var pageTemplate = template.Must(template.New("docs").Parse(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>{{.Title}}</title>
    <style>
      :root { color-scheme: light dark; }
      body {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
        margin: 2rem auto;
        max-width: 880px;
        line-height: 1.65;
        padding: 0 1rem;
      }
      a { color: #2563eb; }
      pre {
        background: color-mix(in oklab, canvas 92%, black 8%);
        border-radius: 10px;
        overflow-x: auto;
        padding: 0.75rem;
      }
      code { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
      iframe { background: white; }
    </style>
  </head>
  <body>
    {{.Body}}
  </body>
</html>
`))

var mdHrefRe = regexp.MustCompile(`href="([^"]+?)\.md"`)

func main() {
	targets := []string{
		"docs/site/index.md",
		"docs/site/components/index.md",
	}

	for _, path := range targets {
		if err := generate(path); err != nil {
			fmt.Fprintf(os.Stderr, "generate %s: %v\n", path, err)
			os.Exit(1)
		}
	}
}

func generate(path string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	var body bytes.Buffer
	if err := md.Convert(src, &body); err != nil {
		return err
	}

	renderedBody := template.HTML(rewriteMarkdownLinks(body.String()))
	title := docTitle(src, filepath.Base(path))

	var out bytes.Buffer
	if err := pageTemplate.Execute(&out, map[string]any{
		"Title": title,
		"Body":  renderedBody,
	}); err != nil {
		return err
	}

	htmlPath := path[:len(path)-len(filepath.Ext(path))] + ".html"
	return os.WriteFile(htmlPath, out.Bytes(), 0o644)
}

func rewriteMarkdownLinks(s string) string {
	return mdHrefRe.ReplaceAllString(s, `href="$1.html"`)
}

func docTitle(src []byte, fallback string) string {
	for _, line := range bytes.Split(src, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("# ")) {
			return string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("# "))))
		}
	}
	return fallback
}
