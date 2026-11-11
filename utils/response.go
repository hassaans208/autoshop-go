package utils

import "github.com/gin-gonic/gin"

func OK(c *gin.Context, data any) {
	c.JSON(200, data)
}

func Created(c *gin.Context, data any) {
	c.JSON(201, data)
}

func ErrBadRequest(c *gin.Context, msg string, details ...any) {
	body := gin.H{"error": msg, "code": "BAD_REQUEST"}
	if len(details) > 0 {
		body["details"] = details[0]
	}
	c.AbortWithStatusJSON(400, body)
}

func ErrUnauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(401, gin.H{"error": msg, "code": "UNAUTHORIZED"})
}

func ErrForbidden(c *gin.Context) {
	c.AbortWithStatusJSON(403, gin.H{"error": "Insufficient permissions", "code": "FORBIDDEN"})
}

func ErrNotFound(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(404, gin.H{"error": msg, "code": "NOT_FOUND"})
}

func ErrConflict(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(409, gin.H{"error": msg, "code": "CONFLICT"})
}

func ErrUnprocessable(c *gin.Context, msg string, details ...any) {
	body := gin.H{"error": msg, "code": "UNPROCESSABLE"}
	if len(details) > 0 {
		body["details"] = details[0]
	}
	c.AbortWithStatusJSON(422, body)
}

func ErrInternal(c *gin.Context, err error) {
	c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error", "code": "INTERNAL", "details": err.Error()})
}
