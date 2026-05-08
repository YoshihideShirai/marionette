package assets

import (
	"fmt"
	"strings"
)

// AssetPolicy controls which asset URLs Marionette may emit into generated shells.
type AssetPolicy struct {
	// ForbidExternalURLs rejects http:// and https:// asset URLs so pages can be audited for offline use.
	ForbidExternalURLs bool
}

// IsExternalURL reports whether value is an absolute network URL that can require online access.
func IsExternalURL(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")
}

// ValidateURL returns an error when value violates policy.
func ValidateURL(policy AssetPolicy, kind, value string) error {
	if !policy.ForbidExternalURLs || !IsExternalURL(value) {
		return nil
	}
	if kind = strings.TrimSpace(kind); kind == "" {
		kind = "asset"
	}
	return fmt.Errorf("asset policy forbids external URL for %s: %s", kind, value)
}
