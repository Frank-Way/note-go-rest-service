package auth

import (
	"context"
	"github.com/Frank-Way/note-go-rest-service/user_service/internal/user"
)

type Service interface {
	GenerateToken(ctx context.Context, dto user.AuthUserDTO) (string, error)
}
