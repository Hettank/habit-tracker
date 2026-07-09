package dto

import "time"

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type AuthResult struct {
	AccessToken   string
	RefreshToken  string
	RefreshExpiry time.Time
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}
