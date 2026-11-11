package dto

import "time"

type CreateEstimateRequest struct {
	Date           time.Time      `json:"date" binding:"required"`
	CustomerID     string         `json:"customerId" binding:"required"`
	CarID          string         `json:"carId" binding:"required"`
	JobTypes       []string       `json:"jobTypes"`
	Subtypes       []string       `json:"subtypes"`
	MechanicID     *string        `json:"mechanicId"`
	Status         string         `json:"status"`
	Notes          string         `json:"notes"`
	Items          []LineItemDTO  `json:"items"`
	Images         []ImageItemDTO `json:"images"`
	DiscountAmount float64        `json:"discountAmount"`
	TaxPct         float64        `json:"taxPct"`
	TaxAmount      float64        `json:"taxAmount"`
}

type UpdateEstimateRequest struct {
	Date               *time.Time `json:"date"`
	JobTypes           []string   `json:"jobTypes"`
	Subtypes           []string   `json:"subtypes"`
	MechanicID         *string    `json:"mechanicId"`
	Status             string     `json:"status"`
	Notes              string     `json:"notes"`
	Items              []LineItemDTO  `json:"items"`
	Images             []ImageItemDTO `json:"images"`
	DiscountAmount     *float64   `json:"discountAmount"`
	TaxPct             *float64   `json:"taxPct"`
	TaxAmount          *float64   `json:"taxAmount"`
	ConvertedInvoiceID *string    `json:"convertedInvoiceId"`
}

type ConvertEstimateRequest struct {
	DueDate      *time.Time `json:"dueDate"`
	DeliveryDate *time.Time `json:"deliveryDate"`
	FinalMileage int        `json:"finalMileage"`
	Warranty     string     `json:"warranty"`
}
