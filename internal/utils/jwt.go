package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/Hettank/habit-tracker/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`

	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey []byte
	accessTTL time.Duration
}

func NewJWTManager(
	secret string,
	accessTTL time.Duration,
) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secret),
		accessTTL: accessTTL,
	}
}

func (j *JWTManager) GenerateAccessToken(
	user *models.User,
) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,

		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTTL)),
			Subject:   strconv.FormatInt(user.ID, 10),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		claims,
	)

	return token.SignedString(j.secretKey)
}

func (j *JWTManager) ValidateAccessToken(
	tokenString string,
) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return j.secretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
