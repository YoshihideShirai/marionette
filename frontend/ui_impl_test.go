package frontend

import "testing"

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
