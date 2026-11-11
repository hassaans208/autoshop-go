package controllers

import (
	"autoworkshop/dto"
	"autoworkshop/services"
	"autoworkshop/utils"

	"github.com/gin-gonic/gin"
)

type CustomerController struct{ svc services.CustomerService }

func NewCustomerController(svc services.CustomerService) *CustomerController {
	return &CustomerController{svc}
}

func (ctrl *CustomerController) List(c *gin.Context) {
	p := utils.ParsePage(c)
	list, total, err := ctrl.svc.List(utils.ClientCompanyID(c), p)
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, utils.Paged(list, total, p))
}

func (ctrl *CustomerController) GetByID(c *gin.Context) {
	cust, err := ctrl.svc.GetByID(utils.ClientCompanyID(c), c.Param("id"))
	if err != nil {
		utils.ErrNotFound(c, "Customer not found")
		return
	}
	utils.OK(c, cust)
}

func (ctrl *CustomerController) GetByPhone(c *gin.Context) {
	cust, err := ctrl.svc.GetByPhone(utils.ClientCompanyID(c), c.Param("phone"))
	if err != nil {
		utils.ErrNotFound(c, "Customer not found")
		return
	}
	utils.OK(c, cust)
}

func (ctrl *CustomerController) Create(c *gin.Context) {
	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	cust, err := ctrl.svc.Create(utils.ClientCompanyID(c), req)
	if err != nil {
		utils.ErrConflict(c, err.Error())
		return
	}
	utils.Created(c, cust)
}

func (ctrl *CustomerController) Update(c *gin.Context) {
	var req dto.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	cust, err := ctrl.svc.Update(utils.ClientCompanyID(c), c.Param("id"), req)
	if err != nil {
		utils.ErrNotFound(c, err.Error())
		return
	}
	utils.OK(c, cust)
}

func (ctrl *CustomerController) Delete(c *gin.Context) {
	if err := ctrl.svc.Delete(utils.ClientCompanyID(c), c.Param("id")); err != nil {
		utils.ErrConflict(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"ok": true})
}
