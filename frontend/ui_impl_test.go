package frontend

import (
	"strings"
	"sync"
	"testing"
)

func TestHeadingHelpersRenderExpectedTags(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want string
	}{
		{name: "h1", node: H1(Text("Title 1")), want: `<h1><span>Title 1</span></h1>`},
		{name: "h2", node: H2(Text("Title 2")), want: `<h2><span>Title 2</span></h2>`},
		{name: "h3", node: H3(Text("Title 3")), want: `<h3><span>Title 3</span></h3>`},
		{name: "h4", node: H4(Text("Title 4")), want: `<h4><span>Title 4</span></h4>`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := tt.node.Render()
			if err != nil {
				t.Fatalf("render failed: %v", err)
			}
			if got := string(html); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestFileUploadRendersFileInput(t *testing.T) {
	html, err := FileUpload("attachment", true).Render()
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	got := string(html)
	for _, want := range []string{`name="attachment"`, `type="file"`, `required`} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
}

func TestLinkRendersExpectedAttributes(t *testing.T) {
	tests := []struct {
		name    string
		node    Node
		want    []string
		notWant []string
	}{
		{
			name: "escapes label and href",
			node: Link(LinkProps{Label: `Docs <now>`, Href: `/docs?team=a&role=admin`}),
			want: []string{`class="link link-hover w-fit"`, `href="/docs?team=a&amp;role=admin"`, `<span>Docs &lt;now&gt;</span>`},
		},
		{
			name: "external defaults target and rel",
			node: ExternalLink("Docs", "https://example.com/docs", ComponentProps{}),
			want: []string{`href="https://example.com/docs"`, `target="_blank"`, `rel="noopener noreferrer"`},
		},
		{
			name: "external icon link",
			node: ExternalIconLink("↗", "Open docs", "https://example.com/docs", ComponentProps{Variant: "ghost", Size: "sm"}),
			want: []string{`class="btn w-fit btn-ghost btn-sm btn-square"`, `href="https://example.com/docs"`, `target="_blank"`, `rel="noopener noreferrer"`, `aria-label="Open docs"`, `<span class="ui-icon material-icons" aria-hidden="true">↗</span>`},
		},
		{
			name: "download filename implies download",
			node: DownloadLink("CSV", "/assets/users.csv", "users report.csv", ComponentProps{Variant: "primary", Size: "sm"}),
			want: []string{`class="btn w-fit btn-primary btn-sm"`, `href="/assets/users.csv"`, `download="users report.csv"`, `<span>CSV</span>`},
		},
		{
			name:    "disabled link is inert",
			node:    Link(LinkProps{Label: "Disabled", Href: "/danger", Props: ComponentProps{Disabled: true}}),
			want:    []string{`href="#"`, `aria-disabled="true"`, `tabindex="-1"`, `pointer-events-none`, `cursor-not-allowed`, `opacity-50`},
			notWant: []string{`href="/danger"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := tt.node.Render()
			if err != nil {
				t.Fatalf("render failed: %v", err)
			}
			got := string(html)
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Fatalf("expected %q in %q", want, got)
				}
			}
			for _, notWant := range tt.notWant {
				if strings.Contains(got, notWant) {
					t.Fatalf("did not expect %q in %q", notWant, got)
				}
			}
		})
	}
}

func TestLoadComponentTemplatesCachesParsedTemplates(t *testing.T) {
	cachedTemplates = nil
	cachedTemplatesErr = nil
	componentTemplatesOnce = sync.Once{}

	first, err := loadComponentTemplates()
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}
	second, err := loadComponentTemplates()
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}
	if first == nil || second == nil {
		t.Fatal("expected non-nil template sets")
	}
	if first != second {
		t.Fatalf("expected cached template pointer reuse, got %p and %p", first, second)
	}
}
