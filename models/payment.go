package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID              string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ClientCompanyID string    `gorm:"type:varchar(36);not null;index" json:"-"`
	InvoiceID       string    `gorm:"type:varchar(36);not null;index" json:"invoiceId"`
	Date            time.Time `gorm:"not null" json:"date"`
	Amount          float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Method          string    `gorm:"type:enum('Cash','Bank Transfer','Cheque','Card','Split');not null" json:"method"`
	Note            string    `gorm:"type:varchar(255)" json:"note"`
	CreatedAt       time.Time `json:"createdAt"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}
