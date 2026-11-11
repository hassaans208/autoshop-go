package models

import "time"

type Settings struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientCompanyID string    `gorm:"type:varchar(36);uniqueIndex;not null" json:"-"`
	WorkshopName    string    `gorm:"type:varchar(255)" json:"workshopName"`
	Logo            string    `gorm:"type:longtext" json:"logo"`
	Tagline         string    `gorm:"type:varchar(255)" json:"tagline"`
	Address         string    `gorm:"type:text" json:"address"`
	Phone           string    `gorm:"type:varchar(20)" json:"phone"`
	Email           string    `gorm:"type:varchar(100)" json:"email"`
	Whatsapp        string    `gorm:"type:varchar(20)" json:"whatsapp"`
	NTN             string    `gorm:"type:varchar(30)" json:"ntn"`
	DefaultTaxPct   float64   `gorm:"type:decimal(5,2);default:0" json:"defaultTaxPct"`
	Currency        string    `gorm:"type:varchar(10);default:'PKR'" json:"currency"`
	Terms           string    `gorm:"type:text" json:"terms"`
	WarrantyDefault string    `gorm:"type:varchar(255)" json:"warrantyDefault"`
	Footer          string    `gorm:"type:varchar(255)" json:"footer"`
	EstPrefix       string    `gorm:"type:varchar(10);default:'EST'" json:"estPrefix"`
	EstStart        int       `gorm:"default:1" json:"estStart"`
	InvPrefix       string    `gorm:"type:varchar(10);default:'INV'" json:"invPrefix"`
	InvStart        int       `gorm:"default:1" json:"invStart"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type DocCounter struct {
	ID              uint   `gorm:"primaryKey;autoIncrement"`
	ClientCompanyID string `gorm:"type:varchar(36);not null;uniqueIndex:idx_counter_cc"`
	Prefix          string `gorm:"type:varchar(10);not null;uniqueIndex:idx_counter_cc"`
	Year            int    `gorm:"type:smallint;not null;uniqueIndex:idx_counter_cc"`
	Last            int    `gorm:"not null;default:0"`
}
