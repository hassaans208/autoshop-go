package dto

import "time"

type CreateInvoiceRequest struct {
	Date           time.Time      `json:"date" binding:"required"`
	DueDate        *time.Time     `json:"dueDate"`
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
	DeliveryDate   *time.Time     `json:"deliveryDate"`
	FinalMileage   int            `json:"finalMileage"`
	Warranty       string         `json:"warranty"`
	FromEstimateID *string        `json:"fromEstimateId"`
}

type UpdateInvoiceRequest struct {
	Date           *time.Time     `json:"date"`
	DueDate        *time.Time     `json:"dueDate"`
	JobTypes       []string       `json:"jobTypes"`
	Subtypes       []string       `json:"subtypes"`
	MechanicID     *string        `json:"mechanicId"`
	Status         string         `json:"status"`
	Notes          string         `json:"notes"`
	Items          []LineItemDTO  `json:"items"`
	Images         []ImageItemDTO `json:"images"`
	DiscountAmount *float64       `json:"discountAmount"`
	TaxPct         *float64       `json:"taxPct"`
	TaxAmount      *float64       `json:"taxAmount"`
	DeliveryDate   *time.Time     `json:"deliveryDate"`
	FinalMileage   *int           `json:"finalMileage"`
	Warranty       string         `json:"warranty"`
}

type CreatePaymentRequest struct {
	Date   time.Time `json:"date" binding:"required"`
	Amount float64   `json:"amount" binding:"required,gt=0"`
	Method string    `json:"method" binding:"required,oneof=Cash 'Bank Transfer' Cheque Card Split"`
	Note   string    `json:"note"`
}
