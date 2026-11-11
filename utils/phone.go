package utils

import "regexp"

var nonDigit = regexp.MustCompile(`\D`)

func NormalizePhone(phone string) string {
	return nonDigit.ReplaceAllString(phone, "")
}
