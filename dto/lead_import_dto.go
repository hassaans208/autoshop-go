package dto

// LeadImportRow maps one spreadsheet row (workshop lead → tenant + owner).
type LeadImportRow struct {
	BusinessName string `json:"businessName" binding:"required"`
	City         string `json:"city"`
	Country      string `json:"country"`
	Category     string `json:"category"`
	Phone        string `json:"phone"`
	WhatsApp     string `json:"whatsapp"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password" binding:"required"`
	Website      string `json:"website"`
}

type LeadImportRequest struct {
	Rows []LeadImportRow `json:"rows" binding:"required,min=1,dive"`
}

type LeadImportResultRow struct {
	BusinessName string `json:"businessName"`
	LoginEmail   string `json:"loginEmail"`
	CompanyID    string `json:"companyId,omitempty"`
	StaffID      string `json:"staffId,omitempty"`
	AutologinURL string `json:"autologinUrl,omitempty"`
	Status       string `json:"status"`
	Error        string `json:"error,omitempty"`
}

type LeadImportResponse struct {
	Results []LeadImportResultRow `json:"results"`
	Created int                   `json:"created"`
	Updated int                   `json:"updated"`
	Failed  int                   `json:"failed"`
}
