package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"online-quiz/internal/domain"
)

type JWTService interface {
	GenerateToken(user *domain.User) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type JWTCustomClaims struct {
	UserID   uint        `json:"user_id"`
	Username string      `json:"username"`
	Role     domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey []byte
	issuer    string
}

func NewJWTService(secret string) JWTService {
	return &jwtService{
		secretKey: []byte(secret),
		issuer:    "online-quiz",
	}
}

func (j *jwtService) GenerateToken(user *domain.User) (string, error) {
	claims := &JWTCustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(encodedToken, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signature method")
		}
		return j.secretKey, nil
	})
}
