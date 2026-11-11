package services

import (
	"autoworkshop/dto"
	"autoworkshop/models"
	"autoworkshop/repositories"
	"autoworkshop/utils"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type VehicleService interface {
	List(tenantID string, p utils.PageParams, customerID string) ([]models.Car, int64, error)
	GetByID(tenantID, id string) (*models.Car, error)
	GetByCustomer(tenantID, customerID string) ([]models.Car, error)
	Create(tenantID string, req dto.CreateCarRequest) (*models.Car, error)
	Update(tenantID, id string, req dto.UpdateCarRequest) (*models.Car, error)
	Delete(tenantID, id string) error
}

type vehicleService struct {
	repo repositories.VehicleRepository
}

func NewVehicleService(repo repositories.VehicleRepository) VehicleService {
	return &vehicleService{repo}
}

func (s *vehicleService) List(tenantID string, p utils.PageParams, customerID string) ([]models.Car, int64, error) {
	return s.repo.List(tenantID, p, customerID)
}

func (s *vehicleService) GetByID(tenantID, id string) (*models.Car, error) {
	return s.repo.FindByID(tenantID, id)
}

func (s *vehicleService) GetByCustomer(tenantID, customerID string) ([]models.Car, error) {
	return s.repo.FindByCustomer(tenantID, customerID)
}

func (s *vehicleService) Create(tenantID string, req dto.CreateCarRequest) (*models.Car, error) {
	plate := strings.ToUpper(req.Plate)
	if req.Company != "" {
		_ = s.repo.EnsureCompany(tenantID, req.Company)
	}
	car := &models.Car{
		ClientCompanyID: tenantID,
		CustomerID:      req.CustomerID,
		Plate:           plate,
		Chassis:         req.Chassis,
		Company:         req.Company,
		Model:           req.Model,
		Year:            req.Year,
		Color:           req.Color,
		ColorLabel:      req.ColorLabel,
		Mileage:         req.Mileage,
		Notes:           req.Notes,
	}
	return car, s.repo.Create(car)
}

func (s *vehicleService) Update(tenantID, id string, req dto.UpdateCarRequest) (*models.Car, error) {
	car, err := s.repo.FindByID(tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	if req.Plate != "" {
		car.Plate = strings.ToUpper(req.Plate)
	}
	if req.Chassis != "" {
		car.Chassis = req.Chassis
	}
	if req.Company != "" {
		car.Company = req.Company
		_ = s.repo.EnsureCompany(tenantID, req.Company)
	}
	if req.Model != "" {
		car.Model = req.Model
	}
	if req.Year != 0 {
		car.Year = req.Year
	}
	if req.Color != "" {
		car.Color = req.Color
	}
	if req.ColorLabel != "" {
		car.ColorLabel = req.ColorLabel
	}
	if req.Mileage != 0 {
		car.Mileage = req.Mileage
	}
	if req.Notes != "" {
		car.Notes = req.Notes
	}
	return car, s.repo.Update(car)
}

func (s *vehicleService) Delete(tenantID, id string) error {
	has, err := s.repo.HasInvoices(tenantID, id)
	if err != nil {
		return err
	}
	if has {
		return errors.New("cannot delete: car has linked invoices")
	}
	return s.repo.Delete(tenantID, id)
}
