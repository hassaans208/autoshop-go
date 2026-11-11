package services

import (
	"autoworkshop/dto"
	"autoworkshop/models"
	"autoworkshop/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StaffService interface {
	List(tenantID string) ([]models.Staff, error)
	GetByID(tenantID, id string) (*models.Staff, error)
	Create(tenantID string, req dto.CreateStaffRequest) (*models.Staff, error)
	Update(tenantID, id string, req dto.UpdateStaffRequest) (*models.Staff, error)
	Delete(tenantID, id string) error
}

type staffService struct{ repo repositories.StaffRepository }

func NewStaffService(repo repositories.StaffRepository) StaffService {
	return &staffService{repo}
}

func (s *staffService) List(tenantID string) ([]models.Staff, error) {
	return s.repo.List(tenantID)
}

func (s *staffService) GetByID(tenantID, id string) (*models.Staff, error) {
	return s.repo.FindByID(tenantID, id)
}

func (s *staffService) Create(tenantID string, req dto.CreateStaffRequest) (*models.Staff, error) {
	existing, err := s.repo.FindByEmail(req.Email)
	if err == nil && existing.ID != "" && existing.ClientCompanyID == tenantID {
		return nil, errors.New("email already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	active := true
	if req.Active != nil {
		active = *req.Active
	}
	staff := &models.Staff{
		ClientCompanyID: tenantID,
		Email:           req.Email,
		PasswordHash:    string(hash),
		Name:            req.Name,
		Role:            req.Role,
		Phone:           req.Phone,
		Salary:          req.Salary,
		Active:          active,
	}
	return staff, s.repo.Create(staff)
}

func (s *staffService) Update(tenantID, id string, req dto.UpdateStaffRequest) (*models.Staff, error) {
	staff, err := s.repo.FindByID(tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	if req.Name != "" {
		staff.Name = req.Name
	}
	if req.Email != "" {
		staff.Email = req.Email
	}
	if req.Role != "" {
		staff.Role = req.Role
	}
	if req.Phone != "" {
		staff.Phone = req.Phone
	}
	if req.Salary != nil {
		staff.Salary = *req.Salary
	}
	if req.Active != nil {
		staff.Active = *req.Active
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		staff.PasswordHash = string(hash)
	}
	return staff, s.repo.Update(staff)
}

func (s *staffService) Delete(tenantID, id string) error {
	return s.repo.Delete(tenantID, id)
}
