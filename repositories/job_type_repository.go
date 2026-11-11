package repositories

import (
	"autoworkshop/models"
	"encoding/json"

	"gorm.io/gorm"
)

type JobTypeRepository interface {
	List(tenantID string) ([]models.JobType, error)
	ReplaceAll(tenantID string, items []models.JobType) ([]models.JobType, error)
}

type jobTypeRepo struct{ db *gorm.DB }

func NewJobTypeRepository(db *gorm.DB) JobTypeRepository { return &jobTypeRepo{db} }

func (r *jobTypeRepo) List(tenantID string) ([]models.JobType, error) {
	var list []models.JobType
	err := scopeTenant(r.db, tenantID).Order("name ASC").Find(&list).Error
	return list, err
}

func (r *jobTypeRepo) ReplaceAll(tenantID string, items []models.JobType) ([]models.JobType, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := scopeTenant(tx, tenantID).Delete(&models.JobType{}).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].ClientCompanyID = tenantID
			if items[i].Subtypes == "" {
				items[i].Subtypes = "[]"
			}
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.List(tenantID)
}

func SubtypesToJSON(subtypes []string) string {
	if subtypes == nil {
		return "[]"
	}
	b, _ := json.Marshal(subtypes)
	return string(b)
}

func SubtypesFromJSON(raw string) []string {
	if raw == "" {
		return []string{}
	}
	var out []string
	_ = json.Unmarshal([]byte(raw), &out)
	return out
}
