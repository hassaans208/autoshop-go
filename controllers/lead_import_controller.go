package controllers

import (
	"autoworkshop/dto"
	"autoworkshop/services"
	"autoworkshop/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type LeadImportController struct{ svc services.LeadImportService }

func NewLeadImportController(svc services.LeadImportService) *LeadImportController {
	return &LeadImportController{svc}
}

func (ctrl *LeadImportController) Import(c *gin.Context) {
	var rows []dto.LeadImportRow
	var parseErr error

	if fh, err := c.FormFile("file"); err == nil && fh != nil {
		f, err := fh.Open()
		if err != nil {
			utils.ErrBadRequest(c, "Could not open uploaded file")
			return
		}
		defer f.Close()
		rows, parseErr = services.ParseLeadsCSV(f)
	} else if strings.Contains(c.GetHeader("Content-Type"), "application/json") {
		var req dto.LeadImportRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ErrBadRequest(c, "Validation failed", err.Error())
			return
		}
		rows = req.Rows
	} else {
		utils.ErrBadRequest(c, "Upload CSV as multipart field \"file\" (same column headers as your sheet), or POST JSON {\"rows\": [...]}")
		return
	}

	if parseErr != nil {
		utils.ErrBadRequest(c, parseErr.Error())
		return
	}

	resp, err := ctrl.svc.Import(rows)
	if err != nil {
		if err.Error() == "FRONTEND_URL is not configured" {
			utils.ErrUnprocessable(c, err.Error())
			return
		}
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, resp)
}
