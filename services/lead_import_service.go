package services

import (
	"autoworkshop/dto"
	"autoworkshop/models"
	"autoworkshop/repositories"
	"autoworkshop/utils"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LeadImportService interface {
	Import(rows []dto.LeadImportRow) (*dto.LeadImportResponse, error)
}

type leadImportService struct {
	db           *gorm.DB
	staffRepo    repositories.StaffRepository
	frontendURL  string
}

func NewLeadImportService(db *gorm.DB, staffRepo repositories.StaffRepository, frontendURL string) LeadImportService {
	return &leadImportService{db: db, staffRepo: staffRepo, frontendURL: strings.TrimRight(strings.TrimSpace(frontendURL), "/")}
}

func (s *leadImportService) Import(rows []dto.LeadImportRow) (*dto.LeadImportResponse, error) {
	if s.frontendURL == "" {
		return nil, errors.New("FRONTEND_URL is not configured")
	}
	out := &dto.LeadImportResponse{Results: make([]dto.LeadImportResultRow, 0, len(rows))}
	for _, row := range rows {
		res := s.importOne(row)
		out.Results = append(out.Results, res)
		switch res.Status {
		case "created":
			out.Created++
		case "updated":
			out.Updated++
		default:
			out.Failed++
		}
	}
	return out, nil
}

func (s *leadImportService) importOne(row dto.LeadImportRow) dto.LeadImportResultRow {
	businessName := utils.TruncateBusinessName(row.BusinessName)
	loginSource := strings.TrimSpace(row.Username)
	if loginSource == "" {
		loginSource = strings.TrimSpace(row.Email)
	}
	loginEmail := utils.NormalizeLoginEmail(loginSource)
	result := dto.LeadImportResultRow{
		BusinessName: businessName,
		LoginEmail:   loginEmail,
	}

	if businessName == "" {
		result.Status = "error"
		result.Error = "businessName is required"
		return result
	}
	if loginEmail == "" {
		result.Status = "error"
		result.Error = "username or email required for login"
		return result
	}
	if len(strings.TrimSpace(row.Password)) < 6 {
		result.Status = "error"
		result.Error = "password must be at least 6 characters"
		return result
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(row.Password), bcrypt.DefaultCost)
	if err != nil {
		result.Status = "error"
		result.Error = "could not hash password"
		return result
	}
	hashStr := string(hash)

	existing, err := s.staffRepo.FindByEmail(loginEmail)
	if err == nil && existing.ID != "" {
		return s.updateExisting(existing, row, businessName, loginEmail, hashStr, result)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		result.Status = "error"
		result.Error = fmt.Sprintf("lookup staff: %v", err)
		return result
	}

	ownerID := uuid.NewString()
	companyID := uuid.NewString()
	company := &models.ClientCompany{
		ID: companyID, Name: businessName, OwnerID: ownerID, Active: true,
	}
	owner := &models.Staff{
		ID:              ownerID,
		ClientCompanyID: companyID,
		Email:           loginEmail,
		PasswordHash:    hashStr,
		Name:            businessName,
		Role:            "owner",
		Phone:           strings.TrimSpace(row.Phone),
		Active:          true,
	}
	settings := settingsFromRow(row, companyID, businessName, loginEmail)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(company).Error; err != nil {
			return err
		}
		if err := tx.Create(owner).Error; err != nil {
			return err
		}
		return tx.Create(settings).Error
	})
	if err != nil {
		result.Status = "error"
		result.Error = fmt.Sprintf("database error: %v", err)
		return result
	}

	result.Status = "created"
	result.CompanyID = companyID
	result.StaffID = ownerID
	result.AutologinURL = utils.BuildAutologinURL(s.frontendURL, loginEmail, hashStr)
	return result
}

func (s *leadImportService) updateExisting(staff *models.Staff, row dto.LeadImportRow, businessName, loginEmail, hashStr string, result dto.LeadImportResultRow) dto.LeadImportResultRow {
	tenantID := staff.ClientCompanyID
	if tenantID == "" {
		result.Status = "error"
		result.Error = "existing staff has no tenant"
		return result
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		staff.PasswordHash = hashStr
		staff.Name = businessName
		staff.Phone = strings.TrimSpace(row.Phone)
		staff.Active = true
		if err := tx.Save(staff).Error; err != nil {
			return err
		}

		var company models.ClientCompany
		if err := tx.First(&company, "id = ?", tenantID).Error; err != nil {
			return err
		}
		company.Name = businessName
		company.Active = true
		if err := tx.Save(&company).Error; err != nil {
			return err
		}

		settings := settingsFromRow(row, tenantID, businessName, loginEmail)
		var existing models.Settings
		err := tx.Where("client_company_id = ?", tenantID).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Create(settings).Error
		}
		if err != nil {
			return err
		}
		settings.ID = existing.ID
		settings.ClientCompanyID = tenantID
		settings.Currency = coalesce(settings.Currency, existing.Currency, "PKR")
		settings.EstPrefix = coalesce(settings.EstPrefix, existing.EstPrefix, "EST")
		settings.InvPrefix = coalesce(settings.InvPrefix, existing.InvPrefix, "INV")
		if existing.EstStart > 0 {
			settings.EstStart = existing.EstStart
		}
		if existing.InvStart > 0 {
			settings.InvStart = existing.InvStart
		}
		return tx.Save(settings).Error
	})
	if err != nil {
		result.Status = "error"
		result.Error = fmt.Sprintf("update failed: %v", err)
		return result
	}

	result.Status = "updated"
	result.CompanyID = tenantID
	result.StaffID = staff.ID
	result.AutologinURL = utils.BuildAutologinURL(s.frontendURL, loginEmail, hashStr)
	return result
}

func settingsFromRow(row dto.LeadImportRow, tenantID, businessName, loginEmail string) *models.Settings {
	contactEmail := strings.TrimSpace(row.Email)
	if contactEmail == "" {
		contactEmail = loginEmail
	}
	tagline := strings.TrimSpace(row.Category)
	if tagline == "" {
		tagline = "Quality Auto Repair & Service"
	}
	footer := strings.TrimSpace(row.Website)
	return &models.Settings{
		ClientCompanyID: tenantID,
		WorkshopName:    businessName,
		Tagline:         tagline,
		Address:         joinNonEmpty(row.City, row.Country),
		Phone:           strings.TrimSpace(row.Phone),
		Email:           contactEmail,
		Whatsapp:        strings.TrimSpace(row.WhatsApp),
		Currency:        "PKR",
		EstPrefix:       "EST",
		InvPrefix:       "INV",
		EstStart:        1,
		InvStart:        1,
		Footer:          footer,
	}
}

func coalesce(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func joinNonEmpty(parts ...string) string {
	var kept []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, ", ")
}
