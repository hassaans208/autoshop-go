package services

import (
	"autoworkshop/dto"
	"crypto/subtle"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AdminAuthService interface {
	Login(email, password string) (string, dto.UserResponse, error)
}

type adminAuthService struct {
	email     string
	password  string
	jwtSecret string
}

func NewAdminAuthService(email, password, jwtSecret string) AdminAuthService {
	return &adminAuthService{email, password, jwtSecret}
}

func (s *adminAuthService) Login(email, password string) (string, dto.UserResponse, error) {
	if s.email == "" || s.password == "" {
		return "", dto.UserResponse{}, errors.New("admin login is not configured")
	}
	if subtle.ConstantTimeCompare([]byte(email), []byte(s.email)) != 1 ||
		subtle.ConstantTimeCompare([]byte(password), []byte(s.password)) != 1 {
		return "", dto.UserResponse{}, errors.New("invalid credentials")
	}
	token, err := s.sign(email)
	if err != nil {
		return "", dto.UserResponse{}, err
	}
	return token, dto.UserResponse{
		ID:     "admin",
		Email:  email,
		Name:   "Platform Admin",
		Role:   "admin",
		Active: true,
	}, nil
}

func (s *adminAuthService) sign(email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   "admin",
		"email": email,
		"role":  "admin",
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtSecret))
}
