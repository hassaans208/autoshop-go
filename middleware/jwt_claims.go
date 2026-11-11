package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
)

// stringFromJWTClaim coerces a MapClaims value to string. JSON unmarshaling or
// library quirks can leave values as json.Number or other types; a plain
// type assertion is fragile and causes "Missing tenant in token" false positives.
func stringFromJWTClaim(v any) string {
	if v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return strings.TrimSpace(s)
	case json.Number:
		return strings.TrimSpace(s.String())
	case float64:
		return strings.TrimSpace(fmt.Sprintf("%.0f", s))
	default:
		return strings.TrimSpace(fmt.Sprint(s))
	}
}
