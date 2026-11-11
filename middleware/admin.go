package middleware

import (
	"autoworkshop/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AdminGate(adminAPIKey, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if adminAPIKey == "" && jwtSecret == "" {
			utils.ErrForbidden(c)
			return
		}

		header := c.GetHeader("Authorization")
		if strings.HasPrefix(header, "Bearer ") {
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			if token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			}); err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if role, _ := claims["role"].(string); role == "admin" {
						c.Set("adminEmail", claims["email"])
						c.Next()
						return
					}
				}
			}
			if adminAPIKey != "" && tokenStr == adminAPIKey {
				c.Next()
				return
			}
		}

		if adminAPIKey != "" {
			key := c.GetHeader("X-Admin-Key")
			if key == adminAPIKey {
				c.Next()
				return
			}
		}

		utils.ErrUnauthorized(c, "Invalid admin credentials")
	}
}
