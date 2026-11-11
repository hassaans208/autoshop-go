package dto

type CreateStaffRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Name     string  `json:"name" binding:"required"`
	Role     string  `json:"role" binding:"required,oneof=owner manager mechanic receptionist accountant"`
	Phone    string  `json:"phone"`
	Salary   float64 `json:"salary"`
	Active   *bool   `json:"active"`
}

type UpdateStaffRequest struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Role     string  `json:"role" binding:"omitempty,oneof=owner manager mechanic receptionist accountant"`
	Phone    string  `json:"phone"`
	Salary   *float64 `json:"salary"`
	Active   *bool   `json:"active"`
	Password string  `json:"password"`
}

type CreateVendorRequest struct {
	Name  string `json:"name" binding:"required"`
	Kind  string `json:"kind" binding:"required,oneof=company vendor"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Notes string `json:"notes"`
}

type UpdateVendorRequest struct {
	Name  string `json:"name"`
	Kind  string `json:"kind" binding:"omitempty,oneof=company vendor"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Notes string `json:"notes"`
}

type CreatePartRequest struct {
	Name                 string  `json:"name" binding:"required"`
	VendorID             *string `json:"vendorId"`
	DefaultPurchasePrice float64 `json:"defaultPurchasePrice"`
	DefaultSalePrice     float64 `json:"defaultSalePrice"`
	Notes                string  `json:"notes"`
}

type UpdatePartRequest struct {
	Name                 string  `json:"name"`
	VendorID             *string `json:"vendorId"`
	DefaultPurchasePrice *float64 `json:"defaultPurchasePrice"`
	DefaultSalePrice     *float64 `json:"defaultSalePrice"`
	Notes                string  `json:"notes"`
}

type JobTypeEntry struct {
	Name     string   `json:"name" binding:"required"`
	Subtypes []string `json:"subtypes"`
}

type ReplaceJobTypesRequest struct {
	List     []JobTypeEntry `json:"list"`
	JobTypes []JobTypeEntry `json:"jobTypes"`
}

func (r ReplaceJobTypesRequest) Entries() []JobTypeEntry {
	if len(r.List) > 0 {
		return r.List
	}
	return r.JobTypes
}

type AddCompanyRequest struct {
	Name string `json:"name" binding:"required"`
}
