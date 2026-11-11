package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientCompanyID string         `gorm:"type:varchar(36);not null;index;uniqueIndex:idx_company_name_cc" json:"-"`
	Name            string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_company_name_cc" json:"name"`
	CreatedAt       time.Time      `json:"createdAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
