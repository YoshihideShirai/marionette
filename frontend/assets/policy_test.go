package assets

import (
	"strings"
	"testing"
)

func TestIsExternalURL(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "https", value: "https://cdn.example.com/app.js", want: true},
		{name: "http with surrounding whitespace", value: "  HTTP://cdn.example.com/app.js\t", want: true},
		{name: "protocol-relative is not absolute network URL", value: "//cdn.example.com/app.js", want: false},
		{name: "local absolute path", value: "/assets/app.js", want: false},
		{name: "relative path", value: "assets/app.js", want: false},
		{name: "empty", value: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExternalURL(tt.value); got != tt.want {
				t.Fatalf("IsExternalURL(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name       string
		policy     AssetPolicy
		kind       string
		value      string
		wantErr    bool
		wantSubstr string
	}{
		{
			name:   "allows external URL by default",
			policy: AssetPolicy{},
			kind:   "script",
			value:  "https://cdn.example.com/app.js",
		},
		{
			name:   "allows local URL when external URLs forbidden",
			policy: AssetPolicy{ForbidExternalURLs: true},
			kind:   "script",
			value:  "/assets/app.js",
		},
		{
			name:       "rejects external URL when forbidden",
			policy:     AssetPolicy{ForbidExternalURLs: true},
			kind:       "script",
			value:      "https://cdn.example.com/app.js",
			wantErr:    true,
			wantSubstr: "external URL for script",
		},
		{
			name:       "uses fallback kind in error",
			policy:     AssetPolicy{ForbidExternalURLs: true},
			kind:       "  ",
			value:      "http://cdn.example.com/app.css",
			wantErr:    true,
			wantSubstr: "external URL for asset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.policy, tt.kind, tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ValidateURL() error = nil, want error")
				}
				if !strings.Contains(err.Error(), tt.wantSubstr) {
					t.Fatalf("ValidateURL() error = %q, want substring %q", err.Error(), tt.wantSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidateURL() error = %v, want nil", err)
			}
		})
	}
}
