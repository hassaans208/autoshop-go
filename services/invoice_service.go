package services

import (
	"autoworkshop/dto"
	"autoworkshop/models"
	"autoworkshop/repositories"
	"autoworkshop/utils"
	"encoding/json"
	"errors"
	"math"

	"gorm.io/gorm"
)

type InvoiceService interface {
	List(tenantID string, p utils.PageParams, filters map[string]string) ([]models.Invoice, int64, error)
	GetByID(tenantID, id string) (*models.Invoice, error)
	GetByCustomer(tenantID, customerID string) ([]models.Invoice, error)
	GetByCar(tenantID, carID string) ([]models.Invoice, error)
	Create(tenantID string, req dto.CreateInvoiceRequest, settings *models.Settings) (*models.Invoice, error)
	Update(tenantID, id string, req dto.UpdateInvoiceRequest) (*models.Invoice, error)
	Delete(tenantID, id string) error
	AddPayment(tenantID, invoiceID string, req dto.CreatePaymentRequest) (*models.Payment, *models.Invoice, error)
	DeletePayment(tenantID, invoiceID, paymentID string) (*models.Invoice, error)
}

type invoiceService struct {
	repo        repositories.InvoiceRepository
	paymentRepo repositories.PaymentRepository
}

func NewInvoiceService(repo repositories.InvoiceRepository, payRepo repositories.PaymentRepository) InvoiceService {
	return &invoiceService{repo, payRepo}
}

func (s *invoiceService) List(tenantID string, p utils.PageParams, filters map[string]string) ([]models.Invoice, int64, error) {
	return s.repo.List(tenantID, p, filters)
}

func (s *invoiceService) GetByID(tenantID, id string) (*models.Invoice, error) {
	return s.repo.FindByID(tenantID, id)
}

func (s *invoiceService) GetByCustomer(tenantID, id string) ([]models.Invoice, error) {
	return s.repo.FindByCustomer(tenantID, id)
}

func (s *invoiceService) GetByCar(tenantID, id string) ([]models.Invoice, error) {
	return s.repo.FindByCar(tenantID, id)
}

func (s *invoiceService) Create(tenantID string, req dto.CreateInvoiceRequest, settings *models.Settings) (*models.Invoice, error) {
	year := req.Date.Year()
	number, err := s.repo.NextNumber(tenantID, settings.InvPrefix, year)
	if err != nil {
		return nil, err
	}

	jtJSON, _ := json.Marshal(req.JobTypes)
	stJSON, _ := json.Marshal(req.Subtypes)

	subtotal, tax, grand := computeTotals(req.Items, req.DiscountAmount, req.TaxPct, req.TaxAmount)

	inv := &models.Invoice{
		ClientCompanyID: tenantID,
		Number:          number,
		PublicSlug:      utils.PublicSlug(),
		Date:            req.Date,
		DueDate:         req.DueDate,
		CustomerID:      req.CustomerID,
		CarID:           req.CarID,
		JobTypes:        string(jtJSON),
		Subtypes:        string(stJSON),
		MechanicID:      req.MechanicID,
		Notes:           req.Notes,
		DiscountAmount:  req.DiscountAmount,
		TaxPct:          req.TaxPct,
		TaxAmount:       req.TaxAmount,
		Subtotal:        subtotal,
		Tax:             tax,
		Grand:           grand,
		Paid:            0,
		Balance:         grand,
		DeliveryDate:    req.DeliveryDate,
		FinalMileage:    req.FinalMileage,
		Warranty:        req.Warranty,
		FromEstimateID:  req.FromEstimateID,
		Items:           mapInvoiceItems(req.Items, ""),
		Images:          mapInvoiceImages(req.Images),
	}
	inv.Status = computeInvoiceStatus(inv)
	return inv, s.repo.Create(inv)
}

func (s *invoiceService) Update(tenantID, id string, req dto.UpdateInvoiceRequest) (*models.Invoice, error) {
	inv, err := s.repo.FindByID(tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	if req.Date != nil {
		inv.Date = *req.Date
	}
	if req.DueDate != nil {
		inv.DueDate = req.DueDate
	}
	if req.JobTypes != nil {
		jtJSON, _ := json.Marshal(req.JobTypes)
		inv.JobTypes = string(jtJSON)
	}
	if req.Subtypes != nil {
		stJSON, _ := json.Marshal(req.Subtypes)
		inv.Subtypes = string(stJSON)
	}
	if req.MechanicID != nil {
		inv.MechanicID = req.MechanicID
	}
	if req.Notes != "" {
		inv.Notes = req.Notes
	}
	if req.DiscountAmount != nil {
		inv.DiscountAmount = *req.DiscountAmount
	}
	if req.TaxPct != nil {
		inv.TaxPct = *req.TaxPct
	}
	if req.TaxAmount != nil {
		inv.TaxAmount = *req.TaxAmount
	}
	if req.DeliveryDate != nil {
		inv.DeliveryDate = req.DeliveryDate
	}
	if req.FinalMileage != nil {
		inv.FinalMileage = *req.FinalMileage
	}
	if req.Warranty != "" {
		inv.Warranty = req.Warranty
	}
	if req.Items != nil {
		inv.Items = mapInvoiceItems(req.Items, inv.ID)
	}
	if req.Images != nil {
		inv.Images = mapInvoiceImages(req.Images)
		for i := range inv.Images {
			inv.Images[i].InvoiceID = inv.ID
		}
	}

	inv.Subtotal, inv.Tax, inv.Grand = recomputeFromInvoiceItems(inv.Items, inv.DiscountAmount, inv.TaxPct, inv.TaxAmount)
	inv.Paid = sumPayments(inv.Payments)
	inv.Balance = math.Max(0, inv.Grand-inv.Paid)
	inv.Status = computeInvoiceStatus(inv)

	return inv, s.repo.Update(inv)
}

func (s *invoiceService) Delete(tenantID, id string) error {
	return s.repo.Delete(tenantID, id)
}

func (s *invoiceService) AddPayment(tenantID, invoiceID string, req dto.CreatePaymentRequest) (*models.Payment, *models.Invoice, error) {
	inv, err := s.repo.FindByID(tenantID, invoiceID)
	if err != nil {
		return nil, nil, errors.New("invoice not found")
	}

	p := &models.Payment{
		ClientCompanyID: tenantID,
		InvoiceID:       invoiceID,
		Date:            req.Date,
		Amount:          req.Amount,
		Method:          req.Method,
		Note:            req.Note,
	}
	if err := s.paymentRepo.Create(p); err != nil {
		return nil, nil, err
	}

	inv.Payments = append(inv.Payments, *p)
	inv.Paid = sumPayments(inv.Payments)
	inv.Balance = math.Max(0, inv.Grand-inv.Paid)
	inv.Status = computeInvoiceStatus(inv)

	if err := s.repo.Update(inv); err != nil {
		return nil, nil, err
	}
	return p, inv, nil
}

func (s *invoiceService) DeletePayment(tenantID, invoiceID, paymentID string) (*models.Invoice, error) {
	inv, err := s.repo.FindByID(tenantID, invoiceID)
	if err != nil {
		return nil, errors.New("invoice not found")
	}
	if err := s.paymentRepo.Delete(paymentID); err != nil {
		return nil, err
	}
	remaining := make([]models.Payment, 0)
	for _, p := range inv.Payments {
		if p.ID != paymentID {
			remaining = append(remaining, p)
		}
	}
	inv.Payments = remaining
	inv.Paid = sumPayments(inv.Payments)
	inv.Balance = math.Max(0, inv.Grand-inv.Paid)
	inv.Status = computeInvoiceStatus(inv)
	return inv, s.repo.Update(inv)
}
