package controllers

import (
	"autoworkshop/dto"
	"autoworkshop/models"
	"autoworkshop/services"
	"autoworkshop/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SettingsController struct{ svc services.SettingsService }

func NewSettingsController(svc services.SettingsService) *SettingsController {
	return &SettingsController{svc}
}

func (ctrl *SettingsController) Get(c *gin.Context) {
	s, err := ctrl.svc.Get(utils.ClientCompanyID(c))
	if err != nil {
		utils.OK(c, models.Settings{EstPrefix: "EST", InvPrefix: "INV", Currency: "PKR"})
		return
	}
	utils.OK(c, s)
}

func (ctrl *SettingsController) Update(c *gin.Context) {
	var s models.Settings
	if err := c.ShouldBindJSON(&s); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	updated, err := ctrl.svc.Update(utils.ClientCompanyID(c), &s)
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, updated)
}

type CompanyController struct{ svc services.CompanyService }

func NewCompanyController(svc services.CompanyService) *CompanyController {
	return &CompanyController{svc}
}

func (ctrl *CompanyController) List(c *gin.Context) {
	names, err := ctrl.svc.List(utils.ClientCompanyID(c))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, names)
}

func (ctrl *CompanyController) Create(c *gin.Context) {
	var req dto.AddCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	company, err := ctrl.svc.Add(utils.ClientCompanyID(c), req)
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.Created(c, gin.H{"name": company.Name})
}

type JobTypeController struct{ svc services.JobTypeService }

func NewJobTypeController(svc services.JobTypeService) *JobTypeController {
	return &JobTypeController{svc}
}

func (ctrl *JobTypeController) List(c *gin.Context) {
	list, err := ctrl.svc.List(utils.ClientCompanyID(c))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, gin.H{"data": list})
}

func (ctrl *JobTypeController) Replace(c *gin.Context) {
	var req dto.ReplaceJobTypesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	if len(req.Entries()) == 0 {
		utils.ErrBadRequest(c, "Validation failed", "list or jobTypes is required")
		return
	}
	list, err := ctrl.svc.Replace(utils.ClientCompanyID(c), req)
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, gin.H{"data": list})
}

type VendorController struct{ svc services.VendorService }

func NewVendorController(svc services.VendorService) *VendorController {
	return &VendorController{svc}
}

func (ctrl *VendorController) List(c *gin.Context) {
	p := utils.ParsePage(c)
	list, total, err := ctrl.svc.List(utils.ClientCompanyID(c), p, c.Query("kind"))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, utils.Paged(list, total, p))
}

func (ctrl *VendorController) GetByID(c *gin.Context) {
	v, err := ctrl.svc.GetByID(utils.ClientCompanyID(c), c.Param("id"))
	if err != nil {
		utils.ErrNotFound(c, "Vendor not found")
		return
	}
	utils.OK(c, v)
}

func (ctrl *VendorController) Create(c *gin.Context) {
	var req dto.CreateVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	v, err := ctrl.svc.Create(utils.ClientCompanyID(c), req)
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.Created(c, v)
}

func (ctrl *VendorController) Update(c *gin.Context) {
	var req dto.UpdateVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	v, err := ctrl.svc.Update(utils.ClientCompanyID(c), c.Param("id"), req)
	if err != nil {
		utils.ErrNotFound(c, err.Error())
		return
	}
	utils.OK(c, v)
}

func (ctrl *VendorController) Delete(c *gin.Context) {
	if err := ctrl.svc.Delete(utils.ClientCompanyID(c), c.Param("id")); err != nil {
		utils.ErrConflict(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"ok": true})
}

type PartController struct{ svc services.PartService }

func NewPartController(svc services.PartService) *PartController {
	return &PartController{svc}
}

func (ctrl *PartController) List(c *gin.Context) {
	p := utils.ParsePage(c)
	list, total, err := ctrl.svc.List(utils.ClientCompanyID(c), p, c.Query("vendorId"))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, utils.Paged(list, total, p))
}

func (ctrl *PartController) GetByID(c *gin.Context) {
	part, err := ctrl.svc.GetByID(utils.ClientCompanyID(c), c.Param("id"))
	if err != nil {
		utils.ErrNotFound(c, "Part not found")
		return
	}
	utils.OK(c, part)
}

func (ctrl *PartController) Create(c *gin.Context) {
	var req dto.CreatePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	part, err := ctrl.svc.Create(utils.ClientCompanyID(c), req)
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.Created(c, part)
}

func (ctrl *PartController) Update(c *gin.Context) {
	var req dto.UpdatePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	part, err := ctrl.svc.Update(utils.ClientCompanyID(c), c.Param("id"), req)
	if err != nil {
		utils.ErrNotFound(c, err.Error())
		return
	}
	utils.OK(c, part)
}

func (ctrl *PartController) Delete(c *gin.Context) {
	if err := ctrl.svc.Delete(utils.ClientCompanyID(c), c.Param("id")); err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, gin.H{"ok": true})
}

type ReportController struct{ svc services.ReportService }

func NewReportController(svc services.ReportService) *ReportController {
	return &ReportController{svc}
}

func (ctrl *ReportController) Dashboard(c *gin.Context) {
	data, err := ctrl.svc.Dashboard(utils.ClientCompanyID(c))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, data)
}

func (ctrl *ReportController) Revenue(c *gin.Context) {
	data, err := ctrl.svc.Revenue(utils.ClientCompanyID(c), c.Query("from"), c.Query("to"), c.Query("groupBy"))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, data)
}

func (ctrl *ReportController) TopCustomers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	data, err := ctrl.svc.TopCustomers(utils.ClientCompanyID(c), limit, c.Query("from"), c.Query("to"))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, data)
}

func (ctrl *ReportController) MechanicPerformance(c *gin.Context) {
	data, err := ctrl.svc.MechanicPerformance(utils.ClientCompanyID(c), c.Query("from"), c.Query("to"))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, data)
}

type DashboardController struct{ svc services.DashboardService }

func NewDashboardController(svc services.DashboardService) *DashboardController {
	return &DashboardController{svc}
}

func (ctrl *DashboardController) Financials(c *gin.Context) {
	data, err := ctrl.svc.Financials(utils.ClientCompanyID(c), c.Query("month"))
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, data)
}

type PublicController struct{ svc services.PublicService }

func NewPublicController(svc services.PublicService) *PublicController {
	return &PublicController{svc}
}

func (ctrl *PublicController) GetEstimate(c *gin.Context) {
	data, err := ctrl.svc.GetEstimate(c.Param("slug"))
	if err != nil {
		utils.ErrNotFound(c, "Estimate not found")
		return
	}
	utils.OK(c, data)
}

func (ctrl *PublicController) GetInvoice(c *gin.Context) {
	data, err := ctrl.svc.GetInvoice(c.Param("slug"))
	if err != nil {
		utils.ErrNotFound(c, "Invoice not found")
		return
	}
	utils.OK(c, data)
}
