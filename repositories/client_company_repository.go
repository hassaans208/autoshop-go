package repositories

import (
	"autoworkshop/models"

	"gorm.io/gorm"
)

type ClientCompanyRepository interface {
	Create(c *models.ClientCompany) error
	FindByID(id string) (*models.ClientCompany, error)
	List() ([]models.ClientCompany, error)
}

type clientCompanyRepo struct{ db *gorm.DB }

func NewClientCompanyRepository(db *gorm.DB) ClientCompanyRepository {
	return &clientCompanyRepo{db}
}

func (r *clientCompanyRepo) Create(c *models.ClientCompany) error { return r.db.Create(c).Error }

func (r *clientCompanyRepo) FindByID(id string) (*models.ClientCompany, error) {
	var c models.ClientCompany
	return &c, r.db.First(&c, "id = ?", id).Error
}

func (r *clientCompanyRepo) List() ([]models.ClientCompany, error) {
	var list []models.ClientCompany
	return list, r.db.Order("created_at DESC").Find(&list).Error
}
