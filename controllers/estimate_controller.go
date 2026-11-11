package controllers

import (
	"autoworkshop/dto"
	"autoworkshop/services"
	"autoworkshop/utils"

	"github.com/gin-gonic/gin"
)

type EstimateController struct {
	svc        services.EstimateService
	settingsSvc services.SettingsService
}

func NewEstimateController(svc services.EstimateService, settingsSvc services.SettingsService) *EstimateController {
	return &EstimateController{svc, settingsSvc}
}

func (ctrl *EstimateController) List(c *gin.Context) {
	p := utils.ParsePage(c)
	filters := map[string]string{
		"status": c.Query("status"), "customerId": c.Query("customerId"),
		"carId": c.Query("carId"), "from": c.Query("from"), "to": c.Query("to"),
	}
	list, total, err := ctrl.svc.List(utils.ClientCompanyID(c), p, filters)
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, utils.Paged(list, total, p))
}

func (ctrl *EstimateController) GetByID(c *gin.Context) {
	e, err := ctrl.svc.GetByID(utils.ClientCompanyID(c), c.Param("id"))
	if err != nil {
		utils.ErrNotFound(c, "Estimate not found")
		return
	}
	utils.OK(c, e)
}

func (ctrl *EstimateController) GetByCar(c *gin.Context) {
	list, err := ctrl.svc.GetByCar(utils.ClientCompanyID(c), c.Param("id"))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, list)
}

func (ctrl *EstimateController) Create(c *gin.Context) {
	var req dto.CreateEstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	settings, err := ctrl.settingsSvc.Get(utils.ClientCompanyID(c))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	e, err := ctrl.svc.Create(utils.ClientCompanyID(c), req, settings)
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.Created(c, e)
}

func (ctrl *EstimateController) Update(c *gin.Context) {
	var req dto.UpdateEstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	e, err := ctrl.svc.Update(utils.ClientCompanyID(c), c.Param("id"), req)
	if err != nil {
		utils.ErrNotFound(c, err.Error())
		return
	}
	utils.OK(c, e)
}

func (ctrl *EstimateController) Delete(c *gin.Context) {
	if err := ctrl.svc.Delete(utils.ClientCompanyID(c), c.Param("id")); err != nil {
		utils.ErrConflict(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"ok": true})
}

func (ctrl *EstimateController) Convert(c *gin.Context) {
	var req dto.ConvertEstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	settings, err := ctrl.settingsSvc.Get(utils.ClientCompanyID(c))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	inv, err := ctrl.svc.Convert(utils.ClientCompanyID(c), c.Param("id"), req, settings)
	if err != nil {
		utils.ErrConflict(c, err.Error())
		return
	}
	utils.Created(c, inv)
}
