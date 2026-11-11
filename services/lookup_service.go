package services

import (
	"autoworkshop/dto"
	"autoworkshop/models"
	"autoworkshop/repositories"
	"errors"

	"gorm.io/gorm"
)

type CompanyService interface {
	List(tenantID string) ([]string, error)
	Add(tenantID string, req dto.AddCompanyRequest) (*models.Company, error)
}

type companyService struct {
	repo repositories.CompanyRepository
}

func NewCompanyService(repo repositories.CompanyRepository) CompanyService {
	return &companyService{repo}
}

func (s *companyService) List(tenantID string) ([]string, error) {
	list, err := s.repo.List(tenantID)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(list))
	for i, c := range list {
		names[i] = c.Name
	}
	return names, nil
}

func (s *companyService) Add(tenantID string, req dto.AddCompanyRequest) (*models.Company, error) {
	return s.repo.AddIfNotExists(tenantID, req.Name)
}

type JobTypeService interface {
	List(tenantID string) ([]dto.JobTypeEntry, error)
	Replace(tenantID string, req dto.ReplaceJobTypesRequest) ([]dto.JobTypeEntry, error)
}

type jobTypeService struct {
	repo repositories.JobTypeRepository
}

func NewJobTypeService(repo repositories.JobTypeRepository) JobTypeService {
	return &jobTypeService{repo}
}

func (s *jobTypeService) List(tenantID string) ([]dto.JobTypeEntry, error) {
	list, err := s.repo.List(tenantID)
	if err != nil {
		return nil, err
	}
	return toJobTypeEntries(list), nil
}

func (s *jobTypeService) Replace(tenantID string, req dto.ReplaceJobTypesRequest) ([]dto.JobTypeEntry, error) {
	items := make([]models.JobType, len(req.Entries()))
	for i, e := range req.Entries() {
		items[i] = models.JobType{
			Name:     e.Name,
			Subtypes: repositories.SubtypesToJSON(e.Subtypes),
		}
	}
	list, err := s.repo.ReplaceAll(tenantID, items)
	if err != nil {
		return nil, err
	}
	return toJobTypeEntries(list), nil
}

func toJobTypeEntries(list []models.JobType) []dto.JobTypeEntry {
	out := make([]dto.JobTypeEntry, len(list))
	for i, jt := range list {
		out[i] = dto.JobTypeEntry{Name: jt.Name, Subtypes: repositories.SubtypesFromJSON(jt.Subtypes)}
	}
	return out
}

type SettingsService interface {
	Get(tenantID string) (*models.Settings, error)
	Update(tenantID string, s *models.Settings) (*models.Settings, error)
}

type settingsService struct {
	repo repositories.SettingsRepository
}

func NewSettingsService(repo repositories.SettingsRepository) SettingsService {
	return &settingsService{repo}
}

func (s *settingsService) Get(tenantID string) (*models.Settings, error) {
	settings, err := s.repo.Get(tenantID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &models.Settings{EstPrefix: "EST", InvPrefix: "INV", Currency: "PKR", ClientCompanyID: tenantID}, nil
	}
	return settings, err
}

func (s *settingsService) Update(tenantID string, settings *models.Settings) (*models.Settings, error) {
	if err := s.repo.Save(tenantID, settings); err != nil {
		return nil, err
	}
	return s.repo.Get(tenantID)
}
