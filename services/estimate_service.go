package services

import (
	"autoworkshop/dto"
	"autoworkshop/models"
	"autoworkshop/repositories"
	"autoworkshop/utils"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type EstimateService interface {
	List(tenantID string, p utils.PageParams, filters map[string]string) ([]models.Estimate, int64, error)
	GetByID(tenantID, id string) (*models.Estimate, error)
	GetByCar(tenantID, carID string) ([]models.Estimate, error)
	Create(tenantID string, req dto.CreateEstimateRequest, settings *models.Settings) (*models.Estimate, error)
	Update(tenantID, id string, req dto.UpdateEstimateRequest) (*models.Estimate, error)
	Delete(tenantID, id string) error
	Convert(tenantID, id string, req dto.ConvertEstimateRequest, settings *models.Settings) (*models.Invoice, error)
}

type estimateService struct {
	repo    repositories.EstimateRepository
	invRepo repositories.InvoiceRepository
}

func NewEstimateService(repo repositories.EstimateRepository, invRepo repositories.InvoiceRepository) EstimateService {
	return &estimateService{repo, invRepo}
}

func (s *estimateService) List(tenantID string, p utils.PageParams, filters map[string]string) ([]models.Estimate, int64, error) {
	return s.repo.List(tenantID, p, filters)
}

func (s *estimateService) GetByID(tenantID, id string) (*models.Estimate, error) {
	return s.repo.FindByID(tenantID, id)
}

func (s *estimateService) GetByCar(tenantID, carID string) ([]models.Estimate, error) {
	return s.repo.FindByCar(tenantID, carID)
}

func (s *estimateService) Create(tenantID string, req dto.CreateEstimateRequest, settings *models.Settings) (*models.Estimate, error) {
	year := req.Date.Year()
	number, err := s.repo.NextNumber(tenantID, settings.EstPrefix, year)
	if err != nil {
		return nil, err
	}

	jtJSON, _ := json.Marshal(req.JobTypes)
	stJSON, _ := json.Marshal(req.Subtypes)
	status := req.Status
	if status == "" {
		status = "Draft"
	}

	subtotal, tax, grand := computeTotals(req.Items, req.DiscountAmount, req.TaxPct, req.TaxAmount)

	e := &models.Estimate{
		ClientCompanyID: tenantID,
		Number:          number,
		PublicSlug:      utils.PublicSlug(),
		Date:            req.Date,
		CustomerID:      req.CustomerID,
		CarID:           req.CarID,
		JobTypes:        string(jtJSON),
		Subtypes:        string(stJSON),
		MechanicID:      req.MechanicID,
		Status:          status,
		Notes:           req.Notes,
		DiscountAmount:  req.DiscountAmount,
		TaxPct:          req.TaxPct,
		TaxAmount:       req.TaxAmount,
		Subtotal:        subtotal,
		Tax:             tax,
		Grand:           grand,
		Items:           mapEstimateItems(req.Items),
		Images:          mapEstimateImages(req.Images),
	}
	return e, s.repo.Create(e)
}

func (s *estimateService) Update(tenantID, id string, req dto.UpdateEstimateRequest) (*models.Estimate, error) {
	e, err := s.repo.FindByID(tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	if req.Date != nil {
		e.Date = *req.Date
	}
	if req.JobTypes != nil {
		jtJSON, _ := json.Marshal(req.JobTypes)
		e.JobTypes = string(jtJSON)
	}
	if req.Subtypes != nil {
		stJSON, _ := json.Marshal(req.Subtypes)
		e.Subtypes = string(stJSON)
	}
	if req.MechanicID != nil {
		e.MechanicID = req.MechanicID
	}
	if req.Status != "" {
		e.Status = req.Status
	}
	if req.ConvertedInvoiceID != nil {
		e.ConvertedInvoiceID = req.ConvertedInvoiceID
	}
	if req.Notes != "" {
		e.Notes = req.Notes
	}
	if req.DiscountAmount != nil {
		e.DiscountAmount = *req.DiscountAmount
	}
	if req.TaxPct != nil {
		e.TaxPct = *req.TaxPct
	}
	if req.TaxAmount != nil {
		e.TaxAmount = *req.TaxAmount
	}
	if req.Items != nil {
		e.Items = mapEstimateItems(req.Items)
		for i := range e.Items {
			e.Items[i].EstimateID = e.ID
		}
	}
	if req.Images != nil {
		e.Images = mapEstimateImages(req.Images)
		for i := range e.Images {
			e.Images[i].EstimateID = e.ID
		}
	}
	e.Subtotal, e.Tax, e.Grand = recomputeFromEstimateItems(e.Items, e.DiscountAmount, e.TaxPct, e.TaxAmount)
	return e, s.repo.Update(e)
}

func (s *estimateService) Delete(tenantID, id string) error {
	e, err := s.repo.FindByID(tenantID, id)
	if err != nil {
		return errors.New("not found")
	}
	if e.Status == "Converted" {
		return errors.New("cannot delete a converted estimate")
	}
	return s.repo.Delete(tenantID, id)
}

func (s *estimateService) Convert(tenantID, id string, req dto.ConvertEstimateRequest, settings *models.Settings) (*models.Invoice, error) {
	e, err := s.repo.FindByID(tenantID, id)
	if err != nil {
		return nil, errors.New("not found")
	}
	if e.Status == "Converted" {
		return nil, errors.New("already converted")
	}

	year := time.Now().Year()
	number, err := s.invRepo.NextNumber(tenantID, settings.InvPrefix, year)
	if err != nil {
		return nil, err
	}

	estID := e.ID
	inv := &models.Invoice{
		ClientCompanyID: tenantID,
		Number:          number,
		PublicSlug:      utils.PublicSlug(),
		Date:            time.Now(),
		DueDate:         req.DueDate,
		CustomerID:      e.CustomerID,
		CarID:           e.CarID,
		JobTypes:        e.JobTypes,
		Subtypes:        e.Subtypes,
		MechanicID:      e.MechanicID,
		Notes:           e.Notes,
		DiscountAmount:  e.DiscountAmount,
		TaxPct:          e.TaxPct,
		TaxAmount:       e.TaxAmount,
		Subtotal:        e.Subtotal,
		Tax:             e.Tax,
		Grand:           e.Grand,
		Balance:         e.Grand,
		DeliveryDate:    req.DeliveryDate,
		FinalMileage:    req.FinalMileage,
		Warranty:        req.Warranty,
		FromEstimateID:  &estID,
		Items:           estimateItemsToInvoiceItems(e.Items),
		Images:          estimateImagesToInvoiceImages(e.Images),
	}
	inv.Status = computeInvoiceStatus(inv)

	if err := s.invRepo.Create(inv); err != nil {
		return nil, err
	}

	invIDCopy := inv.ID
	e.Status = "Converted"
	e.ConvertedInvoiceID = &invIDCopy
	_ = s.repo.Update(e)

	return inv, nil
}
