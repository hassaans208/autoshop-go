package controllers

import (
	"autoworkshop/dto"
	"autoworkshop/services"
	"autoworkshop/utils"

	"github.com/gin-gonic/gin"
)

type VehicleController struct{ svc services.VehicleService }

func NewVehicleController(svc services.VehicleService) *VehicleController {
	return &VehicleController{svc}
}

func (ctrl *VehicleController) List(c *gin.Context) {
	p := utils.ParsePage(c)
	list, total, err := ctrl.svc.List(utils.ClientCompanyID(c), p, c.Query("customerId"))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, utils.Paged(list, total, p))
}

func (ctrl *VehicleController) GetByID(c *gin.Context) {
	car, err := ctrl.svc.GetByID(utils.ClientCompanyID(c), c.Param("id"))
	if err != nil {
		utils.ErrNotFound(c, "Car not found")
		return
	}
	utils.OK(c, car)
}

func (ctrl *VehicleController) GetByCustomer(c *gin.Context) {
	cars, err := ctrl.svc.GetByCustomer(utils.ClientCompanyID(c), c.Param("id"))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, cars)
}

func (ctrl *VehicleController) Create(c *gin.Context) {
	var req dto.CreateCarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	car, err := ctrl.svc.Create(utils.ClientCompanyID(c), req)
	if err != nil {
		utils.ErrConflict(c, err.Error())
		return
	}
	utils.Created(c, car)
}

func (ctrl *VehicleController) Update(c *gin.Context) {
	var req dto.UpdateCarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	car, err := ctrl.svc.Update(utils.ClientCompanyID(c), c.Param("id"), req)
	if err != nil {
		utils.ErrNotFound(c, err.Error())
		return
	}
	utils.OK(c, car)
}

func (ctrl *VehicleController) Delete(c *gin.Context) {
	if err := ctrl.svc.Delete(utils.ClientCompanyID(c), c.Param("id")); err != nil {
		utils.ErrConflict(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"ok": true})
}
