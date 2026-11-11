package repositories

import (
	"autoworkshop/models"
	"autoworkshop/utils"

	"gorm.io/gorm"
)

type VehicleRepository interface {
	List(tenantID string, p utils.PageParams, customerID string) ([]models.Car, int64, error)
	FindByID(tenantID, id string) (*models.Car, error)
	FindByCustomer(tenantID, customerID string) ([]models.Car, error)
	Create(c *models.Car) error
	Update(c *models.Car) error
	Delete(tenantID, id string) error
	HasInvoices(tenantID, id string) (bool, error)
	EnsureCompany(tenantID, name string) error
	Count(tenantID string) (int64, error)
}

type vehicleRepo struct{ db *gorm.DB }

func NewVehicleRepository(db *gorm.DB) VehicleRepository { return &vehicleRepo{db} }

func (r *vehicleRepo) List(tenantID string, p utils.PageParams, customerID string) ([]models.Car, int64, error) {
	var list []models.Car
	var total int64
	q := scopeTenant(r.db.Model(&models.Car{}), tenantID)
	if customerID != "" {
		q = q.Where("customer_id = ?", customerID)
	}
	if p.Q != "" {
		like := "%" + p.Q + "%"
		q = q.Where("plate LIKE ? OR company LIKE ? OR model LIKE ?", like, like, like)
	}
	q.Count(&total)
	err := q.Offset(utils.Offset(p)).Limit(p.Limit).Find(&list).Error
	return list, total, err
}

func (r *vehicleRepo) FindByID(tenantID, id string) (*models.Car, error) {
	var c models.Car
	return &c, scopeTenant(r.db, tenantID).First(&c, "id = ?", id).Error
}

func (r *vehicleRepo) FindByCustomer(tenantID, customerID string) ([]models.Car, error) {
	var list []models.Car
	return list, scopeTenant(r.db, tenantID).Where("customer_id = ?", customerID).Find(&list).Error
}

func (r *vehicleRepo) Create(c *models.Car) error { return r.db.Create(c).Error }

func (r *vehicleRepo) Update(c *models.Car) error { return r.db.Save(c).Error }

func (r *vehicleRepo) Delete(tenantID, id string) error {
	return scopeTenant(r.db, tenantID).Delete(&models.Car{}, "id = ?", id).Error
}

func (r *vehicleRepo) HasInvoices(tenantID, id string) (bool, error) {
	var count int64
	err := scopeTenant(r.db.Model(&models.Invoice{}), tenantID).Where("car_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *vehicleRepo) EnsureCompany(tenantID, name string) error {
	_, err := NewCompanyRepository(r.db).AddIfNotExists(tenantID, name)
	return err
}

func (r *vehicleRepo) Count(tenantID string) (int64, error) {
	var count int64
	err := scopeTenant(r.db.Model(&models.Car{}), tenantID).Count(&count).Error
	return count, err
}
