package components

import (
	"strings"
	"testing"

	lowhtml "github.com/YoshihideShirai/marionette/frontend/html"
	shared "github.com/YoshihideShirai/marionette/frontend/shared"
)

func TestGridRendersDefaultClassesAndChildren(t *testing.T) {
	rendered, err := Grid(shared.GridProps{}, lowhtml.Text("content")).Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	got := string(rendered)
	wantClass := `class="grid gap-4 grid-cols-1 md:grid-cols-2 xl:grid-cols-3"`
	if !strings.Contains(got, wantClass) {
		t.Fatalf("rendered Grid class missing %q in %q", wantClass, got)
	}
	if !strings.Contains(got, `<span>content</span>`) {
		t.Fatalf("rendered Grid child text missing in %q", got)
	}
}

func TestGridGapClassVariants(t *testing.T) {
	tests := []struct {
		name string
		gap  string
		want string
	}{
		{name: "none", gap: "none", want: "gap-0"},
		{name: "zero", gap: "0", want: "gap-0"},
		{name: "xs", gap: "xs", want: "gap-1"},
		{name: "sm", gap: "sm", want: "gap-2"},
		{name: "lg", gap: "lg", want: "gap-6"},
		{name: "xl", gap: "xl", want: "gap-8"},
		{name: "unknown", gap: "huge", want: "gap-4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered, err := Grid(shared.GridProps{Gap: tt.gap}, lowhtml.Text("content")).Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			got := string(rendered)
			wantClass := `class="grid ` + tt.want + ` grid-cols-1 md:grid-cols-2 xl:grid-cols-3"`
			if !strings.Contains(got, wantClass) {
				t.Fatalf("rendered Grid with Gap %q missing class %q in %q", tt.gap, wantClass, got)
			}
		})
	}
}

func TestGridColumnsClassVariants(t *testing.T) {
	tests := []struct {
		name    string
		columns string
		want    string
	}{
		{name: "one", columns: "1", want: "grid-cols-1"},
		{name: "two", columns: "2", want: "grid-cols-1 md:grid-cols-2"},
		{name: "four", columns: "4", want: "grid-cols-1 sm:grid-cols-2 xl:grid-cols-4"},
		{name: "unknown", columns: "many", want: "grid-cols-1 md:grid-cols-2 xl:grid-cols-3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered, err := Grid(shared.GridProps{Columns: tt.columns}, lowhtml.Text("content")).Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			got := string(rendered)
			wantClass := `class="grid gap-4 ` + tt.want + `"`
			if !strings.Contains(got, wantClass) {
				t.Fatalf("rendered Grid with Columns %q missing class %q in %q", tt.columns, wantClass, got)
			}
		})
	}
}

func TestGridMinColumnWidthTakesPrecedence(t *testing.T) {
	tests := []struct {
		name           string
		minColumnWidth string
		want           string
	}{
		{name: "sm", minColumnWidth: "sm", want: "grid-cols-[repeat(auto-fit,minmax(14rem,1fr))]"},
		{name: "md", minColumnWidth: "md", want: "grid-cols-[repeat(auto-fit,minmax(18rem,1fr))]"},
		{name: "lg", minColumnWidth: "lg", want: "grid-cols-[repeat(auto-fit,minmax(22rem,1fr))]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered, err := Grid(shared.GridProps{Columns: "4", MinColumnWidth: tt.minColumnWidth}, lowhtml.Text("content")).Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			got := string(rendered)
			wantClass := `class="grid gap-4 ` + tt.want + `"`
			if !strings.Contains(got, wantClass) {
				t.Fatalf("rendered Grid with MinColumnWidth %q missing class %q in %q", tt.minColumnWidth, wantClass, got)
			}
			if strings.Contains(got, "xl:grid-cols-4") {
				t.Fatalf("rendered Grid with MinColumnWidth %q should not include Columns-derived class in %q", tt.minColumnWidth, got)
			}
		})
	}
}

func TestGridRenderSkipsNilChildren(t *testing.T) {
	rendered, err := Grid(shared.GridProps{}, lowhtml.Text("before"), nil, lowhtml.Text("after")).Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	got := string(rendered)
	for _, want := range []string{`class="grid gap-4 grid-cols-1 md:grid-cols-2 xl:grid-cols-3"`, `<span>before</span>`, `<span>after</span>`} {
		if !strings.Contains(got, want) {
			t.Fatalf("rendered Grid missing %q in %q", want, got)
		}
	}
}

func TestRegionRendersIDClassAndChildren(t *testing.T) {
	rendered, err := Region(
		shared.RegionProps{ID: " main ", Props: shared.ComponentProps{Class: " p-4 "}},
		lowhtml.Text("content"),
	).Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := `<div class="p-4" id="main"><span>content</span></div>`
	if got := string(rendered); got != want {
		t.Fatalf("Render() = %q, want current trimmed id/class behavior %q", got, want)
	}
}
