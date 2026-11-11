package repositories

import (
	"autoworkshop/models"
	"autoworkshop/utils"
	"fmt"

	"gorm.io/gorm"
)

type EstimateRepository interface {
	List(tenantID string, p utils.PageParams, filters map[string]string) ([]models.Estimate, int64, error)
	FindByID(tenantID, id string) (*models.Estimate, error)
	FindBySlug(slug string) (*models.Estimate, error)
	FindByCar(tenantID, carID string) ([]models.Estimate, error)
	Create(e *models.Estimate) error
	Update(e *models.Estimate) error
	Delete(tenantID, id string) error
	NextNumber(tenantID, prefix string, year int) (string, error)
	CountByStatus(tenantID, status string) (int64, error)
}

type estimateRepo struct{ db *gorm.DB }

func NewEstimateRepository(db *gorm.DB) EstimateRepository { return &estimateRepo{db} }

func (r *estimateRepo) List(tenantID string, p utils.PageParams, filters map[string]string) ([]models.Estimate, int64, error) {
	var list []models.Estimate
	var total int64
	q := scopeTenant(r.db.Model(&models.Estimate{}), tenantID)
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
	err := q.Preload("Items").Preload("Images").
		Offset(utils.Offset(p)).Limit(p.Limit).
		Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (r *estimateRepo) FindByID(tenantID, id string) (*models.Estimate, error) {
	var e models.Estimate
	err := scopeTenant(r.db, tenantID).Preload("Items").Preload("Images").First(&e, "id = ?", id).Error
	return &e, err
}

func (r *estimateRepo) FindBySlug(slug string) (*models.Estimate, error) {
	var e models.Estimate
	err := r.db.Preload("Items").Preload("Images").Where("public_slug = ?", slug).First(&e).Error
	return &e, err
}

func (r *estimateRepo) FindByCar(tenantID, carID string) ([]models.Estimate, error) {
	var list []models.Estimate
	err := scopeTenant(r.db, tenantID).Preload("Items").Preload("Images").
		Where("car_id = ?", carID).Find(&list).Error
	return list, err
}

func (r *estimateRepo) Create(e *models.Estimate) error { return r.db.Create(e).Error }

func (r *estimateRepo) Update(e *models.Estimate) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
}

func (r *estimateRepo) Delete(tenantID, id string) error {
	return scopeTenant(r.db, tenantID).Delete(&models.Estimate{}, "id = ?", id).Error
}

func (r *estimateRepo) NextNumber(tenantID, prefix string, year int) (string, error) {
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

func (r *estimateRepo) CountByStatus(tenantID, status string) (int64, error) {
	var count int64
	err := scopeTenant(r.db.Model(&models.Estimate{}), tenantID).Where("status = ?", status).Count(&count).Error
	return count, err
}
