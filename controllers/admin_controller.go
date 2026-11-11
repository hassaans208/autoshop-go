package controllers

import (
	"autoworkshop/dto"
	"autoworkshop/services"
	"autoworkshop/utils"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	svc     services.AdminService
	authSvc services.AdminAuthService
}

func NewAdminController(svc services.AdminService, authSvc services.AdminAuthService) *AdminController {
	return &AdminController{svc, authSvc}
}

func (ctrl *AdminController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	token, user, err := ctrl.authSvc.Login(req.Email, req.Password)
	if err != nil {
		utils.ErrUnauthorized(c, err.Error())
		return
	}
	utils.OK(c, dto.LoginResponse{Token: token, User: user})
}

func (ctrl *AdminController) CreateCompany(c *gin.Context) {
	var req dto.AdminCreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	resp, err := ctrl.svc.CreateCompany(req)
	if err != nil {
		utils.ErrConflict(c, err.Error())
		return
	}
	utils.Created(c, resp)
}

func (ctrl *AdminController) ListCompanies(c *gin.Context) {
	list, err := ctrl.svc.ListCompanies()
	if err != nil {
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, gin.H{"data": list})
}
