package utils

import (
	"strings"
	"unicode/utf8"
)

const MaxBusinessNameLen = 255

// TruncateBusinessName trims whitespace and limits length for DB name columns.
func TruncateBusinessName(s string) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) <= MaxBusinessNameLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:MaxBusinessNameLen-1]) + "…"
}
