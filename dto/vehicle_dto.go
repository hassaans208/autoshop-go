package dto

type CreateCarRequest struct {
	CustomerID string `json:"customerId" binding:"required"`
	Plate      string `json:"plate" binding:"required"`
	Chassis    string `json:"chassis"`
	Company    string `json:"company"`
	Model      string `json:"model"`
	Year       int    `json:"year"`
	Color      string `json:"color"`
	ColorLabel string `json:"colorLabel"`
	Mileage    int    `json:"mileage"`
	Notes      string `json:"notes"`
}

type UpdateCarRequest struct {
	Plate      string `json:"plate"`
	Chassis    string `json:"chassis"`
	Company    string `json:"company"`
	Model      string `json:"model"`
	Year       int    `json:"year"`
	Color      string `json:"color"`
	ColorLabel string `json:"colorLabel"`
	Mileage    int    `json:"mileage"`
	Notes      string `json:"notes"`
}
