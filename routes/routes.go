package routes

import (
	"autoworkshop/controllers"
	"autoworkshop/middleware"
	"autoworkshop/repositories"
	"autoworkshop/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, jwtSecret, adminAPIKey, adminEmail, adminPassword, frontendURL, autologinLinkSecret, importLeadsSecret string) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Repositories
	staffRepo := repositories.NewStaffRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	vehicleRepo := repositories.NewVehicleRepository(db)
	estimateRepo := repositories.NewEstimateRepository(db)
	invoiceRepo := repositories.NewInvoiceRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	settingsRepo := repositories.NewSettingsRepository(db)
	companyRepo := repositories.NewCompanyRepository(db)
	jobTypeRepo := repositories.NewJobTypeRepository(db)
	vendorRepo := repositories.NewVendorRepository(db)
	partRepo := repositories.NewPartRepository(db)
	clientCompanyRepo := repositories.NewClientCompanyRepository(db)

	// Services
	authSvc := services.NewAuthService(staffRepo, jwtSecret, frontendURL)
	staffSvc := services.NewStaffService(staffRepo)
	customerSvc := services.NewCustomerService(customerRepo)
	vehicleSvc := services.NewVehicleService(vehicleRepo)
	estimateSvc := services.NewEstimateService(estimateRepo, invoiceRepo)
	invoiceSvc := services.NewInvoiceService(invoiceRepo, paymentRepo)
	settingsSvc := services.NewSettingsService(settingsRepo)
	companySvc := services.NewCompanyService(companyRepo)
	jobTypeSvc := services.NewJobTypeService(jobTypeRepo)
	vendorSvc := services.NewVendorService(vendorRepo)
	partSvc := services.NewPartService(partRepo)
	reportSvc := services.NewReportService(invoiceRepo, estimateRepo, customerRepo, staffRepo)
	dashboardSvc := services.NewDashboardService(staffRepo, customerRepo, vehicleRepo, estimateRepo, invoiceRepo)
	publicSvc := services.NewPublicService(estimateRepo, invoiceRepo, settingsRepo, customerRepo, vehicleRepo)
	adminSvc := services.NewAdminService(db, clientCompanyRepo, staffRepo)
	adminAuthSvc := services.NewAdminAuthService(adminEmail, adminPassword, jwtSecret)
	leadImportSvc := services.NewLeadImportService(db, staffRepo, frontendURL)

	// Controllers
	authCtrl := controllers.NewAuthController(authSvc)
	staffCtrl := controllers.NewStaffController(staffSvc)
	customerCtrl := controllers.NewCustomerController(customerSvc)
	vehicleCtrl := controllers.NewVehicleController(vehicleSvc)
	estimateCtrl := controllers.NewEstimateController(estimateSvc, settingsSvc)
	invoiceCtrl := controllers.NewInvoiceController(invoiceSvc, settingsSvc)
	settingsCtrl := controllers.NewSettingsController(settingsSvc)
	companyCtrl := controllers.NewCompanyController(companySvc)
	jobTypeCtrl := controllers.NewJobTypeController(jobTypeSvc)
	vendorCtrl := controllers.NewVendorController(vendorSvc)
	partCtrl := controllers.NewPartController(partSvc)
	reportCtrl := controllers.NewReportController(reportSvc)
	dashboardCtrl := controllers.NewDashboardController(dashboardSvc)
	publicCtrl := controllers.NewPublicController(publicSvc)
	adminCtrl := controllers.NewAdminController(adminSvc, adminAuthSvc)
	leadImportCtrl := controllers.NewLeadImportController(leadImportSvc)

	r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })

	v1 := r.Group("/api/v1")

	// Public (no auth) — frontend uses /estimates/slug/:slug and /invoices/slug/:slug
	v1.POST("/auth/login", authCtrl.Login)
	v1.POST("/auth/autologin", authCtrl.Autologin)
	if autologinLinkSecret != "" {
		v1.POST("/:linkSecret/autologin-link", middleware.AutologinLinkGate(autologinLinkSecret), authCtrl.AutologinLink)
	}
	if importLeadsSecret != "" {
		v1.POST("/:linkSecret/import-leads", middleware.SecretPathGate(importLeadsSecret), leadImportCtrl.Import)
	}
	v1.GET("/estimates/slug/:slug", publicCtrl.GetEstimate)
	v1.GET("/invoices/slug/:slug", publicCtrl.GetInvoice)
	v1.GET("/public/estimates/:slug", publicCtrl.GetEstimate)
	v1.GET("/public/invoices/:slug", publicCtrl.GetInvoice)

	v1.POST("/admin/auth/login", adminCtrl.Login)

	admin := v1.Group("/admin")
	admin.Use(middleware.AdminGate(adminAPIKey, jwtSecret))
	{
		admin.GET("/companies", adminCtrl.ListCompanies)
		admin.POST("/companies", adminCtrl.CreateCompany)
	}

	api := v1.Group("")
	api.Use(middleware.AuthGate(jwtSecret))

	api.GET("/auth/me", authCtrl.Me)
	api.POST("/auth/logout", authCtrl.Logout)

	api.GET("/settings", settingsCtrl.Get)
	api.PUT("/settings", middleware.RequireRole("owner"), settingsCtrl.Update)

	api.GET("/companies", companyCtrl.List)
	api.POST("/companies", companyCtrl.Create)

	api.GET("/job-types", jobTypeCtrl.List)
	api.PUT("/job-types", middleware.RequireRole("owner", "manager"), jobTypeCtrl.Replace)

	staff := api.Group("/staff")
	staff.Use(middleware.RequireRole("owner", "manager"))
	{
		staff.GET("", staffCtrl.List)
		staff.POST("", middleware.RequireRole("owner"), staffCtrl.Create)
		staff.PUT("/:id", staffCtrl.Update)
		staff.DELETE("/:id", middleware.RequireRole("owner"), staffCtrl.Delete)
	}

	vendors := api.Group("/vendors")
	vendors.Use(middleware.RequireRole("owner", "manager"))
	{
		vendors.GET("", vendorCtrl.List)
		vendors.GET("/:id", vendorCtrl.GetByID)
		vendors.POST("", vendorCtrl.Create)
		vendors.PUT("/:id", vendorCtrl.Update)
		vendors.DELETE("/:id", vendorCtrl.Delete)
	}

	parts := api.Group("/parts")
	parts.Use(middleware.RequireRole("owner", "manager"))
	{
		parts.GET("", partCtrl.List)
		parts.GET("/:id", partCtrl.GetByID)
		parts.POST("", partCtrl.Create)
		parts.PUT("/:id", partCtrl.Update)
		parts.DELETE("/:id", partCtrl.Delete)
	}

	cust := api.Group("/customers")
	cust.Use(middleware.RequireRole("owner", "manager", "receptionist", "accountant"))
	{
		cust.GET("", customerCtrl.List)
		cust.GET("/:id", customerCtrl.GetByID)
		cust.GET("/by-phone/:phone", customerCtrl.GetByPhone)
		cust.POST("", customerCtrl.Create)
		cust.PUT("/:id", customerCtrl.Update)
		cust.DELETE("/:id", middleware.RequireRole("owner", "manager"), customerCtrl.Delete)
		cust.GET("/:id/cars", vehicleCtrl.GetByCustomer)
		cust.GET("/:id/invoices", invoiceCtrl.GetByCustomer)
	}

	cars := api.Group("/cars")
	cars.Use(middleware.RequireRole("owner", "manager", "receptionist", "accountant", "mechanic"))
	{
		cars.GET("", vehicleCtrl.List)
		cars.GET("/:id", vehicleCtrl.GetByID)
		cars.POST("", middleware.RequireRole("owner", "manager", "receptionist"), vehicleCtrl.Create)
		cars.PUT("/:id", middleware.RequireRole("owner", "manager", "receptionist"), vehicleCtrl.Update)
		cars.DELETE("/:id", middleware.RequireRole("owner", "manager"), vehicleCtrl.Delete)
		cars.GET("/:id/estimates", estimateCtrl.GetByCar)
		cars.GET("/:id/invoices", invoiceCtrl.GetByCar)
	}

	est := api.Group("/estimates")
	est.Use(middleware.RequireRole("owner", "manager", "receptionist"))
	{
		est.GET("", estimateCtrl.List)
		est.GET("/:id", estimateCtrl.GetByID)
		est.POST("", estimateCtrl.Create)
		est.PUT("/:id", estimateCtrl.Update)
		est.DELETE("/:id", estimateCtrl.Delete)
		est.POST("/:id/convert", estimateCtrl.Convert)
	}

	inv := api.Group("/invoices")
	inv.Use(middleware.RequireRole("owner", "manager", "receptionist", "accountant"))
	{
		inv.GET("", invoiceCtrl.List)
		inv.GET("/:id", invoiceCtrl.GetByID)
		inv.POST("", middleware.RequireRole("owner", "manager", "receptionist"), invoiceCtrl.Create)
		inv.PUT("/:id", middleware.RequireRole("owner", "manager", "receptionist"), invoiceCtrl.Update)
		inv.DELETE("/:id", middleware.RequireRole("owner", "manager"), invoiceCtrl.Delete)
		inv.POST("/:id/payments", invoiceCtrl.AddPayment)
		inv.DELETE("/:id/payments/:paymentId", middleware.RequireRole("owner", "manager"), invoiceCtrl.DeletePayment)
	}

	reports := api.Group("/reports")
	reports.Use(middleware.RequireRole("owner", "manager", "accountant"))
	{
		reports.GET("/dashboard", reportCtrl.Dashboard)
		reports.GET("/revenue", reportCtrl.Revenue)
		reports.GET("/top-customers", reportCtrl.TopCustomers)
		reports.GET("/mechanic-performance", reportCtrl.MechanicPerformance)
	}

	api.GET("/dashboard/financials", middleware.RequireRole("owner", "manager", "accountant"), dashboardCtrl.Financials)

	return r
}
