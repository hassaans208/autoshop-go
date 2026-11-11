package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Car struct {
	ID              string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	ClientCompanyID string         `gorm:"type:varchar(36);not null;index;uniqueIndex:idx_car_plate_cc" json:"-"`
	CustomerID      string         `gorm:"type:varchar(36);not null;index" json:"customerId"`
	Plate           string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_car_plate_cc" json:"plate"`
	Chassis         string         `gorm:"type:varchar(50)" json:"chassis"`
	Company         string         `gorm:"type:varchar(100)" json:"company"`
	Model           string         `gorm:"type:varchar(100)" json:"model"`
	Year            int            `gorm:"type:smallint" json:"year"`
	Color           string         `gorm:"type:varchar(20)" json:"color"`
	ColorLabel      string         `gorm:"type:varchar(50)" json:"colorLabel"`
	Mileage         int            `gorm:"type:int;default:0" json:"mileage"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Car) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}
