package utils

import "github.com/gin-gonic/gin"

const (
	ctxClientCompanyID = "clientCompanyID"
	ctxStaffID         = "staffID"
)

func ClientCompanyID(c *gin.Context) string {
	v, _ := c.Get(ctxClientCompanyID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func StaffID(c *gin.Context) string {
	v, _ := c.Get(ctxStaffID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
