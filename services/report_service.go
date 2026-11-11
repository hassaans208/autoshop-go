package services

import (
	"autoworkshop/repositories"
	"autoworkshop/utils"
	"math"
	"sort"
	"time"
)

type DashboardService interface {
	Financials(tenantID, month string) (map[string]any, error)
}

type dashboardService struct {
	staffRepo repositories.StaffRepository
	custRepo  repositories.CustomerRepository
	carRepo   repositories.VehicleRepository
	estRepo   repositories.EstimateRepository
	invRepo   repositories.InvoiceRepository
}

func NewDashboardService(
	staffRepo repositories.StaffRepository,
	custRepo repositories.CustomerRepository,
	carRepo repositories.VehicleRepository,
	estRepo repositories.EstimateRepository,
	invRepo repositories.InvoiceRepository,
) DashboardService {
	return &dashboardService{staffRepo, custRepo, carRepo, estRepo, invRepo}
}

func (s *dashboardService) Financials(tenantID, month string) (map[string]any, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	start, _ := time.Parse("2006-01", month)
	end := start.AddDate(0, 1, 0)

	salaryObligation, _ := s.staffRepo.SumActiveSalaries(tenantID)
	invoices, _, _ := s.invRepo.List(tenantID, bigPage(), map[string]string{
		"from": start.Format(time.RFC3339),
		"to":   end.Format(time.RFC3339),
	})

	var revenue, partsPurchaseCost, materialCost float64
	for _, inv := range invoices {
		for _, p := range inv.Payments {
			if !p.Date.Before(start) && p.Date.Before(end) {
				revenue += p.Amount
			}
		}
		for _, it := range inv.Items {
			if !inv.Date.Before(start) && inv.Date.Before(end) {
				partsPurchaseCost += it.PurchasePrice * it.Qty
				materialCost += it.MaterialPrice
			}
		}
	}

	remaining := math.Max(0, salaryObligation-revenue)
	profit := math.Max(0, revenue-salaryObligation-partsPurchaseCost-materialCost)
	customers, _ := s.custRepo.Count(tenantID)
	cars, _ := s.carRepo.Count(tenantID)
	openEst, _ := s.estRepo.CountByStatus(tenantID, "Draft")
	unpaidInv, _ := s.invRepo.CountByStatus(tenantID, "Unpaid", "Partially Paid", "Overdue")
	overdueInv, _ := s.invRepo.CountByStatus(tenantID, "Overdue")

	return map[string]any{
		"month": month, "salaryObligation": salaryObligation, "revenue": revenue,
		"remainingToCoverSalaries": remaining, "partsPurchaseCost": partsPurchaseCost,
		"materialCost": materialCost, "profit": profit,
		"counts": map[string]int64{
			"customers": customers, "cars": cars, "openEstimates": openEst,
			"unpaidInvoices": unpaidInv, "overdueInvoices": overdueInv,
		},
	}, nil
}

type ReportService interface {
	Dashboard(tenantID string) (map[string]any, error)
	Revenue(tenantID, from, to, groupBy string) (map[string]any, error)
	TopCustomers(tenantID string, limit int, from, to string) ([]map[string]any, error)
	MechanicPerformance(tenantID, from, to string) ([]map[string]any, error)
}

type reportService struct {
	invRepo   repositories.InvoiceRepository
	estRepo   repositories.EstimateRepository
	custRepo  repositories.CustomerRepository
	staffRepo repositories.StaffRepository
}

func NewReportService(
	invRepo repositories.InvoiceRepository,
	estRepo repositories.EstimateRepository,
	custRepo repositories.CustomerRepository,
	staffRepo repositories.StaffRepository,
) ReportService {
	return &reportService{invRepo, estRepo, custRepo, staffRepo}
}

func (s *reportService) Dashboard(tenantID string) (map[string]any, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	invoices, _, _ := s.invRepo.List(tenantID, bigPage(), map[string]string{})
	var todayRevenue, monthRevenue float64
	var openInvoices, overdueInvoices int64
	for _, inv := range invoices {
		if inv.Status == "Unpaid" || inv.Status == "Partially Paid" || inv.Status == "Overdue" {
			openInvoices++
		}
		if inv.Status == "Overdue" {
			overdueInvoices++
		}
		for _, p := range inv.Payments {
			if !p.Date.Before(todayStart) {
				todayRevenue += p.Amount
			}
			if !p.Date.Before(monthStart) {
				monthRevenue += p.Amount
			}
		}
	}

	estimatesPending, _ := s.estRepo.CountByStatus(tenantID, "Draft")
	recentInv, _ := s.invRepo.Recent(tenantID, 5)
	recentEst, _, _ := s.estRepo.List(tenantID, pageN(5), map[string]string{"status": "Draft"})

	return map[string]any{
		"todayRevenue": todayRevenue, "monthRevenue": monthRevenue,
		"openInvoices": openInvoices, "overdueInvoices": overdueInvoices,
		"estimatesPending": estimatesPending, "activeJobs": 0,
		"recentInvoices": recentInv, "recentEstimates": recentEst,
	}, nil
}

func (s *reportService) Revenue(tenantID, from, to, groupBy string) (map[string]any, error) {
	filters := map[string]string{}
	if from != "" {
		filters["from"] = from
	}
	if to != "" {
		filters["to"] = to
	}
	invoices, _, _ := s.invRepo.List(tenantID, bigPage(), filters)

	type bucket struct {
		revenue, paid float64
		invoices      int
	}
	buckets := map[string]*bucket{}
	var totalRev, totalPaid float64

	for _, inv := range invoices {
		key := inv.Date.Format("2006-01-02")
		if groupBy == "month" {
			key = inv.Date.Format("2006-01")
		}
		if buckets[key] == nil {
			buckets[key] = &bucket{}
		}
		buckets[key].revenue += inv.Grand
		buckets[key].invoices++
		totalRev += inv.Grand
		for _, p := range inv.Payments {
			buckets[key].paid += p.Amount
			totalPaid += p.Amount
		}
	}

	series := make([]map[string]any, 0, len(buckets))
	for k, b := range buckets {
		series = append(series, map[string]any{
			"bucket": k, "revenue": b.revenue, "paid": b.paid, "invoices": b.invoices,
		})
	}
	return map[string]any{
		"series": series,
		"totals": map[string]any{
			"revenue": totalRev, "paid": totalPaid,
			"outstanding": totalRev - totalPaid, "invoices": len(invoices),
		},
	}, nil
}

func (s *reportService) TopCustomers(tenantID string, limit int, from, to string) ([]map[string]any, error) {
	filters := map[string]string{}
	if from != "" {
		filters["from"] = from
	}
	if to != "" {
		filters["to"] = to
	}
	invoices, _, _ := s.invRepo.List(tenantID, bigPage(), filters)
	type agg struct {
		name     string
		invoices int
		revenue  float64
	}
	byCust := map[string]*agg{}
	for _, inv := range invoices {
		if byCust[inv.CustomerID] == nil {
			c, _ := s.custRepo.FindByID(tenantID, inv.CustomerID)
			name := inv.CustomerID
			if c != nil {
				name = c.Name
			}
			byCust[inv.CustomerID] = &agg{name: name}
		}
		byCust[inv.CustomerID].invoices++
		byCust[inv.CustomerID].revenue += inv.Grand
	}
	type row struct {
		id, name string
		inv      int
		rev      float64
	}
	rows := make([]row, 0, len(byCust))
	for id, a := range byCust {
		rows = append(rows, row{id: id, name: a.name, inv: a.invoices, rev: a.revenue})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].rev == rows[j].rev {
			return rows[i].name < rows[j].name
		}
		return rows[i].rev > rows[j].rev
	})
	if limit > 0 && len(rows) > limit {
		rows = rows[:limit]
	}
	out := make([]map[string]any, len(rows))
	for i, r := range rows {
		out[i] = map[string]any{
			"customerId": r.id, "name": r.name, "invoices": r.inv, "revenue": r.rev,
		}
	}
	return out, nil
}

func (s *reportService) MechanicPerformance(tenantID, from, to string) ([]map[string]any, error) {
	filters := map[string]string{}
	if from != "" {
		filters["from"] = from
	}
	if to != "" {
		filters["to"] = to
	}
	invoices, _, _ := s.invRepo.List(tenantID, bigPage(), filters)
	type agg struct {
		name    string
		jobs    int
		revenue float64
	}
	byMech := map[string]*agg{}
	for _, inv := range invoices {
		if inv.MechanicID == nil {
			continue
		}
		mid := *inv.MechanicID
		if byMech[mid] == nil {
			st, _ := s.staffRepo.FindByID(tenantID, mid)
			name := mid
			if st != nil {
				name = st.Name
			}
			byMech[mid] = &agg{name: name}
		}
		byMech[mid].jobs++
		byMech[mid].revenue += inv.Grand
	}
	out := make([]map[string]any, 0, len(byMech))
	for id, a := range byMech {
		out = append(out, map[string]any{
			"mechanicId": id, "name": a.name, "jobs": a.jobs, "revenue": a.revenue,
		})
	}
	return out, nil
}

func bigPage() utils.PageParams    { return utils.PageParams{Page: 1, Limit: 10000} }
func pageN(n int) utils.PageParams { return utils.PageParams{Page: 1, Limit: n} }
