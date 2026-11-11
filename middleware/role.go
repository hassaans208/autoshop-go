package middleware

import (
	"autoworkshop/utils"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		r, _ := role.(string)
		if !allowed[r] {
			utils.ErrForbidden(c)
			return
		}
		c.Next()
	}
}
