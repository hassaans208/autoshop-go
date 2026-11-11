package repositories

import (
	"autoworkshop/models"
	"autoworkshop/utils"
	"fmt"

	"gorm.io/gorm"
)

type InvoiceRepository interface {
	List(tenantID string, p utils.PageParams, filters map[string]string) ([]models.Invoice, int64, error)
	FindByID(tenantID, id string) (*models.Invoice, error)
	FindBySlug(slug string) (*models.Invoice, error)
	FindByCustomer(tenantID, customerID string) ([]models.Invoice, error)
	FindByCar(tenantID, carID string) ([]models.Invoice, error)
	Create(inv *models.Invoice) error
	Update(inv *models.Invoice) error
	Delete(tenantID, id string) error
	NextNumber(tenantID, prefix string, year int) (string, error)
	CountByStatus(tenantID string, statuses ...string) (int64, error)
	Recent(tenantID string, limit int) ([]models.Invoice, error)
}

type invoiceRepo struct{ db *gorm.DB }

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository { return &invoiceRepo{db} }

func (r *invoiceRepo) List(tenantID string, p utils.PageParams, filters map[string]string) ([]models.Invoice, int64, error) {
	var list []models.Invoice
	var total int64
	q := scopeTenant(r.db.Model(&models.Invoice{}), tenantID)
	if v := filters["status"]; v != "" {
		q = q.Where("status = ?", v)
	}
	if v := filters["customerId"]; v != "" {
		q = q.Where("customer_id = ?", v)
	}
	if v := filters["carId"]; v != "" {
		q = q.Where("car_id = ?", v)
	}
	if v := filters["from"]; v != "" {
		q = q.Where("date >= ?", v)
	}
	if v := filters["to"]; v != "" {
		q = q.Where("date <= ?", v)
	}
	if p.Q != "" {
		q = q.Where("number LIKE ?", "%"+p.Q+"%")
	}
	q.Count(&total)
	err := q.Preload("Items").Preload("Images").Preload("Payments").
		Offset(utils.Offset(p)).Limit(p.Limit).
		Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (r *invoiceRepo) FindByID(tenantID, id string) (*models.Invoice, error) {
	var inv models.Invoice
	err := scopeTenant(r.db, tenantID).Preload("Items").Preload("Images").Preload("Payments").
		First(&inv, "id = ?", id).Error
	return &inv, err
}

func (r *invoiceRepo) FindBySlug(slug string) (*models.Invoice, error) {
	var inv models.Invoice
	err := r.db.Preload("Items").Preload("Images").Preload("Payments").
		Where("public_slug = ?", slug).First(&inv).Error
	return &inv, err
}

func (r *invoiceRepo) FindByCustomer(tenantID, customerID string) ([]models.Invoice, error) {
	var list []models.Invoice
	err := scopeTenant(r.db, tenantID).Preload("Items").Preload("Images").Preload("Payments").
		Where("customer_id = ?", customerID).Find(&list).Error
	return list, err
}

func (r *invoiceRepo) FindByCar(tenantID, carID string) ([]models.Invoice, error) {
	var list []models.Invoice
	err := scopeTenant(r.db, tenantID).Preload("Items").Preload("Images").Preload("Payments").
		Where("car_id = ?", carID).Find(&list).Error
	return list, err
}

func (r *invoiceRepo) Create(inv *models.Invoice) error { return r.db.Create(inv).Error }

func (r *invoiceRepo) Update(inv *models.Invoice) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(inv).Error
}

func (r *invoiceRepo) Delete(tenantID, id string) error {
	return scopeTenant(r.db, tenantID).Delete(&models.Invoice{}, "id = ?", id).Error
}

func (r *invoiceRepo) NextNumber(tenantID, prefix string, year int) (string, error) {
	var counter models.DocCounter
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("client_company_id = ? AND prefix = ? AND year = ?", tenantID, prefix, year).
			FirstOrCreate(&counter, models.DocCounter{ClientCompanyID: tenantID, Prefix: prefix, Year: year, Last: 0}).Error; err != nil {
			return err
		}
		counter.Last++
		return tx.Save(&counter).Error
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%d-%05d", prefix, year, counter.Last), nil
}

func (r *invoiceRepo) CountByStatus(tenantID string, statuses ...string) (int64, error) {
	var count int64
	q := scopeTenant(r.db.Model(&models.Invoice{}), tenantID)
	if len(statuses) > 0 {
		q = q.Where("status IN ?", statuses)
	}
	err := q.Count(&count).Error
	return count, err
}

func (r *invoiceRepo) Recent(tenantID string, limit int) ([]models.Invoice, error) {
	var list []models.Invoice
	err := scopeTenant(r.db, tenantID).Preload("Items").Preload("Payments").
		Order("created_at DESC").Limit(limit).Find(&list).Error
	return list, err
}
