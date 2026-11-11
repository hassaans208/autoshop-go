package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID              string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	ClientCompanyID string         `gorm:"type:varchar(36);not null;index;uniqueIndex:idx_cust_phone_cc" json:"-"`
	Name            string         `gorm:"type:varchar(100);not null" json:"name"`
	Phone           string         `gorm:"type:varchar(20);not null" json:"phone"`
	PhoneNorm       string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_cust_phone_cc" json:"-"`
	Email           string         `gorm:"type:varchar(100)" json:"email"`
	Address         string         `gorm:"type:text" json:"address"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}
