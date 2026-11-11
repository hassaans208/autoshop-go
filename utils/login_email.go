package utils

import (
	"strings"
	"unicode"
)

const defaultLoginEmailDomain = "@avinertech.com"

// NormalizeLoginEmail uses the value as-is when it looks like an email; otherwise
// keeps letters, digits, '.', '-', '_' and replaces other runes with '_', then appends @avinertech.com.
func NormalizeLoginEmail(username string) string {
	u := strings.TrimSpace(username)
	if u == "" {
		return ""
	}
	if looksLikeEmail(u) {
		return strings.ToLower(u)
	}
	var b strings.Builder
	for _, r := range u {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '-' || r == '_':
			b.WriteRune(unicode.ToLower(r))
		default:
			b.WriteRune('_')
		}
	}
	s := strings.Trim(b.String(), "._-")
	if s == "" {
		return ""
	}
	return s + defaultLoginEmailDomain
}

func looksLikeEmail(s string) bool {
	at := strings.LastIndex(s, "@")
	if at <= 0 || at >= len(s)-1 {
		return false
	}
	local := s[:at]
	domain := s[at+1:]
	if local == "" || domain == "" {
		return false
	}
	return strings.Contains(domain, ".")
}
