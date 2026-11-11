package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	ID              string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	ClientCompanyID string         `gorm:"type:varchar(36);not null;index;uniqueIndex:idx_inv_num_cc" json:"-"`
	Number          string         `gorm:"type:varchar(30);not null;uniqueIndex:idx_inv_num_cc" json:"number"`
	PublicSlug      string         `gorm:"type:varchar(8);uniqueIndex" json:"publicSlug"`
	Date            time.Time      `gorm:"not null" json:"date"`
	DueDate         *time.Time     `json:"dueDate"`
	CustomerID      string         `gorm:"type:varchar(36);not null;index" json:"customerId"`
	CarID           string         `gorm:"type:varchar(36);not null;index" json:"carId"`
	JobTypes        string         `gorm:"type:json" json:"-"`
	Subtypes        string         `gorm:"type:json" json:"-"`
	MechanicID      *string        `gorm:"type:varchar(36);index" json:"mechanicId"`
	Status          string         `gorm:"type:enum('Unpaid','Partially Paid','Paid','Overdue');default:'Unpaid'" json:"status"`
	Notes           string         `gorm:"type:text" json:"notes"`
	DiscountAmount  float64        `gorm:"type:decimal(12,2);default:0" json:"discountAmount"`
	TaxPct          float64        `gorm:"type:decimal(5,2);default:0" json:"taxPct"`
	TaxAmount       float64        `gorm:"type:decimal(12,2);default:0" json:"taxAmount"`
	Subtotal        float64        `gorm:"type:decimal(12,2);default:0" json:"subtotal"`
	Tax             float64        `gorm:"type:decimal(12,2);default:0" json:"tax"`
	Grand           float64        `gorm:"type:decimal(12,2);default:0" json:"grand"`
	Paid            float64        `gorm:"type:decimal(12,2);default:0" json:"paid"`
	Balance         float64        `gorm:"type:decimal(12,2);default:0" json:"balance"`
	DeliveryDate    *time.Time     `json:"deliveryDate"`
	FinalMileage    int            `gorm:"default:0" json:"finalMileage"`
	Warranty        string         `gorm:"type:varchar(255)" json:"warranty"`
	FromEstimateID  *string        `gorm:"type:varchar(36)" json:"fromEstimateId"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Items    []InvoiceItem  `gorm:"foreignKey:InvoiceID" json:"items"`
	Images   []InvoiceImage `gorm:"foreignKey:InvoiceID" json:"images"`
	Payments []Payment      `gorm:"foreignKey:InvoiceID" json:"payments"`
}

type InvoiceItem struct {
	ID             string  `gorm:"type:varchar(36);primaryKey" json:"id"`
	InvoiceID      string  `gorm:"type:varchar(36);not null;index" json:"invoiceId"`
	Name           string  `gorm:"type:varchar(200);not null" json:"name"`
	Description    string  `gorm:"type:text" json:"description"`
	Qty            float64 `gorm:"type:decimal(10,3);not null" json:"qty"`
	UnitPrice      float64 `gorm:"type:decimal(12,2);not null" json:"unitPrice"`
	DiscountPct    float64 `gorm:"type:decimal(5,2);default:0" json:"discountPct"`
	DiscountAmount float64 `gorm:"type:decimal(12,2);default:0" json:"discountAmount"`
	Kind           string  `gorm:"type:varchar(20)" json:"kind"`
	PurchasePrice  float64 `gorm:"type:decimal(12,2);default:0" json:"purchasePrice,omitempty"`
	VendorID       *string `gorm:"type:varchar(36)" json:"vendorId,omitempty"`
	JobType        string  `gorm:"type:varchar(100)" json:"jobType,omitempty"`
	Subtype        string  `gorm:"type:varchar(100)" json:"subtype,omitempty"`
	Assignees      string  `gorm:"type:json" json:"assignees,omitempty"`
	MaterialPrice  float64 `gorm:"type:decimal(12,2);default:0" json:"materialPrice,omitempty"`
	ServicePrice   float64 `gorm:"type:decimal(12,2);default:0" json:"servicePrice,omitempty"`
	SortOrder      int     `gorm:"default:0" json:"sortOrder"`
}

type InvoiceImage struct {
	ID        string `gorm:"type:varchar(36);primaryKey" json:"id"`
	InvoiceID string `gorm:"type:varchar(36);not null;index" json:"invoiceId"`
	Label     string `gorm:"type:varchar(150)" json:"label"`
	DataURL   string `gorm:"type:longtext;not null" json:"dataUrl"`
	SortOrder int    `gorm:"default:0" json:"sortOrder"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	return nil
}

func (i *InvoiceItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	return nil
}

func (i *InvoiceImage) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	return nil
}
