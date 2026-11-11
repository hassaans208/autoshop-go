package migrations

import (
	"autoworkshop/models"
	"log"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.ClientCompany{},
		&models.Staff{},
		&models.Settings{},
		&models.Company{},
		&models.JobType{},
		&models.Customer{},
		&models.Car{},
		&models.Estimate{},
		&models.EstimateItem{},
		&models.EstimateImage{},
		&models.Invoice{},
		&models.InvoiceItem{},
		&models.InvoiceImage{},
		&models.Payment{},
		&models.DocCounter{},
		&models.Vendor{},
		&models.Part{},
	)
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}
	widenNameColumns(db)
	log.Println("✅ Migrations complete")
}

func widenNameColumns(db *gorm.DB) {
	stmts := []string{
		"ALTER TABLE client_companies MODIFY COLUMN name VARCHAR(255) NOT NULL",
		"ALTER TABLE staffs MODIFY COLUMN name VARCHAR(255) NOT NULL",
		"ALTER TABLE settings MODIFY COLUMN workshop_name VARCHAR(255)",
	}
	for _, sql := range stmts {
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("⚠️  widen column (may already be applied): %v", err)
		}
	}
}
