package repositories

import (
	"autoworkshop/models"
	"autoworkshop/utils"

	"gorm.io/gorm"
)

type VendorRepository interface {
	List(tenantID string, p utils.PageParams, kind string) ([]models.Vendor, int64, error)
	FindByID(tenantID, id string) (*models.Vendor, error)
	Create(v *models.Vendor) error
	Update(v *models.Vendor) error
	Delete(tenantID, id string) error
	HasParts(tenantID, id string) (bool, error)
}

type vendorRepo struct{ db *gorm.DB }

func NewVendorRepository(db *gorm.DB) VendorRepository { return &vendorRepo{db} }

func (r *vendorRepo) List(tenantID string, p utils.PageParams, kind string) ([]models.Vendor, int64, error) {
	var list []models.Vendor
	var total int64
	q := scopeTenant(r.db.Model(&models.Vendor{}), tenantID)
	if kind != "" {
		q = q.Where("kind = ?", kind)
	}
	if p.Q != "" {
		like := "%" + p.Q + "%"
		q = q.Where("name LIKE ? OR phone LIKE ? OR email LIKE ?", like, like, like)
	}
	q.Count(&total)
	err := q.Offset(utils.Offset(p)).Limit(p.Limit).Order("name ASC").Find(&list).Error
	return list, total, err
}

func (r *vendorRepo) FindByID(tenantID, id string) (*models.Vendor, error) {
	var v models.Vendor
	return &v, scopeTenant(r.db, tenantID).First(&v, "id = ?", id).Error
}

func (r *vendorRepo) Create(v *models.Vendor) error { return r.db.Create(v).Error }

func (r *vendorRepo) Update(v *models.Vendor) error { return r.db.Save(v).Error }

func (r *vendorRepo) Delete(tenantID, id string) error {
	return scopeTenant(r.db, tenantID).Delete(&models.Vendor{}, "id = ?", id).Error
}

func (r *vendorRepo) HasParts(tenantID, id string) (bool, error) {
	var count int64
	err := scopeTenant(r.db.Model(&models.Part{}), tenantID).Where("vendor_id = ?", id).Count(&count).Error
	return count > 0, err
}
