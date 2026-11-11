package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vendor struct {
	ID              string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	ClientCompanyID string         `gorm:"type:varchar(36);not null;index" json:"-"`
	Name            string         `gorm:"type:varchar(150);not null" json:"name"`
	Kind            string         `gorm:"type:enum('company','vendor');not null;default:'vendor'" json:"kind"`
	Phone           string         `gorm:"type:varchar(30)" json:"phone"`
	Email           string         `gorm:"type:varchar(100)" json:"email"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (v *Vendor) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.NewString()
	}
	return nil
}

type Part struct {
	ID                   string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	ClientCompanyID      string         `gorm:"type:varchar(36);not null;index" json:"-"`
	Name                 string         `gorm:"type:varchar(200);not null" json:"name"`
	VendorID             *string        `gorm:"type:varchar(36);index" json:"vendorId"`
	DefaultPurchasePrice float64        `gorm:"type:decimal(12,2);default:0" json:"defaultPurchasePrice"`
	DefaultSalePrice     float64        `gorm:"type:decimal(12,2);default:0" json:"defaultSalePrice"`
	Notes                string         `gorm:"type:text" json:"notes"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Part) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}
