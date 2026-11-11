package repositories

import (
	"autoworkshop/models"

	"gorm.io/gorm"
)

type StaffRepository interface {
	FindByEmail(email string) (*models.Staff, error)
	FindByID(tenantID, id string) (*models.Staff, error)
	FindByIDAny(id string) (*models.Staff, error)
	List(tenantID string) ([]models.Staff, error)
	Create(s *models.Staff) error
	Update(s *models.Staff) error
	Delete(tenantID, id string) error
	SumActiveSalaries(tenantID string) (float64, error)
}

type staffRepo struct{ db *gorm.DB }

func NewStaffRepository(db *gorm.DB) StaffRepository { return &staffRepo{db} }

func (r *staffRepo) FindByEmail(email string) (*models.Staff, error) {
	var s models.Staff
	return &s, r.db.Where("email = ?", email).First(&s).Error
}

func (r *staffRepo) FindByID(tenantID, id string) (*models.Staff, error) {
	var s models.Staff
	return &s, scopeTenant(r.db, tenantID).First(&s, "id = ?", id).Error
}

func (r *staffRepo) FindByIDAny(id string) (*models.Staff, error) {
	var s models.Staff
	return &s, r.db.First(&s, "id = ?", id).Error
}

func (r *staffRepo) List(tenantID string) ([]models.Staff, error) {
	var list []models.Staff
	return list, scopeTenant(r.db, tenantID).Find(&list).Error
}

func (r *staffRepo) Create(s *models.Staff) error { return r.db.Create(s).Error }

func (r *staffRepo) Update(s *models.Staff) error { return r.db.Save(s).Error }

func (r *staffRepo) Delete(tenantID, id string) error {
	return scopeTenant(r.db.Model(&models.Staff{}), tenantID).Where("id = ?", id).Update("active", false).Error
}

func (r *staffRepo) SumActiveSalaries(tenantID string) (float64, error) {
	var total float64
	err := scopeTenant(r.db.Model(&models.Staff{}), tenantID).Where("active = ?", true).
		Select("COALESCE(SUM(salary), 0)").Scan(&total).Error
	return total, err
}
