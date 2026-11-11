package repositories

import (
	"autoworkshop/models"
	"errors"

	"gorm.io/gorm"
)

type SettingsRepository interface {
	Get(tenantID string) (*models.Settings, error)
	Save(tenantID string, s *models.Settings) error
}

type settingsRepo struct{ db *gorm.DB }

func NewSettingsRepository(db *gorm.DB) SettingsRepository { return &settingsRepo{db} }

func (r *settingsRepo) Get(tenantID string) (*models.Settings, error) {
	var s models.Settings
	err := scopeTenant(r.db, tenantID).First(&s).Error
	return &s, err
}

func (r *settingsRepo) Save(tenantID string, s *models.Settings) error {
	s.ClientCompanyID = tenantID
	var existing models.Settings
	err := scopeTenant(r.db, tenantID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(s).Error
	}
	s.ID = existing.ID
	return r.db.Save(s).Error
}
