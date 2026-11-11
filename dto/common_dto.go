package dto

type LineItemDTO struct {
	ID             string   `json:"id"`
	Name           string   `json:"name" binding:"required"`
	Description    string   `json:"description"`
	Qty            float64  `json:"qty" binding:"required"`
	UnitPrice      float64  `json:"unitPrice" binding:"required"`
	DiscountPct    float64  `json:"discountPct"`
	DiscountAmount float64  `json:"discountAmount"`
	Kind           string   `json:"kind"`
	PurchasePrice  float64  `json:"purchasePrice"`
	VendorID       *string  `json:"vendorId"`
	JobType        string   `json:"jobType"`
	Subtype        string   `json:"subtype"`
	Assignees      FlexibleStringSlice `json:"assignees"`
	MaterialPrice  float64  `json:"materialPrice"`
	ServicePrice   float64  `json:"servicePrice"`
	SortOrder      int      `json:"sortOrder"`
}

type ImageItemDTO struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	DataURL   string `json:"dataUrl"`
	SortOrder int    `json:"sortOrder"`
}
