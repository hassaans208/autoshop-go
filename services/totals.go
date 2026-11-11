package services

import (
	"autoworkshop/dto"
	"autoworkshop/models"
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
)

func lineTotal(it dto.LineItemDTO) float64 {
	base := it.Qty * it.UnitPrice * (1 - it.DiscountPct/100)
	return math.Max(0, base-it.DiscountAmount)
}

func computeTotals(items []dto.LineItemDTO, discountAmount, taxPct, taxAmount float64) (subtotal, tax, grand float64) {
	for _, it := range items {
		subtotal += lineTotal(it)
	}
	afterDisc := math.Max(0, subtotal-discountAmount)
	if taxAmount > 0 {
		tax = taxAmount
	} else {
		tax = afterDisc * taxPct / 100
	}
	grand = afterDisc + tax
	return
}

func recomputeFromEstimateItems(items []models.EstimateItem, discountAmount, taxPct, taxAmount float64) (subtotal, tax, grand float64) {
	for _, it := range items {
		base := it.Qty * it.UnitPrice * (1 - it.DiscountPct/100)
		subtotal += math.Max(0, base-it.DiscountAmount)
	}
	afterDisc := math.Max(0, subtotal-discountAmount)
	if taxAmount > 0 {
		tax = taxAmount
	} else {
		tax = afterDisc * taxPct / 100
	}
	grand = afterDisc + tax
	return
}

func recomputeFromInvoiceItems(items []models.InvoiceItem, discountAmount, taxPct, taxAmount float64) (subtotal, tax, grand float64) {
	for _, it := range items {
		base := it.Qty * it.UnitPrice * (1 - it.DiscountPct/100)
		subtotal += math.Max(0, base-it.DiscountAmount)
	}
	afterDisc := math.Max(0, subtotal-discountAmount)
	if taxAmount > 0 {
		tax = taxAmount
	} else {
		tax = afterDisc * taxPct / 100
	}
	grand = afterDisc + tax
	return
}

func mapEstimateItems(items []dto.LineItemDTO) []models.EstimateItem {
	out := make([]models.EstimateItem, len(items))
	for i, it := range items {
		id := it.ID
		if id == "" {
			id = uuid.NewString()
		}
		assignees, _ := json.Marshal(it.Assignees.Slice())
		out[i] = models.EstimateItem{
			ID: id, Name: it.Name, Description: it.Description,
			Qty: it.Qty, UnitPrice: it.UnitPrice, DiscountPct: it.DiscountPct,
			DiscountAmount: it.DiscountAmount, Kind: it.Kind,
			PurchasePrice: it.PurchasePrice, VendorID: it.VendorID,
			JobType: it.JobType, Subtype: it.Subtype,
			Assignees: string(assignees), MaterialPrice: it.MaterialPrice,
			ServicePrice: it.ServicePrice, SortOrder: it.SortOrder,
		}
	}
	return out
}

func mapInvoiceItems(items []dto.LineItemDTO, invoiceID string) []models.InvoiceItem {
	out := make([]models.InvoiceItem, len(items))
	for i, it := range items {
		id := it.ID
		if id == "" {
			id = uuid.NewString()
		}
		assignees, _ := json.Marshal(it.Assignees.Slice())
		out[i] = models.InvoiceItem{
			ID: id, InvoiceID: invoiceID, Name: it.Name, Description: it.Description,
			Qty: it.Qty, UnitPrice: it.UnitPrice, DiscountPct: it.DiscountPct,
			DiscountAmount: it.DiscountAmount, Kind: it.Kind,
			PurchasePrice: it.PurchasePrice, VendorID: it.VendorID,
			JobType: it.JobType, Subtype: it.Subtype,
			Assignees: string(assignees), MaterialPrice: it.MaterialPrice,
			ServicePrice: it.ServicePrice, SortOrder: it.SortOrder,
		}
	}
	return out
}

func estimateItemsToInvoiceItems(items []models.EstimateItem) []models.InvoiceItem {
	out := make([]models.InvoiceItem, len(items))
	for i, it := range items {
		out[i] = models.InvoiceItem{
			ID: uuid.NewString(), Name: it.Name, Description: it.Description,
			Qty: it.Qty, UnitPrice: it.UnitPrice, DiscountPct: it.DiscountPct,
			DiscountAmount: it.DiscountAmount, Kind: it.Kind,
			PurchasePrice: it.PurchasePrice, VendorID: it.VendorID,
			JobType: it.JobType, Subtype: it.Subtype,
			Assignees: it.Assignees, MaterialPrice: it.MaterialPrice,
			ServicePrice: it.ServicePrice, SortOrder: it.SortOrder,
		}
	}
	return out
}

func mapEstimateImages(images []dto.ImageItemDTO) []models.EstimateImage {
	out := make([]models.EstimateImage, len(images))
	for i, img := range images {
		id := img.ID
		if id == "" {
			id = uuid.NewString()
		}
		out[i] = models.EstimateImage{ID: id, Label: img.Label, DataURL: img.DataURL, SortOrder: img.SortOrder}
	}
	return out
}

func mapInvoiceImages(images []dto.ImageItemDTO) []models.InvoiceImage {
	out := make([]models.InvoiceImage, len(images))
	for i, img := range images {
		id := img.ID
		if id == "" {
			id = uuid.NewString()
		}
		out[i] = models.InvoiceImage{ID: id, Label: img.Label, DataURL: img.DataURL, SortOrder: img.SortOrder}
	}
	return out
}

func estimateImagesToInvoiceImages(images []models.EstimateImage) []models.InvoiceImage {
	out := make([]models.InvoiceImage, len(images))
	for i, img := range images {
		out[i] = models.InvoiceImage{
			ID: uuid.NewString(), Label: img.Label, DataURL: img.DataURL, SortOrder: img.SortOrder,
		}
	}
	return out
}

func sumPayments(payments []models.Payment) float64 {
	var total float64
	for _, p := range payments {
		total += p.Amount
	}
	return total
}

// invoiceIsPastDue compares calendar dates in now's location: due on that local day is not past due until the next local day.
func invoiceIsPastDue(inv *models.Invoice, now time.Time) bool {
	if inv.DueDate == nil {
		return false
	}
	loc := now.Location()
	due := inv.DueDate.In(loc)
	ny, nm, nd := now.Date()
	dy, dm, dd := due.Date()
	todayStart := time.Date(ny, nm, nd, 0, 0, 0, 0, loc)
	dueStart := time.Date(dy, dm, dd, 0, 0, 0, 0, loc)
	return todayStart.After(dueStart)
}

func computeInvoiceStatus(inv *models.Invoice) string {
	now := time.Now()
	if inv.Paid >= inv.Grand {
		return "Paid"
	}
	balance := math.Max(0, inv.Grand-inv.Paid)
	if balance <= 0 {
		return "Paid"
	}
	if invoiceIsPastDue(inv, now) {
		return "Overdue"
	}
	if inv.Paid > 0 {
		return "Partially Paid"
	}
	return "Unpaid"
}
