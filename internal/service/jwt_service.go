package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/medhansh-32/api-gateway/internal/models"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService interface{
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type JWTServiceImpl struct{
	jwtSecret string
}

func NewJWTService (secret string) (JWTService){
	return &JWTServiceImpl{jwtSecret: secret}
}

func (J JWTServiceImpl) GenerateToken(user *models.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID,10),
			Issuer:    "API-Gateway",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(J.jwtSecret))
}


func (J JWTServiceImpl) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(J.jwtSecret), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, errors.New("token has expired")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, errors.New("invalid token signature")
		default:
			return nil, err
		}
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}