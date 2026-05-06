package twailwindcss

import "testing"

func TestFrameworkScriptsReturnsTailwindBrowserURL(t *testing.T) {
	scripts := FrameworkScripts()
	if len(scripts) != 1 {
		t.Fatalf("expected one TailwindCSS browser script, got %d", len(scripts))
	}
	if scripts[0] != BrowserURL {
		t.Fatalf("expected BrowserURL %q, got %q", BrowserURL, scripts[0])
	}
}

func TestFrameworkStylesheetsIsEmpty(t *testing.T) {
	if stylesheets := FrameworkStylesheets(); len(stylesheets) != 0 {
		t.Fatalf("expected no TailwindCSS stylesheets, got %v", stylesheets)
	}
}
