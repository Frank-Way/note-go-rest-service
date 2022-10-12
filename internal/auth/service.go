package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"time"
)

var _ Service = &service{}

var secret = []byte("smboniudrou5wghius")

type Service interface {
	GenerateToken(ctx context.Context, login string) (string, error)
	ParseToken(ctx context.Context, token string) (string, error)
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserLogin string `json:"user_login"`
}

type service struct {
	logger *logrus.Logger
}

func NewAuthService(logger *logrus.Logger) Service {
	return &service{
		logger: logger,
	}
}

func (s service) GenerateToken(ctx context.Context, login string) (string, error) {
	claims := UserClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		login,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)
	return ss, err
}

func (s service) ParseToken(ctx context.Context, tokenStr string) (string, error) {
	if tokenStr == "" {
		return "", fmt.Errorf("empty token string")
	}
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	)
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims.UserLogin, nil
	} else {
		return "", err
	}
}
