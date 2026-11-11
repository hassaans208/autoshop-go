package migrations

import (
	"time"

	"gorm.io/gorm"
)

// ─── Staff ────────────────────────────────────────────────────────────────────
type Staff struct {
	ID           string `gorm:"type:varchar(36);primaryKey"`
	Name         string `gorm:"type:varchar(100);not null"`
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Role         string `gorm:"type:enum('owner','manager','mechanic','receptionist','accountant');not null"`
	Phone        string `gorm:"type:varchar(20)"`
	Active       bool   `gorm:"default:true;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// ─── Settings (single row, id always = 1) ────────────────────────────────────
type Settings struct {
	ID              uint    `gorm:"primaryKey;autoIncrement"`
	WorkshopName    string  `gorm:"type:varchar(150)"`
	Logo            string  `gorm:"type:longtext"` // base64 or URL
	Tagline         string  `gorm:"type:varchar(255)"`
	Address         string  `gorm:"type:text"`
	Phone           string  `gorm:"type:varchar(20)"`
	Email           string  `gorm:"type:varchar(100)"`
	Whatsapp        string  `gorm:"type:varchar(20)"`
	NTN             string  `gorm:"type:varchar(30)"`
	DefaultTaxPct   float64 `gorm:"type:decimal(5,2);default:0"`
	Currency        string  `gorm:"type:varchar(10);default:'PKR'"`
	Terms           string  `gorm:"type:text"`
	WarrantyDefault string  `gorm:"type:varchar(255)"`
	Footer          string  `gorm:"type:varchar(255)"`
	EstPrefix       string  `gorm:"type:varchar(10);default:'EST'"`
	EstStart        int     `gorm:"default:1"`
	InvPrefix       string  `gorm:"type:varchar(10);default:'INV'"`
	InvStart        int     `gorm:"default:1"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ─── Company ──────────────────────────────────────────────────────────────────
type Company struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// ─── JobType ──────────────────────────────────────────────────────────────────
type JobType struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(100);uniqueIndex;not null"`
}

// ─── Customer ─────────────────────────────────────────────────────────────────
type Customer struct {
	ID        string `gorm:"type:varchar(36);primaryKey"`
	Name      string `gorm:"type:varchar(100);not null"`
	Phone     string `gorm:"type:varchar(20);uniqueIndex;not null"`
	PhoneNorm string `gorm:"type:varchar(20);uniqueIndex;not null"` // digits only, for search
	Email     string `gorm:"type:varchar(100)"`
	Address   string `gorm:"type:text"`
	Notes     string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// ─── Car ──────────────────────────────────────────────────────────────────────
type Car struct {
	ID         string `gorm:"type:varchar(36);primaryKey"`
	CustomerID string `gorm:"type:varchar(36);not null;index"`
	Plate      string `gorm:"type:varchar(20);uniqueIndex;not null"` // stored UPPERCASE
	Chassis    string `gorm:"type:varchar(50)"`
	Company    string `gorm:"type:varchar(100)"` // denormalised — just store the name
	Model      string `gorm:"type:varchar(100)"`
	Year       int    `gorm:"type:smallint"`
	Color      string `gorm:"type:varchar(20)"` // hex
	ColorLabel string `gorm:"type:varchar(50)"`
	Mileage    int    `gorm:"type:int;default:0"`
	Notes      string `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// ─── Estimate ─────────────────────────────────────────────────────────────────
type Estimate struct {
	ID                 string    `gorm:"type:varchar(36);primaryKey"`
	Number             string    `gorm:"type:varchar(30);uniqueIndex;not null"` // EST-2026-00001
	Date               time.Time `gorm:"not null"`
	CustomerID         string    `gorm:"type:varchar(36);not null;index"`
	CarID              string    `gorm:"type:varchar(36);not null;index"`
	JobTypes           string    `gorm:"type:json"` // JSON array of strings
	MechanicID         *string   `gorm:"type:varchar(36);index"`
	Status             string    `gorm:"type:enum('Draft','Sent','Approved','Rejected','Converted');default:'Draft'"`
	Notes              string    `gorm:"type:text"`
	DiscountAmount     float64   `gorm:"type:decimal(12,2);default:0"`
	TaxPct             float64   `gorm:"type:decimal(5,2);default:0"`
	Subtotal           float64   `gorm:"type:decimal(12,2);default:0"` // server computed
	Tax                float64   `gorm:"type:decimal(12,2);default:0"`
	Grand              float64   `gorm:"type:decimal(12,2);default:0"`
	ConvertedInvoiceID *string   `gorm:"type:varchar(36)"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type EstimateItem struct {
	ID          string  `gorm:"type:varchar(36);primaryKey"`
	EstimateID  string  `gorm:"type:varchar(36);not null;index"`
	Name        string  `gorm:"type:varchar(200);not null"`
	Description string  `gorm:"type:text"`
	Qty         float64 `gorm:"type:decimal(10,3);not null"`
	UnitPrice   float64 `gorm:"type:decimal(12,2);not null"`
	DiscountPct float64 `gorm:"type:decimal(5,2);default:0"`
	SortOrder   int     `gorm:"default:0"`
}

type EstimateImage struct {
	ID         string `gorm:"type:varchar(36);primaryKey"`
	EstimateID string `gorm:"type:varchar(36);not null;index"`
	Label      string `gorm:"type:varchar(150)"`
	DataURL    string `gorm:"type:longtext;not null"` // base64 or CDN URL
	SortOrder  int    `gorm:"default:0"`
}

// ─── Invoice ──────────────────────────────────────────────────────────────────
type Invoice struct {
	ID             string    `gorm:"type:varchar(36);primaryKey"`
	Number         string    `gorm:"type:varchar(30);uniqueIndex;not null"` // INV-2026-00001
	Date           time.Time `gorm:"not null"`
	DueDate        *time.Time
	CustomerID     string  `gorm:"type:varchar(36);not null;index"`
	CarID          string  `gorm:"type:varchar(36);not null;index"`
	JobTypes       string  `gorm:"type:json"`
	MechanicID     *string `gorm:"type:varchar(36);index"`
	Status         string  `gorm:"type:enum('Unpaid','Partially Paid','Paid','Overdue');default:'Unpaid'"`
	Notes          string  `gorm:"type:text"`
	DiscountAmount float64 `gorm:"type:decimal(12,2);default:0"`
	TaxPct         float64 `gorm:"type:decimal(5,2);default:0"`
	Subtotal       float64 `gorm:"type:decimal(12,2);default:0"`
	Tax            float64 `gorm:"type:decimal(12,2);default:0"`
	Grand          float64 `gorm:"type:decimal(12,2);default:0"`
	Paid           float64 `gorm:"type:decimal(12,2);default:0"`
	Balance        float64 `gorm:"type:decimal(12,2);default:0"`
	DeliveryDate   *time.Time
	FinalMileage   int     `gorm:"default:0"`
	Warranty       string  `gorm:"type:varchar(255)"`
	FromEstimateID *string `gorm:"type:varchar(36)"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type InvoiceItem struct {
	ID          string  `gorm:"type:varchar(36);primaryKey"`
	InvoiceID   string  `gorm:"type:varchar(36);not null;index"`
	Name        string  `gorm:"type:varchar(200);not null"`
	Description string  `gorm:"type:text"`
	Qty         float64 `gorm:"type:decimal(10,3);not null"`
	UnitPrice   float64 `gorm:"type:decimal(12,2);not null"`
	DiscountPct float64 `gorm:"type:decimal(5,2);default:0"`
	SortOrder   int     `gorm:"default:0"`
}

type InvoiceImage struct {
	ID        string `gorm:"type:varchar(36);primaryKey"`
	InvoiceID string `gorm:"type:varchar(36);not null;index"`
	Label     string `gorm:"type:varchar(150)"`
	DataURL   string `gorm:"type:longtext;not null"`
	SortOrder int    `gorm:"default:0"`
}

// ─── Payment ──────────────────────────────────────────────────────────────────
type Payment struct {
	ID        string    `gorm:"type:varchar(36);primaryKey"`
	InvoiceID string    `gorm:"type:varchar(36);not null;index"`
	Date      time.Time `gorm:"not null"`
	Amount    float64   `gorm:"type:decimal(12,2);not null"`
	Method    string    `gorm:"type:enum('Cash','Bank Transfer','Cheque','Card','Split');not null"`
	Note      string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time
}

// ─── Number counter (atomic sequence per prefix per year) ─────────────────────
type DocCounter struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Prefix string `gorm:"type:varchar(10);not null"`
	Year   int    `gorm:"type:smallint;not null"`
	Last   int    `gorm:"not null;default:0"`
}

func init() {
	// composite unique so EST+2026 and INV+2026 are separate rows
}
