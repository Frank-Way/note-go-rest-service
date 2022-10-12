package auth

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

type Middleware struct {
	authSrv Service
	logger  *logrus.Logger
}

func NewMiddleware(authSrv Service, logger *logrus.Logger) *Middleware {
	return &Middleware{
		authSrv: authSrv,
		logger:  logger,
	}
}

func (m *Middleware) CheckAndParse(ctx context.Context, authStr string) (string, error) {
	m.logger.Debug("check if authStr is empty")
	if authStr == "" {
		m.logger.Debug("authStr is empty")
		return "", fmt.Errorf("authStr string is empty")
	}
	m.logger.Debug("check authStr format")
	authParts := strings.Split(authStr, " ")
	if len(authParts) != 2 || authParts[0] != "Bearer" {
		m.logger.Debug("wrong authStr format")
		return "", fmt.Errorf("wrong authStr string format")
	}
	m.logger.Debug("parse auth token")
	login, err := m.authSrv.ParseToken(ctx, authStr)
	if err != nil {
		m.logger.Debugf("error during token parsing: %v", err)
		return "", err
	}
	return login, nil
}

func (m *Middleware) GetToken(ctx context.Context, login string) (string, error) {
	return m.authSrv.GenerateToken(ctx, login)
}
