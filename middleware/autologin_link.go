package middleware

import (
	"crypto/subtle"

	"github.com/gin-gonic/gin"
)

// SecretPathGate allows the handler only when linkSecret matches the configured value.
// Wrong or missing configuration yields 404 so the endpoint stays unobvious.
func SecretPathGate(expected string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if expected == "" {
			c.AbortWithStatus(404)
			return
		}
		got := c.Param("linkSecret")
		if subtle.ConstantTimeCompare([]byte(got), []byte(expected)) != 1 {
			c.AbortWithStatus(404)
			return
		}
		c.Next()
	}
}

// AutologinLinkGate is an alias for SecretPathGate (autologin-link routes).
func AutologinLinkGate(expected string) gin.HandlerFunc { return SecretPathGate(expected) }
