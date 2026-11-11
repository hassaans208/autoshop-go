package services

import (
	"autoworkshop/dto"
	"autoworkshop/models"
	"autoworkshop/repositories"
	"autoworkshop/utils"
	"errors"

	"gorm.io/gorm"
)

type VendorService interface {
	List(tenantID string, p utils.PageParams, kind string) ([]models.Vendor, int64, error)
	GetByID(tenantID, id string) (*models.Vendor, error)
	Create(tenantID string, req dto.CreateVendorRequest) (*models.Vendor, error)
	Update(tenantID, id string, req dto.UpdateVendorRequest) (*models.Vendor, error)
	Delete(tenantID, id string) error
}

type vendorService struct {
	repo repositories.VendorRepository
}

func NewVendorService(repo repositories.VendorRepository) VendorService {
	return &vendorService{repo}
}

func (s *vendorService) List(tenantID string, p utils.PageParams, kind string) ([]models.Vendor, int64, error) {
	return s.repo.List(tenantID, p, kind)
}

func (s *vendorService) GetByID(tenantID, id string) (*models.Vendor, error) {
	return s.repo.FindByID(tenantID, id)
}

func (s *vendorService) Create(tenantID string, req dto.CreateVendorRequest) (*models.Vendor, error) {
	v := &models.Vendor{
		ClientCompanyID: tenantID,
		Name:            req.Name,
		Kind:            req.Kind,
		Phone:           req.Phone,
		Email:           req.Email,
		Notes:           req.Notes,
	}
	return v, s.repo.Create(v)
}

func (s *vendorService) Update(tenantID, id string, req dto.UpdateVendorRequest) (*models.Vendor, error) {
	v, err := s.repo.FindByID(tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	if req.Name != "" {
		v.Name = req.Name
	}
	if req.Kind != "" {
		v.Kind = req.Kind
	}
	if req.Phone != "" {
		v.Phone = req.Phone
	}
	if req.Email != "" {
		v.Email = req.Email
	}
	if req.Notes != "" {
		v.Notes = req.Notes
	}
	return v, s.repo.Update(v)
}

func (s *vendorService) Delete(tenantID, id string) error {
	has, err := s.repo.HasParts(tenantID, id)
	if err != nil {
		return err
	}
	if has {
		return errors.New("cannot delete: vendor is referenced by parts")
	}
	return s.repo.Delete(tenantID, id)
}

type PartService interface {
	List(tenantID string, p utils.PageParams, vendorID string) ([]models.Part, int64, error)
	GetByID(tenantID, id string) (*models.Part, error)
	Create(tenantID string, req dto.CreatePartRequest) (*models.Part, error)
	Update(tenantID, id string, req dto.UpdatePartRequest) (*models.Part, error)
	Delete(tenantID, id string) error
}

type partService struct {
	repo repositories.PartRepository
}

func NewPartService(repo repositories.PartRepository) PartService {
	return &partService{repo}
}

func (s *partService) List(tenantID string, p utils.PageParams, vendorID string) ([]models.Part, int64, error) {
	return s.repo.List(tenantID, p, vendorID)
}

func (s *partService) GetByID(tenantID, id string) (*models.Part, error) {
	return s.repo.FindByID(tenantID, id)
}

func (s *partService) Create(tenantID string, req dto.CreatePartRequest) (*models.Part, error) {
	p := &models.Part{
		ClientCompanyID:      tenantID,
		Name:                 req.Name,
		VendorID:             req.VendorID,
		DefaultPurchasePrice: req.DefaultPurchasePrice,
		DefaultSalePrice:     req.DefaultSalePrice,
		Notes:                req.Notes,
	}
	return p, s.repo.Create(p)
}

func (s *partService) Update(tenantID, id string, req dto.UpdatePartRequest) (*models.Part, error) {
	p, err := s.repo.FindByID(tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.VendorID != nil {
		p.VendorID = req.VendorID
	}
	if req.DefaultPurchasePrice != nil {
		p.DefaultPurchasePrice = *req.DefaultPurchasePrice
	}
	if req.DefaultSalePrice != nil {
		p.DefaultSalePrice = *req.DefaultSalePrice
	}
	if req.Notes != "" {
		p.Notes = req.Notes
	}
	return p, s.repo.Update(p)
}

func (s *partService) Delete(tenantID, id string) error {
	return s.repo.Delete(tenantID, id)
}
