package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AutologinRequest — username is staff email; passwordHash is either a bcrypt hash or plain password.
type AutologinRequest struct {
	Email        string `json:"email" binding:"required"`
	PasswordHash string `json:"passwordHash" binding:"required"`
}

type UserResponse struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Phone  string `json:"phone"`
	Active bool   `json:"active"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type AutologinLinkRequest struct {
	StaffID string `json:"staffId" binding:"required"`
}

type AutologinLinkResponse struct {
	URL   string `json:"url"`
	Email string `json:"email"`
}
