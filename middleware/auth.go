package middleware

import (
	"autoworkshop/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthGate(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			utils.ErrUnauthorized(c, "Missing or malformed token")
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			utils.ErrUnauthorized(c, "Invalid or expired token")
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		c.Set("staffID", stringFromJWTClaim(claims["sub"]))
		c.Set("role", stringFromJWTClaim(claims["role"]))
		c.Set("email", stringFromJWTClaim(claims["email"]))
		ccID := stringFromJWTClaim(claims["clientCompanyId"])
		if ccID == "" {
			ccID = stringFromJWTClaim(claims["client_company_id"])
		}
		if ccID == "" {
			utils.ErrUnauthorized(c, "Missing tenant in token")
			return
		}
		c.Set("clientCompanyID", ccID)
		c.Next()
	}
}
