package services

import (
	"autoworkshop/models"
	"autoworkshop/repositories"
	"crypto/subtle"
	"errors"
	"autoworkshop/utils"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(email, password string) (string, *models.Staff, error)
	Autologin(email, passwordOrHash string) (string, *models.Staff, error)
	AutologinLink(staffID string) (string, *models.Staff, error)
	Me(tenantID, staffID string) (*models.Staff, error)
}

type authService struct {
	repo        repositories.StaffRepository
	jwtSecret   string
	frontendURL string
}

func NewAuthService(repo repositories.StaffRepository, secret, frontendURL string) AuthService {
	return &authService{
		repo:        repo,
		jwtSecret:   secret,
		frontendURL: strings.TrimRight(strings.TrimSpace(frontendURL), "/"),
	}
}

func (s *authService) Login(email, password string) (string, *models.Staff, error) {
	staff, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid credentials")
		}
		return "", nil, err
	}
	if !staff.Active {
		return "", nil, errors.New("account is disabled")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	if strings.TrimSpace(staff.ClientCompanyID) == "" {
		return "", nil, errors.New("account is not linked to a workshop; ask an admin to assign your tenant")
	}
	token, err := s.sign(staff)
	return token, staff, err
}

func (s *authService) Autologin(email, passwordOrHash string) (string, *models.Staff, error) {
	staff, err := s.repo.FindByEmail(strings.TrimSpace(email))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid credentials")
		}
		return "", nil, err
	}
	if !staff.Active {
		return "", nil, errors.New("invalid credentials")
	}
	if !autologinSecretOK(staff.PasswordHash, passwordOrHash) {
		return "", nil, errors.New("invalid credentials")
	}
	if strings.TrimSpace(staff.ClientCompanyID) == "" {
		return "", nil, errors.New("invalid credentials")
	}
	token, err := s.sign(staff)
	return token, staff, err
}

// autologinSecretOK accepts either the stored bcrypt hash (link embeds hash) or a plain password.
func autologinSecretOK(storedHash, provided string) bool {
	provided = strings.TrimSpace(provided)
	if provided == "" {
		return false
	}
	if strings.HasPrefix(provided, "$2a$") || strings.HasPrefix(provided, "$2b$") || strings.HasPrefix(provided, "$2y$") {
		return subtle.ConstantTimeCompare([]byte(storedHash), []byte(provided)) == 1
	}
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(provided)) == nil
}

func (s *authService) AutologinLink(staffID string) (string, *models.Staff, error) {
	if s.frontendURL == "" {
		return "", nil, errors.New("FRONTEND_URL is not configured")
	}
	staff, err := s.repo.FindByIDAny(staffID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("staff not found")
		}
		return "", nil, err
	}
	if !staff.Active {
		return "", nil, errors.New("staff not found")
	}
	if strings.TrimSpace(staff.PasswordHash) == "" {
		return "", nil, errors.New("staff has no password")
	}
	link := utils.BuildAutologinURL(s.frontendURL, staff.Email, staff.PasswordHash)
	if link == "" {
		return "", nil, errors.New("FRONTEND_URL is not configured")
	}
	return link, staff, nil
}

func (s *authService) Me(tenantID, staffID string) (*models.Staff, error) {
	return s.repo.FindByID(tenantID, staffID)
}

func (s *authService) sign(staff *models.Staff) (string, error) {
	claims := jwt.MapClaims{
		"sub":              staff.ID,
		"email":            staff.Email,
		"role":             staff.Role,
		"clientCompanyId":  staff.ClientCompanyID,
		"exp":              time.Now().Add(24 * time.Hour).Unix(),
		"iat":              time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtSecret))
}
