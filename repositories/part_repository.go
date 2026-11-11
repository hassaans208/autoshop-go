package repositories

import (
	"autoworkshop/models"
	"autoworkshop/utils"

	"gorm.io/gorm"
)

type PartRepository interface {
	List(tenantID string, p utils.PageParams, vendorID string) ([]models.Part, int64, error)
	FindByID(tenantID, id string) (*models.Part, error)
	Create(p *models.Part) error
	Update(p *models.Part) error
	Delete(tenantID, id string) error
}

type partRepo struct{ db *gorm.DB }

func NewPartRepository(db *gorm.DB) PartRepository { return &partRepo{db} }

func (r *partRepo) List(tenantID string, p utils.PageParams, vendorID string) ([]models.Part, int64, error) {
	var list []models.Part
	var total int64
	q := scopeTenant(r.db.Model(&models.Part{}), tenantID)
	if vendorID != "" {
		q = q.Where("vendor_id = ?", vendorID)
	}
	if p.Q != "" {
		q = q.Where("name LIKE ?", "%"+p.Q+"%")
	}
	q.Count(&total)
	err := q.Offset(utils.Offset(p)).Limit(p.Limit).Order("name ASC").Find(&list).Error
	return list, total, err
}

func (r *partRepo) FindByID(tenantID, id string) (*models.Part, error) {
	var part models.Part
	return &part, scopeTenant(r.db, tenantID).First(&part, "id = ?", id).Error
}

func (r *partRepo) Create(p *models.Part) error { return r.db.Create(p).Error }

func (r *partRepo) Update(p *models.Part) error { return r.db.Save(p).Error }

func (r *partRepo) Delete(tenantID, id string) error {
	return scopeTenant(r.db, tenantID).Delete(&models.Part{}, "id = ?", id).Error
}
