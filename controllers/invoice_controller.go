package controllers

import (
	"autoworkshop/dto"
	"autoworkshop/services"
	"autoworkshop/utils"

	"github.com/gin-gonic/gin"
)

type InvoiceController struct {
	svc         services.InvoiceService
	settingsSvc services.SettingsService
}

func NewInvoiceController(svc services.InvoiceService, settingsSvc services.SettingsService) *InvoiceController {
	return &InvoiceController{svc, settingsSvc}
}

func (ctrl *InvoiceController) List(c *gin.Context) {
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

func (ctrl *InvoiceController) GetByID(c *gin.Context) {
	inv, err := ctrl.svc.GetByID(utils.ClientCompanyID(c), c.Param("id"))
	if err != nil {
		utils.ErrNotFound(c, "Invoice not found")
		return
	}
	utils.OK(c, inv)
}

func (ctrl *InvoiceController) GetByCustomer(c *gin.Context) {
	list, err := ctrl.svc.GetByCustomer(utils.ClientCompanyID(c), c.Param("id"))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, list)
}

func (ctrl *InvoiceController) GetByCar(c *gin.Context) {
	list, err := ctrl.svc.GetByCar(utils.ClientCompanyID(c), c.Param("id"))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, list)
}

func (ctrl *InvoiceController) Create(c *gin.Context) {
	var req dto.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	settings, err := ctrl.settingsSvc.Get(utils.ClientCompanyID(c))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	inv, err := ctrl.svc.Create(utils.ClientCompanyID(c), req, settings)
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.Created(c, inv)
}

func (ctrl *InvoiceController) Update(c *gin.Context) {
	var req dto.UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	inv, err := ctrl.svc.Update(utils.ClientCompanyID(c), c.Param("id"), req)
	if err != nil {
		utils.ErrNotFound(c, err.Error())
		return
	}
	utils.OK(c, inv)
}

func (ctrl *InvoiceController) Delete(c *gin.Context) {
	if err := ctrl.svc.Delete(utils.ClientCompanyID(c), c.Param("id")); err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, gin.H{"ok": true})
}

func (ctrl *InvoiceController) AddPayment(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	payment, inv, err := ctrl.svc.AddPayment(utils.ClientCompanyID(c), c.Param("id"), req)
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.Created(c, gin.H{"payment": payment, "invoice": inv})
}

func (ctrl *InvoiceController) DeletePayment(c *gin.Context) {
	inv, err := ctrl.svc.DeletePayment(utils.ClientCompanyID(c), c.Param("id"), c.Param("paymentId"))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, gin.H{"invoice": inv})
}
