package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Staff struct {
	ID              string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	ClientCompanyID string         `gorm:"type:varchar(36);not null;index;uniqueIndex:idx_staff_email_cc" json:"-"`
	Name            string         `gorm:"type:varchar(255);not null" json:"name"`
	Email           string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_staff_email_cc" json:"email"`
	PasswordHash    string         `gorm:"type:varchar(255);not null" json:"-"`
	Role            string         `gorm:"type:enum('owner','manager','mechanic','receptionist','accountant');not null" json:"role"`
	Phone           string         `gorm:"type:varchar(20)" json:"phone"`
	Salary          float64        `gorm:"type:decimal(12,2);default:0" json:"salary"`
	Active          bool           `gorm:"default:true;not null" json:"active"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Staff) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return nil
}
