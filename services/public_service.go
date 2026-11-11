package services

import (
	"autoworkshop/models"
	"autoworkshop/repositories"
	"encoding/json"
)

type PublicService interface {
	GetEstimate(slug string) (map[string]any, error)
	GetInvoice(slug string) (map[string]any, error)
}

type publicService struct {
	estRepo      repositories.EstimateRepository
	invRepo      repositories.InvoiceRepository
	settingsRepo repositories.SettingsRepository
	custRepo     repositories.CustomerRepository
	carRepo      repositories.VehicleRepository
}

func NewPublicService(
	estRepo repositories.EstimateRepository,
	invRepo repositories.InvoiceRepository,
	settingsRepo repositories.SettingsRepository,
	custRepo repositories.CustomerRepository,
	carRepo repositories.VehicleRepository,
) PublicService {
	return &publicService{estRepo, invRepo, settingsRepo, custRepo, carRepo}
}

func (s *publicService) GetEstimate(slug string) (map[string]any, error) {
	e, err := s.estRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}
	settings, _ := s.settingsRepo.Get(e.ClientCompanyID)
	customer, _ := s.custRepo.FindByID(e.ClientCompanyID, e.CustomerID)
	car, _ := s.carRepo.FindByID(e.ClientCompanyID, e.CarID)

	return map[string]any{
		"id":             e.ID,
		"number":         e.Number,
		"status":         e.Status,
		"date":           e.Date,
		"publicSlug":     e.PublicSlug,
		"customerId":     e.CustomerID,
		"carId":          e.CarID,
		"jobTypes":       parseJSONStrings(e.JobTypes),
		"subtypes":       parseJSONStrings(e.Subtypes),
		"items":          sanitizeEstimateItems(e.Items),
		"images":         publicImages(e.Images),
		"discountAmount": e.DiscountAmount,
		"taxPct":         e.TaxPct,
		"taxAmount":      e.TaxAmount,
		"subtotal":       e.Subtotal,
		"grandTotal":     e.Grand,
		"workshop":       workshopPublic(settings),
		"customer":       customerPublic(customer),
		"car":            carPublic(car),
	}, nil
}

func (s *publicService) GetInvoice(slug string) (map[string]any, error) {
	inv, err := s.invRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}
	settings, _ := s.settingsRepo.Get(inv.ClientCompanyID)
	customer, _ := s.custRepo.FindByID(inv.ClientCompanyID, inv.CustomerID)
	car, _ := s.carRepo.FindByID(inv.ClientCompanyID, inv.CarID)

	return map[string]any{
		"id":             inv.ID,
		"number":         inv.Number,
		"status":         inv.Status,
		"date":           inv.Date,
		"publicSlug":     inv.PublicSlug,
		"dueDate":        inv.DueDate,
		"deliveryDate":   inv.DeliveryDate,
		"warranty":       inv.Warranty,
		"finalMileage":   inv.FinalMileage,
		"customerId":     inv.CustomerID,
		"carId":          inv.CarID,
		"jobTypes":       parseJSONStrings(inv.JobTypes),
		"subtypes":       parseJSONStrings(inv.Subtypes),
		"items":          sanitizeInvoiceItems(inv.Items),
		"images":         publicInvoiceImages(inv.Images),
		"discountAmount": inv.DiscountAmount,
		"taxPct":         inv.TaxPct,
		"taxAmount":      inv.TaxAmount,
		"subtotal":       inv.Subtotal,
		"grandTotal":     inv.Grand,
		"paidAmount":     inv.Paid,
		"balanceDue":     inv.Balance,
		"workshop":       workshopPublic(settings),
		"customer":       customerPublic(customer),
		"car":            carPublic(car),
	}, nil
}

func parseJSONStrings(raw string) []string {
	if raw == "" {
		return []string{}
	}
	var out []string
	_ = json.Unmarshal([]byte(raw), &out)
	if out == nil {
		return []string{}
	}
	return out
}

func sanitizeEstimateItems(items []models.EstimateItem) []map[string]any {
	out := make([]map[string]any, len(items))
	for i, it := range items {
		row := map[string]any{
			"id": it.ID, "name": it.Name, "qty": it.Qty,
			"unitPrice": it.UnitPrice, "discountAmount": it.DiscountAmount,
		}
		if it.Kind != "" {
			row["kind"] = it.Kind
		}
		if it.Subtype != "" {
			row["subtype"] = it.Subtype
		}
		if it.ServicePrice > 0 {
			row["servicePrice"] = it.ServicePrice
		}
		out[i] = row
	}
	return out
}

func sanitizeInvoiceItems(items []models.InvoiceItem) []map[string]any {
	out := make([]map[string]any, len(items))
	for i, it := range items {
		row := map[string]any{
			"id": it.ID, "name": it.Name, "qty": it.Qty,
			"unitPrice": it.UnitPrice, "discountAmount": it.DiscountAmount,
		}
		if it.Kind != "" {
			row["kind"] = it.Kind
		}
		if it.Subtype != "" {
			row["subtype"] = it.Subtype
		}
		if it.ServicePrice > 0 {
			row["servicePrice"] = it.ServicePrice
		}
		out[i] = row
	}
	return out
}

func publicImages(images []models.EstimateImage) []map[string]any {
	out := make([]map[string]any, len(images))
	for i, img := range images {
		out[i] = map[string]any{"id": img.ID, "label": img.Label, "dataUrl": img.DataURL}
	}
	return out
}

func publicInvoiceImages(images []models.InvoiceImage) []map[string]any {
	out := make([]map[string]any, len(images))
	for i, img := range images {
		out[i] = map[string]any{"id": img.ID, "label": img.Label, "dataUrl": img.DataURL}
	}
	return out
}

func workshopPublic(s *models.Settings) map[string]any {
	if s == nil {
		return map[string]any{}
	}
	return map[string]any{
		"name": s.WorkshopName, "logo": s.Logo, "phone": s.Phone, "email": s.Email,
		"address": s.Address, "ntn": s.NTN,
		"terms": s.Terms, "warrantyDefault": s.WarrantyDefault, "footer": s.Footer,
	}
}

func customerPublic(c *models.Customer) map[string]any {
	if c == nil {
		return map[string]any{}
	}
	return map[string]any{"name": c.Name, "phone": c.Phone}
}

func carPublic(c *models.Car) map[string]any {
	if c == nil {
		return map[string]any{}
	}
	return map[string]any{
		"plate": c.Plate, "company": c.Company, "model": c.Model, "year": c.Year,
	}
}
