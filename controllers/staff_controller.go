package controllers

import (
	"autoworkshop/dto"
	"autoworkshop/services"
	"autoworkshop/utils"

	"github.com/gin-gonic/gin"
)

type StaffController struct{ svc services.StaffService }

func NewStaffController(svc services.StaffService) *StaffController { return &StaffController{svc} }

func (ctrl *StaffController) List(c *gin.Context) {
	list, err := ctrl.svc.List(utils.ClientCompanyID(c))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, gin.H{"data": list})
}

func (ctrl *StaffController) Create(c *gin.Context) {
	var req dto.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	staff, err := ctrl.svc.Create(utils.ClientCompanyID(c), req)
	if err != nil {
		utils.ErrConflict(c, err.Error())
		return
	}
	utils.Created(c, staff)
}

func (ctrl *StaffController) Update(c *gin.Context) {
	var req dto.UpdateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	staff, err := ctrl.svc.Update(utils.ClientCompanyID(c), c.Param("id"), req)
	if err != nil {
		utils.ErrNotFound(c, err.Error())
		return
	}
	utils.OK(c, staff)
}

func (ctrl *StaffController) Delete(c *gin.Context) {
	if err := ctrl.svc.Delete(utils.ClientCompanyID(c), c.Param("id")); err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, gin.H{"ok": true})
}
