package services

import (
	"autoworkshop/dto"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// Expected CSV headers (same export pattern every time; extra columns ignored).
var leadCSVHeaders = map[string]func(*dto.LeadImportRow, string){
	"business name": func(r *dto.LeadImportRow, v string) { r.BusinessName = v },
	"city":          func(r *dto.LeadImportRow, v string) { r.City = v },
	"country":       func(r *dto.LeadImportRow, v string) { r.Country = v },
	"category":      func(r *dto.LeadImportRow, v string) { r.Category = v },
	"phone":         func(r *dto.LeadImportRow, v string) { r.Phone = v },
	"whatsapp":      func(r *dto.LeadImportRow, v string) { r.WhatsApp = v },
	"email":         func(r *dto.LeadImportRow, v string) { r.Email = v },
	"username":      func(r *dto.LeadImportRow, v string) { r.Username = v },
	"password":      func(r *dto.LeadImportRow, v string) { r.Password = v },
	"website":       func(r *dto.LeadImportRow, v string) { r.Website = v },
}

func normalizeCSVHeader(h string) string {
	h = strings.TrimSpace(h)
	h = strings.TrimPrefix(h, "#")
	h = strings.TrimSpace(h)
	return strings.ToLower(h)
}

func ParseLeadsCSV(r io.Reader) ([]dto.LeadImportRow, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read csv header: %w", err)
	}
	colMap := make([]func(*dto.LeadImportRow, string), len(header))
	for i, h := range header {
		key := normalizeCSVHeader(h)
		if fn, ok := leadCSVHeaders[key]; ok {
			colMap[i] = fn
		}
	}

	var rows []dto.LeadImportRow
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv row: %w", err)
		}
		var row dto.LeadImportRow
		for i, val := range record {
			if i < len(colMap) && colMap[i] != nil {
				colMap[i](&row, strings.TrimSpace(val))
			}
		}
		if strings.TrimSpace(row.BusinessName) == "" && strings.TrimSpace(row.Username) == "" && strings.TrimSpace(row.Email) == "" {
			continue
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("csv has no data rows")
	}
	return rows, nil
}
