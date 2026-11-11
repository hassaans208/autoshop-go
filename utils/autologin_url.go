package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func BuildAutologinURL(frontendBase, email, passwordHash string) string {
	base := strings.TrimRight(strings.TrimSpace(frontendBase), "/")
	if base == "" || email == "" || passwordHash == "" {
		return ""
	}
	return fmt.Sprintf("%s/autologin/%s?PasswordHash=%s",
		base,
		url.PathEscape(email),
		url.QueryEscape(passwordHash),
	)
}
