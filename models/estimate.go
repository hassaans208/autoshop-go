package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Estimate struct {
	ID                 string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	ClientCompanyID    string         `gorm:"type:varchar(36);not null;index;uniqueIndex:idx_est_num_cc" json:"-"`
	Number             string         `gorm:"type:varchar(30);not null;uniqueIndex:idx_est_num_cc" json:"number"`
	PublicSlug         string         `gorm:"type:varchar(8);uniqueIndex" json:"publicSlug"`
	Date               time.Time      `gorm:"not null" json:"date"`
	CustomerID         string         `gorm:"type:varchar(36);not null;index" json:"customerId"`
	CarID              string         `gorm:"type:varchar(36);not null;index" json:"carId"`
	JobTypes           string         `gorm:"type:json" json:"-"`
	Subtypes           string         `gorm:"type:json" json:"-"`
	MechanicID         *string        `gorm:"type:varchar(36);index" json:"mechanicId"`
	Status             string         `gorm:"type:enum('Draft','Sent','Approved','Rejected','Converted');default:'Draft'" json:"status"`
	Notes              string         `gorm:"type:text" json:"notes"`
	DiscountAmount     float64        `gorm:"type:decimal(12,2);default:0" json:"discountAmount"`
	TaxPct             float64        `gorm:"type:decimal(5,2);default:0" json:"taxPct"`
	TaxAmount          float64        `gorm:"type:decimal(12,2);default:0" json:"taxAmount"`
	Subtotal           float64        `gorm:"type:decimal(12,2);default:0" json:"subtotal"`
	Tax                float64        `gorm:"type:decimal(12,2);default:0" json:"tax"`
	Grand              float64        `gorm:"type:decimal(12,2);default:0" json:"grand"`
	ConvertedInvoiceID *string        `gorm:"type:varchar(36)" json:"convertedInvoiceId"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	Items  []EstimateItem  `gorm:"foreignKey:EstimateID" json:"items"`
	Images []EstimateImage `gorm:"foreignKey:EstimateID" json:"images"`
}

type EstimateItem struct {
	ID              string  `gorm:"type:varchar(36);primaryKey" json:"id"`
	EstimateID      string  `gorm:"type:varchar(36);not null;index" json:"estimateId"`
	Name            string  `gorm:"type:varchar(200);not null" json:"name"`
	Description     string  `gorm:"type:text" json:"description"`
	Qty             float64 `gorm:"type:decimal(10,3);not null" json:"qty"`
	UnitPrice       float64 `gorm:"type:decimal(12,2);not null" json:"unitPrice"`
	DiscountPct     float64 `gorm:"type:decimal(5,2);default:0" json:"discountPct"`
	DiscountAmount  float64 `gorm:"type:decimal(12,2);default:0" json:"discountAmount"`
	Kind            string  `gorm:"type:varchar(20)" json:"kind"`
	PurchasePrice   float64 `gorm:"type:decimal(12,2);default:0" json:"purchasePrice,omitempty"`
	VendorID        *string `gorm:"type:varchar(36)" json:"vendorId,omitempty"`
	JobType         string  `gorm:"type:varchar(100)" json:"jobType,omitempty"`
	Subtype         string  `gorm:"type:varchar(100)" json:"subtype,omitempty"`
	Assignees       string  `gorm:"type:json" json:"assignees,omitempty"`
	MaterialPrice   float64 `gorm:"type:decimal(12,2);default:0" json:"materialPrice,omitempty"`
	ServicePrice    float64 `gorm:"type:decimal(12,2);default:0" json:"servicePrice,omitempty"`
	SortOrder       int     `gorm:"default:0" json:"sortOrder"`
}

type EstimateImage struct {
	ID         string `gorm:"type:varchar(36);primaryKey" json:"id"`
	EstimateID string `gorm:"type:varchar(36);not null;index" json:"estimateId"`
	Label      string `gorm:"type:varchar(150)" json:"label"`
	DataURL    string `gorm:"type:longtext;not null" json:"dataUrl"`
	SortOrder  int    `gorm:"default:0" json:"sortOrder"`
}

func (e *Estimate) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	return nil
}

func (e *EstimateItem) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	return nil
}

func (e *EstimateImage) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	return nil
}
