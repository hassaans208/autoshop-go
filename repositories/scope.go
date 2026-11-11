package repositories

import "gorm.io/gorm"

func scopeTenant(db *gorm.DB, tenantID string) *gorm.DB {
	return db.Where("client_company_id = ?", tenantID)
}
