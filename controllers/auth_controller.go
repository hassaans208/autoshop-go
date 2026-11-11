package controllers

import (
	"autoworkshop/dto"
	"autoworkshop/services"
	"autoworkshop/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct{ svc services.AuthService }

func NewAuthController(svc services.AuthService) *AuthController { return &AuthController{svc} }

func (ctrl *AuthController) Autologin(c *gin.Context) {
	var req dto.AutologinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	token, staff, err := ctrl.svc.Autologin(req.Email, req.PasswordHash)
	if err != nil {
		utils.ErrUnauthorized(c, "invalid credentials")
		return
	}
	utils.OK(c, dto.LoginResponse{
		Token: token,
		User:  dto.UserResponse{ID: staff.ID, Email: staff.Email, Name: staff.Name, Role: staff.Role, Phone: staff.Phone, Active: staff.Active},
	})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	token, staff, err := ctrl.svc.Login(req.Email, req.Password)
	if err != nil {
		utils.ErrUnauthorized(c, err.Error())
		return
	}
	utils.OK(c, dto.LoginResponse{
		Token: token,
		User:  dto.UserResponse{ID: staff.ID, Email: staff.Email, Name: staff.Name, Role: staff.Role, Phone: staff.Phone, Active: staff.Active},
	})
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	utils.OK(c, gin.H{"ok": true})
}

func (ctrl *AuthController) AutologinLink(c *gin.Context) {
	var req dto.AutologinLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrBadRequest(c, "Validation failed", err.Error())
		return
	}
	link, staff, err := ctrl.svc.AutologinLink(req.StaffID)
	if err != nil {
		if err.Error() == "staff not found" {
			utils.ErrNotFound(c, err.Error())
			return
		}
		if err.Error() == "FRONTEND_URL is not configured" {
			utils.ErrUnprocessable(c, err.Error())
			return
		}
		utils.ErrInternal(c, err)
		return
	}
	utils.OK(c, dto.AutologinLinkResponse{URL: link, Email: staff.Email})
}

func (ctrl *AuthController) Me(c *gin.Context) {
	tenantID := utils.ClientCompanyID(c)
	staffID := utils.StaffID(c)
	staff, err := ctrl.svc.Me(tenantID, staffID)
	if err != nil {
		utils.ErrNotFound(c, "Staff not found")
		return
	}
	utils.OK(c, dto.UserResponse{ID: staff.ID, Email: staff.Email, Name: staff.Name, Role: staff.Role, Phone: staff.Phone, Active: staff.Active})
}
