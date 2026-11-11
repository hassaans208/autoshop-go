package repositories

import (
	"autoworkshop/models"
	"strings"

	"gorm.io/gorm"
)

type CompanyRepository interface {
	List(tenantID string) ([]models.Company, error)
	AddIfNotExists(tenantID, name string) (*models.Company, error)
}

type companyRepo struct{ db *gorm.DB }

func NewCompanyRepository(db *gorm.DB) CompanyRepository { return &companyRepo{db} }

func (r *companyRepo) List(tenantID string) ([]models.Company, error) {
	var list []models.Company
	err := scopeTenant(r.db, tenantID).Order("name ASC").Find(&list).Error
	return list, err
}

func (r *companyRepo) AddIfNotExists(tenantID, name string) (*models.Company, error) {
	normalized := strings.TrimSpace(name)
	var c models.Company
	err := scopeTenant(r.db, tenantID).Where("LOWER(name) = LOWER(?)", normalized).First(&c).Error
	if err == nil {
		return &c, nil
	}
	c = models.Company{ClientCompanyID: tenantID, Name: normalized}
	return &c, r.db.Create(&c).Error
}
