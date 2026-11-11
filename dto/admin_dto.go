package dto

type AdminCreateOwnerRequest struct {
	Name     string  `json:"name" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Phone    string  `json:"phone"`
	Salary   float64 `json:"salary"`
}

type AdminCreateCompanyRequest struct {
	CompanyName  string                  `json:"companyName" binding:"required"`
	WorkshopName string                  `json:"workshopName"`
	Owner        AdminCreateOwnerRequest `json:"owner" binding:"required"`
}

type AdminCompanyResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"ownerId"`
	Active  bool   `json:"active"`
}

type AdminCreateCompanyResponse struct {
	Company AdminCompanyResponse `json:"company"`
	Owner   UserResponse         `json:"owner"`
}
