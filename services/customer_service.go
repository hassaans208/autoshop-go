package services

import (
	"autoworkshop/dto"
	"autoworkshop/models"
	"autoworkshop/repositories"
	"autoworkshop/utils"
	"errors"

	"gorm.io/gorm"
)

type CustomerService interface {
	List(tenantID string, p utils.PageParams) ([]models.Customer, int64, error)
	GetByID(tenantID, id string) (*models.Customer, error)
	GetByPhone(tenantID, phone string) (*models.Customer, error)
	Create(tenantID string, req dto.CreateCustomerRequest) (*models.Customer, error)
	Update(tenantID, id string, req dto.UpdateCustomerRequest) (*models.Customer, error)
	Delete(tenantID, id string) error
}

type customerService struct {
	repo repositories.CustomerRepository
}

func NewCustomerService(repo repositories.CustomerRepository) CustomerService {
	return &customerService{repo}
}

func (s *customerService) List(tenantID string, p utils.PageParams) ([]models.Customer, int64, error) {
	return s.repo.List(tenantID, p)
}

func (s *customerService) GetByID(tenantID, id string) (*models.Customer, error) {
	return s.repo.FindByID(tenantID, id)
}

func (s *customerService) GetByPhone(tenantID, phone string) (*models.Customer, error) {
	norm := utils.NormalizePhone(phone)
	return s.repo.FindByPhone(tenantID, norm)
}

func (s *customerService) Create(tenantID string, req dto.CreateCustomerRequest) (*models.Customer, error) {
	norm := utils.NormalizePhone(req.Phone)
	existing, err := s.repo.FindByPhone(tenantID, norm)
	if err == nil && existing.ID != "" {
		return nil, errors.New("phone already exists")
	}
	c := &models.Customer{
		ClientCompanyID: tenantID,
		Name:            req.Name,
		Phone:           req.Phone,
		PhoneNorm:       norm,
		Email:           req.Email,
		Address:         req.Address,
		Notes:           req.Notes,
	}
	return c, s.repo.Create(c)
}

func (s *customerService) Update(tenantID, id string, req dto.UpdateCustomerRequest) (*models.Customer, error) {
	c, err := s.repo.FindByID(tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	if req.Name != "" {
		c.Name = req.Name
	}
	if req.Phone != "" {
		c.Phone = req.Phone
		c.PhoneNorm = utils.NormalizePhone(req.Phone)
	}
	if req.Email != "" {
		c.Email = req.Email
	}
	if req.Address != "" {
		c.Address = req.Address
	}
	if req.Notes != "" {
		c.Notes = req.Notes
	}
	return c, s.repo.Update(c)
}

func (s *customerService) Delete(tenantID, id string) error {
	has, err := s.repo.HasCarsOrInvoices(tenantID, id)
	if err != nil {
		return err
	}
	if has {
		return errors.New("cannot delete: customer has linked cars or invoices")
	}
	return s.repo.Delete(tenantID, id)
}
