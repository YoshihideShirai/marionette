package daisyui

import (
	"strings"
	"testing"
)

func TestToggleVariantRendersDaisyUIColorClasses(t *testing.T) {
	tests := []struct {
		name    string
		variant string
		want    string
	}{
		{name: "primary", variant: "primary", want: `class="toggle toggle-primary"`},
		{name: "secondary", variant: "secondary", want: `class="toggle toggle-secondary"`},
		{name: "accent", variant: "accent", want: `class="toggle toggle-accent"`},
		{name: "neutral", variant: "neutral", want: `class="toggle toggle-neutral"`},
		{name: "info", variant: "info", want: `class="toggle toggle-info"`},
		{name: "success", variant: "success", want: `class="toggle toggle-success"`},
		{name: "warning", variant: "warning", want: `class="toggle toggle-warning"`},
		{name: "error", variant: "error", want: `class="toggle toggle-error"`},
		{name: "danger alias", variant: "danger", want: `class="toggle toggle-error"`},
		{name: "normalized", variant: " Primary ", want: `class="toggle toggle-primary"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := ToggleVariant("demo", true, tt.variant).Render()
			if err != nil {
				t.Fatalf("render ToggleVariant: %v", err)
			}
			got := string(html)
			if !strings.Contains(got, tt.want) {
				t.Fatalf("expected %q in %q", tt.want, got)
			}
			if !strings.Contains(got, `checked="checked"`) {
				t.Fatalf("expected checked attribute in %q", got)
			}
		})
	}
}
