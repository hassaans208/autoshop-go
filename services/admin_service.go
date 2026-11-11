package services

import (
	"autoworkshop/dto"
	"autoworkshop/models"
	"autoworkshop/repositories"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminService interface {
	CreateCompany(req dto.AdminCreateCompanyRequest) (*dto.AdminCreateCompanyResponse, error)
	ListCompanies() ([]dto.AdminCompanyResponse, error)
}

type adminService struct {
	db          *gorm.DB
	companyRepo repositories.ClientCompanyRepository
	staffRepo   repositories.StaffRepository
}

func NewAdminService(db *gorm.DB, companyRepo repositories.ClientCompanyRepository, staffRepo repositories.StaffRepository) AdminService {
	return &adminService{db, companyRepo, staffRepo}
}

func (s *adminService) CreateCompany(req dto.AdminCreateCompanyRequest) (*dto.AdminCreateCompanyResponse, error) {
	existing, err := s.staffRepo.FindByEmail(req.Owner.Email)
	if err == nil && existing.ID != "" {
		return nil, errors.New("owner email already in use")
	}

	workshopName := req.WorkshopName
	if workshopName == "" {
		workshopName = req.CompanyName
	}

	ownerID := uuid.NewString()
	companyID := uuid.NewString()

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Owner.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	company := &models.ClientCompany{
		ID:      companyID,
		Name:    req.CompanyName,
		OwnerID: ownerID,
		Active:  true,
	}
	owner := &models.Staff{
		ID:              ownerID,
		ClientCompanyID: companyID,
		Name:            req.Owner.Name,
		Email:           req.Owner.Email,
		PasswordHash:    string(hash),
		Role:            "owner",
		Phone:           req.Owner.Phone,
		Salary:          req.Owner.Salary,
		Active:          true,
	}
	settings := &models.Settings{
		ClientCompanyID: companyID,
		WorkshopName:    workshopName,
		Currency:        "PKR",
		EstPrefix:       "EST",
		InvPrefix:       "INV",
	}

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
		return nil, err
	}

	return &dto.AdminCreateCompanyResponse{
		Company: dto.AdminCompanyResponse{
			ID: company.ID, Name: company.Name, OwnerID: company.OwnerID, Active: company.Active,
		},
		Owner: dto.UserResponse{
			ID: owner.ID, Email: owner.Email, Name: owner.Name,
			Role: owner.Role, Phone: owner.Phone, Active: owner.Active,
		},
	}, nil
}

func (s *adminService) ListCompanies() ([]dto.AdminCompanyResponse, error) {
	list, err := s.companyRepo.List()
	if err != nil {
		return nil, err
	}
	out := make([]dto.AdminCompanyResponse, len(list))
	for i, c := range list {
		out[i] = dto.AdminCompanyResponse{
			ID: c.ID, Name: c.Name, OwnerID: c.OwnerID, Active: c.Active,
		}
	}
	return out, nil
}
