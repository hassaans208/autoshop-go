package repositories

import (
	"autoworkshop/models"
	"autoworkshop/utils"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	List(tenantID string, p utils.PageParams) ([]models.Customer, int64, error)
	FindByID(tenantID, id string) (*models.Customer, error)
	FindByPhone(tenantID, phoneNorm string) (*models.Customer, error)
	Create(c *models.Customer) error
	Update(c *models.Customer) error
	Delete(tenantID, id string) error
	HasCarsOrInvoices(tenantID, id string) (bool, error)
	Count(tenantID string) (int64, error)
}

type customerRepo struct{ db *gorm.DB }

func NewCustomerRepository(db *gorm.DB) CustomerRepository { return &customerRepo{db} }

func (r *customerRepo) List(tenantID string, p utils.PageParams) ([]models.Customer, int64, error) {
	var list []models.Customer
	var total int64
	q := scopeTenant(r.db.Model(&models.Customer{}), tenantID)
	if p.Q != "" {
		norm := utils.NormalizePhone(p.Q)
		q = q.Where("name LIKE ? OR phone_norm LIKE ?", "%"+p.Q+"%", "%"+norm+"%")
	}
	q.Count(&total)
	err := q.Offset(utils.Offset(p)).Limit(p.Limit).Find(&list).Error
	return list, total, err
}

func (r *customerRepo) FindByID(tenantID, id string) (*models.Customer, error) {
	var c models.Customer
	return &c, scopeTenant(r.db, tenantID).First(&c, "id = ?", id).Error
}

func (r *customerRepo) FindByPhone(tenantID, phoneNorm string) (*models.Customer, error) {
	var c models.Customer
	return &c, scopeTenant(r.db, tenantID).Where("phone_norm = ?", phoneNorm).First(&c).Error
}

func (r *customerRepo) Create(c *models.Customer) error { return r.db.Create(c).Error }

func (r *customerRepo) Update(c *models.Customer) error { return r.db.Save(c).Error }

func (r *customerRepo) Delete(tenantID, id string) error {
	return scopeTenant(r.db, tenantID).Delete(&models.Customer{}, "id = ?", id).Error
}

func (r *customerRepo) HasCarsOrInvoices(tenantID, id string) (bool, error) {
	var carCount, invCount int64
	scopeTenant(r.db.Model(&models.Car{}), tenantID).Where("customer_id = ?", id).Count(&carCount)
	scopeTenant(r.db.Model(&models.Invoice{}), tenantID).Where("customer_id = ?", id).Count(&invCount)
	return carCount > 0 || invCount > 0, nil
}

func (r *customerRepo) Count(tenantID string) (int64, error) {
	var count int64
	err := scopeTenant(r.db.Model(&models.Customer{}), tenantID).Count(&count).Error
	return count, err
}
