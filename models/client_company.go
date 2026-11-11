package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ClientCompany is the tenant root — each workshop belongs to one client company.
type ClientCompany struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	OwnerID   string         `gorm:"type:varchar(36);not null;index" json:"ownerId"`
	Active    bool           `gorm:"default:true;not null" json:"active"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *ClientCompany) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}
