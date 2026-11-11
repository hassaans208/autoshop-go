package main

import (
	"autoworkshop/config"
	"autoworkshop/database"
	"autoworkshop/models"
	"autoworkshop/routes"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)
	database.Migrate(db)
	seed(db)

	if cfg.AutologinLinkSecret != "" {
		log.Println("✅ autologin-link API active (POST /api/v1/<AUTOLOGIN_LINK_SECRET>/autologin-link)")
	} else {
		log.Println("⚠️  autologin-link API off — set AUTOLOGIN_LINK_SECRET in .env or environment")
	}
	if cfg.ImportLeadsSecret != "" {
		log.Println("✅ import-leads API active (POST /api/v1/<IMPORT_LEADS_SECRET>/import-leads)")
	} else {
		log.Println("⚠️  import-leads API off — set IMPORT_LEADS_SECRET in .env or environment")
	}

	r := routes.Setup(db, cfg.JWTSecret, cfg.AdminAPIKey, cfg.AdminEmail, cfg.AdminPassword, cfg.FrontendURL, cfg.AutologinLinkSecret, cfg.ImportLeadsSecret)
	r.Run(cfg.Port)
}

func seed(db *gorm.DB) {
	var count int64
	db.Model(&models.ClientCompany{}).Count(&count)
	if count > 0 {
		return
	}

	ownerID := uuid.NewString()
	companyID := uuid.NewString()

	if err := db.Create(&models.ClientCompany{
		ID: companyID, Name: "AutoWorkshop Pro", OwnerID: ownerID, Active: true,
	}).Error; err != nil {
		log.Fatal("seed ClientCompany:", err)
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("owner123"), bcrypt.DefaultCost)
	if err := db.Create(&models.Staff{
		ID:              ownerID,
		ClientCompanyID: companyID,
		Email:           "owner@workshop.com",
		PasswordHash:    string(hash),
		Name:            "Ahmed Khan",
		Role:            "owner",
		Phone:           "0300-1111111",
		Salary:          80000,
		Active:          true,
	}).Error; err != nil {
		log.Fatal("seed Staff:", err)
	}

	db.Create(&models.Settings{
		ClientCompanyID: companyID,
		WorkshopName:    "AutoWorkshop Pro",
		Tagline:         "Quality Auto Repair & Service",
		Currency:        "PKR",
		EstPrefix:       "EST",
		InvPrefix:       "INV",
	})

	log.Println("✅ Seed data created (owner@workshop.com / owner123)")
}
